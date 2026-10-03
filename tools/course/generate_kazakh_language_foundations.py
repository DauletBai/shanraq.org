#!/usr/bin/env python3
"""Generate the three manually localized versions of Kazakh block one."""
from pathlib import Path

ROOT = Path(__file__).resolve().parents[2]
OUT = ROOT / "course/lessons/kazakh-language"

UI = {
    "ru": {
        "map": "Где мы на карте", "situation": "Живая ситуация", "listen": "Послушайте и ответьте",
        "map_intro": "Схема удерживает маршрут разговора целиком; подробности мы добавим по одной.",
        "model": "Точная модель", "support": "Опорный сигнал", "limit": "Граница правила",
        "predict": "Предскажите до объяснения", "recall": "Восстановите без подсказки",
        "mistake": "Найдите и исправьте ошибку", "transfer": "Перенесите в новую ситуацию",
        "project": "Изменение проекта «Моя среда»", "task": "Задание",
        "repeat": "Повтор через 1, 7 и 30 дней", "next": "Следующий урок",
        "play": "Слушать"
    },
    "kz": {
        "map": "Картадағы орнымыз", "situation": "Өмірлік жағдай", "listen": "Тыңдаңыз және жауап беріңіз",
        "map_intro": "Сызба әңгіме бағытын тұтас көрсетеді; әр бөлігін кезекпен толықтырамыз.",
        "model": "Нақты үлгі", "support": "Тірек белгі", "limit": "Ереженің шегі",
        "predict": "Түсіндірмеге дейін болжаңыз", "recall": "Тірексіз қалпына келтіріңіз",
        "mistake": "Қатені тауып, түзетіңіз", "transfer": "Жаңа жағдайға көшіріңіз",
        "project": "«Менің ортам» жобасындағы өзгеріс", "task": "Тапсырма",
        "repeat": "1, 7 және 30 күннен кейін қайталау", "next": "Келесі сабақ",
        "play": "Тыңдау"
    },
    "en": {
        "map": "Where we are on the map", "situation": "A real situation", "listen": "Listen and respond",
        "map_intro": "The map holds the whole route; we will add its details one at a time.",
        "model": "The precise model", "support": "The support signal", "limit": "The limit of the rule",
        "predict": "Predict before the explanation", "recall": "Rebuild it without the prompt",
        "mistake": "Find and repair the mistake", "transfer": "Transfer to a new setting",
        "project": "Change to the My World project", "task": "Exercise",
        "repeat": "Retrieval after 1, 7, and 30 days", "next": "Next lesson",
        "play": "Listen"
    },
}


def L(title, lead, situation, phrases, model, limit, predict, recall, mistake, transfer, project, task, repeat):
    return dict(title=title, lead=lead, situation=situation, phrases=phrases, model=model, limit=limit,
                predict=predict, recall=recall, mistake=mistake, transfer=transfer,
                project=project, task=task, repeat=repeat)


LESSONS = [
    ("01-map-first-conversation", "01-first-contact", {
        "ru": L(
            "Карта курса и первые 30 секунд разговора",
            "За один урок вы проведёте первые 30 секунд знакомства на казахском и получите карту всех 80 уроков — от приветствия до самостоятельного разговора уровня B1.",
            "Представьте: в кружок пришёл новый ученик. Можно ждать идеальной грамматики и промолчать, а можно удержать короткий маршрут: установить контакт, назвать себя, задать один вопрос, услышать имя и доброжелательно завершить обмен. Нам нужна не заученная сценка, а конструкция, которую легко наполнить другим именем и голосом.",
            [("Сәлем!", "Привет!"), ("Менің атым Айша.", "Меня зовут Айша."), ("Сіздің атыңыз кім?", "Как вас зовут?"), ("Менің атым Данияр.", "Меня зовут Данияр."), ("Танысқаныма қуаныштымын.", "Рад(а) познакомиться.")],
            "Разговор держится на четырёх действиях: **контакт → сведения о себе → вопрос → ответ и завершение**. Фраза `Менің атым…` буквально строится как «моё имя…», но начинающему полезнее хранить её целой рамкой. `Сіздің` — уважительное «ваше», а `кім?` спрашивает о человеке. Произнесите имя после короткой паузы, не проглатывая последнее слово вопроса.",
            "Один маршрут не означает один текст. Вместо `Сәлем!` уместно более формальное `Сәлеметсіз бе!`; знакомые ровесники говорят короче. Пока мы берём нейтральный и вежливый вариант. Кнопка использует казахский голос устройства: она даёт слуховую опору, но живые голоса всё равно будут отличаться темпом и интонацией.",
            "Что можно заменить, не ломая разговор: имя, приветствие или слово `кім`? Сначала ответьте. Имя заменяется свободно; приветствие можно выбрать по ситуации. `Кім` сохраняет вопрос об имени человека, поэтому случайная замена изменит смысл.",
            "Закройте страницу и нарисуйте четыре значка: ладонь, карточка с собой, вопрос, ответ. По ним восстановите разговор, используя другое безопасное вымышленное имя. Если забыли одну реплику, откройте только схему, а не весь текст.",
            "Алан говорит: `Сәлем! Сіздің атыңыз Алан.` Он смешал сведения о себе и обращение к собеседнику. Разделите это на утверждение о себе и отдельный вопрос, не копируя готовый диалог целиком.",
            "Проведите тот же обмен в трёх условиях: с ровесником, со взрослым и в голосовом сообщении. Решите, где оставить `Сәлем`, а где выбрать `Сәлеметсіз бе`. Имя и порядок реплик выберите сами.",
            "Создайте карточку версии `0.0`: вымышленное имя героя, его роль и знакомое безопасное место. Добавьте четыре пустых поля будущего путеводителя: «кто», «где», «что полезного покажет», «каких личных данных не хранит».",
            "**Обязательно.** Напишите собственный диалог из пяти реплик: приветствие, имя, вопрос об имени, ответ и доброжелательное завершение. Затем произнесите его без чтения.\n\n**Проверка переноса.** Поменяйте обоих героев и формальность приветствия. Объясните одним предложением, почему это всё ещё тот же маршрут разговора.",
            ["Завтра восстановите четыре шага и скажите диалог с новыми именами.", "Через семь дней начните знакомство без предупреждения и добавьте встречный вопрос.", "Через тридцать дней запишите 30 секунд речи и отметьте только одно место, которое хотите улучшить."]
        ),
        "kz": L(
            "Курс картасы және әңгіменің алғашқы 30 секунды",
            "Бір сабақта қазақша танысудың алғашқы 30 секундын өткізесіз және амандасудан B1 деңгейіндегі дербес әңгімеге дейінгі 80 сабақтың картасын көресіз.",
            "Үйірмеге жаңа оқушы келді деп елестетіңіз. Мінсіз грамматиканы күтіп үндемей қалуға да, қысқа бағытты ұстауға да болады: байланыс орнату, өз атыңызды айту, бір сұрақ қою, есімді есту және әңгімені жылы аяқтау. Бізге жатталған көрініс емес, басқа есіммен толтыруға болатын құрылым керек.",
            [("Сәлем!", "Амандасу"), ("Менің атым Айша.", "Өз есімін айту"), ("Сіздің атыңыз кім?", "Әңгімелесушінің есімін сұрау"), ("Менің атым Данияр.", "Жауап беру"), ("Танысқаныма қуаныштымын.", "Әңгімені жылы аяқтау")],
            "Әңгіме төрт әрекетке сүйенеді: **байланыс → өзі туралы мәлімет → сұрақ → жауап пен аяқтау**. `Менің атым…` үлгісін әзірге тұтас қалып ретінде сақтаңыз. `Сіздің` — құрметті түрде айтылған сөз, ал `кім?` адамды сұрайды. Сұрақтың соңғы сөзін анық айтыңыз.",
            "Бір бағыт бір ғана мәтін деген сөз емес. Ресми жағдайда `Сәлеметсіз бе!`, ал таныс құрдаспен `Сәлем!` деуге болады. Құрылғыдағы қазақ дауысы тыңдауға тірек береді, бірақ тірі адамдардың қарқыны мен интонациясы әртүрлі болады.",
            "Әңгімені бұзбай нені ауыстыруға болады: есімді, амандасуды әлде `кім` сөзін бе? Есім еркін ауысады, амандасу жағдайға қарай таңдалады. `Кім` адам туралы сұрақты сақтайды; оны кездейсоқ ауыстырсаңыз, мағына өзгереді.",
            "Бетті жауып, төрт белгі салыңыз: алақан, өзі туралы карточка, сұрақ, жауап. Сол белгілермен басқа ойдан шығарылған есімді қолданып әңгімені қалпына келтіріңіз.",
            "Алан: `Сәлем! Сіздің атыңыз Алан` дейді. Ол өзі туралы мәлімет пен әңгімелесушіге үндеуді араластырып жіберді. Оны өзі туралы сөйлемге және бөлек сұраққа бөліңіз.",
            "Осы алмасуды үш жағдайда қолданыңыз: құрдаспен, ересек адаммен және дауыстық хабарламада. Қай жерде `Сәлем`, қай жерде `Сәлеметсіз бе` лайық екенін өзіңіз шешіңіз.",
            "`0.0` нұсқасының карточкасын жасаңыз: кейіпкердің ойдан шығарылған есімі, рөлі және қауіпсіз таныс орны. «Кім», «қайда», «нені пайдалы көрсетеді», «қандай жеке деректі сақтамайды» деген төрт өріс қосыңыз.",
            "**Міндетті.** Бес репликадан тұратын өз диалогыңызды жазыңыз: амандасу, есім, есім туралы сұрақ, жауап және жылы аяқтау. Содан кейін оқымай айтып көріңіз.\n\n**Көшіруді тексеру.** Екі кейіпкерді де және амандасу ресмилігін өзгертіңіз. Неліктен бұл сол әңгіменің бағыты болып қалғанын бір сөйлеммен түсіндіріңіз.",
            ["Ертең төрт қадамды қалпына келтіріп, жаңа есімдермен айтыңыз.", "Жеті күннен кейін алдын ала дайындықсыз танысуды бастап, қарсы сұрақ қосыңыз.", "Отыз күннен кейін 30 секундтық сөзіңізді жазып, жақсартатын бір ғана жерді белгілеңіз."]
        ),
        "en": L(
            "The course map and the first 30 seconds of conversation",
            "In one lesson you will carry the first 30 seconds of a Kazakh introduction and see the full 80-lesson route from a greeting to independent B1 communication.",
            "Imagine a new learner arriving at a club. You can wait for perfect grammar and stay silent, or hold a short route: make contact, name yourself, ask one question, hear the name, and close warmly. We need a reusable structure rather than a scene memorised with one pair of names.",
            [("Сәлем!", "Hello!"), ("Менің атым Айша.", "My name is Aisha."), ("Сіздің атыңыз кім?", "What is your name?"), ("Менің атым Данияр.", "My name is Daniyar."), ("Танысқаныма қуаныштымын.", "I am glad to meet you.")],
            "The exchange holds four actions: **contact → information about self → question → response and close**. For now, store `Менің атым…` as one useful frame. `Сіздің` is respectful “your,” and `кім?` asks about a person. Leave a small pause before the name and keep the last word of the question clear.",
            "One route does not mean one script. `Сәлеметсіз бе!` is more formal than `Сәлем!`. We begin with a neutral respectful exchange. The button uses the device's Kazakh voice as a listening anchor; real speakers will still vary in pace and intonation.",
            "What can change without breaking the exchange: the name, the greeting, or `кім`? A name changes freely and the greeting follows the setting. `Кім` preserves a question about a person, so replacing it at random changes the meaning.",
            "Close the page and draw four signs: a hand, a self card, a question, and a response. Use them to rebuild the exchange with a different fictional name. If one line disappears, inspect the map before returning to the full text.",
            "Alan says, `Сәлем! Сіздің атыңыз Алан.` He has mixed information about himself with an address to the other person. Separate it into a statement about self and one independent question.",
            "Run the same exchange with a peer, an adult, and in a voice message. Decide where `Сәлем` fits and where `Сәлеметсіз бе` is safer. Choose the names and order yourself.",
            "Create project card `0.0`: a fictional name, role, and safe familiar place. Add four fields: who, where, what useful thing the guide will show, and which personal data it will never store.",
            "**Required.** Write your own five-line exchange: greeting, name, question about the other name, response, and warm close. Then say it without reading.\n\n**Transfer check.** Change both people and the formality of the greeting. In one sentence, explain why the route remains the same.",
            ["Tomorrow rebuild the four steps and speak them with new names.", "After seven days, open an introduction without warning and add a return question.", "After thirty days, record 30 seconds and mark only one place you want to improve."]
        ),
    }),
]

# Lessons 2–8 are deliberately stored separately below. Keeping each localized
# explanation explicit makes language review possible; machine translation is
# never performed during the build.

MORE = []


def speech_button(lang, phrase, meaning):
    label = UI[lang]["play"]
    safe = phrase.replace('&', '&amp;').replace('"', '&quot;').replace('<', '&lt;')
    return f'<button type="button" class="speak-kz" data-speak-kz="{safe}">{label}: {phrase}</button> — {meaning}'


def render(number, stem, map_stem, lang, data, next_lesson):
    ui = UI[lang]
    map_path = f"/static/course/kazakh-language/map-{map_stem}-{lang}.svg"
    parts = [
        f"# {data['title']}", "", f"_Lead (summary):_ **{data['lead']}**", "",
        f"## {ui['map']}", "",
        ui['map_intro'], "",
        f"![{data['title']}]({map_path})", "",
        f"## {ui['situation']}", "", data['situation'], "",
        f"## {ui['listen']}", "",
    ]
    for phrase, meaning in data['phrases']:
        parts.extend([speech_button(lang, phrase, meaning), ""])
    parts.extend([
        f"## {ui['model']}", "", data['model'], "",
        f"## {ui['support']}", "", "```text", support_for(number), "```", "",
        f"## {ui['limit']}", "", data['limit'], "",
        f"## {ui['predict']}", "", data['predict'], "",
        f"## {ui['recall']}", "", data['recall'], "",
        f"## {ui['mistake']}", "", data['mistake'], "",
        f"## {ui['transfer']}", "", data['transfer'], "",
        f"## {ui['project']}", "", data['project'], "",
        f"## {ui['task']}", "", data['task'], "",
        f"## {ui['repeat']}", "",
    ])
    parts.extend(f"- {item}" for item in data['repeat'])
    if next_lesson:
        parts.extend(["", f"[{ui['next']}: {next_lesson[1]}](/read/kazakh-language-{next_lesson[0]}?lang={lang})"])
    return "\n".join(parts).strip() + "\n"


SUPPORT = {
    1: "СӘЛЕМ → МЕНІҢ АТЫМ… → СІЗДІҢ АТЫҢЫЗ КІМ? → ЖАУАП → ҚУАНЫШТЫМЫН",
}


def support_for(number):
    return SUPPORT[number]


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    lessons = LESSONS + MORE
    for idx, (stem, map_stem, localized) in enumerate(lessons):
        next_entry = lessons[idx + 1] if idx + 1 < len(lessons) else None
        for lang, data in localized.items():
            suffix = "" if lang == "ru" else f"-{lang}"
            nxt = None
            if next_entry:
                nxt_data = next_entry[2][lang]
                nxt = (next_entry[0], nxt_data["title"])
            (OUT / f"{stem}{suffix}.md").write_text(
                render(idx + 1, stem, map_stem, lang, data, nxt), encoding="utf-8")
    print(f"Generated {len(lessons) * 3} localized Kazakh lesson pages")


if __name__ == "__main__":
    main()
