let wenchang =
  fun
  | 1 => "internal error"
  | 2 => "tx parse error"
  | 3 => "invalid sequence"
  | 4 => "unauthorized"
  | 5 => "insufficient funds"
  | 6 => "unknown request"
  | 7 => "invalid address"
  | 8 => "invalid pubkey"
  | 9 => "unknown address"
  | 10 => "insufficient coins"
  | 11 => "invalid coins"
  | 12 => "out of gas"
  | 13 => "memo too large"
  | 14 => "insufficient fee"
  | 15 => "maximum numer of signatures exceeded"
  | 16 => "no signatures supplied"
  | _ => "an error occurred";

let guanyu =
  fun
  | 1 => "internal error"
  | 2 => "tx parse error"
  | 3 => "invalid sequence"
  | 4 => "unauthorized"
  | 5 => "insufficient funds"
  | 6 => "unknown request"
  | 7 => "invalid address"
  | 8 => "invalid pubkey"
  | 9 => "unknown address"
  | 10 => "invalid coins"
  | 11 => "out of gas"
  | 12 => "memo too large"
  | 13 => "insufficient fee"
  | 14 => "maximum numer of signatures exceeded"
  | 15 => "no signatures supplied"
  | 16 => "failed to marshal JSON bytes"
  | 17 => "failed to unmarshal JSON bytes"
  | 18 => "invalid request"
  | 19 => "tx already in mempool"
  | 20 => "mempool is full"
  | 21 => "tx too large"
  | 22 => "key not found"
  | 23 => "invalid account password"
  | 24 => "tx intended signer does not match the given signer"
  | 25 => "invalid gas adjustment"
  | 26 => "invalid height"
  | 27 => "invalid version"
  | 28 => "invalid chain-id"
  | 29 => "invalid type"
  | 111222 => "panic"
  | _ => "an error occurred";

let parse = code => {
  exception WrongNetwork(string);
  switch (Env.network) {
  | "GUANYU38"
  | "GUANYU" => guanyu(code)
  | "WENCHANG" => wenchang(code)
  | _ => raise(WrongNetwork("Incorrect or unspecified NETWORK environment variable"))
  };
};
