type t = {
  bandChain: CosmosJS.t,
  mnemonic: string,
  privKey: JsBuffer.t,
};

let create = mnemonic => {
  // Just use arbitrary rpcUrl, chainID beacuase they didn't use in CosmosJS
  let bandChain = CosmosJS.network("rpcUrl", "chainID");
  bandChain->CosmosJS.setPath("m/44'/494'/0'/0/0");
  bandChain->CosmosJS.setBech32MainPrefix("band");
  let privKey = bandChain |> CosmosJS.getECPairPriv(_, mnemonic);
  {bandChain, mnemonic, privKey};
};

let getAddressAndPubKey = x => {
  (
    x.bandChain |> CosmosJS.getAddress(_, x.mnemonic) |> Address.fromBech32,
    Secp256k1.publicKeyCreate(x.privKey, true) |> JsBuffer.toBase64 |> PubKey.fromBase64,
  );
};

// TODO: sign message
