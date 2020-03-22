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
  ("^$" |> Js.Re.fromString, _ => HomePage),
  ("^blocks$" |> Js.Re.fromString, _ => BlockHomePage),
  ("^scripts$" |> Js.Re.fromString, _ => OracleScriptHomePage),
  ("^sources$" |> Js.Re.fromString, _ => DataSourceHomePage),
  ("^txs$" |> Js.Re.fromString, _ => TxHomePage),
  (
    "^tx$" |> Js.Re.fromString,
    keys => {
      switch (keys->Belt_List.get(1)) {
      | Some(hash) => TxIndexPage(Hash.fromHex(hash))
      | None => NotFound
      };
    },
  ),
  ("^validators$" |> Js.Re.fromString, _ => ValidatorHomePage),
  (
    "B([0-9]+)$" |> Js.Re.fromString,
    keys => {
      switch (keys->Belt_List.get(1)->Belt_Option.getWithDefault("")->int_of_string_opt) {
      | Some(height) => BlockIndexPage(height)
      | None => NotFound
      };
    },
  ),
  (
    "D([0-9]+)$" |> Js.Re.fromString,
    keys => {
      switch (
        keys->Belt_List.get(1)->Belt_Option.getWithDefault("")->int_of_string_opt,
        keys->Belt_List.get(2),
      ) {
      | (Some(id), Some("")) => DataSourceIndexPage(id, DataSourceExecute)
      | (Some(id), Some("code")) => DataSourceIndexPage(id, DataSourceCode)
      | (Some(id), Some("requests")) => DataSourceIndexPage(id, DataSourceRequests)
      | (Some(id), Some("revisions")) => DataSourceIndexPage(id, DataSourceRevisions)
      | _ => NotFound
      };
    },
  ),
  (
    "O([0-9]+)$" |> Js.Re.fromString,
    keys => {
      switch (
        keys->Belt_List.get(1)->Belt_Option.getWithDefault("")->int_of_string_opt,
        keys->Belt_List.get(2),
      ) {
      | (Some(id), Some("")) => OracleScriptIndexPage(id, OracleScriptExecute)
      | (Some(id), Some("code")) => OracleScriptIndexPage(id, OracleScriptCode)
      | (Some(id), Some("requests")) => OracleScriptIndexPage(id, OracleScriptRequests)
      | (Some(id), Some("revisions")) => OracleScriptIndexPage(id, OracleScriptRevisions)
      | _ => NotFound
      };
    },
  ),
  (
    "R([0-9]+)$" |> Js.Re.fromString,
    keys => {
      switch (
        keys->Belt_List.get(1)->Belt_Option.getWithDefault("")->int_of_string_opt,
        keys->Belt_List.get(2),
      ) {
      | (Some(id), Some("")) => RequestIndexPage(id, RequestReportStatus)
      | (Some(id), Some("proof")) => RequestIndexPage(id, RequestProof)
      | _ => NotFound
      };
    },
  ),
  (
    "^bandvaloper([0-9a-z]+)$" |> Js.Re.fromString,
    keys => {
      switch (
        keys->Belt_List.get(0)->Belt_Option.getWithDefault("")->Address.fromBech32Opt,
        keys->Belt_List.get(2),
      ) {
      | (Some(addr), Some("")) => ValidatorIndexPage(addr, ProposedBlocks)
      | (Some(addr), Some("delegators")) => ValidatorIndexPage(addr, Delegators)
      | (Some(addr), Some("reports")) => ValidatorIndexPage(addr, Reports)
      | _ => NotFound
      };
    },
  ),
  (
    "^band([0-9a-z]+)$" |> Js.Re.fromString,
    keys => {
      switch (
        keys->Belt_List.get(0)->Belt_Option.getWithDefault("")->Address.fromBech32Opt,
        keys->Belt_List.get(2),
      ) {
      | (Some(addr), Some("")) => AccountIndexPage(addr, AccountTransactions)
      | (Some(addr), Some("delegations")) => AccountIndexPage(addr, AccountDelegations)
      | _ => NotFound
      };
    },
  ),
];

let fromUrl = (url: ReasonReactRouter.url) => switch (
    regexes
    ->Belt_List.keepMap(((regex, keysToRoute)) =>
        url.path
        ->Belt_List.get(0)
        ->Belt_Option.getWithDefault("")
        ->Js.Re.exec_(regex, _)
        ->Belt_Option.map(result => (result, keysToRoute))
      )
    ->Belt.List.head
  ) {
  | Some((result, keysToRoute)) =>
    keysToRoute(
      result
      ->Js.Re.captures
      ->Belt_Array.keepMap(Js.toOption)
      ->Belt_List.fromArray
      ->Belt_List.concat(url.path->Belt_List.drop(1)->Belt_Option.getWithDefault([]))
      ->Belt_List.concat([url.hash]),
    )
  | None => NotFound
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
      {j|$validatorAddressBech32|j};
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
