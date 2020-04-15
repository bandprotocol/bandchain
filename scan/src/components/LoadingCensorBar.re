module Styles = {
  open Css;

  let main = (w, h, r) =>
    style([
      display(`flex),
      width(`px(w)),
      height(`px(h)),
      borderRadius(`px(r)),
      backgroundColor(Colors.blueGray2),
      overflow(`hidden),
      position(`relative),
      before([
        contentRule(`text("")),
        position(`absolute),
        left(`percent(-250.)),
        width(`percent(500.)),
        height(`percent(100.)),
        backgroundColor(Colors.blueGray2),
        backgroundImage(
          `linearGradient((
            `deg(90.),
            [
              (`percent(0.), Colors.blueGray2),
              (`percent(25.), Colors.blueGray2),
              (`percent(50.), Colors.blueGray1),
              (`percent(75.), Colors.blueGray2),
              (`percent(100.), Colors.blueGray2),
            ],
          )),
        ),
        animation(
          ~duration=1000,
          ~timingFunction=`linear,
          ~iterationCount=`infinite,
          keyframes([
            (0, [transform(translateX(`percent(-30.)))]),
            (100, [transform(translateX(`percent(50.)))]),
          ]),
        ),
      ]),
    ]);
};

[@react.component]
let make = (~width, ~height, ~radius=4) => {
  <div className={Styles.main(width, height, radius)} />;
};
