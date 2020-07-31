use crate::env::Env;
use crate::span::Span;

use owasm::core::error::Error;
use owasm::core::vm;

pub struct VMEnv {
    env: Env,       // The execution environment for callbacks to Golang.
    span_size: i64, // Maximum span size for communication between Rust & Go.
}

impl VMEnv {
    pub fn new(env: Env, span_size: i64) -> VMEnv {
        VMEnv {
            env: env,
            span_size: span_size,
        }
    }
}

impl vm::Env for VMEnv {
    fn get_span_size(&self) -> i64 { self.span_size }

    fn get_calldata(&self) -> Result<Vec<u8>, Error> {
        let mut mem: Vec<u8> = Vec::with_capacity(self.span_size as usize);
        let mut span = Span::create_writable(mem.as_mut_ptr(), self.span_size as usize);
        match (self.env.dis.get_calldata)(self.env.env, &mut span) {
            Error::NoError => {
                unsafe {
                    mem.set_len(span.len);
                }
                Ok(mem)
            }
            err => Err(err),
        }
    }

    fn set_return_data(&self, data: &[u8]) -> Result<(), Error> {
        match (self.env.dis.set_return_data)(self.env.env, Span::create(data)) {
            Error::NoError => Ok(()),
            err => Err(err),
        }
    }

    fn get_ask_count(&self) -> i64 { (self.env.dis.get_ask_count)(self.env.env) }

    fn get_min_count(&self) -> i64 { (self.env.dis.get_min_count)(self.env.env) }

    fn get_ans_count(&self) -> Result<i64, Error> {
        let mut ans_count = 0;
        match (self.env.dis.get_ans_count)(self.env.env, &mut ans_count) {
            Error::NoError => Ok(ans_count),
            err => Err(err),
        }
    }

    fn ask_external_data(&self, eid: i64, did: i64, data: &[u8]) -> Result<(), Error> {
        match (self.env.dis.ask_external_data)(self.env.env, eid, did, Span::create(data)) {
            Error::NoError => Ok(()),
            err => Err(err),
        }
    }

    fn get_external_data_status(&self, eid: i64, vid: i64) -> Result<i64, Error> {
        let mut status = 0;
        match (self.env.dis.get_external_data_status)(self.env.env, eid, vid, &mut status) {
            Error::NoError => Ok(status),
            err => Err(err),
        }
    }

    fn get_external_data(&self, eid: i64, vid: i64) -> Result<Vec<u8>, Error> {
        let mut mem: Vec<u8> = Vec::with_capacity(self.span_size as usize);
        let mut span = Span::create_writable(mem.as_mut_ptr(), self.span_size as usize);
        match (self.env.dis.get_external_data)(self.env.env, eid, vid, &mut span) {
            Error::NoError => {
                unsafe {
                    mem.set_len(span.len);
                }
                Ok(mem)
            }
            err => Err(err),
        }
    }
}
