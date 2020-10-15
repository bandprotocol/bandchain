module Styles = {
  open Css;

  let container = style([position(`relative), paddingBottom(`px(24))]);

  let input =
    style([
      width(`percent(100.)),
      height(`px(37)),
      paddingLeft(`px(9)),
      paddingRight(`px(9)),
      borderRadius(`px(4)),
      fontSize(`px(14)),
      fontWeight(`light),
      border(`px(1), `solid, Colors.gray9),
      placeholder([color(Colors.gray5)]),
      focus([outline(`zero, `none, Colors.white)]),
      fontFamilies([
        `custom("Inter"),
        `custom("-apple-system"),
        `custom("BlinkMacSystemFont"),
        `custom("Segoe UI"),
        `custom("Roboto"),
        `custom("Oxygen"),
        `custom("Ubuntu"),
        `custom("Cantarell"),
        `custom("Fira Sans"),
        `custom("Droid Sans"),
        `custom("Helvatica Neue"),
        `custom("sans-serif"),
      ]),
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

  let errMsg = style([position(`absolute), bottom(`px(7))]);
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
      ~maxValue=?,
      ~width,
      ~code=false,
      ~placeholder="",
      ~inputType="text",
      ~autoFocus=false,
      ~id,
    ) => {
  let (status, setStatus) = React.useState(_ => Untouched);

  let onNewText = newText => {
    let newVal = parse(newText);
    setStatus(_ => Touched(newVal));
    switch (newVal) {
    | Ok(newVal') => setInputData(_ => {text: newText, value: Some(newVal')})
    | Err(_) => setInputData(_ => {text: newText, value: None})
    };
  };

  <div className=Styles.container>
    <Text value=msg size=Text.Md weight=Text.Medium nowrap=true block=true />
    <VSpacing size=Spacing.sm />
    <div className={CssHelper.flexBox(~wrap=`nowrap, ())}>
      <input
        id
        value={inputData.text}
        className={Css.merge([Styles.input, code ? Styles.code : ""])}
        placeholder
        type_=inputType
        spellCheck=false
        autoFocus
        onChange={event => {
          let newText = ReactEvent.Form.target(event)##value;
          onNewText(newText);
        }}
      />
      {switch (maxValue) {
       | Some(maxValue') =>
         <>
           <HSpacing size=Spacing.md />
           <MaxButton
             onClick={_ => onNewText(maxValue')}
             disabled={inputData.text == maxValue'}
           />
         </>
       | None => React.null
       }}
    </div>
    {switch (status) {
     | Touched(Err(errMsg)) =>
       <div className=Styles.errMsg> <Text value=errMsg size=Text.Sm color=Colors.red3 /> </div>
     | _ => React.null
     }}
  </div>;
};

module Loading = {
  [@react.component]
  let make = (~msg, ~width) => {
    <div className=Styles.container>
      <Text value=msg size=Text.Md weight=Text.Medium nowrap=true block=true />
      <VSpacing size=Spacing.sm />
      <LoadingCensorBar width height=37 />
    </div>;
  };
};

module Loading2 = {
  [@react.component]
  let make = (~msg, ~useMax=false, ~code=false, ~placeholder) => {
    <div className=Styles.container>
      <Text value=msg size=Text.Md weight=Text.Medium nowrap=true block=true />
      <VSpacing size=Spacing.sm />
      <div className={CssHelper.flexBox(~wrap=`nowrap, ())}>
        <input
          className={Css.merge([Styles.input, code ? Styles.code : ""])}
          placeholder
          disabled=true
        />
        {useMax
           ? <> <HSpacing size=Spacing.md /> <MaxButton disabled=true onClick={_ => ()} /> </>
           : React.null}
      </div>
    </div>;
  };
};
