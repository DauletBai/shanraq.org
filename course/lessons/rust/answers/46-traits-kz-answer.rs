trait Summary {
    fn summary(&self) -> String;
}

struct Task {
    title: String,
}

impl Summary for Task {
    fn summary(&self) -> String {
        self.title.clone()
    }
}

struct Note {
    text: String,
}

impl Summary for Note {
    fn summary(&self) -> String {
        self.text.clone()
    }
}

fn show<T: Summary>(item: &T) {
    println!("{}", item.summary());
}

fn main() {
    let task = Task {
        title: String::from("Rust оқу"),
    };
    let note = Note {
        text: String::from("Кітап сатып алу"),
    };
    show(&task);
    show(&note);
}
