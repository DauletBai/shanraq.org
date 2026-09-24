fn longer<'a>(left: &'a str, right: &'a str) -> &'a str {
    if left.len() >= right.len() {
        left
    } else {
        right
    }
}

fn main() {
    let short = String::from("task");
    let result;
    {
        let long = String::from("shopping");
        result = longer(&short, &long);
        println!("{result}");
    }
    println!("{short}");
}
