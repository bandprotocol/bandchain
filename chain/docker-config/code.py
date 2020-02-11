import os
import subprocess

import json

CURRENT_PATH = os.path.dirname(os.path.realpath(__file__))
EXAMPLE_DIR = CURRENT_PATH + "/../../owasm/examples/"

scripts = []
for script in os.listdir(EXAMPLE_DIR):
    if script.startswith(".") or script.startswith("yahoo"):
        continue
    subprocess.call(
        [CURRENT_PATH + "/build_wasm.sh %s" % EXAMPLE_DIR + script], shell=True
    )
    a = subprocess.check_output(
        [
            "xxd",
            "-p",
            "-c100000000000",
            EXAMPLE_DIR + script + "/pkg/" + script + "_bg.wasm",
        ]
    )
    scripts.append(
        {
            "code": a.decode()[:-1],
            "name": script,
            "owner": "band1r8u4m0wefrdyp7sdde7drtrt94qczggkv38jrk",
        }
    )

with open(CURRENT_PATH + "/temp.json", "w") as f:
    json.dump(scripts, f)
