type coin_amount_t = {
  value: list(Coin.t),
  hasDenom: bool,
};

type request_count_t = {
  reportedValidators: int,
  minimumValidators: int,
  requestValidators: int,
};

type request_response_t = {
  request: int,
  responseTime: int,
};

type t =
  | Address(Address.t, int, [ | `account | `validator])
  | Height(ID.Block.t)
  | Coin(coin_amount_t)
  | Count(int)
  | DataSource(ID.DataSource.t, string)
  | OracleScript(ID.OracleScript.t, string)
  | RequestID(ID.Request.t)
  | RequestResponse(request_response_t)
  | ProgressBar(request_count_t)
  | Float(float, option(int))
  | KVTableReport(list(string), list(TxSub.RawDataReport.t))
  | KVTableRequest(option(array(Obi.field_key_value_t)))
  | CopyButton(JsBuffer.t)
  | Percentage(float, option(int))
  | Timestamp(MomentRe.Moment.t)
  | TxHash(Hash.t, int)
  | Validator(Address.t, string, string)
  | Messages(Hash.t, list(TxSub.Msg.t), bool, string)
  | PubKey(PubKey.t)
  | Badge(TxSub.Msg.badge_theme_t)
  | VotingPower(Coin.t, float)
  | Uptime(option(float))
  | Loading(int)
  | Text(string)
  | Nothing;

module Styles = {
  open Css;
  let vFlex = style([display(`flex), alignItems(`center)]);
  let addressContainer = w => {
    style([width(`px(w))]);
  };
  let badge = color =>
    style([
      display(`inlineFlex),
      padding2(~v=`px(5), ~h=`px(10)),
      backgroundColor(color),
      borderRadius(`px(15)),
    ]);
};

[@react.component]
let make = (~info) => {
  switch (info) {
  | Address(address, width, accountType) =>
    <div className={Styles.addressContainer(width)}>
      <AddressRender address position=AddressRender.Text clickable=true accountType />
    </div>
  | Height(height) =>
    <div className=Styles.vFlex> <TypeID.Block id=height position=TypeID.Subtitle /> </div>
  | Coin({value, hasDenom}) =>
    <AmountRender coins=value pos={hasDenom ? AmountRender.TxIndex : Fee} />
  | Count(value) =>
    <Text
      value={value |> Format.iPretty}
      size=Text.Md
      weight=Text.Semibold
      spacing={Text.Em(0.02)}
      code=true
    />
  | DataSource(id, name) =>
    <div className=Styles.vFlex>
      <TypeID.DataSource id />
      <HSpacing size=Spacing.sm />
      <Text value=name ellipsis=true />
    </div>
  | OracleScript(id, name) =>
    <div className=Styles.vFlex>
      <TypeID.OracleScript id />
      <HSpacing size=Spacing.sm />
      <Text value=name ellipsis=true />
    </div>
  | RequestID(id) => <TypeID.Request id />
  | RequestResponse({request, responseTime}) =>
    <div className={CssHelper.flexBox()}>
      <Text value={request |> Format.iPretty} weight=Text.Regular block=true ellipsis=true />
      <HSpacing size=Spacing.sm />
      <Text
        value={"(" ++ (responseTime |> Format.iPretty) ++ "ms)"}
        weight=Text.Regular
        block=true
        color=Colors.gray6
      />
    </div>
  | ProgressBar({reportedValidators, minimumValidators, requestValidators}) =>
    <ProgressBar reportedValidators minimumValidators requestValidators />
  | Float(value, digits) =>
    <Text
      value={value |> Format.fPretty(~digits?)}
      size=Text.Md
      spacing={Text.Em(0.02)}
      code=true
    />
  | KVTableReport(heading, rawReports) =>
    <KVTable
      tableWidth=480
      headers=heading
      rows={
        rawReports
        |> Belt_List.map(_, rawReport =>
             [
               KVTable.Value(rawReport.externalDataID |> string_of_int),
               KVTable.Value(rawReport.exitCode |> string_of_int),
               KVTable.Value(rawReport.data |> JsBuffer.toUTF8),
             ]
           )
      }
    />
  | KVTableRequest(calldataKVsOpt) =>
    switch (calldataKVsOpt) {
    | Some(calldataKVs) =>
      <KVTable
        tableWidth=480
        rows={
          calldataKVs
          ->Belt_Array.map(({fieldName, fieldValue}) =>
              [KVTable.Value(fieldName), KVTable.Value(fieldValue)]
            )
          ->Belt_List.fromArray
        }
      />
    | None =>
      <Text
        value="Could not decode calldata."
        spacing={Text.Em(0.02)}
        nowrap=true
        ellipsis=true
        code=true
        block=true
      />
    }
  | CopyButton(calldata) => <CopyButton data=calldata title="Copy as bytes" />
  | Percentage(value, digits) =>
    <Text
      value={value |> Format.fPercent(~digits?)}
      size=Text.Md
      spacing={Text.Em(0.02)}
      code=true
    />
  | Text(text) =>
    <Text value=text spacing={Text.Em(0.02)} nowrap=true ellipsis=true code=true block=true />
  | Timestamp(time) => <Timestamp time size=Text.Md weight=Text.Regular code=true />
  | Validator(address, moniker, identity) =>
    <ValidatorMonikerLink
      validatorAddress=address
      moniker
      size=Text.Md
      identity
      width={`px(230)}
    />
  | PubKey(publicKey) => <PubKeyRender alignLeft=true pubKey=publicKey display=`block />
  | TxHash(txHash, width) => <TxLink txHash width size=Text.Lg />
  | Messages(txHash, messages, success, errMsg) =>
    <TxMessages txHash messages success errMsg width=360 />
  | Badge({text, textColor, bgColor}) =>
    <div className={Styles.badge(bgColor)}>
      <Text value=text size=Text.Xs spacing={Text.Em(0.07)} color=textColor />
    </div>
  | VotingPower(tokens, votingPercent) =>
    <div className=Styles.vFlex>
      <Text
        value={tokens |> Coin.getBandAmountFromCoin |> Format.fPretty(~digits=0)}
        color=Colors.gray7
        code=true
        weight=Text.Regular
        spacing={Text.Em(0.02)}
        block=true
      />
      <HSpacing size=Spacing.sm />
      <Text
        value={"(" ++ (votingPercent |> Format.fPercent(~digits=2)) ++ ")"}
        color=Colors.gray6
        code=true
        weight=Text.Thin
        spacing={Text.Em(0.02)}
        block=true
      />
    </div>
  // Special case for uptime to have loading state inside.
  | Uptime(uptimeOpt) =>
    switch (uptimeOpt) {
    | Some(uptime) =>
      <div className=Styles.vFlex>
        <Text
          value={uptime |> Format.fPercent(~digits=2)}
          spacing={Text.Em(0.02)}
          code=true
          nowrap=true
        />
        <HSpacing size=Spacing.lg />
        <UptimeBar percent=uptime />
      </div>
    | None => <Text value="N/A" spacing={Text.Em(0.02)} code=true nowrap=true />
    }
  | Loading(width) => <LoadingCensorBar width height=21 />
  | Nothing => React.null
  };
};
