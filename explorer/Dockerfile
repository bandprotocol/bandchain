FROM band-meteor:0.4

COPY . /root/explorer

EXPOSE 12000

WORKDIR /root/explorer

RUN meteor npm install

CMD meteor run -p 0.0.0.0:12000 --settings default_settings.json --no-release-check --verbose ; meteor run -p 0.0.0.0:12000 --settings default_settings.json --no-release-check --verbose
