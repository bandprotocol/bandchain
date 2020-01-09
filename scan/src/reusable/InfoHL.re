type t =
  | Height(int)
  | Count(int)
  | Timestamp(MomentRe.Moment.t)
  | Fee(float);

module Styles = {
  open Css;

  let hFlex = style([display(`flex), flexDirection(`column), alignItems(`flexStart)]);
  let vFlex = style([display(`flex), alignItems(`center)]);
  let headerContainer = style([lineHeight(`px(25))]);
};

[@react.component]
let make = (~info, ~header) => {
  <div className=Styles.hFlex>
    <div className=Styles.headerContainer> <Text value=header color=Colors.darkerGrayText /> </div>
    {switch (info) {
     | Height(height) =>
       <div className=Styles.vFlex>
         <Text value="#" size=Text.Lg weight=Text.Semibold color={Css.hex("806BFF")} />
         <HSpacing size=Spacing.xs />
         <Text value={height |> Format.iPretty} size=Text.Lg weight=Text.Semibold />
       </div>
     | Count(count) => <Text value={count |> Format.iPretty} size=Text.Lg weight=Text.Semibold />
     | Timestamp(time) =>
       <Text
         value={time |> MomentRe.Moment.format("MMM-DD-YYYY hh:mm:ss A [GMT]Z")}
         size=Text.Lg
         weight=Text.Semibold
       />
     | Fee(fee) =>
       <Text value={(fee |> Format.fPretty) ++ " BAND"} size=Text.Lg weight=Text.Semibold />
     }}
  </div>;
};
