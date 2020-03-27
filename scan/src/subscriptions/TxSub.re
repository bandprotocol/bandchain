module RawDataReport = {
  type t = {
    externalDataID: int,
    data: JsBuffer.t,
  };

  let decode = json =>
    JsonUtils.Decode.{
      externalDataID: json |> field("externalDataID", int),
      data: json |> field("data", string) |> JsBuffer.fromBase64,
    };
};

module Coin = {
  type t = {
    denom: string,
    amount: float,
  };

  let decodeCoin = json =>
    JsonUtils.Decode.{
      denom: json |> field("denom", string),
      amount: json |> field("amount", uamount),
    };

  let newCoin = (denom, amount) => {denom, amount};

  let getBandAmountFromCoins = coins =>
    coins
    ->Belt_List.keep(coin => coin.denom == "uband")
    ->Belt_List.get(0)
    ->Belt_Option.mapWithDefault(0., coin => coin.amount /. 1e6);

  let getDescription = coin => {
    (coin.amount |> Format.fPretty)
    ++ " "
    ++ (
      switch (coin.denom.[0]) {
      | 'u' => coin.denom->String.sub(_, 1, (coin.denom |> String.length) - 1) |> String.uppercase
      | _ => coin.denom
      }
    );
  };

  let toCoinsString = coins => {
    coins
    ->Belt_List.map(coin => coin->getDescription)
    ->Belt_List.reduceWithIndex("", (des, acc, i) =>
        acc ++ des ++ (i + 1 < coins->Belt_List.size ? ", " : "")
      );
  };

  let getFeeAmount = coins => {
    let coinOpt = coins->Belt_List.get(0);
    switch (coinOpt) {
    | Some(coin) => coin.amount
    | None => 0.
    };
  };
};

module Msg = {
  module Send = {
    type t = {
      fromAddress: Address.t,
      toAddress: Address.t,
      amount: list(Coin.t),
    };

    let decode = json =>
      JsonUtils.Decode.{
        fromAddress: json |> field("from_address", string) |> Address.fromBech32,
        toAddress: json |> field("to_address", string) |> Address.fromBech32,
        amount: json |> field("amount", list(Coin.decodeCoin)),
      };
  };

  module CreateDataSource = {
    type t = {
      id: int,
      owner: Address.t,
      name: string,
      fee: list(Coin.t),
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: 0, // TODO , use id from events (not available right now)
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
        fee: json |> field("fee", list(Coin.decodeCoin)),
        executable: json |> field("executable", string) |> JsBuffer.fromBase64,
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module EditDataSource = {
    type t = {
      id: int,
      owner: Address.t,
      name: string,
      fee: list(Coin.t),
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> field("dataSourceID", intstr),
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
        fee: json |> field("fee", list(Coin.decodeCoin)),
        executable: json |> field("executable", string) |> JsBuffer.fromBase64,
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module CreateOracleScript = {
    type t = {
      id: ID.OracleScript.t,
      owner: Address.t,
      name: string,
      code: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> field("oracleScriptID", ID.OracleScript.fromJson),
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
        code: json |> field("code", string) |> JsBuffer.fromBase64,
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module EditOracleScript = {
    type t = {
      id: ID.OracleScript.t,
      owner: Address.t,
      name: string,
      code: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> field("oracleScriptID", ID.OracleScript.fromJson),
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
        code: json |> field("code", string) |> JsBuffer.fromBase64,
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module Request = {
    type t = {
      id: ID.Request.t,
      oracleScriptID: ID.OracleScript.t,
      calldata: JsBuffer.t,
      requestedValidatorCount: int,
      sufficientValidatorCount: int,
      expiration: int,
      prepareGas: int,
      executeGas: int,
      sender: Address.t,
    };

    let decode = json => {
      Js.Console.log(json);
      JsonUtils.Decode.{
        id: json |> field("requestID", ID.Request.fromJson),
        oracleScriptID: json |> field("oracleScriptID", ID.OracleScript.fromJson),
        calldata: json |> field("calldata", string) |> JsBuffer.fromBase64,
        requestedValidatorCount: json |> field("requestedValidatorCount", int),
        sufficientValidatorCount: json |> field("sufficientValidatorCount", int),
        expiration: json |> field("expiration", int),
        prepareGas: json |> field("prepareGas", int),
        executeGas: json |> field("executeGas", int),
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
    };
  };

  module Report = {
    type t = {
      requestID: ID.Request.t,
      dataSet: list(RawDataReport.t),
      validator: Address.t,
      reporter: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        requestID: json |> field("requestID", ID.Request.fromJson),
        dataSet: json |> field("dataSet", list(RawDataReport.decode)),
        validator: json |> field("validator", string) |> Address.fromBech32,
        reporter: json |> field("reporter", string) |> Address.fromBech32,
      };
  };
  module AddOracleAddress = {
    type t = {
      validator: string,
      reporterAddress: Address.t,
      sender: Address.t,
    };
  };

  module RemoveOracleAddress = {
    type t = {
      validator: string,
      reporterAddress: Address.t,
      sender: Address.t,
    };
  };

  module CreateValidator = {
    type t = {
      moniker: string,
      identity: string,
      website: string,
      details: string,
      commissionRate: float,
      commissionMaxRate: float,
      commissionMaxChange: float,
      delegatorAddress: Address.t,
      validatorAddress: Address.t,
      publicKey: PubKey.t,
      minSelfDelegation: list(Coin.t),
      selfDelegation: list(Coin.t),
      sender: Address.t,
    };
  };

  module EditValidator = {
    type t = {
      moniker: string,
      identity: string,
      website: string,
      details: string,
      commissionRate: float,
      validatorAddress: Address.t,
      minSelfDelegation: list(Coin.t),
      sender: Address.t,
    };
  };

  type t =
    | Unknown
    | Send(Send.t)
    | CreateDataSource(CreateDataSource.t)
    | EditDataSource(EditDataSource.t)
    | CreateOracleScript(CreateOracleScript.t)
    | EditOracleScript(EditOracleScript.t)
    | Request(Request.t)
    | Report(Report.t)
    | AddOracleAddress(AddOracleAddress.t)
    | RemoveOracleAddress(RemoveOracleAddress.t)
    | CreateValidator(CreateValidator.t)
    | EditValidator(EditValidator.t);

  let getCreator = msg => {
    switch (msg) {
    | Send(send) => send.fromAddress
    | CreateDataSource(dataSource) => dataSource.sender
    | EditDataSource(dataSource) => dataSource.sender
    | CreateOracleScript(oracleScript) => oracleScript.sender
    | EditOracleScript(oracleScript) => oracleScript.sender
    | Request(request) => request.sender
    | Report(report) => report.reporter
    | AddOracleAddress(address) => address.sender
    | RemoveOracleAddress(address) => address.sender
    | CreateValidator(validator) => validator.sender
    | EditValidator(validator) => validator.sender
    | Unknown => "" |> Address.fromHex
    };
  };

  let getDescription = msg => {
    switch (msg) {
    | Send(send) =>
      (send.amount |> Coin.toCoinsString) ++ "->" ++ (send.toAddress |> Address.toBech32)
    | CreateDataSource(dataSource) => dataSource.name
    | EditDataSource(dataSource) => dataSource.name
    | CreateOracleScript(oracleScript) => oracleScript.name
    | EditOracleScript(oracleScript) => oracleScript.name
    | Request(_) => "yo"
    // switch (msg.events->Event.getValueOfKey("request.code_name")) {
    // | Some(value) =>
    //   switch (msg.events->Event.getValueOfKey("request.id")) {
    //   | Some(id) => "#" ++ id ++ " " ++ value
    //   | None => ""
    //   }
    // | None => "?"
    // }
    | Report(report) => "yo"
    //   switch (msg.events->Event.getValueOfKey("report.code_name")) {
    //   | Some(value) => "#" ++ (report.requestID |> string_of_int) ++ " " ++ value
    //   | None => "?"
    //   }
    | AddOracleAddress(_) => "ADDORACLEADDRESS DESCRIPTION"
    | RemoveOracleAddress(_) => "REMOVEORACLEADDRESS DESCRIPTION"
    | CreateValidator(_) => "CREATEVALIDATOR DESCRIPTION"
    | EditValidator(_) => "EDITVALIDATOR DESCRIPTION"
    | Unknown => "Unknown"
    };
  };

  let decodeAction = json => {
    Js.Console.log(json);
    JsonUtils.Decode.(
      switch (json |> field("type", string)) {
      | "send" => Send(json |> Send.decode)
      | "createDataSource" => CreateDataSource(json |> CreateDataSource.decode)
      | "editDataSource" => EditDataSource(json |> EditDataSource.decode)
      | "createOracleScript" => CreateOracleScript(json |> CreateOracleScript.decode)
      | "editOracleScript" => EditOracleScript(json |> EditOracleScript.decode)
      | "request" => Request(json |> Request.decode)
      | "report" => Report(json |> Report.decode)
      | _ => Unknown
      }
    );
  };
  let decodeActions = json => json |> JsonUtils.Decode.list(decodeAction);

  let getRoute = msg =>
    switch (msg) {
    | Send(_) => None
    // TODO: Route to each data source and oracle script page
    | CreateDataSource(_) => None
    // switch (msg.events->Event.getValueOfKey("store_code.codehash")) {
    // | Some(value) => Some(Route.ScriptIndexPage(value |> Hash.fromHex, ScriptTransactions))
    // | None => None
    // }
    | EditDataSource(_) => None
    | CreateOracleScript(_) => None
    | EditOracleScript(_) => None
    | Request(_) => None
    // switch (msg.events->Event.getValueOfKey("request.id")) {
    // | Some(value) => Some(Route.RequestIndexPage(value->int_of_string, RequestReportStatus))
    // | None => None
    // }
    | Report(_) => None
    // switch (msg.events->Event.getValueOfKey("report.id")) {
    // | Some(value) => Some(Route.RequestIndexPage(value->int_of_string, RequestReportStatus))
    // | None => None
    // }
    | AddOracleAddress(_) => None
    | RemoveOracleAddress(_) => None
    | CreateValidator(_) => None
    | EditValidator(_) => None
    | Unknown => None
    };
};

type t = {
  txHash: Hash.t,
  blockHeight: ID.Block.t,
  success: bool,
  gasFee: list(TxHook.Coin.t),
  gasLimit: int,
  gasUsed: int,
  sender: Address.t,
  timestamp: MomentRe.Moment.t,
  messages: list(Msg.t),
};

module Mini = {
  type t = {
    txHash: Hash.t,
    blockHeight: ID.Block.t,
    timestamp: MomentRe.Moment.t,
  };
};

module SingleConfig = [%graphql
  {|
  subscription Transaction($tx_hash:bytea!) {
    transactions_by_pk(tx_hash: $tx_hash) @bsRecord {
      txHash : tx_hash @bsDecoder(fn: "GraphQLParser.hash")
      blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
      success
      gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
      gasLimit: gas_limit @bsDecoder(fn: "GraphQLParser.int64")
      gasUsed : gas_used @bsDecoder(fn: "GraphQLParser.int64")
      sender  @bsDecoder(fn: "Address.fromBech32")
      timestamp  @bsDecoder(fn: "GraphQLParser.time")
      messages @bsDecoder(fn: "Msg.decodeActions")
    }
  },
|}
];

module MultiConfig = [%graphql
  {|
  subscription Transaction($limit: Int!, $offset: Int!) {
    transactions(offset: $offset, limit: $limit, order_by: {block_height: desc}) @bsRecord {
      txHash : tx_hash @bsDecoder(fn: "GraphQLParser.hash")
      blockHeight: block_height @bsDecoder(fn: "ID.Block.fromJson")
      success
      gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
      gasLimit: gas_limit @bsDecoder(fn: "GraphQLParser.int64")
      gasUsed : gas_used @bsDecoder(fn: "GraphQLParser.int64")
      sender  @bsDecoder(fn: "Address.fromBech32")
      timestamp  @bsDecoder(fn: "GraphQLParser.time")
      messages @bsDecoder(fn: "Msg.decodeActions")
    }
  }
|}
];

module TxCountConfig = [%graphql
  {|
  subscription Transaction {
    transactions_aggregate{
      aggregate{
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

module Event = {
  type t = {
    key: string,
    value: string,
  };

  let decode = (prefix, json) =>
    JsonUtils.Decode.{
      key: prefix ++ "." ++ (json |> field("key", string)),
      value: json |> field("value", string),
    };

  let decodeEvent = json =>
    JsonUtils.Decode.(
      {
        let prefix = json |> field("type", string);
        json |> field("attributes", list(decode(prefix)));
      }
    );

  let decodeEvents = json =>
    List.flatten(JsonUtils.Decode.(json |> field("events", list(decodeEvent))));

  let getValueOfKey = (events: list(t), key) =>
    events
    ->Belt.List.keepMap(event => event.key == key ? Some(event.value) : None)
    ->Belt.List.get(0);
};

let get = txHash => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=
        SingleConfig.makeVariables(
          ~tx_hash=txHash |> Hash.toHex |> (x => "\x" ++ x) |> Js.Json.string,
          (),
        ),
    );
  let%Sub x = result;
  switch (x##transactions_by_pk) {
  | Some(data) => Sub.resolve(data)
  | None => NoData
  };
};

let getList = (~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiConfig.definition,
      ~variables=MultiConfig.makeVariables(~limit=pageSize, ~offset, ()),
    );
  Js.Console.log(result);

  result |> Sub.map(_, x => x##transactions);
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(TxCountConfig.definition);
  result
  |> Sub.map(_, x => x##transactions_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};
