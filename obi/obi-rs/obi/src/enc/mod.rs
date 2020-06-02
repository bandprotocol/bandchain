use std::mem::size_of;
use std::{io, io::Write};

const DEFAULT_ENCODE_CAPACITY: usize = 1024;

/// A data-structure that can be encoded into binary format by NBOR.
pub trait OBIEncode {
    fn encode<W: Write>(&self, writer: &mut W) -> io::Result<()>;

    /// Encode this instance into a vector of bytes.
    fn try_to_vec(&self) -> io::Result<Vec<u8>> {
        let mut result = Vec::with_capacity(DEFAULT_ENCODE_CAPACITY);
        self.encode(&mut result)?;
        Ok(result)
    }

    /// Whether Self is u8.
    /// NOTE: `Vec<u8>` is the most common use-case for encode and decode, it's
    /// worth handling it as a special case to improve performance.
    /// It's a workaround for specific `Vec<u8>` implementation versus generic `Vec<T>`
    /// implementation. See https://github.com/rust-lang/rfcs/pull/1210 for details.
    #[inline]
    fn is_u8() -> bool {
        false
    }
}

impl OBIEncode for u8 {
    #[inline]
    fn encode<W: Write>(&self, writer: &mut W) -> io::Result<()> {
        writer.write_all(std::slice::from_ref(self))
    }

    #[inline]
    fn is_u8() -> bool {
        true
    }
}

macro_rules! impl_for_integer {
    ($type: ident) => {
        impl OBIEncode for $type {
            #[inline]
            fn encode<W: Write>(&self, writer: &mut W) -> io::Result<()> {
                writer.write_all(&self.to_be_bytes())
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

impl OBIEncode for bool {
    #[inline]
    fn encode<W: Write>(&self, writer: &mut W) -> io::Result<()> {
        (if *self { 1u8 } else { 0u8 }).encode(writer)
    }
}

impl OBIEncode for String {
    #[inline]
    fn encode<W: Write>(&self, writer: &mut W) -> io::Result<()> {
        writer.write_all(&(self.len() as u32).to_be_bytes())?;
        writer.write_all(self.as_bytes())?;
        Ok(())
    }
}

#[cfg(feature = "std")]
impl<T> OBIEncode for Vec<T>
where
    T: OBIEncode,
{
    #[inline]
    fn encode<W: Write>(&self, writer: &mut W) -> io::Result<()> {
        writer.write_all(&(self.len() as u32).to_be_bytes())?;
        if T::is_u8() && size_of::<T>() == size_of::<u8>() {
            // The code below uses unsafe memory representation from `&[T]` to `&[u8]`.
            // The size of the memory should match because `size_of::<T>() == size_of::<u8>()`.
            //
            // `T::is_u8()` is a workaround for not being able to implement `Vec<u8>` separately.
            let buf = unsafe { std::slice::from_raw_parts(self.as_ptr() as *const u8, self.len()) };
            writer.write_all(buf)?;
        } else {
            for item in self {
                item.encode(writer)?;
            }
        }
        Ok(())
    }
}
