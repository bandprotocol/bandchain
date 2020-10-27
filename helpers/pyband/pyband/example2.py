from pyband import Client
from pyband.wallet import PrivateKey
from pyband.message import MsgRequest
from pyband.transaction import Transaction

mnemonic = "s"


def main():
    privkey = PrivateKey.from_mnemonic(mnemonic)
    c = Client("https://d3n.bandprotocol.com/rest")

    addr = privkey.to_pubkey().to_address()
    acc = c.get_account(addr.to_acc_bech32())

    msgs = [
        MsgRequest(
            oracle_script_id=1,
            calldata=bytes.fromhex("000000034254430000000000000001"),
            ask_count=4,
            min_count=3,
            client_id="from_pyband",
            sender=addr,
        )
    ]

    t = Transaction.create(msgs)
    t = (
        t.with_account_num(acc.account_number)
        .with_sequence(acc.sequence)
        .with_chain_id("bandchain")
    )

    raw_tx = t.as_json(privkey.sign(t.get_sign_data()))

    result = c.send_tx(raw_tx, "block")
    print(result)


if __name__ == "__main__":
    main()
