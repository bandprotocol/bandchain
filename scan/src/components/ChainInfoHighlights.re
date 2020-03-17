module Styles = {
  open Css;

  let card =
    style([
      display(`flex),
      width(`px(210)),
      height(`px(115)),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(0, 0, 0, 0.1))),
    ]);

  let innerCard =
    style([
      display(`flex),
      width(`percent(100.)),
      flexDirection(`column),
      justifyContent(`spaceBetween),
      alignItems(`flexStart),
      position(`relative),
      zIndex(1),
      padding4(~top=`px(13), ~bottom=`px(13), ~left=`px(14), ~right=`px(10)),
    ]);

  let labelContainer =
    style([display(`flex), justifyContent(`spaceBetween), width(`percent(100.))]);

  let bandPriceExtra =
    style([display(`flex), width(`percent(100.)), justifyContent(`spaceBetween)]);

  let vFlex = style([display(`flex), flexDirection(`row)]);

  let bgCard = (url: string) =>
    style([
      backgroundImage(`url(url)),
      backgroundPosition(`center),
      backgroundSize(`contain),
      width(`percent(100.)),
    ]);
};

module HighlightCard = {
  [@react.component]
  let make = (~label, ~valueComponent, ~extraComponent, ~extraTopRight=?, ~bg=?) => {
    <div className=Styles.card>
      // <img src=Images.graphBG className=Styles.bgCard />

        <div className=Styles.innerCard>
          <div className=Styles.labelContainer>
            <Text value=label color=Colors.bandBlue spacing={Text.Em(0.05)} />
            {extraTopRight->Belt.Option.getWithDefault(React.null)}
          </div>
          valueComponent
          extraComponent
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
      <Row justify=Row.Between>
        <HighlightCard
          label="BAND PRICE"
          valueComponent={
                           let bandPriceInUSD = "$" ++ info.financial.usdPrice->Format.fPretty;
                           <Text
                             value=bandPriceInUSD
                             size=Text.Xxxl
                             weight=Text.Semibold
                             color=Colors.gray8
                             code=true
                           />;
                         }
          extraComponent={
                           let bandPriceInBTC = info.financial.btcPrice;
                           let usd24HrChange = info.financial.usd24HrChange;
                           <div className=Styles.bandPriceExtra>
                             <div className=Styles.vFlex>
                               <Text
                                 value={bandPriceInBTC->Format.fPretty}
                                 color=Colors.gray7
                                 weight=Text.Thin
                                 code=true
                                 spacing={Text.Em(0.01)}
                               />
                               <HSpacing size=Spacing.xs />
                               <Text
                                 value="BTC"
                                 color=Colors.gray7
                                 weight=Text.Thin
                                 spacing={Text.Em(0.01)}
                               />
                             </div>
                             <Text
                               value={usd24HrChange->Format.fPercent}
                               color={usd24HrChange >= 0. ? Colors.green4 : Colors.red5}
                               weight=Text.Semibold
                               code=true
                             />
                           </div>;
                         }
        />
        <HighlightCard
          label="MARKET CAP"
          valueComponent={
                           let marketcap = "$" ++ info.financial.usdMarketCap->Format.fPretty;
                           <Text
                             value=marketcap
                             size=Text.Xxxl
                             weight=Text.Semibold
                             color=Colors.gray8
                             code=true
                           />;
                         }
          extraComponent={
                           let marketcap = info.financial.btcMarketCap;
                           <div className=Styles.vFlex>
                             <Text value={marketcap->Format.fPretty} code=true />
                             <HSpacing size=Spacing.xs />
                             <Text
                               value="BTC"
                               color=Colors.gray7
                               weight=Text.Thin
                               spacing={Text.Em(0.01)}
                             />
                           </div>;
                         }
        />
        <HighlightCard
          label="LATEST BLOCK"
          extraTopRight={<TimeAgos time={info.latestBlock.timestamp} />}
          valueComponent={
                           let latestBlock = info.latestBlock.height->Format.iPretty;
                           <Text
                             value=latestBlock
                             size=Text.Xxxl
                             weight=Text.Semibold
                             color=Colors.gray8
                             code=true
                           />;
                         }
          extraComponent={<Text value="mock" code=true />}
        />
        <HighlightCard
          label="ACTIVE VALIDATORS"
          valueComponent={
                           let activeValidators =
                             validators->Belt_List.size->Format.iPretty ++ " Nodes";
                           <Text
                             value=activeValidators
                             size=Text.Xxxl
                             weight=Text.Semibold
                             color=Colors.gray8
                             code=true
                           />;
                         }
          extraComponent={
                           let bondedAmount = bandBonded->Format.fPretty ++ " BAND Bonded";
                           <Text value=bondedAmount code=true />;
                         }
        />
      </Row>,
    );
  }
  ->Belt.Option.getWithDefault(React.null);
