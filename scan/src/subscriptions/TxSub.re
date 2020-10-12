module RawDataReport = {
  type t = {
    externalDataID: int,
    exitCode: int,
    data: JsBuffer.t,
  };

  let decode = json =>
    JsonUtils.Decode.{
      externalDataID: json |> intWithDefault(field("external_id")),
      exitCode: json |> intWithDefault(field("exit_code")),
      data: json |> bufferWithDefault(field("data")),
    };
};

module Msg = {
  type badge_t =
    | SendBadge
    | ReceiveBadge
    | CreateDataSourceBadge
    | EditDataSourceBadge
    | CreateOracleScriptBadge
    | EditOracleScriptBadge
    | RequestBadge
    | ReportBadge
    | AddReporterBadge
    | RemoveReporterBadge
    | CreateValidatorBadge
    | EditValidatorBadge
    | CreateClientBadge
    | UpdateClientBadge
    | SubmitClientMisbehaviourBadge
    | ConnectionOpenInitBadge
    | ConnectionOpenTryBadge
    | ConnectionOpenAckBadge
    | ConnectionOpenConfirmBadge
    | ChannelOpenInitBadge
    | ChannelOpenTryBadge
    | ChannelOpenAckBadge
    | ChannelOpenConfirmBadge
    | ChannelCloseInitBadge
    | ChannelCloseConfirmBadge
    | PacketBadge
    | AcknowledgementBadge
    | TimeoutBadge
    | DelegateBadge
    | UndelegateBadge
    | RedelegateBadge
    | WithdrawRewardBadge
    | UnjailBadge
    | SetWithdrawAddressBadge
    | SubmitProposalBadge
    | DepositBadge
    | VoteBadge
    | WithdrawCommissionBadge
    | MultiSendBadge
    | ActivateBadge
    | UnknownBadge;

  type msg_cat_t =
    | TokenMsg
    | ValidatorMsg
    | ProposalMsg
    | DataMsg
    | UnknownMsg;

  let getBadgeVariantFromString = badge => {
    switch (badge) {
    | "send" => SendBadge
    | "receive" => raise(Not_found)
    | "create_data_source" => CreateDataSourceBadge
    | "edit_data_source" => EditDataSourceBadge
    | "create_oracle_script" => CreateOracleScriptBadge
    | "edit_oracle_script" => EditOracleScriptBadge
    | "request" => RequestBadge
    | "report" => ReportBadge
    | "add_reporter" => AddReporterBadge
    | "remove_reporter" => RemoveReporterBadge
    | "create_validator" => CreateValidatorBadge
    | "edit_validator" => EditValidatorBadge
    | "create_client" => CreateClientBadge
    | "update_client" => UpdateClientBadge
    | "submit_client_misbehaviour" => SubmitClientMisbehaviourBadge
    | "connection_open_init" => ConnectionOpenInitBadge
    | "connection_open_try" => ConnectionOpenTryBadge
    | "connection_open_ack" => ConnectionOpenAckBadge
    | "connection_open_confirm" => ConnectionOpenConfirmBadge
    | "channel_open_init" => ChannelOpenInitBadge
    | "channel_open_try" => ChannelOpenTryBadge
    | "channel_open_ack" => ChannelOpenAckBadge
    | "channel_open_confirm" => ChannelOpenConfirmBadge
    | "channel_close_init" => ChannelCloseInitBadge
    | "channel_close_confirm" => ChannelCloseConfirmBadge
    | "ics04/opaque" => PacketBadge
    | "ics04/timeout" => TimeoutBadge
    | "acknowledgement" => AcknowledgementBadge
    | "delegate" => DelegateBadge
    | "begin_unbonding" => UndelegateBadge
    | "begin_redelegate" => RedelegateBadge
    | "withdraw_delegator_reward" => WithdrawRewardBadge
    | "unjail" => UnjailBadge
    | "set_withdraw_address" => SetWithdrawAddressBadge
    | "submit_proposal" => SubmitProposalBadge
    | "deposit" => DepositBadge
    | "vote" => VoteBadge
    | "withdraw_validator_commission" => WithdrawCommissionBadge
    | "multisend" => MultiSendBadge
    | "activate" => ActivateBadge
    | _ => UnknownBadge
    };
  };

  module Send = {
    type t = {
      fromAddress: Address.t,
      toAddress: Address.t,
      amount: list(Coin.t),
    };

    let decode = json =>
      JsonUtils.Decode.{
        fromAddress: json |> at(["msg", "from_address"], string) |> Address.fromBech32,
        toAddress: json |> at(["msg", "to_address"], string) |> Address.fromBech32,
        amount: json |> at(["msg", "amount"], list(Coin.decodeCoin)),
      };
  };

  module Receive = {
    type t = {
      fromAddress: Address.t,
      toAddress: Address.t,
      amount: list(Coin.t),
    };
  };

  module CreateDataSource = {
    type success_t = {
      id: ID.DataSource.t,
      owner: Address.t,
      name: string,
      executable: JsBuffer.t,
      sender: Address.t,
    };

    type fail_t = {
      owner: Address.t,
      name: string,
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decodeSuccess = json =>
      JsonUtils.Decode.{
        id: json |> at(["extra", "id"], ID.DataSource.fromJson),
        owner: json |> at(["msg", "owner"], string) |> Address.fromBech32,
        name: json |> at(["msg", "name"], string),
        executable: json |> at(["msg", "executable"], string) |> JsBuffer.fromBase64,
        sender: json |> at(["msg", "sender"], string) |> Address.fromBech32,
      };

    let decodeFail = json =>
      JsonUtils.Decode.{
        owner: json |> at(["msg", "owner"], string) |> Address.fromBech32,
        name: json |> at(["msg", "name"], string),
        executable: json |> at(["msg", "executable"], string) |> JsBuffer.fromBase64,
        sender: json |> at(["msg", "sender"], string) |> Address.fromBech32,
      };
  };

  module EditDataSource = {
    type t = {
      id: ID.DataSource.t,
      owner: Address.t,
      name: string,
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> at(["msg", "data_source_id"], ID.DataSource.fromJson),
        owner: json |> at(["msg", "owner"], string) |> Address.fromBech32,
        name: json |> at(["msg", "name"], string),
        executable: json |> at(["msg", "executable"], string) |> JsBuffer.fromBase64,
        sender: json |> at(["msg", "sender"], string) |> Address.fromBech32,
      };
  };

  module CreateOracleScript = {
    type success_t = {
      id: ID.OracleScript.t,
      owner: Address.t,
      name: string,
      code: JsBuffer.t,
      sender: Address.t,
    };

    type fail_t = {
      owner: Address.t,
      name: string,
      code: JsBuffer.t,
      sender: Address.t,
    };

    let decodeSuccess = json =>
      JsonUtils.Decode.{
        id: json |> at(["extra", "id"], ID.OracleScript.fromJson),
        owner: json |> at(["msg", "owner"], string) |> Address.fromBech32,
        name: json |> at(["msg", "name"], string),
        code: json |> at(["msg", "code"], string) |> JsBuffer.fromBase64,
        sender: json |> at(["msg", "sender"], string) |> Address.fromBech32,
      };

    let decodeFail = json =>
      JsonUtils.Decode.{
        owner: json |> at(["msg", "owner"], string) |> Address.fromBech32,
        name: json |> at(["msg", "name"], string),
        code: json |> at(["msg", "code"], string) |> JsBuffer.fromBase64,
        sender: json |> at(["msg", "sender"], string) |> Address.fromBech32,
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
        id: json |> at(["msg", "oracle_script_id"], ID.OracleScript.fromJson),
        owner: json |> at(["msg", "owner"], string) |> Address.fromBech32,
        name: json |> at(["msg", "name"], string),
        code: json |> at(["msg", "code"], string) |> JsBuffer.fromBase64,
        sender: json |> at(["msg", "sender"], string) |> Address.fromBech32,
      };
  };

  module Request = {
    type success_t = {
      id: ID.Request.t,
      oracleScriptID: ID.OracleScript.t,
      oracleScriptName: string,
      calldata: JsBuffer.t,
      askCount: int,
      minCount: int,
      schema: string,
      sender: Address.t,
    };

    type fail_t = {
      oracleScriptID: ID.OracleScript.t,
      calldata: JsBuffer.t,
      askCount: int,
      minCount: int,
      sender: Address.t,
    };

    let decodeSuccess = json => {
      JsonUtils.Decode.{
        id: json |> at(["extra", "id"], ID.Request.fromJson),
        oracleScriptID: json |> at(["msg", "oracle_script_id"], ID.OracleScript.fromJson),
        oracleScriptName: json |> at(["extra", "name"], string),
        calldata: json |> bufferWithDefault(at(["msg", "calldata"])),
        askCount: json |> at(["msg", "ask_count"], int),
        minCount: json |> at(["msg", "min_count"], int),
        schema: json |> at(["extra", "schema"], string),
        sender: json |> at(["msg", "sender"], string) |> Address.fromBech32,
      };
    };

    let decodeFail = json => {
      JsonUtils.Decode.{
        oracleScriptID: json |> at(["msg", "oracle_script_id"], ID.OracleScript.fromJson),
        calldata: json |> bufferWithDefault(at(["msg", "calldata"])),
        askCount: json |> at(["msg", "ask_count"], int),
        minCount: json |> at(["msg", "min_count"], int),
        sender: json |> at(["msg", "sender"], string) |> Address.fromBech32,
      };
    };
  };

  module Report = {
    type t = {
      requestID: ID.Request.t,
      rawReports: list(RawDataReport.t),
      validator: Address.t,
      reporter: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        requestID: json |> at(["msg", "request_id"], ID.Request.fromJson),
        rawReports: json |> at(["msg", "raw_reports"], list(RawDataReport.decode)),
        validator: json |> at(["msg", "validator"], string) |> Address.fromBech32,
        reporter: json |> at(["msg", "reporter"], string) |> Address.fromBech32,
      };
  };
  module AddReporter = {
    type success_t = {
      validator: Address.t,
      reporter: Address.t,
      validatorMoniker: string,
    };

    type fail_t = {
      validator: Address.t,
      reporter: Address.t,
    };

    let decodeSuccess = json =>
      JsonUtils.Decode.{
        validator: json |> at(["msg", "validator"], string) |> Address.fromBech32,
        reporter: json |> at(["msg", "reporter"], string) |> Address.fromBech32,
        validatorMoniker: json |> at(["extra", "validator_moniker"], string),
      };

    let decodeFail = json =>
      JsonUtils.Decode.{
        validator: json |> at(["msg", "validator"], string) |> Address.fromBech32,
        reporter: json |> at(["msg", "reporter"], string) |> Address.fromBech32,
      };
  };

  module RemoveReporter = {
    type success_t = {
      validator: Address.t,
      reporter: Address.t,
      validatorMoniker: string,
    };

    type fail_t = {
      validator: Address.t,
      reporter: Address.t,
    };

    let decodeSuccess = json =>
      JsonUtils.Decode.{
        validator: json |> at(["msg", "validator"], string) |> Address.fromBech32,
        reporter: json |> at(["msg", "reporter"], string) |> Address.fromBech32,
        validatorMoniker: json |> at(["extra", "validator_moniker"], string),
      };

    let decodeFail = json =>
      JsonUtils.Decode.{
        validator: json |> at(["msg", "validator"], string) |> Address.fromBech32,
        reporter: json |> at(["msg", "reporter"], string) |> Address.fromBech32,
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
      minSelfDelegation: Coin.t,
      selfDelegation: Coin.t,
    };
    let decode = json =>
      JsonUtils.Decode.{
        moniker: json |> at(["msg", "description", "moniker"], string),
        identity: json |> at(["msg", "description", "identity"], string),
        website: json |> at(["msg", "description", "website"], string),
        details: json |> at(["msg", "description", "details"], string),
        commissionRate: json |> at(["msg", "commission", "rate"], floatstr),
        commissionMaxRate: json |> at(["msg", "commission", "max_rate"], floatstr),
        commissionMaxChange: json |> at(["msg", "commission", "max_change_rate"], floatstr),
        delegatorAddress: json |> at(["msg", "delegator_address"], string) |> Address.fromBech32,
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
        publicKey: json |> at(["msg", "pubkey"], string) |> PubKey.fromBech32,
        minSelfDelegation:
          json |> at(["msg", "min_self_delegation"], floatstr) |> Coin.newUBANDFromAmount,
        selfDelegation: json |> at(["msg", "value"], Coin.decodeCoin),
      };
  };

  module EditValidator = {
    type t = {
      moniker: string,
      identity: string,
      website: string,
      details: string,
      commissionRate: option(float),
      sender: Address.t,
      minSelfDelegation: option(Coin.t),
    };

    let decode = json => {
      exception WrongNetwork(string);
      switch (Env.network) {
      | "GUANYU"
      | "GUANYU38" =>
        JsonUtils.Decode.{
          moniker: json |> at(["msg", "description", "moniker"], string),
          identity: json |> at(["msg", "description", "identity"], string),
          website: json |> at(["msg", "description", "website"], string),
          details: json |> at(["msg", "description", "details"], string),
          commissionRate: json |> optional(at(["msg", "commission_rate"], floatstr)),
          sender: json |> at(["msg", "address"], string) |> Address.fromBech32,
          minSelfDelegation:
            json
            |> optional(at(["msg", "min_self_delegation"], floatstr))
            |> Belt.Option.map(_, Coin.newUBANDFromAmount),
        }
      | "WENCHANG" =>
        JsonUtils.Decode.{
          moniker: json |> at(["msg", "moniker"], string),
          identity: json |> at(["msg", "identity"], string),
          website: json |> at(["msg", "website"], string),
          details: json |> at(["msg", "details"], string),
          commissionRate: json |> optional(at(["msg", "commission_rate"], floatstr)),
          sender: json |> at(["msg", "address"], string) |> Address.fromBech32,
          minSelfDelegation:
            json
            |> optional(at(["msg", "min_self_delegation"], floatstr))
            |> Belt.Option.map(_, Coin.newUBANDFromAmount),
        }
      | _ => raise(WrongNetwork("Incorrect or unspecified NETWORK environment variable"))
      };
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
      timeoutHeight: int,
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
          timeoutHeight: json |> at(["packet", "timeout_height"], int),
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
          timeoutHeight: 999999,
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
          timeoutHeight: json |> at(["packet", "timeout_height"], int),
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
    };
    let decode = json =>
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        common: json |> ConnectionCommon.decode,
        clientID: json |> field("client_id", string),
      };
  };

  module ConnectionOpenTry = {
    type t = {
      signer: Address.t,
      common: ConnectionCommon.t,
      clientID: string,
    };
    let decode = json =>
      JsonUtils.Decode.{
        signer: json |> field("signer", string) |> Address.fromBech32,
        clientID: json |> field("client_id", string),
        common: json |> ConnectionCommon.decode,
      };
  };

  module ConnectionOpenAck = {
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

  module Delegate = {
    type t = {
      validatorAddress: Address.t,
      delegatorAddress: Address.t,
      amount: Coin.t,
    };

    let decode = json => {
      JsonUtils.Decode.{
        delegatorAddress: json |> at(["msg", "delegator_address"], string) |> Address.fromBech32,
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
        amount: json |> at(["msg", "amount"], Coin.decodeCoin),
      };
    };
  };

  module Undelegate = {
    type t = {
      validatorAddress: Address.t,
      delegatorAddress: Address.t,
      amount: Coin.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        delegatorAddress: json |> at(["msg", "delegator_address"], string) |> Address.fromBech32,
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
        amount: json |> at(["msg", "amount"], Coin.decodeCoin),
      };
    };
  };

  module Redelegate = {
    type t = {
      validatorSourceAddress: Address.t,
      validatorDestinationAddress: Address.t,
      delegatorAddress: Address.t,
      amount: Coin.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        validatorSourceAddress:
          json |> at(["msg", "validator_src_address"], string) |> Address.fromBech32,
        validatorDestinationAddress:
          json |> at(["msg", "validator_dst_address"], string) |> Address.fromBech32,
        delegatorAddress: json |> at(["msg", "delegator_address"], string) |> Address.fromBech32,
        amount: json |> at(["msg", "amount"], Coin.decodeCoin),
      };
    };
  };

  module WithdrawReward = {
    type success_t = {
      validatorAddress: Address.t,
      delegatorAddress: Address.t,
      amount: list(Coin.t),
    };
    type fail_t = {
      validatorAddress: Address.t,
      delegatorAddress: Address.t,
    };

    let decodeSuccess = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
        delegatorAddress: json |> at(["msg", "delegator_address"], string) |> Address.fromBech32,
        amount: json |> at(["extra", "reward_amount"], string) |> GraphQLParser.coins,
      };
    };

    let decodeFail = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
        delegatorAddress: json |> at(["msg", "delegator_address"], string) |> Address.fromBech32,
      };
    };
  };

  module Unjail = {
    type t = {address: Address.t};

    let decode = json =>
      JsonUtils.Decode.{address: json |> at(["msg", "address"], string) |> Address.fromBech32};
  };
  module SetWithdrawAddress = {
    type t = {
      delegatorAddress: Address.t,
      withdrawAddress: Address.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        delegatorAddress: json |> at(["msg", "delegator_address"], string) |> Address.fromBech32,
        withdrawAddress: json |> at(["msg", "withdraw_address"], string) |> Address.fromBech32,
      };
    };
  };
  module SubmitProposal = {
    type t = {
      proposer: Address.t,
      title: string,
      description: string,
      initialDeposit: list(Coin.t),
    };
    let decode = json => {
      JsonUtils.Decode.{
        proposer: json |> at(["msg", "proposer"], string) |> Address.fromBech32,
        title: json |> at(["msg", "content", "title"], string),
        description: json |> at(["msg", "content", "description"], string),
        initialDeposit: json |> at(["msg", "initial_deposit"], list(Coin.decodeCoin)),
      };
    };
  };

  module Deposit = {
    type t = {
      depositor: Address.t,
      proposalID: int,
      amount: list(Coin.t),
    };
    let decode = json => {
      JsonUtils.Decode.{
        depositor: json |> at(["msg", "depositor"], string) |> Address.fromBech32,
        proposalID: json |> at(["msg", "proposal_id"], int),
        amount: json |> at(["msg", "amount"], list(Coin.decodeCoin)),
      };
    };
  };
  module Vote = {
    type t = {
      voterAddress: Address.t,
      proposalID: int,
      option: string,
    };
    let decode = json => {
      JsonUtils.Decode.{
        voterAddress: json |> at(["msg", "voter"], string) |> Address.fromBech32,
        proposalID: json |> at(["msg", "proposal_id"], int),
        option: json |> at(["msg", "option"], string),
      };
    };
  };
  module WithdrawCommission = {
    type success_t = {
      validatorAddress: Address.t,
      amount: list(Coin.t),
    };
    type fail_t = {validatorAddress: Address.t};
    let decodeSuccess = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
        amount: json |> at(["extra", "commission_amount"], string) |> GraphQLParser.coins,
      };
    };

    let decodeFail = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
      };
    };
  };

  module MultiSend = {
    type send_tx_t = {
      address: Address.t,
      coins: list(Coin.t),
    };
    type t = {
      inputs: list(send_tx_t),
      outputs: list(send_tx_t),
    };
    let decodeSendTx = json => {
      JsonUtils.Decode.{
        address: json |> field("address", string) |> Address.fromBech32,
        coins: json |> field("coins", list(Coin.decodeCoin)),
      };
    };
    let decode = json => {
      JsonUtils.Decode.{
        inputs: json |> at(["msg", "inputs"], list(decodeSendTx)),
        outputs: json |> at(["msg", "outputs"], list(decodeSendTx)),
      };
    };
  };

  module Activate = {
    type t = {validatorAddress: Address.t};

    let decode = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> at(["msg", "validator"], string) |> Address.fromBech32,
      };
    };
  };

  type t =
    | SendMsgSuccess(Send.t)
    | SendMsgFail(Send.t)
    | ReceiveMsg(Receive.t)
    | CreateDataSourceMsgSuccess(CreateDataSource.success_t)
    | CreateDataSourceMsgFail(CreateDataSource.fail_t)
    | EditDataSourceMsgSuccess(EditDataSource.t)
    | EditDataSourceMsgFail(EditDataSource.t)
    | CreateOracleScriptMsgSuccess(CreateOracleScript.success_t)
    | CreateOracleScriptMsgFail(CreateOracleScript.fail_t)
    | EditOracleScriptMsgSuccess(EditOracleScript.t)
    | EditOracleScriptMsgFail(EditOracleScript.t)
    | RequestMsgSuccess(Request.success_t)
    | RequestMsgFail(Request.fail_t)
    | ReportMsgSuccess(Report.t)
    | ReportMsgFail(Report.t)
    | AddReporterMsgSuccess(AddReporter.success_t)
    | AddReporterMsgFail(AddReporter.fail_t)
    | RemoveReporterMsgSuccess(RemoveReporter.success_t)
    | RemoveReporterMsgFail(RemoveReporter.fail_t)
    | CreateValidatorMsgSuccess(CreateValidator.t)
    | CreateValidatorMsgFail(CreateValidator.t)
    | EditValidatorMsgSuccess(EditValidator.t)
    | EditValidatorMsgFail(EditValidator.t)
    | CreateClientMsg(CreateClient.t)
    | UpdateClientMsg(UpdateClient.t)
    | SubmitClientMisbehaviourMsg(SubmitClientMisbehaviour.t)
    | ConnectionOpenInitMsg(ConnectionOpenInit.t)
    | ConnectionOpenTryMsg(ConnectionOpenTry.t)
    | ConnectionOpenAckMsg(ConnectionOpenAck.t)
    | ConnectionOpenConfirmMsg(ConnectionOpenConfirm.t)
    | ChannelOpenInitMsg(ChannelOpenInit.t)
    | ChannelOpenTryMsg(ChannelOpenTry.t)
    | ChannelOpenAckMsg(ChannelOpenAck.t)
    | ChannelOpenConfirmMsg(ChannelOpenConfirm.t)
    | ChannelCloseInitMsg(ChannelCloseInit.t)
    | ChannelCloseConfirmMsg(ChannelCloseConfirm.t)
    | PacketMsg(Packet.t)
    | AcknowledgementMsg(Acknowledgement.t)
    | TimeoutMsg(Timeout.t)
    | DelegateMsgSuccess(Delegate.t)
    | DelegateMsgFail(Delegate.t)
    | UndelegateMsgSuccess(Undelegate.t)
    | UndelegateMsgFail(Undelegate.t)
    | RedelegateMsgSuccess(Redelegate.t)
    | RedelegateMsgFail(Redelegate.t)
    | WithdrawRewardMsgSuccess(WithdrawReward.success_t)
    | WithdrawRewardMsgFail(WithdrawReward.fail_t)
    | UnjailMsgSuccess(Unjail.t)
    | UnjailMsgFail(Unjail.t)
    | SetWithdrawAddressMsgSuccess(SetWithdrawAddress.t)
    | SetWithdrawAddressMsgFail(SetWithdrawAddress.t)
    | SubmitProposalMsgSuccess(SubmitProposal.t)
    | SubmitProposalMsgFail(SubmitProposal.t)
    | DepositMsgSuccess(Deposit.t)
    | DepositMsgFail(Deposit.t)
    | VoteMsgSuccess(Vote.t)
    | VoteMsgFail(Vote.t)
    | WithdrawCommissionMsgSuccess(WithdrawCommission.success_t)
    | WithdrawCommissionMsgFail(WithdrawCommission.fail_t)
    | MultiSendMsgSuccess(MultiSend.t)
    | MultiSendMsgFail(MultiSend.t)
    | ActivateMsgSuccess(Activate.t)
    | ActivateMsgFail(Activate.t)
    | UnknownMsg;

  let getCreator = msg => {
    switch (msg) {
    | ReceiveMsg(receive) => receive.fromAddress
    | SendMsgSuccess(send)
    | SendMsgFail(send) => send.fromAddress
    | CreateDataSourceMsgSuccess(dataSource) => dataSource.sender
    | CreateDataSourceMsgFail(dataSource) => dataSource.sender
    | EditDataSourceMsgSuccess(dataSource) => dataSource.sender
    | EditDataSourceMsgFail(dataSource) => dataSource.sender
    | CreateOracleScriptMsgSuccess(oracleScript) => oracleScript.sender
    | CreateOracleScriptMsgFail(oracleScript) => oracleScript.sender
    | EditOracleScriptMsgSuccess(oracleScript) => oracleScript.sender
    | EditOracleScriptMsgFail(oracleScript) => oracleScript.sender
    | RequestMsgSuccess(request) => request.sender
    | RequestMsgFail(request) => request.sender
    | ReportMsgSuccess(report)
    | ReportMsgFail(report) => report.reporter
    | AddReporterMsgSuccess(address) => address.validator
    | AddReporterMsgFail(address) => address.validator
    | RemoveReporterMsgSuccess(address) => address.validator
    | RemoveReporterMsgFail(address) => address.validator
    | CreateValidatorMsgSuccess(validator)
    | CreateValidatorMsgFail(validator) => validator.delegatorAddress
    | EditValidatorMsgSuccess(validator)
    | EditValidatorMsgFail(validator) => validator.sender
    | DelegateMsgSuccess(delegation)
    | DelegateMsgFail(delegation) => delegation.delegatorAddress
    | UndelegateMsgSuccess(delegation)
    | UndelegateMsgFail(delegation) => delegation.delegatorAddress
    | RedelegateMsgSuccess(delegation)
    | RedelegateMsgFail(delegation) => delegation.delegatorAddress
    | WithdrawRewardMsgSuccess(withdrawal) => withdrawal.delegatorAddress
    | WithdrawRewardMsgFail(withdrawal) => withdrawal.delegatorAddress
    | UnjailMsgSuccess(validator)
    | UnjailMsgFail(validator) => validator.address
    | SetWithdrawAddressMsgSuccess(set)
    | SetWithdrawAddressMsgFail(set) => set.delegatorAddress
    | SubmitProposalMsgSuccess(proposal)
    | SubmitProposalMsgFail(proposal) => proposal.proposer
    | DepositMsgSuccess(deposit)
    | DepositMsgFail(deposit) => deposit.depositor
    | VoteMsgSuccess(vote)
    | VoteMsgFail(vote) => vote.voterAddress
    | WithdrawCommissionMsgSuccess(withdrawal) => withdrawal.validatorAddress
    | WithdrawCommissionMsgFail(withdrawal) => withdrawal.validatorAddress
    | MultiSendMsgSuccess(tx)
    | MultiSendMsgFail(tx) =>
      let firstInput = tx.inputs |> Belt_List.getExn(_, 0);
      firstInput.address;
    | ActivateMsgSuccess(activator)
    | ActivateMsgFail(activator) => activator.validatorAddress
    //TODO: Revisit IBC msg
    | CreateClientMsg(client) => client.address
    | UpdateClientMsg(client) => client.address
    | SubmitClientMisbehaviourMsg(client) => client.address
    | ConnectionOpenInitMsg(connection) => connection.signer
    | ConnectionOpenTryMsg(connection) => connection.signer
    | ConnectionOpenAckMsg(connection) => connection.signer
    | ConnectionOpenConfirmMsg(connection) => connection.signer
    | ChannelOpenInitMsg(channel) => channel.signer
    | ChannelOpenTryMsg(channel) => channel.signer
    | ChannelOpenAckMsg(channel) => channel.signer
    | ChannelOpenConfirmMsg(channel) => channel.signer
    | ChannelCloseInitMsg(channel) => channel.signer
    | ChannelCloseConfirmMsg(channel) => channel.signer
    | PacketMsg(packet) => packet.sender
    | AcknowledgementMsg(ack) => ack.sender
    | TimeoutMsg(timeout) => timeout.sender
    | _ => "" |> Address.fromHex
    };
  };

  type badge_theme_t = {
    name: string,
    category: msg_cat_t,
  };

  let getBadge = badgeVariant => {
    switch (badgeVariant) {
    | SendBadge => {name: "SEND", category: TokenMsg}
    | ReceiveBadge => {name: "RECEIVE", category: TokenMsg}
    | CreateDataSourceBadge => {name: "CREATE DATA SOURCE", category: DataMsg}
    | EditDataSourceBadge => {name: "EDIT DATA SOURCE", category: DataMsg}
    | CreateOracleScriptBadge => {name: "CREATE ORACLE SCRIPT", category: DataMsg}
    | EditOracleScriptBadge => {name: "EDIT ORACLE SCRIPT", category: DataMsg}
    | RequestBadge => {name: "REQUEST", category: DataMsg}
    | ReportBadge => {name: "REPORT", category: DataMsg}
    | AddReporterBadge => {name: "ADD REPORTER", category: ValidatorMsg}
    | RemoveReporterBadge => {name: "REMOVE REPORTER", category: ValidatorMsg}
    | CreateValidatorBadge => {name: "CREATE VALIDATOR", category: ValidatorMsg}
    | EditValidatorBadge => {name: "EDIT VALIDATOR", category: ValidatorMsg}
    | DelegateBadge => {name: "DELEGATE", category: TokenMsg}
    | UndelegateBadge => {name: "UNDELEGATE", category: TokenMsg}
    | RedelegateBadge => {name: "REDELEGATE", category: TokenMsg}
    | VoteBadge => {name: "VOTE", category: ProposalMsg}
    | WithdrawRewardBadge => {name: "WITHDRAW REWARD", category: TokenMsg}
    | UnjailBadge => {name: "UNJAIL", category: ValidatorMsg}
    | SetWithdrawAddressBadge => {name: "SET WITHDRAW ADDRESS", category: ValidatorMsg}
    | SubmitProposalBadge => {name: "SUBMIT PROPOSAL", category: ProposalMsg}
    | DepositBadge => {name: "DEPOSIT", category: ProposalMsg}
    | WithdrawCommissionBadge => {name: "WITHDRAW COMMISSION", category: TokenMsg}
    | MultiSendBadge => {name: "MULTI SEND", category: TokenMsg}
    | ActivateBadge => {name: "ACTIVATE", category: ValidatorMsg}
    | UnknownBadge => {name: "UNKNOWN", category: TokenMsg}
    //TODO: Revisit IBC msg
    | CreateClientBadge => {name: "CREATE CLIENT", category: TokenMsg}
    | UpdateClientBadge => {name: "UPDATE CLIENT", category: TokenMsg}
    | SubmitClientMisbehaviourBadge => {name: "SUBMIT CLIENT MISBEHAVIOUR", category: TokenMsg}
    | ConnectionOpenInitBadge => {name: "CONNECTION OPEN INIT", category: TokenMsg}
    | ConnectionOpenTryBadge => {name: "CONNECTION OPEN TRY", category: TokenMsg}
    | ConnectionOpenAckBadge => {name: "CONNECTION OPEN ACK", category: TokenMsg}
    | ConnectionOpenConfirmBadge => {name: "CONNECTION OPEN CONFIRM", category: TokenMsg}
    | ChannelOpenInitBadge => {name: "CHANNEL OPEN INIT", category: TokenMsg}
    | ChannelOpenTryBadge => {name: "CHANNEL OPEN TRY", category: TokenMsg}
    | ChannelOpenAckBadge => {name: "CHANNEL OPEN ACK", category: TokenMsg}
    | ChannelOpenConfirmBadge => {name: "CHANNEL OPEN CONFIRM", category: TokenMsg}
    | ChannelCloseInitBadge => {name: "CHANNEL CLOSE INIT", category: TokenMsg}
    | ChannelCloseConfirmBadge => {name: "CHANNEL CLOSE CONFIRM", category: TokenMsg}
    | PacketBadge => {name: "PACKET", category: TokenMsg}
    | AcknowledgementBadge => {name: "ACKNOWLEDGEMENT", category: TokenMsg}
    | TimeoutBadge => {name: "TIMEOUT", category: TokenMsg}
    };
  };

  let getBadgeTheme = msg => {
    switch (msg) {
    | SendMsgSuccess(_)
    | SendMsgFail(_) => getBadge(SendBadge)
    | ReceiveMsg(_) => getBadge(ReceiveBadge)
    | CreateDataSourceMsgSuccess(_)
    | CreateDataSourceMsgFail(_) => getBadge(CreateDataSourceBadge)
    | EditDataSourceMsgSuccess(_)
    | EditDataSourceMsgFail(_) => getBadge(EditDataSourceBadge)
    | CreateOracleScriptMsgSuccess(_)
    | CreateOracleScriptMsgFail(_) => getBadge(CreateOracleScriptBadge)
    | EditOracleScriptMsgSuccess(_)
    | EditOracleScriptMsgFail(_) => getBadge(EditOracleScriptBadge)
    | RequestMsgSuccess(_)
    | RequestMsgFail(_) => getBadge(RequestBadge)
    | ReportMsgSuccess(_)
    | ReportMsgFail(_) => getBadge(ReportBadge)
    | AddReporterMsgSuccess(_)
    | AddReporterMsgFail(_) => getBadge(AddReporterBadge)
    | RemoveReporterMsgSuccess(_)
    | RemoveReporterMsgFail(_) => getBadge(RemoveReporterBadge)
    | CreateValidatorMsgSuccess(_)
    | CreateValidatorMsgFail(_) => getBadge(CreateValidatorBadge)
    | EditValidatorMsgSuccess(_)
    | EditValidatorMsgFail(_) => getBadge(EditValidatorBadge)
    | DelegateMsgSuccess(_)
    | DelegateMsgFail(_) => getBadge(DelegateBadge)
    | UndelegateMsgSuccess(_)
    | UndelegateMsgFail(_) => getBadge(UndelegateBadge)
    | RedelegateMsgSuccess(_)
    | RedelegateMsgFail(_) => getBadge(RedelegateBadge)
    | VoteMsgSuccess(_)
    | VoteMsgFail(_) => getBadge(VoteBadge)
    | WithdrawRewardMsgSuccess(_)
    | WithdrawRewardMsgFail(_) => getBadge(WithdrawRewardBadge)
    | UnjailMsgSuccess(_)
    | UnjailMsgFail(_) => getBadge(UnjailBadge)
    | SetWithdrawAddressMsgSuccess(_)
    | SetWithdrawAddressMsgFail(_) => getBadge(SetWithdrawAddressBadge)
    | SubmitProposalMsgSuccess(_)
    | SubmitProposalMsgFail(_) => getBadge(SubmitProposalBadge)
    | DepositMsgSuccess(_)
    | DepositMsgFail(_) => getBadge(DepositBadge)
    | WithdrawCommissionMsgSuccess(_)
    | WithdrawCommissionMsgFail(_) => getBadge(WithdrawCommissionBadge)
    | MultiSendMsgSuccess(_) => getBadge(MultiSendBadge)
    | MultiSendMsgFail(_) => getBadge(MultiSendBadge)
    | ActivateMsgSuccess(_) => getBadge(ActivateBadge)
    | ActivateMsgFail(_) => getBadge(ActivateBadge)
    | UnknownMsg => getBadge(UnknownBadge)
    //TODO: Revisit IBC msg
    | CreateClientMsg(_) => getBadge(CreateClientBadge)
    | UpdateClientMsg(_) => getBadge(UpdateClientBadge)
    | SubmitClientMisbehaviourMsg(_) => getBadge(SubmitClientMisbehaviourBadge)
    | ConnectionOpenInitMsg(_) => getBadge(ConnectionOpenInitBadge)
    | ConnectionOpenTryMsg(_) => getBadge(ConnectionOpenTryBadge)
    | ConnectionOpenAckMsg(_) => getBadge(ConnectionOpenAckBadge)
    | ConnectionOpenConfirmMsg(_) => getBadge(ConnectionOpenConfirmBadge)
    | ChannelOpenInitMsg(_) => getBadge(ChannelOpenInitBadge)
    | ChannelOpenTryMsg(_) => getBadge(ChannelOpenTryBadge)
    | ChannelOpenAckMsg(_) => getBadge(ChannelOpenAckBadge)
    | ChannelOpenConfirmMsg(_) => getBadge(ChannelOpenConfirmBadge)
    | ChannelCloseInitMsg(_) => getBadge(ChannelCloseInitBadge)
    | ChannelCloseConfirmMsg(_) => getBadge(ChannelCloseConfirmBadge)
    | PacketMsg(_) => getBadge(PacketBadge)
    | AcknowledgementMsg(_) => getBadge(AcknowledgementBadge)
    | TimeoutMsg(_) => getBadge(TimeoutBadge)
    };
  };

  let decodeAction = json => {
    JsonUtils.Decode.(
      switch (json |> field("type", string) |> getBadgeVariantFromString) {
      | SendBadge => SendMsgSuccess(json |> Send.decode)
      | ReceiveBadge => raise(Not_found)
      | CreateDataSourceBadge =>
        CreateDataSourceMsgSuccess(json |> CreateDataSource.decodeSuccess)
      | EditDataSourceBadge => EditDataSourceMsgSuccess(json |> EditDataSource.decode)
      | CreateOracleScriptBadge =>
        CreateOracleScriptMsgSuccess(json |> CreateOracleScript.decodeSuccess)
      | EditOracleScriptBadge => EditOracleScriptMsgSuccess(json |> EditOracleScript.decode)
      | RequestBadge => RequestMsgSuccess(json |> Request.decodeSuccess)
      | ReportBadge => ReportMsgSuccess(json |> Report.decode)
      | AddReporterBadge => AddReporterMsgSuccess(json |> AddReporter.decodeSuccess)
      | RemoveReporterBadge => RemoveReporterMsgSuccess(json |> RemoveReporter.decodeSuccess)
      | CreateValidatorBadge => CreateValidatorMsgSuccess(json |> CreateValidator.decode)
      | EditValidatorBadge => EditValidatorMsgSuccess(json |> EditValidator.decode)
      | DelegateBadge => DelegateMsgSuccess(json |> Delegate.decode)
      | UndelegateBadge => UndelegateMsgSuccess(json |> Undelegate.decode)
      | RedelegateBadge => RedelegateMsgSuccess(json |> Redelegate.decode)
      | WithdrawRewardBadge => WithdrawRewardMsgSuccess(json |> WithdrawReward.decodeSuccess)
      | UnjailBadge => UnjailMsgSuccess(json |> Unjail.decode)
      | SetWithdrawAddressBadge => SetWithdrawAddressMsgSuccess(json |> SetWithdrawAddress.decode)
      | SubmitProposalBadge => SubmitProposalMsgSuccess(json |> SubmitProposal.decode)
      | DepositBadge => DepositMsgSuccess(json |> Deposit.decode)
      | VoteBadge => VoteMsgSuccess(json |> Vote.decode)
      | WithdrawCommissionBadge =>
        WithdrawCommissionMsgSuccess(json |> WithdrawCommission.decodeSuccess)
      | MultiSendBadge => MultiSendMsgSuccess(json |> MultiSend.decode)
      | ActivateBadge => ActivateMsgSuccess(json |> Activate.decode)
      | UnknownBadge => UnknownMsg
      //TODO: Revisit IBC msg
      | CreateClientBadge => CreateClientMsg(json |> CreateClient.decode)
      | UpdateClientBadge => UpdateClientMsg(json |> UpdateClient.decode)
      | SubmitClientMisbehaviourBadge =>
        SubmitClientMisbehaviourMsg(json |> SubmitClientMisbehaviour.decode)
      | ConnectionOpenInitBadge => ConnectionOpenInitMsg(json |> ConnectionOpenInit.decode)
      | ConnectionOpenTryBadge => ConnectionOpenTryMsg(json |> ConnectionOpenTry.decode)
      | ConnectionOpenAckBadge => ConnectionOpenAckMsg(json |> ConnectionOpenAck.decode)
      | ConnectionOpenConfirmBadge =>
        ConnectionOpenConfirmMsg(json |> ConnectionOpenConfirm.decode)
      | ChannelOpenInitBadge => ChannelOpenInitMsg(json |> ChannelOpenInit.decode)
      | ChannelOpenTryBadge => ChannelOpenTryMsg(json |> ChannelOpenTry.decode)
      | ChannelOpenAckBadge => ChannelOpenAckMsg(json |> ChannelOpenAck.decode)
      | ChannelOpenConfirmBadge => ChannelOpenConfirmMsg(json |> ChannelOpenConfirm.decode)
      | ChannelCloseInitBadge => ChannelCloseInitMsg(json |> ChannelCloseInit.decode)
      | ChannelCloseConfirmBadge => ChannelCloseConfirmMsg(json |> ChannelCloseConfirm.decode)
      | PacketBadge => PacketMsg(json |> Packet.decode)
      | TimeoutBadge => TimeoutMsg(json |> Timeout.decode)
      // TODO: handle case correctly
      | AcknowledgementBadge => AcknowledgementMsg(json |> Acknowledgement.decode)
      }
    );
  };

  let decodeFailAction = json => {
    JsonUtils.Decode.(
      switch (json |> field("type", string) |> getBadgeVariantFromString) {
      | SendBadge => SendMsgFail(json |> Send.decode)
      | ReceiveBadge => raise(Not_found)
      | CreateDataSourceBadge => CreateDataSourceMsgFail(json |> CreateDataSource.decodeFail)
      | EditDataSourceBadge => EditDataSourceMsgFail(json |> EditDataSource.decode)
      | CreateOracleScriptBadge =>
        CreateOracleScriptMsgFail(json |> CreateOracleScript.decodeFail)
      | EditOracleScriptBadge => EditOracleScriptMsgFail(json |> EditOracleScript.decode)
      | RequestBadge => RequestMsgFail(json |> Request.decodeFail)
      | ReportBadge => ReportMsgFail(json |> Report.decode)
      | AddReporterBadge => AddReporterMsgFail(json |> AddReporter.decodeFail)
      | RemoveReporterBadge => RemoveReporterMsgFail(json |> RemoveReporter.decodeFail)
      | CreateValidatorBadge => CreateValidatorMsgFail(json |> CreateValidator.decode)
      | EditValidatorBadge => EditValidatorMsgFail(json |> EditValidator.decode)
      | DelegateBadge => DelegateMsgFail(json |> Delegate.decode)
      | UndelegateBadge => UndelegateMsgFail(json |> Undelegate.decode)
      | RedelegateBadge => RedelegateMsgFail(json |> Redelegate.decode)
      | WithdrawRewardBadge => WithdrawRewardMsgFail(json |> WithdrawReward.decodeFail)
      | UnjailBadge => UnjailMsgFail(json |> Unjail.decode)
      | SetWithdrawAddressBadge => SetWithdrawAddressMsgFail(json |> SetWithdrawAddress.decode)
      | SubmitProposalBadge => SubmitProposalMsgFail(json |> SubmitProposal.decode)
      | DepositBadge => DepositMsgFail(json |> Deposit.decode)
      | VoteBadge => VoteMsgFail(json |> Vote.decode)
      | WithdrawCommissionBadge =>
        WithdrawCommissionMsgFail(json |> WithdrawCommission.decodeFail)
      | MultiSendBadge => MultiSendMsgFail(json |> MultiSend.decode)
      | ActivateBadge => ActivateMsgFail(json |> Activate.decode)
      | UnknownBadge => UnknownMsg
      | _ => UnknownMsg
      }
    );
  };
};

type block_t = {timestamp: MomentRe.Moment.t};

type t = {
  id: int,
  txHash: Hash.t,
  blockHeight: ID.Block.t,
  success: bool,
  gasFee: list(Coin.t),
  gasLimit: int,
  gasUsed: int,
  sender: Address.t,
  timestamp: MomentRe.Moment.t,
  messages: list(Msg.t),
  memo: string,
  errMsg: string,
};

type internal_t = {
  id: int,
  txHash: Hash.t,
  blockHeight: ID.Block.t,
  success: bool,
  gasFee: list(Coin.t),
  gasLimit: int,
  gasUsed: int,
  sender: Address.t,
  block: block_t,
  messages: Js.Json.t,
  memo: string,
  errMsg: option(string),
};

type account_transaction_t = {transaction: internal_t};

module Mini = {
  type block_t = {timestamp: MomentRe.Moment.t};
  type t = {
    hash: Hash.t,
    blockHeight: ID.Block.t,
    block: block_t,
    gasFee: list(Coin.t),
  };
};

let toExternal =
    (
      {
        id,
        txHash,
        blockHeight,
        success,
        gasFee,
        gasLimit,
        gasUsed,
        sender,
        memo,
        block,
        messages,
        errMsg,
      },
    ) => {
  id,
  txHash,
  blockHeight,
  success,
  gasFee,
  gasLimit,
  gasUsed,
  sender,
  memo,
  timestamp: block.timestamp,
  messages: {
    let msg = messages |> Js.Json.decodeArray |> Belt.Option.getExn |> Belt.List.fromArray;
    msg->Belt.List.map(success ? Msg.decodeAction : Msg.decodeFailAction);
  },
  errMsg: errMsg->Belt.Option.getWithDefault(""),
};

module SingleConfig = [%graphql
  {|
  subscription Transaction($tx_hash: bytea!) {
    transactions_by_pk(hash: $tx_hash) @bsRecord {
      id
      txHash: hash @bsDecoder(fn: "GraphQLParser.hash")
      blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
      success
      memo
      gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
      gasLimit: gas_limit
      gasUsed: gas_used
      sender  @bsDecoder(fn: "Address.fromBech32")
      messages
      errMsg: err_msg
      block @bsRecord {
        timestamp @bsDecoder(fn: "GraphQLParser.timestamp")
      }
    }
  },
|}
];

module MultiConfig = [%graphql
  {|
  subscription Transactions($limit: Int!, $offset: Int!) {
    transactions(offset: $offset, limit: $limit, order_by: {id: desc}) @bsRecord {
      id
      txHash: hash @bsDecoder(fn: "GraphQLParser.hash")
      blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
      success
      memo
      gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
      gasLimit: gas_limit
      gasUsed: gas_used
      sender  @bsDecoder(fn: "Address.fromBech32")
      messages
      errMsg: err_msg
      block @bsRecord {
        timestamp  @bsDecoder(fn: "GraphQLParser.timestamp")
      }
    }
  }
|}
];

module MultiByHeightConfig = [%graphql
  {|
  subscription TransactionsByHeight($height: Int!, $limit: Int!, $offset: Int!) {
    transactions(where: {block_height: {_eq: $height}}, offset: $offset, limit: $limit, order_by: {id: desc}) @bsRecord {
      id
      txHash: hash @bsDecoder(fn: "GraphQLParser.hash")
      blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
      success
      memo
      gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
      gasLimit: gas_limit
      gasUsed: gas_used
      sender  @bsDecoder(fn: "Address.fromBech32")
      messages
      errMsg: err_msg
      block @bsRecord {
        timestamp  @bsDecoder(fn: "GraphQLParser.timestamp")
      }
    }
  }
|}
];

module MultiBySenderConfig = [%graphql
  {|
  subscription TransactionsBySender($sender: String!, $limit: Int!, $offset: Int!) {
    accounts_by_pk(address: $sender) {
      account_transactions(offset: $offset, limit: $limit, order_by: {transaction_id: desc}) @bsRecord{
        transaction @bsRecord {
          id
          txHash: hash @bsDecoder(fn: "GraphQLParser.hash")
          blockHeight: block_height @bsDecoder(fn: "ID.Block.fromInt")
          success
          memo
          gasFee: gas_fee @bsDecoder(fn: "GraphQLParser.coins")
          gasLimit: gas_limit
          gasUsed: gas_used
          sender  @bsDecoder(fn: "Address.fromBech32")
          messages
          errMsg: err_msg
          block @bsRecord {
            timestamp  @bsDecoder(fn: "GraphQLParser.timestamp")
          }
        }
      }
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
    accounts_by_pk(address: $sender) {
      account_transactions_aggregate {
        aggregate {
          count @bsDecoder(fn: "Belt_Option.getExn")
        }
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
  | Some(data) => Sub.resolve(data |> toExternal)
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
  result |> Sub.map(_, x => x##transactions->Belt_Array.map(toExternal));
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
  result
  |> Sub.map(_, x => {
       switch (x##accounts_by_pk) {
       | Some(x') =>
         x'##account_transactions->Belt_Array.map(({transaction}) => transaction->toExternal)
       | None => [||]
       }
     });
};

let getListByBlockHeight = (height, ~page, ~pageSize, ()) => {
  let offset = (page - 1) * pageSize;
  let (result, _) =
    ApolloHooks.useSubscription(
      MultiByHeightConfig.definition,
      ~variables=
        MultiByHeightConfig.makeVariables(
          ~height=height |> ID.Block.toInt,
          ~limit=pageSize,
          ~offset,
          (),
        ),
    );
  result |> Sub.map(_, x => x##transactions->Belt_Array.map(toExternal));
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
  |> Sub.map(_, a => {
       switch (a##accounts_by_pk) {
       | Some(account) =>
         account##account_transactions_aggregate##aggregate
         |> Belt_Option.getExn
         |> (y => y##count)
       | None => 0
       }
     });
};
