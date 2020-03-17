module Styles = {
  open Css;

  let container = style([padding2(~h=`px(20), ~v=`px(20))]);

  let paramsContainer = style([display(`flex), flexDirection(`column)]);

  let listContainer = style([marginBottom(`px(25))]);

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
      outline(`px(0), `none, white),
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

  let resultWrapper = (w, h, paddingV, overflow_choice) =>
    style([
      width(w),
      height(h),
      display(`flex),
      flexDirection(`column),
      padding2(paddingV, `px(0)),
      justifyContent(`center),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      overflow(overflow_choice),
    ]);

  let buttonWrapper = color =>
    style([
      backgroundColor(color),
      padding2(~h=`px(8), ~v=`px(4)),
      display(`flex),
      width(`px(103)),
      height(`px(25)),
      borderRadius(`px(6)),
      cursor(`pointer),
      alignItems(`center),
      justifyContent(`center),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), rgba(20, 32, 184, 0.2))),
    ]);

  let logo = style([width(`px(15))]);
};

let parameterInput = (name, placeholder, index, setCalldataArr) => {
  <div className=Styles.listContainer key=name>
    <Text value=name size=Text.Md color=Colors.gray6 />
    <VSpacing size=Spacing.xs />
    <input
      className=Styles.input
      type_="text"
      onChange={event => {
        let newVal = ReactEvent.Form.target(event)##value;
        setCalldataArr(prev => {
          prev->Belt_Array.set(index, newVal);
          prev;
        });
      }}
      placeholder
    />
  </div>;
};

let copyButton = (~data) => {
  <div
    className={Styles.buttonWrapper(Colors.blue1)}
    onClick={_ => {Copy.copy(data |> JsBuffer.toHex(~with0x=false))}}>
    <img src=Images.copy className=Styles.logo />
    <HSpacing size=Spacing.sm />
    <Text value="Copy Proof" size=Text.Sm block=true color=Colors.bandBlue nowrap=true />
  </div>;
};

let extLinkButton = () => {
  <a href="https://twitter.com/bandprotocol" target="_blank" rel="noopener">
    <div className={Styles.buttonWrapper(Colors.gray4)}>
      <img src=Images.externalLink className=Styles.logo />
      <HSpacing size=Spacing.sm />
      <Text value="What is Proof ?" size=Text.Sm block=true color=Colors.gray7 nowrap=true />
    </div>
  </a>;
};

type result_t =
  | Nothing
  | Loading
  | Error(string)
  | Success(string);

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
      <div className={Styles.resultWrapper(`percent(100.), `px(90), `px(0), `scroll)}>
        <Text value=err />
      </div>
    </>
  | Success(output) =>
    let isFinish = Js.Math.random_int(0, 2) > 0;
    let kvs = [["Price", "866825"], ["Random", "135730902915"]];
    let proof =
      "0x0000000000000000000434000000009024900000000000b0a0df0000000fd070a00b0becd989f8989af9c80000fd070a00b0becd989f8989af"
      |> JsBuffer.fromHex;
    <>
      <VSpacing size=Spacing.lg />
      <div className={Styles.resultWrapper(`percent(100.), `auto, `px(30), `auto)}>
        <div className={Styles.hFlex(`auto)}>
          <HSpacing size=Spacing.lg />
          <div className={Styles.resultWrapper(`px(220), `px(12), `px(0), `auto)}>
            <Text value="EXIT STATUS" size=Text.Sm color=Colors.gray6 weight=Text.Semibold />
          </div>
          <Text value="0" />
        </div>
        <VSpacing size=Spacing.lg />
        <div className={Styles.hFlex(`auto)}>
          <HSpacing size=Spacing.lg />
          <div className={Styles.resultWrapper(`px(220), `px(12), `px(0), `auto)}>
            <Text value="REQUEST ID" size=Text.Sm color=Colors.gray6 weight=Text.Semibold />
          </div>
          <TypeID.Request id={ID.Request.ID(8)} />
        </div>
        <VSpacing size=Spacing.lg />
        <div className={Styles.hFlex(`auto)}>
          <HSpacing size=Spacing.lg />
          <div className={Styles.resultWrapper(`px(220), `px(12), `px(0), `auto)}>
            <Text value="TX HASH" size=Text.Sm color=Colors.gray6 weight=Text.Semibold />
          </div>
          <TxLink
            txHash={
              "D0023B6243CBBC6BC72C2543C87D55345257229868ED40C01C967A649B6F9BFD" |> Hash.fromHex
            }
            width=500
          />
        </div>
        <VSpacing size=Spacing.lg />
        {isFinish
           ? <>
               <div className={Styles.hFlex(`auto)}>
                 <HSpacing size=Spacing.lg />
                 <div className={Styles.vFlex(`px(220), `px(20 * (kvs |> Belt_List.length)))}>
                   <Text
                     value="OUTPUT"
                     size=Text.Sm
                     color=Colors.gray6
                     weight=Text.Semibold
                     height={Text.Px(20)}
                   />
                 </div>
                 <div className={Styles.vFlex(`auto, `auto)}>
                   {kvs->Belt_List.map(entry =>
                      <div className={Styles.hFlex(`px(20))}>
                        <div className={Styles.vFlex(`px(220), `auto)}>
                          <Text value={entry->Belt_List.getExn(0)} color=Colors.gray8 />
                        </div>
                        <div className={Styles.vFlex(`px(440), `auto)}>
                          <Text
                            value={entry->Belt_List.getExn(1)}
                            code=true
                            color=Colors.gray8
                            weight=Text.Bold
                          />
                        </div>
                      </div>
                    )
                    |> Belt_List.toArray
                    |> React.array}
                 </div>
               </div>
               <VSpacing size=Spacing.lg />
               <div className={Styles.hFlex(`auto)}>
                 <HSpacing size=Spacing.lg />
                 <div className={Styles.vFlex(`px(220), `auto)}>
                   <Text
                     value="PROOF OF VALIDITY"
                     size=Text.Sm
                     color=Colors.gray6
                     weight=Text.Semibold
                     height={Text.Px(15)}
                   />
                 </div>
                 <div className={Styles.vFlex(`px(660), `auto)}>
                   <Text
                     value={proof |> JsBuffer.toHex}
                     height={Text.Px(15)}
                     code=true
                     ellipsis=true
                   />
                 </div>
               </div>
               <VSpacing size=Spacing.md />
               <div className={Styles.hFlex(`auto)}>
                 <HSpacing size=Spacing.lg />
                 <div className={Styles.vFlex(`px(220), `auto)} />
                 {copyButton(proof)}
                 <HSpacing size=Spacing.lg />
                 {extLinkButton()}
               </div>
             </>
           : <div className={Styles.hFlex(`auto)}>
               <HSpacing size=Spacing.lg />
               <div className={Styles.resultWrapper(`px(220), `px(12), `px(0), `auto)}>
                 <Text
                   value="WAITING FOR OUTPUT AND PROOF"
                   size=Text.Sm
                   color=Colors.gray6
                   weight=Text.Semibold
                 />
               </div>
               <div className={Styles.resultWrapper(`px(660), `px(40), `px(0), `auto)}>
                 <ProgressBar reportedValidators=3 minimumValidators=4 totalValidators=5 />
               </div>
             </div>}
      </div>
    </>;
  };
};

[@react.component]
let make = (~code: JsBuffer.t) => {
  let params = ["Symbol", "Multiplier"]; // TODO, replace this mock by the real deal
  let numParams = params->Belt_List.length;

  let (calldataArr, setCalldataArr) = React.useState(_ => params->Belt_List.toArray);

  let (result, setResult) = React.useState(_ => Nothing);

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
           {params
            ->Belt_List.mapWithIndex((i, param) => parameterInput(param, "", i, setCalldataArr))
            ->Belt_List.toArray
            ->React.array}
         </div>
       : React.null}
    <VSpacing size=Spacing.md />
    <div className=Styles.buttonContainer>
      <button
        className={Styles.button(result == Loading)}
        onClick={_ =>
          if (result != Loading) {
            setResult(_ => Loading);
            let _ =
              AxiosRequest.request(
                AxiosRequest.t(
                  ~executable=code->JsBuffer.toHex,
                  ~calldata={
                    calldataArr
                    ->Belt_Array.reduce("", (acc, calldata) => acc ++ " " ++ calldata)
                    ->String.trim;
                  },
                ),
              )
              |> Js.Promise.then_(res => {
                   setResult(_ => Success(res##data##result));
                   Js.Promise.resolve();
                 })
              |> Js.Promise.catch(err => {
                   //  let errorValue =
                   //    Js.Json.stringifyAny(err)->Belt_Option.getWithDefault("Unknown");
                   //  setResult(_ => Error(errorValue));
                   setResult(_ => Success("test"));
                   Js.Promise.resolve();
                 });
            ();
          }
        }>
        {(result == Loading ? "Sending Request ... " : "Request") |> React.string}
      </button>
    </div>
    {resultRender(result)}
  </div>;
};
