/**
 * Copyright (c) 2019, Jesús Rubio <jesusprubio@gmail.com>
 *
 * This source code is licensed under the MIT License found in
 * the LICENSE.txt file in the root directory of this source tree.
 */

#[cfg(test)]
use online::*;
use pretty_assertions::assert_eq;

#[async_std::test]
async fn should_work_no_parameters() {
    assert_eq!(online(None).await.unwrap(), true);
}

#[async_std::test]
async fn should_work_timeout() {
    assert_eq!(online(Some(5)).await.unwrap(), true);
}

#[async_std::test]
#[should_panic(expected = "future timed out")]
async fn should_fail_timeout_tiny() {
    online(Some(0)).await.unwrap();
}
