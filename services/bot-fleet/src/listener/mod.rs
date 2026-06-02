use config::CONFIG;
use kafka_service::KafkaConfig;
use serde;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Debug)]
struct KafkaMessageResponse {
    user_id: i32,
    submission_id: i32,
    created_time: i64,
}

pub async fn listen_from_kafka() {
    let kafka_config = KafkaConfig::new(&CONFIG.kafka_address, &CONFIG.kafka_topic).unwrap();

    let _res = kafka_config
        .consume(|message| {
            let response: KafkaMessageResponse = serde_json::from_str(&message).unwrap();
        })
        .await;
}
