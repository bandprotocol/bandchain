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
      backgroundColor(Colors.gray7),
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
  let balanceOpt = AccountHook.getBalance(address);
  let priceOpt = PriceHook.get();
  let balanceStakeOpt = AccountHook.getBalanceStake(address);

  let delegations = AccountHook.getDelegations(address) |> Belt_Option.getWithDefault(_, []);
  let rewardOpt = AccountHook.getReward(address);
  let usdPrice = {
    let%Opt price = priceOpt;
    Some(1. /. price.usdPrice);
  };

  let avialableBalance = {
    let%Opt balanceStake = balanceStakeOpt;
    let%Opt balance = balanceOpt;
    Some(balance -. balanceStake);
  };

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
    <Row justify=Row.Between>
      <Col size=0.75> <img src=Images.pieChart className=Styles.graph /> </Col>
      <Col size=1.>
        {switch (avialableBalance, usdPrice) {
         | (Some(avBalance), Some(price)) =>
           balanceDetail(
             "AVAILABLE BALANCE",
             avBalance |> Format.fPretty,
             avBalance *. price |> Format.fPretty,
             "5269FF",
           )
         | _ => balanceDetail("AVAILABLE BALANCE", "?", "?", "5269FF")
         }}
        <VSpacing size=Spacing.xl />
        <VSpacing size=Spacing.md />
        {switch (balanceStakeOpt, usdPrice) {
         | (Some(balanceStake), Some(price)) =>
           balanceDetail(
             "BALANCE AT STAKE",
             balanceStake |> Format.fPretty,
             balanceStake *. price |> Format.fPretty,
             "ABB6FF",
           )
         | _ => balanceDetail("BALANCE AT STAKE", "?", "?", "ABB6FF")
         }}
        <VSpacing size=Spacing.xl />
        <VSpacing size=Spacing.md />
        {switch (rewardOpt, usdPrice) {
         | (Some(reward), Some(price)) =>
           balanceDetail(
             "REWARD",
             reward |> Format.fPretty,
             reward *. price |> Format.fPretty,
             "000C5C",
           )
         | _ => balanceDetail("REWARD", "?", "?", "000C5C")
         }}
      </Col>
      <div className=Styles.separatorLine />
      <Col size=1. alignSelf=Col.Start>
        <div className=Styles.totalContainer>
          {switch (balanceOpt, usdPrice) {
           | (Some(balance), Some(price)) =>
             <>
               {totalBalance("TOTAL BAND BALANCE", balance |> Format.fPretty, "BAND")}
               {totalBalance(
                  "TOTAL BAND IN USD ($3.42 / BAND)",
                  balance *. price |> Format.fPretty,
                  "USD",
                )}
             </>
           | _ =>
             <>
               {totalBalance("TOTAL BAND BALANCE", "?", "BAND")}
               {totalBalance("TOTAL BAND IN USD ($3.42 / BAND)", "?", "USD")}
             </>
           }}
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
        {name: "DELEGATIONS", route: Route.AccountIndexPage(address, Route.AccountDelegations)},
      |]
      currentRoute={Route.AccountIndexPage(address, hashtag)}>
      {switch (hashtag) {
       | AccountTransactions => <AccountIndexTransactions />
       | AccountDelegations => <AccountIndexDelegations delegations />
       }}
    </Tab>
  </div>;
};
