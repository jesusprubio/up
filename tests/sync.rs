#[cfg(test)]
use online::sync::check;
use pretty_assertions::assert_eq;

#[test]
fn should_work_no_parameters() {
    assert_eq!(check(None).is_ok(), true);
}

#[test]
fn should_work_timeout() {
    assert_eq!(check(Some(5)).is_ok(), true);
}

#[test]
#[should_panic(expected = "cannot set a 0 duration timeout")]
fn should_fail_timeout_zero() {
    check(Some(0)).unwrap();
}
