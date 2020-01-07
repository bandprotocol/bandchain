module Styles = {
  open Css;

  let container =
    style([width(`percent(100.)), height(`percent(100.)), position(`relative)]);

  let innerContainer =
    style([
      maxWidth(`px(1100)),
      marginLeft(`auto),
      marginRight(`auto),
      paddingLeft(Spacing.md),
      paddingRight(Spacing.md),
    ]);

  /* Main bar */
  let mainBar = style([display(`flex), paddingTop(Spacing.lg)]);
  let version =
    style([
      display(`flex),
      borderRadius(`px(10)),
      backgroundColor(Colors.pinkLight),
      padding4(~top=`px(0), ~bottom=`px(0), ~left=Spacing.sm, ~right=Spacing.sm),
      height(`px(20)),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(Spacing.xs),
      marginTop(`px(1)),
    ]);

  let uFlex = style([display(`flex), flexDirection(`row)]);

  let logo = style([width(`px(35))]);
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <NavBar />
    <div className=Styles.innerContainer>
      <div className=Styles.mainBar>
        <Row>
          <Col size=1.> <img src=Images.bandLogo className=Styles.logo /> </Col>
          <Col size=4.>
            <div className=Styles.uFlex>
              <Text value="D3N" size=Text.Xxl weight=Text.Bold nowrap=true />
              <div className=Styles.version>
                <Text value="v1.0 TESTNET" size=Text.Sm color=Colors.pink nowrap=true />
              </div>
            </div>
            <Text value="Data Request Explorer" nowrap=true />
          </Col>
        </Row>
        <SearchBar />
      </div>
      /* route handle */
      {switch (ReasonReactRouter.useUrl() |> Route.fromUrl) {
       | HomePage => <HomePage />
       | ScriptHomePage => <ScriptHomePage />
       | ScriptIndexPage(codeHash, hashtag) => <ScriptIndexPage codeHash hashtag />
       | TxHomePage => <TxHomePage />
       | TxIndexPage(txHash, hashtag) => <TxIndexPage txHash hashtag />
       | BlockHomePage => <BlockHomePage />
       | BlockIndexPage(height, hashtag) => <BlockIndexPage height hashtag />
       | NotFound =>
         ReasonReactRouter.replace("/");
         React.null;
       }}
    </div>
    <Footer />
  </div>;
};
