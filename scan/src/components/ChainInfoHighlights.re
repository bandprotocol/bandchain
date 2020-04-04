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
      position(`relative),
    ]);

  let innerCard =
    style([
      display(`flex),
      width(`percent(100.)),
      flexDirection(`column),
      justifyContent(`spaceBetween),
      alignItems(`flexStart),
      position(`relative),
      zIndex(2),
      padding4(~top=`px(13), ~bottom=`px(13), ~left=`px(14), ~right=`px(10)),
    ]);

  let labelContainer =
    style([display(`flex), justifyContent(`spaceBetween), width(`percent(100.))]);

  let bandPriceExtra =
    style([display(`flex), width(`percent(100.)), justifyContent(`spaceBetween)]);

  let vFlex = style([display(`flex), flexDirection(`row)]);

  let withWidth = (w: int) => style([width(`px(w))]);

  let bgCard = (url: string) =>
    style([
      backgroundImage(`url(url)),
      top(`px(12)),
      backgroundSize(`contain),
      backgroundRepeat(`noRepeat),
      width(`percent(100.)),
      height(`percent(100.)),
      position(`absolute),
      zIndex(1),
      opacity(0.4),
    ]);
};

module HighlightCard = {
  [@react.component]
  let make = (~loading=true, ~label, ~valueComponent, ~extraComponent, ~extraTopRight=?, ~bgUrl=?) => {
    <div className=Styles.card>
      {switch (bgUrl, loading) {
       | (Some(url), false) => <div className={Styles.bgCard(url)} />
       | _ => React.null
       }}
      <div className=Styles.innerCard>
        <div className=Styles.labelContainer>
          {loading
             ? <LoadingCensorBar width=90 height=18 />
             : <Text value=label color=Colors.bandBlue spacing={Text.Em(0.05)} />}
          {extraTopRight->Belt.Option.getWithDefault(React.null)}
        </div>
        {loading ? <LoadingCensorBar width=120 height=20 /> : valueComponent}
        <div className={Styles.withWidth(170)}>
          {loading ? <LoadingCensorBar width=75 height=15 /> : extraComponent}
        </div>
      </div>
    </div>;
  };
};

[@react.component]
let make = () => {
  let latestBlockSub = BlockSub.getLatest();
  let infoOpt = React.useContext(GlobalContext.context);
  let validatorsSub = ValidatorSub.getList();

  switch (latestBlockSub, validatorsSub) {
  | (Error(_) | NoData, Error(_) | NoData) => React.null
  | (latestBlock_, validators_) =>
    let loading =
      switch (latestBlock_, validators_) {
      | (ApolloHooks.Subscription.Data(_), ApolloHooks.Subscription.Data(_)) => false
      | _ => true
      };

    let (lastProcessedHeight, moniker) =
      switch (latestBlock_) {
      | ApolloHooks.Subscription.Data({height, validator: {moniker}}) => (height, moniker)
      | _ => (ID.Block.ID(0), "")
      };
    let validators =
      switch (validators_) {
      | ApolloHooks.Subscription.Data(vs) => vs
      | _ => [||]
      };

    // TODO replace this Mock finance.
    let mockFinance: PriceHook.Price.t = {
      usdPrice: 0.,
      usdMarketCap: 0.,
      usd24HrChange: 0.,
      btcPrice: 0.,
      btcMarketCap: 0.,
      btc24HrChange: 0.,
      circulatingSupply: 0.,
    };
    let financial = infoOpt->Belt_Option.mapWithDefault(mockFinance, info => info.financial);
    let bandBonded = validators->Belt_Array.map(x => x.tokens)->Belt_Array.reduce(0.0, (+.));

    <Row justify=Row.Between>
      <HighlightCard
        loading
        label="BAND PRICE"
        bgUrl=Images.graphBG
        valueComponent={
                         let bandPriceInUSD = "$" ++ financial.usdPrice->Format.fPretty;
                         <Text
                           value=bandPriceInUSD
                           size=Text.Xxxl
                           weight=Text.Semibold
                           color=Colors.gray8
                           code=true
                         />;
                       }
        extraComponent={
                         let bandPriceInBTC = financial.btcPrice;
                         let usd24HrChange = financial.usd24HrChange;
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
        loading
        label="MARKET CAP"
        valueComponent={
                         let marketcap = "$" ++ financial.usdMarketCap->Format.fPretty;
                         <Text
                           value=marketcap
                           size=Text.Xxxl
                           weight=Text.Semibold
                           color=Colors.gray8
                           code=true
                         />;
                       }
        extraComponent={
                         let marketcap = financial.circulatingSupply;
                         <div className=Styles.vFlex>
                           <Text value={marketcap->Format.fPretty} code=true weight=Text.Thin />
                           <HSpacing size=Spacing.xs />
                           <Text
                             value="BAND"
                             color=Colors.gray7
                             weight=Text.Thin
                             spacing={Text.Em(0.01)}
                           />
                         </div>;
                       }
      />
      <HighlightCard
        loading
        label="LATEST BLOCK"
        valueComponent={<TypeID.Block id=lastProcessedHeight position=TypeID.Landing />}
        extraComponent={<Text value=moniker nowrap=true ellipsis=true block=true />}
      />
      <HighlightCard
        loading
        label="ACTIVE VALIDATORS"
        valueComponent={
                         let activeValidators =
                           validators->Belt_Array.size->Format.iPretty ++ " Nodes";
                         <Text
                           value=activeValidators
                           size=Text.Xxxl
                           weight=Text.Semibold
                           color=Colors.gray8
                         />;
                       }
        extraComponent={
          <div className=Styles.vFlex>
            <Text value={bandBonded->Format.fPretty} code=true />
            <HSpacing size=Spacing.sm />
            <Text value=" BAND Bonded" />
          </div>
        }
      />
    </Row>;
  };
};
