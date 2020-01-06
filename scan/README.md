# Scan

## Installation

```
yarn install
yarn global add local-cors-proxy
```

## Running App

In 3 separate tabs:

```sh
lcp --proxyUrl http://d3n.bandprotocol.com:1317/ --proxyPartial '' # Proxy server
yarn bsb -make-world -w -ws _ # ReasonML compiler
yarn parcel index.html # Serve to localhost:1234
```
