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
    type t = {
      id: ID.DataSource.t,
      owner: Address.t,
      name: string,
      executable: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> at(["extra", "id"], ID.DataSource.fromJson),
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
    type t = {
      id: ID.OracleScript.t,
      owner: Address.t,
      name: string,
      code: JsBuffer.t,
      sender: Address.t,
    };

    let decode = json =>
      JsonUtils.Decode.{
        id: json |> at(["extra", "id"], ID.OracleScript.fromJson),
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
    type t = {
      validator: Address.t,
      reporter: Address.t,
      validatorMoniker: string,
    };
    let decode = json =>
      JsonUtils.Decode.{
        validator: json |> at(["msg", "validator"], string) |> Address.fromBech32,
        reporter: json |> at(["msg", "reporter"], string) |> Address.fromBech32,
        validatorMoniker: json |> at(["extra", "validator_moniker"], string),
      };
  };

  module RemoveReporter = {
    type t = {
      validator: Address.t,
      reporter: Address.t,
      validatorMoniker: string,
    };
    let decode = json =>
      JsonUtils.Decode.{
        validator: json |> at(["msg", "validator"], string) |> Address.fromBech32,
        reporter: json |> at(["msg", "reporter"], string) |> Address.fromBech32,
        validatorMoniker: json |> at(["extra", "validator_moniker"], string),
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

  module FailMessage = {
    type t = {
      sender: Address.t,
      message: badge_t,
    };
    let decode = (json, sender) =>
      JsonUtils.Decode.{
        sender,
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
    type t = {
      validatorAddress: Address.t,
      delegatorAddress: Address.t,
      amount: list(Coin.t),
    };
    let decode = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
        delegatorAddress: json |> at(["msg", "delegator_address"], string) |> Address.fromBech32,
        amount: json |> at(["extra", "reward_amount"], string) |> GraphQLParser.coins,
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
    type t = {
      validatorAddress: Address.t,
      amount: list(Coin.t),
    };
    let decode = json => {
      JsonUtils.Decode.{
        validatorAddress: json |> at(["msg", "validator_address"], string) |> Address.fromBech32,
        amount: json |> at(["extra", "commission_amount"], string) |> GraphQLParser.coins,
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
    | SendMsg(Send.t)
    | ReceiveMsg(Receive.t)
    | CreateDataSourceMsg(CreateDataSource.t)
    | EditDataSourceMsg(EditDataSource.t)
    | CreateOracleScriptMsg(CreateOracleScript.t)
    | EditOracleScriptMsg(EditOracleScript.t)
    | RequestMsg(Request.t)
    | ReportMsg(Report.t)
    | AddReporterMsg(AddReporter.t)
    | RemoveReporterMsg(RemoveReporter.t)
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
    | ActivateMsg(Activate.t)
    | UnknownMsg;

  let getCreator = msg => {
    switch (msg) {
    | ReceiveMsg(receive) => receive.fromAddress
    | SendMsg(send) => send.fromAddress
    | CreateDataSourceMsg(dataSource) => dataSource.sender
    | EditDataSourceMsg(dataSource) => dataSource.sender
    | CreateOracleScriptMsg(oracleScript) => oracleScript.sender
    | EditOracleScriptMsg(oracleScript) => oracleScript.sender
    | RequestMsg(request) => request.sender
    | ReportMsg(report) => report.reporter
    | AddReporterMsg(address) => address.validator
    | RemoveReporterMsg(address) => address.validator
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
    | ActivateMsg(activator) => activator.validatorAddress
    | _ => "" |> Address.fromHex
    };
  };

  let getCatVarientbyMsgType =
    fun
    | DelegateMsg(_)
    | UndelegateMsg(_)
    | RedelegateMsg(_)
    | WithdrawRewardMsg(_)
    | WithdrawCommissionMsg(_)
    | SendMsg(_)
    | MultiSendMsg(_)
    | ReceiveMsg(_) => TokenMsg
    | _ => UnknownMsg;

  let getNameByMsgType =
    fun
    | DelegateMsg(_) => "Delegate"
    | UndelegateMsg(_) => "Undelegate"
    | RedelegateMsg(_) => "Redelegate"
    | WithdrawRewardMsg(_) => "Withdraw Reward"
    | WithdrawCommissionMsg(_) => "Withdraw Comission"
    | SendMsg(_) => "Send"
    | MultiSendMsg(_) => "Multisend"
    | ReceiveMsg(_) => "Receive"
    | _ => "";

  type badge_theme_t = {
    text: string,
    textColor: Css.Types.Color.t,
    bgColor: Css.Types.Color.t,
  };

  let getBadge = badgeVariant => {
    switch (badgeVariant) {
    | SendBadge => {text: "SEND", textColor: Colors.blue7, bgColor: Colors.blue1}
    | ReceiveBadge => {text: "RECEIVE", textColor: Colors.green1, bgColor: Colors.green7}
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
    | AddReporterBadge => {
        text: "ADD REPORTER",
        textColor: Colors.purple6,
        bgColor: Colors.purple1,
      }
    | RemoveReporterBadge => {
        text: "REMOVE REPORTER",
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
    | ActivateBadge => {text: "ACTIVATE", textColor: Colors.blue7, bgColor: Colors.blue1}
    | UnknownBadge => {text: "UNKNOWN", textColor: Colors.gray7, bgColor: Colors.gray4}
    };
  };

  let getBadgeTheme = msg => {
    switch (msg) {
    | SendMsg(_) => getBadge(SendBadge)
    | ReceiveMsg(_) => getBadge(ReceiveBadge)
    | CreateDataSourceMsg(_) => getBadge(CreateDataSourceBadge)
    | EditDataSourceMsg(_) => getBadge(EditDataSourceBadge)
    | CreateOracleScriptMsg(_) => getBadge(CreateOracleScriptBadge)
    | EditOracleScriptMsg(_) => getBadge(EditOracleScriptBadge)
    | RequestMsg(_) => getBadge(RequestBadge)
    | ReportMsg(_) => getBadge(ReportBadge)
    | AddReporterMsg(_) => getBadge(AddReporterBadge)
    | RemoveReporterMsg(_) => getBadge(RemoveReporterBadge)
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
    | ActivateMsg(_) => getBadge(ActivateBadge)
    | FailMsg(msg) => getBadge(msg.message)
    | UnknownMsg => getBadge(UnknownBadge)
    };
  };

  let decodeAction = json => {
    JsonUtils.Decode.(
      switch (json |> field("type", string) |> getBadgeVariantFromString) {
      | SendBadge => SendMsg(json |> Send.decode)
      | ReceiveBadge => raise(Not_found)
      | CreateDataSourceBadge => CreateDataSourceMsg(json |> CreateDataSource.decode)
      | EditDataSourceBadge => EditDataSourceMsg(json |> EditDataSource.decode)
      | CreateOracleScriptBadge => CreateOracleScriptMsg(json |> CreateOracleScript.decode)
      | EditOracleScriptBadge => EditOracleScriptMsg(json |> EditOracleScript.decode)
      | RequestBadge => RequestMsg(json |> Request.decode)
      | ReportBadge => ReportMsg(json |> Report.decode)
      | AddReporterBadge => AddReporterMsg(json |> AddReporter.decode)
      | RemoveReporterBadge => RemoveReporterMsg(json |> RemoveReporter.decode)
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
      | ActivateBadge => ActivateMsg(json |> Activate.decode)
      | UnknownBadge => UnknownMsg
      }
    );
  };

  let decodeFailAction = (json, sender): t => FailMsg(json |> FailMessage.decode(_, sender));
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
    success
      ? msg->Belt.List.map(Msg.decodeAction)
      : msg->Belt.List.map(each => Msg.decodeFailAction(each, sender));
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
