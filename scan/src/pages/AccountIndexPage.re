module Styles = {
  open Css;

  let pageContainer = style([paddingTop(`px(40))]);

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let graph = style([width(`px(186))]);

  let cFlex = style([display(`flex), flexDirection(`column), alignItems(`flexEnd)]);

  let rFlex = style([display(`flex), flexDirection(`row)]);

  let separatorLine =
    style([
      width(`px(1)),
      height(`px(200)),
      backgroundColor(Colors.mediumGray),
      marginLeft(`px(20)),
      opacity(0.3),
    ]);

  let ovalIcon = color =>
    style([
      width(`px(17)),
      height(`px(17)),
      backgroundColor(Css.hex(color)),
      borderRadius(`percent(50.)),
    ]);

  let balance = style([minWidth(`px(150)), justifyContent(`flexEnd)]);

  let totalContainer =
    style([
      display(`flex),
      flexDirection(`column),
      justifyContent(`spaceBetween),
      alignItems(`flexEnd),
      height(`px(190)),
      padding2(~v=`px(12), ~h=`zero),
    ]);

  let totalBalance = style([display(`flex), flexDirection(`column), alignItems(`flexEnd)]);
};

let balanceDetail = (title, amount, amountUsd, color) => {
  <Row alignItems=Css.flexStart>
    <Col size=0.25> <div className={Styles.ovalIcon(color)} /> </Col>
    <Col size=1.2>
      <Text value=title size=Text.Sm height={Text.Px(18)} spacing={Text.Em(0.03)} nowrap=true />
    </Col>
    <Col size=0.6>
      <div className=Styles.cFlex>
        <div className=Styles.rFlex>
          <Text
            value=amount
            size=Text.Lg
            weight=Text.Semibold
            spacing={Text.Em(0.02)}
            nowrap=true
            code=true
          />
          <HSpacing size=Spacing.sm />
          <Text
            value="BAND"
            size=Text.Lg
            code=true
            weight=Text.Thin
            spacing={Text.Em(0.02)}
            nowrap=true
          />
        </div>
        <VSpacing size=Spacing.xs />
        <div className={Css.merge([Styles.rFlex, Styles.balance])}>
          <Text
            value=amountUsd
            size=Text.Sm
            spacing={Text.Em(0.02)}
            weight=Text.Thin
            nowrap=true
            code=true
          />
          <HSpacing size=Spacing.sm />
          <Text
            value="USD"
            size=Text.Sm
            code=true
            spacing={Text.Em(0.02)}
            weight=Text.Thin
            nowrap=true
          />
        </div>
      </div>
    </Col>
  </Row>;
};

let totalBalance = (title, amount, symbol) => {
  <div className=Styles.totalBalance>
    <Text value=title size=Text.Md spacing={Text.Em(0.03)} height={Text.Px(18)} />
    <VSpacing size=Spacing.md />
    <div className=Styles.rFlex>
      <Text
        value=amount
        size=Text.Xxl
        weight=Text.Semibold
        code=true
        spacing={Text.Em(0.02)}
        nowrap=true
      />
      <HSpacing size=Spacing.sm />
      <Text value=symbol size=Text.Xxl weight=Text.Thin spacing={Text.Em(0.02)} code=true />
    </div>
  </div>;
};

[@react.component]
let make = (~address, ~hashtag: Route.account_tab_t) => {
  <div className=Styles.pageContainer>
    <Row justify=Row.Between>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.accountLogo className=Styles.logo />
          <Text
            value="ACCOUNT DETAIL"
            weight=Text.Medium
            size=Text.Md
            spacing={Text.Em(0.06)}
            height={Text.Px(15)}
            nowrap=true
            color=Colors.mediumGray
            block=true
          />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.sm />
    <div className=Styles.vFlex> <AddressRender address position=AddressRender.Title /> </div>
    <VSpacing size=Spacing.xxl />
    <Row justify=Row.Between>
      <Col size=0.75> <img src=Images.pieChart className=Styles.graph /> </Col>
      <Col size=1.>
        {balanceDetail("AVAILABLE BALANCE", "10,547,434.89", "4,829,360.21", "5269FF")}
        <VSpacing size=Spacing.xl />
        <VSpacing size=Spacing.md />
        {balanceDetail("BALANCE AT STAKE", "1,800,000.00", "158,303.88", "ABB6FF")}
        <VSpacing size=Spacing.xl />
        <VSpacing size=Spacing.md />
        {balanceDetail("REWARD", "61,301.04", "60.31", "000C5C")}
      </Col>
      <div className=Styles.separatorLine />
      <Col size=1. alignSelf=Col.FlexStart>
        <div className=Styles.totalContainer>
          {totalBalance("TOTAL BAND BALANCE", "12,408,746.93", "BAND")}
          {totalBalance("TOTAL BAND IN USD ($3.42 / BAND)", "37,226,240.79", "USD")}
        </div>
      </Col>
    </Row>
  </div>;
};
