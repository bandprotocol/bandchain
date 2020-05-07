module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      alignItems(`center),
      flexDirection(`column),
      margin2(~v=`vh(12.), ~h=`px(20)),
      textAlign(`center),
    ]);

  let bandLogo = style([width(`vw(20.))]);
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <img src=Images.bandLogo className=Styles.bandLogo />
    <VSpacing size=Spacing.xl />
    <Text
      value="BandChain Explorer is currently not supported on mobile platforms."
      align=Text.Center
      size=Text.Xl
    />
    <VSpacing size=Spacing.md />
    <Text
      value="You can use desktop browsers to interact with Explorer."
      align=Text.Center
      size=Text.Xl
    />
    <VSpacing size=Spacing.xxl />
  </div>;
};
