use online::tokio::check;

#[tokio::main]
async fn main() {
    println!("Online? {}", check(None).await.is_ok());
}
