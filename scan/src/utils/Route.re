type data_source_tab_t =
  | DataSourceExecute
  | DataSourceCode
  | DataSourceRequests
  | DataSourceRevisions;

let dataSourceTab = key =>
  switch (key) {
  | "" => Some(DataSourceExecute)
  | "code" => Some(DataSourceCode)
  | "requests" => Some(DataSourceRequests)
  | "revisions" => Some(DataSourceRevisions)
  | _ => None
  };

type oracle_script_tab_t =
  | OracleScriptExecute
  | OracleScriptCode
  | OracleScriptRequests
  | OracleScriptRevisions;

let oracleScriptTab = key =>
  switch (key) {
  | "" => Some(OracleScriptExecute)
  | "code" => Some(OracleScriptCode)
  | "requests" => Some(OracleScriptRequests)
  | "revisions" => Some(OracleScriptRevisions)
  | _ => None
  };

type request_tab_t =
  | RequestReportStatus
  | RequestProof;

let requestTab = key =>
  switch (key) {
  | "" => Some(RequestReportStatus)
  | "proof" => Some(RequestProof)
  | _ => None
  };

type account_tab_t =
  | AccountTransactions
  | AccountDelegations;

let accountTab = key =>
  switch (key) {
  | "" => Some(AccountTransactions)
  | "delegations" => Some(AccountDelegations)
  | _ => None
  };

type validator_tab_t =
  | ProposedBlocks
  | Delegators
  | Reports;

let validatorTab = key =>
  switch (key) {
  | "" => Some(ProposedBlocks)
  | "delegators" => Some(Delegators)
  | "reports" => Some(Reports)
  | _ => None
  };

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

type url_parts_t = {
  first: list(string),
  remaining: list(string),
  urlHash: string,
};

let regexes = [
  ("^$", _ => HomePage),
  ("^blocks$", _ => BlockHomePage),
  ("^scripts$", _ => OracleScriptHomePage),
  ("^sources$", _ => DataSourceHomePage),
  ("^txs$", _ => TxHomePage),
  (
    "^tx$",
    url_parts => {
      switch (url_parts.remaining->Belt_List.get(0)) {
      | Some(txHash) => TxIndexPage(Hash.fromHex(txHash))
      | None => NotFound
      };
    },
  ),
  ("^validators$", _ => ValidatorHomePage),
  (
    "B([0-9]+)$",
    url_parts => {
      switch (
        url_parts.first->Belt_List.get(1)->Belt_Option.getWithDefault("")->int_of_string_opt
      ) {
      | Some(height) => BlockIndexPage(height)
      | None => NotFound
      };
    },
  ),
  (
    "D([0-9]+)$",
    url_parts =>
      dataSourceTab(url_parts.urlHash)
      ->Belt_Option.mapWithDefault(NotFound, tab =>
          DataSourceIndexPage(url_parts.first->Belt_List.getExn(1)->int_of_string, tab)
        ),
  ),
  (
    "O([0-9]+)$",
    url_parts =>
      oracleScriptTab(url_parts.urlHash)
      ->Belt_Option.mapWithDefault(NotFound, tab =>
          OracleScriptIndexPage(url_parts.first->Belt_List.getExn(1)->int_of_string, tab)
        ),
  ),
  (
    "R([0-9]+)$",
    url_parts =>
      requestTab(url_parts.urlHash)
      ->Belt_Option.mapWithDefault(NotFound, tab =>
          RequestIndexPage(url_parts.first->Belt_List.getExn(1)->int_of_string, tab)
        ),
  ),
  (
    "^bandvaloper([0-9a-z]+)$",
    url_parts => {
      switch (
        url_parts.first->Belt_List.head->Belt_Option.getWithDefault("")->Address.fromBech32Opt,
        url_parts.urlHash->validatorTab,
      ) {
      | (Some(addr), Some(tab)) => ValidatorIndexPage(addr, tab)
      | _ => NotFound
      };
    },
  ),
  (
    "^band([0-9a-z]+)$",
    url_parts => {
      switch (
        url_parts.first->Belt_List.head->Belt_Option.getWithDefault("")->Address.fromBech32Opt,
        url_parts.urlHash->accountTab,
      ) {
      | (Some(addr), Some(tab)) => AccountIndexPage(addr, tab)
      | _ => NotFound
      };
    },
  ),
];

let fromUrl = (url: ReasonReactRouter.url) =>
  switch (
    regexes
    ->Belt_List.keepMap(((regex, keysToRoute)) =>
        url.path
        ->Belt_List.head
        ->Belt_Option.getWithDefault("")
        ->Js.Re.exec_(regex->Js.Re.fromString, _)
        ->Belt_Option.map(result => (result, keysToRoute))
      )
    ->Belt.List.head
  ) {
  | Some((result, keysToRoute)) =>
    keysToRoute({
      first: result->Js.Re.captures->Belt_Array.keepMap(Js.toOption)->Belt_List.fromArray,
      remaining: url.path->Belt_List.drop(1)->Belt_Option.getWithDefault([]),
      urlHash: url.hash,
    })
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
