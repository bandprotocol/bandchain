FROM owasm-server:latest

WORKDIR /root/ui/

# install yarn
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | apt-key add - \
  && echo "deb https://dl.yarnpkg.com/debian/ stable main" \
  |  tee /etc/apt/sources.list.d/yarn.list \
  && apt-get update && apt-get install -y yarn

COPY . /root/ui

RUN yarn install

CMD ["yarn", "start"]
