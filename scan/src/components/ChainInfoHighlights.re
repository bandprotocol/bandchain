module Highlights = {
  open Belt.Option;

  module Styles = {
    open Css;
    let highlights = style([textAlign(`center)]);
  };

  [@react.component]
  let make = (~label, ~value, ~valuePrefix=?, ~extraText=?, ~extraComponent=?, ~extraSuffix=?) => {
    <div className=Styles.highlights>
      <div> <Text value=label size=Text.Sm weight=Text.Bold color=Colors.purple /> </div>
      <div className={Css.style([Css.marginTop(Spacing.sm)])}>
        {valuePrefix->getWithDefault(React.string(""))}
        <Text value size=Text.Xxl weight=Text.Bold />
      </div>
      <div>
        {<Text value={extraText->Belt_Option.getWithDefault("")} size=Text.Sm />}
        {extraComponent->Belt_Option.getWithDefault(React.null)}
        {extraSuffix->getWithDefault(React.string(""))}
      </div>
    </div>;
  };
};

[@react.component]
let make = () =>
  {
    let%Opt info = React.useContext(GlobalContext.context);

    let validators = info.validators;
    let bandBonded = validators->Belt_List.map(x => x.tokens)->Belt_List.reduce(0.0, (+.));

    Some(
      <Row>
        <Col size=1.>
          <Highlights
            label="BAND PRICE"
            value={"$" ++ info.financial.usdPrice->Format.fPretty}
            extraText={"@" ++ info.financial.btcPrice->Format.fPretty ++ " BTC "}
            extraSuffix={
              <Text
                value={"(" ++ info.financial.usd24HrChange->Format.fPercent ++ ")"}
                size=Text.Sm
                color={info.financial.usd24HrChange >= 0. ? Colors.green : Colors.red}
              />
            }
          />
        </Col>
        <Col size=1.>
          <Highlights
            label="MARKET CAP"
            value={"$" ++ info.financial.usdMarketCap->Format.fPretty}
            extraText={info.financial.btcMarketCap->Format.fPretty ++ " BTC "}
          />
        </Col>
        <Col size=1.>
          <Highlights
            label="LATEST BLOCK"
            valuePrefix={<Text value="# " size=Text.Xxl weight=Text.Bold color=Colors.pink />}
            value={info.latestBlock.height->Format.iPretty}
            extraComponent={<TimeAgos time={info.latestBlock.timestamp} />}
          />
        </Col>
        <Col size=1.>
          <Highlights
            label="ACTIVE VALIDATORS"
            value={validators->Belt_List.size->Format.iPretty ++ " Nodes"}
            extraText={bandBonded->Format.fPretty ++ " BAND Bonded"}
          />
        </Col>
      </Row>,
    );
  }
  ->Belt.Option.getWithDefault(React.null);
