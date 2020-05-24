from pyobi import PyObi

def main():
    obi = PyObi(""" {
  symbol: string,
  px: u64,
  in: {
    a: u8,
    b: u8
  }
} / string
""")
    print("Hello, World")
    encoded = obi.encode_input({
        "symbol": "BTC",
        "px": 9000,
        "in": {
            "a": 1,
            "b": 2,
        }
    })
    print(obi.decode_input(encoded))


if __name__ == '__main__':
    main()
