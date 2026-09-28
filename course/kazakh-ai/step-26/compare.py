import subprocess
import sys
from pathlib import Path

word = "мектептерімізде"
ours = {"root": "мектеп", "parts": ["тер", "іміз", "де"]}
print("Оқу үлгісі:", word, "→", ours["root"], ours["parts"])
if len(sys.argv) == 1:
    print("adam_fst: қосымша салыстыру іске қосылмады")
else:
    binary = Path(sys.argv[1]).expanduser().resolve()
    if not binary.is_file():
        print("adam_fst: бағдарлама файлы табылмады")
    else:
        run = subprocess.run([str(binary), "analyse", word], cwd=binary.parents[2],
                             text=True, capture_output=True, timeout=10, check=False)
        if run.returncode != 0:
            print("adam_fst: іске қосу қатесі")
        else:
            print("adam_fst:", run.stdout.strip())
