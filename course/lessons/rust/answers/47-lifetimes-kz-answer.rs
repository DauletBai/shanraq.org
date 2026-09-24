fn longer<'a>(left: &'a str, right: &'a str) -> &'a str {
    if left.len() >= right.len() {
        left
    } else {
        right
    }
}

fn main() {
    let short = String::from("іс");
    let result;
    {
        let long = String::from("дүкенге бару");
        result = longer(&short, &long);
        println!("{result}");
    }
    println!("{short}");
    let chosen = longer(short.as_str(), "жоспар");
    println!("{chosen}");
}
