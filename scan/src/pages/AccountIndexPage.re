module Styles = {
  open Css;

  let innerCenter = style([Media.mobile([display(`flex), justifyContent(`center)])]);

  let separatorLine =
    style([
      width(`px(1)),
      height(`px(275)),
      backgroundColor(Colors.gray7),
      marginLeft(`px(20)),
      opacity(0.3),
      Media.mobile([marginLeft(`zero), width(`percent(100.)), height(`px(1))]),
    ]);

  let squareIcon = color =>
    style([width(`px(8)), marginRight(`px(8)), height(`px(8)), backgroundColor(color)]);

  let balance = style([minWidth(`px(150)), justifyContent(`flexEnd)]);

  let totalContainer =
    style([
      display(`flex),
      flexDirection(`column),
      justifyContent(`spaceBetween),
      alignItems(`flexEnd),
      height(`px(200)),
      padding2(~v=`px(12), ~h=`zero),
      Media.mobile([height(`px(100))]),
    ]);

  let infoContainerFullwidth =
    style([
      Media.mobile([
        selector("> div", [flexBasis(`percent(100.))]),
        selector("> div + div", [marginTop(`px(15))]),
      ]),
    ]);

  let totalBalance =
    style([
      display(`flex),
      flexDirection(`column),
      alignItems(`flexEnd),
      Media.mobile([
        flexDirection(`row),
        justifyContent(`spaceBetween),
        alignItems(`center),
        width(`percent(100.)),
      ]),
    ]);

  let button =
    style([
      backgroundColor(Colors.blue1),
      padding2(~h=`px(8), ~v=`px(4)),
      display(`flex),
      borderRadius(`px(6)),
      cursor(`pointer),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(20, 32, 184, `num(0.2)))),
      borderRadius(`px(10)),
    ]);

  let amountBoxes = style([selector("> div + div", [marginTop(`px(24))])]);
};

let balanceDetail = (~title, ~description, ~amount, ~usdPrice, ~color, ~isCountup=false, ()) => {
  <Row.Grid>
    <Col.Grid col=Col.Six>
      <div className={CssHelper.flexBox()}>
        <div className={Styles.squareIcon(color)} />
        <Text
          value=title
          size=Text.Lg
          weight=Text.Semibold
          tooltipItem={description |> React.string}
          tooltipPlacement=Text.AlignBottomStart
        />
      </div>
    </Col.Grid>
    <Col.Grid col=Col.Six>
      <div className={CssHelper.flexBox(~direction=`column, ~align=`flexEnd, ())}>
        <div className={CssHelper.flexBox()}>
          {isCountup
             ? <NumberCountup
                 value=amount
                 size=Text.Lg
                 weight=Text.Regular
                 spacing={Text.Em(0.)}
               />
             : <Text
                 value={amount |> Format.fPretty}
                 size=Text.Lg
                 weight=Text.Regular
                 spacing={Text.Em(0.)}
                 nowrap=true
                 code=true
               />}
          <HSpacing size=Spacing.sm />
          <Text
            value="BAND"
            size=Text.Lg
            code=true
            weight=Text.Thin
            spacing={Text.Em(0.)}
            nowrap=true
          />
        </div>
        <VSpacing size=Spacing.xs />
        <div className={Css.merge([CssHelper.flexBox(), Styles.balance])}>
          {isCountup
             ? <NumberCountup
                 value={amount *. usdPrice}
                 size=Text.Md
                 weight=Text.Thin
                 spacing={Text.Em(0.02)}
                 color=Colors.gray6
               />
             : <Text
                 value={amount *. usdPrice |> Format.fPretty}
                 size=Text.Md
                 spacing={Text.Em(0.02)}
                 weight=Text.Thin
                 nowrap=true
                 code=true
                 color=Colors.gray6
               />}
          <HSpacing size=Spacing.sm />
          <Text
            value="USD"
            size=Text.Md
            code=true
            spacing={Text.Em(0.02)}
            weight=Text.Thin
            nowrap=true
          />
        </div>
      </div>
    </Col.Grid>
  </Row.Grid>;
};

let totalBalanceRender = (isMobile, rawTitle, amount, symbol) => {
  let titles = isMobile ? rawTitle->Js.String2.split("\n") : [|rawTitle|];

  <div className=Styles.totalBalance>
    <div className={CssHelper.flexBox(~direction=`column, ())}>
      {titles
       ->Belt_Array.mapWithIndex((i, title) =>
           <Text
             key={i->string_of_int ++ title}
             value=title
             size={isMobile ? Text.Sm : Text.Md}
             spacing={Text.Em(0.03)}
             height={Text.Px(18)}
           />
         )
       ->React.array}
    </div>
    <VSpacing size=Spacing.md />
    <div className={CssHelper.flexBox()}>
      <NumberCountup
        value=amount
        size={isMobile ? Text.Lg : Text.Xxl}
        weight=Text.Semibold
        spacing={Text.Em(0.02)}
      />
      <HSpacing size=Spacing.sm />
      <Text
        value=symbol
        size={isMobile ? Text.Lg : Text.Xxl}
        weight=Text.Thin
        spacing={Text.Em(0.02)}
        code=true
      />
    </div>
  </div>;
};

let totalBalanceRenderNew = (isMobile, amountBAND, usdPrice) => {
  <>
    <div className={CssHelper.flexBox()}>
      <NumberCountup
        value=amountBAND
        size=Text.Xxxl
        weight=Text.Regular
        spacing={Text.Em(0.)}
        color=Colors.bandBlue
        smallNumber=true
      />
      <HSpacing size=Spacing.sm />
      <Text
        value="BAND"
        color=Colors.bandBlue
        size={isMobile ? Text.Lg : Text.Xxl}
        code=false
        weight=Text.Thin
      />
    </div>
    <div className={CssHelper.flexBox()}>
      <NumberCountup
        value={amountBAND *. usdPrice}
        size=Text.Lg
        weight=Text.Regular
        spacing={Text.Em(0.)}
        color=Colors.gray7
      />
      <HSpacing size=Spacing.sm />
      <Text
        value={"USD " ++ "($" ++ (usdPrice |> string_of_float) ++ " / BAND)"}
        color=Colors.gray6
        size=Text.Lg
        weight=Text.Thin
      />
    </div>
  </>;
};

[@react.component]
let make = (~address, ~hashtag: Route.account_tab_t) => {
  let currentTime =
    React.useContext(TimeContext.context) |> MomentRe.Moment.format(Config.timestampUseFormat);
  let isMobile = Media.isMobile();
  let accountSub = AccountSub.get(address);
  let trackingSub = TrackingSub.use();
  let balanceAtStakeSub = DelegationSub.getTotalStakeByDelegator(address);
  let unbondingSub = UnbondingSub.getUnbondingBalance(address, currentTime);
  let infoSub = React.useContext(GlobalContext.context);
  let (_, dispatchModal) = React.useContext(ModalContext.context);
  let (accountOpt, _) = React.useContext(AccountContext.context);

  let topPartAllSub = Sub.all5(infoSub, accountSub, balanceAtStakeSub, unbondingSub, trackingSub);

  let sumBalance = (balance, amount, unbonding, reward, commission) => {
    let availableBalance = balance->Coin.getBandAmountFromCoins;
    let balanceAtStakeAmount = amount->Coin.getBandAmountFromCoin;
    let unbondingAmount = unbonding->Coin.getBandAmountFromCoin;
    let rewardAmount = reward->Coin.getBandAmountFromCoin;
    let commissionAmount = commission->Coin.getBandAmountFromCoins;

    availableBalance +. balanceAtStakeAmount +. rewardAmount +. unbondingAmount +. commissionAmount;
  };
  let send = chainID => {
    switch (accountOpt) {
    | Some({address: sender}) =>
      let openSendModal = () => Some(address)->SubmitMsg.Send->SubmitTx->OpenModal->dispatchModal;
      if (sender == address) {
        Webapi.Dom.(window |> Window.confirm("Are you sure you want to send tokens to yourself?"))
          ? openSendModal() : ();
      } else {
        openSendModal();
      };
    | None => dispatchModal(OpenModal(Connect(chainID)))
    };
  };

  <Section pbSm=0>
    <div className=CssHelper.container>
      <Row.Grid marginBottom=40 marginBottomSm=24>
        <Col.Grid> <Heading value="Account Detail" size=Heading.H4 /> </Col.Grid>
      </Row.Grid>
      // <div className={CssHelper.flexBox()}>
      //   {switch (topPartAllSub) {
      //    | Data((_, _, _, _, {chainID})) =>
      //      <>
      //        <AddressRender address position=AddressRender.Title copy=true clickable=false />
      //        {isMobile
      //           ? React.null
      //           : <>
      //               <HSpacing size=Spacing.md />
      //               <div
      //                 className={CssHelper.btn(~px=13, ~fsize=10, ~py=5, ())}
      //                 onClick={_ => {send(chainID)}}>
      //                 <Text
      //                   value="Send BAND"
      //                   size=Text.Lg
      //                   block=true
      //                   color=Colors.white
      //                   nowrap=true
      //                 />
      //               </div>
      //             </>}
      //      </>
      //    | _ => <LoadingCensorBar width=600 height=20 />
      //    }}
      // </div>
      // <VSpacing size={isMobile ? Spacing.lg : Spacing.xxl} />
      // <Row justify=Row.Between alignItems=`flexStart wrap=true style=Styles.infoContainerFullwidth>
      //   <Col size=0.75>
      //     <div className=Styles.innerCenter>
      //       {switch (topPartAllSub) {
      //        | Data((_, {balance, commission}, {amount, reward}, unbonding, _)) =>
      //          let availableBalance = balance->Coin.getBandAmountFromCoins;
      //          let balanceAtStakeAmount = amount->Coin.getBandAmountFromCoin;
      //          let unbondingAmount = unbonding->Coin.getBandAmountFromCoin;
      //          let rewardAmount = reward->Coin.getBandAmountFromCoin;
      //          let commissionAmount = commission->Coin.getBandAmountFromCoins;
      //          <PieChart
      //            size={isMobile ? 160 : 187}
      //            availableBalance
      //            balanceAtStake=balanceAtStakeAmount
      //            reward=rewardAmount
      //            unbonding=unbondingAmount
      //            commission=commissionAmount
      //          />;
      //        | _ => <LoadingCensorBar width=160 height=160 radius=160 />
      //        }}
      //     </div>
      //   </Col>
      //   <Col size=1.>
      //     <VSpacing size=Spacing.md />
      //     {switch (topPartAllSub) {
      //      | Data(({financial}, {balance}, _, _, _)) =>
      //        balanceDetail(
      //          ~title="Available Balance",
      //          ~description="Balance available to send, delegate, etc",
      //          ~amount={
      //            balance->Coin.getBandAmountFromCoins;
      //          },
      //          ~usdPrice=financial.usdPrice,
      //          ~color=Colors.bandBlue,
      //          (),
      //        )
      //      | _ => <LoadingCensorBar width=338 height=20 />
      //      }}
      //     <VSpacing size=Spacing.lg />
      //     <VSpacing size=Spacing.md />
      //     {switch (topPartAllSub) {
      //      | Data(({financial}, _, {amount}, _, _)) =>
      //        balanceDetail(
      //          ~title="Balance At Stake",
      //          ~description="Balance currently delegated to validators",
      //          ~amount={
      //            amount->Coin.getBandAmountFromCoin;
      //          },
      //          ~usdPrice=financial.usdPrice,
      //          ~color=Colors.chartBalanceAtStake,
      //          (),
      //        )
      //      | _ => <LoadingCensorBar width=338 height=20 />
      //      }}
      //     <VSpacing size=Spacing.lg />
      //     <VSpacing size=Spacing.md />
      //     {switch (topPartAllSub) {
      //      | Data(({financial}, _, _, unbonding, _)) =>
      //        balanceDetail(
      //          ~title="Unbonding Amount",
      //          ~description="Amount undelegated from validators awaiting 21 days lockup period",
      //          ~amount={
      //            unbonding->Coin.getBandAmountFromCoin;
      //          },
      //          ~usdPrice=financial.usdPrice,
      //          ~color=Colors.blue4,
      //          (),
      //        )
      //      | _ => <LoadingCensorBar width=338 height=20 />
      //      }}
      //     <VSpacing size=Spacing.lg />
      //     <VSpacing size=Spacing.md />
      //     {switch (topPartAllSub) {
      //      | Data(({financial}, _, {reward}, _, _)) =>
      //        balanceDetail(
      //          ~title="Reward",
      //          ~description="Reward from staking to validators",
      //          ~amount={
      //            reward->Coin.getBandAmountFromCoin;
      //          },
      //          ~usdPrice=financial.usdPrice,
      //          ~color=Colors.chartReward,
      //          ~isCountup=true,
      //          (),
      //        )
      //      | _ => <LoadingCensorBar width=338 height=20 />
      //      }}
      //     {switch (topPartAllSub) {
      //      | Data(({financial}, {commission}, _, _, _)) =>
      //        let commissionAmount = commission->Coin.getBandAmountFromCoins;
      //        commissionAmount == 0.
      //          ? React.null
      //          : <>
      //              <VSpacing size=Spacing.lg />
      //              <VSpacing size=Spacing.md />
      //              {balanceDetail(
      //                 ~title="Commission",
      //                 ~description="Reward commission from delegator's reward",
      //                 ~amount=commissionAmount,
      //                 ~usdPrice=financial.usdPrice,
      //                 ~color=Colors.gray6,
      //                 ~isCountup=true,
      //                 (),
      //               )}
      //              <VSpacing size=Spacing.lg />
      //            </>;
      //      | _ =>
      //        <>
      //          <VSpacing size=Spacing.lg />
      //          <VSpacing size=Spacing.md />
      //          <LoadingCensorBar width=338 height=20 />
      //        </>
      //      }}
      //   </Col>
      //   <div className=Styles.separatorLine />
      //   <Col size=1. alignSelf=Col.Start>
      //     <div className=Styles.totalContainer>
      //       {switch (topPartAllSub) {
      //        | Data(({financial}, {balance, commission}, {amount, reward}, unbonding, _)) =>
      //          totalBalanceRender(
      //            isMobile,
      //            "Total BAND In USD \n($"
      //            ++ (financial.usdPrice |> Format.fPretty(~digits=2))
      //            ++ " / BAND)",
      //            sumBalance(balance, amount, unbonding, reward, commission) *. financial.usdPrice,
      //            "USD",
      //          )
      //        | _ => <LoadingCensorBar width=200 height=20 />
      //        }}
      //     </div>
      //   </Col>
      // </Row>
      <Row.Grid>
        <Col.Grid col=Col.Six>
          <Row.Grid marginBottom=24>
            <Col.Grid>
              <div className=CssHelper.infoContainer>
                <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
                  {switch (topPartAllSub) {
                   | Data((_, _, _, _, {chainID})) =>
                     <>
                       {isMobile
                          ? React.null
                          : <>
                              <HSpacing size=Spacing.md />
                              <div
                                className={CssHelper.btn(~variant=Outline, ~fsize=10, ~py=5, ())}
                                onClick={_ => {send(chainID)}}>
                                <Text
                                  value="Send BAND"
                                  block=true
                                  weight=Text.Semibold
                                  color=Colors.bandBlue
                                  nowrap=true
                                />
                              </div>
                            </>}
                     </>
                   | _ => <LoadingCensorBar width=600 height=20 />
                   }}
                </div>
                <Heading size=Heading.H4 value="Address" marginBottom=5 />
                <div className={CssHelper.flexBox()}>
                  <AddressRender
                    address
                    position=AddressRender.Subtitle
                    copy=true
                    clickable=false
                  />
                </div>
              </div>
            </Col.Grid>
          </Row.Grid>
          <Row.Grid>
            <Col.Grid>
              <div className=CssHelper.infoContainer>
                <Heading size=Heading.H4 value="Total Balance" marginBottom=5 />
                {switch (topPartAllSub) {
                 | Data(({financial}, {balance, commission}, {amount, reward}, unbonding, _)) =>
                   totalBalanceRenderNew(
                     isMobile,
                     sumBalance(balance, amount, unbonding, reward, commission),
                     financial.usdPrice,
                   )
                 | _ => <LoadingCensorBar width=200 height=20 />
                 }}
              </div>
            </Col.Grid>
          </Row.Grid>
        </Col.Grid>
        <Col.Grid col=Col.Six>
          <div className=CssHelper.infoContainer>
            {switch (topPartAllSub) {
             | Data((_, {balance, commission}, {amount, reward}, unbonding, _)) =>
               let availableBalance = balance->Coin.getBandAmountFromCoins;
               let balanceAtStakeAmount = amount->Coin.getBandAmountFromCoin;
               let unbondingAmount = unbonding->Coin.getBandAmountFromCoin;
               let rewardAmount = reward->Coin.getBandAmountFromCoin;
               let commissionAmount = commission->Coin.getBandAmountFromCoins;
               <AccountBarChart
                 availableBalance
                 balanceAtStake=balanceAtStakeAmount
                 reward=rewardAmount
                 unbonding=unbondingAmount
                 commission=commissionAmount
               />;
             | _ => <LoadingCensorBar fullWidth=true height=10 />
             }}
            <div className=Styles.amountBoxes>
              <div>
                {switch (topPartAllSub) {
                 | Data(({financial}, {balance}, _, _, _)) =>
                   balanceDetail(
                     ~title="Available Balance",
                     ~description="Balance available to send, delegate, etc",
                     ~amount={
                       balance->Coin.getBandAmountFromCoins;
                     },
                     ~usdPrice=financial.usdPrice,
                     ~color=Colors.bandBlue,
                     (),
                   )
                 | _ => <LoadingCensorBar width=338 height=20 />
                 }}
              </div>
              <div>
                {switch (topPartAllSub) {
                 | Data(({financial}, _, {amount}, _, _)) =>
                   balanceDetail(
                     ~title="Balance At Stake",
                     ~description="Balance currently delegated to validators",
                     ~amount={
                       amount->Coin.getBandAmountFromCoin;
                     },
                     ~usdPrice=financial.usdPrice,
                     ~color=Colors.chartBalanceAtStake,
                     (),
                   )
                 | _ => <LoadingCensorBar width=338 height=20 />
                 }}
              </div>
              <div>
                {switch (topPartAllSub) {
                 | Data(({financial}, _, _, unbonding, _)) =>
                   balanceDetail(
                     ~title="Unbonding Amount",
                     ~description=
                       "Amount undelegated from validators awaiting 21 days lockup period",
                     ~amount={
                       unbonding->Coin.getBandAmountFromCoin;
                     },
                     ~usdPrice=financial.usdPrice,
                     ~color=Colors.blue4,
                     (),
                   )
                 | _ => <LoadingCensorBar width=338 height=20 />
                 }}
              </div>
              <div>
                {switch (topPartAllSub) {
                 | Data(({financial}, _, {reward}, _, _)) =>
                   balanceDetail(
                     ~title="Reward",
                     ~description="Reward from staking to validators",
                     ~amount={
                       reward->Coin.getBandAmountFromCoin;
                     },
                     ~usdPrice=financial.usdPrice,
                     ~color=Colors.chartReward,
                     ~isCountup=true,
                     (),
                   )
                 | _ => <LoadingCensorBar width=338 height=20 />
                 }}
              </div>
              {switch (topPartAllSub) {
               | Data(({financial}, {commission}, _, _, _)) =>
                 let commissionAmount = commission->Coin.getBandAmountFromCoins;
                 commissionAmount == 0.
                   ? React.null
                   : <div>
                       {balanceDetail(
                          ~title="Commission",
                          ~description="Reward commission from delegator's reward",
                          ~amount=commissionAmount,
                          ~usdPrice=financial.usdPrice,
                          ~color=Colors.gray6,
                          ~isCountup=true,
                          (),
                        )}
                     </div>;

               | _ => <div> <LoadingCensorBar width=338 height=20 /> </div>
               }}
            </div>
          </div>
        </Col.Grid>
      </Row.Grid>
      <VSpacing size=Spacing.xl />
      <Tab
        tabs=[|
          {
            name: "Transactions",
            route: Route.AccountIndexPage(address, Route.AccountTransactions),
          },
          {
            name: "Delegations",
            route: Route.AccountIndexPage(address, Route.AccountDelegations),
          },
          {name: "Unbonding", route: Route.AccountIndexPage(address, Route.AccountUnbonding)},
          {name: "Redelegate", route: Route.AccountIndexPage(address, Route.AccountRedelegate)},
        |]
        currentRoute={Route.AccountIndexPage(address, hashtag)}>
        {switch (hashtag) {
         | AccountTransactions => <AccountIndexTransactions accountAddress=address />
         | AccountDelegations => <AccountIndexDelegations address />
         | AccountUnbonding => <AccountIndexUnbonding address />
         | AccountRedelegate => <AccountIndexRedelegate address />
         }}
      </Tab>
    </div>
  </Section>;
};
