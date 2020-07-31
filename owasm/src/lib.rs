//! # Oracle WebAssembly
//!
//! Owasm (o-wah-some) is the standard library for writing and encoding oracle logic to be for
//! deterministic execution on public ledgers that involve fetching data from external sources.
//! Initially developed by [Band Protocol](https://bandprotocol.com), it is currently used as the
//! go-to standard for writing oracle scripts in the [Decentralized Data Delivery Network].
//!
//! ## Design
//!
//! A D3N oracle script must allow the host to ask for the external data that the script wants
//! (by invoking `__prepare`) and perform aggregation on data points collected from them
//! (by invoking `__execute`).
//!
//! ## Code Structure
//!
//! Owasm library consists of two primary modules:
//! - [owasm/core] is the backbone of the crate. It defines the methods to encode and decode both
//! input commands and data outputs, with the `Oracle` trait to allow arbitrary data types, once
//! implemented the trait, to be able to get converted to/from external data.
//! - [owasm/ext] is the extension library providing convenient ways to write oracle scripts that
//! connect to various public APIs. The library is growing to support more use cases and it open
//! for public contribution!
//!
//! [Band Protocol]: https://bandprotocol.com
//! [Decentralized Data Delivery Network]: https://github.com/bandprotocol/d3n
//! [owasm/core]: core/index.html
//! [owasm/ext]: ext/index.html

#[macro_use]
mod macros;

pub mod core;
pub mod ext;
pub mod oei;
