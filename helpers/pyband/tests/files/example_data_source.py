#!/usr/bin/env python3

import sys
import json
import requests

HEADERS = {"Content-Type": "application/json"}

SYMBOLS_WITH_WEIGHTS = [
    ["BTC", 0.004909683753206196],
    ["ETH", 0.13310406316230874],
    ["BNB", 0.14811355573986523],
    ["EOS", 8.953653071560828],
    ["XTZ", 9.094515911727951],
    ["ATOM", 1.7587624274324516],
    ["YFI", 0.0010102624194177228],
    ["BCH", 0.02036657498252182],
    ["TRX", 205.4816659023913],
    ["HT", 0.6857067994908471],
    ["OKB", 0.5980410214630887],
    ["ZIL", 87.2055371354679],
    ["CRO", 69.95450022346525],
    ["XLM", 23.053273306237845],
    ["ADA", 7.179120672821946],
    ["DOT", 0.27893252164401366],
    ["SNX", 0.5011997921585194],
    ["ALGO", 9.192726225525655],
    ["OMG", 2.1996607471011007],
    ["COMP", 0.023306593928026486],
    ["NEO", 0.19421422582454298],
    ["VET", 171.95165365836635],
    ["BAT", 13.126425414588228],
    ["AAVE", 0.019770547127068678],
    ["LINK", 0.2764242681323155],
    ["THETA", 2.176840201683421],
    ["LEO", 3.5708266638194557],
    ["KNC", 4.3561169623294544],
    ["LTC", 0.0413220565629011],
    ["ZRX", 5.640930140168719],
    ["XRP", 16.427604258162788],
    ["MKR", 0.003358561719070429],
    ["CELO", 1.7828863706544964],
    ["FTT", 0.5340837373118201],
    ["BTT", 12062.322952053712],
    ["REN", 13.314915914819395],
    ["QTUM", 2.8585514788368362],
    ["DGB", 62.22667438038412],
    ["ONT", 3.754157138230833],
    ["ICX", 2.2741324674258427],
    ["KSM", 0.01489803116735643],
    ["UMA", 0.16281776279745758],
    ["XEM", 5.280971599079332],
    ["MIOTA", 2.828955811666102],
]


def main(url):
    symbols, weights = zip(*SYMBOLS_WITH_WEIGHTS)
    payload = {"symbols": symbols, "min_count": 10, "ask_count": 16}
    result = requests.request("POST", url, headers=HEADERS, json=payload).json()
    prices = [float(px["px"]) / float(px["multiplier"]) for px in result["result"]]
    acc = 0
    for (p, w) in zip(prices, weights):
        acc += float(p) * w
    return acc


if __name__ == "__main__":
    try:
        print(main(*sys.argv[1:]))
    except Exception as e:
        print(str(e), file=sys.stderr)
        sys.exit(1)
