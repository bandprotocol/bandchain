type t =
  | Height(int)
  | Count(int)
  | Float(float)
  | Text(string)
  | Timestamp(MomentRe.Moment.t)
  | Fee(float)
  | DataSources(list(ID.DataSource.t))
  | OracleScript(ID.OracleScript.t, string)
  | Hash(Hash.t, Css.Types.Color.t)
  | Address(Address.t, int)
  | Fraction(int, int, bool)
  | FloatWithSuffix(float, string);

module Styles = {
  open Css;

  let hFlex = isLeft =>
    style([
      display(`flex),
      flexDirection(`column),
      alignItems(isLeft ? `flexStart : `flexEnd),
    ]);
  let vFlex = style([display(`flex), alignItems(`center)]);
  let addressContainer = maxwidth_ => style([alignItems(`center), maxWidth(`px(maxwidth_))]);
  let datasourcesContainer = style([display(`flex), alignItems(`center), flexWrap(`wrap)]);
  let headerContainer = style([lineHeight(`px(25))]);
  let sourceContainer =
    style([
      display(`inlineFlex),
      alignItems(`center),
      marginRight(`px(40)),
      marginTop(`px(13)),
    ]);
  let sourceIcon = style([width(`px(16)), marginRight(`px(8))]);
};

[@react.component]
let make = (~info, ~header, ~isLeft=true) => {
  let infoOpt = React.useContext(GlobalContext.context);
  <div className={Styles.hFlex(isLeft)}>
    <div className=Styles.headerContainer>
      <Text
        value=header
        color=Colors.gray7
        size=Text.Sm
        weight=Text.Thin
        height={Text.Px(13)}
        spacing={Text.Em(0.03)}
      />
    </div>
    {switch (info) {
     | Height(height) =>
       <div className=Styles.vFlex>
         <TypeID.Block id={ID.Block.ID(height)} position=TypeID.Subtitle />
       </div>
     | Float(value) =>
       <Text
         value={value |> Js.Float.toString}
         size=Text.Lg
         weight=Text.Semibold
         spacing={Text.Em(0.02)}
         code=true
       />
     | FloatWithSuffix(value, suffix) =>
       <Text
         value={(value |> Js.Float.toString) ++ suffix}
         size=Text.Lg
         weight=Text.Semibold
         spacing={Text.Em(0.02)}
         code=true
       />
     | Count(value) =>
       <Text
         value={value |> Format.iPretty}
         size=Text.Lg
         weight=Text.Semibold
         spacing={Text.Em(0.02)}
         code=true
       />
     | Text(text) =>
       <Text value=text size=Text.Lg weight=Text.Semibold code=true spacing={Text.Em(0.02)} />
     | Timestamp(time) =>
       <div className=Styles.vFlex>
         <Text
           value={
             time
             |> MomentRe.Moment.format("MMM-DD-YYYY  hh:mm:ss A [+UTC]")
             |> String.uppercase_ascii
           }
           size=Text.Lg
           weight=Text.Semibold
           spacing={Text.Em(0.02)}
           code=true
         />
         <HSpacing size=Spacing.sm />
         <TimeAgos
           time
           prefix="("
           suffix=")"
           size=Text.Md
           weight=Text.Thin
           spacing={Text.Em(0.06)}
           color=Colors.gray7
         />
       </div>
     | Fee(fee) =>
       <div className=Styles.vFlex>
         <Text value={fee |> Format.fPretty} size=Text.Lg weight=Text.Bold code=true />
         <HSpacing size=Spacing.md />
         <Text value="BAND" size=Text.Lg weight=Text.Regular spacing={Text.Em(0.02)} code=true />
         <HSpacing size=Spacing.xs />
         <HSpacing size=Spacing.xs />
         {switch (infoOpt) {
          | Some(info) =>
            let feeInUsd =
              info.financial.usdPrice *. fee |> Js.Float.toFixedWithPrecision(~digits=2);
            <Text
              value={j|(\$$feeInUsd)|j}
              size=Text.Lg
              weight=Text.Regular
              spacing={Text.Em(0.02)}
              code=true
            />;
          | None => React.null
          }}
       </div>
     | DataSources(ids) =>
       <div className=Styles.datasourcesContainer>
         {ids
          ->Belt.List.map(id =>
              <> <TypeID.DataSource id position=TypeID.Subtitle /> <HSpacing size=Spacing.sm /> </>
            )
          ->Array.of_list
          ->React.array}
       </div>
     | OracleScript(id, name) =>
       <div className=Styles.datasourcesContainer>
         <TypeID.OracleScript id position=TypeID.Subtitle />
         <HSpacing size=Spacing.sm />
         <Text value=name size=Text.Lg weight=Text.Regular spacing={Text.Em(0.02)} code=true />
       </div>
     | Hash(hash, textColor) =>
       <Text
         value={hash |> Hash.toHex(~with0x=true, ~upper=true)}
         size=Text.Lg
         weight=Text.Semibold
         color=textColor
       />
     | Fraction(x, y, space) =>
       <Text
         value={(x |> Format.iPretty) ++ (space ? " / " : "/") ++ (y |> Format.iPretty)}
         size=Text.Lg
         weight=Text.Semibold
         spacing={Text.Em(0.02)}
         code=true
       />
     | Address(address, maxWidth) =>
       <div className={Styles.addressContainer(maxWidth)}>
         <AddressRender address position=AddressRender.Subtitle />
       </div>
     }}
  </div>;
};
