type placement =
  | Left
  | Right
  | Top
  | Bottom;

module Styles = {
  open Css;
  let tooltipContainer =
    style([
      cursor(`pointer),
      position(`relative),
      display(`block),
      hover([selector("> div:nth-child(1)", [opacity(1.), pointerEvents(`unset)])]),
    ]);
  let tooltipItem = width_ =>
    style([
      position(`absolute),
      display(`block),
      backgroundColor(Colors.gray7),
      borderRadius(`px(4)),
      width(`px(width_)),
      color(Colors.white),
      fontSize(`px(14)),
      padding(`px(16)),
      opacity(0.),
      pointerEvents(`none),
      transition(~duration=200, "all"),
      before([
        contentRule(`text("")),
        display(`block),
        position(`absolute),
        border(`px(5), `solid, Colors.transparent),
      ]),
    ]);

  let placement =
    fun
    | Left =>
      style([
        top(`percent(50.)),
        right(`percent(100.)),
        transform(`translateY(`percent(-50.))),
        before([
          top(`percent(50.)),
          transform(`translateY(`percent(-50.))),
          right(`px(-10)),
          borderLeftColor(Colors.gray7),
        ]),
      ])
    | Right =>
      style([
        top(`percent(50.)),
        left(`percent(100.)),
        transform(`translateY(`percent(-50.))),
        before([
          top(`percent(50.)),
          transform(`translateY(`percent(-50.))),
          left(`px(-10)),
          borderRightColor(Colors.gray7),
        ]),
      ])
    | Top =>
      style([
        bottom(`percent(100.)),
        left(`percent(50.)),
        transform(`translateX(`percent(-50.))),
        before([
          transform(`translateX(`percent(-50.))),
          left(`percent(50.)),
          bottom(`px(-10)),
          borderTopColor(Colors.gray7),
        ]),
      ])
    | Bottom =>
      style([
        top(`percent(100.)),
        left(`percent(50.)),
        transform(`translateX(`percent(-50.))),
        before([
          transform(`translateX(`percent(-50.))),
          left(`percent(50.)),
          top(`px(-10)),
          borderBottomColor(Colors.gray7),
        ]),
      ]);
};

[@react.component]
let make = (~width=150, ~tooltipText="", ~tooltipPlacement=Bottom, ~children) => {
  <div className=Styles.tooltipContainer>
    <div className={Css.merge([Styles.tooltipItem(width), Styles.placement(tooltipPlacement)])}>
      {tooltipText |> React.string}
    </div>
    children
  </div>;
};
