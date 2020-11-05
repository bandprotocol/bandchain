from pyband.utils import parse_datetime
from datetime import datetime


def test_parse_datetime_to_epoch():
    raw_datetime = "2020-11-05T09:15:18.445494105Z"
    print(parse_datetime(raw_datetime))
    assert parse_datetime(raw_datetime) == 16045425181


def test_parse_epoch_to_datetime():
    raw_datetime = "2020-11-05T09:15:18.445494105Z"
    assert datetime.fromtimestamp(float(parse_datetime(raw_datetime))) == datetime(
        2020, 11, 5, 9, 15, 18
    )
