module Styles = {
  open Css;

  let bg =
    style([
      width(`percent(100.)),
      height(`px(300)),
      left(`px(0)),
      bottom(`px(0)),
      position(`relative),
      background(hex("F6F3FA")),
      before([
        position(`absolute),
        contentRule(""),
        background(hex("F6F3FA")),
        width(`percent(100.)),
        height(`px(300)),
        transform(`skewY(`deg(6.))),
        zIndex(-1),
        top(`px(-150)),
      ]),
    ]);

  let footerContainer =
    style([
      marginTop(Spacing.xl),
      maxWidth(`px(1100)),
      marginLeft(`auto),
      marginRight(`auto),
      paddingLeft(Spacing.md),
      paddingRight(Spacing.md),
    ]);

  let vFlex = style([display(`flex), flexDirection(`column)]);

  let logo = style([width(`px(35))]);

  let logoContainer = style([width(`px(50))]);

  let messageBlock = style([maxWidth(`px(220)), lineHeight(`px(20))]);

  let subHeader = style([height(`px(40)), display(`flex), alignItems(`center)]);

  let commonLink = style([marginBottom(Spacing.md), textDecoration(`none), display(`block)]);
};

let renderSubHeader = text =>
  <>
    <div className=Styles.subHeader>
      <Text value=text size=Text.Md weight=Text.Bold color=Colors.grayText />
    </div>
    <VSpacing size=Spacing.lg />
  </>;

let renderCommonLink = (url, text) =>
  <a href=url target="_blank" rel="noopener" key=url className=Styles.commonLink>
    <Text value=text size=Text.Md weight=Text.Semibold color=Colors.purple />
  </a>;

[@react.component]
let make = () => {
  <div className=Styles.bg>
    <div className=Styles.footerContainer>
      <Row justify=Row.Between>
        <Col size=1.5>
          <div className=Styles.subHeader>
            <Row>
              <div className=Styles.logoContainer>
                <Col size=1.> <img src=Images.bandLogo className=Styles.logo /> </Col>
              </div>
              <Col size=4.>
                <div className=Styles.vFlex>
                  <Text value="POWERED BY" size=Text.Md weight=Text.Regular nowrap=true />
                  <Text value="Band Protocol" size=Text.Xxl weight=Text.Bold nowrap=true />
                </div>
              </Col>
            </Row>
          </div>
          <VSpacing size=Spacing.lg />
          <div className=Styles.messageBlock>
            <Text
              value="D3N is a block explorer and an analytic platform for Band Protocol, a decentralized platform for reliable and provably-secure data on blockchain."
              size=Text.Md
              weight=Text.Regular
              block=true
            />
          </div>
        </Col>
        <Col size=1.>
          {renderSubHeader("D3N PROJECT")}
          {[|
             ("https://bandprotocol.com", "Band Protocol Website"),
             ("https://google.com", "D3N Blockchain"),
             ("https://metamask.com", "D3N Wallet"),
           |]
           ->Belt.Array.map(((url, text)) => {renderCommonLink(url, text)})
           |> React.array}
        </Col>
        <Col size=1.>
          {renderSubHeader("COMMUNITY")}
          {[|
             ("https://github.com", "Open Source Repositories"),
             ("https://google.com", "Developer Docs"),
             ("https://etherscan.com", "Network Status"),
           |]
           ->Belt.Array.map(((url, text)) => {renderCommonLink(url, text)})
           |> React.array}
        </Col>
        <Col size=0.5>
          {renderSubHeader("SOCIAL LINKS")}
          {[|
             ("https://twitter.com", "Twitter"),
             ("https://telegram.org", "Telegram"),
             ("https://medium.com", "Medium"),
           |]
           ->Belt.Array.map(((url, text)) => {renderCommonLink(url, text)})
           |> React.array}
        </Col>
      </Row>
    </div>
  </div>;
};
