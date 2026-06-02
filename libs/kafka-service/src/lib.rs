use futures::StreamExt;
use rdkafka::{
    ClientConfig, Message,
    consumer::{Consumer, StreamConsumer},
};
use std::error::Error;

pub struct KafkaConfig {
    pub consumer: StreamConsumer,
}

impl KafkaConfig {
    pub fn new(kafka_address: &str, kafka_topic: &str) -> Result<KafkaConfig, Box<dyn Error>> {
        let consumer: StreamConsumer = ClientConfig::new()
            .set("group.id", "submission-consumer")
            .set("bootstrap.servers", kafka_address)
            .set("auto.offset.reset", "earliest")
            .create()?;

        consumer.subscribe(&[kafka_topic])?;

        Ok(KafkaConfig { consumer })
    }

    pub async fn consume<F>(&self, mut handler: F) -> Result<(), Box<dyn Error>>
    where
        F: FnMut(String),
    {
        let mut stream = self.consumer.stream();

        while let Some(message) = stream.next().await {
            match message {
                Ok(msg) => {
                    let res = msg.payload();

                    match res {
                        Some(value) => {
                            let decoded = std::str::from_utf8(value);

                            match decoded {
                                Ok(val) => {
                                    handler(val.to_string());
                                }
                                _ => continue,
                            }
                        }
                        _ => continue,
                    }
                }
                Err(e) => {
                    eprintln!("Kafka error: {}", e);
                }
            }
        }

        Ok(())
    }
}
