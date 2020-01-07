module Styles = {
  open Css;

  let nav = style([paddingLeft(Spacing.md), cursor(`pointer)]);

  let navContainer =
    style([
      paddingTop(Spacing.md),
      paddingBottom(Spacing.md),
      maxWidth(`px(1100)),
      marginLeft(`auto),
      marginRight(`auto),
      paddingLeft(Spacing.md),
      paddingRight(Spacing.md),
    ]);
};

[@react.component]
let make = () => {
  <div className=Styles.navContainer>
    <Row>
      <Col size=1.> <Text color=Colors.grayText value="Made with <3 by Band Protocol" /> </Col>
      <Col>
        <Row justify=Row.Right>
          {[
             ("Validators", Route.HomePage),
             ("Blocks", BlockHomePage),
             ("Transactions", TxHomePage),
             ("Request Scripts", ScriptHomePage),
             ("Data Providers", HomePage),
             ("OWASM Studio", HomePage),
           ]
           ->Belt.List.map(((v, route)) =>
               <Col key=v>
                 <div className=Styles.nav onClick={_ => route |> Route.redirect}>
                   <Text color=Colors.grayText value=v nowrap=true />
                 </div>
               </Col>
             )
           ->Array.of_list
           ->React.array}
        </Row>
      </Col>
    </Row>
  </div>;
};
