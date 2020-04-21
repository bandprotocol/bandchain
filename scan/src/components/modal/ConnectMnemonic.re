module Styles = {
  open Css;

  let outer = style([marginTop(`px(27))]);

  let inputBar =
    style([
      width(`px(290)),
      height(`px(30)),
      paddingLeft(`px(9)),
      borderRadius(`px(8)),
      boxShadow(
        Shadow.box(~inset=true, ~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(11, 29, 142, 0.1)),
      ),
    ]);

  let mnemonicHelper =
    style([
      width(`px(130)),
      height(`px(16)),
      display(`flex),
      justifyContent(`spaceBetween),
      alignContent(`center),
      color(Css.hex("5269FF")),
    ]);

  let connectBtn =
    style([
      width(`px(100)),
      height(`px(30)),
      display(`flex),
      justifySelf(`right),
      justifyContent(`center),
      alignItems(`center),
      backgroundImage(
        `linearGradient((
          `deg(90.),
          [(`percent(0.), Css.hex("142ABB")), (`percent(100.), Css.hex("5269FF"))],
        )),
      ),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(82, 105, 255, 0.25))),
      borderRadius(`px(4)),
      cursor(`pointer),
    ]);

  let bottom =
    style([
      width(`px(290)),
      display(`flex),
      justifyContent(`spaceBetween),
      alignItems(`center),
    ]);
};

[@react.component]
let make = () => {
  let (mnemonic, setMnemonic) = React.useState(_ => "");
  <div className=Styles.outer>
    <Row> <Text value="Enter Your Mnemonic" size=Text.Md weight=Text.Medium /> </Row>
    <VSpacing size=Spacing.sm />
    <Row>
      <input
        value=mnemonic
        className=Styles.inputBar
        onChange={event => setMnemonic(ReactEvent.Form.target(event)##value)}
      />
    </Row>
    <VSpacing size={`px(35)} />
    <Row>
      <div className=Styles.bottom>
        <Col>
          <div className=Styles.mnemonicHelper>
            <Text value="What is Mnemonic" />
            <img src=Images.linkIcon />
          </div>
        </Col>
        <Col>
          <div className=Styles.connectBtn>
            <Text value="Connect" weight=Text.Bold size=Text.Md color=Colors.white />
          </div>
        </Col>
      </div>
    </Row>
  </div>;
};
