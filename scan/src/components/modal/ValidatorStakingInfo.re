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
};

let stakingBalanceDetail = (title, amount, usdPrice) => {
  <Row alignItems=Css.flexStart>
    <Col size=1.2>
      <Text value=title size=Text.Sm height={Text.Px(18)} spacing={Text.Em(0.03)} nowrap=true />
    </Col>
    <Col size=0.6>
      <div className=Styles.cFlex>
        <div className=Styles.rFlex>
          <Text
            value={amount->Coin.getBandAmountFromCoin |> Format.fPretty}
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
            value={amount->Coin.getBandAmountFromCoin *. usdPrice |> Format.fPretty}
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

module StakingInfo = {
  [@react.component]
  let make = (~delegatorAddress, ~validatorAddress) =>
    {
      let (_, dispatchModal) = React.useContext(ModalContext.context);

      let infoSub = React.useContext(GlobalContext.context);
      let balanceAtStakeSub =
        DelegationSub.getStakeByValiator(delegatorAddress, validatorAddress);
      let unbondingSub =
        UnbondingSub.getUnbondingBalanceByValidator(delegatorAddress, validatorAddress);
      let unbondingListSub = UnbondingSub.getUnbondingList(delegatorAddress, validatorAddress);

      let%Sub info = infoSub;
      let%Sub balanceAtStake = balanceAtStakeSub;
      let%Sub unbonding = unbondingSub;
      let%Sub unbondingList = unbondingListSub;

      let balanceAtStakeAmount = balanceAtStake.amount;
      let unbondingAmount = unbonding;
      let rewardAmount = balanceAtStake.reward;
      let usdPrice = info.financial.usdPrice;

      let delegate = () =>
        dispatchModal(OpenModal(SubmitTx(SubmitMsg.Delegate(validatorAddress))));
      let undelegate = () =>
        dispatchModal(OpenModal(SubmitTx(SubmitMsg.Undelegate(validatorAddress))));
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
                 <Text value="Warning: You have reached the maximum number (7) of pending delegation unbonding entries." />
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
          <button className={Styles.button(100)} onClick={_ => {delegate()}}>
            <Text value="Delegate" />
          </button>
          <HSpacing size=Spacing.md />
          <button
            className={Styles.button(100)}
            onClick={_ => {undelegate()}}
            disabled={balanceAtStakeAmount.amount == 0. || isReachUnbondingLimit}>
            <Text value="Undelegate" />
          </button>
          <HSpacing size=Spacing.md />
          <button
            className={Styles.button(150)}
            onClick={_ => {withdrawReward()}}
            disabled={rewardAmount.amount < 1.}>
            <Text value="Withdraw Reward" />
          </button>
        </Row>
        <VSpacing size=Spacing.lg />
        {stakingBalanceDetail("BALANCE AT STAKE", balanceAtStakeAmount, usdPrice)}
        <VSpacing size=Spacing.lg />
        {stakingBalanceDetail("UNBONDING AMOUNT", unbondingAmount, usdPrice)}
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
                   ->Belt_Array.map(({completionTime, balance}) =>
                       [
                         KVTable.Value(balance |> Coin.getBandAmountFromCoin |> Format.fPretty),
                         KVTable.Value(
                           completionTime
                           |> MomentRe.Moment.format("MMM-DD-YYYY  hh:mm:ss A [+UTC]")
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
        {stakingBalanceDetail("REWARD", rewardAmount, usdPrice)}
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
      <img src=Images.connectIcon className=Styles.logo />
    </div>;
  };
};

[@react.component]
let make = (~validatorAddress) => {
  let metadataSub = MetadataSub.use();
  let (accountOpt, _) = React.useContext(AccountContext.context);
  let (_, dispatchModal) = React.useContext(ModalContext.context);

  let connect = chainID => dispatchModal(OpenModal(Connect(chainID)));

  <div className=Styles.topPartWrapper>
    <Text value="DELEGATION INFO" size=Text.Lg weight=Text.Semibold />
    <VSpacing size=Spacing.md />
    {switch (accountOpt) {
     | Some({address}) => <StakingInfo validatorAddress delegatorAddress=address />
     | None =>
       switch (metadataSub) {
       | Data({chainID}) =>
         <div>
           <Row>
             <Col> <Text value="Please connect to see delegation info." /> </Col>
             <Col> <ConnectBtn connect={_ => connect(chainID)} /> </Col>
           </Row>
         </div>
       | Error(err) =>
         // log for err details
         Js.Console.log(err);
         <Text value="chain id not found" />;
       | _ => <LoadingCensorBar width=300 height=18 />
       }
     }}
  </div>;
};
