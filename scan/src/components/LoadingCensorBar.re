module Styles = {
  open Css;

  let main = (~w, ~h, ~r, ~colorBase, ~colorLighter, ()) =>
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

  let alignRight = style([marginLeft(`auto)]);
  let fullWidth = style([width(`percent(100.))]);
  let mt = (~size, ()) => style([marginTop(`px(size))]);
  let mb = (~mb, ~mbSm, ()) =>
    style([marginBottom(`px(mb)), Media.mobile([marginBottom(`px(mbSm))])]);
};

[@react.component]
let make =
    (
      ~width=100,
      ~height,
      ~fullWidth=false,
      ~radius=4,
      ~colorBase=Colors.blueGray2,
      ~colorLighter=Colors.blueGray1,
      ~isRight=false,
      ~mt=0,
      ~mb=0,
      ~mbSm=mb,
    ) => {
  <div
    className={Css.merge([
      Styles.main(~w=width, ~h=height, ~r=radius, ~colorBase, ~colorLighter, ()),
      Styles.mt(~size=mt, ()),
      Styles.mb(~mb, ~mbSm, ()),
      isRight ? Styles.alignRight : "",
      fullWidth ? Styles.fullWidth : "",
    ])}
  />;
};
