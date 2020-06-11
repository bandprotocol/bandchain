use crate::error::Error;

#[derive(Copy, Clone)]
#[repr(C)]
pub struct Span {
    pub ptr: *mut u8,
    pub len: usize,
    pub cap: usize,
}

impl Span {
    // TODO
    pub fn create(data: &[u8]) -> Span {
        Span {
            ptr: data.as_ptr() as *mut u8,
            len: data.len(),
            cap: data.len(),
        }
    }

    /// TODO
    pub fn read(&self) -> &[u8] {
        unsafe { std::slice::from_raw_parts(self.ptr, self.len) }
    }

    /// TODO
    pub fn write(&mut self, data: &[u8]) -> Error {
        if self.len + data.len() > self.cap {
            return Error::SpanExceededCapacityError
        }
        unsafe { std::ptr::copy(data.as_ptr(), self.ptr.offset(self.len as isize), data.len()) }
        self.len += data.len();
        Error::NoError
    }
}

pub fn add(a: i32, b: i32) -> i32 {
    a + b
}

#[cfg(test)]
mod test {
    use super::*;
    
    #[test]
    fn test_add() {
        assert_eq!(add(1, 2), 3);
    }

    #[test]
    fn test_success_write_span_ok() {
        let mut empty_space = vec![0u8; 32];
        let mut span = Span{ ptr:  empty_space.as_mut_ptr(), len: 0, cap: 32 };

        let data: Vec<u8> = vec![1, 2, 3, 4, 5];
        assert_eq!(span.write(data.as_slice()), Error::NoError);
        assert_eq!(span.len, 5);
        assert_eq!(span.cap, 32);
        assert_eq!(empty_space[0], 1);
        assert_eq!(empty_space[5], 0);
        
        assert_eq!(span.write(data.as_slice()), Error::NoError);
        assert_eq!(span.len, 10);
        assert_eq!(span.cap, 32);
        assert_eq!(empty_space[0], 1);
        assert_eq!(empty_space[5], 1);
        assert_eq!(empty_space[9], 5);
    }
    
    #[test]
    fn test_success_write_span_fail() {
        let mut empty_space = vec![0u8; 3];
        let mut span = Span{ ptr:  empty_space.as_mut_ptr(), len: 0, cap: 3 };

        let data: Vec<u8> = vec![1, 2, 3, 4, 5];
        span.write(data.as_slice());

        assert_eq!(span.write(data.as_slice()), Error::SpanExceededCapacityError);
        println!("{:?}", empty_space);
        // assert_eq!(add(1, 2), 3);
    }
}
