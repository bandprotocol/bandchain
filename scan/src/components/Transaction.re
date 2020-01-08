type txType =
  | DataRequest(string)
  | NewScript(string);

let txColor =
  fun
  | DataRequest(_) => Colors.purple
  | NewScript(_) => Colors.pink;

module Styles = {
  open Css;

  let container =
    style([
      padding(Spacing.md),
      paddingTop(`px(9)),
      paddingBottom(`px(15)),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
      marginBottom(`px(1)),
      cursor(`pointer),
      hover([backgroundColor(Colors.purpleLighter)]),
    ]);

  let header =
    style([
      paddingTop(`px(0)),
      paddingBottom(`px(0)),
      cursor(`default),
      hover([backgroundColor(white)]),
    ]);

  let txTypeContainer = txType =>
    style([
      marginLeft(`px(-2)),
      display(`inlineFlex),
      justifyContent(`center),
      alignItems(`center),
      borderRadius(`px(15)),
      padding2(~v=Spacing.xs, ~h=Spacing.sm),
      color(txType->txColor),
      backgroundColor(
        switch (txType) {
        | DataRequest(_) => Colors.purpleLight
        | NewScript(_) => Colors.pinkLight
        },
      ),
    ]);

  let txIcon =
    style([
      width(`px(28)),
      height(`px(28)),
      marginTop(`px(5)),
      marginLeft(Spacing.md),
      marginRight(Spacing.xl),
    ]);

  let hashCol = style([width(`px(250))]);
  let feeCol = style([width(`px(80))]);
};

let txLabel =
  fun
  | DataRequest(_) => "DATA REQUEST"
  | NewScript(_) => "NEW SCRIPT";

let txSource =
  fun
  | DataRequest(source) => source
  | NewScript(source) => source;

let txIcon =
  fun
  | DataRequest(_) => Images.dataRequest
  | NewScript(_) => Images.newScript;

let renderDataType = txType =>
  <>
    <div className={Styles.txTypeContainer(txType)}>
      <Text value={txType->txLabel} size=Text.Xs block=true />
    </div>
    <VSpacing size=Spacing.xs />
    <Text value={txType->txSource} size=Text.Lg weight=Text.Semibold block=true />
  </>;

let renderHeader = () =>
  <div className={Css.merge([Styles.container, Styles.header])}>
    <Row>
      <Col> <div className=Styles.txIcon /> </Col>
      <Col>
        <div className=Styles.hashCol>
          <Text block=true value="TX HASH" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </div>
      </Col>
      <Col size=2.>
        <Text
          block=true
          value="SOURCE & TYPE"
          size=Text.Sm
          weight=Text.Bold
          color=Colors.grayText
        />
      </Col>
      <Col>
        <div className=Styles.feeCol>
          <Text block=true value="FEE" size=Text.Sm weight=Text.Bold color=Colors.grayText />
        </div>
      </Col>
    </Row>
  </div>;

[@react.component]
let make = (~type_, ~hash, ~timestamp, ~fee) => {
  <div className=Styles.container>
    <Row>
      <Col> <img src={type_->txIcon} className=Styles.txIcon /> </Col>
      <Col>
        <div className=Styles.hashCol>
          <VSpacing size={`px(9)} />
          timestamp
          <VSpacing size={`px(6)} />
          <Text block=true code=true value=hash size=Text.Lg weight=Text.Bold />
        </div>
      </Col>
      <Col size=2.> {renderDataType(type_)} </Col>
      <Col>
        <div className=Styles.feeCol>
          <VSpacing size={`px(4)} />
          <Text size=Text.Sm block=true value="$0.002" color=Colors.grayText />
          <VSpacing size={`px(4)} />
          <Text value=fee color={Css.hex("555")} weight=Text.Semibold />
        </div>
      </Col>
    </Row>
  </div>;
};
