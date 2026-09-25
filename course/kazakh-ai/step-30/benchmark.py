import json
import multiprocessing
import os
import platform
import statistics
import subprocess
import sys
import time
import tracemalloc
from concurrent.futures import ProcessPoolExecutor
from datetime import date
from pathlib import Path


HERE = Path(__file__).resolve().parent
os.chdir(HERE)

from checks import club_forms
from engine import answer


MODEL_FILES = ["main.py", "engine.py", "checks.py", "examples.json", "facts.json"]
CASES = json.loads((HERE / "quality_cases.json").read_text(encoding="utf-8"))


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


def run_requests(repeats):
    start = time.perf_counter_ns()
    for repeat in range(repeats):
        case = CASES[repeat % len(CASES)]
        answer(case["question"], date.fromisoformat(case["date"]))
    elapsed = time.perf_counter_ns() - start
    return {"pid": os.getpid(), "requests": repeats, "elapsed_ns": elapsed,
            "peak_rss_bytes": process_peak_bytes()}


def worker_counts(available):
    maximum = max(1, available or 1)
    return sorted({1, min(2, maximum), min(4, maximum), min(8, maximum)})


def measure_scaling(available):
    rows = []
    for workers in worker_counts(available):
        batch_rates = []
        batch_rss = []
        with ProcessPoolExecutor(max_workers=workers) as pool:
            list(pool.map(run_requests, [200] * workers))
            for batch in range(3):
                start = time.perf_counter_ns()
                results = list(pool.map(run_requests, [10_000] * workers))
                elapsed = time.perf_counter_ns() - start
                batch_rates.append(sum(row["requests"] for row in results)
                                   / (elapsed / 1_000_000_000))
                peaks = [row["peak_rss_bytes"] for row in results]
                if all(value is not None for value in peaks):
                    batch_rss.append(sum(peaks))
        rows.append({
            "worker_processes": workers,
            "parallel_requests": workers,
            "throughput_requests_per_second": round(statistics.median(batch_rates), 1),
            "summed_worker_peak_rss_bytes": (round(statistics.median(batch_rss))
                                             if batch_rss else None),
        })
    baseline = rows[0]["throughput_requests_per_second"]
    for row in rows:
        speedup = row["throughput_requests_per_second"] / baseline
        row["speedup_vs_one_worker"] = round(speedup, 3)
        row["parallel_efficiency"] = round(speedup / row["worker_processes"], 3)
    return rows


def build_report():
    cold_times = []
    for repeat in range(10):
        start = time.perf_counter_ns()
        subprocess.run(
            [sys.executable, "main.py", "2026-09-24", "Шахмат қашан?"],
            cwd=HERE, check=True, capture_output=True, text=True,
        )
        cold_times.append(time.perf_counter_ns() - start)

    for case in CASES:
        answer(case["question"], date.fromisoformat(case["date"]))

    times = []
    throughputs = []
    for batch in range(5):
        start_batch = time.perf_counter_ns()
        for repeat in range(2000):
            case = CASES[repeat % len(CASES)]
            start = time.perf_counter_ns()
            answer(case["question"], date.fromisoformat(case["date"]))
            times.append(time.perf_counter_ns() - start)
        batch_ns = time.perf_counter_ns() - start_batch
        throughputs.append(2000 / (batch_ns / 1_000_000_000))

    tracemalloc.start()
    for repeat in range(2000):
        case = CASES[repeat % len(CASES)]
        answer(case["question"], date.fromisoformat(case["date"]))
    unused_current_bytes, traced_peak_bytes = tracemalloc.get_traced_memory()
    tracemalloc.stop()

    correct = 0
    correct_answers = 0
    expected_answers = 0
    correct_refusals = 0
    expected_refusals = 0
    unsupported_confident_answers = 0
    for case in CASES:
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
    available = os.cpu_count() or 1
    return {
        "measured_at": date.today().isoformat(),
        "machine": platform.platform(),
        "python": platform.python_version(),
        "available_logical_cpus": available,
        "worker_processes": 1,
        "parallel_requests": 1,
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
        "parallel_scaling": measure_scaling(available),
        "quality_cases": len(CASES),
        "scenario_accuracy": round(correct / len(CASES), 4),
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


def main():
    report = build_report()
    if "--json" in sys.argv:
        print(json.dumps(report, ensure_ascii=False, indent=2))
    else:
        for key, value in report.items():
            print(f"{key}: {value}")


if __name__ == "__main__":
    multiprocessing.freeze_support()
    main()
