use serde::{Deserialize, Serialize, de};
pub use serde_json::Error;
#[derive(Serialize, Deserialize, Debug, Clone, Copy)]
pub struct KafkaMessageResponse {
    user_id: i32,
    submission_id: i32,
    created_time: i64,
}

pub fn parse_json<'a, T>(str: &'a str) -> Result<T, Error>
where
    T: de::Deserialize<'a>,
{
    let parsed_value: Result<T, Error> = serde_json::from_str(str);
    parsed_value
}

