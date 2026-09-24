use std::collections::{HashMap, HashSet};

fn main() {
    let tags = ["work", "home", "work"];
    let mut counts: HashMap<&str, usize> = HashMap::new();
    let mut unique: HashSet<&str> = HashSet::new();
    for tag in tags {
        let old = match counts.get(tag) {
            Some(value) => *value,
            None => 0,
        };
        counts.insert(tag, old + 1);
        unique.insert(tag);
    }
    let mut sorted: Vec<&str> = unique.into_iter().collect();
    sorted.sort();
    println!("Tags: {:?}", sorted);
    println!("Work: {}", counts["work"]);
}
