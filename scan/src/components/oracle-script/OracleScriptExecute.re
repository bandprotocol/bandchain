module Styles = {
  open Css;

  let container =
    style([
      padding2(~v=`px(40), ~h=`px(45)),
      Media.mobile([padding2(~v=`px(20), ~h=`zero)]),
    ]);

  let upperTextCotainer = style([marginBottom(`px(24))]);

  let paramsContainer = style([display(`flex), flexDirection(`column)]);

  let listContainer = style([marginBottom(`px(25))]);

  let input =
    style([
      width(`percent(100.)),
      background(white),
      padding2(~v=`zero, ~h=`px(16)),
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

  let separatorLine =
    style([
      borderStyle(`none),
      backgroundColor(Colors.gray9),
      height(`px(1)),
      margin3(~top=`px(10), ~h=`zero, ~bottom=`px(20)),
    ]);
  let titleSpacing = style([marginBottom(`px(8))]);
  let mobileBlockContainer = style([padding2(~v=`px(24), ~h=`zero)]);
  let mobileBlock =
    style([
      backgroundColor(Colors.connectBG),
      borderRadius(`px(4)),
      minHeight(`px(164)),
      selector("> i", [marginBottom(`px(16))]),
    ]);
};

module ConnectPanel = {
  module Styles = {
    open Css;
    let connectContainer =
      style([backgroundColor(Colors.connectBG), borderRadius(`px(4)), padding(`px(24))]);
    let connectInnerContainer = style([width(`percent(100.)), maxWidth(`px(370))]);
  };
  [@react.component]
  let make = (~connect) => {
    <div
      className={Css.merge([Styles.connectContainer, CssHelper.flexBox(~justify=`center, ())])}>
      <div
        className={Css.merge([
          Styles.connectInnerContainer,
          CssHelper.flexBox(~justify=`spaceBetween, ()),
        ])}>
        <Icon name="fal fa-link" size=32 color=Colors.bandBlue />
        <Text value="Please connect to make request" size=Text.Lg nowrap=true block=true />
        <Button px=20 py=5 onClick={_ => {connect()}}>
          <Text value="Connect" weight=Text.Medium nowrap=true block=true />
        </Button>
      </div>
    </div>;
  };
};

let parameterInput = (Obi.{fieldName, fieldType}, index, setCalldataArr) => {
  let fieldName = Js.String.replaceByRe([%re "/[_]/g"], " ", fieldName);
  <div className=Styles.listContainer key=fieldName>
    <div className={CssHelper.flexBox()}>
      <Text value=fieldName color=Colors.gray7 weight=Text.Semibold transform=Text.Capitalize />
      <HSpacing size=Spacing.xs />
      <Text value={j|($fieldType)|j} color=Colors.gray7 weight=Text.Semibold />
    </div>
    <VSpacing size=Spacing.sm />
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

module CountInputs = {
  [@react.component]
  let make = (~askCount, ~setAskCount, ~setMinCount, ~validatorCount) => {
    <Row.Grid marginBottom=24>
      <Col.Grid col=Col.Two colSm=Col.Six>
        <div className={Css.merge([CssHelper.flexBox(), Styles.titleSpacing])}>
          <Text
            value="Ask Count"
            color=Colors.gray7
            weight=Text.Semibold
            transform=Text.Capitalize
          />
          <HSpacing size=Spacing.xs />
          <CTooltip
            tooltipPlacementSm=CTooltip.BottomLeft
            tooltipText="The number of validators that are requested to respond to this request">
            <Icon name="fal fa-info-circle" size=10 />
          </CTooltip>
        </div>
        <div className={CssHelper.selectWrapper()}>
          <select
            className=Styles.input
            onChange={event => {
              let newVal = ReactEvent.Form.target(event)##value;
              setAskCount(_ => newVal);
            }}>
            {Belt.Array.makeBy(validatorCount, i => i + 1)
             |> Belt.Array.map(_, index =>
                  <option
                    key={(index |> string_of_int) ++ "askCount"} value={index |> string_of_int}>
                    {index |> string_of_int |> React.string}
                  </option>
                )
             |> React.array}
          </select>
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Two colSm=Col.Six>
        <div className={Css.merge([CssHelper.flexBox(), Styles.titleSpacing])}>
          <Text
            value="Min Count"
            color=Colors.gray7
            weight=Text.Semibold
            transform=Text.Capitalize
          />
          <HSpacing size=Spacing.xs />
          <CTooltip
            tooltipPlacementSm=CTooltip.BottomLeft
            tooltipText="The minimum number of validators necessary for the request to proceed to the execution phase">
            <Icon name="fal fa-info-circle" size=10 />
          </CTooltip>
        </div>
        <div className={CssHelper.selectWrapper()}>
          <select
            className=Styles.input
            onChange={event => {
              let newVal = ReactEvent.Form.target(event)##value;
              setMinCount(_ => newVal);
            }}>
            {Belt.Array.makeBy(askCount |> int_of_string, i => i + 1)
             |> Belt.Array.map(_, index =>
                  <option
                    key={(index |> string_of_int) ++ "minCount"} value={index |> string_of_int}>
                    {index |> string_of_int |> React.string}
                  </option>
                )
             |> React.array}
          </select>
        </div>
      </Col.Grid>
    </Row.Grid>;
  };
};

let clientIDInput = (clientID, setClientID) => {
  <div className=Styles.listContainer>
    <div className={CssHelper.flexBox()}>
      <Text value="Client ID" color=Colors.gray7 weight=Text.Semibold transform=Text.Capitalize />
      <HSpacing size=Spacing.xs />
      <Text value="(default value is from_scan)" color=Colors.gray7 weight=Text.Semibold />
    </div>
    <VSpacing size=Spacing.sm />
    <input
      className=Styles.input
      type_="text"
      onChange={event => {
        let newVal = ReactEvent.Form.target(event)##value;
        setClientID(_ => newVal);
      }}
      value=clientID
    />
  </div>;
};

type result_t =
  | Nothing
  | Loading
  | Error(string)
  | Success(TxCreator.tx_response_t);

let loadingRender = (wDiv, wImg, h) => {
  <div className={Styles.withWH(wDiv, h)}> <Loading width=wImg height=h /> </div>;
};

let resultRender = (result, schema) => {
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
  | Success(txResponse) => <OracleScriptExecuteResponse txResponse schema />
  };
};

module ExecutionPart = {
  [@react.component]
  let make = (~id: ID.OracleScript.t, ~schema: string, ~paramsInput: array(Obi.field_key_type_t)) => {
    let isMobile = Media.isMobile();
    let (accountOpt, dispatch) = React.useContext(AccountContext.context);
    let (_, dispatchModal) = React.useContext(ModalContext.context);
    let trackingSub = TrackingSub.use();
    let connect = chainID => dispatchModal(OpenModal(Connect(chainID)));
    let numParams = paramsInput->Belt_Array.size;

    let validatorCount = ValidatorSub.countByActive(true);

    let (callDataArr, setCallDataArr) = React.useState(_ => Belt_Array.make(numParams, ""));
    let (clientID, setClientID) = React.useState(_ => "from_scan");
    let (askCount, setAskCount) = React.useState(_ => "1");
    let (minCount, setMinCount) = React.useState(_ => "1");
    let (result, setResult) = React.useState(_ => Nothing);

    // TODO: Change when input can be empty
    let isUnused = {
      let field = paramsInput->Belt_Array.getExn(0);
      field.fieldName |> Js.String.startsWith("_");
    };
    React.useEffect0(() => {
      if (isUnused) {
        setCallDataArr(_ => [|"0"|]);
      };
      None;
    });

    let requestCallback =
      React.useCallback0(requestPromise => {
        ignore(
          requestPromise
          |> Js.Promise.then_(res =>
               switch (res) {
               | TxCreator.Tx(txResponse) =>
                 setResult(_ => Success(txResponse));
                 Js.Promise.resolve();
               | _ =>
                 setResult(_ =>
                   Error("Fail to sign message, please connect with mnemonic or ledger first")
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
    isMobile
      ? <div className=Styles.mobileBlockContainer>
          <div
            className={Css.merge([
              Styles.mobileBlock,
              CssHelper.flexBox(~justify=`center, ~direction=`column, ()),
            ])}>
            <Icon name="fal fa-exclamation-circle" size=32 color=Colors.bandBlue />
            <Text
              value="Oracle request"
              color=Colors.gray7
              size=Text.Lg
              align=Text.Center
              block=true
            />
            <Text
              value="not available on mobile"
              color=Colors.gray7
              size=Text.Lg
              align=Text.Center
              block=true
            />
          </div>
        </div>
      : <Row.Grid>
          <Col.Grid>
            <div className=Styles.container>
              {isUnused
                 ? React.null
                 : <div>
                     <div className={Css.merge([CssHelper.flexBox(), Styles.upperTextCotainer])}>
                       <Text
                         value="This oracle script requires the following"
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
                     <VSpacing size=Spacing.lg />
                     <div className=Styles.paramsContainer>
                       {paramsInput
                        ->Belt_Array.mapWithIndex((i, param) =>
                            parameterInput(param, i, setCallDataArr)
                          )
                        ->React.array}
                     </div>
                   </div>}
              <div> {clientIDInput(clientID, setClientID)} </div>
              <hr className=Styles.separatorLine />
              {switch (validatorCount) {
               | Data(count) =>
                 let limitCount = count > 16 ? 16 : count;
                 <CountInputs askCount setAskCount setMinCount validatorCount=limitCount />;
               | _ => React.null
               }}
              {switch (accountOpt) {
               | Some(_) =>
                 <>
                   <Button
                     fsize=14
                     px=25
                     py=13
                     style={Styles.button(result == Loading)}
                     onClick={_ =>
                       if (result != Loading) {
                         switch (
                           Obi.encode(
                             schema,
                             "input",
                             paramsInput
                             ->Belt_Array.map(({fieldName}) => fieldName)
                             ->Belt_Array.zip(callDataArr)
                             ->Belt_Array.map(((fieldName, fieldValue)) =>
                                 Obi.{fieldName, fieldValue}
                               ),
                           )
                         ) {
                         | Some(encoded) =>
                           setResult(_ => Loading);
                           dispatch(
                             AccountContext.SendRequest({
                               oracleScriptID: id,
                               calldata: encoded,
                               callback: requestCallback,
                               askCount,
                               minCount,
                               clientID: {
                                 switch (clientID |> String.trim == "") {
                                 | false => clientID |> String.trim
                                 | true => "from_scan"
                                 };
                               },
                             }),
                           );
                           ();
                         | None =>
                           setResult(_ =>
                             Error("Encoding fail, please check each parameter's type")
                           )
                         };
                         ();
                       }
                     }>
                     {(result == Loading ? "Sending Request ... " : "Request") |> React.string}
                   </Button>
                   {resultRender(result, schema)}
                 </>
               | None =>
                 switch (trackingSub) {
                 | Data({chainID}) => <ConnectPanel connect={_ => connect(chainID)} />
                 | Error(err) =>
                   // log for err details
                   Js.Console.log(err);
                   <Text value="chain id not found" />;
                 | _ => <LoadingCensorBar fullWidth=true height=120 />
                 }
               }}
            </div>
          </Col.Grid>
        </Row.Grid>;
  };
};

[@react.component]
let make = (~id: ID.OracleScript.t, ~schema: string) =>
  {
    let%Opt paramsInput = schema->Obi.extractFields("input");
    Some(<ExecutionPart id schema paramsInput />);
  }
  |> Belt.Option.getWithDefault(
       _,
       <div className=Styles.mobileBlockContainer>
         <div
           className={Css.merge([
             Styles.mobileBlock,
             CssHelper.flexBox(~justify=`center, ~direction=`column, ()),
           ])}>
           <Icon name="fal fa-exclamation-circle" size=32 color=Colors.bandBlue />
           <Text value="Schema not found" color=Colors.gray7 />
         </div>
       </div>,
     );
