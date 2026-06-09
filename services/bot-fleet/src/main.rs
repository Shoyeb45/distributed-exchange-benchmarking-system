mod listener;
mod worker;
mod validator;

use std::collections::HashMap;
use std::sync::Arc;
use tokio::sync::Mutex;
use protocol::server::Server;
use crate::worker::{Session, ActiveSessions};

#[tokio::main]
async fn main() {
    let mut server = Server::new();
    let sessions: ActiveSessions = Arc::new(Mutex::new(HashMap::new()));

    let sessions_connect = sessions.clone();
    server.on_connect(move |id| {
        println!("[server] client connected, id = {id}");
    });

    let sessions_msg = sessions.clone();
    server.on_message(move |id, msg| {
        let sessions = sessions_msg.clone();
        let id_clone = id.clone();
        tokio::spawn(async move {
            let mut map = sessions.lock().await;
            if let Some(session) = map.get_mut(&id_clone) {
                session.handle_message(&msg).await;
            } else {
                // Potential first message (handshake)
                if msg.starts_with("HELLO") {
                    println!("[server] Handshake started for {id_clone}");
                }
            }
        });
    });

    let sessions_disconnect = sessions.clone();
    server.on_disconnect(move |id| {
        let sessions = sessions_disconnect.clone();
        tokio::spawn(async move {
            sessions.lock().await.remove(&id);
            println!("[server] client disconnected → id = {id}");
        });
    });

    let handle = server.start("5052").await.expect("failed to start");
    println!("[server] running on 5052");

    // Integration: For this implementation, we'll manually trigger a session for testing
    // In a real system, this would be triggered by a Kafka message or API call
    let handle_clone = handle.clone();
    tokio::spawn(async move {
        loop {
            tokio::time::sleep(std::time::Duration::from_secs(1)).await;
            let clients = handle_clone.connected_clients().await;
            for id in clients {
                let mut map = sessions.lock().await;
                if !map.contains_key(&id) {
                    println!("[server] Initializing session for {id}");
                    map.insert(id.clone(), Session::new(id.clone(), handle_clone.clone(), vec!["AAPL".to_string()]));
                    
                    // Trigger benchmark simulation
                    let session_for_run = Arc::new(Mutex::new(Session::new(id.clone(), handle_clone.clone(), vec!["AAPL".to_string()])));
                    tokio::spawn(crate::worker::run_benchmark(session_for_run));
                }
            }
        }
    });

    listener::listen_from_kafka().await;

    tokio::signal::ctrl_c()
        .await
        .expect("failed to listen for ctrl+c");

    println!("shutting down...");
}
