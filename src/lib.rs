/**
 * Copyright (c) 2019, Jesús Rubio <jesusprubio@gmail.com>
 *
 * This source code is licensed under the MIT License found in
 * the LICENSE.txt file in the root directory of this source tree.
 */
use std::net::{SocketAddr, TcpStream};
use std::time::Duration;

use simple_error::SimpleError;

const DEFAULT_TIMEOUT: u64 = 3;

fn connect(addr: &SocketAddr, timeout: Option<Duration>) -> Result<bool, SimpleError> {
    let duration = match timeout {
        Some(tout) => tout,
        _ => Duration::new(DEFAULT_TIMEOUT, 0),
    };

    match TcpStream::connect_timeout(addr, duration) {
        Ok(_) => Ok(true),
        Err(e) => Err(SimpleError::from(e)),
    }
}

/// It uses HTTP and DNS as fallback.
///
/// * `timeout` - Number of seconds to wait for a response (default: 3)
pub fn online(timeout: Option<Duration>) -> Result<bool, SimpleError> {
    //! ```rust
    //! use std::time::Duration;
    //!
    //! use online::*;
    //!
    //! assert_eq!(online(None), Ok(true));
    //!
    //! // with timeout
    //! let timeout = Duration::new(6, 0);
    //! assert_eq!(online(Some(timeout)), Ok(true));
    //! ```

    // Chrome captive portal detection.
    // http://clients3.google.com/generate_204
    let addr = SocketAddr::from(([216, 58, 201, 174], 80));

    match connect(&addr, timeout) {
        Ok(_) => Ok(true),
        Err(e) => match e.as_str() {
            "Network is unreachable (os error 101)" => Ok(false),
            "connection timed out" => {
                // Firefox captive portal detection.
                // http://detectportal.firefox.com/success.txt.
                let addr_fallback = SocketAddr::from(([2, 22, 126, 57], 80));

                match connect(&addr_fallback, timeout) {
                    Ok(_) => Ok(true),
                    Err(err) => match err.as_str() {
                        "connection timed out" => Ok(false),
                        _ => Err(err),
                    },
                }
            }
            _ => Err(e),
        },
    }
}

#[cfg(test)]
mod connect {
    use std::time::Duration;

    use super::connect;
    use std::net::SocketAddr;

    #[test]
    fn should_work_no_parameters() {
        let addr = SocketAddr::from(([8, 8, 8, 8], 53));

        assert_eq!(connect(&addr, None), Ok(true));
    }

    #[test]
    fn should_work_timeout() {
        let addr = SocketAddr::from(([8, 8, 8, 8], 53));
        let timeout = Duration::new(6, 0);

        assert_eq!(connect(&addr, Some(timeout)), Ok(true));
    }

    #[test]
    #[should_panic(expected = "connection timed out")]
    fn should_fail_unreachable() {
        let addr = SocketAddr::from(([8, 8, 8, 8], 8888));

        connect(&addr, None).unwrap();
    }

    #[test]
    #[should_panic(expected = "connection timed out")]
    fn should_fail_unreachable_timeout() {
        let addr = SocketAddr::from(([8, 8, 8, 8], 8888));
        let timeout = Duration::new(6, 0);

        connect(&addr, Some(timeout)).unwrap();
    }
}
