type data_source_tab_t =
  | DataSourceExecute
  | DataSourceCode
  | DataSourceRequests
  | DataSourceRevisions;

type oracle_script_tab_t =
  | OracleScriptExecute
  | OracleScriptCode
  | OracleScriptRequests
  | OracleScriptRevisions;

type request_tab_t =
  | RequestReportStatus
  | RequestProof;

type account_tab_t =
  | AccountTransactions
  | AccountDelegations;

type validator_tab_t =
  | ProposedBlocks
  | Delegators
  | Reports;

type t =
  | NotFound
  | HomePage
  | DataSourceHomePage
  | DataSourceIndexPage(int, data_source_tab_t)
  | OracleScriptHomePage
  | OracleScriptIndexPage(int, oracle_script_tab_t)
  | TxHomePage
  | TxIndexPage(Hash.t)
  | BlockHomePage
  | BlockIndexPage(int)
  | RequestIndexPage(int, request_tab_t)
  | AccountIndexPage(Address.t, account_tab_t)
  | ValidatorHomePage
  | ValidatorIndexPage(Address.t, validator_tab_t);

let chars_to_int = chars =>
  switch (chars->Belt_List.toArray->Js.Array.joinWith("", _)->int_of_string_opt) {
  | Some(id) => id
  | None => 0
  };

let id_to_int = id_str =>
  switch (id_str->Js.String.sliceToEnd(1)->int_of_string_opt) {
  | Some(id_) => id_
  | None => 0
  };

let fromUrl = (url: ReasonReactRouter.url) =>
  switch (url.path, url.hash) {
  | (["sources"], _) => DataSourceHomePage
  | (["data-source", dataSourceID], "code") =>
    DataSourceIndexPage(dataSourceID |> int_of_string, DataSourceCode)
  | (["data-source", dataSourceID], "requests") =>
    DataSourceIndexPage(dataSourceID |> int_of_string, DataSourceRequests)
  | (["data-source", dataSourceID], "revisions") =>
    DataSourceIndexPage(dataSourceID |> int_of_string, DataSourceRevisions)
  | (["data-source", dataSourceID], _) =>
    DataSourceIndexPage(dataSourceID |> int_of_string, DataSourceExecute)
  | (["scripts"], _) => OracleScriptHomePage
  | ([oracleScriptID], "code") =>
    OracleScriptIndexPage(oracleScriptID |> id_to_int, OracleScriptCode)
  | ([oracleScriptID], "requests") =>
    OracleScriptIndexPage(oracleScriptID |> id_to_int, OracleScriptRequests)
  | ([oracleScriptID], "revisions") =>
    OracleScriptIndexPage(oracleScriptID |> id_to_int, OracleScriptRevisions)
  | (["txs"], _) => TxHomePage
  | (["tx", txHash], _) => TxIndexPage(Hash.fromHex(txHash))
  | (["validators"], _) => ValidatorHomePage
  | (["blocks"], _) => BlockHomePage
  | (["block", blockHeight], _) =>
    let blockHeightIntOpt = blockHeight |> int_of_string_opt;
    BlockIndexPage(blockHeightIntOpt->Belt_Option.getWithDefault(0));
  | (["request", reqID], "proof") => RequestIndexPage(reqID |> int_of_string, RequestProof)
  | (["request", reqID], _) => RequestIndexPage(reqID |> int_of_string, RequestReportStatus)
  | ([address], "delegations") =>
    AccountIndexPage(address |> Address.fromBech32, AccountDelegations)
  | ([address], "delegators") => ValidatorIndexPage(address |> Address.fromBech32, Delegators)
  | ([address], "reports") => ValidatorIndexPage(address |> Address.fromBech32, Reports)
  | ([path], _) =>
    let chars = (path |> Js.String.split(""))->Belt_List.fromArray;
    switch (chars) {
    | ["O", ...rest] => OracleScriptIndexPage(rest |> chars_to_int, OracleScriptExecute)
    | ["b", "a", "n", "d", "v", "a", "l", "o", "p", "e", "r", ..._] =>
      ValidatorIndexPage(path |> Address.fromBech32, ProposedBlocks)
    | ["b", "a", "n", "d", ..._] =>
      AccountIndexPage(path |> Address.fromBech32, AccountTransactions)
    | _ => NotFound
    };
  | ([], _) => HomePage
  | (_, _) => NotFound
  };

let toString =
  fun
  | DataSourceHomePage => "/sources"
  | DataSourceIndexPage(dataSourceID, DataSourceExecute) => {j|/data-source/$dataSourceID|j}
  | DataSourceIndexPage(dataSourceID, DataSourceCode) => {j|/data-source/$dataSourceID#code|j}
  | DataSourceIndexPage(dataSourceID, DataSourceRequests) => {j|/data-source/$dataSourceID#requests|j}
  | DataSourceIndexPage(dataSourceID, DataSourceRevisions) => {j|/data-source/$dataSourceID#revisions|j}
  | OracleScriptHomePage => "/scripts"
  | OracleScriptIndexPage(oracleScriptID, OracleScriptExecute) => {j|/O$oracleScriptID|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptCode) => {j|/O$oracleScriptID#code|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptRequests) => {j|/O$oracleScriptID#requests|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptRevisions) => {j|/O$oracleScriptID#revisions|j}
  | TxHomePage => "/txs"
  | TxIndexPage(txHash) => {j|/tx/$txHash|j}
  | ValidatorHomePage => "/validators"
  | BlockHomePage => "/blocks"
  | BlockIndexPage(height) => {j|/block/$height|j}
  | RequestIndexPage(reqID, RequestReportStatus) => {j|/request/$reqID|j}
  | RequestIndexPage(reqID, RequestProof) => {j|/request/$reqID#proof|j}
  | AccountIndexPage(address, AccountTransactions) => {
      let addressBech32 = address |> Address.toBech32;
      {j|$addressBech32|j};
    }
  | AccountIndexPage(address, AccountDelegations) => {
      let addressBech32 = address |> Address.toBech32;
      {j|$addressBech32#delegations|j};
    }
  | ValidatorIndexPage(validatorAddress, Delegators) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|$validatorAddressBech32#delegators|j};
    }
  | ValidatorIndexPage(validatorAddress, Reports) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|$validatorAddressBech32#reports|j};
    }
  | ValidatorIndexPage(validatorAddress, ProposedBlocks) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|$validatorAddressBech32#proposed-blocks|j};
    }
  | HomePage => "/"
  | NotFound => "/notfound";

let redirect = (route: t) => ReasonReactRouter.push(route |> toString);

let search = (str: string) => {
  let len = str |> String.length;
  let capStr = str |> String.capitalize_ascii;

  (
    switch (str |> int_of_string_opt) {
    | Some(blockID) => Some(BlockIndexPage(blockID))
    | None =>
      if (str |> Js.String.startsWith("bandvaloper")) {
        Some(ValidatorIndexPage(str |> Address.fromBech32, ProposedBlocks));
      } else if (str |> Js.String.startsWith("band")) {
        Some(AccountIndexPage(str |> Address.fromBech32, AccountTransactions));
      } else if (capStr |> Js.String.startsWith("B")) {
        let%Opt blockID = str |> String.sub(_, 1, len - 1) |> int_of_string_opt;
        Some(BlockIndexPage(blockID));
      } else if (capStr |> Js.String.startsWith("D")) {
        let%Opt dataSourceID = str |> String.sub(_, 1, len - 1) |> int_of_string_opt;
        Some(DataSourceIndexPage(dataSourceID, DataSourceExecute));
      } else if (capStr |> Js.String.startsWith("R")) {
        let%Opt requestID = str |> String.sub(_, 1, len - 1) |> int_of_string_opt;
        Some(RequestIndexPage(requestID, RequestReportStatus));
      } else if (capStr |> Js.String.startsWith("O")) {
        let%Opt oracleScriptID = str |> String.sub(_, 1, len - 1) |> int_of_string_opt;
        Some(OracleScriptIndexPage(oracleScriptID, OracleScriptExecute));
      } else if (len == 64 || str |> Js.String.startsWith("0x") && len == 66) {
        Some(TxIndexPage(str |> Hash.fromHex));
      } else {
        None;
      }
    }
  )
  |> Belt_Option.getWithDefault(_, NotFound);
};
