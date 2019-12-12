# The tag here should match the Meteor version of your app, per .meteor/release
FROM geoffreybooth/meteor-base:1.8.1

ENV METEOR_ALLOW_SUPERUSER true

# Copy source into container
COPY . /root/explorer

EXPOSE 12000

WORKDIR /root/explorer

RUN meteor npm install

CMD meteor run -p 0.0.0.0:12000 --settings default_settings.json --no-release-check
