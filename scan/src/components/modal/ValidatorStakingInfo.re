module Styles = {
  open Css;
  let topPartWrapper =
    style([
      width(`percent(100.0)),
      display(`flex),
      flexDirection(`column),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      padding2(~v=`px(35), ~h=`px(30)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(8), Css.rgba(0, 0, 0, 0.08))),
    ]);

  let cFlex = style([display(`flex), flexDirection(`column), alignItems(`flexEnd)]);
  let rFlex = style([display(`flex), flexDirection(`row)]);
  let balance = style([minWidth(`px(150)), justifyContent(`flexEnd)]);
  let button = wid =>
    style([
      display(`flex),
      backgroundColor(Colors.purple1),
      width(`px(wid)),
      height(`px(25)),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(4)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(11, 29, 142, 0.1))),
      border(`zero, `solid, Colors.blueGray1),
      color(Colors.purple7),
      cursor(`pointer),
      disabled([backgroundColor(Colors.gray3), color(Colors.gray6), cursor(`default)]),
      focus([outline(`zero, `none, Colors.white)]),
    ]);
  let logo = style([width(`px(10))]);

  let connectBtn =
    style([
      backgroundColor(Colors.green1),
      padding2(~h=`px(8), ~v=`px(2)),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(10)),
      cursor(`pointer),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(17, 85, 78, 0.1))),
    ]);

  let reminder =
    style([
      padding(`px(10)),
      color(Colors.blue5),
      backgroundColor(Colors.blue1),
      border(`px(1), `solid, Colors.blue6),
      borderRadius(`px(4)),
    ]);

  let warning =
    style([
      padding(`px(10)),
      color(Colors.yellow5),
      backgroundColor(Colors.yellow1),
      border(`px(1), `solid, Colors.yellow6),
      borderRadius(`px(4)),
    ]);

  let connectContainer =
    style([
      height(`px(200)),
      display(`flex),
      flexDirection(`column),
      justifyContent(`center),
      alignItems(`center),
      backgroundColor(Colors.profileBG),
    ]);

  let infoContainer =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      padding(`px(24)),
      Media.mobile([padding(`px(16))]),
    ]);
  let infoHeader =
    style([
      borderBottom(`px(1), `solid, Colors.gray9),
      paddingBottom(`px(16)),
      marginBottom(`px(16)),
    ]);
  let loadingBox = style([width(`percent(100.))]);
  let rewardContainer =
    style([backgroundColor(Colors.profileBG), padding2(~v=`px(16), ~h=`px(24))]);
};

let stakingBalanceDetail = (~title, ~amount, ~usdPrice, ~tooltipItem, ~isCountup=false, ()) => {
  <Row alignItems=Css.flexStart>
    <Col size=1.2>
      <Text
        value=title
        size=Text.Sm
        height={Text.Px(18)}
        spacing={Text.Em(0.03)}
        nowrap=true
        tooltipItem={tooltipItem |> React.string}
        tooltipPlacement=Text.AlignRight
      />
    </Col>
    <Col size=0.6>
      <div className=Styles.cFlex>
        <div className=Styles.rFlex>
          {isCountup
             ? <NumberCountup
                 value={amount->Coin.getBandAmountFromCoin}
                 size=Text.Lg
                 weight=Text.Semibold
                 spacing={Text.Em(0.02)}
               />
             : <Text
                 value={amount->Coin.getBandAmountFromCoin |> Format.fPretty}
                 size=Text.Lg
                 weight=Text.Semibold
                 spacing={Text.Em(0.02)}
                 nowrap=true
                 code=true
               />}
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
          {isCountup
             ? <NumberCountup
                 value={amount->Coin.getBandAmountFromCoin *. usdPrice}
                 size=Text.Sm
                 weight=Text.Thin
                 spacing={Text.Em(0.02)}
               />
             : <Text
                 value={amount->Coin.getBandAmountFromCoin *. usdPrice |> Format.fPretty}
                 size=Text.Sm
                 spacing={Text.Em(0.02)}
                 weight=Text.Thin
                 nowrap=true
                 code=true
               />}
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

module ButtonSection = {
  [@react.component]
  let make = (~validatorAddress) => {
    let (_, dispatchModal) = React.useContext(ModalContext.context);
    let validatorInfoSub = ValidatorSub.get(validatorAddress);

    let delegate = () =>
      dispatchModal(OpenModal(SubmitTx(SubmitMsg.Delegate(validatorAddress))));
    let undelegate = () =>
      dispatchModal(OpenModal(SubmitTx(SubmitMsg.Undelegate(validatorAddress))));
    let redelegate = () =>
      dispatchModal(OpenModal(SubmitTx(SubmitMsg.Redelegate(validatorAddress))));

    switch (validatorInfoSub) {
    | Data(validatorInfo) =>
      <div className={CssHelper.flexBox()}>
        <button
          className={CssHelper.btn(~px=20, ~py=5, ())}
          onClick={_ => {
            validatorInfo.commission == 100.
              ? Window.alert("Delegation to foundation validator nodes is not advised.")
              : delegate()
          }}>
          <Text value="Delegate" weight=Text.Medium nowrap=true block=true />
        </button>
        <HSpacing size=Spacing.md />
        <button
          className={CssHelper.btn(~variant=Outline, ~px=20, ~py=5, ())}
          onClick={_ => {undelegate()}}>
          <Text value="Undelegate" weight=Text.Medium nowrap=true block=true />
        </button>
        <HSpacing size=Spacing.md />
        <button
          className={CssHelper.btn(~variant=Outline, ~px=20, ~py=5, ())}
          onClick={_ => {redelegate()}}>
          <Text value="Redelegate" weight=Text.Medium nowrap=true block=true />
        </button>
      </div>
    | _ => React.null
    };
  };
};

module StakingInfo = {
  [@react.component]
  let make = (~delegatorAddress, ~validatorAddress) => {
    let currentTime =
      React.useContext(TimeContext.context) |> MomentRe.Moment.format(Config.timestampUseFormat);
    let (_, dispatchModal) = React.useContext(ModalContext.context);

    let infoSub = React.useContext(GlobalContext.context);
    let validatorInfoSub = ValidatorSub.get(validatorAddress);
    let balanceAtStakeSub = DelegationSub.getStakeByValiator(delegatorAddress, validatorAddress);
    let unbondingSub =
      UnbondingSub.getUnbondingBalanceByValidator(
        delegatorAddress,
        validatorAddress,
        currentTime,
      );

    let withdrawReward = () =>
      dispatchModal(OpenModal(SubmitTx(SubmitMsg.WithdrawReward(validatorAddress))));

    <>
      <Row.Grid marginBottom=24>
        <Col.Grid>
          <Text
            value="Note: You have non-zero pending reward on this validator. Any additional staking actions will automatically withdraw that reward your balance."
            color=Colors.gray6
            weight=Text.Thin
          />
        </Col.Grid>
      </Row.Grid>
      <Row.Grid marginBottom=24>
        <Col.Grid col=Col.Six>
          <div>
            <Heading value="Balance at stake" size=Heading.H5 />
            <VSpacing size={`px(8)} />
            <Text value="0.500000 BAND" size=Text.Lg color=Colors.gray7 weight=Text.Thin />
          </div>
        </Col.Grid>
        <Col.Grid col=Col.Six>
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
            <Text value="0.500000 BAND" size=Text.Lg color=Colors.gray7 weight=Text.Thin />
          </div>
        </Col.Grid>
      </Row.Grid>
      <Row.Grid style=Styles.rewardContainer alignItems=Row.Center>
        <Col.Grid>
          <div className={CssHelper.flexBox(~justify=`spaceBetween, ())}>
            <div>
              <Heading value="Reward" size=Heading.H5 />
              <VSpacing size={`px(8)} />
              <Text value="0.500000 BAND" size=Text.Lg color=Colors.gray7 weight=Text.Thin />
            </div>
            <button className={CssHelper.btn(~px=20, ~py=5, ())} onClick={_ => withdrawReward()}>
              <Text value="Withdraw Reward" weight=Text.Medium nowrap=true block=true />
            </button>
          </div>
        </Col.Grid>
      </Row.Grid>
    </>;
  };
};

module StakingInfo23 = {
  [@react.component]
  let make = (~delegatorAddress, ~validatorAddress) =>
    {
      let currentTime =
        React.useContext(TimeContext.context)
        |> MomentRe.Moment.format(Config.timestampUseFormat);
      let (_, dispatchModal) = React.useContext(ModalContext.context);

      let infoSub = React.useContext(GlobalContext.context);
      let validatorInfoSub = ValidatorSub.get(validatorAddress);
      let balanceAtStakeSub =
        DelegationSub.getStakeByValiator(delegatorAddress, validatorAddress);
      let unbondingSub =
        UnbondingSub.getUnbondingBalanceByValidator(
          delegatorAddress,
          validatorAddress,
          currentTime,
        );
      let unbondingListSub =
        UnbondingSub.getUnbondingList(delegatorAddress, validatorAddress, currentTime);

      let%Sub info = infoSub;
      let%Sub balanceAtStake = balanceAtStakeSub;
      let%Sub unbonding = unbondingSub;
      let%Sub unbondingList = unbondingListSub;
      let%Sub validatorInfo = validatorInfoSub;

      let balanceAtStakeAmount = balanceAtStake.amount;
      let unbondingAmount = unbonding;
      let rewardAmount = balanceAtStake.reward;
      let usdPrice = info.financial.usdPrice;

      let delegate = () =>
        dispatchModal(OpenModal(SubmitTx(SubmitMsg.Delegate(validatorAddress))));
      let undelegate = () =>
        dispatchModal(OpenModal(SubmitTx(SubmitMsg.Undelegate(validatorAddress))));
      let redelegate = () =>
        dispatchModal(OpenModal(SubmitTx(SubmitMsg.Redelegate(validatorAddress))));
      let withdrawReward = () =>
        dispatchModal(OpenModal(SubmitTx(SubmitMsg.WithdrawReward(validatorAddress))));
      let isReachUnbondingLimit = unbondingList |> Belt_Array.length == 7;

      <div>
        <VSpacing size=Spacing.md />
        {rewardAmount.amount > 1.
           ? <div>
               <div className=Styles.reminder>
                 <Text
                   value="Note: You have non-zero pending reward on this validator. Any additional staking actions will automatically withdraw that reward your balance."
                 />
               </div>
               <VSpacing size=Spacing.lg />
             </div>
           : React.null}
        {isReachUnbondingLimit
           ? <div>
               <div className=Styles.warning>
                 <Text
                   value="Warning: You have reached the maximum number (7) of pending delegation unbonding entries."
                 />
               </div>
               <VSpacing size=Spacing.lg />
             </div>
           : React.null}
        <Row>
          <Col size=1.2>
            <Text
              value="ACTIONS:"
              size=Text.Sm
              height={Text.Px(18)}
              spacing={Text.Em(0.03)}
              nowrap=true
            />
          </Col>
          <HSpacing size=Spacing.md />
          <button
            onClick={_ => {
              validatorInfo.commission == 100.
                ? Window.alert("Delegation to foundation validator nodes is not advised.")
                : delegate()
            }}>
            <Text value="Delegate" />
          </button>
          <HSpacing size=Spacing.md />
          <button
            onClick={_ => {undelegate()}}
            disabled={balanceAtStakeAmount.amount == 0. || isReachUnbondingLimit}>
            <Text value="Undelegate" />
          </button>
          <HSpacing size=Spacing.md />
          <button onClick={_ => {redelegate()}} disabled={balanceAtStakeAmount.amount == 0.}>
            <Text value="Redelegate" />
          </button>
          <HSpacing size=Spacing.md />
          <button onClick={_ => {withdrawReward()}} disabled={rewardAmount.amount < 1.}>
            <Text value="Withdraw Reward" />
          </button>
        </Row>
        <VSpacing size=Spacing.lg />
        {stakingBalanceDetail(
           ~title="BALANCE AT STAKE",
           ~amount=balanceAtStakeAmount,
           ~usdPrice,
           ~tooltipItem="Balance currently delegated to validators",
           (),
         )}
        <VSpacing size=Spacing.lg />
        {stakingBalanceDetail(
           ~title="UNBONDING AMOUNT",
           ~amount=unbondingAmount,
           ~usdPrice,
           ~tooltipItem="Amount undelegated from validators awaiting 21 days lockup period",
           (),
         )}
        {unbondingList |> Belt_Array.length > 0
           ? <>
               <VSpacing size=Spacing.md />
               <Col size=1.2>
                 <Text
                   value="Breakdown:"
                   size=Text.Sm
                   height={Text.Px(18)}
                   spacing={Text.Em(0.03)}
                   nowrap=true
                 />
               </Col>
               <VSpacing size=Spacing.md />
               <KVTable
                 tableWidth=470
                 headers=["AMOUNT (BAND)", "UNBONDED AT"]
                 rows={
                   unbondingList
                   ->Belt_Array.map(({completionTime, amount}) =>
                       [
                         KVTable.Value(amount |> Coin.getBandAmountFromCoin |> Format.fPretty),
                         KVTable.Value(
                           completionTime
                           |> MomentRe.Moment.format(Config.timestampDisplayFormat)
                           |> String.uppercase_ascii,
                         ),
                       ]
                     )
                   ->Belt_List.fromArray
                 }
               />
             </>
           : React.null}
        <VSpacing size=Spacing.lg />
        {stakingBalanceDetail(
           ~title="REWARD",
           ~amount=rewardAmount,
           ~usdPrice,
           ~tooltipItem="Reward from staking to validators",
           ~isCountup=true,
           (),
         )}
      </div>
      |> Sub.resolve;
    }
    |> Sub.default(_, React.null);
};

module ConnectBtn = {
  [@react.component]
  let make = (~connect) => {
    <div className=Styles.connectBtn onClick={_ => connect()}>
      <Text
        value="connect"
        size=Text.Xs
        weight=Text.Medium
        color=Colors.green7
        nowrap=true
        height={Text.Px(10)}
        spacing={Text.Em(0.03)}
        block=true
      />
      <HSpacing size=Spacing.sm />
      // TODO: change it later
      <Icon name="fal fa-link" color=Colors.gray1 />
    </div>;
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
        //TODO: remove mock message later
        <CTooltip tooltipText="Lorem ipsum, or lipsum as it is sometimes known.">
          <Icon name="fal fa-info-circle" size=10 />
        </CTooltip>
      </div>
      {switch (accountOpt) {
       | Some(_) => <ButtonSection validatorAddress />
       | None => React.null
       }}
    </div>
    {switch (accountOpt) {
     | Some({address}) => <StakingInfo validatorAddress delegatorAddress=address />
     | None =>
       switch (trackingSub) {
       | Data({chainID}) =>
         <div className=Styles.connectContainer>
           <Icon name="fal fa-link" size=32 color=Colors.bandBlue />
           <VSpacing size={`px(16)} />
           <Text value="Please connect to make request" size=Text.Lg nowrap=true block=true />
           <VSpacing size={`px(16)} />
           <button className={CssHelper.btn(~px=20, ~py=5, ())} onClick={_ => connect(chainID)}>
             <Text value="Connect" weight=Text.Medium nowrap=true block=true />
           </button>
         </div>
       | Error(err) =>
         // log for err details
         Js.Console.log(err);
         <Text value="chain id not found" />;
       | _ => <LoadingCensorBar width=100 height=200 style=Styles.loadingBox />
       }
     }}
  </div>;
};
