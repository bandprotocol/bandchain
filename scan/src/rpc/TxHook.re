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

  let getDescription = coin => (coin.amount |> Format.fPretty) ++ " " ++ coin.denom;
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
      owner: Address.t,
      name: string,
      fee: list(Coin.t),
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
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
      owner: Address.t,
      name: string,
      code: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
        code: json |> field("code", string) |> JsBuffer.fromBase64,
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module EditOracleScript = {
    type t = {
      id: int,
      owner: Address.t,
      name: string,
      code: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> field("oracleScriptID", intstr),
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
        code: json |> field("code", string) |> JsBuffer.fromBase64,
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module Request = {
    type t = {
      oracleScriptID: int,
      calldata: JsBuffer.t,
      requestedValidatorCount: int,
      sufficientValidatorCount: int,
      expiration: int,
      prepareGas: int,
      executeGas: int,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        oracleScriptID: json |> field("oracleScriptID", intstr),
        calldata: json |> field("calldata", string) |> JsBuffer.fromBase64,
        requestedValidatorCount: json |> field("requestedValidatorCount", intstr),
        sufficientValidatorCount: json |> field("sufficientValidatorCount", intstr),
        expiration: json |> field("expiration", intstr),
        prepareGas: json |> field("prepareGas", intstr),
        executeGas: json |> field("executeGas", intstr),
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module Report = {
    type t = {
      requestID: int,
      dataSet: list(RequestHook.RawDataReport.t),
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        requestID: json |> field("requestID", intstr),
        dataSet: json |> field("dataSet", list(RequestHook.RawDataReport.decode)),
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  type action_t =
    | Unknown
    | Send(Send.t)
    | CreateDataSource(CreateDataSource.t)
    | EditDataSource(EditDataSource.t)
    | CreateOracleScript(CreateOracleScript.t)
    | EditOracleScript(EditOracleScript.t)
    | Request(Request.t)
    | Report(Report.t);

  type t = {
    action: action_t,
    events: list(Event.t),
  };

  let getCreator = msg => {
    switch (msg.action) {
    | Send(send) => send.fromAddress
    | CreateDataSource(dataSource) => dataSource.owner
    | EditDataSource(dataSource) => dataSource.owner
    | CreateOracleScript(oracleScript) => oracleScript.owner
    | EditOracleScript(oracleScript) => oracleScript.owner
    | Request(request) => request.sender
    | Report(report) => report.sender
    | Unknown => "" |> Address.fromHex
    };
  };

  let getDescription = msg => {
    switch (msg.action) {
    | Send(send) =>
      send.amount
      ->Belt_List.map(coin => coin->Coin.getDescription)
      ->Belt_List.reduceWithIndex("", (des, acc, i) =>
          acc ++ des ++ (i + 1 < send.amount->Belt_List.size ? ", " : "")
        )
      ++ "->"
      ++ (send.toAddress |> Address.toBech32)
    | CreateDataSource(dataSource) => dataSource.name
    | EditDataSource(dataSource) => dataSource.name
    | CreateOracleScript(oracleScript) => oracleScript.name
    | EditOracleScript(oracleScript) => oracleScript.name
    | Request(_) =>
      switch (msg.events->Event.getValueOfKey("request.code_name")) {
      | Some(value) =>
        switch (msg.events->Event.getValueOfKey("request.id")) {
        | Some(id) => "#" ++ id ++ " " ++ value
        | None => ""
        }
      | None => "?"
      }
    | Report(report) =>
      switch (msg.events->Event.getValueOfKey("report.code_name")) {
      | Some(value) => "#" ++ (report.requestID |> string_of_int) ++ " " ++ value
      | None => "?"
      }
    | Unknown => "Unknown"
    };
  };

  let decodeAction = json =>
    JsonUtils.Decode.(
      switch (json |> field("type", string)) {
      | "cosmos-sdk/MsgSend" => Send(json |> field("value", Send.decode))
      | "zoracle/CreateDataSource" =>
        CreateDataSource(json |> field("value", CreateDataSource.decode))
      | "zoracle/EditDataSource" => EditDataSource(json |> field("value", EditDataSource.decode))
      | "zoracle/CreateOracleScript" =>
        CreateOracleScript(json |> field("value", CreateOracleScript.decode))
      | "zoracle/EditOracleScript" =>
        EditOracleScript(json |> field("value", EditOracleScript.decode))
      | "zoracle/Request" => Request(json |> field("value", Request.decode))
      | "zoracle/Report" => Report(json |> field("value", Report.decode))
      | _ => Unknown
      }
    );

  let getRoute = msg =>
    switch (msg.action) {
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
    | Request(_) =>
      switch (msg.events->Event.getValueOfKey("request.id")) {
      | Some(value) => Some(Route.RequestIndexPage(value->int_of_string, RequestReportStatus))
      | None => None
      }
    | Report(_) =>
      switch (msg.events->Event.getValueOfKey("report.id")) {
      | Some(value) => Some(Route.RequestIndexPage(value->int_of_string, RequestReportStatus))
      | None => None
      }
    | Unknown => None
    };
};

module Signature = {
  type t = {
    pubKey: PubKey.t,
    pubKeyType: string,
    signature: JsBuffer.t,
  };

  let decode = json =>
    JsonUtils.Decode.{
      pubKey: json |> at(["pub_key", "value"], string) |> PubKey.fromBase64,
      pubKeyType: json |> at(["pub_key", "type"], string),
      signature: json |> field("signature", string) |> JsBuffer.fromBase64,
    };
};

module Tx = {
  type t = {
    sender: Address.t,
    blockHeight: int,
    hash: Hash.t,
    timestamp: MomentRe.Moment.t,
    gasWanted: int,
    gasUsed: int,
    messages: list(Msg.t),
  };

  let decodeTx = json =>
    JsonUtils.Decode.{
      sender:
        json
        |> at(["tx", "value", "signatures"], list(Signature.decode))
        |> Belt_List.getExn(_, 0)
        |> ((firstSignature: Signature.t) => firstSignature.pubKey)
        |> PubKey.toAddress,
      blockHeight: json |> field("height", intstr),
      hash: json |> field("txhash", string) |> Hash.fromHex,
      timestamp: json |> field("timestamp", moment),
      gasWanted: json |> field("gas_wanted", intstr),
      gasUsed: json |> field("gas_used", intstr),
      messages: {
        let actions = json |> at(["tx", "value", "msg"], list(Msg.decodeAction));
        let eventDoubleLists =
          json
          |> optional(field("logs", list(Event.decodeEvents)))
          |> Belt.Option.getWithDefault(_, actions->Belt_List.map(_ => []));
        Belt.List.zip(actions, eventDoubleLists)
        ->Belt.List.map(((action, events)) => Msg.{action, events});
      },
    };

  let getDescription = tx => tx.messages->Belt_List.getExn(0)->Msg.getDescription;
};

module Txs = {
  type t = {
    totalCount: int,
    txs: list(Tx.t),
  };

  let decodeTxs = json =>
    JsonUtils.Decode.{
      totalCount: json |> field("total_count", intstr),
      txs: json |> field("txs", list(Tx.decodeTx)),
    };
};

let atHash = txHash => {
  let txHashHex = txHash->Hash.toHex;
  let json = AxiosHooks.use({j|txs/$txHashHex|j});
  json |> Belt.Option.map(_, Tx.decodeTx);
};

let atHeight = (height, ~page=1, ~limit=25, ()) => {
  let json = AxiosHooks.use({j|txs?tx.height=$height&page=$page&limit=$limit|j});
  json |> Belt.Option.map(_, Txs.decodeTxs);
};

let latest = (~page=1, ~limit=10, ()) => {
  let json = AxiosHooks.use({j|d3n/txs/latest?page=$page&limit=$limit|j});
  json |> Belt.Option.map(_, Txs.decodeTxs);
};

let withCodehash = (~codeHash, ~page=1, ~limit=10, ()) => {
  let codeHashHex = codeHash->Hash.toHex;
  let json = AxiosHooks.use({j|txs?request.codehash=$codeHashHex&page=$page&limit=$limit|j});
  json |> Belt.Option.map(_, Txs.decodeTxs);
};
