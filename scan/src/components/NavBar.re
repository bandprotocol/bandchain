module Styles = {
  open Css;

  let nav =
    style([
      paddingRight(Spacing.md),
      cursor(`pointer),
      color(Colors.blueGray4),
      fontSize(`px(11)),
      hover([color(Colors.blueGray4)]),
      active([color(Colors.blueGray4)]),
    ]);

  let navContainer =
    style([
      paddingTop(Spacing.md),
      paddingBottom(Spacing.md),
      maxWidth(`px(970)),
      marginLeft(`auto),
      marginRight(`auto),
      minHeight(`px(70)),
    ]);
};

[@react.component]
let make = () => {
  <div className=Styles.navContainer>
    <Row justify=Row.Between alignItems=`flexStart>
      <Col>
        <Row>
          {[
             ("Home", Route.HomePage),
             ("Validators", ValidatorHomePage),
             ("Blocks", BlockHomePage),
             ("Transactions", TxHomePage),
             ("Data Sources", DataSourceHomePage),
             ("Oracle Scripts", OracleScriptHomePage),
             ("Requests", RequestHomePage),
            //  ("IBCs", IBCHomePage),
           ]
           ->Belt.List.map(((v, route)) =>
               <Col key=v> <Link className=Styles.nav route> {v |> React.string} </Link> </Col>
             )
           ->Array.of_list
           ->React.array}
        </Row>
      </Col>
      <Col> <UserAccount /> </Col>
    </Row>
  </div>;
};
