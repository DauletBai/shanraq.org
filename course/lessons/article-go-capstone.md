# Сборка: блокнот в консоли из всего, что вы уже знаете

_Лид (summary):_ **Восемнадцатый урок курса по Go и конец модуля о языке. Одна законченная программа из знакомых деталей: структура, срез, словарь, методы с указателем, ошибки, свой пакет и тесты. Ничего нового, кроме чтения команд с клавиатуры: собрать выученное в работающую вещь — отдельное умение.**

## Зачем это нужно

Пятнадцать уроков вы учили детали по одной. Сегодня они собираются в программу, которая работает целиком: принимает команды, хранит записи, ищет, считает и удаляет.

Собрать — отдельное умение, и его тренируют отдельно. Человек, знающий все конструкции языка, может не суметь сложить из них рабочую вещь: там, где по одной всё понятно, вместе появляется вопрос «а куда это положить».

> **Образ.** Велосипед из деталей, которые уже лежат на столе. Ни одной новой детали не привезли. Умение не в том, чтобы достать ещё одну, а в том, чтобы собрать эти.

## Сразу целиком

Папка `sabaq-16`, внутри — папка `notes`. Четыре файла.

```
sabaq-16/
  go.mod
  main.go
  notes/
    note.go
    store.go
    store_test.go
```

`notes/note.go` — сама запись:

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
	return fmt.Sprintf("%d. %s (слов: %d · букв: %d)", n.ID, n.Text, n.Words, n.Letters())
}
```

`notes/store.go` — хранилище со всеми действиями:

```go
package notes

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

var (
	ErrEmpty    = errors.New("запись пустая")
	ErrNotFound = errors.New("запись не найдена")
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

`main.go` — разговор с человеком:

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

	fmt.Println("блокнот. команды: add, list, find, top, del, help, exit")

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
				fmt.Println("ошибка:", err)
				continue
			}
			fmt.Println("добавлено:", n)

		case "list":
			if store.Len() == 0 {
				fmt.Println("пусто")
				continue
			}
			for _, n := range store.All() {
				fmt.Println(n)
			}

		case "find":
			found := store.Find(arg)
			if len(found) == 0 {
				fmt.Println("ничего не найдено")
				continue
			}
			for _, n := range found {
				fmt.Println(n)
			}

		case "del":
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("ошибка: нужен номер")
				continue
			}
			if err := store.Delete(id); err != nil {
				if errors.Is(err, notes.ErrNotFound) {
					fmt.Println("ошибка: такого номера нет")
					continue
				}
				fmt.Println("ошибка:", err)
				continue
			}
			fmt.Println("удалено:", id)

		case "top":
			for _, line := range store.Top() {
				fmt.Println(line)
			}

		case "help":
			fmt.Println("add текст · list · find часть · top · del номер · exit")

		case "exit", "":
			return

		default:
			fmt.Printf("неизвестная команда: %q, наберите help\n", cmd)
		}
	}
}
```

Соберите и запустите:

```
go mod init sabaq16
go run .
```

Прежде чем набирать команды — скажите вслух, что программа ответит на `add заметка о степи`. Потом наберите и сверьте:

```
блокнот. команды: add, list, find, top, del, help, exit
> add заметка о степи
добавлено: 1. заметка о степи (слов: 3 · букв: 15)
> add что такое шанырак
добавлено: 2. что такое шанырак (слов: 3 · букв: 17)
> add
ошибка: запись пустая
> top
заметка — 1
о — 1
степи — 1
такое — 1
что — 1
шанырак — 1
> del 9
ошибка: такого номера нет
> exit
```

## Разбор

### Что здесь откуда

В программе нет ни одной конструкции, которую вы видите впервые. Проверьте себя по списку:

| Что в коде | Из какого урока |
|---|---|
| `type Note struct` | структуры |
| `[]Note` и `append` | срез |
| `map[string]int` в `Top` | словарь |
| `func (n Note) String()` | методы |
| `func (s *Store) Add(...)` | указатели |
| `ErrEmpty`, `errors.Is` | ошибки |
| `package notes`, заглавные буквы | пакеты |
| `utf8.RuneCountInString` | руны |
| `store_test.go` | первый тест |

Если что-то в таблице не узнаётся — вернитесь к тому уроку, не к этому.

### Единственное новое: чтение команд

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

`os.Stdin` — то, что человек набирает в терминале. `bufio.Scanner` разбивает поток на строки и отдаёт по одной: `Scan` читает следующую и возвращает `false`, когда строки кончились.

> **Образ.** Турникет. Поток людей идёт сплошной массой, а через турникет проходит по одному, и с каждым можно разобраться отдельно.

`strings.Cut` разрезает строку по первому пробелу: слева команда, справа всё остальное. Для `add заметка о степи` получится `add` и `Дала туралы жазба` — текст с пробелами не рассыпается.

### `switch` вместо лестницы из `if`

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

Когда одно значение сравнивают с несколькими вариантами, `switch` читается лучше цепочки `if — else if`. В Go он не проваливается в следующий случай: `break` писать не нужно.

`default` — ответ на неизвестную команду. Программа, которая молчит в ответ на опечатку, кажется сломанной.

### Пакет `notes` и его дверь наружу

Наружу выходят `Note`, `Store`, `NewStore`, действия и две ошибки. Внутри остаются поля `next` и `items` — снаружи к ним не подобраться.

Проверьте: `main.go` нигде не трогает список напрямую. Он **просит** — `store.Add`, `store.All`, `store.Delete`. Поэтому завтра список можно заменить на базу данных, и `main.go` не заметит.

### Ошибки: две на весь пакет

```go
var (
	ErrEmpty    = errors.New("запись пустая")
	ErrNotFound = errors.New("запись не найдена")
)
```

`Delete` заворачивает `ErrNotFound`, добавляя номер, а `main.go` узнаёт её через `errors.Is` и печатает человеческое «мұндай нөмір жоқ». Пользователю — понятная фраза, программисту — подробность в тексте ошибки.

### Тест не знает про клавиатуру

```go
package notes

import (
	"errors"
	"testing"
)

func TestAddAndDelete(t *testing.T) {
	s := NewStore()

	n, err := s.Add("О степи")
	if err != nil {
		t.Fatalf("правильная запись не принята: %v", err)
	}
	if n.ID != 1 || n.Words != 2 {
		t.Errorf("получили ID %d, %d слов; ждали 1 и 2", n.ID, n.Words)
	}

	if _, err := s.Add("   "); !errors.Is(err, ErrEmpty) {
		t.Errorf("для пустой записи получили %v, ждали ErrEmpty", err)
	}

	if err := s.Delete(99); !errors.Is(err, ErrNotFound) {
		t.Errorf("для чужого номера получили %v, ждали ErrNotFound", err)
	}

	if err := s.Delete(1); err != nil {
		t.Fatalf("удаление не удалось: %v", err)
	}
	if s.Len() != 0 {
		t.Errorf("после удаления осталось записей: %d", s.Len())
	}
}

func TestLettersCountsLettersNotBytes(t *testing.T) {
	n := Note{Text: "шанырак"}
	if got := n.Letters(); got != 7 {
		t.Errorf("получили %d букв, ждали 7", got)
	}
}
```

Тест проверяет `Store`, не запуская `main`. Ввод, вывод и разбор команд его не касаются — и это не хитрость, а следствие того, что логика лежит в пакете, а не в `main`.

> **Образ.** Двигатель на стенде. Его гоняют отдельно от машины: так проще найти неисправность и не надо каждый раз выезжать на дорогу.

Запустите `go test ./...` — проверяются оба файла пакета сразу.

### Почему записи живут только до выхода

Закроете программу — записи пропадут. Это не недоделка: файлы и база будут в следующих модулях, а сегодня важно, что **вся логика уже на месте** и не изменится, когда хранилище станет настоящим.

Ровно так вы и добавите его позже: `Store` останется с теми же методами, поменяется только то, что внутри.

## Карта урока

![Карта урока: шесть знакомых деталей собираются в одну программу](/static/course/go/map-capstone-ru.svg)

## Скажите своими словами

Не подглядывая, ответьте вслух или на бумаге. Ответы — в конце урока.

1. Почему тест лежит в пакете `notes`, а не рядом с `main.go`?
2. Что произойдёт, если `main.go` попробует написать `store.items`?
3. Зачем `Delete` заворачивает ошибку, если можно было вернуть `ErrNotFound` как есть?

## Задание

**Обязательное.** Добавьте записям тег. Поле `Tag string` в `Note`, команда `tag нөмір слово`, которая ставит тег существующей записи, и команда `bytag слово`, которая показывает записи с этим тегом. Ошибку для несуществующего номера возьмите готовую — `ErrNotFound`. Допишите тест на обе команды и добейтесь `go test ./...` без ошибок.

**По желанию.**

- Команда `stats`: сколько всего записей, слов и букв.
- Сделайте так, чтобы `find` искал и по тегу тоже.
- Выложите готовый блокнот на GitHub отдельным репозиторием с README — это первая законченная программа, которую можно показать.

## Куда это встанет в блоге

Модуль о языке на этом закончен. Дальше — веб: тот же `Store`, но вместо `add` и `list` его будут звать обработчики HTTP, а вместо терминала отвечать браузеру.

Всё, что вы написали сегодня, останется. Изменится только то, кто задаёт вопросы.

## Ответы

1. Потому что он проверяет `Store`, а `Store` живёт в пакете `notes`. Тест в том же пакете видит и внутренние поля, и ничего не нужно открывать наружу ради проверки. Заодно так он не зависит от ввода с клавиатуры.
2. Программа не соберётся: `items` написано со строчной буквы и видно только внутри пакета `notes`. Компилятор скажет `cannot refer to unexported field items`.
3. Чтобы в сообщении остался номер, которого не нашли: `delete 9: жазба табылмады`. Обёртка добавляет подробность и при этом сохраняет саму ошибку, поэтому `errors.Is` по-прежнему её узнаёт.

## Источники

- [Пакет bufio — документация](https://pkg.go.dev/bufio)
- [strings.Cut — документация](https://pkg.go.dev/strings#Cut)
- [Go spec: Switch statements](https://go.dev/ref/spec#Switch_statements)
- [Effective Go: switch](https://go.dev/doc/effective_go#switch)
