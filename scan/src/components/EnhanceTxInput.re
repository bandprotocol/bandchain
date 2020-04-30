module Styles = {
  open Css;

  let rowContainer =
    style([
      display(`flex),
      alignItems(`center),
      justifyContent(`spaceBetween),
      position(`relative),
    ]);

  let input = wid =>
    style([
      width(`px(wid)),
      height(`px(30)),
      paddingLeft(`px(9)),
      paddingRight(`px(9)),
      borderRadius(`px(8)),
      fontSize(`px(12)),
      textAlign(`right),
      boxShadow(
        Shadow.box(
          ~inset=false,
          ~x=`zero,
          ~y=`px(3),
          ~blur=`px(4),
          Css.rgba(11, 29, 142, 0.1),
        ),
      ),
      focus([outlineColor(Colors.white)]),
    ]);

  let errMsg = style([position(`absolute), top(`px(20))]);
};

type input_t('a) = {
  text: string,
  value: option('a),
};

type status =
  | Ok
  | Error;

let empty = {text: "", value: None};

[@react.component]
let make = (~inputData, ~setInputData, ~msg, ~errMsg, ~parse, ~width) => {
  let (status, setStatus) = React.useState(_ => Ok);

  <div className=Styles.rowContainer>
    <Text value=msg size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
    <input
      value={inputData.text}
      className={Styles.input(width)}
      onChange={event => {
        let newText = ReactEvent.Form.target(event)##value;
        let newVal = parse(newText);
        switch (newVal) {
        | Some(newVal') =>
          setStatus(_ => Ok);
          setInputData(_ => {text: newText, value: Some(newVal')});
        | None =>
          setStatus(_ => Error);
          setInputData(_ => {text: newText, value: None});
        };
      }}
    />
    {status == Error
       ? <div className=Styles.errMsg> <Text value=errMsg color=Colors.red3 size=Text.Sm /> </div>
       : React.null}
  </div>;
};
