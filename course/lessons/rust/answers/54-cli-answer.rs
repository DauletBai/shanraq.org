fn run(args: &[String]) -> i32 {
    match args.first().map(|word| word.as_str()) {
        None | Some("help") | Some("--help") if args.len() <= 1 => {
            println!("Команды: help, list, status");
            0
        }
        Some("list") if args.len() == 1 => {
            println!("Задач пока нет");
            0
        }
        Some("status") if args.len() == 1 => {
            println!("Задач: 0");
            0
        }
        _ => {
            eprintln!("Неизвестная команда");
            2
        }
    }
}

fn main() {
    let args: Vec<String> = std::env::args().skip(1).collect();
    let code = run(&args);
    if code != 0 {
        std::process::exit(code);
    }
}
