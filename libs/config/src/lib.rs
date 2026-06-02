use dotenvy::dotenv;
use std::env;
use once_cell::sync::Lazy;

#[derive(Debug, Clone)]
pub struct Config {
    pub kafka_address: String,
    pub kafka_topic: String,
    pub database_url: String,
}

impl Config {
    pub fn load() -> Self {
        dotenv().ok();
        println!("Loading environment variable");

        Self {
            kafka_address: Self::get_env(
                "KAFKA_ADDRESS".to_string(),
                "KAFKA_ADDRESS is needed".to_string(),
                "localhost:9092".to_string(),
            ),
            kafka_topic: Self::get_env(
                "KAFKA_TOPIC".to_string(),
                "".to_string(),
                "submission-created".to_string(),
            ),
            database_url: Self::get_env(
                String::from("DATABASE_URL"),
                String::from("DATABASE_URL is required"),
                String::from(""),
            ),
        }
    }

    fn get_env(key: String, message: String, default_value: String) -> String {
        let res = env::var(key);

        match res {
            Ok(value) => {
                if value == "" {
                    panic!("{}", message);
                }
                return value;
            }
            Err(_err) => {
                if default_value != "" {
                    return default_value;
                }
                panic!("{}", message);
            }
        }
    }
}


pub static CONFIG: Lazy<Config> = Lazy::new(Config::load);