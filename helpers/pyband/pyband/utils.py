from datetime import datetime, timezone


def parse_datetime(raw_datetime) -> int:
    splitted_datetime = raw_datetime.split(".")
    formatted_datetime = splitted_datetime[0] + ".0Z"
    parsed_datetime = datetime.strptime(formatted_datetime, "%Y-%m-%dT%H:%M:%S.%fZ")
    return int(parsed_datetime.replace(tzinfo=timezone.utc).timestamp())
