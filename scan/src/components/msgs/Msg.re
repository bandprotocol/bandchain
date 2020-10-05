module Styles = {
  open Css;
  let rowWithWidth = (w: int) =>
    style([
      width(`px(w)),
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      Media.mobile([
        width(`auto),
        flexWrap(`wrap),
        selector("> div:nth-child(1)", [width(`px(90)), marginBottom(`px(10))]),
        selector(
          "> .labelContainer",
          [
            display(`flex),
            flexBasis(`calc((`sub, `percent(100.), `px(100)))),
            marginBottom(`px(10)),
          ],
        ),
      ]),
      Media.smallMobile([
        selector("> div:nth-child(1)", [width(`px(68)), marginBottom(`px(10))]),
      ]),
    ]);
  let withWidth = (w: int) => style([width(`px(w))]);
  let withBg = (color: Types.Color.t, mw: int) =>
    style([
      minWidth(`px(mw)),
      height(`px(16)),
      backgroundColor(color),
      borderRadius(`px(100)),
      margin2(~v=`zero, ~h=`px(5)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
    ]);

  let addressWrapper = style([width(`px(120))]);

  let msgContainer =
    style([
      Media.mobile([
        selector("> div", [width(`percent(100.))]),
        selector("> div + div", [marginTop(`px(8))]),
      ]),
    ]);
};

let makeBadge = (name, length, color1, color2) =>
  <div className="labelContainer">
    <div className={Styles.withBg(color1, length)}>
      <Text value=name size=Text.Xs spacing={Text.Em(0.07)} weight=Text.Medium color=color2 />
    </div>
  </div>;

[@react.component]
let make = (~msg: TxSub.Msg.t, ~width: int) => {
  <div
    className={Css.merge([
      CssHelper.flexBox(~wrap=`nowrap, ()),
      CssHelper.flexBoxSm(~wrap=`wrap, ()),
      CssHelper.overflowHidden,
      Styles.msgContainer,
    ])}>
    <MsgFront
      msgType=msg
      name={msg |> TxSub.Msg.getNameByMsgType}
      fromAddress={msg |> TxSub.Msg.getCreator}
    />
    {switch (msg) {
     | SendMsg({toAddress, amount}) => <TokenMsg.SendMsg toAddress amount />
     | ReceiveMsg({fromAddress, amount}) => <TokenMsg.ReceiveMsg fromAddress amount />
     | CreateDataSourceMsg({id, sender, name}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.yellow1, 110)}>
             <Text
               value="CREATE DATASOURCE"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.yellow6
             />
           </div>
         </div>
         <TypeID.DataSource id />
         <HSpacing size=Spacing.sm />
         <Text
           value=name
           color=Colors.gray7
           weight=Text.Medium
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | EditDataSourceMsg({id, sender, name}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.yellow1, 100)}>
             <Text
               value="EDIT DATASOURCE"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.yellow6
             />
           </div>
         </div>
         <TypeID.DataSource id />
         {name == Config.doNotModify
            ? React.null
            : <>
                <HSpacing size=Spacing.sm />
                <Text
                  value=name
                  color=Colors.gray7
                  weight=Text.Medium
                  nowrap=true
                  block=true
                  ellipsis=true
                />
              </>}
       </div>
     | CreateOracleScriptMsg({id, sender, name}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.pink1, 120)}>
             <Text
               value="CREATE ORACLE SCRIPT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.pink6
             />
           </div>
         </div>
         <div className={Styles.rowWithWidth(200)}>
           <TypeID.OracleScript id />
           <HSpacing size=Spacing.sm />
           <Text
             value=name
             color=Colors.gray7
             weight=Text.Medium
             nowrap=true
             block=true
             ellipsis=true
           />
         </div>
       </div>
     | EditOracleScriptMsg({id, sender, name}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.pink1, 110)}>
             <Text
               value="EDIT ORACLE SCRIPT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.pink6
             />
           </div>
         </div>
         <div className={Styles.rowWithWidth(210)}>
           <TypeID.OracleScript id />
           {name == Config.doNotModify
              ? React.null
              : <>
                  <HSpacing size=Spacing.sm />
                  <Text
                    value=name
                    color=Colors.gray7
                    weight=Text.Medium
                    nowrap=true
                    block=true
                    ellipsis=true
                  />
                </>}
         </div>
       </div>
     | RequestMsg({id, oracleScriptID, oracleScriptName, sender}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.orange1, 60)}>
             <Text
               value="REQUEST"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.orange6
             />
           </div>
         </div>
         <TypeID.Request id />
         <HSpacing size=Spacing.sm />
         <Icon name="far fa-arrow-right" color=Colors.black />
         <HSpacing size=Spacing.sm />
         <TypeID.OracleScript id=oracleScriptID />
         <HSpacing size=Spacing.sm />
         <Text
           value=oracleScriptName
           color=Colors.gray7
           weight=Text.Medium
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ReportMsg({requestID, reporter}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=reporter /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.orange1, 50)}>
             <Text
               value="REPORT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.orange6
             />
           </div>
         </div>
         <Icon name="far fa-arrow-right" color=Colors.black />
         <HSpacing size=Spacing.sm />
         <TypeID.Request id=requestID />
       </div>
     | AddReporterMsg({validator, reporter}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=validator /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.purple1, 80)}>
             <Text
               value="ADD REPORTER"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.purple6
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <div className={Styles.withWidth(120)}> <AddressRender address=reporter /> </div>
       </div>
     | RemoveReporterMsg({validator, reporter}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=validator /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.purple1, 100)}>
             <Text
               value="REMOVE REPORTER"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.purple6
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <div className={Styles.withWidth(120)}> <AddressRender address=reporter /> </div>
       </div>
     | CreateValidatorMsg({delegatorAddress, moniker}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=delegatorAddress /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.purple1, 97)}>
             <Text
               value="CREATE VALIDATOR"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.purple6
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <div className={Styles.withWidth(width / 2 - 5)}>
           <Text
             value=moniker
             color=Colors.gray7
             weight=Text.Regular
             code=true
             nowrap=true
             block=true
             ellipsis=true
           />
         </div>
       </div>
     | EditValidatorMsg({sender, moniker}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.purple1, 85)}>
             <Text
               value="EDIT VALIDATOR"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.purple6
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <div className={Styles.withWidth(width / 2 - 5)}>
           {moniker == Config.doNotModify
              ? <AddressRender address=sender accountType=`validator />
              : <Text
                  value=moniker
                  color=Colors.gray7
                  weight=Text.Regular
                  code=true
                  nowrap=true
                  block=true
                  ellipsis=true
                />}
         </div>
       </div>
     | CreateClientMsg({address, clientID, chainID}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="CREATE CLIENT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=clientID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | UpdateClientMsg({address, clientID, chainID}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="UPDATE CLIENT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=clientID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | SubmitClientMisbehaviourMsg({address, clientID, chainID}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="SUBMIT CLIENT MISBEHAVIOUR"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=clientID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ConnectionOpenInitMsg({signer, common: {connectionID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 140)}>
             <Text
               value="CONNECTION OPEN INIT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=connectionID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ConnectionOpenTryMsg({signer, common: {connectionID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 120)}>
             <Text
               value="CONNECTION OPEN TRY"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=connectionID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ConnectionOpenAckMsg({signer, common: {connectionID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 130)}>
             <Text
               value="CONNECTION OPEN ACK"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=connectionID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ConnectionOpenConfirmMsg({signer, common: {connectionID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 140)}>
             <Text
               value="CONNECTION OPEN CONFIRM"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=connectionID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelOpenInitMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="CHANNEL OPEN INIT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelOpenTryMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="CHANNEL OPEN TRY"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelOpenAckMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="CHANNEL OPEN ACK"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelOpenConfirmMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 130)}>
             <Text
               value="CHANNEL OPEN CONFIRM"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelCloseInitMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="CHANNEL CLOSE INIT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | ChannelCloseConfirmMsg({signer, common: {channelID, chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=signer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="CHANNEL CLOSE CONFIRM"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value=channelID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | PacketMsg({sender, data, common: {chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 50)}>
             <Text
               value="PACKET"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value="data"
           color=Colors.gray7
           weight=Text.Semibold
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <HSpacing size=Spacing.sm />
         <div className={Styles.withWidth(110)}>
           <Text value=data color=Colors.gray7 code=true nowrap=true block=true ellipsis=true />
         </div>
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | TimeoutMsg({sender, common: {chainID}}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 85)}>
             <Text
               value="TIMEOUT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <HSpacing size=Spacing.sm />
         <Text
           value="at"
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
         <Text value={j|||j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
         <Text
           value=chainID
           color=Colors.gray7
           weight=Text.Regular
           code=true
           nowrap=true
           block=true
           ellipsis=true
         />
       </div>
     | DelegateMsg({amount, delegatorAddress}) => <> <AmountRender coins=[amount] /> </>
     | UndelegateMsg({amount, delegatorAddress}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=delegatorAddress /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.purple1, 73)}>
             <Text
               value="UNDELEGATE"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.purple6
             />
           </div>
         </div>
         <AmountRender coins=[amount] />
       </div>
     | RedelegateMsg({amount, delegatorAddress}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=delegatorAddress /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.purple1, 73)}>
             <Text
               value="REDELEGATE"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.purple6
             />
           </div>
         </div>
         <AmountRender coins=[amount] />
       </div>
     | WithdrawRewardMsg({delegatorAddress, amount}) => <AmountRender coins=amount />
     | UnjailMsg({address}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 50)}>
             <Text
               value="UNJAIL"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <div className={Styles.withWidth(width / 2)}>
           <AddressRender address accountType=`validator />
         </div>
       </div>
     | SetWithdrawAddressMsg({delegatorAddress, withdrawAddress}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=delegatorAddress /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.purple1, 130)}>
             <Text
               value="SET WITHDRAW ADDRESS"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.purple6
             />
           </div>
         </div>
         <div className={Styles.withWidth(width / 3)}>
           <AddressRender address=withdrawAddress />
         </div>
       </div>
     | SubmitProposalMsg({proposer, title}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=proposer /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 100)}>
             <Text
               value="SUBMIT PROPOSAL"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <div className={Styles.rowWithWidth(200)}>
           <Text value=title weight=Text.Regular code=true nowrap=true block=true />
         </div>
       </div>
     | DepositMsg({depositor, amount, proposalID}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=depositor /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 50)}>
             <Text
               value="DEPOSIT"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <AmountRender coins=amount />
         <HSpacing size=Spacing.sm />
         <Icon name="far fa-arrow-right" color=Colors.black />
         <HSpacing size=Spacing.sm />
         <div className={Styles.rowWithWidth(200)}>
           <Text
             value={"Proposal " ++ (proposalID |> string_of_int)}
             weight=Text.Regular
             code=true
             nowrap=true
             block=true
           />
         </div>
       </div>
     | VoteMsg({voterAddress, proposalID, option}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=voterAddress /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 40)}>
             <Text
               value="VOTE"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <Text value=option weight=Text.Regular code=true nowrap=true block=true />
         <HSpacing size=Spacing.sm />
         <Icon name="far fa-arrow-right" color=Colors.black />
         <HSpacing size=Spacing.sm />
         <div className={Styles.rowWithWidth(200)}>
           <Text
             value={"Proposal " ++ (proposalID |> string_of_int)}
             weight=Text.Regular
             code=true
             nowrap=true
             block=true
           />
         </div>
       </div>
     | WithdrawCommissionMsg({validatorAddress, amount}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=validatorAddress /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.purple1, 130)}>
             <Text
               value="WITHDRAW COMMISSION"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.purple6
             />
           </div>
         </div>
         <AmountRender coins=amount />
       </div>
     | MultiSendMsg({inputs, outputs}) =>
       let firstInput = inputs |> Belt_List.getExn(_, 0);
       let firstSender = firstInput.address;
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=firstSender /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 70)}>
             <Text
               value="MULTI SEND"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
         <div className={Styles.rowWithWidth(300)}>
           <Text value={inputs |> Belt_List.length |> string_of_int} weight=Text.Semibold />
           <HSpacing size=Spacing.sm />
           <Text value="Inputs" />
           <HSpacing size=Spacing.sm />
           <Icon name="far fa-arrow-right" color=Colors.black />
           <HSpacing size=Spacing.sm />
           <Text value={outputs |> Belt_List.length |> string_of_int} weight=Text.Semibold />
           <HSpacing size=Spacing.sm />
           <Text value="Outputs" />
         </div>
       </div>;
     | ActivateMsg({validatorAddress}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=validatorAddress /> </div>
         <div className="labelContainer">
           <div className={Styles.withBg(Colors.blue1, 65)}>
             <Text
               value="ACTIVATE"
               size=Text.Xs
               spacing={Text.Em(0.07)}
               weight=Text.Medium
               color=Colors.blue7
             />
           </div>
         </div>
       </div>
     | UnknownMsg => makeBadge("UNKNOWN", 70, Colors.gray1, Colors.gray6)
     | FailMsg({sender, message}) =>
       <div className={Styles.rowWithWidth(width)}>
         <div className={Styles.withWidth(120)}> <AddressRender address=sender /> </div>
         {switch (message) {
          | SendBadge => makeBadge("SEND", 40, Colors.blue1, Colors.blue7)
          | ReceiveBadge => makeBadge("RECEIVE", 50, Colors.green1, Colors.green7)
          | CreateDataSourceBadge =>
            makeBadge("CREATE DATASOURCE", 110, Colors.yellow1, Colors.yellow6)
          | EditDataSourceBadge =>
            makeBadge("EDIT DATASOURCE", 100, Colors.yellow1, Colors.yellow6)
          | CreateOracleScriptBadge =>
            makeBadge("CREATE ORACLE SCRIPT", 120, Colors.pink1, Colors.pink6)
          | EditOracleScriptBadge =>
            makeBadge("EDIT ORACLE SCRIPT", 110, Colors.pink1, Colors.pink6)
          | RequestBadge => makeBadge("REQUEST", 60, Colors.orange1, Colors.orange6)
          | ReportBadge => makeBadge("REPORT", 50, Colors.orange1, Colors.orange6)
          | AddReporterBadge => makeBadge("ADD REPORTER", 114, Colors.purple1, Colors.purple6)
          | RemoveReporterBadge =>
            makeBadge("REMOVE REPORTER", 133, Colors.purple1, Colors.purple6)
          | CreateValidatorBadge =>
            makeBadge("CREATE VALIDATOR", 97, Colors.purple1, Colors.purple6)
          | EditValidatorBadge => makeBadge("EDIT VALIDATOR", 85, Colors.purple1, Colors.purple6)
          | CreateClientBadge => makeBadge("CREATE CLIENT", 85, Colors.blue1, Colors.blue7)
          | UpdateClientBadge => makeBadge("UPDATE CLIENT", 85, Colors.blue1, Colors.blue7)
          | SubmitClientMisbehaviourBadge =>
            makeBadge("SUBMIT CLIENT BEHAVIOUR", 140, Colors.blue1, Colors.blue7)
          | ConnectionOpenInitBadge =>
            makeBadge("CONNECTION OPEN INIT", 140, Colors.blue1, Colors.blue7)
          | ConnectionOpenTryBadge =>
            makeBadge("CONNECTION OPEN TRY", 130, Colors.blue1, Colors.blue7)
          | ConnectionOpenAckBadge =>
            makeBadge("CONNECTION OPEN ACK", 130, Colors.blue1, Colors.blue7)
          | ConnectionOpenConfirmBadge =>
            makeBadge("CONNECTION OPEN CONFIRM ", 140, Colors.blue1, Colors.blue7)
          | ChannelOpenInitBadge =>
            makeBadge("CHANNEL OPEN INIT", 100, Colors.blue1, Colors.blue7)
          | ChannelOpenTryBadge => makeBadge("CHANNEL OPEN TRY", 100, Colors.blue1, Colors.blue7)
          | ChannelOpenAckBadge => makeBadge("CHANNEL OPEN ACK", 100, Colors.blue1, Colors.blue7)
          | ChannelOpenConfirmBadge =>
            makeBadge("CHANNEL OPEN CONFIRM", 130, Colors.blue1, Colors.blue7)
          | ChannelCloseInitBadge =>
            makeBadge("CHANNEL CLOSE INIT", 85, Colors.blue1, Colors.blue7)
          | ChannelCloseConfirmBadge =>
            makeBadge("CHANNEL CLOSE CONFIRM", 100, Colors.blue1, Colors.blue7)
          | PacketBadge => makeBadge("PACKET", 50, Colors.blue1, Colors.blue7)
          | AcknowledgementBadge => makeBadge("ACKNOWLEDGEMENT", 100, Colors.blue1, Colors.blue7)
          | TimeoutBadge => makeBadge("TIMEOUT", 50, Colors.blue1, Colors.blue7)
          | DelegateBadge => makeBadge("DELEGATE", 60, Colors.purple1, Colors.purple6)
          | UndelegateBadge => makeBadge("UNDELEGATE", 73, Colors.purple1, Colors.purple6)
          | RedelegateBadge => makeBadge("REDELEGATE", 73, Colors.purple1, Colors.purple6)
          | WithdrawRewardBadge =>
            makeBadge("WITHDRAW REWARD", 100, Colors.purple1, Colors.purple6)
          | SetWithdrawAddressBadge =>
            makeBadge("SET WITHDRAW ADDRESS", 170, Colors.purple1, Colors.purple6)
          | SubmitProposalBadge => makeBadge("SUBMIT PROPOSAL", 100, Colors.blue1, Colors.blue7)
          | DepositBadge => makeBadge("DEPOSIT", 50, Colors.blue1, Colors.blue7)
          | UnjailBadge => makeBadge("UNJAIL", 50, Colors.blue1, Colors.blue7)
          | VoteBadge => makeBadge("VOTE", 40, Colors.blue1, Colors.blue7)
          | WithdrawCommissionBadge =>
            makeBadge("WITHDRAW COMMISSION", 100, Colors.purple1, Colors.purple6)
          | MultiSendBadge => makeBadge("MULTI SEND", 70, Colors.blue1, Colors.blue7)
          | ActivateBadge => makeBadge("Activate", 70, Colors.blue1, Colors.blue7)
          | UnknownBadge => makeBadge("UNKNOWN", 70, Colors.gray1, Colors.gray6)
          }}
       </div>
     | _ => React.null
     }}
  </div>;
};
