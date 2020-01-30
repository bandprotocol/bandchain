module Styles = {
  open Css;

  let container = style([padding(`px(20)), background(Colors.lighterGray)]);

  let paramsContainer = style([display(`inlineBlock)]);

  let listContainer =
    style([
      display(`grid),
      gridColumnGap(`px(15)),
      gridTemplateColumns([`auto, `px(280)]),
      background(Colors.white),
      border(`px(1), `solid, Colors.lightGray),
      alignItems(`center),
    ]);

  let keyContainer =
    style([
      marginLeft(`px(25)),
      marginRight(`px(25)),
      display(`flex),
      justifyContent(`flexEnd),
    ]);

  let input =
    style([
      width(`percent(100.)),
      background(white),
      padding(Spacing.md),
      paddingLeft(`px(10)),
      fontSize(`px(14)),
      outline(`px(1), `none, white),
    ]);

  let buttonContainer = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let button =
    style([
      width(`px(110)),
      backgroundColor(Colors.btnGreen),
      borderRadius(`px(4)),
      border(`px(0), `solid, Colors.green),
      fontSize(`px(12)),
      fontWeight(`medium),
      color(`hex("127658")),
      cursor(`pointer),
      outline(`px(1), `none, white),
      padding2(~v=Css.px(10), ~h=Css.px(10)),
      whiteSpace(`nowrap),
    ]);

  let resultLink = style([cursor(`pointer)]);

  let selectPadding = style([padding(`px(10))]);
};

let parameterInput = (name, dataType, value, updateData) => {
  <div className=Styles.listContainer key=name>
    <div className=Styles.keyContainer> <Text value=name size=Text.Lg /> </div>
    {switch (dataType) {
     | "coins::Coins" =>
       <div className=Styles.selectPadding>
         <select
           value
           onChange={event => {
             let newVal = ReactEvent.Form.target(event)##value;
             updateData(name, newVal);
           }}>
           <option value=""> {"Select token" |> React.string} </option>
           {[|"ADA", "BAND", "BCH", "BNB", "BSV", "BTC", "EOS", "ETH", "LTC", "USDT", "XRP"|]
            ->Belt_Array.map(symbol => <option value=symbol> {symbol |> React.string} </option>)
            |> React.array}
         </select>
       </div>
     | _ =>
       <input
         className=Styles.input
         type_="text"
         value
         placeholder="Input Parameter here"
         onChange={event => {
           let newVal = ReactEvent.Form.target(event)##value;
           updateData(name, newVal);
         }}
       />
     }}
  </div>;
};

type result_t =
  | Nothing
  | Loading
  | Error(string)
  | Success(Hash.t);

type action =
  | DispatchSuccess(Hash.t)
  | DispatchError(string)
  | DispatchLoading;

let reducer = _state =>
  fun
  | DispatchLoading => Loading
  | DispatchError(err) => Error(err)
  | DispatchSuccess(txHash) => Success(txHash);

[@react.component]
let make = (~script: ScriptHook.Script.t) => {
  let params = script.info.params;
  let (data, setData) =
    React.useState(_ => params->Belt.List.map(({name, dataType}) => (name, dataType, "")));
  let (result, dispatch) = React.useReducer(reducer, Nothing);

  let updateData = (targetName, newVal) => {
    let newData =
      data->Belt.List.map(((name, dataType, value)) =>
        if (name == targetName) {
          (name, dataType, newVal);
        } else {
          (name, dataType, value);
        }
      );
    setData(_ => newData);
  };

  <div className=Styles.container>
    <Text value="Request Data with Parameters" color=Colors.darkGrayText size=Text.Lg />
    <VSpacing size=Spacing.md />
    <div className=Styles.paramsContainer>
      {data
       ->Belt.List.map(((name, dataType, value)) =>
           parameterInput(name, dataType, value, updateData)
         )
       ->Array.of_list
       ->React.array}
    </div>
    <VSpacing size=Spacing.md />
    <div className=Styles.buttonContainer>
      <button
        className=Styles.button
        onClick={_ => {
          DispatchLoading |> dispatch;
          let _ =
            AxiosRequest.execute(
              AxiosRequest.t(
                ~codeHash={
                  script.info.codeHash |> Hash.toHex;
                },
                ~params=
                  data->Belt_List.map(((name, _, value)) => (name, value))->Js.Dict.fromList,
              ),
            )
            |> Js.Promise.then_(res => {
                 DispatchSuccess(res##data##txHash |> Hash.fromHex) |> dispatch;
                 Js.Promise.resolve();
               })
            |> Js.Promise.catch(err => {
                 let errorValue =
                   Js.Json.stringifyAny(err)->Belt_Option.getWithDefault("Unknown");
                 DispatchError("An error occured: " ++ errorValue) |> dispatch;
                 Js.Promise.resolve();
               });
          ();
        }}>
        {"Send Request" |> React.string}
      </button>
      <HSpacing size=Spacing.xl />
      {switch (result) {
       | Nothing => React.null
       | Loading => <Text value="Loading..." />
       | Error(error) => <Text value=error color=Colors.red />
       | Success(txHash) =>
         <div
           className=Styles.resultLink onClick={_ => Route.redirect(Route.TxIndexPage(txHash))}>
           <Text value={txHash |> Hash.toHex} color=Colors.green />
         </div>
       }}
    </div>
  </div>;
};
