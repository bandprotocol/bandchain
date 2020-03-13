type t =
  | Height(int)
  | Count(int)
  | Float(float)
  | Text(string)
  | Timestamp(MomentRe.Moment.t)
  | Fee(float)
  | DataSources(list(string))
  | Hash(Hash.t, Css.Types.Color.t)
  | Address(Address.t)
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
  let addressContainer = style([display(`flex), alignItems(`center), maxWidth(`px(290))]);
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
  <div className={Styles.hFlex(isLeft)}>
    <div className=Styles.headerContainer>
      <Text
        value=header
        color=Colors.mediumGray
        size=Text.Sm
        weight=Text.Thin
        height={Text.Px(13)}
        spacing={Text.Em(0.03)}
      />
    </div>
    {switch (info) {
     | Height(height) =>
       <div className=Styles.vFlex>
         <Text value="#" size=Text.Lg weight=Text.Semibold color=Colors.brightPurple />
         <HSpacing size=Spacing.xs />
         <Text value={height |> Format.iPretty} size=Text.Lg weight=Text.Semibold />
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
     | Text(text) => <Text value=text size=Text.Lg weight=Text.Semibold />
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
         <Text
           value="(9 hrs 2 mins ago)"
           size=Text.Lg
           spacing={Text.Em(0.02)}
           weight=Text.Thin
           code=true
         />
       </div>
     | Fee(fee) =>
       <div className=Styles.vFlex>
         <Text value={fee |> Format.fPretty} size=Text.Lg weight=Text.Bold code=true />
         <HSpacing size=Spacing.md />
         <Text value="BAND" size=Text.Lg weight=Text.Regular spacing={Text.Em(0.02)} code=true />
         <HSpacing size=Spacing.xs />
         <HSpacing size=Spacing.xs />
         <Text
           value="($0.3)"
           size=Text.Lg
           weight=Text.Regular
           spacing={Text.Em(0.02)}
           code=true
         />
       </div>
     | DataSources(sources) =>
       <div className=Styles.datasourcesContainer>
         {sources
          ->Belt.List.map(source =>
              <div key=source className=Styles.sourceContainer>
                <img src=Images.source className=Styles.sourceIcon />
                <Text value=source weight=Text.Bold size=Text.Lg />
              </div>
            )
          ->Array.of_list
          ->React.array}
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

     | Address(address) =>
       <div className=Styles.addressContainer>
         <AddressRender address position=AddressRender.Subtitle />
       </div>
     }}
  </div>;
};
