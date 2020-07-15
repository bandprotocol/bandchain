FROM ubuntu:18.04

COPY requirements.txt .

RUN apt-get update
RUN apt-get install -y python3-pip python3-dev

RUN pip3 install --upgrade pip
RUN pip3 install -r requirements.txt
