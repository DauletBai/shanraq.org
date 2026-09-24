import sys
from datetime import date

from engine import answer

if len(sys.argv) == 3:
    show_trace = False
    raw_date = sys.argv[1]
    question = sys.argv[2]
elif len(sys.argv) == 4 and sys.argv[1] == "--trace":
    show_trace = True
    raw_date = sys.argv[2]
    question = sys.argv[3]
else:
    print('Қолдану: python3 main.py [--trace] YYYY-MM-DD "Сұрақ"')
    sys.exit(2)
try:
    on_date = date.fromisoformat(raw_date)
except ValueError:
    print("Қате күн: YYYY-MM-DD түрінде жазыңыз")
    sys.exit(2)
result = answer(question, on_date)
print(result["text"])
if show_trace:
    for line in result["trace"]:
        print("Із:", line)
