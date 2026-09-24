mod organizer {
    pub struct Task {
        id: u32,
        title: String,
    }
    impl Task {
        pub fn new(id: u32, title: &str) -> Self {
            Self {
                id,
                title: String::from(title),
            }
        }
        pub fn print(&self) {
            println!("Task {}: {}", self.id, self.title);
        }
        pub fn id(&self) -> u32 {
            self.id
        }
    }
}

use crate::organizer::Task;

fn main() {
    let task = Task::new(1, "Read Rust");
    task.print();
    println!("Number: {}", task.id());
}
