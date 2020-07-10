module Styles = {
  open Css;

  let vFlex = align => style([display(`flex), flexDirection(`row), alignItems(align)]);

  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);

  let icon = style([width(`px(80)), height(`px(80))]);
  let iconWrapper =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
    ]);

  let withWidth = w => style([width(`px(w))]);

  let fillLeft = style([marginLeft(`auto)]);
  let pagination = style([height(`px(30))]);
};

module Header = {
  [@react.component]
  let make = () => {
    <THead>
      <Row>
        <Col> <HSpacing size=Spacing.lg /> </Col>
        <Col size=1.4>
          <Text
            block=true
            value="DELEGATOR"
            size=Text.Sm
            weight=Text.Semibold
            color=Colors.gray6
            spacing={Text.Em(0.05)}
          />
        </Col>
        <Col size=1.45>
          <div className={Styles.vFlex(`flexEnd)}>
            <div className=Styles.fillLeft />
            <Text
              block=true
              value="SHARE (%)"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col size=1.45>
          <div className={Styles.vFlex(`flexEnd)}>
            <div className=Styles.fillLeft />
            <Text
              block=true
              value="AMOUNT (BAND)"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.gray6
              spacing={Text.Em(0.05)}
            />
          </div>
        </Col>
        <Col> <HSpacing size=Spacing.lg /> </Col>
      </Row>
    </THead>;
  };
};

module Loading = {
  [@react.component]
  let make = () => {
    <>
      <Header />
      {Belt_Array.make(
         10,
         <Row>
           <Col> <HSpacing size=Spacing.lg /> </Col>
           <Col size=1.4> <LoadingCensorBar width=300 height=16 /> </Col>
           <Col size=1.30>
             <div className={Styles.vFlex(`flexEnd)}>
               <div className=Styles.fillLeft />
               <LoadingCensorBar width=70 height=16 />
             </div>
           </Col>
           <Col size=1.45>
             <div className={Styles.vFlex(`flexEnd)}>
               <div className=Styles.fillLeft />
               <LoadingCensorBar width=70 height=16 />
             </div>
           </Col>
           <Col> <HSpacing size=Spacing.lg /> </Col>
         </Row>,
       )
       ->Belt.Array.mapWithIndex((i, e) => {<TBody key={i |> string_of_int}> e </TBody>})
       ->React.array}
      <VSpacing size=Spacing.lg />
      <div className=Styles.pagination />
    </>;
  };
};

[@react.component]
let make = (~address) => {
  let (page, setPage) = React.useState(_ => 1);
  let pageSize = 10;

  let delegatorsSub = DelegationSub.getDelegatorsByValidator(address, ~pageSize, ~page, ());
  let delegatorCountSub = DelegationSub.getDelegatorCountByValidator(address);

  let allSub = Sub.all2(delegatorCountSub, delegatorsSub);

  <div className=Styles.tableWrapper>
    <Row>
      <HSpacing size={`px(25)} />
      {switch (allSub) {
       | Data((delegatorCount, _)) =>
         <>
           <Text value={delegatorCount |> string_of_int} weight=Text.Bold />
           <HSpacing size={`px(5)} />
           <Text value="Delegators" />
         </>
       | _ => <LoadingCensorBar width=100 height=20 />
       }}
    </Row>
    <VSpacing size=Spacing.lg />
    {switch (allSub) {
     | Data((delegatorCount, delegators)) =>
       let pageCount = Page.getPageCount(delegatorCount, pageSize);
       delegators->Belt_Array.length > 0
         ? <>
             <Header />
             {delegators
              ->Belt.Array.map(({amount, sharePercentage, delegatorAddress}) => {
                  <TBody key={delegatorAddress |> Address.toBech32}>
                    <Row>
                      <Col> <HSpacing size=Spacing.lg /> </Col>
                      <Col size=1.4> <AddressRender address=delegatorAddress /> </Col>
                      <Col size=1.30>
                        <div className={Styles.vFlex(`flexEnd)}>
                          <div className=Styles.fillLeft />
                          <Text
                            block=true
                            value={sharePercentage |> Format.fPretty}
                            size=Text.Md
                            weight=Text.Regular
                            color=Colors.gray7
                            spacing={Text.Em(0.05)}
                            code=true
                          />
                        </div>
                      </Col>
                      <Col size=1.45>
                        <div className={Styles.vFlex(`flexEnd)}>
                          <div className=Styles.fillLeft />
                          <Text
                            block=true
                            value={amount |> Coin.getBandAmountFromCoin |> Format.fPretty}
                            size=Text.Md
                            weight=Text.Regular
                            color=Colors.gray7
                            spacing={Text.Em(0.05)}
                            code=true
                          />
                        </div>
                      </Col>
                      <Col> <HSpacing size=Spacing.lg /> </Col>
                    </Row>
                  </TBody>
                })
              ->React.array}
             <VSpacing size=Spacing.lg />
             <Pagination
               currentPage=page
               pageCount
               onPageChange={newPage => setPage(_ => newPage)}
             />
           </>
         : <div className=Styles.iconWrapper>
             <VSpacing size={`px(30)} />
             <img src=Images.noRequestIcon className=Styles.icon />
             <VSpacing size={`px(40)} />
             <Text block=true value="NO DELEGATORS" weight=Text.Regular color=Colors.blue4 />
             <VSpacing size={`px(15)} />
           </div>;
     | _ => <Loading />
     }}
  </div>;
};
