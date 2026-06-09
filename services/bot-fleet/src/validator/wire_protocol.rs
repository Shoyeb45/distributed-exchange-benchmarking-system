use crate::validator::types::{Command, Response, Fill, ErrorReason};

pub fn serialize_command(cmd: &Command) -> String {
    match cmd {
        Command::New { order_id, symbol, side, price, qty } => {
            format!("NEW {} {} {} {} {}", order_id, symbol, side, price, qty)
        }
        Command::Mkt { order_id, symbol, side, qty } => {
            format!("MKT {} {} {} {}", order_id, symbol, side, qty)
        }
        Command::Cxl { order_id } => {
            format!("CXL {}", order_id)
        }
        Command::Bbo { symbol } => {
            format!("BBO {}", symbol)
        }
    }
}

pub fn parse_response(input: &str) -> Option<Response> {
    let parts: Vec<&str> = input.split_whitespace().collect();
    if parts.is_empty() {
        return None;
    }

    match parts[0] {
        "ACK" => {
            if parts.len() < 2 { return None; }
            let order_id = parts[1].parse().ok()?;
            let mut fills = Vec::new();
            
            let mut i = 2;
            while i + 4 < parts.len() {
                if parts[i] == "FILL" {
                    let maker_id = parts[i+1].parse().ok()?;
                    let taker_id = parts[i+2].parse().ok()?;
                    let price    = parts[i+3].parse().ok()?;
                    let qty      = parts[i+4].parse().ok()?;
                    fills.push(Fill { maker_id, taker_id, price, qty });
                    i += 5;
                } else {
                    break;
                }
            }
            Some(Response::Ack { order_id, fills })
        }
        "CXL_ACK" => {
            if parts.len() < 2 { return None; }
            let order_id = parts[1].parse().ok()?;
            Some(Response::CxlAck { order_id })
        }
        "BBO_RESP" => {
            if parts.len() < 4 { return None; }
            let symbol = parts[1].to_string();
            let best_bid = parts[2].parse().ok()?;
            let best_ask = parts[3].parse().ok()?;
            Some(Response::BboResp { symbol, best_bid, best_ask })
        }
        "ERR" => {
            if parts.len() < 3 { return None; }
            let order_id = parts[1].parse().ok()?;
            let reason = match parts[2] {
                "UNKNOWN_SYMBOL"    => ErrorReason::UnknownSymbol,
                "DUPLICATE_ID"      => ErrorReason::DuplicateId,
                "UNKNOWN_ORDER"     => ErrorReason::UnknownOrder,
                "ALREADY_FILLED"    => ErrorReason::AlreadyFilled,
                "ALREADY_CANCELLED" => ErrorReason::AlreadyCancelled,
                "SELF_TRADE"        => ErrorReason::SelfTrade,
                "BAD_FORMAT"        => ErrorReason::BadFormat,
                _ => return None,
            };
            Some(Response::Err { order_id, reason })
        }
        _ => None,
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::validator::types::Side;

    #[test]
    fn test_serialize_new() {
        let cmd = Command::New {
            order_id: 1001,
            symbol: "AAPL".to_string(),
            side: Side::Buy,
            price: 1005000,
            qty: 50,
        };
        assert_eq!(serialize_command(&cmd), "NEW 1001 AAPL BUY 1005000 50");
    }

    #[test]
    fn test_parse_ack_fills() {
        let input = "ACK 1003 FILL 1001 1003 1005000 30 FILL 1002 1003 1005500 70";
        let resp = parse_response(input).unwrap();
        if let Response::Ack { order_id, fills } = resp {
            assert_eq!(order_id, 1003);
            assert_eq!(fills.len(), 2);
            assert_eq!(fills[0].maker_id, 1001);
            assert_eq!(fills[1].price, 1005500);
        } else {
            panic!("Expected Ack");
        }
    }
}
