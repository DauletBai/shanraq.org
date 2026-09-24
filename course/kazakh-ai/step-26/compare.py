import json
import subprocess
import sys
from pathlib import Path

word = "мектептерімізде"
ours = {"root": "мектеп", "parts": ["тер", "іміз", "де"]}
print("Оқу үлгісі:", word, "→", ours["root"], ours["parts"])
if len(sys.argv) == 1:
    print("qazaq-ir: қосымша салыстыру іске қосылмады")
else:
    binary = Path(sys.argv[1])
    if not binary.is_file():
        print("qazaq-ir: бағдарлама файлы табылмады")
    else:
        run = subprocess.run([str(binary), "analyze", "--format", "compact", word],
                             text=True, capture_output=True, timeout=10, check=False)
        if run.returncode != 0:
            print("qazaq-ir: іске қосу қатесі")
        else:
            result = json.loads(run.stdout)
            for token in result["tokens"]:
                print("qazaq-ir:", token["surface"], "→", token["root"],
                      token["analysis_status"])
