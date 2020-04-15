module Styles = {
  open Css;

  let main = (w, h, r, colorBase, colorLighter) =>
    style([
      display(`flex),
      width(`px(w)),
      height(`px(h)),
      borderRadius(`px(r)),
      backgroundColor(colorBase),
      overflow(`hidden),
      position(`relative),
      before([
        contentRule(`text("")),
        position(`absolute),
        left(`percent(-250.)),
        width(`percent(500.)),
        height(`percent(100.)),
        backgroundColor(colorBase),
        backgroundImage(
          `linearGradient((
            `deg(90.),
            [
              (`percent(0.), colorBase),
              (`percent(25.), colorBase),
              (`percent(50.), colorLighter),
              (`percent(75.), colorBase),
              (`percent(100.), colorBase),
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
let make =
    (~width, ~height, ~radius=4, ~colorBase=Colors.blueGray2, ~colorLighter=Colors.blueGray1) => {
  <div className={Styles.main(width, height, radius, colorBase, colorLighter)} />;
};
