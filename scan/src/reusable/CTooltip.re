type placement =
  | Left
  | Right
  | Top
  | Bottom
  | BottomLeft
  | BottomRight;

module Styles = {
  open Css;
  let tooltipContainer =
    style([
      cursor(`pointer),
      position(`relative),
      display(`block),
      hover([
        selector("> div:nth-child(1)", [opacity(1.), zIndex(99), pointerEvents(`unset)]),
      ]),
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
      Media.mobile([padding(`px(12)), fontSize(`px(12))]),
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
      ])
    | BottomLeft =>
      style([
        top(`percent(100.)),
        left(`percent(50.)),
        transform(`translateX(`percent(-30.))),
        before([
          transform(`translateX(`percent(-50.))),
          left(`percent(30.)),
          top(`px(-10)),
          borderBottomColor(Colors.gray7),
        ]),
      ])
    | BottomRight =>
      style([
        top(`percent(100.)),
        left(`percent(50.)),
        transform(`translateX(`percent(-70.))),
        before([
          transform(`translateX(`percent(-70.))),
          left(`percent(70.)),
          top(`px(-10)),
          borderBottomColor(Colors.gray7),
        ]),
      ]);

  let placementSm =
    fun
    | Left =>
      style([
        Media.mobile([
          top(`percent(50.)),
          right(`percent(100.)),
          transform(`translateY(`percent(-50.))),
          before([
            top(`percent(50.)),
            transform(`translateY(`percent(-50.))),
            right(`px(-10)),
            borderLeftColor(Colors.gray7),
          ]),
        ]),
      ])
    | Right =>
      style([
        Media.mobile([
          top(`percent(50.)),
          left(`percent(100.)),
          transform(`translateY(`percent(-50.))),
          before([
            top(`percent(50.)),
            transform(`translateY(`percent(-50.))),
            left(`px(-10)),
            borderRightColor(Colors.gray7),
          ]),
        ]),
      ])
    | Top =>
      style([
        Media.mobile([
          bottom(`percent(100.)),
          left(`percent(50.)),
          transform(`translateX(`percent(-50.))),
          before([
            transform(`translateX(`percent(-50.))),
            left(`percent(50.)),
            bottom(`px(-10)),
            borderTopColor(Colors.gray7),
          ]),
        ]),
      ])
    | Bottom =>
      style([
        Media.mobile([
          top(`percent(100.)),
          left(`percent(50.)),
          transform(`translateX(`percent(-50.))),
          before([
            transform(`translateX(`percent(-50.))),
            left(`percent(50.)),
            top(`px(-10)),
            borderBottomColor(Colors.gray7),
          ]),
        ]),
      ])
    | BottomLeft =>
      style([
        Media.mobile([
          top(`percent(100.)),
          left(`percent(50.)),
          transform(`translateX(`percent(-30.))),
          before([
            transform(`translateX(`percent(-50.))),
            left(`percent(30.)),
            top(`px(-10)),
            borderBottomColor(Colors.gray7),
          ]),
        ]),
      ])
    | BottomRight =>
      style([
        Media.mobile([
          top(`percent(100.)),
          left(`percent(50.)),
          transform(`translateX(`percent(-70.))),
          before([
            transform(`translateX(`percent(-50.))),
            left(`percent(70.)),
            top(`px(-10)),
            borderBottomColor(Colors.gray7),
          ]),
        ]),
      ]);
};

[@react.component]
let make =
    (~width=150, ~tooltipText="", ~tooltipPlacement=Bottom, ~tooltipPlacementSm=Bottom, ~children) => {
  <div className=Styles.tooltipContainer>
    <div
      className={Css.merge([
        Styles.tooltipItem(width),
        Styles.placement(tooltipPlacement),
        Styles.placementSm(tooltipPlacementSm),
      ])}>
      {tooltipText |> React.string}
    </div>
    children
  </div>;
};
