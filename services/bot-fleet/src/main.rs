mod listener;
mod worker;
use protocol::server::Server;

#[tokio::main]
async fn main() {
    // start ws server
    let mut server = Server::new();

    server.on_connect(|id| {
        println!("[server] client connected, id = {id}");
    });

    server.on_message(|id, msg| {
        println!("[server] {id} says: {msg}");
    });

    server.on_disconnect(|id| {
        println!("[server] client disconnected → id = {id}");
    });

    let _handle = server.start("5052").await.expect("failed to start");
    println!("[server] running on 5052");

    listener::listen_from_kafka().await;

    tokio::signal::ctrl_c()
        .await
        .expect("failed to listen for ctrl+c");

    println!("shutting down...");
}
