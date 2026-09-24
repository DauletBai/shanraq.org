fn version_text() -> String {
    format!("Нұсқа: {}", env!("CARGO_PKG_VERSION"))
}

fn run(args: &[String]) -> Result<String, String> {
    match args {
        [] => Ok(version_text()),
        [flag] if flag == "--version" => Ok(version_text()),
        _ => Err(String::from("Белгісіз аргумент")),
    }
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    match run(&args) {
        Ok(message) => println!("{message}"),
        Err(message) => {
            eprintln!("{message}");
            std::process::exit(2);
        }
    }
}

#[test]
fn prints_version_for_flag() {
    assert!(run(&[String::from("--version")]).unwrap().contains("0.1.0"));
}

#[test]
fn rejects_unknown_argument() {
    assert!(run(&[String::from("--unknown")]).is_err());
}
