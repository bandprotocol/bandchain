module Styles = {
  open Css;

  let tableLowerContainer = style([position(`relative)]);
  let tableWrapper =
    style([padding(`px(24)), Media.mobile([padding2(~v=`px(20), ~h=`zero)])]);

  let codeImage = style([width(`px(20)), marginRight(`px(10))]);
  let titleSpacing = style([marginBottom(`px(8))]);
  let scriptContainer =
    style([
      fontSize(`px(12)),
      lineHeight(`px(20)),
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

  let padding = style([padding(`px(20))]);

  let selectWrapper =
    style([
      backgroundColor(Colors.white),
      border(`px(1), `solid, Colors.gray9),
      borderRadius(`px(4)),
      display(`flex),
      flexDirection(`row),
      padding2(~v=`px(10), ~h=`px(16)),
      width(`percent(100.)),
      maxWidth(`px(200)),
      minHeight(`px(37)),
      Media.mobile([padding2(~v=`px(10), ~h=`px(8))]),
    ]);

  let selectContent =
    style([
      backgroundColor(Colors.transparent),
      borderStyle(`none),
      color(Colors.gray7),
      width(`percent(100.)),
      outlineStyle(`none),
    ]);

  let iconBody = style([width(`px(20)), height(`px(20))]);
  let copyContainer =
    style([
      position(`absolute),
      top(`px(10)),
      right(`px(10)),
      zIndex(2),
      Media.mobile([right(`zero), top(`px(-35))]),
    ]);
};

let renderCode = content => {
  <div className=Styles.scriptContainer>
    <ReactHighlight className=Styles.padding> {content |> React.string} </ReactHighlight>
  </div>;
};

type target_platform_t =
  | Ethereum
  | CosmosIBC
  | Web;
/*   | Kadena; */

type language_t =
  | Solidity
  /*   | Vyper */
  | Go
  | NodeJS;
/*   | PACT; */

exception WrongLanugageChoice(string);
exception WrongPlatformChoice(string);

let toLanguageVariant =
  fun
  | "Solidity" => Solidity
  /* | "Vyper" => Vyper */
  | "Go" => Go
  | "NodeJS" => NodeJS
  /*   | "PACT" => PACT */
  | _ => raise(WrongLanugageChoice("Chosen language does not exist"));

let toPlatformVariant =
  fun
  | "Ethereum" => Ethereum
  | "Cosmos IBC" => CosmosIBC
  | "Web" => Web
  /*   | "Kadena" => Kadena */
  | _ => raise(WrongPlatformChoice("Chosen platform does not exist"));

let toLanguageString =
  fun
  | Solidity => "Solidity"
  /*   | Vyper => "Vyper" */
  | Go => "Go"
  | NodeJS => "NodeJS";
/*  | PACT => "PACT"; */

let toPlatformString =
  fun
  | Ethereum => "Ethereum"
  | CosmosIBC => "Cosmos IBC"
  | Web => "Web";
/*   | Kadena => "Kadena"; */

let getLanguagesByPlatform =
  fun
  //TODO: Add back Vyper
  | Ethereum => [|Solidity|]
  | CosmosIBC => [|Go|]
  | Web => [|NodeJS|];
/*   | Kadena => [|PACT|]; */

module TargetPlatformIcon = {
  [@react.component]
  let make = (~icon) => {
    <div className={CssHelper.flexBox(~justify=`center, ())}>
      <img
        className=Styles.iconBody
        src={
          switch (icon) {
          | Ethereum => Images.ethereumIcon
          | CosmosIBC => Images.cosmosIBCIcon
          | Web => Images.webIcon
          /*           | Kadena => Images.kadenaIcon */
          }
        }
      />
    </div>;
  };
};

module LanguageIcon = {
  [@react.component]
  let make = (~icon) => {
    <div className={CssHelper.flexBox(~justify=`center, ())}>
      <img
        className=Styles.iconBody
        src={
          switch (icon) {
          | Solidity => Images.solidityIcon
          /*           | Vyper => Images.vyperIcon */
          | Go => Images.golangIcon
          | NodeJS => Images.nodeJSIcon
          /*           | PACT => Images.pactIcon */
          }
        }
      />
    </div>;
  };
};

let getFileNameFromLanguage = (~language, ~dataType) => {
  let dataTypeString = dataType |> Obi.dataTypeToString;
  switch (language) {
  | Solidity => "Decoders.sol"
  | Go => {j|$(dataTypeString)Decoder.go|j}
  | NodeJS => "index.js"
  };
};

let getCodeFromSchema = (~schema, ~oracleScriptID, ~language, ~dataType) => {
  switch (language) {
  | Solidity => Obi.generateDecoderSolidity(schema)
  | Go => Obi.generateDecoderGo("main", schema, dataType)
  | NodeJS => Obi.generateNodeJS(oracleScriptID, schema, dataType)
  };
};

module GenerateDecodeCode = {
  [@react.component]
  let make = (~language, ~oracleScriptID, ~schema, ~dataType) => {
    let codeOpt = getCodeFromSchema(~schema, ~oracleScriptID, ~language, ~dataType);
    let code =
      switch (codeOpt) {
      | Some(code) => code
      | _ => "Code is not available."
      };
    <>
      <Row.Grid marginBottom=24 marginTop=24 marginTopSm=12 marginBottomSm=12>
        <Col.Grid>
          <div className={CssHelper.flexBox()}>
            <Icon name="fal fa-file" size=16 />
            <HSpacing size=Spacing.sm />
            <Text
              value={getFileNameFromLanguage(~language, ~dataType)}
              weight=Text.Semibold
              size=Text.Lg
              block=true
              color=Colors.gray7
            />
          </div>
        </Col.Grid>
      </Row.Grid>
      <div className=Styles.tableLowerContainer>
        <div className=Styles.copyContainer> <CopyButton data=code title="Copy Code" /> </div>
        {code |> renderCode}
      </div>
    </>;
  };
};

[@react.component]
let make = (~oracleScriptID, ~schema) => {
  let (targetPlatform, setTargetPlatform) = React.useState(_ => Ethereum);
  let (language, setLanguage) = React.useState(_ => Solidity);
  <div className=Styles.tableWrapper>
    <Row.Grid marginBottom=24>
      <Col.Grid col=Col.Three colSm=Col.Six>
        <div className={Css.merge([CssHelper.flexBox(), Styles.titleSpacing])}>
          <Heading size=Heading.H5 value="Target Platform" />
          <HSpacing size=Spacing.xs />
          <CTooltip
            tooltipPlacementSm=CTooltip.BottomLeft
            tooltipText="The target platform to which to generate the code for">
            <Icon name="fal fa-info-circle" size=10 />
          </CTooltip>
        </div>
        <div className=Styles.selectWrapper>
          <TargetPlatformIcon icon=targetPlatform />
          <div className={CssHelper.selectWrapper(~pRight=0, ())}>
            <select
              className=Styles.selectContent
              onChange={event => {
                let newPlatform = ReactEvent.Form.target(event)##value |> toPlatformVariant;
                setTargetPlatform(_ => newPlatform);
                let newLanguage =
                  newPlatform |> getLanguagesByPlatform |> Belt_Array.getExn(_, 0);
                setLanguage(_ => newLanguage);
              }}>
              {[|Ethereum, Web|]
               ->Belt_Array.map(symbol =>
                   <option value={symbol |> toPlatformString}>
                     {symbol |> toPlatformString |> React.string}
                   </option>
                 )
               |> React.array}
            </select>
          </div>
        </div>
      </Col.Grid>
      <Col.Grid col=Col.Three colSm=Col.Six>
        <div className={Css.merge([CssHelper.flexBox(), Styles.titleSpacing])}>
          <Heading size=Heading.H5 value="Language" />
          <HSpacing size=Spacing.xs />
          <CTooltip tooltipText="The programming language">
            <Icon name="fal fa-info-circle" size=10 />
          </CTooltip>
        </div>
        <div className=Styles.selectWrapper>
          <LanguageIcon icon=language />
          <div className={CssHelper.selectWrapper(~pRight=0, ())}>
            <select
              className=Styles.selectContent
              onChange={event => {
                let newLanguage = ReactEvent.Form.target(event)##value |> toLanguageVariant;
                setLanguage(_ => newLanguage);
              }}>
              {targetPlatform
               |> getLanguagesByPlatform
               |> Belt.Array.map(_, symbol =>
                    <option value={symbol |> toLanguageString}>
                      {symbol |> toLanguageString |> React.string}
                    </option>
                  )
               |> React.array}
            </select>
          </div>
        </div>
      </Col.Grid>
    </Row.Grid>
    <Row.Grid marginBottom=24 marginBottomSm=12>
      <Col.Grid>
        <div className={CssHelper.flexBox()}>
          <Icon name="fal fa-file" size=16 />
          <HSpacing size=Spacing.sm />
          <Text
            value="Oracle Script Schema"
            weight=Text.Semibold
            size=Text.Lg
            block=true
            color=Colors.gray7
          />
        </div>
      </Col.Grid>
    </Row.Grid>
    <div className=Styles.tableLowerContainer>
      <div className=Styles.copyContainer> <CopyButton data=schema title="Copy Code" /> </div>
      {schema |> renderCode}
    </div>
    <GenerateDecodeCode language oracleScriptID schema dataType=Obi.Params />
  </div>;
};
