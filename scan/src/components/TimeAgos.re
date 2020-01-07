let setMomentLocale: unit => unit = [%bs.raw
  {|
function() {
  const moment = require("moment");
  moment.updateLocale('en', {
    relativeTime: {
      s: x => x == 1 ? '1 second' : x + ' seconds',
    }
  });
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
          () => {setDisplayTime(_ => time->MomentRe.Moment.fromNow(~withoutSuffix=None))},
          100,
        );
      Some(() => Js.Global.clearInterval(intervalId));
    },
    [|time|],
  );

  <Text value=displayTime size=Text.Sm />;
};
