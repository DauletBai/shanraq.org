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
            return Err(String::from("Атау бос"));
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
            None => Err(String::from("Нөмірлер таусылды")),
        }
    }
    fn remove(&mut self, id: u32) -> Result<(), String> {
        for index in 0..self.tasks.len() {
            if self.tasks[index].id == id {
                self.tasks.remove(index);
                return Ok(());
            }
        }
        Err(String::from("ID табылмады"))
    }
}
fn main() {
    let mut o = Organizer::new();
    o.add("Оқу").expect("fixed title");
    o.add("Жазу").expect("fixed title");
    o.remove(1).expect("fixed ID");
    println!("Қалды: {}", o.tasks[0].title);
}

#[test]
fn ids_survive_removal() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Оқу"), Ok(1));
    assert_eq!(o.add("Жазу"), Ok(2));
    assert_eq!(o.remove(1), Ok(()));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 2);
    assert_eq!(o.tasks[0].title, "Жазу");
    assert_eq!(o.next_id, 3);
}

#[test]
fn overflow_changes_nothing() {
    let mut o = Organizer::new();
    assert_eq!(o.add("Оқу"), Ok(1));
    o.next_id = u32::MAX;
    assert_eq!(o.add("Жазу"), Err(String::from("Нөмірлер таусылды")));
    assert_eq!(o.tasks.len(), 1);
    assert_eq!(o.tasks[0].id, 1);
    assert_eq!(o.next_id, u32::MAX);
}
