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
        pub fn id(&self) -> u32 {
            self.id
        }
    }
}

use crate::organizer::Task;

fn main() {
    let task = Task::new(1, "Rust оқу");
    task.print();
    println!("Нөмір: {}", task.id());
}
