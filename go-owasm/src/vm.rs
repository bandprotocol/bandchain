use crate::env::Env;
use crate::span::Span;

pub struct VMLogic {
    gas_left: u32,
    env: Env,
}

impl VMLogic {
    pub fn new(env: Env) -> VMLogic {
        VMLogic {
            env: env,
            gas_left: 100,
        }
    }

    pub fn consume_gas(&mut self, gas: u32) -> Result<(), ()> {
        if self.gas_left <= gas {
            Err(())
        } else {
            self.gas_left -= gas;
            Ok(())
        }
    }

    pub fn get_calldata(&self) -> Span {
        (self.env.dis.get_calldata)(self.env.env)
    }

    pub fn set_return_data(&self, data: &[u8]) {
        (self.env.dis.set_return_data)(self.env.env, Span::create(data))
    }

    pub fn get_ask_count(&self) -> i64 {
        (self.env.dis.get_ask_count)(self.env.env)
    }

    pub fn get_min_count(&self) -> i64 {
        (self.env.dis.get_min_count)(self.env.env)
    }

    pub fn get_ans_count(&self) -> i64 {
        (self.env.dis.get_ans_count)(self.env.env)
    }

    pub fn ask_external_data(&self, eid: i64, did: i64, data: &[u8]) {
        (self.env.dis.ask_external_data)(self.env.env, eid, did, Span::create(data))
    }

    pub fn get_external_data_status(&self, eid: i64, vid: i64) -> i64 {
        (self.env.dis.get_external_data_status)(self.env.env, eid, vid)
    }

    pub fn get_external_data(&self, eid: i64, vid: i64) -> Span {
        (self.env.dis.get_external_data)(self.env.env, eid, vid)
    }
}
