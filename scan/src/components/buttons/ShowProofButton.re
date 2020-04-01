module Styles = {
  open Css;

  let button =
    style([
      backgroundColor(Colors.green1),
      padding2(~h=`px(8), ~v=`px(4)),
      display(`flex),
      width(`px(120)),
      borderRadius(`px(6)),
      cursor(`pointer),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(20, 32, 184, 0.2))),
    ]);

  let withHeight = showProof =>
    style([maxHeight(`px(12)), transform(`rotate(`deg(showProof ? 180. : 0.)))]);
};

[@react.component]
let make = (~showProof: bool, ~setShowProof) => {
  <div className=Styles.button onClick={_ => setShowProof(_ => !showProof)}>
    <img src=Images.showProofArrow className={Styles.withHeight(showProof)} />
    <HSpacing size=Spacing.sm />
    <Text
      value={(showProof ? "Hide" : "Show") ++ " Proof JSON"}
      size=Text.Sm
      block=true
      color=Colors.green5
      nowrap=true
    />
  </div>;
};
