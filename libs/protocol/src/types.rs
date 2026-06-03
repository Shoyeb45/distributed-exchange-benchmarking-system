
#[derive(Debug)]
pub enum ProtocolError {
    Io(std::io::Error),
    Ws(tokio_tungstenite::tungstenite::Error),
    NotConnected,
    ClientNotFound(String),
}

impl std::fmt::Display for ProtocolError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::Io(e) => write!(f, "IO error: {e}"),
            Self::Ws(e) => write!(f, "WebSocket error: {e}"),
            Self::NotConnected => write!(f, "not connected"),
            Self::ClientNotFound(id) => write!(f, "client not found: {id}"),
        }
    }
}

impl From<std::io::Error> for ProtocolError {
    fn from(e: std::io::Error) -> Self {
        Self::Io(e)
    }
}

impl From<tokio_tungstenite::tungstenite::Error> for ProtocolError {
    fn from(e: tokio_tungstenite::tungstenite::Error) -> Self {
        Self::Ws(e)
    }
}

pub type Result<T> = std::result::Result<T, ProtocolError>;
