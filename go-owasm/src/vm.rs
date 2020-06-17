use crate::env::Env;
use crate::error::Error;
use crate::span::Span;

pub struct VMLogic {
    env: Env,
    gas_left: u32,
    span_size: usize,
}

impl VMLogic {
    pub fn new(env: Env, gas: u32) -> VMLogic {
        VMLogic {
            env: env,
            gas_left: gas,
            span_size: 1000, // TODO: Get span size from env
        }
    }

    pub fn get_span_size(&self) -> usize {
        self.span_size
    }

    pub fn consume_gas(&mut self, gas: u32) -> Result<(), Error> {
        if self.gas_left <= gas {
            Err(Error::GasLimitExceedError)
        } else {
            self.gas_left -= gas;
            Ok(())
        }
    }

    pub fn get_calldata(&self, calldata: &mut Span) -> Result<(), Error> {
        let err: Error = (self.env.dis.get_calldata)(self.env.env, calldata).into();
        match err {
            Error::NoError => Ok(()),
            _ => Err(err),
        }
    }

    pub fn set_return_data(&self, data: &[u8]) -> Result<(), Error> {
        let err: Error = (self.env.dis.set_return_data)(self.env.env, Span::create(data)).into();
        match err {
            Error::NoError => Ok(()),
            _ => Err(err),
        }
    }

    pub fn get_ask_count(&self) -> i64 {
        (self.env.dis.get_ask_count)(self.env.env)
    }

    pub fn get_min_count(&self) -> i64 {
        (self.env.dis.get_min_count)(self.env.env)
    }

    pub fn get_ans_count(&self) -> Result<i64, Error> {
        let mut ans_count = 0;
        let err: Error = (self.env.dis.get_ans_count)(self.env.env, &mut ans_count).into();
        match err {
            Error::NoError => Ok(ans_count),
            _ => Err(err),
        }
    }

    pub fn ask_external_data(&self, eid: i64, did: i64, data: &[u8]) -> Result<(), Error> {
        let err: Error =
            (self.env.dis.ask_external_data)(self.env.env, eid, did, Span::create(data)).into();
        match err {
            Error::NoError => Ok(()),
            _ => Err(err),
        }
    }

    pub fn get_external_data_status(&self, eid: i64, vid: i64) -> Result<i64, Error> {
        let mut status = 0;
        let err: Error =
            (self.env.dis.get_external_data_status)(self.env.env, eid, vid, &mut status).into();
        match err {
            Error::NoError => Ok(status),
            _ => Err(err),
        }
    }

    pub fn get_external_data(&self, eid: i64, vid: i64, data: &mut Span) -> Result<(), Error> {
        let err: Error = (self.env.dis.get_external_data)(self.env.env, eid, vid, data).into();
        match err {
            Error::NoError => Ok(()),
            _ => Err(err),
        }
    }
}
