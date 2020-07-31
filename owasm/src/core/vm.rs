use crate::core::error::Error;

pub trait Env {
    /// Returns the maximum span size value.
    fn get_span_size(&self) -> i64;
    /// Returns user calldata, or returns error from VM runner.
    fn get_calldata(&self) -> Result<Vec<u8>, Error>;
    /// Sends the desired return `data` to VM runner, or returns error from VM runner.
    fn set_return_data(&self, data: &[u8]) -> Result<(), Error>;
    /// Returns the current "ask count" value.
    fn get_ask_count(&self) -> i64;
    /// Returns the current "min count" value.
    fn get_min_count(&self) -> i64;
    /// Returns the current "ans count" value, or error from VM runner if called on wrong period.
    fn get_ans_count(&self) -> Result<i64, Error>;
    /// Issues a new external data request to VM runner, with the specified ids and calldata.
    fn ask_external_data(&self, eid: i64, did: i64, data: &[u8]) -> Result<(), Error>;
    /// Returns external data status for data id `eid` from validator index `vid`.
    fn get_external_data_status(&self, eid: i64, vid: i64) -> Result<i64, Error>;
    /// Returns data span with the data id `eid` from validator index `vid`.
    fn get_external_data(&self, eid: i64, vid: i64) -> Result<Vec<u8>, Error>;
}

/// A `VMLogic` encapsulates the runtime logic of Owasm scripts.
pub struct VMLogic<'a, E>
where
    E: Env,
{
    pub env: &'a E,     // The execution environment.
    pub gas_limit: u32, // Amount of gas allowed for total execution.
    pub gas_used: u32,  // Amount of gas used in this execution.
}

impl<'a, E> VMLogic<'a, E>
where
    E: Env,
{
    /// Creates a new `VMLogic` instance.
    pub fn new(env: &'a E, gas: u32) -> VMLogic<'a, E> {
        VMLogic { env: env, gas_limit: gas, gas_used: 0 }
    }

    /// Consumes the given amount of gas. Return `OutOfGasError` error if run out of gas.
    pub fn consume_gas(&mut self, gas: u32) -> Result<(), Error> {
        self.gas_used = self.gas_used.saturating_add(gas);
        if self.gas_used > self.gas_limit {
            Err(Error::OutOfGasError)
        } else {
            Ok(())
        }
    }
}
