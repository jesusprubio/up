use online::sync::check;

fn main() {
    println!("Online? {}", check(None).is_ok());

    println!("Online (`Result`)? {:?}", check(None).unwrap());
}
