use online::check;

#[async_std::main]
async fn main() {
    println!("Online? {}", check(None).await.is_ok());

    println!("Online (`Result`)? {:?}", check(None).await.unwrap());
}
