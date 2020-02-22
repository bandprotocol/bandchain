import requests
import time
import subprocess
from datetime import datetime

DOCKER_ID = "333824c9a0db"


def main():
    while True:
        try:
            r = requests.get("http://localhost:26657")
            if r.status_code == 200:
                print("Last ok at", datetime.now(), end="\r")
        except KeyboardInterrupt:
            break
        except:
            print("\nNode downed")
            subprocess.run(["docker", "logs", "--tail", "20", DOCKER_ID])
            subprocess.run(["docker", "restart", "-t", "5", DOCKER_ID])
            print("Recover at", datetime.now())
        time.sleep(5)


if __name__ == "__main__":
    main()

