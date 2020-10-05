[@react.component]
let make = (~showProof: bool, ~setShowProof) => {
  <Button px=20 py=12 pxSm=12 pySm=10 onClick={_ => setShowProof(_ => !showProof)}>
    <div className={CssHelper.flexBox()}>
      <Icon
        name={showProof ? "fal fa-long-arrow-up" : "fal fa-long-arrow-down"}
        color=Colors.white
      />
      <HSpacing size=Spacing.sm />
      <Text
        value={(showProof ? "Hide" : "Show") ++ (Media.isMobile() ? " Proof" : " Proof JSON")}
        weight=Text.Medium
        block=true
        nowrap=true
      />
    </div>
  </Button>;
};
