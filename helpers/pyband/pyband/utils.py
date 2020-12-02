from datetime import datetime, timezone


def parse_epoch_time(raw_datetime: str) -> int:
    split_datetime = raw_datetime.split(".")
    formatted_datetime = split_datetime[0] + ".0Z"
    parsed_datetime = datetime.strptime(formatted_datetime, "%Y-%m-%dT%H:%M:%S.%fZ")
    return int(parsed_datetime.replace(tzinfo=timezone.utc).timestamp())
