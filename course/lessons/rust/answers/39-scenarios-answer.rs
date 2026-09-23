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
            return Err(String::from("Пустое название"));
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
            None => Err(String::from("Номера закончились")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID не найден"))
    }
}
fn main() {
    let mut o = Organizer::new();
    o.add("Читать").expect("fixed title");
    o.add("Писать").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Осталась: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Читать"), Ok(1));
    assert_eq!(o.add("Писать"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Писать");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Читать"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Писать"), Err(String::from("Номера закончились")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}

#[test]
fn invalid_remove_keeps_ids() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Читать"), Ok(1));
    assert_eq!(o.add("Писать"), Ok(2));
    assert_eq!(o.remove(99), Err(String::from("ID не найден")));
    assert_eq!(o.tasks.len(), 2);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.tasks[1].id, 2);
    assert_eq!(o.add("Читать"), Ok(3));
}
