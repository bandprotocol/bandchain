module Styles = {
  open Css;

  let container =
    style([
      backgroundColor(Colors.white),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      Media.mobile([margin2(~h=`px(-12), ~v=`zero)]),
    ]);
  let header =
    style([
      backgroundColor(Colors.white),
      padding2(~v=`zero, ~h=`px(24)),
      borderBottom(`px(1), `solid, Colors.gray4),
      selector("> a + a", [marginLeft(`px(40))]),
      Media.mobile([overflow(`auto), padding2(~v=`px(1), ~h=`px(15))]),
    ]);

  let buttonContainer = active =>
    style([
      height(`px(40)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      cursor(`pointer),
      padding3(~top=`px(24), ~h=`zero, ~bottom=`px(20)),
      borderBottom(`px(4), `solid, active ? Colors.bandBlue : Colors.white),
      Media.mobile([whiteSpace(`nowrap)]),
    ]);

  let childrenContainer =
    style([
      backgroundColor(Colors.blueGray1),
      Media.mobile([padding2(~h=`px(16), ~v=`zero)]),
    ]);
};

let button = (~name, ~route, ~active) => {
  <Link key=name isTab=true className={Styles.buttonContainer(active)} route>
    <Text
      value=name
      weight={active ? Text.Semibold : Text.Regular}
      size=Text.Lg
      color=Colors.gray6
    />
  </Link>;
};

type t = {
  name: string,
  route: Route.t,
};

[@react.component]
let make = (~tabs: array(t), ~currentRoute, ~children) => {
  <div className=Styles.container>
    <div className={Css.merge([Styles.header, CssHelper.flexBox(~wrap=`nowrap, ())])}>
      {tabs
       ->Belt.Array.map(({name, route}) => button(~name, ~route, ~active=route == currentRoute))
       ->React.array}
    </div>
    <div className=Styles.childrenContainer> children </div>
  </div>;
};
