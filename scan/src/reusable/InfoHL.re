type t =
  | Height(int)
  | MsgCount(int)
  | Timestamp(MomentRe.Moment.t)
  | Fee(float);

let header =
  fun
  | Height(_) => "HEIGHT"
  | MsgCount(_) => "MESSAGES"
  | Timestamp(_) => "TIMESTAMP"
  | Fee(_) => "FEE";

module Styles = {
  open Css;

  let hFlex = style([display(`flex), flexDirection(`column), alignItems(`flexStart)]);
  let vFlex = style([display(`flex), alignItems(`center)]);
  let headerContainer = style([lineHeight(`px(25))]);
};

[@react.component]
let make = (~info) => {
  <div className=Styles.hFlex>
    <div className=Styles.headerContainer>
      <Text value={info |> header} color={Css.hex("555555")} />
    </div>
    {switch (info) {
     | Height(height) =>
       <div className=Styles.vFlex>
         <Text value="#" size=Text.Lg weight=Text.Semibold color={Css.hex("806BFF")} />
         <HSpacing size=Spacing.xs />
         <Text value={height |> Format.iPretty} size=Text.Lg weight=Text.Semibold />
       </div>
     | MsgCount(count) =>
       <Text value={count |> Format.iPretty} size=Text.Lg weight=Text.Semibold />
     | Timestamp(time) =>
       <Text
         value={time |> MomentRe.Moment.format("MMM-DD-YYYY hh:mm:ss a +UTC")}
         size=Text.Lg
         weight=Text.Semibold
       />
     | Fee(fee) =>
       <Text value={(fee |> Format.fPretty) ++ " BAND"} size=Text.Lg weight=Text.Semibold />
     }}
  </div>;
};
