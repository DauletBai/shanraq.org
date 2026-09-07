# A free programming course in Kazakh: Go from scratch

_Лид (summary):_ **A course of 50 lessons: from "how the internet works" to a blog of your own, written in Go and live on the internet at your own domain. The course is written in full, and reading it needs neither an account nor a payment. Kazakh first, with Russian and English beside it. Inside: what comes out at the end, why Go, and who writes in it.**

## What we have opened

"Go: from zero to your own blog" is a course of [50 lessons](https://shanraq.org/course/go?lang=en) that ends with a working blog of your own, live on the internet at your own domain. Not a teaching example in a sandbox — a site whose link you can hand to somebody.

The course is written in full, from "how the internet works" to generics. The blog itself lies beside the lessons: 34 steps in an open repository, a folder per lesson, and any step starts with one command.

The whole text is open always, in Kazakh, Russian and English.

![Fifty lessons: from the first program to a blog on the internet](/static/course/go/map-launch-en.svg)

This is what comes out at the end. It is the front page of the blog you will build: your mark, your name, your entries, the form for a new article and the footer — and all of it your own Go code, without a single third-party library for the site.

![The finished blog from the course: the front page with the list of articles and the publishing form](/static/course/go/blog-final.png)

## Why Go

Go was made at Google for what the existing languages handled badly: services that hold thousands of simultaneous connections and stay up.

**Who writes in it.** The language's own site keeps a page of case studies, and they are not startups: Uber, Netflix, PayPal, Microsoft, Dropbox, Twitch, American Express, Capital One, Salesforce, ByteDance — the company behind TikTok — and MercadoLibre, the largest shop in Latin America. For many of them Go became the main back-end language: PayPal and Salesforce moved to it from other languages and wrote down why.

Separately, there is the infrastructure the internet runs on today. Kubernetes, Docker, Terraform, Prometheus, etcd, containerd, Traefik, Caddy and CockroachDB are written in Go; that is not an opinion but the repository language visible on GitHub. From the recent crop: Ollama, which runs language models on your own machine, is Go as well, with 180 thousand stars.

For Kazakhstan that has a practical consequence. Back-end work is the part of the job that can be done from anywhere, and Go is one of the most sought-after languages in it. The way in is short, too: the language deliberately has few constructs, and it can be learned entire.

And one more reason for this particular choice: the site you are reading is written in Go. The course is assembled from what we do every day, not retold from somebody else's textbook.

## How the lessons are built

The course follows the principles of **Viktor Shatalov**, a Soviet teacher whose system let schoolchildren cover the syllabus several times faster by leaning on understanding rather than memorising.

We took four things from it.

**The whole first, the parts after.** Every lesson opens with a working program. You run it before understanding any of it, and only then take it apart. In the [first lesson](https://shanraq.org/read/internet-qalai-zhumys-istejdi?lang=en) one command shows you a real HTTP response before you have read a line about it.

**A memory anchor.** Each lesson closes its explanation with a map: one screen where everything that matters is drawn as solid shapes and links between them. It can be photographed and kept on a phone; the point is that you could redraw it by hand.

**A picture instead of a definition.** Every unfamiliar term is explained by something from ordinary life. PATH is a courier's list of addresses: the shop may be one street over, but if it is not on the list the courier comes back and says there is no such place. HTTP headers are the writing on an envelope. A status code is the stamp on a form you handed in: accepted, no such case here, there is a fire.

**The right not to get it first time.** No marks, no failed tests, no "you did not pass". There is one required exercise per lesson — small and always doable — and two optional ones.

## What an account is for

None of the text is hidden behind it: the lessons read without one. It is there for something else.

**Having your exercise reviewed.** You paste your solution and it is read as your own code, not compared against a model answer. The reviewer does not hand back a corrected version: it names the line that does not do what you think and says what it does instead. The fixing is yours.

Three checks per lesson. Not to save money — three cost under half a cent — but because a fourth turns the exercise into a guessing game against a reviewer. When they run out, what is offered instead of the machine is the lesson again, or people in the comments.

**Tidying your code.** A button of its own puts your code in order with the same `gofmt` that sits in your editor, coloured the way the lessons are. It spends no attempts. It also catches an unclosed bracket before the review: code like that never reaches the reviewer and costs you nothing.

**Marks for what you have done.** An accepted exercise marks the lesson passed, and the course map shows how far you have come. The mark is earned by the exercise, not by scrolling: a page can be flicked through in a second, and a course that congratulates you for that teaches you to flick.

**Comments.** Questions about a lesson, and answers for other readers.

All of it free, like the course itself.

## What comes next

Plainly about what does not exist yet. One more thing is being built and is not promised as ready: **help as you go** — asking a question on the page itself instead of waiting for an answer in the comments. It will be free; we will write separately when it works.

## How this differs from what already exists

We are neither a state body nor a private school.

**We spend no public money.** The project takes no funding from the treasury, takes part in no state programme, and disburses nothing.

**We take on no obligations towards you and demand none from you.** No intake, no deadlines, no expulsion, no "you must finish in three months". No age limit and no entrance exam. Start, drop it, come back in six months — that is your business.

**We sell nothing at the end.** The course is not a funnel towards a paid tier: there is no paid tier.

This is not a reproach to those arranged differently. We simply have other values — and other freedoms.

## Where to start

Go in order; the lessons are linked.

1. [How the internet works: HTTP request and response in 15 min](https://shanraq.org/read/internet-qalai-zhumys-istejdi?lang=en)
2. [Installing Go and setting up VS Code: your workspace](https://shanraq.org/read/go-ornatu-vs-code-baptau?lang=en)
3. [Your first Go program: a net/http web server in 20 minutes](https://shanraq.org/read/go-dagy-algashqy-bagdarlama?lang=en)
4. [Variables and types in Go: string, int, float64 and bool](https://shanraq.org/read/go-aynymalylar-men-tipter?lang=en)
5. [Functions in Go: parameters, return and two values at once](https://shanraq.org/read/go-funkciyalar-parametrler?lang=en)

Or open the [course map](https://shanraq.org/course/go?lang=en) in full — it shows the whole road to a blog of your own.

What you need to start: a computer, the ability to install a program, and an hour and a half for the first five lessons. Nothing else.

## Sources

- [Go: the language's official site](https://go.dev/)
- [Why Go — the authors on the problems the language was made for](https://go.dev/talks/2012/splash.article)
- [Go: case studies of the companies writing in it](https://go.dev/solutions/case-studies)
- [Course map: "Go: from zero to your own blog"](https://shanraq.org/course/go?lang=en)
