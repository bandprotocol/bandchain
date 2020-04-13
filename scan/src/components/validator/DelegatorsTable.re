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
};

[@react.component]
let make = (~address) =>
  {
    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 10;

    let delegatorsSub = DelegationSub.getDelegatorsByValidator(address, ~pageSize, ~page, ());
    let delegatorCountSub = DelegationSub.getDelegatorCountByValidator(address);

    let%Sub delegators = delegatorsSub;
    let%Sub delegatorCount = delegatorCountSub;

    let pageCount = Page.getPageCount(delegatorCount, pageSize);

    <div className=Styles.tableWrapper>
      <Row>
        <HSpacing size={`px(25)} />
        <Text value={delegatorCount |> string_of_int} weight=Text.Bold />
        <HSpacing size={`px(5)} />
        <Text value="Delegators" />
      </Row>
      <VSpacing size=Spacing.lg />
      {delegators->Belt_Array.length > 0
         ? <>
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
             </THead>
             {delegators
              ->Belt.Array.map(({amount, sharePercentage, delegatorAddress}) => {
                  <TBody>
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
                            value={amount |> Format.fPretty}
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
           </div>}
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
