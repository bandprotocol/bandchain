type t =
  | NotFound
  | HomePage
  | ScriptHomePage
  | ScriptIndexPage(string, string)
  | TxHomePage
  | TxIndexPage(string, string)
  | BlockHomePage
  | BlockIndexPage(string, string);

let fromUrl = (url: ReasonReactRouter.url) =>
  switch (url.path, url.hash) {
  | ([], _) => HomePage
  | (["script"], _) => ScriptHomePage
  | (["script", codeHash], hashtag) => ScriptIndexPage(codeHash, hashtag)
  | (["tx"], _) => TxHomePage
  | (["tx", txHash], hashtag) => TxIndexPage(txHash, hashtag)
  | (["block"], _) => BlockHomePage
  | (["block", blockHeight], hashtag) => BlockIndexPage(blockHeight, hashtag)
  | (_, _) => NotFound
  };

let toString =
  fun
  | ScriptHomePage => "/script"
  | ScriptIndexPage(codeHash, "") => {j|/script/$codeHash|j}
  | ScriptIndexPage(codeHash, hashtag) => {j|/script/$codeHash#$hashtag|j}
  | TxHomePage => "/tx"
  | TxIndexPage(txHash, "") => {j|/tx/$txHash|j}
  | TxIndexPage(txHash, hashtag) => {j|/tx/$txHash#$hashtag|j}
  | BlockHomePage => "/block"
  | BlockIndexPage(height, "") => {j|/block/$height|j}
  | BlockIndexPage(height, hashtag) => {j|/block/$height#$hashtag|j}
  | _ => "/";
