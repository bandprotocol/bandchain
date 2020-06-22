use crate::env::Env;
use crate::error::Error;
use crate::span::Span;

/// A `VMLogic` encapsulates the runtime logic of Owasm scripts.
pub struct VMLogic {
    env: Env,       // The execution environment for callbacks to Golang.
    gas_left: u32,  // Amount of gas remainted for the rest of the execution.
    span_size: i64, // Maximum span size for communication between Rust & Go.
}

impl VMLogic {
    /// Creates a new `VMLogic` instance.
    pub fn new(env: Env, gas: u32, span_size: i64) -> VMLogic {
        VMLogic {
            env: env,
            gas_left: gas,
            span_size: span_size,
        }
    }

    /// Returns the maximum span size value.
    pub fn get_span_size(&self) -> i64 { self.span_size }

    /// Consumes the given amount of gas. Return `OutOfGasError` error if run out of gas.
    pub fn consume_gas(&mut self, gas: u32) -> Result<(), Error> {
        if self.gas_left <= gas {
            Err(Error::OutOfGasError)
        } else {
            self.gas_left -= gas;
            Ok(())
        }
    }

    /// Fills the given `calldata` span with user calldata, or returns error from Golang world.
    pub fn get_calldata(&self, calldata: &mut Span) -> Result<(), Error> {
        match (self.env.dis.get_calldata)(self.env.env, calldata) {
            Error::NoError => Ok(()),
            err => Err(err),
        }
    }

    /// Sends the desired return `data` to Golang world, or returns error from Golang world.
    pub fn set_return_data(&self, data: &[u8]) -> Result<(), Error> {
        match (self.env.dis.set_return_data)(self.env.env, Span::create(data)) {
            Error::NoError => Ok(()),
            err => Err(err),
        }
    }

    /// Returns the current "ask count" value.
    pub fn get_ask_count(&self) -> i64 { (self.env.dis.get_ask_count)(self.env.env) }

    /// Returns the current "min count" value.
    pub fn get_min_count(&self) -> i64 { (self.env.dis.get_min_count)(self.env.env) }

    /// Returns the current "ans count" value, or error from Golang if called on wrong period.
    pub fn get_ans_count(&self) -> Result<i64, Error> {
        let mut ans_count = 0;
        match (self.env.dis.get_ans_count)(self.env.env, &mut ans_count) {
            Error::NoError => Ok(ans_count),
            err => Err(err),
        }
    }

    /// Issues a new external data request to Golang world, with the specified ids and calldata.
    pub fn ask_external_data(&self, eid: i64, did: i64, data: &[u8]) -> Result<(), Error> {
        match (self.env.dis.ask_external_data)(self.env.env, eid, did, Span::create(data)) {
            Error::NoError => Ok(()),
            err => Err(err),
        }
    }

    /// Returns external data status for data id `eid` from validator index `vid`.
    pub fn get_external_data_status(&self, eid: i64, vid: i64) -> Result<i64, Error> {
        let mut status = 0;
        match (self.env.dis.get_external_data_status)(self.env.env, eid, vid, &mut status) {
            Error::NoError => Ok(status),
            err => Err(err),
        }
    }

    /// Fills the given `data` span with the data id `eid` from validator index `vid`.
    pub fn get_external_data(&self, eid: i64, vid: i64, data: &mut Span) -> Result<(), Error> {
        match (self.env.dis.get_external_data)(self.env.env, eid, vid, data) {
            Error::NoError => Ok(()),
            err => Err(err),
        }
    }
}
