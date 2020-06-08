import re
import sys
import os


def generate_substitute(content, header, number, component, changeType, description):
    return re.sub(
        f"### {header}\n",
        f"### {header}\n\n- {changeType} [\#{number}](https://https://github.com/bandprotocol/bandchain/pull/{number}) {description}",
        content,
        count=1,
    )


def append_changelog(number, component, changeType, description):
    read = open("DEMOCHANGELOG.md", "r+")
    content = read.read()
    read.close()
    with open("DEMOCHANGELOG.md", "w") as f:
        content_new = ""
        component = component.lower()
        if component == "chain":
            f.write(
                generate_substitute(
                    content, "Chain", number, component, changeType, description
                )
            )
        elif component == "scan":
            f.write(
                generate_substitute(
                    content, "Scan", number, component, changeType, description
                )
            )
        elif component == "bridges":
            f.write(
                generate_substitute(
                    content, "Bridges", number, component, changeType, description
                )
            )
        elif component == "owasm":
            f.write(
                generate_substitute(
                    content, "Owasm", number, component, changeType, description
                )
            )
        elif component == "obi":
            print("obi")
            f.write(
                generate_substitute(
                    content,
                    "Oracle Binary Encoding \(OBI\)",
                    number,
                    component,
                    changeType,
                    description,
                )
            )
        elif component == "misc":
            f.write(
                generate_substitute(
                    content, "MISC", number, component, changeType, description
                )
            )

if __name__ == "__main__":
    append_changelog(
        os.getenv("PR_NUMBER"),
        os.getenv("PR_BODY").split("\n")[2].split(":")[0],
        os.getenv("PR_BODY").split("\n")[1],
        os.getenv("PR_BODY").split("\n")[2].split(":")[1],
    )
    ## append_changelog(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4])
