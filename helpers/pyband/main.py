from pyband import Client, PyObi


def main():
    c = Client("http://guanyu-devnet.bandchain.org/rest")
    req_info = c.get_latest_request(67, bytes.fromhex("0000000000000064"), 4, 4)
    oracle_script = c.get_oracle_script(67)
    obi = PyObi(oracle_script.schema)
    print(obi.decode_output(req_info.result.ResponsePacketData.result))


if __name__ == "__main__":
    main()
