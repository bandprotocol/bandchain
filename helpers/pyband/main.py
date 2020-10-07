from pyband import Client, PyObi
from pyband.wallet import PrivateKey


def main():
    c = Client("http://poa-api.bandchain.org")
    # req_info = c.get_latest_request(
    #     6, bytes.fromhex("000000045041584700000003555344000000003b9aca00"), 4, 4
    # )
    oracle_script = c.get_oracle_script(99)
    # x = c.get_account("band1c8sgs7j5rmjjtnv284senhcw9d0g7m5qxlf9js")
    print(x)
    # obi = PyObi(oracle_script.schema)
    # print(obi.decode_output(req_info.result.response_packet_data.result))

    # _, priv = PrivateKey.generate()
    # print(priv.to_pubkey().to_acc_bech32(), priv.to_pubkey().to_address().to_acc_bech32())


if __name__ == "__main__":
    main()
