use std::collections::{BTreeMap, VecDeque, HashMap};
use crate::validator::types::{Side, Fill, Response, ErrorReason};

/// Internal representation of an order resting in the book.
#[derive(Debug, Clone)]
pub struct InternalOrder {
    pub order_id: u64,
    pub price: u64,
    pub qty: u64,
    pub remaining_qty: u64,
    pub side: Side,
}

/// A Reference Central Limit Order Book (CLOB).
/// This serves as the 'Ground Truth' for validating contestant matching logic.
pub struct Clob {
    pub symbol: String,
    /// Bids sorted descending (highest price first).
    /// Value is a FIFO queue of orders at that price level.
    pub bids: BTreeMap<u64, VecDeque<InternalOrder>>,
    /// Asks sorted ascending (lowest price first).
    /// Value is a FIFO queue of orders at that price level.
    pub asks: BTreeMap<u64, VecDeque<InternalOrder>>,
    /// Fast lookup to find the side of a resting order by its ID.
    pub order_lookup: HashMap<u64, Side>,
}

impl Clob {
    pub fn new(symbol: String) -> Self {
        Self {
            symbol,
            bids: BTreeMap::new(),
            asks: BTreeMap::new(),
            order_lookup: HashMap::new(),
        }
    }

    /// Processes a NEW limit order.
    /// Simultaneously matches against the opposite book and appends remainders to the resting book.
    pub fn apply_new(&mut self, order_id: u64, side: Side, price: u64, qty: u64) -> Response {
        let mut fills = Vec::new();
        let mut remaining_taker_qty = qty;

        match side {
            Side::Buy => {
                // Match against asks (lowest first: BTreeMap naturally sorts ascending)
                let mut prices_to_remove = Vec::new();
                for (&ask_price, queue) in self.asks.iter_mut() {
                    if ask_price > price { break; } // Price doesn't cross
                    
                    while let Some(mut maker) = queue.pop_front() {
                        let fill_qty = std::cmp::min(remaining_taker_qty, maker.remaining_qty);
                        fills.push(Fill {
                            maker_id: maker.order_id,
                            taker_id: order_id,
                            price: ask_price, // Match occurs at Maker's limit price
                            qty: fill_qty,
                        });

                        maker.remaining_qty -= fill_qty;
                        remaining_taker_qty -= fill_qty;

                        if maker.remaining_qty > 0 {
                            // Partial fill on maker: put it back at the front of its FIFO queue
                            queue.push_front(maker);
                            break;
                        } else {
                            // Full fill on maker: remove from lookup
                            self.order_lookup.remove(&maker.order_id);
                        }

                        if remaining_taker_qty == 0 { break; }
                    }

                    if queue.is_empty() {
                        prices_to_remove.push(ask_price);
                    }
                    if remaining_taker_qty == 0 { break; }
                }

                for p in prices_to_remove {
                    self.asks.remove(&p);
                }

                if remaining_taker_qty > 0 {
                    // Taker not fully filled: rest the remainder in the Bids book
                    let order = InternalOrder { order_id, price, qty, remaining_qty: remaining_taker_qty, side };
                    self.bids.entry(price).or_default().push_back(order);
                    self.order_lookup.insert(order_id, Side::Buy);
                }
            }
            Side::Sell => {
                // Match against bids (highest first: using .rev() on BTreeMap)
                let mut prices_to_remove = Vec::new();
                for (&bid_price, queue) in self.bids.iter_mut().rev() {
                    if bid_price < price { break; }

                    while let Some(mut maker) = queue.pop_front() {
                        let fill_qty = std::cmp::min(remaining_taker_qty, maker.remaining_qty);
                        fills.push(Fill {
                            maker_id: maker.order_id,
                            taker_id: order_id,
                            price: bid_price,
                            qty: fill_qty,
                        });

                        maker.remaining_qty -= fill_qty;
                        remaining_taker_qty -= fill_qty;

                        if maker.remaining_qty > 0 {
                            queue.push_front(maker);
                            break;
                        } else {
                            self.order_lookup.remove(&maker.order_id);
                        }

                        if remaining_taker_qty == 0 { break; }
                    }

                    if queue.is_empty() {
                        prices_to_remove.push(bid_price);
                    }
                    if remaining_taker_qty == 0 { break; }
                }

                for p in prices_to_remove {
                    self.bids.remove(&p);
                }

                if remaining_taker_qty > 0 {
                    let order = InternalOrder { order_id, price, qty, remaining_qty: remaining_taker_qty, side };
                    self.asks.entry(price).or_default().push_back(order);
                    self.order_lookup.insert(order_id, Side::Sell);
                }
            }
        }

        Response::Ack { order_id, fills }
    }

    /// Processes a MARKET order.
    /// Sweeps the opposite book immediately. Remainders are cancelled as per spec.
    pub fn apply_mkt(&mut self, order_id: u64, side: Side, qty: u64) -> Response {
        let mut fills = Vec::new();
        let mut remaining_taker_qty = qty;

        match side {
            Side::Buy => {
                let mut prices_to_remove = Vec::new();
                for (&ask_price, queue) in self.asks.iter_mut() {
                    while let Some(mut maker) = queue.pop_front() {
                        let fill_qty = std::cmp::min(remaining_taker_qty, maker.remaining_qty);
                        fills.push(Fill {
                            maker_id: maker.order_id,
                            taker_id: order_id,
                            price: ask_price,
                            qty: fill_qty,
                        });

                        maker.remaining_qty -= fill_qty;
                        remaining_taker_qty -= fill_qty;

                        if maker.remaining_qty > 0 {
                            queue.push_front(maker);
                            break;
                        } else {
                            self.order_lookup.remove(&maker.order_id);
                        }

                        if remaining_taker_qty == 0 { break; }
                    }
                    if queue.is_empty() {
                        prices_to_remove.push(ask_price);
                    }
                    if remaining_taker_qty == 0 { break; }
                }
                for p in prices_to_remove { self.asks.remove(&p); }
            }
            Side::Sell => {
                let mut prices_to_remove = Vec::new();
                for (&bid_price, queue) in self.bids.iter_mut().rev() {
                    while let Some(mut maker) = queue.pop_front() {
                        let fill_qty = std::cmp::min(remaining_taker_qty, maker.remaining_qty);
                        fills.push(Fill {
                            maker_id: maker.order_id,
                            taker_id: order_id,
                            price: bid_price,
                            qty: fill_qty,
                        });

                        maker.remaining_qty -= fill_qty;
                        remaining_taker_qty -= fill_qty;

                        if maker.remaining_qty > 0 {
                            queue.push_front(maker);
                            break;
                        } else {
                            self.order_lookup.remove(&maker.order_id);
                        }

                        if remaining_taker_qty == 0 { break; }
                    }
                    if queue.is_empty() {
                        prices_to_remove.push(bid_price);
                    }
                    if remaining_taker_qty == 0 { break; }
                }
                for p in prices_to_remove { self.bids.remove(&p); }
            }
        }

        Response::Ack { order_id, fills }
    }

    /// Processes a CANCEL command.
    pub fn apply_cxl(&mut self, order_id: u64) -> Response {
        if let Some(side) = self.order_lookup.remove(&order_id) {
            match side {
                Side::Buy => {
                    for queue in self.bids.values_mut() {
                        if let Some(pos) = queue.iter().position(|o| o.order_id == order_id) {
                            queue.remove(pos);
                            return Response::CxlAck { order_id };
                        }
                    }
                }
                Side::Sell => {
                    for queue in self.asks.values_mut() {
                        if let Some(pos) = queue.iter().position(|o| o.order_id == order_id) {
                            queue.remove(pos);
                            return Response::CxlAck { order_id };
                        }
                    }
                }
            }
        }
        Response::Err { order_id, reason: ErrorReason::UnknownOrder }
    }

    /// Returns the current Best Bid and Offer.
    pub fn get_bbo(&self) -> (u64, u64) {
        let best_bid = self.bids.keys().next_back().cloned().unwrap_or(0);
        let best_ask = self.asks.keys().next().cloned().unwrap_or(0);
        (best_bid, best_ask)
    }
}
