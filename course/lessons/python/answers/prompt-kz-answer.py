"""48-сабақ тапсырмасы: сұранысты тексеріп, тәуелсіз көшірмесін жасау."""

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
    {"role": "system", "content": "Жаңа факт қоспаңыз."},
    {"role": "user", "content": "<data>2025: 11.4 %</data>"}],
    "stream": False,
    "options": {"temperature": 0.0, "seed": 42, "num_predict": 120}}
settings = validate_request(request)
copy = json.loads(json.dumps(request))
copy["options"]["temperature"] = 0.7
print("модель:", settings["model"])
print("рөлдер:", settings["roles"])
print("temperature:", settings["temperature"])
print("seed:", settings["seed"])
print("шек:", settings["limit"])
print("көшірме тәуелсіз:", request["options"]["temperature"] == 0.0)
differences = sum(request["options"][key] != copy["options"][key]
                  for key in request["options"])
print("бір өріс өзгерді:", differences == 1)
