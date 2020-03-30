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
      border(`px(1), `solid, hex("F9F7FA")),
    ]);

  let topBarInner =
    style([display(`flex), width(`percent(100.)), justifyContent(`spaceBetween)]);

  let version =
    style([
      display(`flex),
      borderRadius(`px(10)),
      backgroundColor(`hex("EBECFF")),
      padding2(~v=`pxFloat(4.3), ~h=`px(8)),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(Spacing.xs),
      marginTop(`px(1)),
    ]);

  let rFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let bandLogo = style([width(`px(40))]);

  let logoContainer = style([display(`flex), alignItems(`center), cursor(`pointer)]);

  let routeContainer = style([minHeight(`calc((`sub, `vh(100.), `px(300))))]);
};

module TopBar = {
  [@react.component]
  let make = () =>
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
                  color={Css.hex("777777")}
                  spacing={Text.Em(0.03)}
                />
                <HSpacing size=Spacing.xs />
                <div className=Styles.version>
                  <Text
                    value="v1.0 TESTNET"
                    size=Text.Xs
                    color={Css.hex("535BBF")}
                    nowrap=true
                    weight=Text.Semibold
                    spacing={Text.Em(0.03)}
                  />
                </div>
              </div>
            </Col>
          </Row>
        </div>
        <SearchBar />
      </div>
    </div>;
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <TopBar />
    <div className={Css.merge([Styles.innerContainer, Styles.pageWidth])}>
      <NavBar />
      /* route handle */
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
         | RequestIndexPage(reqID, hashtag) => <RequestIndexPage reqID _hashtag=hashtag />
         | AccountIndexPage(address, hashtag) => <AccountIndexPage address hashtag />
         | ValidatorIndexPage(address, hashtag) => <ValidatorIndexPage address hashtag />
         | NotFound => <NotFound />
         }}
      </div>
    </div>
    // <Footer />
  </div>;
};
