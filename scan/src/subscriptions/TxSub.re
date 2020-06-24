module RawDataReport = {
  type t = {
    externalDataID: int,
    data: JsBuffer.t,
  };

  let decode = json =>
    JsonUtils.Decode.{
      externalDataID: json |> field("external_id", int),
      data: json |> field("data", string) |> JsBuffer.fromBase64,
    };
};

module Msg = {
  type badge_t =
    | SendBadge
    | CreateDataSourceBadge
    | EditDataSourceBadge
    | CreateOracleScriptBadge
    | EditOracleScriptBadge
    | RequestBadge
    | ReportBadge
    | AddOracleAddressBadge
    | RemoveOracleAddressBadge
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
    | UnknownBadge;

  let getBadgeVariantFromString = badge => {
    switch (badge) {
    | "send" => SendBadge
    | "create_data_source" => CreateDataSourceBadge
    | "edit_data_source" => EditDataSourceBadge
    | "create_oracle_script" => CreateOracleScriptBadge
    | "edit_oracle_script" => EditOracleScriptBadge
    | "request" => RequestBadge
    | "report" => ReportBadge
    | "add_reporter" => AddOracleAddressBadge
    | "remove_reporter" => RemoveOracleAddressBadge
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
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> field("data_source_id", ID.DataSource.fromJson),
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
        executable: json |> field("executable", string) |> JsBuffer.fromBase64,
        sender: json |> field("sender", string) |> Address.fromBech32,
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
        id: json |> field("data_source_id", ID.DataSource.fromJson),
        owner: json |> field("owner", string) |> Address.fromBech32,
        name: json |> field("name", string),
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
      askCount: int,
      minCount: int,
      schema: string,
      sender: Address.t,
    };

    let decode = json => {
      JsonUtils.Decode.{
        id: json |> field("request_id", ID.Request.fromJson),
        oracleScriptID: json |> field("oracle_script_id", ID.OracleScript.fromJson),
        oracleScriptName: json |> field("oracle_script_name", string),
        calldata: json |> field("calldata", string) |> JsBuffer.fromBase64,
        askCount: json |> field("ask_count", int),
        minCount: json |> field("min_count", int),
        schema: json |> field("schema", string),
        sender: json |> field("sender", string) |> Address.fromBech32,
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
        requestID: json |> field("request_id", ID.Request.fromJson),
        rawReports: json |> field("raw_reports", list(RawDataReport.decode)),
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
      minSelfDelegation: Coin.t,
      selfDelegation: Coin.t,
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
        minSelfDelegation:
          json |> field("min_self_delegation", floatstr) |> Coin.newUBANDFromAmount,
        selfDelegation: json |> field("value", Coin.decodeCoin),
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
    let decode = json =>
      JsonUtils.Decode.{
        moniker: json |> field("moniker", string),
        identity: json |> field("identity", string),
        website: json |> field("website", string),
        details: json |> field("details", string),
        commissionRate: json |> optional(field("commission_rate", floatstr)),
        sender: json |> field("address", string) |> Address.fromBech32,
        minSelfDelegation:
          json
          |> optional(field("min_self_delegation", floatstr))
          |> Belt.Option.map(_, Coin.newUBANDFromAmount),
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

  module FailMessage = {
    type t = {
      sender: Address.t,
      message: badge_t,
    };
    let decode = json =>
      JsonUtils.Decode.{
        sender: json |> field("sender", string) |> Address.fromBech32,
        message: json |> field("type", string) |> getBadgeVariantFromString,
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
        delegatorAddress: json |> field("delegator_address", string) |> Address.fromBech32,
        validatorAddress: json |> field("validator_address", string) |> Address.fromBech32,
        amount: json |> field("amount", Coin.decodeCoin),
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
        delegatorAddress: json |> field("delegator_address", string) |> Address.fromBech32,
        validatorAddress: json |> field("validator_address", string) |> Address.fromBech32,
        amount: json |> field("amount", Coin.decodeCoin),
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
          json |> field("validator_src_address", string) |> Address.fromBech32,
        validatorDestinationAddress:
          json |> field("validator_dst_address", string) |> Address.fromBech32,
        delegatorAddress: json |> field("delegator_address", string) |> Address.fromBech32,
        amount: json |> field("amount", Coin.decodeCoin),
      };
    };
  };

  module WithdrawReward = {
    type t = {
      validatorAddress: Address.t,
      delegatorAddress: Address.t,
      amount: list(Coin.t),
    };
    let decode = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> field("validator_address", string) |> Address.fromBech32,
        delegatorAddress: json |> field("delegator_address", string) |> Address.fromBech32,
        amount: {
          exception WrongNetwork(string);
          switch (Env.network) {
          | "GUANYU"
          | "WENCHANG38" =>
            json
            |> field("reward_amount", array(string))
            |> Belt.Array.getExn(_, 0)
            |> GraphQLParser.coins
          | "WENCHANG" => json |> field("reward_amount", string) |> GraphQLParser.coins
          | _ => raise(WrongNetwork("Incorrect or unspecified NETWORK environment variable"))
          };
        },
      };
    };
  };

  module Unjail = {
    type t = {address: Address.t};

    let decode = json =>
      JsonUtils.Decode.{address: json |> field("address", string) |> Address.fromBech32};
  };
  module SetWithdrawAddress = {
    type t = {
      delegatorAddress: Address.t,
      withdrawAddress: Address.t,
    };
    let decode = json => {
      JsonUtils.Decode.{
        delegatorAddress: json |> field("delegator_address", string) |> Address.fromBech32,
        withdrawAddress: json |> field("withdraw_address", string) |> Address.fromBech32,
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
        proposer: json |> field("proposer", string) |> Address.fromBech32,
        title: json |> at(["content", "title"], string),
        description: json |> at(["content", "description"], string),
        initialDeposit: json |> field("initial_deposit", list(Coin.decodeCoin)),
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
        depositor: json |> field("depositor", string) |> Address.fromBech32,
        proposalID: json |> field("proposal_id", int),
        amount: json |> field("amount", list(Coin.decodeCoin)),
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
        voterAddress: json |> field("voter", string) |> Address.fromBech32,
        proposalID: json |> field("proposal_id", int),
        option: json |> field("option", string),
      };
    };
  };
  module WithdrawCommission = {
    type t = {
      validatorAddress: Address.t,
      amount: list(Coin.t),
    };
    let decode = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> field("validator_address", string) |> Address.fromBech32,
        amount: json |> field("commission_amount", string) |> GraphQLParser.coins,
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
        inputs: json |> field("inputs", list(decodeSendTx)),
        outputs: json |> field("outputs", list(decodeSendTx)),
      };
    };
  };

  type t =
    | SendMsg(Send.t)
    | CreateDataSourceMsg(CreateDataSource.t)
    | EditDataSourceMsg(EditDataSource.t)
    | CreateOracleScriptMsg(CreateOracleScript.t)
    | EditOracleScriptMsg(EditOracleScript.t)
    | RequestMsg(Request.t)
    | ReportMsg(Report.t)
    | AddOracleAddressMsg(AddOracleAddress.t)
    | RemoveOracleAddressMsg(RemoveOracleAddress.t)
    | CreateValidatorMsg(CreateValidator.t)
    | EditValidatorMsg(EditValidator.t)
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
    | FailMsg(FailMessage.t)
    | DelegateMsg(Delegate.t)
    | UndelegateMsg(Undelegate.t)
    | RedelegateMsg(Redelegate.t)
    | WithdrawRewardMsg(WithdrawReward.t)
    | UnjailMsg(Unjail.t)
    | SetWithdrawAddressMsg(SetWithdrawAddress.t)
    | SubmitProposalMsg(SubmitProposal.t)
    | DepositMsg(Deposit.t)
    | VoteMsg(Vote.t)
    | WithdrawCommissionMsg(WithdrawCommission.t)
    | MultiSendMsg(MultiSend.t)
    | UnknownMsg;

  let getCreator = msg => {
    switch (msg) {
    | SendMsg(send) => send.fromAddress
    | CreateDataSourceMsg(dataSource) => dataSource.sender
    | EditDataSourceMsg(dataSource) => dataSource.sender
    | CreateOracleScriptMsg(oracleScript) => oracleScript.sender
    | EditOracleScriptMsg(oracleScript) => oracleScript.sender
    | RequestMsg(request) => request.sender
    | ReportMsg(report) => report.reporter
    | AddOracleAddressMsg(address) => address.validator
    | RemoveOracleAddressMsg(address) => address.validator
    | CreateValidatorMsg(validator) => validator.delegatorAddress
    | EditValidatorMsg(validator) => validator.sender
    | FailMsg(fail) => fail.sender
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
    | DelegateMsg(delegation) => delegation.delegatorAddress
    | UndelegateMsg(delegation) => delegation.delegatorAddress
    | RedelegateMsg(delegation) => delegation.delegatorAddress
    | WithdrawRewardMsg(withdrawal) => withdrawal.delegatorAddress
    | UnjailMsg(validator) => validator.address
    | SetWithdrawAddressMsg(set) => set.delegatorAddress
    | SubmitProposalMsg(proposal) => proposal.proposer
    | DepositMsg(deposit) => deposit.depositor
    | VoteMsg(vote) => vote.voterAddress
    | WithdrawCommissionMsg(withdrawal) => withdrawal.validatorAddress
    | MultiSendMsg(tx) =>
      let firstInput = tx.inputs |> Belt_List.getExn(_, 0);
      firstInput.address;
    | _ => "" |> Address.fromHex
    };
  };

  type badge_theme_t = {
    text: string,
    textColor: Css.Types.Color.t,
    bgColor: Css.Types.Color.t,
  };

  let getBadge = badgeVariant => {
    switch (badgeVariant) {
    | SendBadge => {text: "SEND", textColor: Colors.blue7, bgColor: Colors.blue1}
    | CreateDataSourceBadge => {
        text: "CREATE DATA SOURCE",
        textColor: Colors.yellow5,
        bgColor: Colors.yellow1,
      }
    | EditDataSourceBadge => {
        text: "EDIT DATA SOURCE",
        textColor: Colors.yellow5,
        bgColor: Colors.yellow1,
      }
    | CreateOracleScriptBadge => {
        text: "CREATE ORACLE SCRIPT",
        textColor: Colors.pink6,
        bgColor: Colors.pink1,
      }
    | EditOracleScriptBadge => {
        text: "EDIT ORACLE SCRIPT",
        textColor: Colors.pink6,
        bgColor: Colors.pink1,
      }
    | RequestBadge => {text: "REQUEST", textColor: Colors.orange6, bgColor: Colors.orange1}
    | ReportBadge => {text: "REPORT", textColor: Colors.orange6, bgColor: Colors.orange1}
    | AddOracleAddressBadge => {
        text: "ADD ORACLE ADDRESS",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | RemoveOracleAddressBadge => {
        text: "REMOVE ORACLE ADDRESS",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | CreateValidatorBadge => {
        text: "CREATE VALIDATOR",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | EditValidatorBadge => {
        text: "EDIT VALIDATOR",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | CreateClientBadge => {text: "CREATE CLIENT", textColor: Colors.blue7, bgColor: Colors.blue1}
    | UpdateClientBadge => {text: "UPDATE CLIENT", textColor: Colors.blue7, bgColor: Colors.blue1}
    | SubmitClientMisbehaviourBadge => {
        text: "SUBMIT CLIENT MISBEHAVIOUR",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ConnectionOpenInitBadge => {
        text: "CONNECTION OPEN INIT",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ConnectionOpenTryBadge => {
        text: "CONNECTION OPEN TRY",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ConnectionOpenAckBadge => {
        text: "CONNECTION OPEN ACK",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ConnectionOpenConfirmBadge => {
        text: "CONNECTION OPEN CONFIRM",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelOpenInitBadge => {
        text: "CHANNEL OPEN INIT",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelOpenTryBadge => {
        text: "CHANNEL OPEN TRY",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelOpenAckBadge => {
        text: "CHANNEL OPEN ACK",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelOpenConfirmBadge => {
        text: "CHANNEL OPEN CONFIRM",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelCloseInitBadge => {
        text: "CHANNEL CLOSE INIT",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | ChannelCloseConfirmBadge => {
        text: "CHANNEL CLOSE CONFIRM",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | PacketBadge => {text: "PACKET", textColor: Colors.blue7, bgColor: Colors.blue1}
    | AcknowledgementBadge => {
        text: "ACKNOWLEDGEMENT",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | TimeoutBadge => {text: "TIMEOUT", textColor: Colors.blue7, bgColor: Colors.blue1}
    | DelegateBadge => {text: "DELEGATE", textColor: Colors.purple6, bgColor: Colors.purple1}
    | UndelegateBadge => {text: "UNDELEGATE", textColor: Colors.purple6, bgColor: Colors.purple1}
    | RedelegateBadge => {text: "REDELEGATE", textColor: Colors.purple6, bgColor: Colors.purple1}
    | VoteBadge => {text: "VOTE", textColor: Colors.blue7, bgColor: Colors.blue1}
    | WithdrawRewardBadge => {
        text: "WITHDRAW REWARD",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | UnjailBadge => {text: "UNJAIL", textColor: Colors.blue7, bgColor: Colors.blue1}
    | SetWithdrawAddressBadge => {
        text: "SET WITHDRAW ADDRESS",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | SubmitProposalBadge => {
        text: "SUBMIT PROPOSAL",
        textColor: Colors.blue7,
        bgColor: Colors.blue1,
      }
    | DepositBadge => {text: "DEPOSIT", textColor: Colors.blue7, bgColor: Colors.blue1}
    | WithdrawCommissionBadge => {
        text: "WITHDRAW COMMISSION",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | MultiSendBadge => {text: "MULTI SEND", textColor: Colors.blue7, bgColor: Colors.blue1}
    | UnknownBadge => {text: "UNKNOWN", textColor: Colors.gray7, bgColor: Colors.gray4}
    };
  };

  let getBadgeTheme = msg => {
    switch (msg) {
    | SendMsg(_) => getBadge(SendBadge)
    | CreateDataSourceMsg(_) => getBadge(CreateDataSourceBadge)
    | EditDataSourceMsg(_) => getBadge(EditDataSourceBadge)
    | CreateOracleScriptMsg(_) => getBadge(CreateOracleScriptBadge)
    | EditOracleScriptMsg(_) => getBadge(EditOracleScriptBadge)
    | RequestMsg(_) => getBadge(RequestBadge)
    | ReportMsg(_) => getBadge(ReportBadge)
    | AddOracleAddressMsg(_) => getBadge(AddOracleAddressBadge)
    | RemoveOracleAddressMsg(_) => getBadge(RemoveOracleAddressBadge)
    | CreateValidatorMsg(_) => getBadge(CreateValidatorBadge)
    | EditValidatorMsg(_) => getBadge(EditValidatorBadge)
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
    | DelegateMsg(_) => getBadge(DelegateBadge)
    | UndelegateMsg(_) => getBadge(UndelegateBadge)
    | RedelegateMsg(_) => getBadge(RedelegateBadge)
    | VoteMsg(_) => getBadge(VoteBadge)
    | WithdrawRewardMsg(_) => getBadge(WithdrawRewardBadge)
    | UnjailMsg(_) => getBadge(UnjailBadge)
    | SetWithdrawAddressMsg(_) => getBadge(SetWithdrawAddressBadge)
    | SubmitProposalMsg(_) => getBadge(SubmitProposalBadge)
    | DepositMsg(_) => getBadge(DepositBadge)
    | WithdrawCommissionMsg(_) => getBadge(WithdrawCommissionBadge)
    | MultiSendMsg(_) => getBadge(MultiSendBadge)
    | FailMsg(msg) => getBadge(msg.message)
    | UnknownMsg => getBadge(UnknownBadge)
    };
  };

  let decodeAction = json => {
    JsonUtils.Decode.(
      switch (json |> field("type", string) |> getBadgeVariantFromString) {
      | SendBadge => SendMsg(json |> Send.decode)
      | CreateDataSourceBadge => CreateDataSourceMsg(json |> CreateDataSource.decode)
      | EditDataSourceBadge => EditDataSourceMsg(json |> EditDataSource.decode)
      | CreateOracleScriptBadge => CreateOracleScriptMsg(json |> CreateOracleScript.decode)
      | EditOracleScriptBadge => EditOracleScriptMsg(json |> EditOracleScript.decode)
      | RequestBadge => RequestMsg(json |> Request.decode)
      | ReportBadge => ReportMsg(json |> Report.decode)
      | AddOracleAddressBadge => AddOracleAddressMsg(json |> AddOracleAddress.decode)
      | RemoveOracleAddressBadge => RemoveOracleAddressMsg(json |> RemoveOracleAddress.decode)
      | CreateValidatorBadge => CreateValidatorMsg(json |> CreateValidator.decode)
      | EditValidatorBadge => EditValidatorMsg(json |> EditValidator.decode)
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
      | DelegateBadge => DelegateMsg(json |> Delegate.decode)
      | UndelegateBadge => UndelegateMsg(json |> Undelegate.decode)
      | RedelegateBadge => RedelegateMsg(json |> Redelegate.decode)
      | WithdrawRewardBadge => WithdrawRewardMsg(json |> WithdrawReward.decode)
      | UnjailBadge => UnjailMsg(json |> Unjail.decode)
      | SetWithdrawAddressBadge => SetWithdrawAddressMsg(json |> SetWithdrawAddress.decode)
      | SubmitProposalBadge => SubmitProposalMsg(json |> SubmitProposal.decode)
      | DepositBadge => DepositMsg(json |> Deposit.decode)
      | VoteBadge => VoteMsg(json |> Vote.decode)
      | WithdrawCommissionBadge => WithdrawCommissionMsg(json |> WithdrawCommission.decode)
      | MultiSendBadge => MultiSendMsg(json |> MultiSend.decode)
      | UnknownBadge => UnknownMsg
      }
    );
  };

  let decodeFailAction = (json): t => FailMsg(json |> FailMessage.decode);

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
  rawLog: string,
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
      timestamp  @bsDecoder(fn: "GraphQLParser.timeMS")
      messages @bsDecoder(fn: "Msg.decodeActions")
      rawLog: raw_log
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
      timestamp  @bsDecoder(fn: "GraphQLParser.timeMS")
      messages @bsDecoder(fn: "Msg.decodeActions")
      rawLog: raw_log
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
      timestamp  @bsDecoder(fn: "GraphQLParser.timeMS")
      messages @bsDecoder(fn: "Msg.decodeActions")
      rawLog: raw_log
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
      timestamp  @bsDecoder(fn: "GraphQLParser.timeMS")
      messages @bsDecoder(fn: "Msg.decodeActions")
      rawLog: raw_log
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
