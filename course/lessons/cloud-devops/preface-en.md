# Before you begin: why Cloud & DevOps now

_Lead (summary):_ **A detailed introduction to this free practical course: why infrastructure skills are growing with cloud, AI, and data centres, what Kazakhstan is building, and what you will learn without a promise of an instant career.**

## Why we chose this course

A cloud service looks intangible, but it always runs on physical processors, disks, networks, cooling, and electricity. As organizations use more AI, streaming data, digital public services, and online products, they need more compute and people who can release, observe, secure, and recover software systems.

The International Energy Agency reports that global data centre electricity use rose 17% in 2025. Its central outlook almost doubles consumption from 485 TWh in 2025 to 950 TWh in 2030, while consumption by AI-focused facilities triples. This is not a forecast of job counts, but it is a physical signal that infrastructure is expanding. [Source: IEA, Key Questions on Energy and AI](https://www.iea.org/reports/key-questions-on-energy-and-ai/executive-summary).

The software layer is mature too. In CNCF's 2025 survey, 82% of container users reported running Kubernetes in production, and 66% of organizations hosting generative models used Kubernetes for at least some inference workloads. The respondents come from the cloud native ecosystem and do not represent every company, but the result shows that container operations are no longer merely experimental. [Source: CNCF Annual Cloud Native Survey](https://www.cncf.io/reports/the-cncf-annual-cloud-native-survey/).

## What is changing in work

Employers surveyed by the World Economic Forum ranked networks and cybersecurity as the second fastest-growing skill group through 2030, after AI and big data. [Source: Future of Jobs Report 2025, PDF](https://reports.weforum.org/docs/WEF_Future_of_Jobs_Report_2025.pdf). That finding must not be turned into a promise that every job titled DevOps will grow.

US labour projections illustrate the shift. For 2025–2035, BLS projects about 8% growth for network architects, 21% for information security analysts, and 10% for software developers, QA analysts, and testers. It projects a 4% decline for traditional network and systems administrators because some routine work is automated, moved to DevOps-focused developers, or outsourced to service providers. These are US figures, not a Kazakhstan forecast, but they explain why this course focuses on code, automation, security, and reliability instead of memorising commands. [Sources: BLS on network architects](https://www.bls.gov/ooh/computer-and-information-technology/computer-network-architects.htm), [information security](https://www.bls.gov/ooh/computer-and-information-technology/information-security-analysts.htm), and [systems administration](https://www.bls.gov/ooh/computer-and-information-technology/network-and-computer-systems-administrators.htm).

AI can draft YAML or suggest a command. It does not take responsibility for downtime, lost data, a cloud bill, or an exposed secret. An engineer still has to state testable requirements, understand access boundaries, read metrics, choose a rollback, and prove recovery. You may use AI as an assistant in this course, but verify every proposal with a command, test, or observable result.

## Why this matters in Kazakhstan

Kazakhstan is turning strategic interest into physical infrastructure. In 2025 it launched the national Alem.Cloud cluster on NVIDIA H200 accelerators; the Ministry of Artificial Intelligence and Digital Development reported performance of about two exaflops at FP8. [Source: the Ministry's 2025 results](https://www.gov.kz/memleket/entities/maidd/press/news/details/1136177?lang=en).

In 2026 the Government described digital infrastructure as a priority and announced agreements around the Data Center Valley in Ekibastuz. Its first 50 MW facility is under construction with commissioning announced for June 2027; the wider campus is being considered with capacity up to 1 GW. These are plans and construction milestones, so this course does not describe future capacity as already operating. [Source: Government project update](https://primeminister.kz/ru/news/olzhas-bektenov-provel-soveshchanie-o-hode-realizacii-proekta-dolina-codov-v-pavlodarskoy-oblasti-31770), [KAZAKH INVEST on the proposed 1 GW scale](https://invest.gov.kz/ru/amp/news/40804/).

This direction may increase demand across several distinct fields: power and cooling, telecommunications, network engineering, physical security, hardware operations, cloud and platform engineering, SRE, DevOps, and cybersecurity. This course covers software operations. It does not train a data centre electrician or facilities designer, and it cannot replace experience on a real production site.

## Why Go appears so often in cloud native infrastructure

Modern infrastructure has no single language. The Linux kernel is mostly C, automation often uses Python and shell, interfaces use JavaScript and TypeScript, and performance-sensitive components may use Rust, C++, or Java. Yet **a substantial part of the cloud native control plane is built in Go**. Docker Engine, Kubernetes, Prometheus, Caddy, OpenTofu, and many Kubernetes operators are developed in it. Their own repositories make that claim inspectable: [Moby/Docker Engine](https://github.com/moby/moby), [Kubernetes](https://github.com/kubernetes/kubernetes), [Prometheus](https://github.com/prometheus/prometheus), [Caddy](https://github.com/caddyserver/caddy), and [OpenTofu](https://github.com/opentofu/opentofu).

Go fits this field for practical reasons: a project can produce one convenient binary, cross-compilation is built into the toolchain, goroutines suit concurrent network services, the standard library has strong HTTP support, and explicit errors suit systems where failure must remain visible. This does not make Go the best language for every task. It explains why reading a small Go service helps an infrastructure engineer understand the tools around it.

CloudLab is therefore written in Go and built as a static binary for a minimal container. **You do not need to know Go before starting**: the application is supplied and this course operates its lifecycle. To understand its code more deeply, take the free [Go: From Zero to Your Own Blog](https://shanraq.org/course/go?lang=en) course.

## The roles you can try

- A **cloud engineer** manages cloud networks, compute, access, and cost.
- A **DevOps engineer** shortens the path from a change to a safe release through automation and shared responsibility.
- A **platform engineer** builds an internal platform and a paved path for developers.
- An **SRE** manages reliability through indicators, service objectives, error budgets, and engineering incident response.
- A **data centre engineer** also works with the physical site, power, cooling, racks, and links; this is a separate field.

Job titles overlap. The course outcome is therefore evidence rather than a label: a container image, CI, a secured server, HTTPS, infrastructure as code, Kubernetes manifests, metrics, a backup, and a recovery record.

## Why there is no diploma for its own sake

Shanraq does not issue a useless certificate merely because a lesson tab was opened. A polished PDF cannot restore a failed service, find a leaked secret, or recover data. Our free courses produce work that can be shown and repeated: a repository with a clear history, a reproducible release, a working endpoint, automated checks, an incident record, and data restored from a backup.

The strongest evidence is your desire to understand and your sustained effort turned into a working result. At the end you have an engineering tool and portfolio that can solve real problems, support a technical interview, and help you build a livelihood. Twenty-four lessons do not guarantee income automatically; practice, responsibility, and the ability to solve another person's problem reliably create value. The course helps you prove those abilities without relying on a decorative credential.

## How we will proceed: the support-signal method

This course is for someone who may be meeting DNS, TCP, TLS, container, or pipeline for the first time. We do not treat these words as obvious. A new word first connects to a familiar image, then receives a precise definition and a place on a short support map, and only then appears in a command.

We adapt principles from Viktor Shatalov's teaching system to independent online learning:

1. **See the whole first.** Four supports at the start keep individual commands attached to a larger meaning.
2. **Familiar image, then term.** DNS becomes a directory, an IP a building address, and a port an office number. We then state the technical limit of the analogy.
3. **Support signal.** A few words and symbols encode the lesson logic. You expand them in your own words instead of memorising a slogan.
4. **Repeated recall.** The next lesson briefly returns to the previous signal. At the end, you draw the chain from memory and explain it aloud.
5. **Predict before observing.** Before a command, you predict its result. Comparing prediction with evidence turns a label into understanding.
6. **Frequent feedback.** Small checks reveal the exact link that is missing. An error points to the block to revisit; it is not punishment or a permanent gate.
7. **Open perspective.** Every ending shows why the next lesson matters and which complete project is taking shape in your hands.

Read each lesson in two passes. On the first, follow the story and map without trying to memorise everything. On the second, work through the terms, prediction, command, and verification. Then hide the map and reconstruct its four supports on paper. If you cannot explain a link, revisit that block rather than blindly rereading the whole page.

## Zero setup: if you have never opened a terminal

A terminal is an ordinary application with a text window. Open **Terminal** on macOS or Ubuntu. On Windows, install **WSL2 with Ubuntu** and enable WSL integration in Docker Desktop. This is your safe learning workshop. Do not rent a server or pay for anything at this stage.

Install Git, Docker, Go, and curl using the official instructions for your system. Open a new terminal window and run these commands one at a time:

```shell
git --version
docker version
go version
curl --version
```

Each command should print a version. `docker version` should show both Client and Server information. If Server is unavailable, start Docker Desktop or Docker Engine first. Then obtain the learning project:

```shell
git clone https://github.com/DauletBai/shanraq.org.git
cd shanraq.org/course/cloud-devops-lab
pwd
```

The output of `pwd` must end in `course/cloud-devops-lab`. Do not run later lesson commands until it does; correct the working directory first. Official installation pages: [Git](https://git-scm.com/downloads), [Docker](https://docs.docker.com/get-started/get-docker/), [Go](https://go.dev/doc/install), and [WSL](https://learn.microsoft.com/windows/wsl/install).

## The project and course rules

We operate a small service called CloudLab. It saves notes to disk and exposes health and readiness endpoints, a request metric, and structured logs. The application is deliberately simple so the deployment path and failure behaviour stay visible. The first part is local and free. For public DNS and TLS, you may rent a small virtual machine from any provider for a short time; no brand is required. Delete it after the exercise.

The 24 lessons cover Linux, networking, Docker, Compose, a registry, CI/CD, secrets, cloud access and cost, SSH, firewall rules, HTTPS, OpenTofu, Ansible, Kubernetes, observability, backups, and a recovery drill. Read commands before running them, never commit real tokens, and do not experiment on infrastructure you do not own or on a production system.

The course is free in Kazakh, Russian, and English. Reading needs no account. Practice needs a computer, a terminal, Git, and permission to install Docker. You can run Linux in a virtual machine or WSL2. The checked reference files live in [`course/cloud-devops-lab`](https://github.com/DauletBai/shanraq.org/tree/main/course/cloud-devops-lab).

[Open the course contents](/course/cloud-devops?lang=en)
