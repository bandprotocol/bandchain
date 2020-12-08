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
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(0, 0, 0, `num(0.1)))),
    ]);

  let linkToHome = style([display(`flex), alignItems(`center), cursor(`pointer)]);

  let rightArrow = style([width(`px(20)), filter([`saturate(50.0), `brightness(70.0)])]);

  let logo = style([width(`px(180)), marginRight(`px(10))]);
};

[@react.component]
let make = () => {
  <Section>
    <div className=CssHelper.container>
      <VSpacing size=Spacing.xxl />
      <div className=Styles.pageContainer>
        <div className={CssHelper.flexBox()}>
          <img src=Images.notFoundBg className=Styles.logo />
        </div>
        <VSpacing size=Spacing.xxl />
        <Text
          value="Oops! We cannot find the page you're looking for."
          size=Text.Lg
          color=Colors.blueGray6
        />
        <VSpacing size=Spacing.lg />
        <Link className=Styles.linkToHome route=Route.HomePage>
          <Text value="Back to Homepage" weight=Text.Bold size=Text.Md color=Colors.blueGray6 />
          <HSpacing size=Spacing.md />
          <img src=Images.rightArrow className=Styles.rightArrow />
        </Link>
        <VSpacing size=Spacing.xxl />
      </div>
    </div>
  </Section>;
};
