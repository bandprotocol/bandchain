module Styles = {
  open Css;

  let tableWrapper = style([Media.mobile([padding2(~v=`px(16), ~h=`zero)])]);
  let noDataImage = style([width(`auto), height(`px(70)), marginBottom(`px(16))]);
};

let transform = (account, msg: TxSub.Msg.t) => {
  switch (msg) {
  | SendMsg({toAddress, fromAddress, amount}) when toAddress == account =>
    TxSub.Msg.ReceiveMsg({toAddress, fromAddress, amount})
  | _ => msg
  };
};

[@react.component]
let make = (~accountAddress: Address.t) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let txsSub = TxSub.getListBySender(accountAddress, ~pageSize, ~page, ());
  let txsCountSub = TxSub.countBySender(accountAddress);

  let isMobile = Media.isMobile();

  <div className=Styles.tableWrapper>
    {isMobile
       ? <Row.Grid marginBottom=16>
           <Col.Grid>
             {switch (txsCountSub) {
              | Data(txsCount) =>
                <div className={CssHelper.flexBox()}>
                  <Text
                    block=true
                    value={txsCount |> string_of_int}
                    weight=Text.Semibold
                    color=Colors.gray7
                  />
                  <HSpacing size=Spacing.xs />
                  <Text block=true value="Transactions" weight=Text.Semibold color=Colors.gray7 />
                </div>
              | _ => <LoadingCensorBar width=100 height=15 />
              }}
           </Col.Grid>
         </Row.Grid>
       : <THead.Grid>
           <Row.Grid alignItems=Row.Center>
             <Col.Grid col=Col.Two>
               {switch (txsCountSub) {
                | Data(txsCount) =>
                  <div className={CssHelper.flexBox()}>
                    <Text
                      block=true
                      value={txsCount |> string_of_int}
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                    <HSpacing size=Spacing.xs />
                    <Text
                      block=true
                      value="Transactions"
                      weight=Text.Semibold
                      color=Colors.gray7
                    />
                  </div>
                | _ => <LoadingCensorBar width=100 height=15 />
                }}
             </Col.Grid>
             <Col.Grid col=Col.One>
               <Text block=true value="Block" weight=Text.Semibold color=Colors.gray7 />
             </Col.Grid>
             <Col.Grid col=Col.One>
               <Text
                 block=true
                 value="Status"
                 size=Text.Md
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Center
               />
             </Col.Grid>
             <Col.Grid col=Col.Two>
               <Text
                 block=true
                 value="Gas Fee (BAND)"
                 weight=Text.Semibold
                 color=Colors.gray7
                 align=Text.Center
               />
             </Col.Grid>
             <Col.Grid col=Col.Six>
               <Text block=true value="Actions" weight=Text.Semibold color=Colors.gray7 />
             </Col.Grid>
           </Row.Grid>
         </THead.Grid>}
    {switch (txsSub) {
     | Data(txs) =>
       txs->Belt.Array.size > 0
         ? <TxsTable txsSub msgTransform={transform(accountAddress)} />
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
     | _ => <TxsTable txsSub msgTransform={transform(accountAddress)} />
     }}
    // <TxsTable txsSub msgTransform={transform(accountAddress)} />
    {switch (txsCountSub) {
     | Data(txsCount) =>
       let pageCount = Page.getPageCount(txsCount, pageSize);
       <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />;
     | _ => React.null
     }}
  </div>;
};
