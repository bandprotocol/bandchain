type t;

[@bs.module "@cosmostation/cosmosjs"] external network: (string, string) => t = "network";

[@bs.send] external setPath: (t, string) => unit = "setPath";

[@bs.send] external setBech32MainPrefix: (t, string) => unit = "setBech32MainPrefix";

[@bs.send] external getAddress: (t, string) => string = "getAddress";

[@bs.send] external getECPairPriv: (t, string) => JsBuffer.t = "getECPairPriv";

let cosmos = network("bandchain", "http://d3n-debug.bandprotocol.com:1317");

cosmos->setPath("m/44'/494'/0'/0/0");
cosmos->setBech32MainPrefix("band");

Js.Console.log(
  cosmos->getAddress(
    "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comi",
  ),
);

Js.Console.log(
  cosmos->getECPairPriv(
    "smile stem oven genius cave resource better lunar nasty moon company ridge brass rather supply used horn three panic put venue analyst leader comi",
  ),
);
