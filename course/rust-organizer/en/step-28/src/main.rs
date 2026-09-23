struct Task {
    title: String,
    minutes: u32,
    done: bool,
}

impl Task {
    fn new(title: String, minutes: u32) -> Task {
        Task {
            title,
            minutes,
            done: false,
        }
    }

    fn is_done(&self) -> bool {
        self.done
    }

    fn finish(&mut self) {
        self.done = true;
    }
}

fn main() {
    let mut task = Task::new(String::from("Reading"), 20);
    println!(
        "{}: {} min, done {}",
        task.title,
        task.minutes,
        task.is_done()
    );
    task.finish();
    println!("After: {}", task.is_done());
}
