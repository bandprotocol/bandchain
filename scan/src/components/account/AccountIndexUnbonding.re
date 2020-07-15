module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);

  let alignRight = style([display(`flex), justifyContent(`flexEnd)]);
};

[@react.component]
let make = (~address) =>
  {
    let currentTime =
      React.useContext(TimeContext.context)
      |> MomentRe.Moment.defaultUtc
      |> MomentRe.Moment.format("YYYY-MM-DDTHH:mm:ss.SSSSSS");

    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 10;

    let unbondingListSub =
      UnbondingSub.getUnbondingByDelegator(address, currentTime, ~pageSize, ~page, ());
    let unbondingCountSub = UnbondingSub.getUnbondingCountByDelegator(address, currentTime);

    let%Sub unbondingCount = unbondingCountSub;
    let%Sub unbondingList = unbondingListSub;

    let pageCount = Page.getPageCount(unbondingCount, pageSize);

    <div>
      <VSpacing size=Spacing.md />
      <div className=Styles.hFlex>
        <HSpacing size=Spacing.lg />
        <Text value={unbondingCount |> string_of_int} weight=Text.Semibold />
        <HSpacing size=Spacing.xs />
        <Text value="Unbonding Entries" />
      </div>
      <VSpacing size=Spacing.lg />
      <>
        <THead>
          <Row>
            <Col> <HSpacing size=Spacing.lg /> </Col>
            <Col size=1.>
              <Text
                block=true
                value="VALIDATOR"
                size=Text.Sm
                weight=Text.Bold
                spacing={Text.Em(0.05)}
                color=Colors.gray6
              />
            </Col>
            <Col size=0.6>
              <div className=Styles.alignRight>
                <Text
                  block=true
                  value="AMOUNT (BAND)"
                  size=Text.Sm
                  weight=Text.Bold
                  spacing={Text.Em(0.05)}
                  color=Colors.gray6
                />
              </div>
            </Col>
            <Col size=1.>
              <div className=Styles.alignRight>
                <Text
                  block=true
                  value="UNBONDED AT"
                  size=Text.Sm
                  spacing={Text.Em(0.05)}
                  weight=Text.Bold
                  color=Colors.gray6
                />
              </div>
            </Col>
            <Col> <HSpacing size=Spacing.lg /> </Col>
          </Row>
        </THead>
        {unbondingList
         ->Belt.Array.map(unbondingEntry => {
             <TBody
               key={unbondingEntry.validator.operatorAddress |> Address.toBech32} minHeight=50>
               <Row>
                 <Col> <HSpacing size=Spacing.lg /> </Col>
                 <Col size=1.>
                   <div className=Styles.hFlex>
                     <ValidatorMonikerLink
                       validatorAddress={unbondingEntry.validator.operatorAddress}
                       moniker={unbondingEntry.validator.moniker}
                       identity={unbondingEntry.validator.identity}
                       width={`px(300)}
                     />
                   </div>
                 </Col>
                 <Col size=0.6>
                   <div className=Styles.alignRight>
                     <Text
                       value={
                         unbondingEntry.amount |> Coin.getBandAmountFromCoin |> Format.fPretty
                       }
                       code=true
                     />
                   </div>
                 </Col>
                 <Col size=1.>
                   <div className=Styles.alignRight>
                     <Text
                       value={
                         unbondingEntry.completionTime
                         |> MomentRe.Moment.format(Config.timestampFormat)
                         |> String.uppercase_ascii
                       }
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
        <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
      </>
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
