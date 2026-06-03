use protocol::client::Client;

#[tokio::main]
async fn main() {
    let mut client = Client::new();

    client.on_message(|msg| {
        println!("[client] 📨 server says: {msg}");
    });

    client.on_disconnect(|| {
        println!("[client] ❌ disconnected from server");
    });

    let handle = client
        .connect("ws://127.0.0.1:5052")
        .await
        .expect("could not connect — is the server running?");

    println!("[client] ✅ connected! Type a message and press Enter. Ctrl+C to quit.");

    tokio::signal::ctrl_c()
        .await
        .expect("failed to listen for ctrl+c");

    println!("shutting down...");
}
