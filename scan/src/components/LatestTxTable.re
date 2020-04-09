module Styles = {
  open Css;

  let topicBar =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`row),
      justifyContent(`spaceBetween),
    ]);
  let seeAll = style([display(`flex), flexDirection(`row), cursor(`pointer)]);
  let cFlex = style([display(`flex), flexDirection(`column)]);
  let rightArrow = style([width(`px(25)), marginTop(`px(17)), marginLeft(`px(16))]);

  let hScale = 20;
  let fullWidth = style([width(`percent(100.0)), display(`flex)]);
  let blockContainer = style([minWidth(`px(60))]);
  let statusContainer =
    style([
      minWidth(`percent(100.)),
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      justifyContent(`center),
    ]);

  let logo = style([width(`px(20))]);

  let heightByMsgsNum = (numMsgs, mt) =>
    style([
      minHeight(numMsgs <= 1 ? `auto : `px(numMsgs * hScale)),
      marginTop(`px(numMsgs <= 1 ? 0 : mt)),
    ]);
};

let txBodyRender = (reserveIndex: int, txSub: ApolloHooks.Subscription.variant(TxSub.t)) => {
  <TBody
    key={
      switch (txSub) {
      | Data({txHash}) => txHash |> Hash.toHex
      | _ => reserveIndex |> string_of_int
      }
    }>
    <Row minHeight={`px(30)}>
      <Col> <HSpacing size={`px(12)} /> </Col>
      <Col size=1.2>
        {switch (txSub) {
         | Data({messages, txHash}) =>
           <div className={Styles.heightByMsgsNum(messages->Belt_List.size, 0)}>
             <TxLink txHash width=110 />
           </div>
         | _ => <LoadingCensorBar width=100 height=10 />
         }}
      </Col>
      <Col size=0.68>
        {switch (txSub) {
         | Data({messages, blockHeight}) =>
           <div
             className={Css.merge([
               Styles.heightByMsgsNum(messages->Belt_List.size, -4),
               Styles.blockContainer,
             ])}>
             <TypeID.Block id=blockHeight />
           </div>
         | _ => <LoadingCensorBar width=50 height=10 />
         }}
      </Col>
      <Col size=1.>
        {switch (txSub) {
         | Data({messages, success}) =>
           <div className={Styles.heightByMsgsNum(messages->Belt_List.size, -8)}>
             <div className=Styles.statusContainer>
               <img src={success ? Images.success : Images.fail} className=Styles.logo />
             </div>
           </div>
         | _ =>
           <div className=Styles.statusContainer>
             <LoadingCensorBar width=20 height=20 radius=20 />
           </div>
         }}
      </Col>
      <Col size=3.8>
        {switch (txSub) {
         | Data({messages, txHash}) =>
           messages
           ->Belt_List.toArray
           ->Belt_Array.mapWithIndex((i, msg) =>
               <React.Fragment key={(txHash |> Hash.toHex) ++ (i |> string_of_int)}>
                 <VSpacing size=Spacing.sm />
                 <Msg msg width=350 />
                 <VSpacing size=Spacing.sm />
               </React.Fragment>
             )
           ->React.array
         | _ => <LoadingCensorBar width=360 height=10 />
         }}
      </Col>
      <Col> <HSpacing size={`px(12)} /> </Col>
    </Row>
  </TBody>;
};

[@react.component]
let make = () => {
  let allSub = Sub.all2(TxSub.getList(~page=1, ~pageSize=10, ()), TxSub.count());

  <>
    <div className=Styles.topicBar>
      <Text
        value="Latest Transactions"
        size=Text.Xxl
        weight=Text.Bold
        block=true
        color=Colors.gray8
      />
      <div className=Styles.seeAll onClick={_ => Route.redirect(Route.TxHomePage)}>
        <div className=Styles.cFlex>
          {switch (allSub) {
           | Data((_, totalCount)) =>
             <Text
               value={totalCount |> Format.iPretty}
               size=Text.Xxl
               color=Colors.gray8
               height={Text.Px(24)}
               weight=Text.Bold
             />
           | _ => <LoadingCensorBar width=90 height=18 />
           }}
          <VSpacing size=Spacing.xs />
          <Text
            value="ALL TRANSACTIONS"
            size=Text.Sm
            color=Colors.bandBlue
            spacing={Text.Em(0.05)}
            weight=Text.Medium
          />
        </div>
        <img src=Images.rightArrow className=Styles.rightArrow />
      </div>
    </div>
    <VSpacing size=Spacing.lg />
    <THead>
      <Row>
        <Col> <HSpacing size={`px(12)} /> </Col>
        <Col size=1.2>
          <div className=Styles.fullWidth>
            <Text
              value="TX HASH"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col size=0.68>
          <div className={Css.merge([Styles.fullWidth, Styles.blockContainer])}>
            <Text
              value="BLOCK"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col size=1.>
          <div className=Styles.statusContainer>
            <Text
              value="STATUS"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col size=3.8>
          <div className=Styles.fullWidth>
            <Text
              value="ACTIONS"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col> <HSpacing size={`px(12)} /> </Col>
      </Row>
    </THead>
    {switch (allSub) {
     | Data((txs, _)) =>
       txs->Belt_Array.mapWithIndex((i, e) => txBodyRender(i, Sub.resolve(e)))->React.array
     | _ =>
       Belt_Array.make(10, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) => txBodyRender(i, noData))
       ->React.array
     }}
  </>;
};
