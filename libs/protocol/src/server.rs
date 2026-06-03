use futures_util::{SinkExt, StreamExt};
use std::{collections::HashMap, sync::Arc};
use tokio::{
    net::TcpListener,
    sync::{RwLock, mpsc},
};
use crate::types::{ProtocolError, Result};
use tokio_tungstenite::{accept_async, tungstenite::Message};

type ClientSender = mpsc::UnboundedSender<Message>;
type ClientMap    = Arc<RwLock<HashMap<String, ClientSender>>>;

type ConnectCb    = Arc<dyn Fn(String) + Send + Sync + 'static>;
type MessageCb    = Arc<dyn Fn(String, String) + Send + Sync + 'static>;
type DisconnectCb = Arc<dyn Fn(String) + Send + Sync + 'static>;

pub struct Server {
    on_connect:    Option<ConnectCb>,
    on_message:    Option<MessageCb>,
    on_disconnect: Option<DisconnectCb>,
}

#[derive(Clone)]
pub struct ServerHandle {
    clients: ClientMap,
}

impl Server {
    pub fn new() -> Self {
        Self {
            on_connect:    None,
            on_message:    None,
            on_disconnect: None,
        }
    }

    /// Called once when a new client connects. Receives the client's unique ID.
    pub fn on_connect<F>(&mut self, f: F)
    where
        F: Fn(String) + Send + Sync + 'static,
    {
        self.on_connect = Some(Arc::new(f));
    }

    /// Called each time a text message arrives. Receives (client_id, message_text).
    pub fn on_message<F>(&mut self, f: F)
    where
        F: Fn(String, String) + Send + Sync + 'static,
    {
        self.on_message = Some(Arc::new(f));
    }

    /// Called when a client disconnects (cleanly or not).
    pub fn on_disconnect<F>(&mut self, f: F)
    where
        F: Fn(String) + Send + Sync + 'static,
    {
        self.on_disconnect = Some(Arc::new(f));
    }

    pub async fn start(self, port: &str) -> Result<ServerHandle> {
        let addr = format!("0.0.0.0:{port}");
        let listener = TcpListener::bind(&addr).await?;
        println!("[protocol] server listening on {addr}");

        let clients: ClientMap = Arc::new(RwLock::new(HashMap::new()));
        let handle = ServerHandle { clients: clients.clone() };

        // Wrap callbacks in Arc so they can be cheaply cloned per-connection
        let on_connect    = self.on_connect;
        let on_message    = self.on_message;
        let on_disconnect = self.on_disconnect;

        tokio::spawn(async move {
            while let Ok((stream, addr)) = listener.accept().await {
                let id = uuid::Uuid::new_v4().to_string();
                let clients       = clients.clone();
                let on_connect    = on_connect.clone();
                let on_message    = on_message.clone();
                let on_disconnect = on_disconnect.clone();

                tokio::spawn(async move {
                    let ws = match accept_async(stream).await {
                        Ok(ws) => ws,
                        Err(e) => {
                            eprintln!("[protocol] handshake failed for {addr}: {e}");
                            return;
                        }
                    };

                    let (mut ws_write, mut ws_read) = ws.split();

                    // Per-client outgoing channel so any task can send without locking
                    let (tx, mut rx) = mpsc::unbounded_channel::<Message>();
                    clients.write().await.insert(id.clone(), tx);

                    // Drain outgoing channel → socket
                    let write_id = id.clone();
                    tokio::spawn(async move {
                        while let Some(msg) = rx.recv().await {
                            if ws_write.send(msg).await.is_err() {
                                break;
                            }
                        }
                        let _ = ws_write.close().await;
                        println!("[protocol] write task ended for {write_id}");
                    });

                    // Fire on_connect
                    if let Some(cb) = &on_connect {
                        cb(id.clone());
                    }

                    // Read loop
                    loop {
                        match ws_read.next().await {
                            Some(Ok(Message::Text(text))) => {
                                if let Some(cb) = &on_message {
                                    cb(id.clone(), text.to_string());
                                }
                            }
                            Some(Ok(Message::Binary(bytes))) => {
                                if let Ok(text) = String::from_utf8(bytes.to_vec()) {
                                    if let Some(cb) = &on_message {
                                        cb(id.clone(), text);
                                    }
                                }
                            }
                            Some(Ok(Message::Ping(_))) => { /* auto-pong by tungstenite */ }
                            Some(Ok(Message::Close(_))) | None => break,
                            Some(Err(e)) => {
                                eprintln!("[protocol] error from {id}: {e}");
                                break;
                            }
                            _ => {}
                        }
                    }

                    // Cleanup
                    clients.write().await.remove(&id);
                    if let Some(cb) = &on_disconnect {
                        cb(id.clone());
                    }
                });
            }
        });

        Ok(handle)
    }
}

impl Default for Server {
    fn default() -> Self { Self::new() }
}

impl ServerHandle {
    /// Send a text message to one specific client by ID.
    pub async fn send_to(&self, client_id: &str, text: impl Into<String>) -> Result<()> {
        let clients = self.clients.read().await;
        match clients.get(client_id) {
            Some(tx) => {
                tx.send(Message::Text(text.into().into()))
                    .map_err(|_| ProtocolError::ClientNotFound(client_id.to_string()))?;
                Ok(())
            }
            None => Err(ProtocolError::ClientNotFound(client_id.to_string())),
        }
    }

    /// Send a text message to every connected client.
    pub async fn broadcast(&self, text: impl Into<String>) -> Result<()> {
        let msg = Message::Text(text.into().into());
        let clients = self.clients.read().await;
        for tx in clients.values() {
            let _ = tx.send(msg.clone());
        }
        Ok(())
    }

    /// Send raw bytes to one specific client.
    pub async fn send_bytes_to(&self, client_id: &str, bytes: Vec<u8>) -> Result<()> {
        let clients = self.clients.read().await;
        match clients.get(client_id) {
            Some(tx) => {
                tx.send(Message::Binary(bytes.into()))
                    .map_err(|_| ProtocolError::ClientNotFound(client_id.to_string()))?;
                Ok(())
            }
            None => Err(ProtocolError::ClientNotFound(client_id.to_string())),
        }
    }

    /// All currently connected client IDs.
    pub async fn connected_clients(&self) -> Vec<String> {
        self.clients.read().await.keys().cloned().collect()
    }

    /// Force-disconnect a specific client.
    pub async fn kick(&self, client_id: &str) -> Result<()> {
        let clients = self.clients.read().await;
        match clients.get(client_id) {
            Some(tx) => { let _ = tx.send(Message::Close(None)); Ok(()) }
            None => Err(ProtocolError::ClientNotFound(client_id.to_string())),
        }
    }
}