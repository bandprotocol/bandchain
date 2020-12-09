module Styles = {
  open Css;

  let connectContainer = style([height(`px(200)), backgroundColor(Colors.profileBG)]);

  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      padding(`px(24)),
      Media.mobile([padding(`px(16))]),
    ]);
  let infoHeader =
    style([
      borderBottom(`px(1), `solid, Colors.gray9),
      paddingBottom(`px(12)),
      marginBottom(`px(16)),
      minHeight(`px(41)),
    ]);
  let rewardContainer =
    style([backgroundColor(Colors.profileBG), padding2(~v=`px(16), ~h=`px(24))]);
};

module ButtonSection = {
  [@react.component]
  let make = (~validatorAddress) => {
    let (_, dispatchModal) = React.useContext(ModalContext.context);
    let validatorInfoSub = ValidatorSub.get(validatorAddress);

    let delegate = () => validatorAddress->SubmitMsg.Delegate->SubmitTx->OpenModal->dispatchModal;
    let undelegate = () =>
      validatorAddress->SubmitMsg.Undelegate->SubmitTx->OpenModal->dispatchModal;
    let redelegate = () =>
      validatorAddress->SubmitMsg.Redelegate->SubmitTx->OpenModal->dispatchModal;

    switch (validatorInfoSub) {
    | Data(validatorInfo) =>
      <div className={CssHelper.flexBox()} id="validatorDelegationinfoDlegate">
        <Button
          px=20
          py=5
          onClick={_ => {
            validatorInfo.commission == 100.
              ? Webapi.Dom.(
                  window
                  |> Window.alert("Delegation to foundation validator nodes is not advised.")
                )
              : delegate()
          }}>
          <Text value="Delegate" weight=Text.Medium nowrap=true block=true />
        </Button>
        <HSpacing size=Spacing.md />
        <Button px=20 py=5 variant=Button.Outline onClick={_ => {undelegate()}}>
          <Text value="Undelegate" weight=Text.Medium nowrap=true block=true />
        </Button>
        <HSpacing size=Spacing.md />
        <Button px=20 py=5 variant=Button.Outline onClick={_ => {redelegate()}}>
          <Text value="Redelegate" weight=Text.Medium nowrap=true block=true />
        </Button>
      </div>
    | _ => React.null
    };
  };
};

module DisplayBalance = {
  module Loading = {
    [@react.component]
    let make = () => {
      <>
        <LoadingCensorBar width=120 height=15 />
        <VSpacing size=Spacing.xs />
        <LoadingCensorBar width=80 height=15 />
      </>;
    };
  };

  [@react.component]
  let make = (~amount, ~usdPrice, ~isCountup=false) => {
    <>
      <div className={CssHelper.flexBox()}>
        {isCountup
           ? <NumberCountup
               value={amount->Coin.getBandAmountFromCoin}
               size=Text.Lg
               weight=Text.Regular
               spacing={Text.Em(0.)}
               code=false
             />
           : <Text
               value={amount->Coin.getBandAmountFromCoin |> Format.fPretty(~digits=6)}
               size=Text.Lg
               color=Colors.gray7
               block=true
             />}
        <HSpacing size=Spacing.sm />
        <Text value="BAND" size=Text.Lg color=Colors.gray7 block=true />
      </div>
      <div className={CssHelper.flexBox()}>
        {isCountup
           ? <NumberCountup
               value={amount->Coin.getBandAmountFromCoin *. usdPrice}
               size=Text.Md
               weight=Text.Regular
               color=Colors.gray6
               code=false
               spacing={Text.Em(0.)}
             />
           : <Text
               value={amount->Coin.getBandAmountFromCoin *. usdPrice |> Format.fPretty(~digits=6)}
               size=Text.Md
               color=Colors.gray6
               block=true
             />}
        <HSpacing size=Spacing.sm />
        <Text value="USD" size=Text.Md color=Colors.gray6 block=true />
      </div>
    </>;
  };
};

module StakingInfo = {
  [@react.component]
  let make = (~delegatorAddress, ~validatorAddress) => {
    let currentTime =
      React.useContext(TimeContext.context) |> MomentRe.Moment.format(Config.timestampUseFormat);
    let (_, dispatchModal) = React.useContext(ModalContext.context);

    let infoSub = React.useContext(GlobalContext.context);
    let balanceAtStakeSub = DelegationSub.getStakeByValidator(delegatorAddress, validatorAddress);
    let unbondingSub =
      UnbondingSub.getUnbondingBalanceByValidator(
        delegatorAddress,
        validatorAddress,
        currentTime,
      );

    let allSub = Sub.all3(infoSub, balanceAtStakeSub, unbondingSub);

    let withdrawReward = () => {
      validatorAddress->SubmitMsg.WithdrawReward->SubmitTx->OpenModal->dispatchModal;
    };

    let reinvest = reward =>
      (validatorAddress, reward)->SubmitMsg.Reinvest->SubmitTx->OpenModal->dispatchModal;
    <>
      <Row marginBottom=24>
        <Col>
          <Text
            value="Note: You have non-zero pending reward on this validator. Any additional staking actions will automatically withdraw that reward your balance."
            color=Colors.gray6
            weight=Text.Thin
          />
        </Col>
      </Row>
      <Row marginBottom=24>
        <Col col=Col.Six>
          <div>
            <Heading value="Balance at Stake" size=Heading.H5 />
            <VSpacing size={`px(8)} />
            {switch (allSub) {
             | Data(({financial: {usdPrice}}, balanceAtStake, _)) =>
               <DisplayBalance amount={balanceAtStake.amount} usdPrice />
             | _ => <DisplayBalance.Loading />
             }}
          </div>
        </Col>
        <Col col=Col.Six>
          <div>
            <div className={CssHelper.flexBox()}>
              <Heading value="Unbonding Amount" size=Heading.H5 />
              <HSpacing size=Spacing.sm />
              <Link
                className={CssHelper.flexBox()}
                route={Route.AccountIndexPage(delegatorAddress, Route.AccountUnbonding)}>
                <Text value="View Entries" color=Colors.bandBlue weight=Text.Medium />
              </Link>
            </div>
            <VSpacing size={`px(8)} />
            {switch (allSub) {
             | Data(({financial: {usdPrice}}, _, unbonding)) =>
               <DisplayBalance amount=unbonding usdPrice />
             | _ => <DisplayBalance.Loading />
             }}
          </div>
        </Col>
      </Row>
      <Row style=Styles.rewardContainer alignItems=Row.Center>
        <Col>
          <div className={CssHelper.flexBox(~justify=`spaceBetween, ())}>
            <div>
              <Heading value="Reward" size=Heading.H5 />
              <VSpacing size={`px(8)} />
              {switch (allSub) {
               | Data(({financial: {usdPrice}}, balanceAtStake, _)) =>
                 <DisplayBalance amount={balanceAtStake.reward} usdPrice isCountup=true />
               | _ => <DisplayBalance.Loading />
               }}
            </div>
            <div className={CssHelper.flexBox()} id="withdrawRewardContainer">
              {let (disable, reward) =
                 switch (allSub) {
                 | Data((_, balanceAtStake, _)) => (
                     balanceAtStake.reward.amount <= 0.,
                     balanceAtStake.reward.amount,
                   )
                 | _ => (true, 0.)
                 };

               <>
                 <Button px=20 py=5 onClick={_ => withdrawReward()} disabled=disable>
                   <Text value="Withdraw Reward" weight=Text.Medium nowrap=true block=true />
                 </Button>
                 <HSpacing size=Spacing.sm />
                 <Button px=20 py=5 onClick={_ => reinvest(reward)} disabled=disable>
                   <Text value="Reinvest" weight=Text.Medium nowrap=true block=true />
                 </Button>
               </>}
            </div>
          </div>
        </Col>
      </Row>
    </>;
  };
};

[@react.component]
let make = (~validatorAddress) => {
  let trackingSub = TrackingSub.use();
  let (accountOpt, _) = React.useContext(AccountContext.context);
  let (_, dispatchModal) = React.useContext(ModalContext.context);

  let connect = chainID => dispatchModal(OpenModal(Connect(chainID)));

  <div className=Styles.infoContainer>
    <div
      className={Css.merge([CssHelper.flexBox(~justify=`spaceBetween, ()), Styles.infoHeader])}>
      <div className={CssHelper.flexBox()}>
        <Heading value="Your Delegation Info" size=Heading.H4 />
        <HSpacing size=Spacing.xs />
        <CTooltip tooltipText="Your delegation stats on this validators">
          <Icon name="fal fa-info-circle" size=10 />
        </CTooltip>
      </div>
      {switch (accountOpt) {
       | Some(_) => <ButtonSection validatorAddress />
       | None => <VSpacing size={`px(28)} />
       }}
    </div>
    {switch (accountOpt) {
     | Some({address}) => <StakingInfo validatorAddress delegatorAddress=address />
     | None =>
       switch (trackingSub) {
       | Data({chainID}) =>
         <div
           className={Css.merge([
             CssHelper.flexBox(~direction=`column, ~justify=`center, ()),
             Styles.connectContainer,
           ])}>
           <Icon name="fal fa-link" size=32 color=Colors.bandBlue />
           <VSpacing size={`px(16)} />
           <Text value="Please connect to make request" size=Text.Lg nowrap=true block=true />
           <VSpacing size={`px(16)} />
           <Button px=20 py=5 onClick={_ => connect(chainID)}>
             <Text value="Connect" weight=Text.Medium nowrap=true block=true />
           </Button>
         </div>
       | Error(err) =>
         // log for err details
         Js.Console.log(err);
         <Text value="chain id not found" />;
       | _ => <LoadingCensorBar fullWidth=true height=200 />
       }
     }}
  </div>;
};
