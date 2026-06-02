mod listener;

#[tokio::main]
async fn main() {
    listener::listen_from_kafka().await;
}
