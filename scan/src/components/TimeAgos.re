[@react.component]
let make = (~time) =>
  <Text value={time->MomentRe.Moment.fromNow(~withoutSuffix=None)} size=Text.Sm />;
