module Styles = {
  open Css;

  let squareIcon = color =>
    style([width(`px(8)), marginRight(`px(8)), height(`px(8)), backgroundColor(color)]);

  let balance = style([minWidth(`px(150)), justifyContent(`flexEnd)]);

  let infoHeader =
    style([borderBottom(`px(1), `solid, Colors.gray9), paddingBottom(`px(16))]);

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

  let infoLeft =
    style([
      height(`percent(100.)),
      selector(
        "> div",
        [
          height(`calc((`sub, `percent(50.), `px(12)))),
          width(`percent(100.)),
          Media.mobile([height(`auto)]),
        ],
      ),
    ]);

  let amountBoxes = style([selector("> div + div", [marginTop(`px(18))])]);

  let qrContainer = style([width(`percent(100.)), Media.mobile([width(`auto)])]);

  let qrCode =
    style([
      backgroundColor(Colors.bandBlue),
      borderRadius(`px(4)),
      padding(`px(10)),
      cursor(`pointer),
      Media.mobile([marginRight(`px(8))]),
    ]);

  let addressContainer =
    style([Media.mobile([width(`calc((`sub, `percent(100.), `px(50))))])]);
};

let balanceDetail = (~title, ~description, ~amount, ~usdPrice, ~color, ~isCountup=false, ()) => {
  <Row>
    <Col.Grid col=Col.Six colSm=Col.Five>
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
    <Col.Grid col=Col.Six colSm=Col.Seven>
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
            color=Colors.gray6
          />
        </div>
      </div>
    </Col.Grid>
  </Row>;
};

module BalanceDetailLoading = {
  [@react.component]
  let make = () => {
    <Row>
      <Col.Grid col=Col.Six colSm=Col.Five> <LoadingCensorBar width=130 height=18 /> </Col.Grid>
      <Col.Grid col=Col.Six colSm=Col.Seven>
        <div className={CssHelper.flexBox(~direction=`column, ~align=`flexEnd, ())}>
          <LoadingCensorBar width=120 height=20 />
          <VSpacing size=Spacing.xs />
          <LoadingCensorBar width=120 height=16 />
        </div>
      </Col.Grid>
    </Row>;
  };
};

let totalBalanceRender = (amountBAND, usdPrice) => {
  <>
    <div
      className={Css.merge([CssHelper.flexBox(~align=`flexEnd, ()), CssHelper.mb(~size=5, ())])}>
      <NumberCountup
        value=amountBAND
        size=Text.Xxxl
        weight=Text.Regular
        spacing={Text.Em(0.)}
        color=Colors.bandBlue
        smallNumber=true
      />
      <HSpacing size=Spacing.sm />
      <Text value="BAND" color=Colors.bandBlue size=Text.Lg code=false weight=Text.Thin />
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
        value={"USD " ++ "($" ++ (usdPrice |> Js.Float.toString) ++ " / BAND)"}
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

  let qrCode = () => {
    address->QRCode->OpenModal->dispatchModal;
  };

  <Section pbSm=0>
    <div className=CssHelper.container>
      <Row marginBottom=40 marginBottomSm=24>
        <Col.Grid> <Heading value="Account Detail" size=Heading.H4 /> </Col.Grid>
      </Row>
      <Row>
        <Col.Grid col=Col.Six>
          <div
            className={Css.merge([
              CssHelper.flexBox(~direction=`column, ~justify=`spaceBetween, ~align=`stretch, ()),
              Styles.infoLeft,
            ])}>
            <div
              className={Css.merge([
                CssHelper.infoContainer,
                CssHelper.flexBox(~direction=`column, ~justify=`center, ~align=`stretch, ()),
                CssHelper.flexBoxSm(~direction=`row, ~align=`center, ~justify=`flexStart, ()),
                CssHelper.mb(~size=24, ()),
              ])}>
              <div
                className={Css.merge([
                  CssHelper.flexBox(~justify=`spaceBetween, ~align=`flexStart, ()),
                  CssHelper.mb(~size=24, ()),
                  CssHelper.mbSm(~size=0, ()),
                  Styles.qrContainer,
                ])}>
                <div className=Styles.qrCode onClick={_ => {qrCode()}}>
                  <Icon size=20 name="far fa-qrcode" color=Colors.white />
                </div>
                {isMobile
                   ? React.null
                   : {
                     switch (topPartAllSub) {
                     | Data((_, _, _, _, {chainID})) =>
                       <Button variant=Button.Outline py=5 px=11 onClick={_ => {send(chainID)}}>
                         <Text
                           value="Send BAND"
                           block=true
                           weight=Text.Semibold
                           color=Colors.bandBlue
                           nowrap=true
                         />
                       </Button>
                     | _ => <LoadingCensorBar width=90 height=26 />
                     };
                   }}
              </div>
              <div className=Styles.addressContainer>
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
            </div>
            <div
              className={Css.merge([
                CssHelper.infoContainer,
                CssHelper.flexBox(~direction=`column, ~justify=`center, ~align=`stretch, ()),
                CssHelper.mbSm(~size=24, ()),
              ])}>
              <Heading size=Heading.H4 value="Total Balance" marginBottom=8 />
              {switch (topPartAllSub) {
               | Data(({financial}, {balance, commission}, {amount, reward}, unbonding, _)) =>
                 totalBalanceRender(
                   sumBalance(balance, amount, unbonding, reward, commission),
                   financial.usdPrice,
                 )
               | _ =>
                 <>
                   <LoadingCensorBar width=200 height=22 mb=10 />
                   <LoadingCensorBar width=220 height=16 />
                 </>
               }}
            </div>
          </div>
        </Col.Grid>
        <Col.Grid col=Col.Six>
          <div className=CssHelper.infoContainer>
            <Heading value="Balance" size=Heading.H4 style=Styles.infoHeader marginBottom=24 />
            <div className=Styles.amountBoxes>
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
               | _ => <LoadingCensorBar fullWidth=true height=12 radius=50 />
               }}
              <div>
                {switch (topPartAllSub) {
                 | Data(({financial}, {balance}, _, _, _)) =>
                   balanceDetail(
                     ~title="Available",
                     ~description="Balance available to send, delegate, etc",
                     ~amount={
                       balance->Coin.getBandAmountFromCoins;
                     },
                     ~usdPrice=financial.usdPrice,
                     ~color=Colors.bandBlue,
                     (),
                   )
                 | _ => <BalanceDetailLoading />
                 }}
              </div>
              <div>
                {switch (topPartAllSub) {
                 | Data(({financial}, _, {amount}, _, _)) =>
                   balanceDetail(
                     ~title="Delegated",
                     ~description="Balance currently delegated to validators",
                     ~amount={
                       amount->Coin.getBandAmountFromCoin;
                     },
                     ~usdPrice=financial.usdPrice,
                     ~color=Colors.chartBalanceAtStake,
                     (),
                   )
                 | _ => <BalanceDetailLoading />
                 }}
              </div>
              <div>
                {switch (topPartAllSub) {
                 | Data(({financial}, _, _, unbonding, _)) =>
                   balanceDetail(
                     ~title="Unbonding",
                     ~description=
                       "Amount undelegated from validators awaiting 21 days lockup period",
                     ~amount={
                       unbonding->Coin.getBandAmountFromCoin;
                     },
                     ~usdPrice=financial.usdPrice,
                     ~color=Colors.blue4,
                     (),
                   )
                 | _ => <BalanceDetailLoading />
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
                 | _ => <BalanceDetailLoading />
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

               | _ => React.null
               }}
            </div>
          </div>
        </Col.Grid>
      </Row>
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
