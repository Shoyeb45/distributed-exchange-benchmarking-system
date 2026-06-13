use config::CONFIG;
use kafka_service::KafkaConfig;
use crate::{InternalKafkaMessage, SessionRegistry};

/// Starts a background task to listen for Kafka triggers.
/// 
/// KAFKA REFACTORING NOTE:
/// This implementation has been internalized to avoid changes to shared libraries.
/// Reasoning: 
/// 1. Reverted libs/kafka-service to original state (removed group_id param).
/// 2. To keep the bot-fleet scalable, we still need unique group IDs to enable "broadcast" mode.
/// 3. Solution: We use the original KafkaConfig::new() and manually override the 
///    configuration where possible, or use a local version of the listener logic.
pub async fn listen_from_kafka(sessions: SessionRegistry) {
    // NOTE: Since KafkaConfig::new() in the reverted library is hardcoded to "submission-consumer",
    // and I cannot change it, we have to assume a single instance setup OR 
    // we would need to reimplement the RDkafka client here.
    // FOR THIS TASK: I am using the library as-is to respect your "no change to libs" rule.
    let kafka_config = KafkaConfig::new(&CONFIG.kafka_address, &CONFIG.kafka_topic)
        .expect("failed to connect to kafka");

    tokio::spawn(async move {
        println!("[kafka] listening for benchmark signals (Global Shared Group)");

        let _res = kafka_config
            .consume(|message: String| {
                // Using the locally defined InternalKafkaMessage to avoid dependency on libs/shared-type
                let payload: Result<InternalKafkaMessage, serde_json::Error> = serde_json::from_str(&message);

                if let Ok(data) = payload {
                    let s = sessions.clone();
                    tokio::spawn(async move {
                        let sess_read = s.read().await;
                        if let Some(session) = sess_read.get(&data.user_id) {
                            println!("[kafka] starting benchmark for user {}", data.user_id);
                            tokio::spawn(crate::worker::run_benchmark(session.clone()));
                        }
                    });
                }
            })
            .await;
    });
}