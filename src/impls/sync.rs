use std::{
    io::{Error, ErrorKind},
    net::{TcpStream, ToSocketAddrs},
    time::Duration,
};

#[path = "../utils.rs"]
mod utils;

// - The signature `connect_timeout` is different to `connect` (no DNS resolution and 1 unique target).
// - We choose to go with `connect_timeout` (vs `connect` + `io::timeout`) when available due to its better precision.
fn connect_timeout(addr_str: &str, timeout: Duration) -> Result<(), Error> {
    let addrs_iter = addr_str.to_socket_addrs()?;

    if let Some(addr) = addrs_iter.into_iter().next() {
        if let Err(e) = TcpStream::connect_timeout(&addr, timeout) {
            Err(e)
        } else {
            Ok(())
        }
    } else {
        Err(Error::from(ErrorKind::NotFound))
    }
}

// TODO: Update after new additions
/// Synchronous implementation.
///
/// * `timeout` - Number of seconds to wait for a response (default: OS dependent)
pub fn check(timeout: Option<u64>) -> Result<(), Error> {
    if let Some(t) = timeout {
        let dur = utils::parse_timeout(t)?;

        return if connect_timeout(utils::ADDRS[0], dur).is_ok() {
            Ok(())
        } else {
            connect_timeout(utils::ADDRS[1], dur)
        };
    }

    if TcpStream::connect(utils::ADDRS[0]).is_ok() {
        Ok(())
    } else if let Err(e) = TcpStream::connect(utils::ADDRS[1]) {
        Err(e)
    } else {
        Ok(())
    }
}
