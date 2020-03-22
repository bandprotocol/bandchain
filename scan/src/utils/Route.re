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

let regexes = [
  "^blocks$" |> Js.Re.fromString,
  "^scripts$" |> Js.Re.fromString,
  "^sources$" |> Js.Re.fromString,
  "^tx$" |> Js.Re.fromString,
  "^txs$" |> Js.Re.fromString,
  "^validators$" |> Js.Re.fromString,
  "B([0-9]+)$" |> Js.Re.fromString,
  "D([0-9]+)$" |> Js.Re.fromString,
  "O([0-9]+)$" |> Js.Re.fromString,
  "R([0-9]+)$" |> Js.Re.fromString,
  "^bandvaloper([0-9a-z]+)" |> Js.Re.fromString,
  "^band([0-9a-z]+)" |> Js.Re.fromString,
];

let fromUrl = (url: ReasonReactRouter.url) => {
  switch (
    regexes
    ->Belt_List.keepMap(regex =>
        url.path->Belt_List.get(0)->Belt_Option.getWithDefault("")->Js.Re.exec_(regex, _)
      )
    ->Belt.List.head
  ) {
  | Some(result) =>
    let x = result->Js.Re.captures->Belt_Array.keepMap(Js.toOption);
    Js.Console.log2(x, url.hash);
    NotFound;
  | None => NotFound
  };

  switch (url.path, url.hash) {
  | (["sources"], _) => DataSourceHomePage
  | (["scripts"], _) => OracleScriptHomePage
  | (["txs"], _) => TxHomePage
  | (["validators"], _) => ValidatorHomePage
  | (["blocks"], _) => BlockHomePage
  | (["tx", txHash], _) => TxIndexPage(Hash.fromHex(txHash))
  | ([address], "delegations") =>
    AccountIndexPage(address |> Address.fromBech32, AccountDelegations)
  | ([address], "delegators") => ValidatorIndexPage(address |> Address.fromBech32, Delegators)
  | ([address], "reports") => ValidatorIndexPage(address |> Address.fromBech32, Reports)
  | ([path], tab) =>
    switch ((path |> Js.String.split(""))->Belt_List.fromArray, tab) {
    | (["B", ...rest], _) => BlockIndexPage(rest |> chars_to_int)
    | (["O", ...rest], "code") => OracleScriptIndexPage(rest |> chars_to_int, OracleScriptCode)
    | (["O", ...rest], "requests") =>
      OracleScriptIndexPage(rest |> chars_to_int, OracleScriptRequests)
    | (["O", ...rest], "revisions") =>
      OracleScriptIndexPage(rest |> chars_to_int, OracleScriptRevisions)
    | (["O", ...rest], _) => OracleScriptIndexPage(rest |> chars_to_int, OracleScriptExecute)
    | (["D", ...rest], "code") => DataSourceIndexPage(rest |> chars_to_int, DataSourceCode)
    | (["D", ...rest], "requests") =>
      DataSourceIndexPage(rest |> chars_to_int, DataSourceRequests)
    | (["D", ...rest], "revisions") =>
      DataSourceIndexPage(rest |> chars_to_int, DataSourceRevisions)
    | (["D", ...rest], _) => DataSourceIndexPage(rest |> chars_to_int, DataSourceExecute)
    | (["R", ...rest], "proof") => RequestIndexPage(rest |> chars_to_int, RequestProof)
    | (["R", ...rest], _) => RequestIndexPage(rest |> chars_to_int, RequestReportStatus)
    | (["b", "a", "n", "d", "v", "a", "l", "o", "p", "e", "r", ..._], _) =>
      ValidatorIndexPage(path |> Address.fromBech32, ProposedBlocks)
    | (["b", "a", "n", "d", ..._], _) =>
      AccountIndexPage(path |> Address.fromBech32, AccountTransactions)
    | _ => NotFound
    }
  | ([], _) => HomePage
  | (_, _) => NotFound
  };
};

let toString =
  fun
  | DataSourceHomePage => "/sources"
  | DataSourceIndexPage(dataSourceID, DataSourceExecute) => {j|/D$dataSourceID|j}
  | DataSourceIndexPage(dataSourceID, DataSourceCode) => {j|/D$dataSourceID#code|j}
  | DataSourceIndexPage(dataSourceID, DataSourceRequests) => {j|/D$dataSourceID#requests|j}
  | DataSourceIndexPage(dataSourceID, DataSourceRevisions) => {j|/D$dataSourceID#revisions|j}
  | OracleScriptHomePage => "/scripts"
  | OracleScriptIndexPage(oracleScriptID, OracleScriptExecute) => {j|/O$oracleScriptID|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptCode) => {j|/O$oracleScriptID#code|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptRequests) => {j|/O$oracleScriptID#requests|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptRevisions) => {j|/O$oracleScriptID#revisions|j}
  | TxHomePage => "/txs"
  | TxIndexPage(txHash) => {j|/tx/$txHash|j}
  | ValidatorHomePage => "/validators"
  | BlockHomePage => "/blocks"
  | BlockIndexPage(height) => {j|/B$height|j}
  | RequestIndexPage(reqID, RequestReportStatus) => {j|/R$reqID|j}
  | RequestIndexPage(reqID, RequestProof) => {j|/R$reqID#proof|j}
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
