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
  let make = (~insights, ~title, ~color, ~textColor, ~onClick) => {
    <div className={Styles.featuredBox(color)} onClick>
      <div className={style([opacity(0.7)])}>
        <Text value=insights size=Text.Sm color=textColor ellipsis=true block=true />
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
  let make = (~title, ~hash, ~createdAt, ~onClick) => {
    <div className=Styles.recentBox onClick>
      <div className={style([paddingRight(Spacing.lg)])}>
        <Text value=title size=Text.Lg weight=Text.Semibold />
      </div>
      <VSpacing size=Spacing.md />
      <Text block=true code=true value={hash->Hash.toHex} color=Colors.pink ellipsis=true />
      <VSpacing size=Spacing.sm />
      <TimeAgos time=createdAt />
    </div>;
  };
};

let renderFeatured = (recentScripts, index, color, textColor) => {
  let {ScriptHook.Script.info, txHash} =
    recentScripts->Belt.List.getExn(index mod recentScripts->Belt.List.length);
  <Featured
    title={info.name}
    insights={txHash->Hash.toHex}
    color
    textColor
    onClick={_ =>
      Route.redirect(Route.ScriptIndexPage(txHash |> Hash.toHex, Route.ScriptTransactions))
    }
  />;
};

let renderScript = (recentScripts, index) => {
  let {ScriptHook.Script.info, txHash, createdAtTime} =
    recentScripts->Belt.List.getExn(index mod recentScripts->Belt.List.length);
  <Col>
    <Recent
      title={info.name}
      hash=txHash
      createdAt=createdAtTime
      onClick={_ =>
        Route.redirect(Route.ScriptIndexPage(txHash |> Hash.toHex, Route.ScriptTransactions))
      }
    />
  </Col>;
};

[@react.component]
let make = () =>
  {
    let%Opt recentScripts = ScriptHook.getScriptList(~limit=9, ());

    Some(
      <Row>
        <Col size=2.>
          {renderFeatured(recentScripts, 0, Colors.yellow, Css.hex("333333"))}
          {renderFeatured(recentScripts, 1, Colors.orange, Css.white)}
          {renderFeatured(recentScripts, 2, Colors.pink, Css.white)}
        </Col>
        <HSpacing size=Spacing.xl />
        <Col size=5. alignSelf=Col.FlexStart>
          <Text value="Recent Data Scripts" size=Text.Xl weight=Text.Bold block=true />
          <Row wrap=true alignItems=`flexStart>
            {renderScript(recentScripts, 3)}
            {renderScript(recentScripts, 4)}
            {renderScript(recentScripts, 5)}
            {renderScript(recentScripts, 6)}
            {renderScript(recentScripts, 7)}
            {renderScript(recentScripts, 8)}
          </Row>
        </Col>
      </Row>,
    );
  }
  /*<Col>
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
    </Col>*/
  ->Belt.Option.getWithDefault(React.null);
