module Styles = {
  open Css;

  let pageWidth = style([maxWidth(`px(970))]);

  let container =
    style([width(`percent(100.)), height(`percent(100.)), position(`relative)]);

  let innerContainer = style([marginLeft(`auto), marginRight(`auto)]);

  let topBarContainer =
    style([
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      width(`percent(100.)),
      paddingTop(Spacing.lg),
      paddingBottom(Spacing.lg),
      backgroundColor(Colors.white),
      border(`px(2), `solid, Colors.blueGray1),
    ]);

  let topBarInner =
    style([display(`flex), width(`percent(100.)), justifyContent(`spaceBetween)]);

  let rFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let bandLogo = style([width(`px(40))]);

  let logoContainer = style([display(`flex), alignItems(`center), cursor(`pointer)]);

  let socialLink =
    style([
      display(`flex),
      flexDirection(`row),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(`px(10)),
    ]);

  let twitterLogo = style([width(`px(15))]);
  let telegramLogo = style([width(`px(15))]);

  let routeContainer =
    style([minHeight(`calc((`sub, `vh(100.), `px(200)))), paddingBottom(`px(20))]);
};

module TopBar = {
  [@react.component]
  let make = () => {
    <div className=Styles.topBarContainer>
      <div className={Css.merge([Styles.topBarInner, Styles.pageWidth])}>
        <div className=Styles.logoContainer onClick={_ => Route.redirect(Route.HomePage)}>
          <Row>
            <Col> <img src=Images.bandLogo className=Styles.bandLogo /> </Col>
            <Col> <HSpacing size=Spacing.sm /> </Col>
            <Col>
              <Text
                value="BandChain"
                size=Text.Xxl
                weight=Text.Bold
                nowrap=true
                color=Colors.gray8
                spacing={Text.Em(0.05)}
              />
              <VSpacing size=Spacing.xs />
              <div className=Styles.rFlex>
                <Text
                  value="EXPLORER"
                  nowrap=true
                  size=Text.Sm
                  weight=Text.Semibold
                  color=Colors.gray6
                  spacing={Text.Em(0.03)}
                />
                <HSpacing size=Spacing.xs />
                <MetadataBadge />
              </div>
            </Col>
            <Col alignSelf=Col.End>
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
        </div>
        <SearchBar />
      </div>
    </div>;
  };
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <TopBar />
    <div className={Css.merge([Styles.innerContainer, Styles.pageWidth])}>
      <NavBar />
      <div className=Styles.routeContainer>
        {switch (ReasonReactRouter.useUrl() |> Route.fromUrl) {
         | HomePage => <HomePage />
         | DataSourceHomePage => <DataSourceHomePage />
         | DataSourceIndexPage(dataSourceID, hashtag) =>
           <DataSourceIndexPage dataSourceID={ID.DataSource.ID(dataSourceID)} hashtag />
         | OracleScriptHomePage => <OracleScriptHomePage />
         | OracleScriptIndexPage(oracleScriptID, hashtag) =>
           <OracleScriptIndexPage oracleScriptID={ID.OracleScript.ID(oracleScriptID)} hashtag />
         | TxHomePage => <TxHomePage />
         | TxIndexPage(txHash) => <TxIndexPage txHash />
         | BlockHomePage => <BlockHomePage />
         | BlockIndexPage(height) => <BlockIndexPage height={ID.Block.ID(height)} />
         | ValidatorHomePage => <ValidatorHomePage />
         | ValidatorIndexPage(address, hashtag) => <ValidatorIndexPage address hashtag />
         | RequestHomePage => <RequestHomePage />
         | RequestIndexPage(reqID) => <RequestIndexPage reqID={ID.Request.ID(reqID)} />
         | AccountIndexPage(address, hashtag) => <AccountIndexPage address hashtag />
         | IBCHomePage => <IBCHomePage />
         | NotFound => <NotFound />
         }}
      </div>
    </div>
  </div>;
};
