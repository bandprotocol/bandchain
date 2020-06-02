use std::convert::TryInto;
use std::io;
use std::mem::{forget, size_of};

mod hint;

const ERROR_NOT_ALL_BYTES_READ: &str = "Not all bytes read";
const ERROR_UNEXPECTED_LENGTH_OF_INPUT: &str = "Unexpected length of input";

/// A data-structure that can be de-serialized from binary format by NBOR.
pub trait OBIDecode: Sized {
    /// Decodes this instance from a given slice of bytes.
    /// Updates the buffer to point at the remaining bytes.
    fn decode(buf: &mut &[u8]) -> io::Result<Self>;

    /// Decode this instance from a slice of bytes.
    fn try_from_slice(v: &[u8]) -> io::Result<Self> {
        let mut v_mut = v;
        let result = Self::decode(&mut v_mut)?;
        if !v_mut.is_empty() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidData,
                ERROR_NOT_ALL_BYTES_READ,
            ));
        }
        Ok(result)
    }

    /// Whether Self is u8.
    /// NOTE: `Vec<u8>` is the most common use-case for serialization and deserialization, it's
    /// worth handling it as a special case to improve performance.
    /// It's a workaround for specific `Vec<u8>` implementation versus generic `Vec<T>`
    /// implementation. See https://github.com/rust-lang/rfcs/pull/1210 for details.
    #[inline]
    fn is_u8() -> bool {
        false
    }
}

impl OBIDecode for u8 {
    #[inline]
    fn decode(buf: &mut &[u8]) -> io::Result<Self> {
        if buf.is_empty() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                ERROR_UNEXPECTED_LENGTH_OF_INPUT,
            ));
        }
        let res = buf[0];
        *buf = &buf[1..];
        Ok(res)
    }

    #[inline]
    fn is_u8() -> bool {
        true
    }
}

macro_rules! impl_for_integer {
    ($type: ident) => {
        impl OBIDecode for $type {
            #[inline]
            fn decode(buf: &mut &[u8]) -> io::Result<Self> {
                if buf.len() < size_of::<$type>() {
                    return Err(io::Error::new(
                        io::ErrorKind::InvalidInput,
                        ERROR_UNEXPECTED_LENGTH_OF_INPUT,
                    ));
                }
                let res = $type::from_be_bytes(buf[..size_of::<$type>()].try_into().unwrap());
                *buf = &buf[size_of::<$type>()..];
                Ok(res)
            }
        }
    };
}

impl_for_integer!(i8);
impl_for_integer!(i16);
impl_for_integer!(i32);
impl_for_integer!(i64);
impl_for_integer!(i128);
impl_for_integer!(u16);
impl_for_integer!(u32);
impl_for_integer!(u64);
impl_for_integer!(u128);

impl OBIDecode for bool {
    #[inline]
    fn decode(buf: &mut &[u8]) -> io::Result<Self> {
        if buf.is_empty() {
            return Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                ERROR_UNEXPECTED_LENGTH_OF_INPUT,
            ));
        }
        let b = buf[0];
        *buf = &buf[1..];
        if b == 0 {
            Ok(false)
        } else if b == 1 {
            Ok(true)
        } else {
            Err(io::Error::new(
                io::ErrorKind::InvalidInput,
                format!("Invalid bool representation: {}", b),
            ))
        }
    }
}

#[cfg(feature = "std")]
impl<T> OBIDecode for Vec<T>
where
    T: OBIDecode,
{
    #[inline]
    fn decode(buf: &mut &[u8]) -> io::Result<Self> {
        let len = u32::decode(buf)?;
        if len == 0 {
            Ok(Vec::new())
        } else if T::is_u8() && size_of::<T>() == size_of::<u8>() {
            let len = len as usize;
            if buf.len() < len {
                return Err(io::Error::new(
                    io::ErrorKind::InvalidInput,
                    ERROR_UNEXPECTED_LENGTH_OF_INPUT,
                ));
            }
            let result = buf[..len].to_vec();
            *buf = &buf[len..];
            // See comment from https://doc.rust-lang.org/std/mem/fn.transmute.html
            // The no-copy, unsafe way, still using transmute, but not UB.
            // This is equivalent to the original, but safer, and reuses the
            // same `Vec` internals. Therefore, the new inner type must have the
            // exact same size, and the same alignment, as the old type.
            //
            // The size of the memory should match because `size_of::<T>() == size_of::<u8>()`.
            //
            // `T::is_u8()` is a workaround for not being able to implement `Vec<u8>` separately.
            let result = unsafe {
                // Ensure the original vector is not dropped.
                let mut v_clone = std::mem::ManuallyDrop::new(result);
                Vec::from_raw_parts(
                    v_clone.as_mut_ptr() as *mut T,
                    v_clone.len(),
                    v_clone.capacity(),
                )
            };
            Ok(result)
        } else if size_of::<T>() == 0 {
            let mut result = Vec::new();
            result.push(T::decode(buf)?);

            let p = result.as_mut_ptr();
            unsafe {
                forget(result);
                let len = len as usize;
                let result = Vec::from_raw_parts(p, len, len);
                Ok(result)
            }
        } else {
            // TODO(16): return capacity allocation when we can safely do that.
            let mut result = Vec::with_capacity(hint::cautious::<T>(len));
            for _ in 0..len {
                result.push(T::decode(buf)?);
            }
            Ok(result)
        }
    }
}

impl OBIDecode for String {
    #[inline]
    fn decode(buf: &mut &[u8]) -> io::Result<Self> {
        String::from_utf8(Vec::<u8>::decode(buf)?)
            .map_err(|err| io::Error::new(io::ErrorKind::InvalidData, err.to_string()))
    }
}
