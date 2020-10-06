type data_source_tab_t =
  | DataSourceExecute
  | DataSourceCode
  | DataSourceRequests
  | DataSourceRevisions;

type oracle_script_tab_t =
  | OracleScriptExecute
  | OracleScriptCode
  | OracleScriptBridgeCode
  | OracleScriptRequests
  | OracleScriptRevisions;

type account_tab_t =
  | AccountTransactions
  | AccountDelegations
  | AccountUnbonding
  | AccountRedelegate;

type validator_tab_t =
  | ProposedBlocks
  | Delegators
  | Reports
  | Reporters;

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
  | RequestHomePage
  | RequestIndexPage(int)
  | AccountIndexPage(Address.t, account_tab_t)
  | ValidatorHomePage
  | ValidatorIndexPage(Address.t, validator_tab_t)
  | ProposalHomePage
  | ProposalIndexPage(int)
  | IBCHomePage;

let fromUrl = (url: ReasonReactRouter.url) =>
  // TODO: We'll handle the NotFound case for Datasources and Oraclescript later
  switch (url.path, url.hash) {
  | (["data-sources"], _) => DataSourceHomePage
  | (["data-source", dataSourceID], hash) =>
    let urlHash = (
      fun
      | "code" => DataSourceCode
      | "execute" => DataSourceExecute
      | "revisions" => DataSourceRevisions
      | _ => DataSourceRequests
    );
    switch (dataSourceID |> int_of_string_opt) {
    | None => NotFound
    | Some(dataSourceIDInt) => DataSourceIndexPage(dataSourceIDInt, urlHash(hash))
    };

  | (["oracle-scripts"], _) => OracleScriptHomePage
  | (["oracle-script", oracleScriptID], hash) =>
    let urlHash = (
      fun
      | "code" => OracleScriptCode
      | "bridge" => OracleScriptBridgeCode
      | "execute" => OracleScriptExecute
      | "revisions" => OracleScriptRequests
      | _ => OracleScriptRequests
    );
    switch (oracleScriptID |> int_of_string_opt) {
    | None => NotFound
    | Some(oracleScriptIDInt) => OracleScriptIndexPage(oracleScriptIDInt, urlHash(hash))
    };

  | (["txs"], _) => TxHomePage
  | (["tx", txHash], _) => TxIndexPage(Hash.fromHex(txHash))
  | (["validators"], _) => ValidatorHomePage
  | (["blocks"], _) => BlockHomePage
  | (["block", blockHeight], _) =>
    let blockHeightIntOpt = blockHeight |> int_of_string_opt;
    switch (blockHeightIntOpt) {
    | Some(block) => BlockIndexPage(block)
    | None => NotFound
    };

  | (["requests"], _) => RequestHomePage
  | (["request", reqID], _) => RequestIndexPage(reqID |> int_of_string)
  | (["account", address], hash) =>
    let urlHash = (
      fun
      | "delegations" => AccountDelegations
      | "unbonding" => AccountUnbonding
      | "redelegate" => AccountRedelegate
      | _ => AccountTransactions
    );
    switch (address |> Address.fromBech32Opt) {
    | Some(address) => AccountIndexPage(address, urlHash(hash))
    | None => NotFound
    };
  | (["validator", address], hash) =>
    let urlHash = (
      fun
      | "delegators" => Delegators
      | "reporters" => Reporters
      | "proposed-blocks" => ProposedBlocks
      | _ => Reports
    );
    switch (address |> Address.fromBech32Opt) {
    | Some(address) => ValidatorIndexPage(address, urlHash(hash))
    | None => NotFound
    };
  | (["proposals"], _) => ProposalHomePage
  | (["proposal", proposalID], _) => ProposalIndexPage(proposalID |> int_of_string)
  | (["ibcs"], _) => IBCHomePage
  | ([], _) => HomePage
  | (_, _) => NotFound
  };

let toString =
  fun
  | DataSourceHomePage => "/data-sources"
  | DataSourceIndexPage(dataSourceID, DataSourceRequests) => {j|/data-source/$dataSourceID|j}
  | DataSourceIndexPage(dataSourceID, DataSourceCode) => {j|/data-source/$dataSourceID#code|j}
  | DataSourceIndexPage(dataSourceID, DataSourceExecute) => {j|/data-source/$dataSourceID#execute|j}
  | DataSourceIndexPage(dataSourceID, DataSourceRevisions) => {j|/data-source/$dataSourceID#revisions|j}
  | OracleScriptHomePage => "/oracle-scripts"
  | OracleScriptIndexPage(oracleScriptID, OracleScriptRequests) => {j|/oracle-script/$oracleScriptID|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptCode) => {j|/oracle-script/$oracleScriptID#code|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptBridgeCode) => {j|/oracle-script/$oracleScriptID#bridge|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptExecute) => {j|/oracle-script/$oracleScriptID#execute|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptRevisions) => {j|/oracle-script/$oracleScriptID#revisions|j}
  | TxHomePage => "/txs"
  | TxIndexPage(txHash) => {j|/tx/$txHash|j}
  | ValidatorHomePage => "/validators"
  | BlockHomePage => "/blocks"
  | BlockIndexPage(height) => {j|/block/$height|j}
  | RequestHomePage => "/requests"
  | RequestIndexPage(reqID) => {j|/request/$reqID|j}
  | AccountIndexPage(address, AccountTransactions) => {
      let addressBech32 = address |> Address.toBech32;
      {j|/account/$addressBech32|j};
    }
  | AccountIndexPage(address, AccountDelegations) => {
      let addressBech32 = address |> Address.toBech32;
      {j|/account/$addressBech32#delegations|j};
    }
  | AccountIndexPage(address, AccountUnbonding) => {
      let addressBech32 = address |> Address.toBech32;
      {j|/account/$addressBech32#unbonding|j};
    }
  | AccountIndexPage(address, AccountRedelegate) => {
      let addressBech32 = address |> Address.toBech32;
      {j|/account/$addressBech32#redelegate|j};
    }
  | ValidatorIndexPage(validatorAddress, Delegators) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|/validator/$validatorAddressBech32#delegators|j};
    }
  | ValidatorIndexPage(validatorAddress, Reports) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|/validator/$validatorAddressBech32#reports|j};
    }
  | ValidatorIndexPage(validatorAddress, Reporters) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|/validator/$validatorAddressBech32#reporters|j};
    }
  | ValidatorIndexPage(validatorAddress, ProposedBlocks) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|/validator/$validatorAddressBech32#proposed-blocks|j};
    }
  | ProposalHomePage => "/proposals"
  | ProposalIndexPage(proposalID) => {j|/proposal/$proposalID|j}
  | IBCHomePage => "/ibcs"
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
      } else if (len == 64 || str |> Js.String.startsWith("0x") && len == 66) {
        Some(TxIndexPage(str |> Hash.fromHex));
      } else if (capStr |> Js.String.startsWith("B")) {
        let%Opt blockID = str |> String.sub(_, 1, len - 1) |> int_of_string_opt;
        Some(BlockIndexPage(blockID));
      } else if (capStr |> Js.String.startsWith("D")) {
        let%Opt dataSourceID = str |> String.sub(_, 1, len - 1) |> int_of_string_opt;
        Some(DataSourceIndexPage(dataSourceID, DataSourceExecute));
      } else if (capStr |> Js.String.startsWith("R")) {
        let%Opt requestID = str |> String.sub(_, 1, len - 1) |> int_of_string_opt;
        Some(RequestIndexPage(requestID));
      } else if (capStr |> Js.String.startsWith("O")) {
        let%Opt oracleScriptID = str |> String.sub(_, 1, len - 1) |> int_of_string_opt;
        Some(OracleScriptIndexPage(oracleScriptID, OracleScriptExecute));
      } else {
        None;
      }
    }
  )
  |> Belt_Option.getWithDefault(_, NotFound);
};
