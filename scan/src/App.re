module Styles = {
  open Css;

  let container =
    style([
      width(`percent(100.)),
      height(`percent(100.)),
      position(`relative),
    ]);

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
      padding4(
        ~top=`px(0),
        ~bottom=`px(0),
        ~left=Spacing.sm,
        ~right=Spacing.sm,
      ),
      height(`px(20)),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(Spacing.xs),
      marginTop(`px(1)),
    ]);

  let uFlex = style([display(`flex), flexDirection(`row)]);

  let logo35 = style([width(`px(35))]);
  let logo20 = style([width(`px(20))]);

  let skipRight = style([marginLeft(`px(30))]);

  let socialLink =
    style([
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      width(`px(50)),
    ]);
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <NavBar />
    <div className=Styles.innerContainer>
      <div className=Styles.mainBar>
        <Row>
          <Col size=1.>
            <img src=Images.bandLogo className=Styles.logo35 />
          </Col>
          <Col size=4.>
            <div className=Styles.uFlex>
              <Text value="D3N" size=Text.Xxl weight=Text.Bold nowrap=true />
              <div className=Styles.version>
                <Text
                  value="v1.0 TESTNET"
                  size=Text.Sm
                  color=Colors.pink
                  nowrap=true
                />
              </div>
            </div>
            <Text value="Data Request Explorer" nowrap=true />
          </Col>
        </Row>
        <SearchBar />
        <div className=Styles.skipRight />
        <div className=Styles.socialLink>
          <a
            href="https://twitter.com/bandprotocol"
            target="_blank"
            rel="noopener">
            <img src=Images.twitterLogo className=Styles.logo20 />
          </a>
        </div>
        <div className=Styles.socialLink>
          <a href="https://t.me/bandprotocol" target="_blank" rel="noopener">
            <img src=Images.telegramLogo className=Styles.logo20 />
          </a>
        </div>
      </div>
      <HomePage />
    </div>
    <Footer />
  </div>;
};
