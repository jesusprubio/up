use online::check;

fn main() {
    println!("Online? {}", check(Some(1)).is_ok());
}
