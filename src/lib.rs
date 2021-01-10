/**
 * Copyright (c) 2019, Jesús Rubio <jesusprubio@gmail.com>
 *
 * This source code is licensed under the MIT License found in
 * the LICENSE.txt file in the root directory of this source tree.
 */
use async_std::{
    io::{timeout as tout, Error},
    net::TcpStream,
};
use std::time::Duration;

// Captive portals:
// - https://developer.mozilla.org/en-US/docs/Mozilla/Add-ons/WebExtensions/API/captivePortal
// - http://clients3.google.com/generate_204
// - http://detectportal.firefox.com/success.txt.
const ADDRS: &str = "clients3.google.com:80";
const ADDRS_BACK: &str = "detectportal.firefox.com:80";

async fn probe(addrs: &str, timeout: Option<Duration>) -> Result<TcpStream, Error> {
    match timeout {
        Some(d) => tout(d, TcpStream::connect(addrs)).await,
        _ => TcpStream::connect(addrs).await,
    }
}
/// Check if the internet connection is up. When the check succeeds, the returned future is resolved to `true`.
///
/// * `timeout` - Number of seconds to wait for a response (default: OS dependent)
pub async fn online(timeout: Option<Duration>) -> Result<bool, Error> {
    match probe(ADDRS, timeout).await {
        Ok(_) => Ok(true),
        Err(e) => match probe(ADDRS_BACK, timeout).await {
            Ok(_) => Ok(true),
            // If both failed we return the first one.
            Err(_) => Err(e),
        },
    }
}
