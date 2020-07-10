FROM python:3-slim

COPY . .
RUN apt-get update && apt-get install -y libpq-dev gcc
RUN pip install -r requirements.txt

# TODO: Using environment variable for db and kafka bootstrap services
CMD python main.py sync --db postgres:postgrespassword@172.18.0.88:5432/postgres -s 172.18.0.31:9092
