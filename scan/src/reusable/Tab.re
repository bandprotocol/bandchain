module Styles = {
  open Css;

  let container =
    style([
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(10), Css.rgba(0, 0, 0, 0.08))),
    ]);
  let header =
    style([
      backgroundColor(Colors.white),
      padding2(~v=`zero, ~h=Spacing.lg),
      borderBottom(`px(1), `solid, Colors.lightGray),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(10), Css.rgba(0, 0, 0, 0.08))),
    ]);

  let buttonContainer = active =>
    style([
      height(`px(40)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      cursor(`pointer),
      padding2(~v=Spacing.md, ~h=Spacing.lg),
      borderBottom(`pxFloat(1.5), `solid, active ? Colors.brightBlue : Colors.white),
      textShadow(Shadow.text(~blur=`pxFloat(active ? 1. : 0.), Colors.brightBlue)),
    ]);
};

let button = (~name, ~route, ~active) => {
  <div key=name className={Styles.buttonContainer(active)} onClick={_ => route |> Route.redirect}>
    <Text
      value=name
      weight=Text.Regular
      size=Text.Md
      color={active ? Colors.brightBlue : Colors.darkGrayText}
    />
  </div>;
};

type t = {
  name: string,
  route: Route.t,
};

[@react.component]
let make = (~tabs: array(t), ~currentRoute, ~children) => {
  <div className=Styles.container>
    <div className=Styles.header>
      <Row>
        {tabs
         ->Belt.Array.map(({name, route}) =>
             button(~name, ~route, ~active=route == currentRoute)
           )
         ->React.array}
      </Row>
    </div>
    children
  </div>;
};
