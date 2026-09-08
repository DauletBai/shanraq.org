# Тег: беттен көрінбейтін үшінші кесте

_Лид (summary):_ **Go курсының отыз үшінші сабағы. Мақалада тег көп, тегте мақала көп — әрі бұл бағандағы тізім емес, бөлек кесте болып сақталады. Құрама кілт, id үшін on conflict, өшіргендегі каскад. Әрі уәде етілген нәрсе: қатенің мәтінін оқымай, кодын сұрау жолы.**

## Не үшін керек

Тег зиянсыз көрінеді: тақырыптың астындағы бірер сөз. Оны баған етіп қою — үтір арқылы тізім жазылған `tags text` — қызықтырады, әрі бұл қате.

Осындай кестеден «go тегі бар мақалаларды бер» деп сұрап көріңіз, сонда мынау басталады:

```sql
select slug, tags from a where tags like '%go%';
```

```
┌────────┬────────────────┐
│  slug  │      tags      │
├────────┼────────────────┤
│ dala   │ go,дала        │
│ golang │ golang,жаңалық │
└────────┴────────────────┘
```

`golang` жауапқа кірді, өйткені оның ішінде `go` бар. Ары қарай тегті бір пәрменмен қайта атауға болмайтыны, оның қанша рет кездескенін санауға болмайтыны және оны екі рет салуға тыйым сала алмайтыныңыз белгілі болады.

Себебі қарапайым: мұндағы байланыс — **көпке-көп**. Мақалада бірнеше тег, тегте бірнеше мақала. Мұндай байланысты бағанмен өрнектеу мүмкін емес — ол үшін үшінші кесте ашады.

## Бірден тұтас

Жаңа қалта, `go mod init sabaq30`:

```go
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	sqlite "modernc.org/sqlite"
)

// uniqueViolation — SQLite unique бұзылғанда қайтаратын код.
// Кодтар тізімі — sqlite.org/rescode.html.
const uniqueViolation = 2067

// ErrTaken — мекенжай бос емес. Өңдеушіге дерекқордың коды емес, осы керек.
var ErrTaken = errors.New("мекенжай бос емес")

func main() {
	os.Remove("blog.db")
	db, err := sql.Open("sqlite", "blog.db?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	setup(db)

	fmt.Println("dala қосылды :", add(db, "dala", "Дала туралы", "go", "дала"))
	fmt.Println("salem қосылды:", add(db, "salem", "Сәлем", "go", "блог"))
	fmt.Println("сол мекенжай :", add(db, "dala", "Басқа", "go"))

	fmt.Println("dala тегтері :", tagsOf(db, "dala"))
	fmt.Println("go тегі бар  :", byTag(db, "go"))
	fmt.Println("дала тегі бар:", byTag(db, "дала"))

	db.Exec(`delete from articles where slug = 'dala'`)
	fmt.Println("dala өшірілгеннен кейін:")
	fmt.Println("  go тегі бар    :", byTag(db, "go"))
	fmt.Println("  барлық байланыс:", count(db, `select count(*) from article_tags`))
	fmt.Println("  барлық тег     :", count(db, `select count(*) from tags`))
}

// add мақаланы және оның тегтерін салады. Екі кесте, бір әрекет — демек,
// транзакция: тегсіз мақала да, мақаласыз тег те ешкімге керек емес.
func add(db *sql.DB, slug, title string, tags ...string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.Exec(`insert into articles (slug, title) values (?, ?)`, slug, title)
	if err != nil {
		return wrap(err, slug)
	}
	articleID, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for _, name := range tags {
		var tagID int64
		// on conflict ... do update — жаңа тегке де, бұрыннан бар тегке де id
		// алудың жалғыз жолы: do nothing бос қайтарар еді.
		err := tx.QueryRow(`insert into tags (name) values (?)
			on conflict (name) do update set name = excluded.name
			returning id`, name).Scan(&tagID)
		if err != nil {
			return err
		}
		if _, err := tx.Exec(
			`insert into article_tags (article_id, tag_id) values (?, ?)`,
			articleID, tagID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// wrap қатенің мәтінін оқымай, кодын сұрайды.
func wrap(err error, slug string) error {
	var serr *sqlite.Error
	if errors.As(err, &serr) && serr.Code() == uniqueViolation {
		return fmt.Errorf("%q: %w", slug, ErrTaken)
	}
	return err
}

func tagsOf(db *sql.DB, slug string) string {
	var out string
	err := db.QueryRow(`select group_concat(t.name, ', ')
		from tags t
		join article_tags at on at.tag_id = t.id
		join articles a on a.id = at.article_id
		where a.slug = ?`, slug).Scan(&out)
	if err != nil {
		return "қате: " + err.Error()
	}
	return out
}

func byTag(db *sql.DB, tag string) string {
	rows, err := db.Query(`select a.slug
		from articles a
		join article_tags at on at.article_id = a.id
		join tags t on t.id = at.tag_id
		where t.name = ?
		order by a.id`, tag)
	if err != nil {
		return "қате: " + err.Error()
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return "қате: " + err.Error()
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return "қате: " + err.Error()
	}
	return fmt.Sprint(out)
}

func count(db *sql.DB, q string) int {
	var n int
	if err := db.QueryRow(q).Scan(&n); err != nil {
		log.Fatal(err)
	}
	return n
}

func setup(db *sql.DB) {
	for _, q := range []string{
		`create table articles (
			id    integer primary key,
			slug  text not null unique,
			title text not null
		) strict`,
		`create table tags (
			id   integer primary key,
			name text not null unique
		) strict`,
		`create table article_tags (
			article_id integer not null references articles (id) on delete cascade,
			tag_id     integer not null references tags (id),
			primary key (article_id, tag_id)
		) strict`,
	} {
		if _, err := db.Exec(q); err != nil {
			log.Fatal(err)
		}
	}
}
```

Шығатыны:

```
dala қосылды : <nil>
salem қосылды: <nil>
сол мекенжай : "dala": мекенжай бос емес
dala тегтері : go, дала
go тегі бар  : [dala salem]
дала тегі бар: [dala]
dala өшірілгеннен кейін:
  go тегі бар    : [salem]
  барлық байланыс: 2
  барлық тег     : 3
```

## Талдау

### Үшінші кесте және құрама кілт

`article_tags` өзінің ешнәрсесін сақтамайды. Ондағы екі баған да — сілтеме: мақалаға және тегке. Бір жол «бұл мақалада бұл тег бар» дегенді білдіреді.

Мұндай кестені байланыстырушы деп атайды, ал оның бастапқы кілті **құрама** — бірден екі бағаннан тұрады:

```sql
primary key (article_id, tag_id)
```

Бұл безендіру ұсағы емес. Құрама кілт «мақала + тег» жұбы тек бір рет кездесе алатынын білдіреді, әрі дерекқор тегті екі рет салғызбайды:

```
Error: stepping, UNIQUE constraint failed: article_tags.article_id, article_tags.tag_id (19)
```

Мұны кодта тексерудің қажеті жоқ — әрі инъекция туралы сабақтан білетініміздей, болмайды да: тексеру мен кірістірудің арасына басқа сұраныс сыналап кіріп үлгереді.

> **Елестетіп көріңіз.** Шақырылғандар тізімі, ондағы әр жол — «мына адам мына тойда». Бұл жолда адам да, той да сипатталмаған: жол тек олардың кездесуін сипаттайды.

### `id` үшін `on conflict ... do update`

Тег жаңа болуы да, бұрыннан болуы да мүмкін — ал `id` екі жағдайда да керек. Сырт көзге оғаш көрінетін жол содан шығады:

```sql
insert into tags (name) values (?)
on conflict (name) do update set name = excluded.name
returning id
```

Былай оқылады: кірістіріп көр; мұндай ат бұрыннан бар болса, оны дәл сол мәнмен «жаңарт»; қалай болғанда да `id` қайтар.

Мұндағы жаңарту әдейі мағынасыз — ол тек жол қозғалды деп саналуы және `returning` оның `id`-ін қайтаруы үшін керек. `do nothing` жазсаңыз болмайды: сонда қайшы келген жол қозғалды деп саналмайды, `returning` ештеңе қайтармайды, ал `Scan` `sql.ErrNoRows` алады. Бұл — белгілі тұзақ, әрі оны дәл осылай айналып өтеді.

### Транзакция мұнда әшекей емес

Өткен сабақта транзакцияны қоятын жер жоқ еді. Енді бар: `add` екі кестеге жазады, ал жұмыстың жартысы ештеңеден жаман. Тегсіз мақала — көнуге болады, ал кірістірілмеген мақала үшін жасалған тег мәңгі қоқыс болып қалады.

Сондықтан `Begin`, `defer tx.Rollback()` және соңында бір `Commit` — дәл талдағанымыз, әрі алғаш рет шынымен керек болып тұр.

### `errors.As`: қатенің мәтінін емес, кодын сұрау

CRUD сабағында уәде еткенмін: қатеден оның кодын сұрау жолы кейін пайда болады деп. Міне ол.

```go
var serr *sqlite.Error
if errors.As(err, &serr) && serr.Code() == uniqueViolation {
	return fmt.Errorf("%q: %w", slug, ErrTaken)
}
```

`errors.Is` «бұл сол қате ме?» деген сұраққа жауап береді. `errors.As` — басқасына: «ішінде мынадай типтегі қате жатқан жоқ па? бар болса, оны маған бер». Ол `%w` тізбегін `Is` сияқты тарқатады, бірақ тапқанын айнымалыға салады, ал сол айнымалының әдістері бар — мұнда `Code()`.

Қатенің мәтінін салыстыруға болмайды: ол нұсқадан нұсқаға да, тілден тілге де өзгереді. `2067` коды ешқашан өзгермейді — бұл `SQLITE_CONSTRAINT_UNIQUE`, әрі ол дерекқор құжаттамасында жазылған.

`wrap` не істеп тұрғанына назар аударыңыз: ол дерекқордың қатесін блогтың қатесіне **аударады**. Сыртқа `ErrTaken` шығады, әрі өңдеушіге SQLite туралы да, сандар туралы да білудің қажеті жоқ. Postgres қаласаңыз — бір функция өзгереді.

Тағы бір нәрсе: драйвер енді `_` арқылы емес, атымен қосылған. SQLite туралы сабақтағы астыңғы сызық «тек тіркелу үшін керек» дегенді білдіретін; бүгін бізге оның типі керек болды, сондықтан аты қайтып келді.

### `on delete cascade`: байланыс кетеді, тег қалады

Мақала өшірілгенде тегтерге не болатыны шығуда көрініп тұр. Байланыс төртеу еді, екеуі қалды. Оларды дерекқор өзі алып тастады, өйткені схемада былай жазылған:

```sql
article_id integer not null references articles (id) on delete cascade
```

Ал тег үшеуі қалды, «дала» енді ешкімге керек болмаса да. Әрі бұл дұрыс: `cascade` тек мақалаға сілтемеде тұр. Тегті соңғы мақаласымен бірге өшіру — дерекқор емес, блог қабылдайтын шешім.

Кері жағы: сілтеме жасалып тұрған тегті дерекқор өшіртпейді:

```
Error: stepping, FOREIGN KEY constraint failed (19)
```

Алдымен байланыстарды алыңыз, сосын тегті. Бұл, әрине, `foreign_keys` прагмасы қосулы тұрғанда ғана жұмыс істейді — SQLite туралы сабақтан бері қосылу жолында тұрған сол прагма.

### Тегтерді қайта жинау: `group_concat`

Мақаланың тегтерін жолмен көрсету — Go-ның емес, дерекқордың жұмысы:

```sql
select group_concat(t.name, ', ') from tags t join article_tags at ...
```

`group_concat` сәйкес келген барлық жолдың мәнін бөлгіш арқылы жабыстырады. `go, дала` шығады — тақырыптың астына шығаратын дәл сол.

Екі нәрседен сақ болыңыз. `group_concat` ішіндегі рет кепілдендірілмеген — ол маңызды болса, оны нақты беріңіз. Әрі тег мүлде болмаса, функция `NULL` қайтарады, ал `Scan` мұны кәдімгі жолға қабылдамайды; блогта бұл `coalesce` немесе `sql.NullString` арқылы емделеді.

## Сабақ картасы

![Сабақ картасы: екі кесте, ал байланыс үшіншісінде](/static/course/go/map-tags-kz.svg)

## Өз сөзіңізбен айтыңыз

Қарамай, дауыстап немесе қағазға жауап беріңіз. Жауаптары — сабақтың соңында.

1. Тегтерді неге мақаланың бағанында тізіммен сақтамайды?
2. `on conflict` ішінде `do nothing` орнына неге мағынасыз жаңарту жазылған?
3. `errors.As` `errors.Is`-тен несімен ерекшеленеді?

## Тапсырма

**Міндетті.** Блогқа тег қосыңыз. Схеманы көші-қон өзгертеді: екі жаңа кесте және байланыстырушысы. Форма тегтерді үтір арқылы жолмен қабылдайды, `Store.Add` мақала мен тегтерді бір транзакциямен салады, мақала беті тегтерді көрсетеді, ал `/tag/{name}` — сол тегі бар мақалалар тізімі.

Мұның бәрі [step-18](https://github.com/DauletBai/shanraq.org/tree/main/course/go-blog/step-18) ішінде жасалып қойған — өзіңіз жазғаннан кейін салыстырыңыз.

**Қалауыңызша.**

- Тегтерді кірістірер алдында бір түрге келтіріңіз: бос орындарды қиып, әріптерді кішірейтіңіз, сонда `Go` мен `go` екі тег болмайды.
- Басты бетте тег бұлтын көрсетіңіз: аты және мақала саны. `group by` мен `count` керек болады.
- `on delete cascade`-ті алып тастап, мақаланы қайта өшіріңіз. Қатені оқып, оны дауыстап түсіндіріңіз.

## Блогта бұл қайда тұрады

Үшінші кесте — блогтағы өз беті жоқ алғашқы кесте. Ол оқырманға емес, байланысқа керек, әрі бұл қалыпты жағдай: нағыз дерекқорларда мұндай кесте көп.

Қарыздар. Тегтер әзірге сөздің бөлігі бойынша ізделмейді — онымен келесі, іздеу туралы сабақ айналысады. Тег бұлты әр жолы қайта саналады. Ал әртүрлі регистрмен енгізілген тегтер қалауыңызша тапсырманы орындағанша бөлек болып қала береді.

## Жауаптар

1. Себебі байланыс көпке-көп: мақалада бірнеше тег, тегте бірнеше мақала. Оны бағанмен өрнектеу мүмкін емес — іздеу `golang`-ты табатын `like '%go%'`-ға айналады, ал қайта атау мен санау мүмкін болмай қалады.
2. Себебі `do nothing` қайшы келген жолды қозғалды деп санамайды, әрі `returning` `id` қайтармайды. Мағынасыз жаңарту жолды қозғалды етеді, сонда `id` жаңа тегке де, бұрыннан барына да келеді.
3. `errors.Is` мәнмен салыстырып, «бұл сол қате ме?» дегенге жауап береді. `errors.As` тізбектен керекті **типтегі** қатені іздеп, оны айнымалыға салады — содан кейін одан, мысалы, кодын сұрауға болады.

## Дереккөздер

- [SQLite: ON CONFLICT және upsert](https://sqlite.org/lang_upsert.html)
- [SQLite: қате кодтары](https://sqlite.org/rescode.html)
- [Go: errors пакеті](https://pkg.go.dev/errors)
