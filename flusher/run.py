from kafka import KafkaConsumer, TopicPartition


def main():
    consumer = KafkaConsumer("bandchain")
    partitions = consumer.partitions_for_topic("bandchain")
    if len(partitions) != 1:
        raise Exception("Only exact 1 partition is supported")
    consumer.seek(TopicPartition("bandchain", partitions.pop()), 25)
    for msg in consumer:
        print(msg)


if __name__ == "__main__":
    main()
