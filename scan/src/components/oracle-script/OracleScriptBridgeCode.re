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
      overflow(`hidden),
    ]);

  let selectContent =
    style([
      background(rgba(255, 255, 255, 1.)),
      border(`px(0), `solid, hex("FFFFFF")),
      width(`px(169)),
      float(`right),
    ]);

  let languageWrapper = style([overflow(`hidden)]);
};

let renderCode = content => {
  <div className=Styles.scriptContainer>
    <ReactHighlight>
      <div className=Styles.padding> {content |> React.string} </div>
    </ReactHighlight>
  </div>;
};

[@react.component]
let make = () => {
  let description = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent aliquet tempor imperdiet. Morbi tincidunt molestie tortor a finibus. Nulla hendrerit iaculis metus, in laoreet tellus eleifend vel. Aliquam pretium porta mi, a efficitur justo ullamcorper sed. Donec interdum accumsan nibh, sed tempor lectus rutrum ac. Morbi et magna in magna varius iaculis. Praesent mollis nulla non arcu ullamcorper, at bibendum nibh pellentesque. Aenean ac quam eget turpis euismod lacinia. Phasellus libero lectus, pulvinar at ipsum non, ullamcorper commodo felis.";
  let codetest = {j|
    pragma solidity ^0.5.0;\n
    import "./Borsch.sol";\n
    library ResultDecoder {\n\t
      using Borsh for Borsh.Data;\n
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
  <div className=Styles.tableWrapper>
    <>
      <VSpacing size={`px(10)} />
      <Row>
        <HSpacing size={`px(15)} />
        <Col>
          <div> <Text value="Target Platform" /> </div>
          <VSpacing size={`px(5)} />
          <div className=Styles.selectWrapper>
            <select className=Styles.selectContent>
              <option value=""> {"Ethereum" |> React.string} </option>
              {[|"Cosmos IBC", "Kadena"|]
               ->Belt_Array.map(symbol => <option value=symbol> {symbol |> React.string} </option>)
               |> React.array}
            </select>
          </div>
        </Col>
        <HSpacing size={`px(370)} />
        <Col>
          <div>
            <Text value="Language" />
            <div className=Styles.selectWrapper>
              <select className=Styles.selectContent>
                <option value=""> {"Solidity" |> React.string} </option>
                {[|"Vyper"|]
                 ->Belt_Array.map(symbol =>
                     <option value=symbol> {symbol |> React.string} </option>
                   )
                 |> React.array}
              </select>
            </div>
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
