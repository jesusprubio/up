use std::{
    io::{Error, ErrorKind},
    time::Duration,
};

// Captive portals: https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/captivePortal
// - http://clients3.google.com/generate_204
pub const ADDRS: &str = "clients3.google.com:80";
// - http://detectportal.firefox.com/success.txt
pub const ADDRS_BACK: &str = "detectportal.firefox.com:80";

pub fn parse_timeout(timeout: u64) -> Result<Duration, Error> {
    if timeout == 0 {
        // To be consistent with the async implementation:
        // https://github.com/rust-lang/rust/blob/e51830b90afd339332892a8f20db1957d43bf086/library/std/src/sys/unix/net.rs#L142
        return Err(Error::new(
            ErrorKind::InvalidInput,
            "cannot set a 0 duration timeout",
        ));
    }

    Ok(Duration::from_secs(timeout))
}
