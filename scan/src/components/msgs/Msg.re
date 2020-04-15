module Styles = {
  open Css;
  let rowWithWidth = (w: int) =>
    style([width(`px(w)), display(`flex), flexDirection(`row), alignItems(`center)]);
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
};

let makeBadge = (name, length, color1, color2) =>
  <div className={Styles.withBg(color1, length)}>
    <Text value=name size=Text.Xs spacing={Text.Em(0.07)} weight=Text.Medium color=color2 />
  </div>;

[@react.component]
let make = (~msg: TxSub.Msg.t, ~width: int) => {
  switch (msg) {
  | Send({fromAddress, toAddress, amount}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=fromAddress /> </div>
      <div className={Styles.withBg(Colors.blue1, 40)}>
        <Text
          value="SEND"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
      </div>
      <div className={Styles.rowWithWidth(200)}>
        <Text
          value={
            amount->Belt_List.get(0)->Belt_Option.getWithDefault(Coin.newCoin("uband", 0.0)).
              amount
            |> Format.fPretty
          }
          weight=Text.Semibold
          code=true
          nowrap=true
          block=true
        />
        <HSpacing size=Spacing.sm />
        <Text value="BAND" weight=Text.Regular code=true nowrap=true block=true />
        <HSpacing size=Spacing.sm />
        <Text value={j|➜|j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
        <HSpacing size=Spacing.sm />
      </div>
      <div className={Styles.withWidth(width / 3)}> <AddressRender address=toAddress /> </div>
    </div>
  | CreateDataSource({id, sender, name}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.yellow1, 110)}>
        <Text
          value="CREATE DATASOURCE"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.yellow6
        />
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
  | EditDataSource({id, sender, name}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.yellow1, 100)}>
        <Text
          value="EDIT DATASOURCE"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.yellow6
        />
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
  | CreateOracleScript({id, sender, name}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.pink1, 120)}>
        <Text
          value="CREATE ORACLE SCRIPT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.pink6
        />
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
  | EditOracleScript({id, sender, name}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.pink1, 110)}>
        <Text
          value="EDIT ORACLE SCRIPT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.pink6
        />
      </div>
      <div className={Styles.rowWithWidth(210)}>
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
  | Request({id, oracleScriptID, oracleScriptName, sender}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.orange1, 60)}>
        <Text
          value="REQUEST"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.orange6
        />
      </div>
      <TypeID.Request id />
      <HSpacing size=Spacing.sm />
      <Text value={j|➜|j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
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
  | Report({requestID, reporter}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(140)}> <AddressRender address=reporter /> </div>
      <div className={Styles.withBg(Colors.orange1, 50)}>
        <Text
          value="REPORT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.orange6
        />
      </div>
      <Text value={j|➜|j} size=Text.Xxl weight=Text.Bold code=true nowrap=true block=true />
      <HSpacing size=Spacing.sm />
      <TypeID.Request id=requestID />
    </div>
  | AddOracleAddress({validator, reporter}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=validator /> </div>
      <div className={Styles.withBg(Colors.purple1, 114)}>
        <Text
          value="ADD ORACLE ADDRESS"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.purple6
        />
      </div>
      <HSpacing size=Spacing.sm />
      <div className={Styles.withWidth(120)}> <AddressRender address=reporter /> </div>
    </div>
  | RemoveOracleAddress({validator, reporter}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=validator /> </div>
      <div className={Styles.withBg(Colors.purple1, 133)}>
        <Text
          value="REMOVE ORACLE ADDRESS"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.purple6
        />
      </div>
      <HSpacing size=Spacing.sm />
      <div className={Styles.withWidth(120)}> <AddressRender address=reporter /> </div>
    </div>
  | CreateValidator({delegatorAddress, moniker}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=delegatorAddress /> </div>
      <div className={Styles.withBg(Colors.purple1, 97)}>
        <Text
          value="CREATE VALIDATOR"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.purple6
        />
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
  | EditValidator({sender, moniker}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.purple1, 85)}>
        <Text
          value="EDIT VALIDATOR"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.purple6
        />
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
  | CreateClient({address, clientID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address /> </div>
      <div className={Styles.withBg(Colors.blue1, 85)}>
        <Text
          value="CREATE CLIENT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | UpdateClient({address, clientID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address /> </div>
      <div className={Styles.withBg(Colors.blue1, 85)}>
        <Text
          value="UPDATE CLIENT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | SubmitClientMisbehaviour({address, clientID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address /> </div>
      <div className={Styles.withBg(Colors.blue1, 85)}>
        <Text
          value="SUBMIT CLIENT MISBEHAVIOUR"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ConnectionOpenInit({signer, connectionID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 140)}>
        <Text
          value="CONNECTION OPEN INIT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ConnectionOpenTry({signer, connectionID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 120)}>
        <Text
          value="CONNECTION OPEN TRY"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ConnectionOpenAck({signer, connectionID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 130)}>
        <Text
          value="CONNECTION OPEN ACK"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ConnectionOpenConfirm({signer, connectionID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 140)}>
        <Text
          value="CONNECTION OPEN CONFIRM"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ChannelOpenInit({signer, channelID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 100)}>
        <Text
          value="CHANNEL OPEN INIT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ChannelOpenTry({signer, channelID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 100)}>
        <Text
          value="CHANNEL OPEN TRY"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ChannelOpenAck({signer, channelID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 100)}>
        <Text
          value="CHANNEL OPEN ACK"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ChannelOpenConfirm({signer, channelID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 130)}>
        <Text
          value="CHANNEL OPEN CONFIRM"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ChannelCloseInit({signer, channelID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 85)}>
        <Text
          value="CHANNEL CLOSE INIT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | ChannelCloseConfirm({signer, channelID, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=signer /> </div>
      <div className={Styles.withBg(Colors.blue1, 100)}>
        <Text
          value="CHANNEL CLOSE CONFIRM"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | Packet({sender, chainID, data}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.blue1, 50)}>
        <Text
          value="PACKET"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | Timeout({sender, chainID}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.blue1, 85)}>
        <Text
          value="TIMEOUT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
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
  | FailMessage({sender, message}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(130)}> <AddressRender address=sender /> </div>
      {switch (message) {
       | "send" => makeBadge("SEND", 40, Colors.blue1, Colors.blue7)
       | "create_data_source" =>
         makeBadge("CREATE DATASOURCE", 110, Colors.yellow1, Colors.yellow6)
       | "edit_data_source" => makeBadge("EDIT DATASOURCE", 100, Colors.yellow1, Colors.yellow6)
       | "create_oracle_script" =>
         makeBadge("CREATE ORACLE SCRIPT", 120, Colors.pink1, Colors.pink6)
       | "edit_oracle_script" => makeBadge("EDIT ORACLE SCRIPT", 110, Colors.pink1, Colors.pink6)
       | "request" => makeBadge("REQUEST", 60, Colors.orange1, Colors.orange6)
       | "report" => makeBadge("REPORT", 50, Colors.orange1, Colors.orange6)
       | "add_oracle_address" =>
         makeBadge("ADD ORACLE ADDRESS", 114, Colors.purple1, Colors.purple6)
       | "remove_oracle_address" =>
         makeBadge("REMOVE ORACLE ADDRESS", 133, Colors.purple1, Colors.purple6)
       | "create_validator" => makeBadge("CREATE VALIDATOR", 97, Colors.purple1, Colors.purple6)
       | "edit_validator" => makeBadge("EDIT VALIDATOR", 85, Colors.purple1, Colors.purple6)
       | "create_client" => makeBadge("CREATE CLIENT", 85, Colors.blue1, Colors.blue7)
       | "update_client" => makeBadge("UPDATE CLIENT", 85, Colors.blue1, Colors.blue7)
       | "submit_client_behaviour" =>
         makeBadge("SUBMIT CLIENT BEHAVIOUR", 140, Colors.blue1, Colors.blue7)
       | "connection_open_init" =>
         makeBadge("CONNECTION OPEN INIT", 140, Colors.blue1, Colors.blue7)
       | "connection_open_try" =>
         makeBadge("CONNECTION OPEN TRY", 130, Colors.blue1, Colors.blue7)
       | "connection_open_ack" =>
         makeBadge("CONNECTION OPEN ACK", 130, Colors.blue1, Colors.blue7)
       | "connection_open_confirm" =>
         makeBadge("CONNECTION OPEN CONFIRM ", 140, Colors.blue1, Colors.blue7)
       | "channel_open_init" => makeBadge("CHANNEL OPEN INIT", 100, Colors.blue1, Colors.blue7)
       | "channel_open_try" => makeBadge("CHANNEL OPEN TRY", 100, Colors.blue1, Colors.blue7)
       | "channel_open_ack" => makeBadge("CHANNEL OPEN ACK", 100, Colors.blue1, Colors.blue7)
       | "channel_open_confirm" =>
         makeBadge("CHANNEL OPEN CONFIRM", 130, Colors.blue1, Colors.blue7)
       | "channel_close_init" => makeBadge("CHANNEL CLOSE INIT", 85, Colors.blue1, Colors.blue7)
       | "channel_close_confirm" =>
         makeBadge("CHANNEL CLOSE CONFIRM", 100, Colors.blue1, Colors.blue7)
       | "ics04/opaque" => makeBadge("PACKET", 50, Colors.blue1, Colors.blue7)
       | "ics04/timeout" => makeBadge("TIMEOUT", 50, Colors.blue1, Colors.blue7)
       | _ => makeBadge("UNKNOWN", 70, Colors.gray1, Colors.gray6)
       }}
    </div>
  | Unknown => React.null
  };
};
