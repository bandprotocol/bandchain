module Styles = {
  open Css;

  let vFlex = align => style([display(`flex), flexDirection(`row), alignItems(align)]);

  let header =
    style([display(`flex), flexDirection(`row), alignItems(`center), height(`px(50))]);

  let logo = style([minWidth(`px(50)), marginRight(`px(10))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.gray7),
    ]);

  let fullWidth = style([width(`percent(100.0)), display(`flex)]);

  let withWidth = w => style([width(`px(w))]);

  let fillLeft = style([marginLeft(`auto)]);
};

let renderBody = (reserveIndex, blockSub: ApolloHooks.Subscription.variant(BlockSub.t)) => {
  <TBody
    key={
      switch (blockSub) {
      | Data({height}) => height |> ID.Block.toString
      | _ => reserveIndex |> string_of_int
      }
    }>
    <Row minHeight={`px(30)}>
      <Col> <HSpacing size=Spacing.md /> </Col>
      <Col size=1.11>
        {switch (blockSub) {
         | Data({height}) => <TypeID.Block id=height />
         | _ => <LoadingCensorBar width=65 height=15 />
         }}
      </Col>
      <Col size=3.93>
        {switch (blockSub) {
         | Data({hash}) =>
           <div className={Styles.withWidth(330)}>
             <Text
               value={hash |> Hash.toHex(~upper=true)}
               weight=Text.Medium
               block=true
               code=true
               ellipsis=true
             />
           </div>
         | _ => <LoadingCensorBar width=300 height=15 />
         }}
      </Col>
      <Col size=2.1>
        {switch (blockSub) {
         | Data({timestamp}) =>
           <Timestamp time=timestamp size=Text.Md weight=Text.Regular code=true />
         | _ => <LoadingCensorBar width=150 height=15 />
         }}
      </Col>
      <Col size=1.5>
        {switch (blockSub) {
         | Data({validator}) =>
           <div className={Styles.withWidth(150)}>
             <ValidatorMonikerLink
               validatorAddress={validator.operatorAddress}
               moniker={validator.moniker}
             />
           </div>
         | _ => <LoadingCensorBar width=150 height=15 />
         }}
      </Col>
      <Col size=1.05>
        <Row>
          <div className=Styles.fillLeft />
          {switch (blockSub) {
           | Data({txn}) => <Text value={txn |> Format.iPretty} code=true weight=Text.Medium />
           | _ => <LoadingCensorBar width=40 height=15 isRight=true />
           }}
        </Row>
      </Col>
      <Col> <HSpacing size=Spacing.md /> </Col>
    </Row>
  </TBody>;
};

[@react.component]
let make = () => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let allSub = Sub.all2(BlockSub.getList(~pageSize, ~page, ()), BlockSub.count());

  <>
    <Row>
      <div className=Styles.header>
        <img src=Images.blockLogo className=Styles.logo />
        <Text
          value="ALL BLOCKS"
          weight=Text.Medium
          size=Text.Md
          spacing={Text.Em(0.06)}
          height={Text.Px(15)}
          nowrap=true
          block=true
          color=Colors.gray7
        />
        {switch (allSub) {
         | Data((_, blocksCount)) =>
           <>
             <div className=Styles.seperatedLine />
             <Text
               value={blocksCount->Format.iPretty ++ " in total"}
               size=Text.Md
               weight=Text.Thin
               spacing={Text.Em(0.06)}
               color=Colors.gray7
               nowrap=true
             />
           </>
         | _ => React.null
         }}
      </div>
    </Row>
    <VSpacing size=Spacing.xl />
    <THead>
      <Row>
        <Col> <HSpacing size=Spacing.md /> </Col>
        {[
           ("BLOCK", 1.11, false),
           ("BLOCK HASH", 3.80, false),
           ("TIMESTAMP", 2.1, false),
           ("PROPOSER", 1.55, false),
           ("TXN", 1.05, true),
         ]
         ->Belt.List.map(((title, size, alignRight)) => {
             <Col size key=title justifyContent=Col.Start>
               <div className={Styles.vFlex(`flexEnd)}>
                 {alignRight ? <div className=Styles.fillLeft /> : React.null}
                 <Text
                   value=title
                   size=Text.Sm
                   weight=Text.Semibold
                   color=Colors.gray6
                   spacing={Text.Em(0.1)}
                 />
               </div>
             </Col>
           })
         ->Array.of_list
         ->React.array}
        <Col> <HSpacing size=Spacing.md /> </Col>
      </Row>
    </THead>
    {switch (allSub) {
     | Data((blocks, blocksCount)) =>
       let pageCount = Page.getPageCount(blocksCount, pageSize);
       <>
         {blocks->Belt_Array.mapWithIndex((i, e) => renderBody(i, Sub.resolve(e)))->React.array}
         <VSpacing size=Spacing.lg />
         <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
         <VSpacing size=Spacing.lg />
       </>;
     | _ =>
       Belt_Array.make(10, ApolloHooks.Subscription.NoData)
       ->Belt_Array.mapWithIndex((i, noData) => renderBody(i, noData))
       ->React.array
     }}
  </>;
};
