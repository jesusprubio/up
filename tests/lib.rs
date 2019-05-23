/**
 * Copyright (c) 2019, Jesús Rubio <jesusprubio@member.fsf.org>
 *
 * This source code is licensed under the MIT License found in
 * the LICENSE.txt file in the root directory of this source tree.
 */
extern crate online;

#[macro_use]
extern crate pretty_assertions;

#[cfg(test)]
use online::*;
use std::time::Duration;

#[test]
fn should_work_no_parameters() {
    assert_eq!(online(None), Ok(true));
}

#[test]
fn should_work_timeout() {
    let timeout = Duration::new(6, 0);

    assert_eq!(online(Some(timeout)), Ok(true));
}
