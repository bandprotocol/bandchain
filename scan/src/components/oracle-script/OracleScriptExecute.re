module Styles = {
  open Css;

  let container = style([padding2(~h=`px(20), ~v=`px(20))]);

  let paramsContainer = style([display(`flex), flexDirection(`column)]);

  let listContainer = style([marginBottom(`px(25))]);

  let withPadding = (h, v) => style([padding2(~h=`px(h), ~v=`px(v))]);

  let input =
    style([
      width(`percent(100.)),
      background(white),
      paddingLeft(`px(20)),
      fontSize(`px(12)),
      fontWeight(`num(500)),
      outline(`px(1), `none, white),
      height(`px(40)),
      borderRadius(`px(4)),
      boxShadow(
        Shadow.box(~inset=true, ~x=`zero, ~y=`zero, ~blur=`px(4), Css.rgba(0, 0, 0, 0.1)),
      ),
    ]);

  let buttonContainer = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let button = isLoading =>
    style([
      width(`px(isLoading ? 150 : 110)),
      backgroundColor(isLoading ? Colors.blueGray3 : Colors.green2),
      borderRadius(`px(6)),
      fontSize(`px(12)),
      fontWeight(`num(600)),
      color(isLoading ? Colors.blueGray7 : Colors.green7),
      cursor(isLoading ? `auto : `pointer),
      padding2(~v=Css.px(10), ~h=Css.px(10)),
      whiteSpace(`nowrap),
      outline(`zero, `none, white),
      boxShadow(
        isLoading
          ? `none : Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.1)),
      ),
      border(`zero, `solid, Colors.white),
    ]);

  let hFlex = h =>
    style([display(`flex), flexDirection(`row), alignItems(`center), height(h)]);

  let vFlex = (w, h) => style([display(`flex), flexDirection(`column), width(w), height(h)]);

  let withWH = (w, h) =>
    style([
      width(w),
      height(h),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
    ]);

  let resultWrapper = (w, h, paddingV, overflowChioce) =>
    style([
      width(w),
      height(h),
      display(`flex),
      flexDirection(`column),
      padding2(~v=paddingV, ~h=`zero),
      justifyContent(`center),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      overflow(overflowChioce),
    ]);

  let logo = style([width(`px(15))]);
};

let parameterInput = (Borsh.{fieldName, fieldType}, index, setCalldataArr) => {
  <div className=Styles.listContainer key=fieldName>
    <Text value={j|$fieldName ($fieldType)|j} size=Text.Md color=Colors.gray6 />
    <VSpacing size=Spacing.xs />
    <input
      className=Styles.input
      type_="text"
      onChange={event => {
        let newVal = ReactEvent.Form.target(event)##value;
        setCalldataArr(prev => {
          prev->Belt_Array.mapWithIndex((i, value) => {index == i ? newVal : value})
        });
      }}
    />
  </div>;
};

type result_t =
  | Nothing
  | Loading
  | Error(string)
  | Success(BandWeb3.tx_response_t, string);

let loadingRender = (wDiv, wImg, h) => {
  <div className={Styles.withWH(wDiv, h)}>
    <img src=Images.loadingCircles className={Styles.withWH(wImg, h)} />
  </div>;
};

let resultRender = result => {
  switch (result) {
  | Nothing => React.null
  | Loading =>
    <>
      <VSpacing size=Spacing.xl />
      {loadingRender(`percent(100.), `px(104), `px(30))}
      <VSpacing size=Spacing.lg />
    </>
  | Error(err) =>
    <>
      <VSpacing size=Spacing.lg />
      <div className={Styles.resultWrapper(`percent(100.), `px(90), `zero, `scroll)}>
        <Text value=err />
      </div>
    </>
  | Success(txResponse, schema) => <OracleScriptExecuteResponse txResponse schema />
  };
};

module ExecutionPart = {
  [@react.component]
  let make =
      (~id: ID.OracleScript.t, ~schema: string, ~paramsInput: array(Borsh.field_key_type_t)) => {
    let (_, dispatch) = React.useContext(AccountContext.context);

    let numParams = paramsInput->Belt_Array.size;

    let (callDataArr, setCallDataArr) = React.useState(_ => Belt_Array.make(numParams, ""));
    let (result, setResult) = React.useState(_ => Nothing);

    let requestCallback =
      React.useCallback0(requestPromise => {
        ignore(
          requestPromise
          |> Js.Promise.then_(res =>
               switch (res) {
               | BandWeb3.Tx(txResponse) =>
                 setResult(_ => Success(txResponse, schema));
                 Js.Promise.resolve();
               | _ =>
                 setResult(_ =>
                   Error("Fail to sign message, please connect with mnemonic first")
                 );
                 Js.Promise.resolve();
               }
             )
          |> Js.Promise.catch(err => {
               switch (Js.Json.stringifyAny(err)) {
               | Some(errorValue) => setResult(_ => Error(errorValue))
               | None => setResult(_ => Error("Can not stringify error"))
               };
               Js.Promise.resolve();
             }),
        );
        ();
      });

    <div className=Styles.container>
      <div className={Styles.hFlex(`auto)}>
        <Text
          value={
            "Request"
            ++ (numParams == 0 ? "" : " with" ++ (numParams == 1 ? " a " : " ") ++ "following")
          }
          color=Colors.gray7
        />
        <HSpacing size=Spacing.sm />
        {numParams == 0
           ? React.null
           : <Text
               value={numParams > 1 ? "parameters" : "parameter"}
               color=Colors.gray7
               weight=Text.Bold
             />}
      </div>
      <VSpacing size=Spacing.lg />
      {numParams > 0
         ? <div className=Styles.paramsContainer>
             {paramsInput
              ->Belt_Array.mapWithIndex((i, param) => parameterInput(param, i, setCallDataArr))
              ->React.array}
           </div>
         : React.null}
      <VSpacing size=Spacing.md />
      <div className=Styles.buttonContainer>
        <button
          className={Styles.button(result == Loading)}
          onClick={_ =>
            if (result != Loading) {
              switch (
                Borsh.encode(
                  schema,
                  "Input",
                  paramsInput
                  ->Belt_Array.map(({fieldName}) => fieldName)
                  ->Belt_Array.zip(callDataArr)
                  ->Belt_Array.map(((fieldName, fieldValue)) => Borsh.{fieldName, fieldValue}),
                )
              ) {
              | Some(encoded) =>
                setResult(_ => Loading);
                dispatch(AccountContext.SendRequest(id, encoded, requestCallback));
                ();
              | None => setResult(_ => Error("Encoding fail, please check each parameter's type"))
              };
              ();
            }
          }>
          {(result == Loading ? "Sending Request ... " : "Request") |> React.string}
        </button>
      </div>
      {resultRender(result)}
    </div>;
  };
};

[@react.component]
let make = (~id: ID.OracleScript.t, ~schemaOpt: option(string)) => {
  switch (
    {
      let%Opt schema = schemaOpt;
      let%Opt paramsInput = schema->Borsh.extractFields("Input");
      Some(<ExecutionPart id schema paramsInput />);
    }
  ) {
  | Some(dom) => dom
  | None =>
    <div className={Styles.withPadding(20, 20)}>
      <Text value="Schema not found" color=Colors.gray7 />
    </div>
  };
};
