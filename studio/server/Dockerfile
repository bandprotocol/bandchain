FROM node:10.17.0-buster

ENV PATH="$PATH:/root/.cargo/bin"

# install rust & friends
RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs > rust_install.sh \
  && chmod +x rust_install.sh \
  && sh rust_install.sh -y --default-toolchain nightly \
  && rustup target add wasm32-unknown-unknown --toolchain \
  && cargo install wasm-gc \
  && cargo install wasm-bindgen-cli \
  && rustup component add rustfmt-preview --toolchain nightly \
  && curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh

COPY . /root/server

WORKDIR /root/server

CMD yarn && yarn start
