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
mod am {
    use online::*;

    #[test]
    fn should_work_no_parameters() {
        assert_eq!(online(None), Ok(true));
    }

    #[test]
    fn should_work_timeout() {
        assert_eq!(online(Some(6)), Ok(true));
    }
}
