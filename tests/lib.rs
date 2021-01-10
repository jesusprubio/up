/**
 * Copyright (c) 2019, Jesús Rubio <jesusprubio@gmail.com>
 *
 * This source code is licensed under the MIT License found in
 * the LICENSE.txt file in the root directory of this source tree.
 */

#[cfg(test)]
use online::*;
use pretty_assertions::assert_eq;
use std::time::Duration;

#[async_std::test]
async fn should_work_no_parameters() {
    assert_eq!(online(None).await.unwrap(), true);
}

#[async_std::test]
async fn should_work_timeout() {
    let timeout = Duration::from_secs(5);

    assert_eq!(online(Some(timeout)).await.unwrap(), true);
}

#[async_std::test]
#[should_panic(expected = "future timed out")]
async fn should_fail_timeout_tiny() {
    let timeout = Duration::from_nanos(1);

    online(Some(timeout)).await.unwrap();
}
