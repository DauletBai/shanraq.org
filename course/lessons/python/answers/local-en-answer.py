"""Lesson 47 exercise: check the shape of a local model response."""


def answer_text(data):
    """Return non-empty text or name one clear contract error."""
    message = data.get("message")
    if not isinstance(message, dict):
        raise ValueError("no model text")
    content = message.get("content")
    if not isinstance(content, str) or not content.strip():
        raise ValueError("no model text")
    return content.strip()


samples = [
    {"message": {"content": "  answer: 42  "}},
    {"message": {"content": "   "}},
    {"message": {"content": None}},
    {},
]

for number, sample in enumerate(samples, 1):
    try:
        print(f"{number}: {answer_text(sample)}")
    except ValueError as error:
        print(f"{number}: error: {error}")
print("responses checked:", len(samples))
