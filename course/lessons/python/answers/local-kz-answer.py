"""47-сабақ тапсырмасы: жергілікті модель жауабының пішінін тексеру."""


def answer_text(data):
    """Бос емес мәтінді қайтару немесе бір түсінікті келісімшарт қатесін атау."""
    message = data.get("message")
    if not isinstance(message, dict):
        raise ValueError("модель мәтіні жоқ")
    content = message.get("content")
    if not isinstance(content, str) or not content.strip():
        raise ValueError("модель мәтіні жоқ")
    return content.strip()


samples = [
    {"message": {"content": "  жауап: 42  "}},
    {"message": {"content": "   "}},
    {"message": {"content": None}},
    {},
]

for number, sample in enumerate(samples, 1):
    try:
        print(f"{number}: {answer_text(sample)}")
    except ValueError as error:
        print(f"{number}: қате: {error}")
print("тексерілген жауап:", len(samples))
