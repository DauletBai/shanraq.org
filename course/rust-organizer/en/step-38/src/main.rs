struct Task {
    id: u32,
    title: String,
    done: bool,
}
struct Organizer {
    tasks: Vec<Task>,
    next_id: u32,
}

impl Organizer {
    fn new() -> Self {
        Self {
            tasks: Vec::new(),
            next_id: 1,
        }
    }
    fn add(&mut self, title: &str) -> Result<u32, String> {
        if title.trim().is_empty() {
            return Err(String::from("Title is empty"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title: String::from(title),
                    done: false,
                });
                self.next_id = next;
                Ok(id)
            }
            None => Err(String::from("No IDs left")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn complete(&mut self, id: u32) -> Result<(), String> {
        for task in &mut self.tasks {
            if task.id == id {
                task.done = true;
                return Ok(());
            }
        }
        Err(String::from("ID not found"))
    }
    fn show(&self) {
        for task in &self.tasks {
            let mark = if task.done { "done" } else { "open" };
            println!("ID {}: {} ({mark})", task.id, task.title);
        }
    }
}

fn main() {
    let mut organizer = Organizer::new();
    let first = organizer.add("Buy bread").expect("fixed valid title");
    let second = organizer.add("Read a book").expect("fixed valid title");
    organizer.remove(first).expect("fixed existing ID");
    let third = organizer.add("Walk").expect("fixed valid title");
    organizer.complete(second).expect("fixed existing ID");
    println!("ID: {second}, {third}");
    organizer.show();
}
