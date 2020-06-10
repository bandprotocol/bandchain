#!/usr/bin/env python3

import json
import urllib.request
import sys

YAHOO_URL = "https://finance.yahoo.com/quote/{}"


def make_json_request(url):
    return urllib.request.urlopen(url).read()


def main(symbol):
    raw = make_json_request(YAHOO_URL.format(symbol)).decode()
    data = "".join(raw.split("\n")).split(
        "root.App.main = ")[1].split(";}(this)")[0]
    return json.loads(data)["context"]["dispatcher"]["stores"]["QuoteSummaryStore"]["price"][
        "regularMarketPrice"
    ]["raw"]


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
