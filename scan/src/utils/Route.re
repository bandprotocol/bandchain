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

let fromUrl = (url: ReasonReactRouter.url) =>
  switch (url.path, url.hash) {
  | (["data-sources"], _) => DataSourceHomePage
  | (["data-source", dataSourceID], "code") =>
    DataSourceIndexPage(dataSourceID |> int_of_string, DataSourceCode)
  | (["data-source", dataSourceID], "requests") =>
    DataSourceIndexPage(dataSourceID |> int_of_string, DataSourceRequests)
  | (["data-source", dataSourceID], "revisions") =>
    DataSourceIndexPage(dataSourceID |> int_of_string, DataSourceRevisions)
  | (["data-source", dataSourceID], _) =>
    DataSourceIndexPage(dataSourceID |> int_of_string, DataSourceExecute)
  | (["oracle-scripts"], _) => OracleScriptHomePage
  | (["oracle-script", oracleScriptID], "code") =>
    OracleScriptIndexPage(oracleScriptID |> int_of_string, OracleScriptCode)
  | (["oracle-script", oracleScriptID], "requests") =>
    OracleScriptIndexPage(oracleScriptID |> int_of_string, OracleScriptRequests)
  | (["oracle-script", oracleScriptID], "revisions") =>
    OracleScriptIndexPage(oracleScriptID |> int_of_string, OracleScriptRevisions)
  | (["oracle-script", oracleScriptID], _) =>
    OracleScriptIndexPage(oracleScriptID |> int_of_string, OracleScriptExecute)
  | (["txs"], _) => TxHomePage
  | (["tx", txHash], _) => TxIndexPage(Hash.fromHex(txHash))
  | (["validators"], _) => ValidatorHomePage
  | (["blocks"], _) => BlockHomePage
  | (["block", blockHeight], _) =>
    let blockHeightIntOpt = blockHeight |> int_of_string_opt;
    BlockIndexPage(blockHeightIntOpt->Belt_Option.getWithDefault(0));
  | (["request", reqID], "proof") => RequestIndexPage(reqID |> int_of_string, RequestProof)
  | (["request", reqID], _) => RequestIndexPage(reqID |> int_of_string, RequestReportStatus)
  | (["account", address], "delegations") =>
    AccountIndexPage(address |> Address.fromBech32, AccountDelegations)
  | (["account", address], _) =>
    AccountIndexPage(address |> Address.fromBech32, AccountTransactions)
  | (["validator", address], "delegators") =>
    ValidatorIndexPage(address |> Address.fromBech32, Delegators)
  | (["validator", address], "reports") =>
    ValidatorIndexPage(address |> Address.fromBech32, Reports)
  | (["validator", address], _) =>
    ValidatorIndexPage(address |> Address.fromBech32, ProposedBlocks)
  | ([], "") => HomePage
  | (_, _) => NotFound
  };

let toString =
  fun
  | DataSourceHomePage => "/data-sources"
  | DataSourceIndexPage(dataSourceID, DataSourceExecute) => {j|/data-source/$dataSourceID|j}
  | DataSourceIndexPage(dataSourceID, DataSourceCode) => {j|/data-source/$dataSourceID#code|j}
  | DataSourceIndexPage(dataSourceID, DataSourceRequests) => {j|/data-source/$dataSourceID#requests|j}
  | DataSourceIndexPage(dataSourceID, DataSourceRevisions) => {j|/data-source/$dataSourceID#revisions|j}
  | OracleScriptHomePage => "/oracle-scripts"
  | OracleScriptIndexPage(oracleScriptID, OracleScriptExecute) => {j|/oracle-script/$oracleScriptID|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptCode) => {j|/oracle-script/$oracleScriptID#code|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptRequests) => {j|/oracle-script/$oracleScriptID#requests|j}
  | OracleScriptIndexPage(oracleScriptID, OracleScriptRevisions) => {j|/oracle-script/$oracleScriptID#revisions|j}
  | TxHomePage => "/txs"
  | TxIndexPage(txHash) => {j|/tx/$txHash|j}
  | ValidatorHomePage => "/validators"
  | BlockHomePage => "/blocks"
  | BlockIndexPage(height) => {j|/block/$height|j}
  | RequestIndexPage(reqID, RequestReportStatus) => {j|/request/$reqID|j}
  | RequestIndexPage(reqID, RequestProof) => {j|/request/$reqID#proof|j}
  | AccountIndexPage(address, AccountTransactions) => {
      let addressBech32 = address |> Address.toBech32;
      {j|/account/$addressBech32|j};
    }
  | AccountIndexPage(address, AccountDelegations) => {
      let addressBech32 = address |> Address.toBech32;
      {j|/account/$addressBech32#delegations|j};
    }
  | ValidatorIndexPage(validatorAddress, Delegators) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|/validator/$validatorAddressBech32#delegators|j};
    }
  | ValidatorIndexPage(validatorAddress, Reports) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|/validator/$validatorAddressBech32#reports|j};
    }
  | ValidatorIndexPage(validatorAddress, ProposedBlocks) => {
      let validatorAddressBech32 = validatorAddress |> Address.toOperatorBech32;
      {j|/validator/$validatorAddressBech32#proposed-blocks|j};
    }
  | HomePage
  | NotFound => "/";

let redirect = (route: t) => ReasonReactRouter.push(route |> toString);
