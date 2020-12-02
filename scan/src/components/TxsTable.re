module Styles = {
  open Css;

  let statusImg = style([width(`px(20)), marginTop(`px(-3))]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

let renderBody =
    (
      reserveIndex,
      txSub: ApolloHooks.Subscription.variant(TxSub.t),
      msgTransform: TxSub.Msg.t => TxSub.Msg.t,
    ) => {
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
         | Data({txHash}) => <TxLink txHash width=140 />
         | _ => <LoadingCensorBar width=140 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.One>
        {switch (txSub) {
         | Data({blockHeight}) => <TypeID.Block id=blockHeight />
         | _ => <LoadingCensorBar width=65 height=15 />
         }}
      </Col.Grid>
      <Col.Grid col=Col.One>
        <div className={CssHelper.flexBox(~justify=`center, ())}>
          {switch (txSub) {
           | Data({success}) =>
             <img src={success ? Images.success : Images.fail} className=Styles.statusImg />
           | _ => <LoadingCensorBar width=20 height=20 radius=20 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Two>
        <div className={CssHelper.flexBox(~justify=`center, ())}>
          {switch (txSub) {
           | Data({gasFee}) =>
             <Text
               block=true
               value={gasFee->Coin.getBandAmountFromCoins->Format.fPretty}
               align=Text.Center
             />
           | _ => <LoadingCensorBar width=65 height=15 />
           }}
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Six>
        {switch (txSub) {
         | Data({messages, txHash, success, errMsg}) =>
           <div>
             <TxMessages
               txHash
               messages={messages->Belt_List.map(msgTransform)}
               success
               errMsg
               width=400
             />
           </div>
         | _ => <> <LoadingCensorBar width=400 height=15 /> </>
         }}
      </Col.Grid>
    </Row>
  </TBody>;
};

let renderBodyMobile =
    (
      reserveIndex,
      txSub: ApolloHooks.Subscription.variant(TxSub.t),
      msgTransform: TxSub.Msg.t => TxSub.Msg.t,
    ) => {
  switch (txSub) {
  | Data({txHash, blockHeight, gasFee, success, messages, errMsg}) =>
    let msgTransform = messages->Belt_List.map(msgTransform);
    <MobileCard
      values=InfoMobileCard.[
        ("TX Hash", TxHash(txHash, Media.isSmallMobile() ? 170 : 200)),
        ("Block", Height(blockHeight)),
        ("Gas Fee\n(BAND)", Coin({value: gasFee, hasDenom: false})),
        ("Actions", Messages(txHash, msgTransform, success, errMsg)),
      ]
      key={txHash |> Hash.toHex}
      idx={txHash |> Hash.toHex}
      status=success
    />;
  | _ =>
    <MobileCard
      values=InfoMobileCard.[
        ("TX Hash", Loading(Media.isSmallMobile() ? 170 : 200)),
        ("Block", Loading(70)),
        ("Gas Fee\n(BAND)", Loading(70)),
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
let make =
    (
      ~txsSub: ApolloHooks.Subscription.variant(array(TxSub.t)),
      ~msgTransform: TxSub.Msg.t => TxSub.Msg.t=x => x,
    ) => {
  let isMobile = Media.isMobile();
  <>
    {switch (txsSub) {
     | Data(txs) =>
       txs->Belt.Array.size > 0
         ? txs
           ->Belt_Array.mapWithIndex((i, e) =>
               isMobile
                 ? renderBodyMobile(i, Sub.resolve(e), msgTransform)
                 : renderBody(i, Sub.resolve(e), msgTransform)
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
       Belt_Array.make(10, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) =>
           isMobile
             ? renderBodyMobile(i, noData, msgTransform) : renderBody(i, noData, msgTransform)
         )
       ->React.array
     }}
  </>;
};
