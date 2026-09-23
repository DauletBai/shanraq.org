struct Task {
    title: String,
    done: bool,
}

impl Task {
    fn new(title: String) -> Task {
        Task { title, done: false }
    }

    fn is_done(&self) -> bool {
        self.done
    }

    fn finish(&mut self) {
        self.done = true;
    }
}

fn main() {
    let mut task = Task::new(String::from("Walk"));
    println!("{}: {}", task.title, task.is_done());
    task.finish();
    println!("{}: {}", task.title, task.is_done());
}
