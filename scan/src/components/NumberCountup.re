[@react.component]
let make = (~value, ~size, ~weight, ~spacing, ~color=Colors.gray7, ~code=true, ~smallNumber=false) => {
  Js.log2("value", value);
  let countup =
    Countup.context(
      Countup.props(
        ~start=value,
        ~end_=value,
        ~delay=0,
        ~decimals=6,
        ~duration=4,
        ~useEasing=false,
        ~separator=",",
        ~redraw=true
      ),
    );

  React.useEffect1(
    _ => {
      // Js.log(countup);
      // Js.log(value);
      Countup.updateGet(countup, value);
      if (value==0.) {
        // countup.reset();
        let a = [%bs.raw {|
          function() {
            console.log(countup)
            countup.resetCountUp();
            console.log('fark');
            // countup.reset();
            // countup.update(999);
          } 
        |}];
        a();
      }
      None;
    },
    [|value|],
  );
  let newVal = Countup.countUpGet(countup) |> Js.Float.toString;
  smallNumber
    ? {
      let adjustedText = newVal->Js.String2.split(".");
      <div className={CssHelper.flexBox(~align=`flexEnd, ())}>
        <Text
          value={adjustedText->Belt.Array.get(0)->Belt.Option.getWithDefault("0")}
          size
          weight
          spacing
          code
          nowrap=true
          color
        />
        <Text
          value={"." ++ adjustedText->Belt.Array.get(1)->Belt.Option.getWithDefault("0")}
          size=Text.Lg
          weight
          spacing
          code
          nowrap=true
          color
        />
      </div>;
    }
    : <Text value=newVal size weight spacing code nowrap=true color />;
};
