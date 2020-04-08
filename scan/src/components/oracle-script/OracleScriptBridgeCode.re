module Styles = {
  open Css;

  let tableLowerContainer = style([padding(`px(8))]);
  let tableWrapper = style([padding2(~v=`px(20), ~h=`px(15))]);
  let codeImage = style([width(`px(20)), marginRight(`px(10))]);
  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

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
      display(`flex),
      flexDirection(`row),
      padding2(~v=`px(3), ~h=`px(8)),
      position(`static),
      width(`px(169)),
      height(`px(30)),
      left(`px(0)),
      top(`px(32)),
      background(rgba(255, 255, 255, 1.)),
      borderRadius(`px(100)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(4), rgba(0, 0, 0, 0.1))),
    ]);

  let selectContent =
    style([
      background(rgba(255, 255, 255, 1.)),
      border(`px(0), `solid, hex("FFFFFF")),
      width(`px(169)),
      float(`right),
    ]);

  let iconWrapper = style([alignItems(`center), justifyContent(`center)]);

  let iconBody = style([width(`px(20))]);
};

let renderCode = content => {
  <div className=Styles.scriptContainer>
    <ReactHighlight>
      <div className=Styles.padding> {content |> React.string} </div>
    </ReactHighlight>
  </div>;
};

module TargetPlatformIcon = {
  [@react.component]
  let make = (~icon) => {
    <div className=Styles.iconWrapper>
      <img
        className=Styles.iconBody
        src={
          switch (icon) {
          | "Ethereum" => Images.ethereumIcon
          | "Cosmos IBC" => Images.cosmosIBCIcon
          | "Kadena" => Images.kadenaIcon
          | _ => Images.missingIcon
          }
        }
      />
    </div>;
  };
};

module LanguageIcon = {
  [@react.component]
  let make = (~icon) => {
    <div className=Styles.iconWrapper>
      <img
        className=Styles.iconBody
        src={
          switch (icon) {
          | "Solidity" => Images.solidityIcon
          | "Vyper" => Images.vyperIcon
          | _ => Images.missingIcon
          }
        }
      />
    </div>;
  };
};

[@react.component]
let make = () => {
  let description = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent aliquet tempor imperdiet. Morbi tincidunt molestie tortor a finibus. Nulla hendrerit iaculis metus, in laoreet tellus eleifend vel. Aliquam pretium porta mi, a efficitur justo ullamcorper sed. Donec interdum accumsan nibh, sed tempor lectus rutrum ac. Morbi et magna in magna varius iaculis. Praesent mollis nulla non arcu ullamcorper, at bibendum nibh pellentesque. Aenean ac quam eget turpis euismod lacinia. Phasellus libero lectus, pulvinar at ipsum non, ullamcorper commodo felis.";
  let codetest = {j|
    pragma solidity ^0.5.0;

    import "./Borsch.sol";

    library ResultDecoder {
      using Borsh for Borsh.Data;

      struct Result {
        string symbol;
        uint64 multiplier;
        uint8 what;
      }

      function decodeResult(bytes memory _data)
        internal
        pure
        returns (Result memory result)
      {
          Borsh.Data memory data = Borsh.from(_data);
          result.symbol = string(data.decodeBytes());
          result.multiplier = data.decodeU64();
          result.what = data.decodeU8();
      }
    }|j};

  let (targetPlatform, setTargetPlatform) = React.useState(_ => "Ethereum");
  let (language, setLanguage) = React.useState(_ => "Solidity");
  <div className=Styles.tableWrapper>
    <>
      <VSpacing size={`px(10)} />
      <Row>
        <HSpacing size={`px(15)} />
        <Col>
          <div> <Text value="Target Platform" /> </div>
          <VSpacing size={`px(5)} />
          <div className=Styles.selectWrapper>
            <TargetPlatformIcon icon=targetPlatform />
            <select
              className=Styles.selectContent
              onChange={event => {
                let newValue = ReactEvent.Form.target(event)##value;
                setTargetPlatform(_ => newValue);
              }}>
              {[|"Ethereum", "Cosmos IBC", "Kadena"|]
               ->Belt_Array.map(symbol => <option value=symbol> {symbol |> React.string} </option>)
               |> React.array}
            </select>
          </div>
        </Col>
        <HSpacing size={`px(370)} />
        <Col>
          <div> <Text value="Language" /> </div>
          <div className=Styles.selectWrapper>
            <LanguageIcon icon=language />
            <select
              className=Styles.selectContent
              onChange={event => {
                let newValue = ReactEvent.Form.target(event)##value;
                setLanguage(_ => newValue);
              }}>
              {[|"Solidity", "Vyper"|]
               ->Belt_Array.map(symbol => <option value=symbol> {symbol |> React.string} </option>)
               |> React.array}
            </select>
          </div>
          <VSpacing size={`px(5)} />
        </Col>
      </Row>
      <VSpacing size={`px(35)} />
      <div className=Styles.tableLowerContainer>
        <div className=Styles.vFlex>
          <Text value="Description" size=Text.Lg color=Colors.gray7 weight=Text.Medium />
        </div>
        <VSpacing size=Spacing.lg />
        <Text value=description size=Text.Lg />
      </div>
      <VSpacing size={`px(35)} />
      <div className=Styles.tableLowerContainer>
        <div className=Styles.vFlex>
          <img src=Images.code className=Styles.codeImage />
          <Text value="ResultDecoder.sol" size=Text.Lg color=Colors.gray7 />
        </div>
        <VSpacing size=Spacing.lg />
        codetest->renderCode
      </div>
    </>
  </div>;
};
