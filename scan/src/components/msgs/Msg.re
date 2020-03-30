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

[@react.component]
let make = (~msg: TxSub.Msg.t, ~width: int, ~success: bool) => {
  switch (msg) {
  | Send({fromAddress, toAddress, amount}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(width / 2 - 20)}>
        <AddressRender address=fromAddress />
      </div>
      <div className={Styles.withBg(Colors.blue1, 40)}>
        <Text
          value="SEND"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.blue7
        />
      </div>
      {success
         ? <>
             <div className={Styles.rowWithWidth(200)}>
               <Text
                 value={
                   amount
                   ->Belt_List.get(0)
                   ->Belt_Option.getWithDefault(TxSub.Coin.newCoin("uband", 0.0)).
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
               <Text
                 value={j|➜|j}
                 size=Text.Xxl
                 weight=Text.Bold
                 code=true
                 nowrap=true
                 block=true
               />
               <HSpacing size=Spacing.sm />
             </div>
             <div className={Styles.withWidth(width / 2 - 18)}>
               <AddressRender address=toAddress />
             </div>
           </>
         : React.null}
    </div>
  | CreateDataSource({id, sender, name}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(150)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.yellow1, 110)}>
        <Text
          value="CREATE DATASOURCE"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.yellow6
        />
      </div>
      {success
         ? <>
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
           </>
         : React.null}
    </div>
  | EditDataSource({id, sender, name}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(150)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.yellow1, 100)}>
        <Text
          value="EDIT DATASOURCE"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.yellow6
        />
      </div>
      {success
         ? <>
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
           </>
         : React.null}
    </div>
  | CreateOracleScript({id, sender, name}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(140)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.pink1, 120)}>
        <Text
          value="CREATE ORACLE SCRIPT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.pink6
        />
      </div>
      {success
         ? <>
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
           </>
         : React.null}
    </div>
  | EditOracleScript({id, sender, name}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(140)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.pink1, 110)}>
        <Text
          value="EDIT ORACLE SCRIPT"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.pink6
        />
      </div>
      {success
         ? <>
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
           </>
         : React.null}
    </div>
  | Request({id, oracleScriptID, sender}) =>
    <div className={Styles.rowWithWidth(width)}>
      <div className={Styles.withWidth(140)}> <AddressRender address=sender /> </div>
      <div className={Styles.withBg(Colors.orange1, 60)}>
        <Text
          value="REQUEST"
          size=Text.Xs
          spacing={Text.Em(0.07)}
          weight=Text.Medium
          color=Colors.orange6
        />
      </div>
      {success
         ? <>
             <TypeID.Request id />
             <HSpacing size=Spacing.sm />
             <Text
               value={j|➜|j}
               size=Text.Xxl
               weight=Text.Bold
               code=true
               nowrap=true
               block=true
             />
             <HSpacing size=Spacing.sm />
             <TypeID.OracleScript id=oracleScriptID />
             <HSpacing size=Spacing.sm />
             <Text
               value="Mock Oracle Script" // TODO , replace with wire up data
               color=Colors.gray7
               weight=Text.Medium
               nowrap=true
               block=true
               ellipsis=true
             />
           </>
         : React.null}
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
      {success
         ? <>
             <Text
               value={j|➜|j}
               size=Text.Xxl
               weight=Text.Bold
               code=true
               nowrap=true
               block=true
             />
             <HSpacing size=Spacing.sm />
             <TypeID.Request id=requestID />
           </>
         : React.null}
    </div>
  | AddOracleAddress(_) => React.null
  | RemoveOracleAddress(_) => React.null
  | CreateValidator(_) => React.null
  | EditValidator(_) => React.null
  | Unknown => React.null
  };
};
