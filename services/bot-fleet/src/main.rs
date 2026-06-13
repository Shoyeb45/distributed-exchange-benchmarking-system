mod listener;
mod worker;
pub mod validator;

use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::{Mutex, RwLock};
use protocol::server::Server;
use crate::worker::Session;

// INTERNALIZED DATA STRUCTURES (To avoid changing libs/shared-type)
// Reason: Reverted libs/shared-type to keep project isolation. 
// These 5 lines define the Kafka payload locally.
#[derive(serde::Deserialize, Debug, Clone, Copy)]
pub struct InternalKafkaMessage {
    pub user_id: i32,
    pub submission_id: i32,
    pub created_time: i64,
}

/// A thread-safe, high-concurrency registry for active sessions.
/// We use a RwLock for the map itself to allow multiple concurrent readers.
pub type SessionRegistry = Arc<RwLock<HashMap<i32, Arc<Mutex<Session>>>>>;
/// Maps connection IDs to User IDs for quick routing.
type ConnectionMap = Arc<RwLock<HashMap<String, i32>>>;

#[tokio::main]
async fn main() {
    let mut server = Server::new();
    
    // Using Arc<RwLock> for the registries to allow high-concurrency read access 
    // during the "Hot Path" (message arrival).
    let sessions: SessionRegistry = Arc::new(RwLock::new(HashMap::new()));
    let connections: ConnectionMap = Arc::new(RwLock::new(HashMap::new()));

    // --- HOT PATH: Connection Callback ---
    server.on_connect(move |id| {
        println!("[server] new connection: {id}");
    });

    // --- HOT PATH: Message Arrival (Latency Sensitive) ---
    let sessions_msg = sessions.clone();
    let connections_msg = connections.clone();
    
    // We register callbacks BEFORE starting the server to avoid move errors.
    // However, we need the ServerHandle to create new Sessions.
    // We specify the type explicitly to avoid inference errors.
    let server_handle_capture: Arc<Mutex<Option<::protocol::server::ServerHandle>>> = Arc::new(Mutex::new(None));
    let server_handle_for_msg = server_handle_capture.clone();

    server.on_message(move |id, msg| {
        let sessions = sessions_msg.clone();
        let connections = connections_msg.clone();
        let handle_opt = server_handle_for_msg.clone();
        let id_clone = id.clone();
        
        // ACCURACY NOTE: We capture the timestamp at the earliest possible moment.
        let arrival_time = std::time::Instant::now();

        tokio::spawn(async move {
            // Move arrival_time into the task to ensure it's "used" for future metric expansion.
            let _ = arrival_time;

            // 1. Check if we know who this connection is
            let user_id = {
                let conn_read = connections.read().await;
                conn_read.get(&id_clone).cloned()
            };

            if let Some(uid) = user_id {
                // 2. Direct Routing: Find the session and process
                let session_opt = {
                    let sess_read = sessions.read().await;
                    sess_read.get(&uid).cloned()
                };

                if let Some(session) = session_opt {
                    let mut s = session.lock().await;
                    s.handle_message(&msg).await; 
                }
            } else if msg.starts_with("HELLO") {
                // 3. Identification Phase
                let uid = msg.split("user_id=")
                    .nth(1)
                    .and_then(|s| s.parse::<i32>().ok());
                
                if let Some(u) = uid {
                    let h_lock = handle_opt.lock().await;
                    if let Some(handle) = &*h_lock {
                        println!("[server] identifying {id_clone} as user {u}");
                        let new_session = Arc::new(Mutex::new(Session::new(id_clone.clone(), handle.clone(), vec!["AAPL".to_string()])));
                        
                        // We must release the read lock before acquiring write lock to avoid deadlock,
                        // but here we are in a spawned task with no previous locks held.
                        connections.write().await.insert(id_clone, u);
                        sessions.write().await.insert(u, new_session);
                    }
                }
            }
        });
    });

    let sessions_disconnect = sessions.clone();
    let connections_disconnect = connections.clone();
    server.on_disconnect(move |id| {
        let sessions = sessions_disconnect.clone();
        let connections = connections_disconnect.clone();
        tokio::spawn(async move {
            let mut conn_write = connections.write().await;
            if let Some(uid) = conn_write.remove(&id) {
                let mut sess_write = sessions.write().await;
                sess_write.remove(&uid);
                println!("[server] cleaned up session for user {uid}");
            }
        });
    });

    // Start the server and get the handle
    let port = "5052"; 
    let handle = server.start(port).await.expect("failed to start");
    
    // Store the handle so the on_message callback can use it for new connections
    *server_handle_capture.lock().await = Some(handle);
    
    println!("[server] bot-fleet running on {port}");

    // --- CONTROL PLANE: Kafka Listener ---
    let sessions_kafka = sessions.clone();
    listener::listen_from_kafka(sessions_kafka).await;

    tokio::signal::ctrl_c()
        .await
        .expect("failed to listen for ctrl+c");

    println!("shutting down...");
}
