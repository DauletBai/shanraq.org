# Жинақтау: білетініңіздің бәрінен консольдегі жазба кітапшасы

_Лид (summary):_ **Go курсының он сегізінші сабағы және тіл туралы модульдің соңы. Таныс бөлшектерден құралған бір аяқталған бағдарлама: құрылым, тілім, сөздік, сілтегішті әдістер, қателер, өз пакеті және тестер. Пернетақтадан пәрмен оқудан басқа жаңа ештеңе жоқ: үйренгенді жинау — бөлек шеберлік.**

## Не үшін керек

Он бес сабақ бойы бөлшектерді бір-бірлеп үйрендіңіз. Бүгін олар тұтас жұмыс істейтін бағдарламаға жиналады: пәрмен қабылдайды, жазба сақтайды, іздейді, санайды және өшіреді.

Жинау — бөлек шеберлік, оны да бөлек жаттықтырады. Тілдің барлық құрылымын білетін адам солардан жұмыс істейтін дүние құрастыра алмауы мүмкін: бір-бірлеп бәрі түсінікті болған жерде, бірге келгенде «мұны қайда қою керек» деген сұрақ пайда болады.

> **Елестетіп көріңіз.** Үстелде жатқан бөлшектерден жиналатын велосипед. Бірде-бір жаңа бөлшек әкелінген жоқ. Шеберлік тағы біреуін табуда емес, осыларды жинай білуде.

## Бірден тұтас

`sabaq-16` қалтасы, ішінде — `notes` қалтасы. Төрт файл.

```
sabaq-16/
  go.mod
  main.go
  notes/
    note.go
    store.go
    store_test.go
```

`notes/note.go` — жазбаның өзі:

```go
package notes

import (
	"fmt"
	"unicode/utf8"
)

type Note struct {
	ID    int
	Text  string
	Words int
}

func (n Note) Letters() int {
	return utf8.RuneCountInString(n.Text)
}

func (n Note) String() string {
	return fmt.Sprintf("%d. %s (%d сөз · %d әріп)", n.ID, n.Text, n.Words, n.Letters())
}
```

`notes/store.go` — барлық әрекеті бар қойма:

```go
package notes

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrEmpty    = errors.New("жазба бос")
	ErrNotFound = errors.New("жазба табылмады")
)

type Store struct {
	next  int
	items []Note
}

func NewStore() *Store {
	return &Store{next: 1}
}

func (s *Store) Add(text string) (Note, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return Note{}, ErrEmpty
	}
	n := Note{ID: s.next, Text: text, Words: len(strings.Fields(text))}
	s.items = append(s.items, n)
	s.next++
	return n, nil
}

func (s *Store) All() []Note { return s.items }

func (s *Store) Find(part string) []Note {
	var out []Note
	for _, n := range s.items {
		if strings.Contains(strings.ToLower(n.Text), strings.ToLower(part)) {
			out = append(out, n)
		}
	}
	return out
}

func (s *Store) Delete(id int) error {
	for i, n := range s.items {
		if n.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("delete %d: %w", id, ErrNotFound)
}

// Top counts how often each word occurs and returns them in alphabetical
// order. The map counts, the slice puts the answer in an order a map does not
// promise -- both lessons in four lines.
func (s *Store) Top() []string {
	count := map[string]int{}
	for _, n := range s.items {
		for _, w := range strings.Fields(strings.ToLower(n.Text)) {
			count[w]++
		}
	}
	words := make([]string, 0, len(count))
	for w := range count {
		words = append(words, w)
	}
	slices.Sort(words)

	out := make([]string, 0, len(words))
	for _, w := range words {
		out = append(out, fmt.Sprintf("%s — %d", w, count[w]))
	}
	return out
}

func (s *Store) Len() int { return len(s.items) }
```

`main.go` — адаммен сөйлесу:

```go
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sabaq16/notes"
)

func main() {
	store := notes.NewStore()
	in := bufio.NewScanner(os.Stdin)

	fmt.Println("жазба кітапшасы. пәрмендер: add, list, find, top, del, help, exit")

	for {
		fmt.Print("> ")
		if !in.Scan() {
			return
		}
		cmd, arg, _ := strings.Cut(strings.TrimSpace(in.Text()), " ")

		switch cmd {
		case "add":
			n, err := store.Add(arg)
			if err != nil {
				fmt.Println("қате:", err)
				continue
			}
			fmt.Println("қосылды:", n)

		case "list":
			if store.Len() == 0 {
				fmt.Println("бос")
				continue
			}
			for _, n := range store.All() {
				fmt.Println(n)
			}

		case "find":
			found := store.Find(arg)
			if len(found) == 0 {
				fmt.Println("ештеңе табылмады")
				continue
			}
			for _, n := range found {
				fmt.Println(n)
			}

		case "del":
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("қате: нөмір керек")
				continue
			}
			if err := store.Delete(id); err != nil {
				if errors.Is(err, notes.ErrNotFound) {
					fmt.Println("қате: мұндай нөмір жоқ")
					continue
				}
				fmt.Println("қате:", err)
				continue
			}
			fmt.Println("өшірілді:", id)

		case "top":
			for _, line := range store.Top() {
				fmt.Println(line)
			}

		case "help":
			fmt.Println("add мәтін · list · find бөлік · top · del нөмір · exit")

		case "exit", "":
			return

		default:
			fmt.Printf("белгісіз пәрмен: %q, help теріңіз\n", cmd)
		}
	}
}
```

Жинап, іске қосыңыз:

```
go mod init sabaq16
go run .
```

Пәрмен термес бұрын — бағдарлама `add Дала туралы жазба` дегенге не жауап беретінін дауыстап айтыңыз. Содан кейін теріп, салыстырыңыз:

```
жазба кітапшасы. пәрмендер: add, list, find, top, del, help, exit
> add Дала туралы жазба
қосылды: 1. Дала туралы жазба (3 сөз · 17 әріп)
> add Шаңырақ деген не
қосылды: 2. Шаңырақ деген не (3 сөз · 16 әріп)
> add
қате: жазба бос
> top
дала — 1
деген — 1
жазба — 1
не — 1
туралы — 1
шаңырақ — 1
> del 9
қате: мұндай нөмір жоқ
> exit
```

## Талдау

### Мұндағының қайсысы қай сабақтан

Бағдарламада сіз бірінші рет көріп тұрған бірде-бір құрылым жоқ. Өзіңізді тізім бойынша тексеріңіз:

| Кодта не тұр | Қай сабақтан |
|---|---|
| `type Note struct` | құрылымдар |
| `[]Note` және `append` | тілім |
| `Top` ішіндегі `map[string]int` | сөздік |
| `func (n Note) String()` | әдістер |
| `func (s *Store) Add(...)` | сілтегіштер |
| `ErrEmpty`, `errors.Is` | қателер |
| `package notes`, бас әріптер | пакеттер |
| `utf8.RuneCountInString` | руналар |
| `store_test.go` | алғашқы тест |

Кестеден танымайтын нәрсе болса — осы сабаққа емес, сол сабаққа қайтыңыз.

### Жалғыз жаңалық: пәрмендерді оқу

```go
in := bufio.NewScanner(os.Stdin)
for {
	fmt.Print("> ")
	if !in.Scan() {
		return
	}
	...
}
```

`os.Stdin` — адамның терминалда теретіні. `bufio.Scanner` ағынды жолдарға бөліп, бір-бірлеп береді: `Scan` келесісін оқиды да, жолдар біткенде `false` қайтарады.

> **Елестетіп көріңіз.** Турникет. Адамдар ағыны тұтас қалың болып келеді, ал турникеттен бір-бірлеп өтеді, әрқайсысымен бөлек айналысуға болады.

`strings.Cut` жолды бірінші бос орын бойынша қияды: сол жақта пәрмен, оң жақта қалғанының бәрі. `add Дала туралы жазба` үшін `add` пен `Дала туралы жазба` шығады — бос орындары бар мәтін шашылып кетпейді.

### `if` баспалдағының орнына `switch`

```go
switch cmd {
case "add":
	...
case "list":
	...
default:
	...
}
```

Бір мәнді бірнеше нұсқамен салыстырғанда, `switch` `if — else if` тізбегінен жақсы оқылады. Go-да ол келесі жағдайға құлап кетпейді: `break` жазудың қажеті жоқ.

`default` — белгісіз пәрменге берілетін жауап. Қате терілген сөзге үндемейтін бағдарлама бұзылған сияқты көрінеді.

### `notes` пакеті және оның сыртқа ашылатын есігі

Сыртқа `Note`, `Store`, `NewStore`, әрекеттер және екі қате шығады. Ішінде `next` пен `items` өрістері қалады — оларға сырттан жету мүмкін емес.

Тексеріңіз: `main.go` тізімге ешқай жерде тікелей тимейді. Ол **сұрайды** — `store.Add`, `store.All`, `store.Delete`. Сондықтан ертең тізімді дерекқорға ауыстыруға болады, ал `main.go` оны байқамайды.

### Қателер: бүкіл пакетке екеу

```go
var (
	ErrEmpty    = errors.New("жазба бос")
	ErrNotFound = errors.New("жазба табылмады")
)
```

`Delete` нөмірді қосып, `ErrNotFound`-ты орайды, ал `main.go` оны `errors.Is` арқылы танып, адамша «мұндай нөмір жоқ» деп басады. Пайдаланушыға — түсінікті сөйлем, бағдарламашыға — қате мәтініндегі егжей-тегжей.

### Тест пернетақта туралы білмейді

```go
package notes

import (
	"errors"
	"testing"
)

func TestAddAndDelete(t *testing.T) {
	s := NewStore()

	n, err := s.Add("Дала туралы")
	if err != nil {
		t.Fatalf("дұрыс жазба қабылданбады: %v", err)
	}
	if n.ID != 1 || n.Words != 2 {
		t.Errorf("алдық ID %d, %d сөз; күттік 1 және 2", n.ID, n.Words)
	}

	if _, err := s.Add("   "); !errors.Is(err, ErrEmpty) {
		t.Errorf("бос жазба үшін алдық %v, күттік ErrEmpty", err)
	}

	if err := s.Delete(99); !errors.Is(err, ErrNotFound) {
		t.Errorf("жоқ нөмір үшін алдық %v, күттік ErrNotFound", err)
	}

	if err := s.Delete(1); err != nil {
		t.Fatalf("өшіру сәтсіз: %v", err)
	}
	if s.Len() != 0 {
		t.Errorf("өшірілгеннен кейін %d жазба қалды", s.Len())
	}
}

func TestLettersCountsLettersNotBytes(t *testing.T) {
	n := Note{Text: "Шаңырақ"}
	if got := n.Letters(); got != 7 {
		t.Errorf("алдық %d әріп, күттік 7", got)
	}
}
```

Тест `main`-ді іске қоспай, `Store`-ды тексереді. Енгізу, шығару және пәрмендерді талдау оған қатысы жоқ — бұл айла емес, логиканың `main`-де емес, пакетте жатқанының салдары.

> **Елестетіп көріңіз.** Стендтегі қозғалтқыш. Оны машинадан бөлек айналдырады: ақауды табу оңай әрі әр жолы жолға шығудың қажеті жоқ.

`go test ./...` жіберіңіз — пакеттің екі файлы да бірден тексеріледі.

### Жазбалар неге шыққанға дейін ғана тұрады

Бағдарламаны жапсаңыз — жазбалар жоғалады. Бұл аяқталмағандық емес: файлдар мен дерекқор келесі модульдерде болады, ал бүгін маңыздысы — **логиканың бәрі орнында** және қойма шынайы болғанда өзгермейді.

Кейін оны дәл солай қосасыз: `Store` сол әдістерімен қалады, тек ішіндегісі өзгереді.

## Сабақ картасы

![Сабақ картасы: алты таныс бөлшек бір бағдарламаға жиналады](/static/course/go/map-capstone-kz.svg)

## Өз сөзіңізбен айтыңыз

Қарамай, дауыстап немесе қағазға жазып жауап беріңіз. Жауаптары — сабақтың соңында.

1. Тест неге `main.go` қасында емес, `notes` пакетінде жатыр?
2. `main.go` `store.items` деп жазып көрсе не болады?
3. `ErrNotFound`-ты сол күйінде қайтаруға болатын болса, `Delete` қатені не үшін орайды?

## Тапсырма

**Міндетті.** Жазбаларға тег қосыңыз. `Note` ішінде `Tag string` өрісі, бар жазбаға тег қоятын `tag нөмір сөз` пәрмені және сол тегі бар жазбаларды көрсететін `bytag сөз` пәрмені. Жоқ нөмірге арналған қатені дайын күйінде алыңыз — `ErrNotFound`. Екі пәрменге де тест жазып, `go test ./...` қатесіз өтетініне жетіңіз.

**Қалауыңызша.**

- `stats` пәрмені: барлығы қанша жазба, сөз және әріп.
- `find` тег бойынша да іздейтін етіңіз.
- Дайын жазба кітапшасын GitHub-қа README-мен бірге бөлек репозиторий етіп жариялаңыз — бұл көрсетуге болатын алғашқы аяқталған бағдарлама.

## Блогта қайда тұрады

Тіл туралы модуль осымен бітті. Әрі қарай — веб: сол `Store`, бірақ оны `add` пен `list` емес, HTTP өңдеушілері шақырады, ал жауапты терминалға емес, браузерге береді.

Бүгін жазғаныңыздың бәрі қалады. Тек сұрақ қоятын кім екені өзгереді.

## Жауаптар

1. Себебі ол `Store`-ды тексереді, ал `Store` `notes` пакетінде тұрады. Сол пакеттегі тест ішкі өрістерді де көреді, тексеру үшін ештеңені сыртқа ашудың қажеті жоқ. Оның үстіне ол пернетақтадан енгізуге тәуелді болмайды.
2. Бағдарлама жиналмайды: `items` кіші әріппен жазылған және тек `notes` пакетінің ішінде көрінеді. Компилятор `cannot refer to unexported field items` дейді.
3. Хабарда табылмаған нөмір қалуы үшін: `delete 9: жазба табылмады`. Орам егжей-тегжей қосады әрі қатенің өзін сақтайды, сондықтан `errors.Is` оны бұрынғысынша таниды.

## Дереккөздер

- [bufio пакеті — құжаттама](https://pkg.go.dev/bufio)
- [strings.Cut — құжаттама](https://pkg.go.dev/strings#Cut)
- [Go spec: Switch statements](https://go.dev/ref/spec#Switch_statements)
- [Effective Go: switch](https://go.dev/doc/effective_go#switch)
