import re
import sys
import os


def generate_substitute(content, header, number, component, changeType, description):
    return re.sub(
        f"### {header}\n",
        f"### {header}\n\n- ({changeType}) [\#{number}](https://https://github.com/bandprotocol/bandchain/pull/{number}) {description}",
        content,
        count=1,
    )


def append_changelog(prNumber, changedComponent, changeType, prDescription):
    read = open("DEMOCHANGELOG.md", "r+")
    content = read.read()
    read.close()
    with open("DEMOCHANGELOG.md", "w") as f:
        content_new = ""
        changedComponent = changedComponent.lower()
        if changedComponent == "chain":
            f.write(
                generate_substitute(
                    content, "Chain", prNumber, changedComponent, changeType, prDescription
                )
            )
        elif changedComponent == "scan":
            f.write(
                generate_substitute(
                    content, "Scan", prNumber, changedComponent, changeType, prDescription
                )
            )
        elif changedComponent == "bridges":
            f.write(
                generate_substitute(
                    content, "Bridges", prNumber, changedComponent, changeType, prDescription
                )
            )
        elif changedComponent == "owasm":
            f.write(
                generate_substitute(
                    content, "Owasm", prNumber, changedComponent, changeType, prDescription
                )
            )
        elif changedComponent == "obi":
            print("obi")
            f.write(
                generate_substitute(
                    content,
                    "Oracle Binary Encoding \(OBI\)",
                    prNumber,
                    changedComponent,
                    changeType,
                    prDescription,
                )
            )
        elif changedComponent == "misc":
            f.write(
                generate_substitute(
                    content, "MISC", prNumber, changedComponent, changeType, prDescription
                )
            )


if __name__ == "__main__":
    prNumber = os.getenv("PR_NUMBER")
    changedComponent = os.getenv("PR_BODY").split("\n")[2].split(":")[0]
    prChangeType = (os.getenv("PR_BODY").split("\n")[1][:-1],)
    prDescription = os.getenv("PR_BODY").split("\n")[2].split(":")[1]
    append_changelog(
        prNumber, changedComponent, prChangeType, prDescription,
    )
    ## append_changelog(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4])
