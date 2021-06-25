use std::{
    io::{Error, ErrorKind},
    net::{TcpStream, ToSocketAddrs},
    time::Duration,
};

#[path = "./utils.rs"]
mod utils;

// - The signature `connect_timeout` is a bit different.
// - We choose to use `connect_timeout` (vs `connect` + `io::timeout`) when available due to its better precision.
fn connect_timeout(addr_str: &str, timeout: Duration) -> Result<(), Error> {
    let addrs_iter = addr_str.to_socket_addrs()?;

    if let Some(addr) = addrs_iter.into_iter().next() {
        if let Err(e) = TcpStream::connect_timeout(&addr, timeout) {
            return Err(e);
        }

        return Ok(());
    }
    return Err(Error::from(ErrorKind::NotFound));
}

/// Synchronous implementation.
///
/// * `timeout` - Number of seconds to wait for a response (default: OS dependent)
pub fn check(timeout: Option<u64>) -> Result<(), Error> {
    if let Some(t) = timeout {
        let dur = utils::parse_timeout(t)?;

        // First try, ignoring error (if any).
        if connect_timeout(utils::ADDRS, dur).is_ok() {
            return Ok(());
        }

        // Fallback.
        return connect_timeout(utils::ADDRS_BACK, dur);
    }

    // Second try.
    if TcpStream::connect(utils::ADDRS).is_ok() {
        return Ok(());
    }

    if let Err(e) = TcpStream::connect(utils::ADDRS_BACK) {
        return Err(e);
    }

    Ok(())
}
