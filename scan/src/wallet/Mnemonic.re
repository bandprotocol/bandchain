type t = {
  bandChain: BandWeb3.t,
  mnemonic: string,
  privKey: JsBuffer.t,
};

let getAddressAndPubKey = x => {
  (
    x.bandChain |> BandWeb3.getAddress(_, x.mnemonic) |> Address.fromBech32,
    BandWeb3.publicKeyCreate(x.privKey, true) |> JsBuffer.toBase64 |> PubKey.fromBase64,
  );
};

// TODO: sign message
