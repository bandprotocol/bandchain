// TODO: No need for this because it doesn't make sense
let rev: string => string = [%bs.raw
  {|
function(x) {
  return x.match(/.{1,2}/g).reverse().join("")
}
  |}
];

// TODO: No need for this because it doesn't make sense
let toSolanaCommand: string => string = [%bs.raw
  {|
  function(x) {
    return ("02" + Number((x.length)>>1).toString(16).padStart(8,"0").match(/.{1,2}/g).reverse().join("")) + x;
  }
  |}
];

// TODO: Replace this mocking format by the real
let createProofFromResult = result => {
  // Just a mock padding
  "00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000580"
  // Validator pubkey is [2u8; 32]
  // TODO: Use real validator
  ++ "0202020202020202020202020202020202020202020202020202020202020202"
  ++ (result |> JsBuffer.toHex(~with0x=false) |> rev)
  |> toSolanaCommand
  |> JsBuffer.fromHex;
};
