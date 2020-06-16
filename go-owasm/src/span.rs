use crate::error::Error;

#[derive(Copy, Clone)]
#[repr(C)]
pub struct Span {
    pub ptr: *mut u8,
    pub len: usize,
    pub cap: usize,
}

impl Span {
    // Create span.
    pub fn create(data: &[u8]) -> Span {
        Span {
            ptr: data.as_ptr() as *mut u8,
            len: data.len(),
            cap: data.len(),
        }
    }

    /// Read data from the span.
    pub fn read(&self) -> &[u8] {
        unsafe { std::slice::from_raw_parts(self.ptr, self.len) }
    }

    /// Write data to the span.
    pub fn write(&mut self, data: &[u8]) -> Error {
        if self.len + data.len() > self.cap {
            return Error::SpanExceededCapacityError;
        }
        unsafe {
            std::ptr::copy(
                data.as_ptr(),
                self.ptr.offset(self.len as isize),
                data.len(),
            )
        }
        self.len += data.len();
        Error::NoError
    }
}

#[cfg(test)]
mod test {
    use super::*;

    #[test]
    fn test_create_and_read_span_ok() {
        let data: Vec<u8> = vec![1, 2, 3, 4, 5];
        let span = Span::create(data.as_slice());
        let span_data = &span.read();
        let span_data_vec = span_data.to_vec();
        assert_eq!(span_data_vec.len(), data.len());
        assert_eq!(span_data_vec[0], data[0]);
        assert_eq!(span_data_vec[1], data[1]);
        assert_eq!(span_data_vec[2], data[2]);
        assert_eq!(span_data_vec[3], data[3]);
        assert_eq!(span_data_vec[4], data[4]);
    }

    #[test]
    fn test_write_span_ok() {
        let mut empty_space = vec![0u8; 32];
        let mut span = Span {
            ptr: empty_space.as_mut_ptr(),
            len: 0,
            cap: 32,
        };

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
    fn test_write_span_fail() {
        let mut empty_space = vec![0u8; 3];
        let mut span = Span {
            ptr: empty_space.as_mut_ptr(),
            len: 0,
            cap: 3,
        };
        let data: Vec<u8> = vec![1, 2, 3, 4, 5];
        span.write(data.as_slice());
        assert_eq!(
            span.write(data.as_slice()),
            Error::SpanExceededCapacityError
        );
    }
}
