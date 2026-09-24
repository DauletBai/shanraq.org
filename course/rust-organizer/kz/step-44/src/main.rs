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
            println!("Тапсырма {}: {}", self.id, self.title);
        }
    }
}

use crate::organizer::Task;

fn main() {
    let task = Task::new(1, "Rust оқу");
    task.print();
}
