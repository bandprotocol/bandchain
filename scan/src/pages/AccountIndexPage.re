module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let cFlex = style([display(`flex), flexDirection(`column), alignItems(`flexEnd)]);

  let rFlex = style([display(`flex), flexDirection(`row)]);

  let separatorLine =
    style([
      width(`px(1)),
      height(`px(200)),
      backgroundColor(Colors.gray7),
      marginLeft(`px(20)),
      opacity(0.3),
    ]);

  let ovalIcon = color =>
    style([
      width(`px(17)),
      height(`px(17)),
      backgroundColor(color),
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

let totalBalanceRender = (title, amount, symbol) => {
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
let make = (~address, ~hashtag: Route.account_tab_t) =>
  {
    let accountSub = AccountSub.get(address);
    let infoSub = React.useContext(GlobalContext.context);
    let balanceAtStakeSub = DelegationSub.getTotalStakeByDelegator(address);

    let%Sub info = infoSub;
    let%Sub account = accountSub;
    let%Sub balanceAtStake = balanceAtStakeSub;

    let availableBalance = account.balance->Coin.getBandAmountFromCoins;
    let usdPrice = info.financial.usdPrice;
    let totalBalance = availableBalance +. balanceAtStake.amount;

    <>
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
              color=Colors.gray7
              block=true
            />
          </div>
        </Col>
      </Row>
      <VSpacing size=Spacing.lg />
      <VSpacing size=Spacing.sm />
      <div className=Styles.vFlex> <AddressRender address position=AddressRender.Title /> </div>
      <VSpacing size=Spacing.xxl />
      <Row justify=Row.Between alignItems=`flexStart>
        <Col size=0.75>
          <PieChart size=187 availableBalance balanceAtStake={balanceAtStake.amount} reward={balanceAtStake.reward} />
        </Col>
        <Col size=1.>
          <VSpacing size=Spacing.md />
          {balanceDetail(
             "AVAILABLE BALANCE",
             availableBalance |> Format.fPretty,
             availableBalance *. usdPrice |> Format.fPretty,
             Colors.bandBlue,
           )}
          <VSpacing size=Spacing.xl />
          <VSpacing size=Spacing.md />
          {balanceDetail(
             "BALANCE AT STAKE",
             balanceAtStake.amount |> Format.fPretty,
             balanceAtStake.amount *. usdPrice |> Format.fPretty,
             Colors.chartBalanceAtStake,
           )}
          <VSpacing size=Spacing.xl />
          <VSpacing size=Spacing.md />
          {balanceDetail(
             "REWARD",
             balanceAtStake.reward |> Format.fPretty,
             balanceAtStake.reward *. usdPrice |> Format.fPretty,
             Colors.chartReward,
           )}
        </Col>
        <div className=Styles.separatorLine />
        <Col size=1. alignSelf=Col.Start>
          <div className=Styles.totalContainer>
            {totalBalanceRender("TOTAL BAND BALANCE", totalBalance |> Format.fPretty, "BAND")}
            {totalBalanceRender(
               "TOTAL BAND IN USD ($" ++ (usdPrice |> Format.fPretty) ++ " / BAND)",
               totalBalance *. usdPrice |> Format.fPretty,
               "USD",
             )}
          </div>
        </Col>
      </Row>
      <VSpacing size=Spacing.xl />
      <Tab
        tabs=[|
          {
            name: "TRANSACTIONS",
            route: Route.AccountIndexPage(address, Route.AccountTransactions),
          },
          {
            name: "DELEGATIONS",
            route: Route.AccountIndexPage(address, Route.AccountDelegations),
          },
        |]
        currentRoute={Route.AccountIndexPage(address, hashtag)}>
        {switch (hashtag) {
         | AccountTransactions => <AccountIndexTransactions accountAddress=address />
         | AccountDelegations => <AccountIndexDelegations address />
         }}
      </Tab>
    </>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
