module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);

  let alignRight = style([display(`flex), justifyContent(`flexEnd)]);
};

[@react.component]
let make = (~address) =>
  {
    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 10;
    let delegationsCountSub = DelegationSub.getStakeCountByDelegator(address);
    let delegationsSub = DelegationSub.getStakeList(address, ~pageSize, ~page, ());

    let%Sub delegationsCount = delegationsCountSub;
    let%Sub delegations = delegationsSub;

    let pageCount = Page.getPageCount(delegationsCount, pageSize);

    <div className=Styles.tableLowerContainer>
      <VSpacing size=Spacing.md />
      <div className=Styles.hFlex>
        <HSpacing size=Spacing.lg />
        <Text value={delegations |> Belt_Array.length |> string_of_int} weight=Text.Semibold />
        <HSpacing size=Spacing.xs />
        <Text value="Delegated Validators" />
      </div>
      <VSpacing size=Spacing.lg />
      <>
        <THead>
          <Row>
            <Col> <HSpacing size=Spacing.lg /> </Col>
            <Col size=0.9>
              <Text
                block=true
                value="VALIDATOR ADDRESS"
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
            // <Col size=0.6>
            //   <div className=Styles.alignRight>
            //     <Text
            //       block=true
            //       value="REWARD (BAND)"
            //       size=Text.Sm
            //       spacing={Text.Em(0.05)}
            //       weight=Text.Bold
            //       color=Colors.gray6
            //     />
            //   </div>
            // </Col>
            <Col> <HSpacing size=Spacing.lg /> </Col>
          </Row>
        </THead>
        {delegations
         ->Belt.Array.map(delegation => {
             <TBody key={delegation.validatorAddress |> Address.toBech32} minHeight=50>
               <Row>
                 <Col> <HSpacing size=Spacing.lg /> </Col>
                 <Col size=0.9>
                   <div className=Styles.hFlex>
                     <AddressRender address={delegation.validatorAddress} validator=true />
                   </div>
                 </Col>
                 <Col size=0.6>
                   <div className=Styles.alignRight>
                     <Text value={delegation.amount |> Format.fPretty} code=true />
                   </div>
                 </Col>
                 //  <Col size=0.6>
                 //    <div className=Styles.alignRight>
                 //      <Text value={0.00 |> Format.fPretty} code=true />
                 //    </div>
                 //  </Col>
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
