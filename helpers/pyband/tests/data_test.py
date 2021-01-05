import pytest
from pyband.data import Coin
from pyband.exceptions import InsufficientCoinError


def test_coin_success():
    coin = Coin(amount=5000, denom="uband")

    assert coin.as_json() == {"amount": "5000", "denom": "uband"}

    assert coin.validate() == True


def test_coin_amount_fail():
    coin = Coin(amount=-5000, denom="uband")

    with pytest.raises(InsufficientCoinError, match="Expect amount more than 0"):
        coin.validate()


def test_coin_denom_fail():
    coin = Coin(amount=5000, denom="")

    with pytest.raises(ValueError, match="Expect denom"):
        coin.validate()
