module Styles = {
  open Css;

  let chartContainer =
    style([
      position(`relative),
      width(`px(208)),
      height(`px(208)),
      Media.mobile([width(`px(140)), height(`px(140))]),
    ]);

  let textContainer =
    style([
      position(`absolute),
      backgroundColor(Colors.white),
      top(`px(8)),
      left(`px(8)),
      right(`px(8)),
      bottom(`px(8)),
      borderRadius(`percent(50.)),
    ]);
  let circle = percent => {
    style([
      width(`percent(100.)),
      height(`percent(100.)),
      backgroundColor(Colors.profileBG),
      borderRadius(`percent(50.)),
      selector(
        "> circle",
        [
          SVG.fill(Colors.profileBG),
          SVG.strokeWidth(`px(16)),
          SVG.stroke(Colors.bandBlue),
          //TODO: it will be remove when the bs-css upgrade to have this proporty
          // 653.45 is from 2 * pi(3.141) * r(104)
          unsafe("stroke-dasharray", {j|calc($percent * 653.45 / 100) 653.45|j}),
          transforms([`rotate(`deg(-90.)), `translateX(`percent(-100.))]),
          Media.mobile([width(`px(140)), height(`px(140))]),
        ],
      ),
    ]);
  };
};

[@react.component]
let make = (~percent) => {
  <div className=Styles.chartContainer>
    <svg
      className={Styles.circle(percent |> int_of_float |> string_of_int)} viewBox="0 0 208 208">
      <circle r="104" cx="104" cy="104" />
    </svg>
    <div
      className={Css.merge([
        Styles.textContainer,
        CssHelper.flexBox(~justify=`center, ~direction=`column, ()),
      ])}>
      <Heading size=Heading.H4 value="Turnout" align=Heading.Center marginBottom=8 />
      <Text size=Text.Xxxl value={percent |> Format.fPercent(~digits=2)} block=true />
    </div>
  </div>;
};
