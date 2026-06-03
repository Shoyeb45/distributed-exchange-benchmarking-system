use config::CONFIG;
use kafka_service::KafkaConfig;
use shared_type::{KafkaMessageResponse, parse_json};

pub async fn listen_from_kafka() {
    let kafka_config = KafkaConfig::new(&CONFIG.kafka_address, &CONFIG.kafka_topic)
        .expect("failed to connect to kafka");

    tokio::spawn(async move {
        let _res = kafka_config
            .consume(|message: String| {
                let response: Result<KafkaMessageResponse, shared_type::Error> =
                    parse_json(&message);

                match response {
                    Ok(payload) => println!("{:?}", payload),
                    Err(err) => eprintln!("parse error: {:?}", err),
                }
            })
            .await;

        // If consume() ever returns, log it so you know why
        eprintln!("[kafka] consumer exited");
    });
}