
import json
import os
from pathlib import Path
import sys
import re

# python convert.py folder/text.txt

hyphens_regex = re.compile(r"-+\n")
result_regex = re.compile(r"(\d+)\) (.*)")

folder = Path(sys.argv[1]).parent
tests = {
    "tests": []
}

with open(sys.argv[1]) as io:
    lines = io.readlines()
    for i, line in enumerate(lines):
        result = result_regex.match(line)
        if hyphens_regex.match(line):
            if len(tests["tests"]) != 0 and len(tests["tests"][-1]["results"]) == 1:
                tests["tests"].pop()
            # Create ne
            filename = folder.joinpath(lines[i-1].strip() + ".nes").as_posix()
            tests["tests"].append({
                "rom": str(filename),
                "frames": 60,
                "output": sys.argv[3] if len(sys.argv) == 4 else "0x00F0",
                "results": [{
                    "code": 1,
                    "pass": True,
                    "message": "Test passed",
                }],
            })
        elif result and result.group(1) != 1:
            tests["tests"][-1]["results"].append({
                "code": int(result.group(1)),
                "pass": False,
                "message": result.group(2),
            })

with open(sys.argv[2], mode="w") as io:
    json.dump(tests, io, indent=2)
