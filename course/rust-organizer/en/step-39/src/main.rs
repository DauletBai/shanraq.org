struct Task {
    id: u32,
    title: String,
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
            return Err(String::from("Empty title"));
        }
        match self.next_id.checked_add(1) {
            Some(next) => {
                let id = self.next_id;
                self.tasks.push(Task {
                    id,
                    title: String::from(title),
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
}
fn main() {
    let mut o = Organizer::new();
    o.add("Read").expect("fixed title");
    o.add("Write").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Remaining: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Read"), Ok(1));
    assert_eq!(o.add("Write"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Write");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Read"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Write"), Err(String::from("No IDs left")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}
