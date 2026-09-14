"""Lesson 48 exercise: validate and independently copy a request."""

import json


def validate_request(payload):
    messages, options = payload["messages"], payload["options"]
    assert payload["model"] == "gemma4"
    assert [item["role"] for item in messages] == ["system", "user"]
    assert payload["stream"] is False
    assert isinstance(options["temperature"], (int, float))
    assert isinstance(options["seed"], int)
    assert options["num_predict"] > 0
    assert "<data>" in messages[1]["content"] and "</data>" in messages[1]["content"]
    return {"model": payload["model"], "roles": "system,user",
            "temperature": options["temperature"], "seed": options["seed"],
            "limit": options["num_predict"]}


request = {"model": "gemma4", "messages": [
    {"role": "system", "content": "Add no facts."},
    {"role": "user", "content": "<data>2025: 11.4 %</data>"}],
    "stream": False,
    "options": {"temperature": 0.0, "seed": 42, "num_predict": 120}}
settings = validate_request(request)
copy = json.loads(json.dumps(request))
copy["options"]["temperature"] = 0.7
print("model:", settings["model"])
print("roles:", settings["roles"])
print("temperature:", settings["temperature"])
print("seed:", settings["seed"])
print("limit:", settings["limit"])
print("copy is independent:", request["options"]["temperature"] == 0.0)
differences = sum(request["options"][key] != copy["options"][key]
                  for key in request["options"])
print("one field changed:", differences == 1)
