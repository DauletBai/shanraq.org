"""Задание урока 47: проверить форму ответа локальной модели."""


def answer_text(data):
    """Вернуть непустой текст или назвать одну ясную ошибку контракта."""
    message = data.get("message")
    if not isinstance(message, dict):
        raise ValueError("нет текста модели")
    content = message.get("content")
    if not isinstance(content, str) or not content.strip():
        raise ValueError("нет текста модели")
    return content.strip()


samples = [
    {"message": {"content": "  ответ: 42  "}},
    {"message": {"content": "   "}},
    {"message": {"content": None}},
    {},
]

for number, sample in enumerate(samples, 1):
    try:
        print(f"{number}: {answer_text(sample)}")
    except ValueError as error:
        print(f"{number}: ошибка: {error}")
print("проверено ответов:", len(samples))
