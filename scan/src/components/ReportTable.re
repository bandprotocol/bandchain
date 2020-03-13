module Styles = {
  open Css;

  let txhash = style([marginLeft(`px(20))]);
};

[@react.component]
let make = (~reports: list(RequestHook.Report.t)) => {
  let infoOpt = React.useContext(GlobalContext.context);
  let validators =
    switch (infoOpt) {
    | Some(info) => info.validators
    | None => []
    };
  <>
    <THead>
      <Row>
        <Col> <div className=Styles.txhash /> </Col>
        <Col size=1.0>
          <div className=TElement.Styles.hashContainer>
            <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
        <Col size=0.35>
          <Text block=true value="BLOCK" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.4>
          <Text block=true value="AGE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=1.0>
          <Text block=true value="FROM" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </Col>
        <Col size=0.6> <div /> </Col>
        <Col size=0.9>
          <div className=TElement.Styles.feeContainer>
            <Text block=true value="VALUE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
          </div>
        </Col>
      </Row>
    </THead>
    {reports
     ->Belt.List.map(({reporter, txHash, reportedAtHeight, reportedAtTime, values}) => {
         let moniker =
           validators
           ->Belt_List.keepMap(validator =>
               validator.operatorAddress == reporter ? Some(validator.moniker) : None
             )
           ->Belt_List.get(0)
           ->Belt_Option.getWithDefault("Unknown");
         <TBody key={txHash |> Hash.toHex}>

             <Row alignItems=Css.flexStart>
               <Col> <div className=Styles.txhash /> </Col>
               <Col size=1.0> <TElement elementType={txHash->TElement.HashWithLink} /> </Col>
               <Col size=0.35> <TElement elementType={reportedAtHeight->TElement.Height} /> </Col>
               <Col size=0.4> <TElement elementType={reportedAtTime->TElement.Timestamp} /> </Col>
               <Col size=1.0>
                 <TElement elementType={reporter->TElement.Address} />
                 <VSpacing size=Spacing.sm />
                 <TElement elementType={moniker->TElement.Detail} />
               </Col>
             </Row>
           </TBody>;
           //  <Col size=0.6>
           //    {values
           //     ->Belt.Array.map(((source, _)) =>
           //         <>
           //           <TElement elementType={source->TElement.Source} />
           //           <VSpacing size=Spacing.sm />
           //         </>
           //       )
           //     ->React.array}
           //  </Col>
           //  <Col size=0.9>
           //    {values
           //     ->Belt.Array.map(((_, value)) =>
           //         <>
           //           <TElement elementType={value->TElement.Value} />
           //           <VSpacing size=Spacing.sm />
           //         </>
           //       )
           //     ->React.array}
           //  </Col>
       })
     ->Array.of_list
     ->React.array}
  </>;
};
