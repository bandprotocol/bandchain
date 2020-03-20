module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);

  let alignRight = style([display(`flex), justifyContent(`flexEnd)]);
};

[@react.component]
let make = (~delegations: list(AccountHook.Account.delegation_t)) => {
  <div className=Styles.tableLowerContainer>
    <VSpacing size=Spacing.md />
    <div className=Styles.hFlex>
      <HSpacing size=Spacing.lg />
      <Text value={delegations |> Belt_List.length |> string_of_int} weight=Text.Semibold />
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
          <Col size=0.6>
            <div className=Styles.alignRight>
              <Text
                block=true
                value="REWARD (BAND)"
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
      {delegations
       ->Belt.List.map(delegation => {
           <TBody key={delegation.validatorAddress} minHeight=50>
             <Row>
               <Col> <HSpacing size=Spacing.lg /> </Col>
               <Col size=0.9>
                 <div className=Styles.hFlex>
                   <Text
                     value="bandvaloper"
                     weight=Text.Semibold
                     spacing={Text.Em(0.02)}
                     code=true
                     block=true
                   />
                   <Text
                     value={delegation.validatorAddress->Js.String.sliceToEnd(~from=11)}
                     spacing={Text.Em(0.02)}
                     block=true
                     code=true
                     ellipsis=true
                     nowrap=true
                   />
                 </div>
               </Col>
               <Col size=0.6>
                 <div className=Styles.alignRight>
                   <Text value={delegation.balance |> Format.fPretty} code=true />
                 </div>
               </Col>
               <Col size=0.6>
                 <div className=Styles.alignRight>
                   <Text value={delegation.reward |> Format.fPretty} code=true />
                 </div>
               </Col>
               <Col> <HSpacing size=Spacing.lg /> </Col>
             </Row>
           </TBody>
         })
       ->Array.of_list
       ->React.array}
      <VSpacing size=Spacing.lg />
    </>
  </div>;
};
