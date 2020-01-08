let setMomentRelativeTimeThreshold: unit => unit = [%bs.raw
  {|
function() {
  const moment = require("moment");
  moment.relativeTimeRounding(Math.floor);
  moment.relativeTimeThreshold('s', 60);
  moment.relativeTimeThreshold('ss', 0);
  moment.relativeTimeThreshold('m', 60);
  moment.relativeTimeThreshold('h', 24);
  moment.relativeTimeThreshold('d', 30);
  moment.relativeTimeThreshold('M', 12);
}
  |}
];

[@react.component]
let make = (~time) => {
  let (displayTime, setDisplayTime) =
    React.useState(_ => time->MomentRe.Moment.fromNow(~withoutSuffix=None));

  React.useEffect1(
    () => {
      let intervalId =
        Js.Global.setInterval(
          () => setDisplayTime(_ => time->MomentRe.Moment.fromNow(~withoutSuffix=None)),
          100,
        );
      Some(() => Js.Global.clearInterval(intervalId));
    },
    [|time|],
  );

  <Text value=displayTime size=Text.Sm />;
};
