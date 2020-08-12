module Styles = {
  open Css;

  let container =
    style([
      padding2(~v=`px(40), ~h=`px(45)),
      Media.mobile([padding2(~v=`px(20), ~h=`zero)]),
    ]);

  let upperTextCotainer = style([marginBottom(`px(24))]);

  let listContainer = style([marginBottom(`px(25))]);

  let input =
    style([
      width(`percent(100.)),
      background(white),
      paddingLeft(`px(20)),
      fontSize(`px(12)),
      fontWeight(`num(500)),
      outline(`px(1), `none, white),
      height(`px(37)),
      borderRadius(`px(4)),
      border(`px(1), `solid, Colors.gray9),
      placeholder([color(Colors.blueGray3)]),
    ]);

  let button = isLoading =>
    style([
      backgroundColor(isLoading ? Colors.blueGray3 : Colors.bandBlue),
      fontWeight(`num(600)),
      color(isLoading ? Colors.blueGray7 : Colors.white),
      cursor(isLoading ? `auto : `pointer),
      outline(`zero, `none, white),
      marginTop(`px(16)),
      border(`zero, `solid, Colors.white),
    ]);

  let hFlex = style([display(`flex), flexDirection(`row)]);

  let withWH = (w, h) =>
    style([
      width(w),
      height(h),
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
    ]);

  let resultWrapper = (w, h, overflowChioce) =>
    style([
      width(w),
      height(h),
      display(`flex),
      flexDirection(`column),
      justifyContent(`center),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      overflow(overflowChioce),
    ]);
};

let parameterInput = (name, index, setCalldataList) => {
  let name = Js.String.replaceByRe([%re "/[_]/g"], " ", name);
  <div className=Styles.listContainer key=name>
    <Text
      value=name
      size=Text.Md
      color=Colors.gray7
      weight=Text.Semibold
      transform=Text.Capitalize
    />
    <VSpacing size=Spacing.sm />
    <input
      className=Styles.input
      type_="text"
      placeholder="Value"
      onChange={event => {
        let newVal = ReactEvent.Form.target(event)##value;
        setCalldataList(prev => {
          prev->Belt_List.mapWithIndex((i, value) => {index == i ? newVal : value})
        });
      }}
    />
  </div>;
};

type result_data_t = {
  returncode: int,
  stdout: string,
  stderr: string,
};

type result_t =
  | Nothing
  | Loading
  | Error(string)
  | Success(result_data_t);

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
      <div className={Styles.resultWrapper(`percent(100.), `px(90), `scroll)}>
        <Text value=err />
      </div>
    </>
  | Success({returncode, stdout, stderr}) =>
    <>
      <VSpacing size=Spacing.lg />
      <div className={Styles.resultWrapper(`percent(100.), `px(95), `auto)}>
        <div className=Styles.hFlex>
          <HSpacing size=Spacing.lg />
          <div className={Styles.resultWrapper(`px(120), `px(12), `auto)}>
            <Text value="Exit Status" color=Colors.gray6 weight=Text.Semibold />
          </div>
          <Text value={returncode |> string_of_int} />
        </div>
        <VSpacing size=Spacing.md />
        <div className=Styles.hFlex>
          <HSpacing size=Spacing.lg />
          <div className={Styles.resultWrapper(`px(120), `px(12), `auto)}>
            <Text value="Output" color=Colors.gray6 weight=Text.Semibold />
          </div>
          <Text value=stdout code=true weight=Text.Semibold />
        </div>
        <VSpacing size=Spacing.md />
        <div className=Styles.hFlex>
          <HSpacing size=Spacing.lg />
          <div className={Styles.resultWrapper(`px(120), `px(12), `auto)}>
            <Text value="Error" color=Colors.gray6 weight=Text.Semibold />
          </div>
          <Text value=stderr code=true weight=Text.Semibold />
        </div>
      </div>
    </>
  };
};

[@react.component]
let make = (~executable: JsBuffer.t) => {
  let params =
    ExecutableParser.parseExecutableScript(executable)->Belt_Option.getWithDefault([]);
  let numParams = params->Belt_List.length;

  let (callDataList, setCalldataList) = React.useState(_ => Belt_List.make(numParams, ""));

  let (result, setResult) = React.useState(_ => Nothing);

  <Row.Grid>
    <Col.Grid>
      <div className=Styles.container>
        <div className={Css.merge([CssHelper.flexBox(), Styles.upperTextCotainer])}>
          <Text
            value={
              "Test data source execution"
              ++ (numParams == 0 ? "" : " with" ++ (numParams == 1 ? " a " : " ") ++ "following")
            }
            color=Colors.gray7
            size=Text.Lg
          />
          <HSpacing size=Spacing.sm />
          {numParams == 0
             ? React.null
             : <Text
                 value={numParams > 1 ? "parameters" : "parameter"}
                 color=Colors.gray7
                 weight=Text.Bold
                 size=Text.Lg
               />}
        </div>
        {numParams > 0
           ? <>
               {params
                ->Belt_List.mapWithIndex((i, param) => parameterInput(param, i, setCalldataList))
                ->Belt_List.toArray
                ->React.array}
             </>
           : React.null}
        <div className="buttonContainer">
          <div className={CssHelper.flexBox()}>
            <Text value="Click" color=Colors.gray7 />
            <HSpacing size=Spacing.sm />
            <Text value=" Test Execution " color=Colors.gray7 weight=Text.Bold />
            <HSpacing size=Spacing.sm />
            <Text value="to test the data source." color=Colors.gray7 />
          </div>
          <button
            className={Css.merge([
              CssHelper.btn(~fsize=14, ()),
              Styles.button(result == Loading),
            ])}
            onClick={_ =>
              if (result != Loading) {
                setResult(_ => Loading);
                let _ =
                  AxiosRequest.execute(
                    AxiosRequest.t(
                      ~executable=executable->JsBuffer.toBase64,
                      ~calldata={
                        callDataList
                        ->Belt_List.reduce("", (acc, calldata) => acc ++ " " ++ calldata)
                        ->String.trim;
                      },
                      ~timeout=5000,
                    ),
                  )
                  |> Js.Promise.then_(res => {
                       setResult(_ =>
                         Success({
                           returncode: res##data##returncode,
                           stdout: res##data##stdout,
                           stderr: res##data##stderr,
                         })
                       );
                       Js.Promise.resolve();
                     })
                  |> Js.Promise.catch(err => {
                       let errorValue =
                         Js.Json.stringifyAny(err)->Belt_Option.getWithDefault("Unknown");
                       setResult(_ => Error(errorValue));
                       Js.Promise.resolve();
                     });
                ();
              }
            }>
            {(result == Loading ? "Executing ... " : "Test Execution") |> React.string}
          </button>
        </div>
        {resultRender(result)}
      </div>
    </Col.Grid>
  </Row.Grid>;
};
