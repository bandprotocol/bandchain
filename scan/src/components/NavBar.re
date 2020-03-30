module Styles = {
  open Css;

  let nav =
    style([
      paddingRight(Spacing.md),
      cursor(`pointer),
      color(Colors.blueGray4),
      fontSize(`px(11)),
    ]);

  let navContainer =
    style([
      paddingTop(Spacing.md),
      paddingBottom(Spacing.md),
      maxWidth(`px(970)),
      marginLeft(`auto),
      marginRight(`auto),
    ]);

  let rFlex = style([display(`flex), flexDirection(`row)]);

  let socialLink =
    style([
      display(`flex),
      flexDirection(`row),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(`px(15)),
    ]);

  let twitterLogo = style([width(`px(20))]);
  let telegramLogo = style([width(`px(20))]);
};

[@react.component]
let make = () => {
  <div className=Styles.navContainer>
    <Row justify=Row.Between>
      <Col>
        <Row>
          {[
             ("Home", Route.HomePage),
             ("Validators", ValidatorHomePage),
             ("Blocks", BlockHomePage),
             ("Transactions", TxHomePage),
             ("Data Sources", DataSourceHomePage),
             ("Oracle Scripts", OracleScriptHomePage),
           ]
           ->Belt.List.map(((v, route)) =>
               <Col key=v>
                 <div className=Styles.nav onClick={_ => route |> Route.redirect}>
                   {v |> React.string}
                 </div>
               </Col>
             )
           ->Array.of_list
           ->React.array}
        </Row>
      </Col>
      <Col>
        <div className=Styles.rFlex>
          <div className=Styles.socialLink>
            <a href="https://twitter.com/bandprotocol" target="_blank" rel="noopener">
              <img src=Images.twitterLogo className=Styles.twitterLogo />
            </a>
          </div>
          <div className=Styles.socialLink>
            <a href="https://t.me/bandprotocol" target="_blank" rel="noopener">
              <img src=Images.telegramLogo className=Styles.telegramLogo />
            </a>
          </div>
        </div>
      </Col>
    </Row>
  </div>;
};
