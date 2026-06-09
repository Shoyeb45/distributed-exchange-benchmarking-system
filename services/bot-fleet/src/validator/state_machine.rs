use std::collections::{HashMap, HashSet};
use crate::validator::types::{Command, Response, Violation, ViolationCode, ErrorReason};
use crate::validator::clob::Clob;

pub struct Validator {
    pub clobs: HashMap<String, Clob>,
    pub pending_orders: HashMap<u64, Response>, // order_id -> expected_response
    pub resting_orders: HashSet<u64>,
    pub filled_orders: HashSet<u64>,
    pub cancelled_orders: HashSet<u64>,
    pub violations: Vec<Violation>,
}

impl Validator {
    pub fn new(symbols: Vec<String>) -> Self {
        let mut clobs = HashMap::new();
        for sym in symbols {
            clobs.insert(sym.clone(), Clob::new(sym));
        }
        Self {
            clobs,
            pending_orders: HashMap::new(),
            resting_orders: HashSet::new(),
            filled_orders: HashSet::new(),
            cancelled_orders: HashSet::new(),
            violations: Vec::new(),
        }
    }

    pub fn on_command_sent(&mut self, cmd: &Command) {
        match cmd {
            Command::New { order_id, symbol, side, price, qty } => {
                if let Some(clob) = self.clobs.get_mut(symbol) {
                    let expected = clob.apply_new(*order_id, *side, *price, *qty);
                    self.pending_orders.insert(*order_id, expected);
                }
            }
            Command::Mkt { order_id, symbol, side, qty } => {
                if let Some(clob) = self.clobs.get_mut(symbol) {
                    let expected = clob.apply_mkt(*order_id, *side, *qty);
                    self.pending_orders.insert(*order_id, expected);
                }
            }
            Command::Cxl { order_id } => {
                // Find which symbol this order belongs to
                let mut expected = None;
                for clob in self.clobs.values_mut() {
                    if clob.order_lookup.contains_key(order_id) {
                        expected = Some(clob.apply_cxl(*order_id));
                        break;
                    }
                }
                let expected = expected.unwrap_or(Response::Err { 
                    order_id: *order_id, 
                    reason: ErrorReason::UnknownOrder 
                });
                self.pending_orders.insert(*order_id, expected);
            }
            Command::Bbo { symbol } => {
                if let Some(clob) = self.clobs.get(symbol) {
                    let (bid, ask) = clob.get_bbo();
                    let _expected = Response::BboResp {
                        symbol: symbol.clone(),
                        best_bid: bid,
                        best_ask: ask,
                    };
                    // BBO doesn't use order_id in wire protocol, but we'll use a dummy ID or symbol-based key
                    // For now, BBO validation might need a different path since it's a query
                }
            }
        }
    }

    pub fn on_response_received(&mut self, actual: Response) {
        let order_id = match &actual {
            Response::Ack { order_id, .. } => *order_id,
            Response::CxlAck { order_id } => *order_id,
            Response::Err { order_id, .. } => *order_id,
            Response::BboResp { .. } => return, // Handle BBO separately
        };

        if let Some(expected) = self.pending_orders.remove(&order_id) {
            self.compare_responses(order_id, expected, actual);
        } else {
            // Unexpected response (duplicate or unsolicited)
        }
    }

    fn compare_responses(&mut self, order_id: u64, expected: Response, actual: Response) {
        match (expected, actual) {
            (Response::Ack { fills: exp_fills, .. }, Response::Ack { fills: act_fills, .. }) => {
                if exp_fills.len() != act_fills.len() {
                    self.add_violation(order_id, ViolationCode::MissedFill, format!("Expected {} fills, got {}", exp_fills.len(), act_fills.len()));
                    return;
                }
                for (i, (exp, act)) in exp_fills.iter().zip(act_fills.iter()).enumerate() {
                    if exp.maker_id != act.maker_id {
                        self.add_violation(order_id, ViolationCode::WrongMaker, format!("Fill {}: Expected maker {}, got {}", i, exp.maker_id, act.maker_id));
                    }
                    if exp.price != act.price {
                        self.add_violation(order_id, ViolationCode::WrongPrice, format!("Fill {}: Expected price {}, got {}", i, exp.price, act.price));
                    }
                    if exp.qty != act.qty {
                        self.add_violation(order_id, ViolationCode::WrongQty, format!("Fill {}: Expected qty {}, got {}", i, exp.qty, act.qty));
                    }
                }
            }
            (Response::CxlAck { .. }, Response::CxlAck { .. }) => {}
            (Response::Err { reason: exp_reason, .. }, Response::Err { reason: act_reason, .. }) => {
                if exp_reason != act_reason {
                    self.add_violation(order_id, ViolationCode::WrongCancelResponse, format!("Expected ERR {}, got {}", exp_reason, act_reason));
                }
            }
            (exp, act) => {
                self.add_violation(order_id, ViolationCode::WrongCancelResponse, format!("Response type mismatch: Expected {:?}, got {:?}", exp, act));
            }
        }
    }

    fn add_violation(&mut self, order_id: u64, code: ViolationCode, details: String) {
        self.violations.push(Violation { code, order_id, details });
    }
}
