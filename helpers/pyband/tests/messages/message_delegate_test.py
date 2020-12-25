import pytest
from pyband.data import Coin
from pyband.message import MsgDelegate
from pyband.wallet import Address
from pyband.error import InsufficientCoinError


def test_msg_delegate_creation_success():
    msg_delegate = MsgDelegate(
        delegator_address=Address.from_acc_bech32("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"),
        validator_address=Address.from_val_bech32("bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"),
        amount=Coin(amount=1000000, denom="uband"),
    )

    assert msg_delegate.validate() == True

    assert msg_delegate.as_json() == {
        "type": "cosmos-sdk/MsgDelegate",
        "value": {
            "delegator_address": "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte",
            "validator_address": "bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst",
            "amount": {"amount": "1000000", "denom": "uband"},
        },
    }

    assert msg_delegate.get_sender().to_acc_bech32() == "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"


def test_msg_delegate_coin_fail():
    msg_delegate = MsgDelegate(
        delegator_address=Address.from_acc_bech32("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"),
        validator_address=Address.from_val_bech32("bandvaloper1j9vk75jjty02elhwqqjehaspfslaem8pr20qst"),
        amount=Coin(amount=-1000000, denom="uband"),
    )

    with pytest.raises(InsufficientCoinError, match="Expect amount more than 0"):
        msg_delegate.amount.validate()
