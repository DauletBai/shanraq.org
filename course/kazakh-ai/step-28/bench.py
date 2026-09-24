import statistics
import subprocess
import sys
import time
import tracemalloc
from pathlib import Path

from checks import classify

question = "Шахмат қашан?"
cold = []
for repeat in range(5):
    start = time.perf_counter_ns()
    subprocess.run([sys.executable, "-c", "from checks import classify; classify('Шахмат қашан?')"],
                   check=True, capture_output=True, text=True)
    cold.append(time.perf_counter_ns() - start)
classify(question)
times = []
tracemalloc.start()
for repeat in range(200):
    start = time.perf_counter_ns()
    classify(question)
    times.append(time.perf_counter_ns() - start)
current, peak = tracemalloc.get_traced_memory()
tracemalloc.stop()
files = Path("checks.py").stat().st_size + Path("examples.json").stat().st_size
print("Сұрау саны:", len(times))
print("Жаңа іске қосу, мс:", round(statistics.median(cold) / 1000000, 3))
print("Ортаңғы уақыт, мс:", round(statistics.median(times) / 1000000, 3))
print("Ең көп бақыланған бөлу, байт:", peak)
print("Код пен дерек, байт:", files)
print("Тікелей API төлемі: 0; құрылғы мен еңбек құны есептелмеді")
