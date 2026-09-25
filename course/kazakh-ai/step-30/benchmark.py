import json
import platform
import statistics
import subprocess
import sys
import time
import tracemalloc
from datetime import date
from pathlib import Path

from checks import club_forms
from engine import answer


HERE = Path(__file__).resolve().parent
MODEL_FILES = ["main.py", "engine.py", "checks.py", "examples.json", "facts.json"]


def percentile(values, fraction):
    ordered = sorted(values)
    position = round((len(ordered) - 1) * fraction)
    return ordered[position]


def process_peak_bytes():
    try:
        import resource
    except ImportError:
        return None
    peak = resource.getrusage(resource.RUSAGE_SELF).ru_maxrss
    if sys.platform == "darwin":
        return peak
    return peak * 1024


with open(HERE / "quality_cases.json", encoding="utf-8") as file:
    cases = json.load(file)

cold_times = []
for repeat in range(10):
    start = time.perf_counter_ns()
    subprocess.run(
        [sys.executable, "main.py", "2026-09-24", "Шахмат қашан?"],
        cwd=HERE, check=True, capture_output=True, text=True,
    )
    cold_times.append(time.perf_counter_ns() - start)

for case in cases:
    answer(case["question"], date.fromisoformat(case["date"]))

times = []
throughputs = []
for batch in range(5):
    start_batch = time.perf_counter_ns()
    for repeat in range(2000):
        case = cases[repeat % len(cases)]
        start = time.perf_counter_ns()
        answer(case["question"], date.fromisoformat(case["date"]))
        times.append(time.perf_counter_ns() - start)
    batch_ns = time.perf_counter_ns() - start_batch
    throughputs.append(2000 / (batch_ns / 1_000_000_000))

tracemalloc.start()
for repeat in range(2000):
    case = cases[repeat % len(cases)]
    answer(case["question"], date.fromisoformat(case["date"]))
current_bytes, traced_peak_bytes = tracemalloc.get_traced_memory()
tracemalloc.stop()

correct = 0
correct_answers = 0
expected_answers = 0
correct_refusals = 0
expected_refusals = 0
unsupported_confident_answers = 0
for case in cases:
    actual = answer(case["question"], date.fromisoformat(case["date"]))["text"]
    matches = actual == case["expected"]
    correct += matches
    if case["kind"] == "answer":
        expected_answers += 1
        correct_answers += matches
    else:
        expected_refusals += 1
        correct_refusals += matches
        if not matches and not actual.startswith(("білмеймін:", "нақтылаңыз:", "тоқта:")):
            unsupported_confident_answers += 1

example_rows = json.loads((HERE / "examples.json").read_text(encoding="utf-8"))
fact_rows = json.loads((HERE / "facts.json").read_text(encoding="utf-8"))
report = {
    "measured_at": date.today().isoformat(),
    "machine": platform.platform(),
    "python": platform.python_version(),
    "architecture": "deterministic rules plus a count-based intent classifier",
    "trainable_neural_parameters": 0,
    "neural_context_window_tokens": None,
    "training_and_tuning_examples": len(example_rows),
    "supported_intent_types": len({row["label"] for row in example_rows
                                    if row["label"] != "белгісіз"}),
    "known_club_surface_forms": len(club_forms),
    "approved_fact_records": len(fact_rows),
    "model_code_and_data_bytes": sum((HERE / name).stat().st_size for name in MODEL_FILES),
    "cold_start_median_ms": round(statistics.median(cold_times) / 1_000_000, 3),
    "warm_request_median_ms": round(statistics.median(times) / 1_000_000, 3),
    "warm_request_p95_ms": round(percentile(times, 0.95) / 1_000_000, 3),
    "warm_throughput_requests_per_second": round(statistics.median(throughputs), 1),
    "python_traced_peak_bytes": traced_peak_bytes,
    "process_peak_rss_bytes": process_peak_bytes(),
    "quality_cases": len(cases),
    "scenario_accuracy": round(correct / len(cases), 4),
    "supported_answer_accuracy": round(correct_answers / expected_answers, 4),
    "correct_refusal_rate": round(correct_refusals / expected_refusals, 4),
    "unsupported_confident_answer_rate": round(
        unsupported_confident_answers / expected_refusals, 4
    ),
    "external_api_requests": 0,
    "direct_api_fee_usd": 0,
    "energy_per_request_joules": None,
    "energy_note": "not measured; requires a hardware power meter",
}

if "--json" in sys.argv:
    print(json.dumps(report, ensure_ascii=False, indent=2))
else:
    for key, value in report.items():
        print(f"{key}: {value}")
