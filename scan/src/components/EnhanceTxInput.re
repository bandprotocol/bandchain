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
      borderRadius(`px(4)),
      fontSize(`px(12)),
      textAlign(`right),
      fontSize(`px(11)),
      boxShadow(
        Shadow.box(
          ~inset=false,
          ~x=`zero,
          ~y=`px(3),
          ~blur=`px(4),
          Css.rgba(11, 29, 142, 0.1),
        ),
      ),
      boxShadow(
        Shadow.box(~inset=true, ~x=`zero, ~y=`px(1), ~blur=`px(4), Css.rgba(11, 29, 142, 0.1)),
      ),
      placeholder([color(Colors.gray5)]),
      focus([outline(`zero, `none, Colors.white)]),
    ]);

  let code =
    style([
      fontFamilies([
        `custom("IBM Plex Mono"),
        `custom("cousine"),
        `custom("sfmono-regular"),
        `custom("Consolas"),
        `custom("Menlo"),
        `custom("liberation mono"),
        `custom("ubuntu mono"),
        `custom("Courier"),
        `monospace,
      ]),
    ]);

  let errMsg = style([position(`absolute), top(`px(20))]);
};

type input_t('a) = {
  text: string,
  value: option('a),
};

type status('a) =
  | Untouched
  | Touched(Result.t('a));

let empty = {text: "", value: None};

[@react.component]
let make =
    (
      ~inputData,
      ~setInputData,
      ~msg,
      ~parse,
      ~width,
      ~code=false,
      ~placeholder="",
      ~inputType="text",
    ) => {
  let (status, setStatus) = React.useState(_ => Untouched);

  <div className=Styles.rowContainer>
    <Text value=msg size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
    <input
      value={inputData.text}
      className={Css.merge([Styles.input(width), code ? Styles.code : ""])}
      placeholder
      type_=inputType
      spellCheck=false
      onChange={event => {
        let newText = ReactEvent.Form.target(event)##value;
        let newVal = parse(newText);
        setStatus(_ => Touched(newVal));
        switch (newVal) {
        | Ok(newVal') => setInputData(_ => {text: newText, value: Some(newVal')})
        | Err(_) => setInputData(_ => {text: newText, value: None})
        };
      }}
    />
    {switch (status) {
     | Touched(Err(errMsg)) =>
       <div className=Styles.errMsg> <Text value=errMsg color=Colors.red3 size=Text.Sm /> </div>
     | _ => React.null
     }}
  </div>;
};

module Loading = {
  [@react.component]
  let make = (~msg, ~width) => {
    <div className=Styles.rowContainer>
      <Text value=msg size=Text.Lg spacing={Text.Em(0.03)} nowrap=true block=true />
      <LoadingCensorBar width height=30 isRight=true />
    </div>;
  };
};
