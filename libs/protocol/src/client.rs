use crate::types::{ProtocolError, Result};
use futures_util::{SinkExt, StreamExt};
use std::sync::Arc;
use tokio::sync::mpsc;
use tokio_tungstenite::{connect_async, tungstenite::Message};
// type ClientMessageCallback = Arc<Mutex<Option<Box<dyn Fn(String) + Send + Sync + 'static>>>>;
type ClientMessageCallback = Arc<dyn Fn(String) + Send + Sync + 'static>;
type ClientDisconnectCallback = Arc<dyn Fn() + Send + Sync + 'static>;

/// Build a client, register event callbacks, then call `.connect()`.
pub struct Client {
    on_message: Option<ClientMessageCallback>,
    on_disconnect: Option<ClientDisconnectCallback>,
}

impl Client {
    pub fn new() -> Self {
        Self {
            on_message: None,
            on_disconnect: None,
        }
    }

    /// Called each time a text message arrives from the server.
    pub fn on_message<F>(&mut self, f: F)
    where
        F: Fn(String) + Send + Sync + 'static,
    {
        self.on_message = Some(Arc::new(f));
    }

    /// Called when the connection closes (cleanly or by error).
    pub fn on_disconnect<F>(&mut self, f: F)
    where
        F: Fn() + Send + Sync + 'static,
    {
        self.on_disconnect = Some(Arc::new(f));
    }

    /// Connect to a WebSocket URL, e.g. `"ws://127.0.0.1:9001"`.
    /// Returns a `ClientHandle` you can use to send messages.
    pub async fn connect(self, url: &str) -> Result<ClientHandle> {
        let (ws_stream, _response) = connect_async(url).await?;
        println!("[protocol] connected to {url}");

        let (mut ws_write, mut ws_read) = ws_stream.split();

        let (tx, mut rx) = mpsc::unbounded_channel::<Message>();

        tokio::spawn(async move {
            while let Some(msg) = rx.recv().await {
                if ws_write.send(msg).await.is_err() {
                    break;
                }
            }
            let _ = ws_write.close().await;
        });

        let on_message = self.on_message.clone();
        let on_disconnect = self.on_disconnect.clone();

        // Read loop
        tokio::spawn(async move {
            loop {
                match ws_read.next().await {
                    Some(Ok(Message::Text(text))) => {
                        if let Some(cb) = &on_message {
                            cb(text.to_string());
                        }
                    }
                    Some(Ok(Message::Binary(bytes))) => {
                        if let Ok(text) = String::from_utf8(bytes.to_vec()) {
                            if let Some(cb) = &on_message {
                                cb(text);
                            }
                        }
                    }
                    Some(Ok(Message::Ping(_))) => { /* auto-pong */ }
                    Some(Ok(Message::Close(_))) | None => break,
                    Some(Err(e)) => {
                        eprintln!("[protocol] connection error: {e}");
                        break;
                    }
                    _ => {}
                }
            }

            if let Some(cb) = &on_disconnect {
                cb();
            }
            println!("[protocol] disconnected");
        });

        Ok(ClientHandle { tx })
    }
}

impl Default for Client {
    fn default() -> Self {
        Self::new()
    }
}

#[derive(Clone)]
pub struct ClientHandle {
    tx: mpsc::UnboundedSender<Message>,
}

impl ClientHandle {
    /// Send a text message to the server.
    pub fn send(&self, text: impl Into<String>) -> Result<()> {
        self.tx
            .send(Message::Text(text.into().into()))
            .map_err(|_| ProtocolError::NotConnected)
    }

    /// Send raw bytes to the server.
    pub fn send_bytes(&self, bytes: Vec<u8>) -> Result<()> {
        self.tx
            .send(Message::Binary(bytes.into()))
            .map_err(|_| ProtocolError::NotConnected)
    }

    /// Send a ping to the server.
    pub fn ping(&self) -> Result<()> {
        self.tx
            .send(Message::Ping(vec![].into()))
            .map_err(|_| ProtocolError::NotConnected)
    }

    /// Gracefully close the connection.
    pub fn disconnect(&self) -> Result<()> {
        self.tx
            .send(Message::Close(None))
            .map_err(|_| ProtocolError::NotConnected)
    }
}
