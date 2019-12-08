## Provider client

Subscribe event from blockchain, run prepare function, and send a transaction for reporting data to chain

## How to run

1. Set 2 environment varaibles
   - `NODE_URI` Node endpoint to subscribe event and send transaction (ex. tcp://localhost:26657)
   - `PRIVATE_KEY` A hexstring represents validator address (ex. 324342ab3.... (64 digits))
2. Run (without build)

```
go run main.go
```
