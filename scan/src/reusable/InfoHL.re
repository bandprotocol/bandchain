type t =
  | Height(int)
  | Count(int)
  | Timestamp(MomentRe.Moment.t)
  | Fee(float)
  | DataSources(list(string))
  | Hash(string, Css.Types.Color.t);

module Styles = {
  open Css;

  let hFlex = style([display(`flex), flexDirection(`column), alignItems(`flexStart)]);
  let vFlex = style([display(`flex), alignItems(`center)]);
  let headerContainer = style([lineHeight(`px(25))]);
  let sourceContainer =
    style([display(`inlineFlex), alignItems(`center), marginRight(`px(40))]);
  let sourceIcon = style([width(`px(16)), marginRight(`px(8))]);
};

[@react.component]
let make = (~info, ~header) => {
  <div className=Styles.hFlex>
    <div className=Styles.headerContainer> <Text value=header color=Colors.grayHeader /> </div>
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
         weight=Text.Bold
       />
     | Fee(fee) =>
       <Text value={(fee |> Format.fPretty) ++ " BAND"} size=Text.Lg weight=Text.Bold />
     | DataSources(sources) =>
       <div className=Styles.vFlex>
         {sources
          ->Belt.List.map(source =>
              <div className=Styles.sourceContainer>
                <img src=Images.source className=Styles.sourceIcon />
                <Text value=source weight=Text.Bold size=Text.Lg />
              </div>
            )
          ->Array.of_list
          ->React.array}
       </div>
     | Hash(hash, textColor) =>
       <Text value=hash size=Text.Lg weight=Text.Semibold color=textColor />
     }}
  </div>;
};
