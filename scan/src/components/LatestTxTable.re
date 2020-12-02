module Styles = {
  open Css;

  let statusImg = style([width(`px(20)), marginTop(`px(-3))]);
};

let renderBody = (reserveIndex, txSub: ApolloHooks.Subscription.variant(TxSub.t)) => {
  <TBody
    key={
      switch (txSub) {
      | Data({txHash}) => txHash |> Hash.toHex
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row alignItems=Row.Start>
      <Col.Grid col=Col.Two>
        {switch (txSub) {
         | Data({txHash}) => <TxLink txHash width=110 />
         | _ => <LoadingCensorBar width=60 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Two>
        {switch (txSub) {
         | Data({blockHeight}) => <TypeID.Block id=blockHeight />
         | _ => <LoadingCensorBar width=50 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.One>
        <div className={CssHelper.flexBox(~justify=`center, ~align=`center, ())}>
          {switch (txSub) {
           | Data({success}) =>
             <img src={success ? Images.success : Images.fail} className=Styles.statusImg />
           | _ => <LoadingCensorBar width=20 height=20 radius=20 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Seven>
        {switch (txSub) {
         | Data({messages, txHash, success, errMsg}) =>
           <TxMessages txHash messages success errMsg width=320 />
         | _ => <LoadingCensorBar width=320 height=15 />
         }}
      </Col.Grid>
    </Row>
  </TBody>;
};

let renderBodyMobile = (reserveIndex, txSub: ApolloHooks.Subscription.variant(TxSub.t)) => {
  switch (txSub) {
  | Data({txHash, blockHeight, success, messages, errMsg}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("Tx Hash", TxHash(txHash, Media.isSmallMobile() ? 170 : 200)),
        ("Block", Height(blockHeight)),
        ("Actions", Messages(txHash, messages, success, errMsg)),
      ]
      key={txHash |> Hash.toHex}
      idx={txHash |> Hash.toHex}
      status=success
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("Tx Hash", Loading(Media.isSmallMobile() ? 170 : 200)),
        ("Block", Loading(70)),
        (
          "Actions",
          Loading(
            {
              Media.isSmallMobile() ? 160 : 230;
            },
          ),
        ),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = () => {
  let isMobile = Media.isMobile();
  let txCount = isMobile ? 5 : 10;
  let txsSub = TxSub.getList(~page=1, ~pageSize=txCount, ());

  <>
    <div
      className={CssHelper.flexBox(~justify=`spaceBetween, ~align=`flexEnd, ())}
      id="latestTransactionsSectionHeader">
      <div>
        <Text
          value="Latest Transactions"
          size=Text.Lg
          block=true
          color=Colors.gray7
          weight=Text.Medium
        />
        <VSpacing size={`px(4)} />
        {switch (txsSub) {
         | ApolloHooks.Subscription.Data(requests) =>
           <Text
             value={
               requests
               ->Belt.Array.get(0)
               ->Belt.Option.mapWithDefault(0, ({id}) => id)
               ->Format.iPretty
             }
             size=Text.Lg
             color=Colors.gray7
             weight=Text.Medium
           />
         | _ => <LoadingCensorBar width=90 height=18 />
         }}
      </div>
      <Link className={CssHelper.flexBox(~align=`center, ())} route=Route.TxHomePage>
        <Text value="All Transactions" color=Colors.bandBlue weight=Text.Medium />
        <HSpacing size=Spacing.md />
        <Icon name="fal fa-angle-right" color=Colors.bandBlue />
      </Link>
    </div>
    <VSpacing size={`px(16)} />
    {isMobile
       ? React.null
       : <THead height=30>
           <Row alignItems=Row.Center>
             <Col.Grid col=Col.Two>
               <div className={CssHelper.flexBox()}>
                 <Text value="TX Hash" size=Text.Sm weight=Text.Semibold color=Colors.gray7 />
               </div>
             </Col.Grid>
             <Col.Grid col=Col.Two>
               <div className={CssHelper.flexBox()}>
                 <Text value="Block" size=Text.Sm weight=Text.Semibold color=Colors.gray7 />
               </div>
             </Col.Grid>
             <Col.Grid col=Col.One>
               <div className={CssHelper.flexBox(~justify=`center, ~align=`center, ())}>
                 <Text value="Status" size=Text.Sm weight=Text.Semibold color=Colors.gray7 />
               </div>
             </Col.Grid>
             <Col.Grid col=Col.Seven>
               <div className={CssHelper.flexBox()}>
                 <Text value="Actions" size=Text.Sm weight=Text.Semibold color=Colors.gray7 />
               </div>
             </Col.Grid>
           </Row>
         </THead>}
    {switch (txsSub) {
     | Data(txs) =>
       txs
       ->Belt_Array.mapWithIndex((i, e) =>
           isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
         )
       ->React.array
     | _ =>
       Belt_Array.make(txCount, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
         )
       ->React.array
     }}
  </>;
};
