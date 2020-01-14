type script_tab_t =
  | ScriptTransactions
  | ScriptCode
  | ScriptIntegration;

type request_tab_t =
  | RequestReportStatus
  | RequestProof;

type t =
  | NotFound
  | HomePage
  | ScriptHomePage
  | ScriptIndexPage(string, script_tab_t)
  | TxHomePage
  | TxIndexPage(string)
  | BlockHomePage
  | BlockIndexPage(int)
  | RequestIndexPage(string, request_tab_t);

let fromUrl = (url: ReasonReactRouter.url) =>
  switch (url.path, url.hash) {
  | (["scripts"], _) => ScriptHomePage
  | (["script", codeHash], "code") => ScriptIndexPage(codeHash, ScriptCode)
  | (["script", codeHash], "integration") => ScriptIndexPage(codeHash, ScriptIntegration)
  | (["script", codeHash], _) => ScriptIndexPage(codeHash, ScriptTransactions)
  | (["txs"], _) => TxHomePage
  | (["tx", txHash], _) => TxIndexPage(txHash)
  | (["blocks"], _) => BlockHomePage
  | (["block", blockHeight], _) =>
    let blockHeightIntOpt = blockHeight |> int_of_string_opt;
    BlockIndexPage(blockHeightIntOpt->Belt_Option.getWithDefault(0));
  | (["request", reqID], "proof") => RequestIndexPage(reqID, RequestProof)
  | (["request", reqID], _) => RequestIndexPage(reqID, RequestReportStatus)
  | ([], "") => HomePage
  | (_, _) => NotFound
  };

let toString =
  fun
  | ScriptHomePage => "/scripts"
  | ScriptIndexPage(codeHash, ScriptTransactions) => {j|/script/$codeHash|j}
  | ScriptIndexPage(codeHash, ScriptCode) => {j|/script/$codeHash#code|j}
  | ScriptIndexPage(codeHash, ScriptIntegration) => {j|/script/$codeHash#integration|j}
  | TxHomePage => "/txs"
  | TxIndexPage(txHash) => {j|/tx/$txHash|j}
  | BlockHomePage => "/blocks"
  | BlockIndexPage(height) => {j|/block/$height|j}
  | RequestIndexPage(reqID, RequestReportStatus) => {j|/request/$reqID|j}
  | RequestIndexPage(reqID, RequestProof) => {j|/request/$reqID#proof|j}
  | HomePage
  | NotFound => "/";

let redirect = (route: t) => ReasonReactRouter.push(route |> toString);
