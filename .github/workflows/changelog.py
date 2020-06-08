import re
import sys


def generate_substitute(content, header, number, component, changeType, description):
    return re.sub(
        f"### {header}\n",
        f"### {header}\n\n- ({changeType}) [\#{number}](https://https://github.com/bandprotocol/bandchain/pull/{number}) {description}",
        content,
        count=1,
    )


def append_changelog(number, component, changeType, description):
    read = open("../../DEMOCHANGELOG.md", "r+")
    content = read.read()
    read.close()
    with open("../../DEMOCHANGELOG.md", "w") as f:
        content_new = ""
        component = component.split()[0][:-1].lower()
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

def temp(body):
    print(body.splitlines())

if __name__ == "__main__":
    temp(sys.argv[1])
    ## append_changelog(sys.argv[1], sys.argv[2], sys.argv[3], sys.argv[4])
