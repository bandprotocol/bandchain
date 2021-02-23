import base64
import pytest
from pyband.data import Coin
from pyband.message import MsgSend
from pyband.wallet import Address
from pyband.exceptions import InsufficientCoinError, NegativeIntegerError


def test_msg_send_creation_success():
    msg_send = MsgSend(
        to_address=Address.from_acc_bech32("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"),
        from_address=Address.from_acc_bech32("band1acavyhqpxmz6jt390xze705620q23e4tx4r5he"),
        amount=[Coin(amount=1000000, denom="uband")],
    )

    assert msg_send.validate() == True

    assert msg_send.as_json() == {
        "type": "cosmos-sdk/MsgSend",
        "value": {
            "to_address": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
            "from_address": "band1acavyhqpxmz6jt390xze705620q23e4tx4r5he",
            "amount": [{"amount": "1000000", "denom": "uband"}],
        },
    }

    assert msg_send.get_sender().to_acc_bech32() == "band1acavyhqpxmz6jt390xze705620q23e4tx4r5he"


def test_msg_send_nocoin_fail():
    msg_send = MsgSend(
        to_address=Address.from_acc_bech32("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"),
        from_address=Address.from_acc_bech32("band1acavyhqpxmz6jt390xze705620q23e4tx4r5he"),
        amount=[],
    )

    with pytest.raises(InsufficientCoinError, match="Expect at least 1 coin"):
        msg_send.validate()


def test_msg_send_validate_coin_fail():
    msg_send = MsgSend(
        to_address=Address.from_acc_bech32("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"),
        from_address=Address.from_acc_bech32("band1acavyhqpxmz6jt390xze705620q23e4tx4r5he"),
        amount=[Coin(amount=-1000, denom="uband")],
    )

    with pytest.raises(NegativeIntegerError, match="Expect amount more than 0"):
        msg_send.validate()
