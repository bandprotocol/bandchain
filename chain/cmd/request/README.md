## Request script

Send new request to chain with example wasm file

## How to run

1. Set 2 environment varaibles
   - `NODE_URI` Node endpoint to subscribe event and send transaction (Default is -> tcp://localhost:26657)
   - `PRIVATE_KEY` A hexstring represents validator address (Default is -> `eedda7a96ad35758f2ffc404d6ccd7be913f149a530c70e95e2e3ee7a952a877` (64 digits))
2. Run (There are 3 tx mode)
   1. Store code will send tx to store code (hardcode)
   ```
   go run main.go store
   ```
   2. Send coin to other account
   ```
    go run main.go send_token
   ```
   3. Request new data
   ```
   go run main.go request [symbol BTC/ETH]
   ```
