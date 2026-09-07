# CRUD: мақала үстіндегі төрт әрекет және тілімсіз өмір

_Лид (summary):_ **Go курсының отызыншы сабағы. Мақалалар тілімнен дерекқорға көшеді әрі қайта қосудан аман қалады. Дерекқорды қалған кодтан жасыратын Store типі, мақала үстіндегі төрт әрекет және бөлімге тұрарлық бір нәрсе: жоқ нәрсені жаңарту — қате емес, оны тек RowsAffected айтады.**

## Не үшін керек

Блог әлі күнге мақалаларды тілімде ұстайды. Формалар сабағындағы форма жадқа жазады, әрі жазылғанның бәрі `Ctrl+C`-ге дейін өмір сүреді.

Өткен сабақта бағдарлама дерекқордан бір жолды оқуды үйренді. Бүгін ол қалғанының бәрін үйренеді — жасау, өзгерту, өшіру, — әрі тілім мүлде кетеді.

Бұл төрт әрекетті бас әріптерімен солай атайды: create, read, update, delete. CRUD — қандай да бір технология емес, жазбамен не істеуге болатынының жай ғана тізімі.

## Бірден тұтас

Жаңа қалта, `go mod init sabaq27`, драйвер сол:

```
go get modernc.org/sqlite
```

`main.go`:

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

type Article struct {
	ID    int64
	Slug  string
	Title string
	Words int
}

var ErrNotFound = errors.New("мақала табылмады")

// Store — блогтың мақалалармен істей алатынының бәрі. Сырттан тек әрекеттер
// көрінеді; іште дерекқор тұрғаны осы типтен тыс шықпайды.
type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("%s ашу: %w", path, err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("дерекқор жауап бермейді: %w", err)
	}
	if _, err := db.Exec(`create table if not exists articles (
		id    integer primary key,
		slug  text    not null unique,
		title text    not null,
		words integer not null default 0
	) strict`); err != nil {
		db.Close()
		return nil, fmt.Errorf("схема: %w", err)
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() error { return s.db.Close() }

// Create — «C»: мақаланы салып, нөмірін білу.
func (s *Store) Create(a Article) (int64, error) {
	res, err := s.db.Exec(
		`insert into articles (slug, title, words) values (?, ?, ?)`,
		a.Slug, a.Title, a.Words)
	if err != nil {
		return 0, fmt.Errorf("%q жасау: %w", a.Slug, err)
	}
	return res.LastInsertId()
}

// Get — бір мақалаға арналған «R».
func (s *Store) Get(slug string) (Article, error) {
	var a Article
	err := s.db.QueryRow(
		`select id, slug, title, words from articles where slug = ?`, slug).
		Scan(&a.ID, &a.Slug, &a.Title, &a.Words)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return Article{}, fmt.Errorf("%q: %w", slug, ErrNotFound)
	case err != nil:
		return Article{}, fmt.Errorf("%q оқу: %w", slug, err)
	}
	return a, nil
}

// All — тізімге арналған «R».
func (s *Store) All() ([]Article, error) {
	rows, err := s.db.Query(`select id, slug, title, words from articles order by id`)
	if err != nil {
		return nil, fmt.Errorf("тізім: %w", err)
	}
	defer rows.Close()

	var out []Article
	for rows.Next() {
		var a Article
		if err := rows.Scan(&a.ID, &a.Slug, &a.Title, &a.Words); err != nil {
			return nil, fmt.Errorf("жолды талдау: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// Update — «U». RowsAffected-ке назар аударыңыз: онсыз өзгертетін ештеңе
// болмағанын білмейміз.
func (s *Store) Update(a Article) error {
	res, err := s.db.Exec(
		`update articles set title = ?, words = ? where slug = ?`,
		a.Title, a.Words, a.Slug)
	if err != nil {
		return fmt.Errorf("%q өзгерту: %w", a.Slug, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%q: %w", a.Slug, ErrNotFound)
	}
	return nil
}

// Delete — «D», әрі RowsAffected-пен сол әңгіме.
func (s *Store) Delete(slug string) error {
	res, err := s.db.Exec(`delete from articles where slug = ?`, slug)
	if err != nil {
		return fmt.Errorf("%q өшіру: %w", slug, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("%q: %w", slug, ErrNotFound)
	}
	return nil
}

func main() {
	os.Remove("blog.db")
	store, err := Open("blog.db")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer store.Close()

	id, err := store.Create(Article{Slug: "dala", Title: "Дала туралы", Words: 400})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("жасадық, нөмірі:", id)

	if _, err := store.Create(Article{Slug: "dala", Title: "Дала туралы тағы"}); err != nil {
		fmt.Println("сол мекенжай екінші рет:", err)
	}

	a, _ := store.Get("dala")
	fmt.Printf("оқыдық: %q, %d сөз\n", a.Title, a.Words)

	if err := store.Update(Article{Slug: "dala", Title: "Көктемдегі дала", Words: 450}); err != nil {
		log.Fatal(err)
	}
	a, _ = store.Get("dala")
	fmt.Printf("түзетуден кейін: %q, %d сөз\n", a.Title, a.Words)

	err = store.Update(Article{Slug: "kokek", Title: "Ештеңе"})
	fmt.Println("жоқ мақаланы түзету:", errors.Is(err, ErrNotFound), "—", err)

	if err := store.Delete("dala"); err != nil {
		log.Fatal(err)
	}
	err = store.Delete("dala")
	fmt.Println("екі рет өшіру:", errors.Is(err, ErrNotFound), "—", err)

	list, _ := store.All()
	fmt.Println("қалған мақала:", len(list))
}
```

`go run .`:

```
жасадық, нөмірі: 1
сол мекенжай екінші рет: "dala" жасау: constraint failed: UNIQUE constraint failed: articles.slug (2067)
оқыдық: "Дала туралы", 400 сөз
түзетуден кейін: "Көктемдегі дала", 450 сөз
жоқ мақаланы түзету: true — "kokek": мақала табылмады
екі рет өшіру: true — "dala": мақала табылмады
қалған мақала: 0
```

## Талдау

### `Store` — дерекқорды жасыратын тип

```go
type Store struct{ db *sql.DB }
```

Бір өріс, әрі ол кіші әріппен — демек, пакеттен тыс көрінбейді. Сыртқа тек әрекеттер шығады: `Create`, `Get`, `All`, `Update`, `Delete`.

Бұл — пакеттер сабағындағы сол тәсіл, бірақ енді оның көрінетін құны мен пайдасы бар. `store.Get("dala")` шақыратын өңдеуші SQL туралы да, SQLite туралы да, дерекқордың бар екені туралы да білмейді. Ертең Postgres керек болса — бір файл өзгереді. Тесттерде жадтағы қойма керек болса — осы бес әрекетті интерфейспен сипаттасаңыз, ол қойылады.

Кері жағы да бар: мұндай тип жасалмайынша, SQL бөліктері өңдеушілерге тарап кетеді, ал оларды кейін қайта жинау бір күнге түседі.

### `Create`: `insert` және дерекқор берген нөмір

```go
res, err := s.db.Exec(
	`insert into articles (slug, title, words) values (?, ?, ?)`,
	a.Slug, a.Title, a.Words)
...
return res.LastInsertId()
```

`Exec` — ештеңе қайтармайтын сұраныстарға арналған: жауапта жол жоқ, тек нәтиже бар. Одан `LastInsertId` алады — дерекқор өзі берген нөмір.

SQLite-та бұл істейді, ал Postgres-те `LastInsertId` мүлде қолдау таппайды: онда `returning id` бар `QueryRow` алады. Көшу тосын болмауы үшін бұны алдын ала білген жөн.

### `Get`: бір жол және өз қатеміз

```go
case errors.Is(err, sql.ErrNoRows):
	return Article{}, fmt.Errorf("%q: %w", slug, ErrNotFound)
```

`sql.ErrNoRows` дегенді `ErrNotFound`-қа аударамыз — мұны өткен сабақта талдадық. Мұнда басқасы маңызды: `ErrNotFound` **осы пакетте**, `Store`-дың қасында жарияланған, әрі өңдеушілер дәл соны тексереді. Әрі қарайғы кодта `sql` сөзі кездеспеуі керек.

### `All`: тізім және циклден кейінгі қате

```go
return out, rows.Err()
```

`rows.Err()` дегенді тізіммен бірге соңғы жолда қайтардық. Оқу ортасында бұзылса, шақырушы жарты тізімді де, қатені де алады — әрі қате бойынша тізімге сенуге болмайтынын түсінеді.

### `Update`: `RowsAffected` неге міндетті

Мінеки сабақтағы ең бастысы, әрі бұл Go туралы емес, дерекқорлар туралы.

Дерекқордан жоқ мақаланы өзгертуді сұраңыз. Қате **болмайды**:

```
update dala    қате=<nil>  RowsAffected=1
update kokek   қате=<nil>  RowsAffected=0
delete dala    қате=<nil>  RowsAffected=1
delete kokek   қате=<nil>  RowsAffected=0
```

Екі жағдайда да `қате=<nil>`. Дерекқор дәл сұралғанды істеді: сәйкес жолдарды тапты — олар нөл болып шықты — әрі сол нөлдің бәрін өзгертті. Оның көзқарасынша бәрі сәтті өтті.

Шындықты білудің жалғыз жолы — қанша жол қозғалғанын сұрау:

```go
n, err := res.RowsAffected()
if n == 0 {
	return fmt.Errorf("%q: %w", a.Slug, ErrNotFound)
}
```

Осы үш жолсыз түзету формасы оқырманға ештеңе сақтамай тұрып «сақталды» дейді. Қатені дерекқор да, Go да, журнал да көрсетпейді — ол бір аптадан кейін, адам өз түзетуін таппағанда білінеді.

> **Елестетіп көріңіз.** Жоқ мекенжайға жіберілген хатты пошта үнсіз алып кетеді. Түбіртек бар, жеткізу жоқ, ал бұл туралы тек «қанша хат жетті» деп сұрағанда білуге болады.

### `Delete`: сол әңгіме

```
delete kokek   қате=<nil>  RowsAffected=0
```

Жоқ мақаланы өшіру де қате емес. Тексеру сол, әрі сол себеппен: «өшіру» түймесі өшіретін ештеңе болмағанда «өшірілді» деп жауап бермеуі керек.

### Бірегейлікті сіз емес, дерекқор тексереді

```
второй раз тот же адрес: создать "dala": constraint failed: UNIQUE constraint failed: articles.slug (2067)
```

Кодта тексергіңіз келеді — алдымен `Get`, таппасаңыз, сосын `Create` — бұл түсінікті әрі дұрыс емес. Екі сұраныстың арасында басқа біреу сол жолды салып үлгереді, әрі тексеру құтқармайды. Дерекқор мұны атомарлы тексереді, өйткені `unique` схемада жазылған.

Дұрысы — байқап көріп, қатені талдау. Оны хабарлама мәтіні бойынша талдау нашар; келесі сабақтарда қатеден оның кодын сұрау тәсілі пайда болады. Әзірге мұнда кім басты екенін түсіну жеткілікті.

## Сабақ картасы

![Сабақ картасы: төрт әрекет, ал біреуі үнсіз жауап береді](/static/course/go/map-crud-kz.svg)

## Өз сөзіңізбен айтыңыз

Қарамай, дауыстап немесе қағазға жауап беріңіз. Жауаптар — сабақтың соңында.

1. `Store` `*sql.DB`-ны экспортталмайтын өріске не үшін жасырады?
2. `update` бірде-бір жол таппаса, `Exec` не қайтарады?
3. Мекенжайдың бірегейлігін неге салудың алдындағы код емес, дерекқор тексереді?

## Тапсырма

**Міндетті.** Блогқа түзету бетін қосыңыз. `GET /edit/{slug}` тақырыбы толтырылған форманы көрсетеді, `POST /edit/{slug}` мақаланы `Update` арқылы өзгертіп, `/read/{slug}` бетіне бағыттаумен жауап береді. Мақала жоқ болса — `404`, әрі мұны `Get`-тен білдіңіз бе, `RowsAffected`-тен бе, маңызды емес.

Блогтың дерекқорға көшуін өзіңіз талдаудың қажеті жоқ: ол [step-15](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-15) ішінде жасалған — түзету бетін жазып болғаннан кейін сонымен салыстырыңыз.

**Қалауыңызша.**

- `Update` ішінен `RowsAffected` тексеруін алып тастап, жоқ мақаланы түзеткенде форма не көрсететінін қараңыз.
- `Delete`-ті жұмсақ етіңіз: өшірудің орнына `deleted_at` қойып, мұндай мақалаларды тізімде көрсетпеңіз.
- Мақала бетіне «өшіру» түймесін қосыңыз. Ол сұранысты қай әдіспен жіберуі керегін және неге `GET` емес екенін ойлаңыз.

## Блогта бұл қайда тұрады

Тілім кетеді, дерекқор қалады: бүкіл модуль осы қадам үшін басталған еді. Блог енді қайта қосудан аман қалады.

Қарыздар. Схеманы әлі де бағдарламаның өзі іске қосылғанда жасайды — кесте оны болып, өзгере бастағанда олай болмайды; бұл — келесі сабақ, миграциялар туралы. Сұраныс мәтіндері әдістердің ішінде жатыр, ал олар жиырма болғанда бөлек орын керек болады. Әрі `Update` барлық өрісті бірден өзгертеді: бір мақаланы екі адам түзетсе, екіншісі біріншісін байқамай өшіріп жібереді.

## Жауаптар

1. Қалған код іште дерекқор бар екенін білмеуі үшін. Сырттан тек әрекеттер көрінеді, сондықтан SQLite-ты Postgres-ке не тесттерге арналған жадтағы қоймаға ауыстыру барлық өңдеушіні емес, бір файлды қозғайды.
2. Жаман ештеңе жоқ: `err` `nil` болады, өйткені сұраныс сәтті орындалды — жай ғана нөл жол сәйкес келді. «Өзгерттік» пен «өзгертетін ештеңе болмады» дегенді тек `res.RowsAffected()` арқылы ажыратуға болады.
3. Өйткені кодтағы тексеру екі сұраныстан тұрады, ал олардың арасында саңылау бар: басқа сұраныс сол жолды салып үлгереді, әрі екі тексеру де өтеді. Схемадағы `unique` атомарлы тексеріледі, сондықтан дұрыс реті — салып көріп, қатені талдау.

## Дереккөздер

- [database/sql: Result](https://pkg.go.dev/database/sql#Result)
- [database/sql: DB.Exec](https://pkg.go.dev/database/sql#DB.Exec)
- [SQLite: деректі өзгерту](https://sqlite.org/lang_update.html)
- [SQLite: қате кодтары](https://sqlite.org/rescode.html)
