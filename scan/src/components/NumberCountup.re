[@react.component]
let make = (~value, ~size, ~weight, ~spacing, ~color=Colors.gray7, ~code=true, ~smallNumber=false) => {
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
      ),
    );

  React.useEffect1(
    _ => {
      Countup.updateGet(countup, value);
      None;
    },
    [|value|],
  );
  let newVal = Countup.countUpGet(countup) |> Js.Float.toString;
  smallNumber
    ? {
      let adjustedText = newVal->Js.String2.split(".");
      <div className={CssHelper.flexBox(~align=`flexEnd, ())}>
        {adjustedText |> Belt_Array.length == 1
           ? <Text
               value={Array.get(adjustedText, 0)}
               size
               weight
               spacing
               code
               nowrap=true
               color
             />
           : <>
               <Text
                 value={Array.get(adjustedText, 0)}
                 size
                 weight
                 spacing
                 code
                 nowrap=true
                 color
               />
               <Text
                 value={"." ++ adjustedText[1]}
                 size=Text.Lg
                 weight
                 spacing
                 code
                 nowrap=true
                 color
               />
             </>}
      </div>;
    }
    : <Text value=newVal size weight spacing code nowrap=true color />;
};
