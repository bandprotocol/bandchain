type t = {
  bandChain: CosmosJS.t,
  mnemonic: string,
  privKey: JsBuffer.t,
};

let getAddressAndPubKey = x => {
  (
    x.bandChain |> CosmosJS.getAddress(_, x.mnemonic) |> Address.fromBech32,
    Secp256k1.publicKeyCreate(x.privKey, true) |> JsBuffer.toBase64 |> PubKey.fromBase64,
  );
};

// TODO: sign message
