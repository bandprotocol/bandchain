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
      id: ID.DataSource.t,
      owner: Address.t,
      name: string,
      fee: list(Coin.t),
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> field("data_source_id", ID.DataSource.fromJson),
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
        fee: json |> field("fee", list(Coin.decodeCoin)),
        executable: json |> field("executable", string) |> JsBuffer.fromBase64,
        sender: json |> field("sender", string) |> Address.fromBech32,
      };
  };

  module EditDataSource = {
    type t = {
      id: ID.DataSource.t,
      owner: Address.t,
      name: string,
      fee: list(Coin.t),
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> field("data_source_id", ID.DataSource.fromJson),
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
        id: json |> field("oracle_script_id", ID.OracleScript.fromJson),
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
        id: json |> field("oracle_script_id", ID.OracleScript.fromJson),
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
      oracleScriptName: string,
      calldata: JsBuffer.t,
      requestedValidatorCount: int,
      sufficientValidatorCount: int,
      sender: Address.t,
    };

    let decode = json => {
      JsonUtils.Decode.{
        id: json |> field("request_id", ID.Request.fromJson),
        oracleScriptID: json |> field("oracleScriptID", ID.OracleScript.fromJson),
        oracleScriptName: json |> field("oracle_script_name", string),
        calldata: json |> field("calldata", string) |> JsBuffer.fromBase64,
        requestedValidatorCount: json |> field("requestedValidatorCount", int),
        sufficientValidatorCount: json |> field("sufficientValidatorCount", int),
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
      validator: Address.t,
      reporter: Address.t,
      validatorMoniker: string,
    };
    let decode = json =>
      JsonUtils.Decode.{
        validator: json |> field("validator", string) |> Address.fromBech32,
        reporter: json |> field("reporter", string) |> Address.fromBech32,
        validatorMoniker: json |> field("validator_moniker", string),
      };
  };

  module RemoveOracleAddress = {
    type t = {
      validator: Address.t,
      reporter: Address.t,
      validatorMoniker: string,
    };
    let decode = json =>
      JsonUtils.Decode.{
        validator: json |> field("validator", string) |> Address.fromBech32,
        reporter: json |> field("reporter", string) |> Address.fromBech32,
        validatorMoniker: json |> field("validator_moniker", string),
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
      minSelfDelegation: float,
    };
    let decode = json =>
      JsonUtils.Decode.{
        moniker: json |> at(["description", "moniker"], string),
        identity: json |> at(["description", "identity"], string),
        website: json |> at(["description", "website"], string),
        details: json |> at(["description", "details"], string),
        commissionRate: json |> at(["commission", "rate"], floatstr),
        commissionMaxRate: json |> at(["commission", "max_rate"], floatstr),
        commissionMaxChange: json |> at(["commission", "max_change_rate"], floatstr),
        delegatorAddress: json |> field("delegator_address", string) |> Address.fromBech32,
        validatorAddress: json |> field("validator_address", string) |> Address.fromBech32,
        publicKey: json |> field("pubkey", string) |> PubKey.fromBech32,
        minSelfDelegation: json |> field("min_self_delegation", floatstr),
      };
  };

  module EditValidator = {
    type t = {
      moniker: string,
      identity: string,
      website: string,
      details: string,
      commissionRate: float,
      sender: Address.t,
      minSelfDelegation: float,
    };
    let decode = json =>
      JsonUtils.Decode.{
        moniker: json |> field("moniker", string),
        identity: json |> field("identity", string),
        website: json |> field("website", string),
        details: json |> field("details", string),
        commissionRate: json |> field("commission_rate", floatstr),
        sender: json |> field("address", string) |> Address.fromBech32,
        minSelfDelegation: json |> field("min_self_delegation", floatstr),
      };
  };

  module CreateClient = {
    type t = {
      address: Address.t,
      clientID: string,
      chainID: string,
      trustingPeriod: MomentRe.Duration.t,
      unbondingPeriod: MomentRe.Duration.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        address: json |> field("address", string) |> Address.fromBech32,
        clientID: json |> field("client_id", string),
        chainID: "band-consumer",
        trustingPeriod:
          (json |> field("trusting_period", JsonUtils.Decode.float))
          /. 1_000_000.
          |> MomentRe.durationMillis,
        unbondingPeriod:
          (json |> field("unbonding_period", JsonUtils.Decode.float))
          /. 1_000_000.
          |> MomentRe.durationMillis,
      };
    };
  };

  module UpdateClient = {
    type t = {
      address: Address.t,
      clientID: string,
      chainID: string,
      validatorHash: Hash.t,
      prevValidatorHash: Hash.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        address: json |> field("address", string) |> Address.fromBech32,
        clientID: json |> field("client_id", string),
        chainID: "band-consumer",
        validatorHash:
          "88A40098ACD2B9DDBD81E1397217BD0BC6D5B90C86D3614927B2CCDF4A3023BD" |> Hash.fromBase64,
        prevValidatorHash:
          "88A40098ACD2B9DDBD81E1397217BD0BC6D5B90C86D3614927B2CCDF4A3023BD" |> Hash.fromBase64,
      };
    };
  };

  module SubmitClientMisbehaviour = {
    type t = {
      address: Address.t,
      clientID: string,
      chainID: string,
      validatorHash: Hash.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        address: json |> field("address", string) |> Address.fromBech32,
        clientID: json |> field("client_id", string),
        chainID: "band-consumer",
        validatorHash:
          "88A40098ACD2B9DDBD81E1397217BD0BC6D5B90C86D3614927B2CCDF4A3023BD" |> Hash.fromBase64,
      };
    };
  };

  module Packet = {
    type common_t = {
      sequence: int,
      sourcePort: string,
      sourceChannel: string,
      destinationPort: string,
      destinationChannel: string,
      timeoutHeight: ID.Block.t,
      chainID: string,
    };

    type t = {
      sender: Address.t,
      data: string,
      common: common_t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        sender: json |> field("signer", string) |> Address.fromBech32,
        data: json |> at(["packet", "data"], string) |> Js.String.toUpperCase,
        common: {
          sequence: json |> at(["packet", "sequence"], int),
          sourcePort: json |> at(["packet", "source_port"], string),
          sourceChannel: json |> at(["packet", "source_channel"], string),
          destinationPort: json |> at(["packet", "destination_port"], string),
          destinationChannel: json |> at(["packet", "destination_channel"], string),
          timeoutHeight: json |> at(["packet", "timeout_height"], ID.Block.fromJson),
          chainID: "band-consumer",
        },
      };
    };
  };

  module Acknowledgement = {
    type t = {
      common: Packet.common_t,
      sender: Address.t,
      acknowledgement: string,
    };
    let decode = json =>
      JsonUtils.Decode.{
        common: {
          sequence: 999,
          sourcePort: "gjdojfpjfp",
          sourceChannel: "gjdojfpjfp",
          destinationPort: "gjdojfpjfp",
          destinationChannel: "gjdojfpjfp",
          timeoutHeight: ID.Block.ID(999999),
          chainID: "band-consumer",
        },
        sender: json |> field("address", string) |> Address.fromBech32,
        acknowledgement: "iKQAmKzSud29geE5che9C8bVuQyG02FJJ7LM...",
      };
  };

  module Timeout = {
    type t = {
      sender: Address.t,
      common: Packet.common_t,
      nextSequenceReceive: int,
    };
    let decode = json =>
      JsonUtils.Decode.{
        sender: json |> field("signer", string) |> Address.fromBech32,
        common: {
          sequence: json |> at(["packet", "sequence"], int),
          sourcePort: json |> at(["packet", "source_port"], string),
          sourceChannel: json |> at(["packet", "source_channel"], string),
          destinationPort: json |> at(["packet", "destination_port"], string),
          destinationChannel: json |> at(["packet", "destination_channel"], string),
          timeoutHeight: json |> at(["packet", "timeout_height"], ID.Block.fromJson),
          chainID: "band-consumer",
        },
        nextSequenceReceive: json |> at(["packet", "next_sequence_receive"], int),
      };
  };

  module ConnectionCommon = {
    type t = {
      chainID: string,
      connectionID: string,
    };
    let decode = json =>
      JsonUtils.Decode.{
        chainID: "band-consumer",
        connectionID: json |> field("connection_id", string),
      };
  };

  module ConnectionOpenInit = {
    type t = {
      signer: Address.t,
      common: ConnectionCommon.t,
      clientID: string,
      consensusHeight: ID.Block.t,
    };
    let decode = json =>
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ConnectionCommon.decode,
        clientID: json |> field("client_id", string),
        consensusHeight: json |> field("consensus_height", ID.Block.fromJson),
      };
  };

  module ConnectionOpenTry = {
    type t = {
      signer: Address.t,
      common: ConnectionCommon.t,
      clientID: string,
      consensusHeight: ID.Block.t,
    };
    let decode = json =>
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        clientID: json |> field("client_id", string),
        common: json |> ConnectionCommon.decode,
        consensusHeight: json |> field("consensus_height", ID.Block.fromJson),
      };
  };

  module ConnectionOpenAck = {
    type t = {
      signer: Address.t,
      common: ConnectionCommon.t,
      consensusHeight: ID.Block.t,
    };
    let decode = json =>
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ConnectionCommon.decode,
        consensusHeight: json |> field("consensus_height", ID.Block.fromJson),
      };
  };

  module ConnectionOpenConfirm = {
    type t = {
      signer: Address.t,
      common: ConnectionCommon.t,
    };
    let decode = json =>
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ConnectionCommon.decode,
      };
  };

  module ChannelCommon = {
    type t = {
      chainID: string,
      portID: string,
      channelID: string,
    };

    let decode = json => {
      JsonUtils.Decode.{
        chainID: "band-consumer",
        portID: json |> field("port_id", string),
        channelID: json |> field("channel_id", string),
      };
    };
  };

  module ChannelOpenInit = {
    type t = {
      signer: Address.t,
      common: ChannelCommon.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ChannelCommon.decode,
      };
    };
  };

  module ChannelOpenTry = {
    type t = {
      signer: Address.t,
      common: ChannelCommon.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ChannelCommon.decode,
      };
    };
  };

  module ChannelOpenAck = {
    type t = {
      signer: Address.t,
      common: ChannelCommon.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ChannelCommon.decode,
      };
    };
  };

  module ChannelOpenConfirm = {
    type t = {
      signer: Address.t,
      common: ChannelCommon.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ChannelCommon.decode,
      };
    };
  };

  module ChannelCloseInit = {
    type t = {
      signer: Address.t,
      common: ChannelCommon.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ChannelCommon.decode,
      };
    };
  };

  module ChannelCloseConfirm = {
    type t = {
      signer: Address.t,
      common: ChannelCommon.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ChannelCommon.decode,
      };
    };
  };

  module FailMessage = {
    type t = {
      sender: Address.t,
      message: string,
    };
    let decode = json =>
      JsonUtils.Decode.{
        sender: json |> field("sender", string) |> Address.fromBech32,
        message: json |> field("type", string),
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
    | EditValidator(EditValidator.t)
    | CreateClient(CreateClient.t)
    | UpdateClient(UpdateClient.t)
    | SubmitClientMisbehaviour(SubmitClientMisbehaviour.t)
    | ConnectionOpenInit(ConnectionOpenInit.t)
    | ConnectionOpenTry(ConnectionOpenTry.t)
    | ConnectionOpenAck(ConnectionOpenAck.t)
    | ConnectionOpenConfirm(ConnectionOpenConfirm.t)
    | ChannelOpenInit(ChannelOpenInit.t)
    | ChannelOpenTry(ChannelOpenTry.t)
    | ChannelOpenAck(ChannelOpenAck.t)
    | ChannelOpenConfirm(ChannelOpenConfirm.t)
    | ChannelCloseInit(ChannelCloseInit.t)
    | ChannelCloseConfirm(ChannelCloseConfirm.t)
    | Packet(Packet.t)
    | Acknowledgement(Acknowledgement.t)
    | Timeout(Timeout.t)
    | FailMessage(FailMessage.t);

  let getCreator = msg => {
    switch (msg) {
    | Send(send) => send.fromAddress
    | CreateDataSource(dataSource) => dataSource.sender
    | EditDataSource(dataSource) => dataSource.sender
    | CreateOracleScript(oracleScript) => oracleScript.sender
    | EditOracleScript(oracleScript) => oracleScript.sender
    | Request(request) => request.sender
    | Report(report) => report.reporter
    | AddOracleAddress(address) => address.validator
    | RemoveOracleAddress(address) => address.validator
    | CreateValidator(validator) => validator.delegatorAddress
    | EditValidator(validator) => validator.sender
    | FailMessage(fail) => fail.sender
    | CreateClient(client) => client.address
    | UpdateClient(client) => client.address
    | SubmitClientMisbehaviour(client) => client.address
    | ConnectionOpenInit(connection) => connection.signer
    | ConnectionOpenTry(connection) => connection.signer
    | ConnectionOpenAck(connection) => connection.signer
    | ConnectionOpenConfirm(connection) => connection.signer
    | ChannelOpenInit(channel) => channel.signer
    | ChannelOpenTry(channel) => channel.signer
    | ChannelOpenAck(channel) => channel.signer
    | ChannelOpenConfirm(channel) => channel.signer
    | ChannelCloseInit(channel) => channel.signer
    | ChannelCloseConfirm(channel) => channel.signer
    | Packet(packet) => packet.sender
    | Acknowledgement(ack) => ack.sender
    | Timeout(timeout) => timeout.sender
    | _ => "" |> Address.fromHex
    };
  };

  type badge_theme_t = {
    text: string,
    textColor: Css.Types.Color.t,
    bgColor: Css.Types.Color.t,
  };

  let getBadgeTheme = msg => {
    switch (msg) {
    | Send(_) => {text: "SEND TOKEN", textColor: Colors.blue7, bgColor: Colors.blue1}
    | CreateDataSource(_) => {
        text: "CREATE DATA SOURCE",
        textColor: Colors.yellow5,
        bgColor: Colors.yellow1,
      }
    | EditDataSource(_) => {
        text: "EDIT DATA SOURCE",
        textColor: Colors.yellow5,
        bgColor: Colors.yellow1,
      }
    | CreateOracleScript(_) => {
        text: "CREATE ORACLE SCRIPT",
        textColor: Colors.pink6,
        bgColor: Colors.pink1,
      }
    | EditOracleScript(_) => {
        text: "EDIT ORACLE SCRIPT",
        textColor: Colors.pink6,
        bgColor: Colors.pink1,
      }
    | Request(_) => {text: "REQUEST", textColor: Colors.orange6, bgColor: Colors.orange1}
    | Report(_) => {text: "REPORT", textColor: Colors.orange6, bgColor: Colors.orange1}
    | AddOracleAddress(_) => {
        text: "ADD ORACLE ADDRESS",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | RemoveOracleAddress(_) => {
        text: "REMOVE ORACLE ADDRESS",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | CreateValidator(_) => {
        text: "CREATE VALIDATOR",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | EditValidator(_) => {
        text: "EDIT VALIDATOR",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | FailMessage(_) => {text: "FAIL", textColor: Colors.red5, bgColor: Colors.red1}
    | CreateClient(_) => {text: "CREATE CLIENT", textColor: Colors.blue7, bgColor: Colors.blue1}
    | UpdateClient(_) => {text: "UPDATE CLIENT", textColor: Colors.blue7, bgColor: Colors.blue1}
    | SubmitClientMisbehaviour(_) => {
        text: "SUBMIT CLIENT MISBEHAVIOUR",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ConnectionOpenInit(_) => {
        text: "CONNECTION OPEN INIT",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ConnectionOpenTry(_) => {
        text: "CONNECTION OPEN TRY",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ConnectionOpenAck(_) => {
        text: "CONNECTION OPEN ACK",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ConnectionOpenConfirm(_) => {
        text: "CONNECTION OPEN CONFIRM",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelOpenInit(_) => {
        text: "CHANNEL OPEN INIT",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelOpenTry(_) => {
        text: "CHANNEL OPEN TRY",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelOpenAck(_) => {
        text: "CHANNEL OPEN ACK",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelOpenConfirm(_) => {
        text: "CHANNEL OPEN CONFIRM",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelCloseInit(_) => {
        text: "CHANNEL CLOSE INIT",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelCloseConfirm(_) => {
        text: "CHANNEL CLOSE CONFIRM",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | Packet(_) => {text: "PACKET", textColor: Colors.blue7, bgColor: Colors.blue1}
    | Acknowledgement(_) => {
        text: "ACKNOWLEDGEMENT",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | Timeout(_) => {text: "TIMEOUT", textColor: Colors.blue7, bgColor: Colors.blue1}
    | _ => {text: "UNKNOWN", textColor: Colors.gray7, bgColor: Colors.gray4}
    };
  };

  let decodeAction = json => {
    JsonUtils.Decode.(
      switch (json |> field("type", string)) {
      | "send" => Send(json |> Send.decode)
      | "create_data_source" => CreateDataSource(json |> CreateDataSource.decode)
      | "edit_data_source" => EditDataSource(json |> EditDataSource.decode)
      | "create_oracle_script" => CreateOracleScript(json |> CreateOracleScript.decode)
      | "edit_oracle_script" => EditOracleScript(json |> EditOracleScript.decode)
      | "request" => Request(json |> Request.decode)
      | "report" => Report(json |> Report.decode)
      | "add_oracle_address" => AddOracleAddress(json |> AddOracleAddress.decode)
      | "remove_oracle_address" => RemoveOracleAddress(json |> RemoveOracleAddress.decode)
      | "create_validator" => CreateValidator(json |> CreateValidator.decode)
      | "edit_validator" => EditValidator(json |> EditValidator.decode)
      | "create_client" => CreateClient(json |> CreateClient.decode)
      | "update_client" => UpdateClient(json |> UpdateClient.decode)
      | "submit_client_misbehaviour" =>
        SubmitClientMisbehaviour(json |> SubmitClientMisbehaviour.decode)
      | "connection_open_init" => ConnectionOpenInit(json |> ConnectionOpenInit.decode)
      | "connection_open_try" => ConnectionOpenTry(json |> ConnectionOpenTry.decode)
      | "connection_open_ack" => ConnectionOpenAck(json |> ConnectionOpenAck.decode)
      | "connection_open_confirm" => ConnectionOpenConfirm(json |> ConnectionOpenConfirm.decode)
      | "channel_open_init" => ChannelOpenInit(json |> ChannelOpenInit.decode)
      | "channel_open_try" => ChannelOpenTry(json |> ChannelOpenTry.decode)
      | "channel_open_ack" => ChannelOpenAck(json |> ChannelOpenAck.decode)
      | "channel_open_confirm" => ChannelOpenConfirm(json |> ChannelOpenConfirm.decode)
      | "channel_close_init" => ChannelCloseInit(json |> ChannelCloseInit.decode)
      | "channel_close_confirm" => ChannelCloseConfirm(json |> ChannelCloseConfirm.decode)
      | "ics04/opaque" => Packet(json |> Packet.decode)
      | "ics04/timeout" => Timeout(json |> Timeout.decode)
      // TODO: handle case correctly
      | "acknowledgement" => Acknowledgement(json |> Acknowledgement.decode)
      | _ => Unknown
      }
    );
  };

  let decodeFailAction = json => FailMessage(json |> FailMessage.decode);

  let decodeUnknowStatus = _ => {
    Js.Console.log("Error status messages field");
    [];
  };

  let decodeActions = json => {
    JsonUtils.Decode.(
      switch (json |> field("status", string)) {
      | "success" => json |> field("messages", list(decodeAction))
      | "failure" => json |> field("messages", list(decodeFailAction))
      | _ => json |> field("messages", decodeUnknowStatus)
      }
    );
  };
};

type t = {
  txHash: Hash.t,
  blockHeight: ID.Block.t,
  success: bool,
  gasFee: list(Coin.t),
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
  subscription Transactions($limit: Int!, $offset: Int!) {
    transactions(offset: $offset, limit: $limit, order_by: {block_height: desc, index: desc}) @bsRecord {
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

module MultiByHeightConfig = [%graphql
  {|
  subscription TransactionsByHeight($height: bigint!, $limit: Int!, $offset: Int!) {
    transactions(where: {block_height: {_eq: $height}}, offset: $offset, limit: $limit, order_by: {index: desc}) @bsRecord {
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

module MultiBySenderConfig = [%graphql
  {|
  subscription TransactionsBySender($sender: String!, $limit: Int!, $offset: Int!) {
    transactions(
      where: {sender: {_eq: $sender}},
      offset: $offset,
      limit: $limit,
      order_by: {block_height: desc,index: desc},
    ) @bsRecord {
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
  subscription TransactionsCount {
    transactions_aggregate {
      aggregate {
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

module TxCountBySenderConfig = [%graphql
  {|
  subscription TransactionsCountBySender($sender: String!) {
    transactions_aggregate(where: {sender: {_eq: $sender}}) {
      aggregate {
        count @bsDecoder(fn: "Belt_Option.getExn")
      }
    }
  }
|}
];

let get = txHash => {
  let (result, _) =
    ApolloHooks.useSubscription(
      SingleConfig.definition,
      ~variables=
        SingleConfig.makeVariables(
          ~tx_hash=txHash |> Hash.toHex |> (x => "\\x" ++ x) |> Js.Json.string,
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
  result |> Sub.map(_, x => x##transactions);
};

let getListBySender = (sender, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiBySenderConfig.definition,
      ~variables=
        MultiBySenderConfig.makeVariables(
          ~sender=sender |> Address.toBech32,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Sub.map(_, x => x##transactions);
};

let getListByBlockHeight = (height, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiByHeightConfig.definition,
      ~variables=
        MultiByHeightConfig.makeVariables(
          ~height=height |> ID.Block.toJson,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Sub.map(_, x => x##transactions);
};

let count = () => {
  let (result, _) = ApolloHooks.useSubscription(TxCountConfig.definition);
  result
  |> Sub.map(_, x => x##transactions_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};

let countBySender = sender => {
  let (result, _) =
    ApolloHooks.useSubscription(
      TxCountBySenderConfig.definition,
      ~variables=TxCountBySenderConfig.makeVariables(~sender=sender |> Address.toBech32, ()),
    );
  result
  |> Sub.map(_, x => x##transactions_aggregate##aggregate |> Belt_Option.getExn |> (y => y##count));
};
