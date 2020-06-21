use crate::error::Error;

#[derive(Copy, Clone)]
#[repr(C)]
/// A `span` is a lightweight struct used to refer to a section of memory. The memory
/// section is not owned by the span, similar to C++'s std::span. The `span`'s creator is
/// responsible for allocating the space and freeing it afterward.
///
/// The primary usecase of `span` is to faciliate communication between Go and Rust.
/// One side allocates space and creates a `span` for the counterpart to read or write
/// without needing to worry about memory management.
pub struct Span {
    pub ptr: *mut u8, // The starting location of Span's memory piece.
    pub len: usize,   // The variable to keep track of how many bytes are writen.
    pub cap: usize,   // The maximum capacity of this span.
}

impl Span {
    /// Creates a read-only `span` from the given memory slice. The result span will be
    /// full, with both `len` and `cap` equal to the provided `data`'s size.
    pub fn create(data: &[u8]) -> Span {
        Span {
            ptr: data.as_ptr() as *mut u8,
            len: data.len(),
            cap: data.len(),
        }
    }

    /// Creates a writable `span` with the provided `ptr` as the starting memory location,
    /// and inital capacity `cap`. The created span has zero length.
    pub fn create_writable(ptr: *mut u8, cap: usize) -> Span {
        Span { ptr, len: 0, cap }
    }

    /// Returns a read-only view of the `span`.
    pub fn read(&self) -> &[u8] {
        unsafe { std::slice::from_raw_parts(self.ptr, self.len) }
    }

    /// Appends data into the `span`. Returns NoError if the write is successful.
    /// The function may fail if the given data exceeeds the capacity of the `span`.
    pub fn write(&mut self, data: &[u8]) -> Error {
        if self.len + data.len() > self.cap {
            return Error::SpanTooSmallError;
        }
        unsafe { std::ptr::copy(data.as_ptr(), self.ptr.offset(self.len as isize), data.len()) }
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
        let mut span = Span::create_writable(empty_space.as_mut_ptr(), 32);
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
        let mut span = Span::create_writable(empty_space.as_mut_ptr(), 3);
        let data: Vec<u8> = vec![1, 2, 3, 4, 5];
        span.write(data.as_slice());
        assert_eq!(span.write(data.as_slice()), Error::SpanTooSmallError);
    }
}
