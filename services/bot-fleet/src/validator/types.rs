use std::fmt;

/// Represents the side of an order in the market.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum Side {
    Buy,
    Sell,
}

impl fmt::Display for Side {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self {
            Side::Buy => write!(f, "BUY"),
            Side::Sell => write!(f, "SELL"),
        }
    }
}

/// A single execution (fill) against a resting order.
#[derive(Debug, Clone, PartialEq, Eq)]
pub struct Fill {
    pub maker_id: u64,
    pub taker_id: u64,
    pub price: u64, // Fixed-point: 1 unit = $0.0001
    pub qty: u64,
}

/// Commands sent from the Platform to the Contestant.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum Command {
    New {
        order_id: u64,
        symbol: String,
        side: Side,
        price: u64,
        qty: u64,
    },
    Mkt {
        order_id: u64,
        symbol: String,
        side: Side,
        qty: u64,
    },
    Cxl {
        order_id: u64,
    },
    Bbo {
        symbol: String,
    },
}

/// Responses sent from the Contestant back to the Platform.
#[derive(Debug, Clone, PartialEq, Eq)]
pub enum Response {
    /// Order accepted, may contain zero or more fills.
    Ack {
        order_id: u64,
        fills: Vec<Fill>,
    },
    /// Cancellation successfully processed.
    CxlAck {
        order_id: u64,
    },
    /// Current Best Bid and Offer for a symbol.
    BboResp {
        symbol: String,
        best_bid: u64,
        best_ask: u64,
    },
    /// Rejection for various reasons (see ErrorReason).
    Err {
        order_id: u64,
        reason: ErrorReason,
    },
}

#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ErrorReason {
    UnknownSymbol,
    DuplicateId,
    UnknownOrder,
    AlreadyFilled,
    AlreadyCancelled,
    SelfTrade,
    BadFormat,
}

impl fmt::Display for ErrorReason {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        let s = match self {
            ErrorReason::UnknownSymbol => "UNKNOWN_SYMBOL",
            ErrorReason::DuplicateId => "DUPLICATE_ID",
            ErrorReason::UnknownOrder => "UNKNOWN_ORDER",
            ErrorReason::AlreadyFilled => "ALREADY_FILLED",
            ErrorReason::AlreadyCancelled => "ALREADY_CANCELLED",
            ErrorReason::SelfTrade => "SELF_TRADE",
            ErrorReason::BadFormat => "BAD_FORMAT",
        };
        write!(f, "{}", s)
    }
}

/// Violation codes as defined by the IICPC Correctness Specification.
/// These represent logic errors in the contestant's matching engine.
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum ViolationCode {
    FillAfterCancel,
    CrossSymbolFill,
    PhantomFill,
    PricePriorityViolation,
    TimePriorityViolation,
    WrongPrice,
    WrongMaker,
    WrongQty,
    MissedFill,
    FillOrderViolation,
    StaleBbo,
    FalseCancel,
    WrongCancelResponse,
    DuplicateAccepted,
    UnknownSymbolAccepted,
}

#[derive(Debug, Clone)]
pub struct Violation {
    pub code: ViolationCode,
    pub order_id: u64,
    pub details: String,
}
