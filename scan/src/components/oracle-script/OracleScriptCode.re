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
  let code1 = {f|[package]
name = "crypto_price"
version = "0.1.0"
authors = ["Band Protocol <dev@bandprotocol.com>"]
edition = "2018"

# See more keys and their definitions at https://doc.rust-lang.org/cargo/reference/manifest.html

[lib]
crate-type = ["cdylib"]

[dependencies]
owasm = { path = "../.." }
|f};
  let code2 = {f|/* use owasm::oei;

fn parse_coingecko_symbol(symbol: &[u8]) -> &[u8] {
  let s = String::from_utf8(symbol.to_vec()).unwrap();
    (match s.as_str() {
      "BTC" => "bitcoin",
      "ETH" => "ethereum",
      _ => panic!("Unsupported coin!"),
        })
      .as_bytes()
  }
  fn parse_float(data: String) -> Option<f64> {
    data.parse::<f64>().ok()
  }
  #[no_mangle]
  pub fn prepare() {
      let calldata = oei::get_calldata();
      // Coingecko data source
      oei::request_external_data(1, 1, parse_coingecko_symbol(&calldata));
      // Crypto compare source
      oei::request_external_data(2, 2, &calldata);
      // Binance source
      oei::request_external_data(3, 3, &calldata);
  } #[no_mangle]
  pub fn execute() {
    let validator_count = oei::get_requested_validator_count();
    let mut sum: f64 = 0.0;
    let mut count: u64 = 0;
    for validator_index in 0..validator_count {
      let mut val = 0.0;
      let mut fail = false;
      for external_id in 1..4 {
        let data = oei::get_external_data(external_id, validator_index);
        if data.is_none() {
          fail = true;
          break;
        }
      let num = parse_float(data.unwrap());
      if num.is_none() {
        fail = true; break;
      }
      val += num.unwrap();
    }
    if !fail { sum += val / 3.0; count += 1; }
  }
  let result = (sum / (count as f64) * 100.0) as u64;
  oei::save_return_data(&result.to_be_bytes())
} |f};
  let codes = [(code1, "Cargo.toml"), (code2, "src/logic.rs")];

  <div className=Styles.tableWrapper>
    <>
      <VSpacing size={`px(10)} />
      <Row>
        <HSpacing size={`px(15)} />
        <Col>
          <div> <Text value="Platform" /> </div>
          <VSpacing size={`px(5)} />
          <div> <Text value="OWASM v0.1" code=true weight=Text.Semibold /> </div>
        </Col>
        <HSpacing size={`px(370)} />
        <Col>
          <div> <Text value="Language" /> </div>
          <VSpacing size={`px(5)} />
          <div> <Text value="Rust 1.40.0" code=true weight=Text.Semibold /> </div>
        </Col>
      </Row>
      <VSpacing size={`px(35)} />
      {codes
       ->Belt_List.map(((co, name)) => {
           React.useMemo1(
             () =>
               <div className=Styles.tableLowerContainer>
                 <div className=Styles.vFlex>
                   <img src=Images.code className=Styles.codeImage />
                   <Text value=name size=Text.Lg color=Colors.gray7 />
                 </div>
                 <VSpacing size=Spacing.lg />
                 {co |> renderCode}
               </div>,
             [||],
           )
         })
       ->Array.of_list
       ->React.array}
    </>
  </div>;
};
