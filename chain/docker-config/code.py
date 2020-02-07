import os
import subprocess

import json

CURRENT_PATH = os.path.dirname(os.path.realpath(__file__))
EXAMPLE_DIR = CURRENT_PATH + "/../../owasm/example/"

scripts = []
for script in os.listdir(EXAMPLE_DIR):
    if script.startswith("."):
        continue
    subprocess.call([CURRENT_PATH+"/build_wasm.sh %s"% EXAMPLE_DIR+script],shell=True)
    a = subprocess.check_output(["xxd", "-p", "-c100000000000",EXAMPLE_DIR+script+"/pkg/"+script+"_bg.wasm"])
    scripts.append({
        "code": a.decode()[:-1],
        "name": script,
        "owner": "band13zmknvkq2sj920spz90g4r9zjan8g584x8qalj",
    })

with open(CURRENT_PATH+"/temp.json","w") as f:
    json.dump(scripts,f)
