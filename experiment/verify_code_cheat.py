PLACES = {
    "alphavantage": "6fd319bd571dcb0828dc2b5317df806bda1f70bdf08e3c3fb2046d5fd9a8982d",
    "binance_price": "6b7be61b150aec5eb853afb3b53e41438959554580d31259a1095e51645bcd28",
    "bitcoin_info": "ebfefd555d5b7b3e5bd7d9c6a9f396cbf22697edd1687a54d3425f8b4e4eb480",
    "coingecko_price": "d10a439e08c31902cfed8d337e68f297c84665c36c1025caa60ce6bbd2767715",
    "coingecko_volume": "afa2fea33afa77df2e61206df0228ba73923d68c7e54970d6ae25b6c12ea2770",
    "cryptocompare_price": "06077e219b9bceb3ca90240d5f9d383e418e9916a9da02fce7aa441d279af2d4",
    "cryptocompare_volume": "d3497a03c318476612524fb6d8f6c761e411fdf67966b08d157b28a83980cf7c",
    "ethereum_gas_price": "31eea00625c5fe2ea80a7fde255469a2f57534dfadc715ed30228e7392a6529c",
    "flight_verification": "5bc4ba4f391f83ebe066387c43cc1576bcf80602671bd4831e008973dfea8bb9",
    "random_u64": "e7944e5e24dc856dcb6d9926460926ec10b9b66cf44b664f9971b5a5e9255989",
    "weather_info": "1aa0bb8d7584921357cd2c6847f7a0cd896e3ae99a6b95908f22a693d9e6c5c1",
}

import requests
import json


def doit(folder, hs):
    print("YO", folder, hs.upper())
    print(
        requests.post(
            "http://rpc.alpha-debug.bandchain.org:8082/upload",
            json={
                "code": json.dumps(
                    [
                        {
                            "name": "Project/Cargo.toml",
                            "content": open(folder + "/Cargo.toml").read(),
                        },
                        {
                            "name": "Project/src/lib.rs",
                            "content": open(folder + "/src/lib.rs").read(),
                        },
                        {
                            "name": "Project/src/logic.rs",
                            "content": open(folder + "/src/logic.rs").read(),
                        },
                    ]
                ),
                "hash": hs.upper(),
            },
        ).text
    )


for k, v in PLACES.items():
    doit(k, v)
