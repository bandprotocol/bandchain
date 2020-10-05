module Styles = {
  open Css;

  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

let renderBody = (reserveIndex, txSub: ApolloHooks.Subscription.variant(TxSub.t)) => {
  <TBody.Grid
    key={
      switch (txSub) {
      | Data({txHash}) => txHash |> Hash.toHex
      | _ => reserveIndex |> string_of_int
      }
    }
    paddingH={`px(24)}>
    <Row.Grid>
      <Col.Grid col=Col.Two>
        {switch (txSub) {
         | Data({txHash}) => <TxLink txHash width=140 />
         | _ => <LoadingCensorBar width=170 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Two>
        {switch (txSub) {
         | Data({gasFee}) =>
           <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
             <Text
               block=true
               code=true
               spacing={Text.Em(0.02)}
               value={gasFee->Coin.getBandAmountFromCoins->Format.fPretty}
               weight=Text.Medium
             />
           </div>
         | _ => <LoadingCensorBar width=30 height=15 isRight=true />
         }}
      </Col.Grid>
      <Col.Grid col=Col.Eight>
        {switch (txSub) {
         | Data({messages, txHash, success, errMsg}) =>
           <TxMessages txHash messages success errMsg width=530 />

         | _ => <LoadingCensorBar width=530 height=15 />
         }}
      </Col.Grid>
    </Row.Grid>
  </TBody.Grid>;
};

let renderBodyMobile = (reserveIndex, txSub: ApolloHooks.Subscription.variant(TxSub.t)) => {
  switch (txSub) {
  | Data({txHash, gasFee, success, messages, errMsg}) =>
    <MobileCard
      values=InfoMobileCard.[
        ("TX Hash", TxHash(txHash, 200)),
        ("Gas Fee\n(BAND)", Coin({value: gasFee, hasDenom: false})),
        ("Actions", Messages(txHash, messages, success, errMsg)),
      ]
      key={txHash |> Hash.toHex}
      idx={txHash |> Hash.toHex}
      status=success
    />
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("TX Hash", Loading(200)),
        ("Gas Fee\n(BAND)", Loading(60)),
        ("Actions", Loading(230)),
      ]
      key={reserveIndex |> string_of_int}
      idx={reserveIndex |> string_of_int}
    />
  };
};

[@react.component]
let make = (~txsSub: ApolloHooks.Subscription.variant(array(TxSub.t))) => {
  let isMobile = Media.isMobile();
  <>
    {isMobile
       ? React.null
       : <THead.Grid>
           <Row.Grid alignItems=Row.Center>
             <Col.Grid col=Col.Two>
               <Text
                 block=true
                 value="TX Hash"
                 size=Text.Md
                 weight=Text.Semibold
                 color=Colors.gray7
               />
             </Col.Grid>
             <Col.Grid col=Col.Two>
               <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
                 <Text
                   block=true
                   value="Gas Fee (BAND)"
                   size=Text.Md
                   weight=Text.Semibold
                   color=Colors.gray7
                 />
               </div>
             </Col.Grid>
             <Col.Grid col=Col.Eight>
               <Text
                 block=true
                 value="Actions"
                 size=Text.Md
                 weight=Text.Semibold
                 color=Colors.gray7
               />
             </Col.Grid>
           </Row.Grid>
         </THead.Grid>}
    {switch (txsSub) {
     | Data(txs) =>
       txs->Belt.Array.size > 0
         ? txs
           ->Belt_Array.mapWithIndex((i, e) =>
               isMobile ? renderBodyMobile(i, Sub.resolve(e)) : renderBody(i, Sub.resolve(e))
             )
           ->React.array
         : <EmptyContainer>
             <img src=Images.noBlock className=Styles.noDataImage />
             <Heading
               size=Heading.H4
               value="No Transaction"
               align=Heading.Center
               weight=Heading.Regular
               color=Colors.bandBlue
             />
           </EmptyContainer>
     | _ =>
       Belt_Array.make(isMobile ? 1 : 10, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile ? renderBodyMobile(i, noData) : renderBody(i, noData)
         )
       ->React.array
     }}
  </>;
};
