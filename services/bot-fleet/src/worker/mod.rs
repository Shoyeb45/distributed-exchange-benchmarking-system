use std::sync::Arc;
use tokio::sync::Mutex;
use crate::validator::Validator;
use crate::validator::types::{Command, Side};
use crate::validator::wire_protocol;
use ::protocol::server::ServerHandle;

pub struct Session {
    pub validator: Validator,
    pub server_handle: ServerHandle,
    pub client_id: String,
}

impl Session {
    pub fn new(client_id: String, server_handle: ServerHandle, symbols: Vec<String>) -> Self {
        Self {
            validator: Validator::new(symbols),
            server_handle,
            client_id,
        }
    }

    pub async fn handle_message(&mut self, msg: &str) {
        if let Some(resp) = wire_protocol::parse_response(msg) {
            self.validator.on_response_received(resp);
        } else {
            // Handle handshake messages (HELLO, SYMBOLS) or invalid format
            if msg.starts_with("HELLO") {
                // Phase 1 Handshake
                let _ = self.server_handle.send_to(&self.client_id, "HELLO version=1 platform=IICPC").await;
            } else if msg.starts_with("SYMBOLS") {
                // Phase 2 Symbol Negotiation
                // SYMBOLS_ACK AAPL GOOGL ...
                let ack = "SYMBOLS_ACK AAPL GOOGL MSFT TSLA AMZN NFLX NVDA META RELIANCE INFY TCS WIPRO HDFC ICICIBANK TATA BAJAJ HIND MARUTI SBI ONGC BPCL COALINDIA";
                let _ = self.server_handle.send_to(&self.client_id, ack).await;
            }
        }
    }

    pub async fn send_command(&mut self, cmd: Command) {
        let msg = wire_protocol::serialize_command(&cmd);
        self.validator.on_command_sent(&cmd);
        let _ = self.server_handle.send_to(&self.client_id, msg).await;
    }

    pub fn report_violations(&self) {
        for v in &self.validator.violations {
            println!("[violation] Order {}: {:?} - {}", v.order_id, v.code, v.details);
        }
    }
}

pub async fn run_benchmark(session: Arc<Mutex<Session>>) {
    // Phase 3: Start Benchmark
    {
        let s = session.lock().await;
        let start_msg = "BENCH_START run_id=test-run duration_secs=10 max_order_rate=100";
        let _ = s.server_handle.send_to(&s.client_id, start_msg).await;
    }

    // Hot Path Loop (Simulation)
    for i in 1..=10 {
        tokio::time::sleep(std::time::Duration::from_millis(100)).await;
        let mut s = session.lock().await;
        
        // Send a NEW order
        let cmd = Command::New {
            order_id: i,
            symbol: "AAPL".to_string(),
            side: if i % 2 == 0 { Side::Sell } else { Side::Buy },
            price: 1000000 + (i * 1000),
            qty: 10,
        };
        s.send_command(cmd).await;
    }

    tokio::time::sleep(std::time::Duration::from_secs(2)).await;
    
    let s = session.lock().await;
    println!("Benchmark finished for {}. Violations: {}", s.client_id, s.validator.violations.len());
    s.report_violations();
}
