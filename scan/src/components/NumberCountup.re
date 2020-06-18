[@react.component]
let make = (~value, ~size, ~weight, ~spacing) => {
  let countup =
    Countup.context(
      Countup.props(
        ~start=0.,
        ~end_=0.,
        ~delay=0,
        ~decimals=6,
        ~duration=1,
        ~useEasing=false,
        ~separator=",",
      ),
    );

  React.useEffect1(
    _ => {
      Countup.updateGet(countup, value);
      None;
    },
    [|value|],
  );

  let newVal = Countup.countUpGet(countup);
  <Text value={newVal |> Js.Float.toString} size weight spacing code=true nowrap=true />;
};
