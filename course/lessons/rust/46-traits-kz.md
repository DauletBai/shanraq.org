# Трейт: ортақ әрекетке қойылатын талап

_Қысқаша (summary):_ **46-сабақ. Екі түрге ортақ әрекет беріп, жалпыланған қызметті шектеу.**

Алдын ала: [жалпыланған түрлер](/read/rust-45-generics?lang=kz) мен [әдістер](/read/rust-28-methods?lang=kz). Бұрынғы жобаны сақтаңыз; мысалдағы дерек әзірше жадта ғана болады.

## Таныс бейне және тірек сызба

Үйірмедегі екі адам өзін таныстыра алады, бірақ мәліметті әрқайсысы өз жолымен сақтайды. «Өзіңді таныстыр» деген талап **трейтке** (trait) ұқсайды: ол орындалуға тиіс әрекетті атайды. Әр қатысушы сол әрекетті өзі атқарады. Ұқсатудың шегі: компилятор әдістің жазылуын тексереді, ал мәтіннің адамға түсінікті болуын кепілдемейді.

**trait — әдіске қойылатын талап → impl — белгілі түрдің орындауы → T: Summary — ортақ қызметке қойылған шарт.** `trait Summary` ішінде `summary(&self) -> String` әдісінің жазылымы беріледі; соңындағы нүктелі үтір дененің бұл жерде жоқ екенін көрсетеді. `impl Summary for Task` әдісті `Task` үшін орындайды. `show<T: Summary>` ішіндегі қос нүкте түрге қойылған шартты білдіреді: `show` осы талапты орындайтын кез келген `T` мәнін қабылдайды. `&T` мәнді тек қарызға алады. `clone()` қайтаратын мәтінге жеке көшірме жасайды; тапсырмадағы атау орнында қалады.

## Бір қызметті тапсырмаға қолданыңыз

```rust
trait Summary {
    fn summary(&self) -> String;
}

struct Task {
    title: String,
}

impl Summary for Task {
    fn summary(&self) -> String {
        self.title.clone()
    }
}

fn show<T: Summary>(item: &T) {
    println!("{}", item.summary());
}

fn main() {
    let task = Task {
        title: String::from("Rust оқу"),
    };
    show(&task);
}
```
```text
Rust оқу
```

`show(&task)` шақыруы `Task` үшін жазылған әдісті қолданады. `impl Summary for Task` алынса, компилятор шақыруды қабылдамайды. `Debug` сияқты дайын трейт үшін кейде `#[derive(Debug)]` атрибуты компилятордан лайықты кодты жасауды сұрайды. Ол тек қолдауы бар трейті және жарамды өрістермен жұмыс істейді; кез келген трейт өздігінен жасалмайды. Мұнда мәтіннің қалай шығатынын өзіміз белгілеп отырмыз.

## derive қалай жұмыс істейді

`Debug` өрістерді әзірлеушіге көрсетеді; бұл қолданушыға арналған дайын мәтін емес. `#[derive(Debug)]` — құрылым жарияланар алдындағы атрибут; ол компилятордан дайын трейтті орындауды сұрайды. `{:?}` сол тексеру көрінісін шығарады.

```rust
#[derive(Debug)]
struct Task {
    title: String,
}

fn main() {
    let task = Task {
        title: String::from("Rust оқу"),
    };
    println!("{task:?}");
}
```
```text
Task { title: "Rust оқу" }
```

## Бетке қарамай жауап беріңіз

1. `trait Summary` не талап етеді?
2. `impl Summary for Task` нені орындайды?
3. Трейті жоқ түрді `show<T: Summary>` неге қабылдамайды?

## Тапсырма

**Міндетті.** `struct Note { text: String }` қосып, оған `Summary` трейтін орындаңыз. «Кітап сатып алу» деген жазба жасап, оны да сол `show` қызметіне беріңіз. Алдымен екі жолдың қандай болатынын болжаңыз.

## Жауаптар

`impl Summary for Task` үлгісіндегі `Task` атауын `Note`, `title` өрісін `text` деп ауыстырыңыз. Жазбаның мәтіні өзінде қалсын.

<!-- task-answer -->
```rust
trait Summary {
    fn summary(&self) -> String;
}

struct Task {
    title: String,
}

impl Summary for Task {
    fn summary(&self) -> String {
        self.title.clone()
    }
}

struct Note {
    text: String,
}

impl Summary for Note {
    fn summary(&self) -> String {
        self.text.clone()
    }
}

fn show<T: Summary>(item: &T) {
    println!("{}", item.summary());
}

fn main() {
    let task = Task {
        title: String::from("Rust оқу"),
    };
    let note = Note {
        text: String::from("Кітап сатып алу"),
    };
    show(&task);
    show(&note);
}
```
```text
Rust оқу
Кітап сатып алу
```

## Тексергеннен кейін

Екі жол да бір `show` арқылы шықты; нақты әдіс дерек түріне байланысты таңдалды. `T` белгісі түсініксіз болса, [45-сабақты](/read/rust-45-generics?lang=kz) қайталаңыз. [Трейт туралы ресми тарау](https://doc.rust-lang.org/book/ch10-02-traits.html).

[Алдыңғы сабақ](/read/rust-45-generics?lang=kz) · [Мазмұны](/course/rust?lang=kz)
