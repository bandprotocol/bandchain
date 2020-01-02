module Styles = {
  open Css;
  let featuredBox = color =>
    style([
      backgroundColor(color),
      borderRadius(`px(8)),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(4), ~blur=`px(20), Css.rgba(0, 0, 0, 0.1))),
      padding(Spacing.lg),
      paddingRight(Spacing.xl),
      marginTop(Spacing.md),
      maxWidth(`px(240)),
      lineHeight(`px(22)),
      cursor(`pointer),
      transition(~duration=100, "transform"),
      hover([transform(translateY(`px(-3)))]),
    ]);

  let recentBox =
    style([
      backgroundColor(white),
      padding(Spacing.lg),
      paddingRight(Spacing.xl),
      marginRight(Spacing.md),
      marginTop(Spacing.lg),
      height(`px(115)),
      width(`px(200)),
      borderRadius(`px(8)),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      cursor(`pointer),
      transition(~duration=100, "margin-top"),
      hover([marginTop(`px(14))]),
    ]);
};

module Featured = {
  open Css;

  [@react.component]
  let make = (~insights, ~title, ~color, ~textColor) => {
    <div className={Styles.featuredBox(color)}>
      <div className={style([opacity(0.7)])}>
        <Text value=insights size=Text.Sm color=textColor />
      </div>
      <div className={style([opacity(0.8)])}>
        <Text value=title size=Text.Xl weight=Text.Semibold color=textColor />
      </div>
    </div>;
  };
};

module Recent = {
  open Css;

  [@react.component]
  let make = (~title, ~hash, ~createdAt) => {
    <div className=Styles.recentBox>
      <div className={style([paddingRight(Spacing.lg)])}>
        <Text value=title size=Text.Lg weight=Text.Semibold />
      </div>
      <VSpacing size=Spacing.md />
      <Text block=true code=true value=hash color=Colors.pink />
      <VSpacing size=Spacing.sm />
      <Text block=true value=createdAt size=Text.Sm color=Colors.grayText />
    </div>;
  };
};

[@react.component]
let make = () => {
  <Row>
    <Col size=2.>
      <Featured
        title="Latest US stock indexes"
        insights="2,384 queries today"
        color=Colors.yellow
        textColor={Css.hex("333333")}
      />
      <Featured
        title="NFL Most Touchdown By Team"
        insights="2,384 queries today"
        color=Colors.orange
        textColor=Css.white
      />
      <Featured
        title="Premier League Scores at Half time"
        insights="2,384 queries today"
        color=Colors.pink
        textColor=Css.white
      />
    </Col>
    <HSpacing size=Spacing.xl />
    <Col size=5.>
      <Text value="Recent Data Scripts" size=Text.Xl weight=Text.Bold block=true />
      <Row wrap=true alignItems=`initial>
        <Col>
          <Recent
            title="Cryptocurrency Price Feed"
            hash="0xe122543771888011"
            createdAt="2 days ago"
          />
        </Col>
        <Col>
          <Recent title="Powerball Lottery" hash="0xe122543771888011" createdAt="2 days ago" />
        </Col>
        <Col>
          <Recent title="Identity Verification" hash="0xe122543771888011" createdAt="2 days ago" />
        </Col>
        <Col>
          <Recent
            title="Cryptocurrency Price Feed"
            hash="0xe122543771888011"
            createdAt="2 days ago"
          />
        </Col>
        <Col>
          <Recent
            title="Cryptocurrency Price Feed"
            hash="0xe122543771888011"
            createdAt="2 days ago"
          />
        </Col>
        <Col>
          <Recent
            title="Cryptocurrency Price Feed"
            hash="0xe122543771888011"
            createdAt="2 days ago"
          />
        </Col>
      </Row>
    </Col>
    <Col>
      <Text block=true value="358" size=Text.Xxl weight=Text.Bold />
      <VSpacing size=Spacing.sm />
      <Text block=true value="DATA SCRIPTS" size=Text.Sm color=Colors.purple />
      <VSpacing size=Spacing.xs />
      <Text block=true value="CREATED" size=Text.Sm color=Colors.purple />
      <VSpacing size=Spacing.sm />
      <VSpacing size=Spacing.xl />
      <Text block=true value="48" size=Text.Xxl weight=Text.Bold />
      <VSpacing size=Spacing.sm />
      <Text block=true value="DATA PROVIDERS" size=Text.Sm color=Colors.purple />
    </Col>
  </Row>;
};
