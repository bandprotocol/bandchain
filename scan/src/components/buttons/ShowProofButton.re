[@react.component]
let make = (~showProof: bool, ~setShowProof) => {
  <div className={CssHelper.btn(~px=20, ~py=12, ())} onClick={_ => setShowProof(_ => !showProof)}>
    <div className={CssHelper.flexBox()}>
      <Icon name="fal fa-long-arrow-down" color=Colors.white />
      <HSpacing size=Spacing.sm />
      <Text
        value={(showProof ? "Hide" : "Show") ++ (Media.isMobile() ? " Proof" : " Proof JSON")}
        weight=Text.Medium
        block=true
        nowrap=true
      />
    </div>
  </div>;
};
