module Styles = {
  open Css;

  let header = (theme: Theme.t) =>
    style([
      paddingTop(Spacing.lg),
      backgroundColor(theme.mainBg),
      borderBottom(`px(2), `solid, Colors.blueGray1),
      zIndex(3),
      Media.mobile([
        padding(Spacing.md),
        marginBottom(`zero),
        position(`sticky),
        top(`zero),
        width(`percent(100.)),
      ]),
    ]);

  let leftContainer = style([display(`flex), alignItems(`center), width(`percent(100.))]);
  let bandLogo = style([width(`px(40)), Media.mobile([width(`px(34))])]);
  let cmcLogo = style([width(`px(15)), height(`px(15))]);
  let blockImage = style([display(`block)]);

  let socialLink = style([marginLeft(`px(10)), display(`flex), textDecoration(`none)]);

  let link = style([cursor(`pointer)]);
};

module LinkToHome = {
  [@react.component]
  let make = (~children) => {
    <Link className=Styles.link route=Route.HomePage> children </Link>;
  };
};

module ToggleThemeButton = {
  open Css;

  module Styles = {
    let button = isDarkMode =>
      style([
        backgroundColor(Colors.white),
        padding2(~v=`px(3), ~h=`px(6)),
        borderRadius(`px(8)),
        border(`px(1), `solid, isDarkMode ? Colors.white : Colors.black),
        marginLeft(`px(5)),
        cursor(`pointer),
        outlineStyle(`none),
      ]);
  };

  [@react.component]
  let make = () => {
    let ({ThemeContext.isDarkMode}, toggle) = React.useContext(ThemeContext.context);

    <button className={Styles.button(isDarkMode)} onClick={_ => toggle()}>
      <Icon name={isDarkMode ? "fas fa-sun" : "fas fa-moon"} size=14 color=Colors.black />
    </button>;
  };
};

module DesktopRender = {
  [@react.component]
  let make = () => {
    let ({ThemeContext.theme}, _) = React.useContext(ThemeContext.context);

    <header className={Styles.header(theme)}>
      <div className="container">
        <Row alignItems=Row.Center marginBottom=12>
          <Col col=Col.Five>
            <div className=Styles.leftContainer>
              <LinkToHome> <img src=Images.bandLogo className=Styles.bandLogo /> </LinkToHome>
              <HSpacing size=Spacing.md />
              <div className={CssHelper.flexBox(~direction=`column, ~align=`flexStart, ())}>
                <LinkToHome>
                  <Text
                    value="BandChain"
                    size=Text.Xxl
                    weight=Text.Bold
                    nowrap=true
                    color=Colors.gray8
                    spacing={Text.Em(0.05)}
                  />
                </LinkToHome>
                <VSpacing size=Spacing.xs />
                <div className={CssHelper.flexBox()}>
                  <LinkToHome>
                    <Text
                      value="CosmoScan"
                      nowrap=true
                      size=Text.Sm
                      weight=Text.Semibold
                      color=Colors.gray6
                      spacing={Text.Em(0.03)}
                    />
                    <HSpacing size=Spacing.xs />
                  </LinkToHome>
                  <HSpacing size=Spacing.xs />
                  <ChainIDBadge />
                  <div className={CssHelper.flexBox(~align=`center, ())}>
                    <a
                      href="https://twitter.com/bandprotocol"
                      target="_blank"
                      rel="noopener"
                      className=Styles.socialLink>
                      <Icon name="fab fa-twitter" color=Colors.bandBlue size=16 />
                    </a>
                    <a
                      href="https://t.me/bandprotocol"
                      target="_blank"
                      rel="noopener"
                      className=Styles.socialLink>
                      <Icon name="fab fa-telegram-plane" color=Colors.bandBlue size=17 />
                    </a>
                    <a
                      href="https://coinmarketcap.com/currencies/band-protocol/"
                      target="_blank"
                      rel="noopener"
                      className=Styles.socialLink>
                      <img
                        src=Images.cmcLogo
                        className={Css.merge([Styles.cmcLogo, Styles.blockImage])}
                      />
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </Col>
          <Col col=Col.Seven> <SearchBar /> </Col>
        </Row>
        <Row alignItems=Row.Center>
          <Col col=Col.Eight> <NavBar /> </Col>
          <Col col=Col.Four>
            <div className={CssHelper.flexBox(~justify=`flexEnd, ())}>
              <UserAccount />
              <ToggleThemeButton />
            </div>
          </Col>
        </Row>
      </div>
    </header>;
  };
};

module MobileRender = {
  [@react.component]
  let make = () => {
    let ({ThemeContext.theme}, _) = React.useContext(ThemeContext.context);

    <header className={Styles.header(theme)}>
      <Row alignItems=Row.Center>
        <Col colSm=Col.Six>
          <div className={CssHelper.flexBox(~align=`flexEnd, ())}>
            <LinkToHome>
              <img
                src=Images.bandLogo
                className={Css.merge([Styles.bandLogo, Styles.blockImage])}
              />
            </LinkToHome>
            <HSpacing size=Spacing.sm />
            <LinkToHome>
              <div className={CssHelper.flexBox(~direction=`column, ~align=`flexStart, ())}>
                <Text
                  value="BandChain"
                  size=Text.Lg
                  weight=Text.Bold
                  nowrap=true
                  color=Colors.gray8
                  spacing={Text.Em(0.05)}
                />
                <VSpacing size=Spacing.xs />
                <Text
                  value="CosmoScan"
                  nowrap=true
                  size=Text.Sm
                  color=Colors.gray6
                  spacing={Text.Em(0.03)}
                />
              </div>
            </LinkToHome>
          </div>
        </Col>
        <Col colSm=Col.Six>
          <div className={CssHelper.flexBox(~justify=`flexEnd, ~wrap=`nowrap, ())}>
            <ChainIDBadge />
            <NavBar />
          </div>
        </Col>
      </Row>
    </header>;
  };
};

[@react.component]
let make = () => {
  Media.isMobile() ? <MobileRender /> : <DesktopRender />;
};
