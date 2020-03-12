module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(10))]);

  let hFlex = style([display(`flex)]);

  let alignRight = style([display(`flex), justifyContent(`flexEnd)]);
};

[@react.component]
let make = () => {
  <div className=Styles.tableLowerContainer>
    <VSpacing size=Spacing.md />
    <div className=Styles.hFlex>
      <HSpacing size=Spacing.lg />
      <Text value="3" weight=Text.Semibold />
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
              color=Colors.grayText
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
                color=Colors.grayText
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
                color=Colors.grayText
              />
            </div>
          </Col>
          <Col> <HSpacing size=Spacing.lg /> </Col>
        </Row>
      </THead>
      {[
         ("bandvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", 30521.534, 2324.23),
         ("bandvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", 30521.534, 2324.23),
         ("bandvaloper1sjllsnramtg3ewxqwwrwjxfgc4n4ef9u2lcnj0", 30521.534, 2324.23),
       ]
       ->Belt.List.map(((validator, amount, reward)) => {
           <TBody key=validator>
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
                     value={validator->Js.String.sliceToEnd(~from=11)}
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
                   <Text value={amount |> Format.fPretty} code=true />
                 </div>
               </Col>
               <Col size=0.6>
                 <div className=Styles.alignRight>
                   <Text value={reward |> Format.fPretty} code=true />
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
