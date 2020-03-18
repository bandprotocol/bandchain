module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row)]);

  let pageContainer =
    style([
      width(`percent(100.)),
      paddingTop(`px(50)),
      minHeight(`px(450)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
      justifyContent(`center),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(0, 0, 0, 0.1))),
    ]);

  let linkToHome = style([display(`flex), alignItems(`center), cursor(`pointer)]);

  let rightArrow = style([width(`px(20))]);

  let logo = style([width(`px(180)), marginRight(`px(10))]);
};

[@react.component]
let make = () => {
  <>
    <VSpacing size=Spacing.xxl />
    <div className=Styles.pageContainer>
      <Col> <img src=Images.notFoundBg className=Styles.logo /> </Col>
      <VSpacing size=Spacing.xxl />
      <Text
        value="Oops! We cannot find the page you're looking for."
        size=Text.Lg
        color=Colors.blueGray6
      />
      <VSpacing size=Spacing.lg />
      <div className=Styles.linkToHome onClick={_ => Route.redirect(Route.HomePage)}>
        <Text value="Back to Homepage" weight=Text.Bold size=Text.Md color=Colors.blueGray6 />
        <HSpacing size=Spacing.md />
        <img src=Images.rightArrow className=Styles.rightArrow />
      </div>
      <VSpacing size=Spacing.xxl />
    </div>
  </>;
};
