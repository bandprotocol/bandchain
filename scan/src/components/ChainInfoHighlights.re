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
  let make = (~label, ~valueAndExtraComponentSub: ApolloHooks.Subscription.variant(_), ~bgUrl=?) => {
    <div className=Styles.card>
      {switch (bgUrl, valueAndExtraComponentSub) {
       | (Some(url), Data(_)) => <div className={Styles.bgCard(url)} />
       | _ => React.null
       }}
      <div className=Styles.innerCard>
        <div className=Styles.labelContainer>
          {switch (valueAndExtraComponentSub) {
           | Data(_) => <Text value=label color=Colors.bandBlue spacing={Text.Em(0.05)} />
           | _ => <LoadingCensorBar width=90 height=18 />
           }}
        </div>
        {switch (valueAndExtraComponentSub) {
         | Data((valueComponent, extraComponent)) => <> valueComponent extraComponent </>
         | _ =>
           <> <LoadingCensorBar width=120 height=20 /> <LoadingCensorBar width=75 height=15 /> </>
         }}
      </div>
    </div>;
  };
};

[@react.component]
let make = (~latestBlockSub: Sub.t(BlockSub.t)) => {
  let infoSub = React.useContext(GlobalContext.context);
  // TODO: Only count active validators
  let validatorCountSub = ValidatorSub.count();

  let allSub = Sub.all3(latestBlockSub, infoSub, validatorCountSub);

  <Row justify=Row.Between>
    <HighlightCard
      label="BAND PRICE"
      bgUrl=Images.graphBG
      valueAndExtraComponentSub={
        let%Sub (_, {financial}, _) = allSub;
        (
          {
            let bandPriceInUSD = "$" ++ financial.usdPrice->Format.fPretty;
            <Text
              value=bandPriceInUSD
              size=Text.Xxxl
              weight=Text.Semibold
              color=Colors.gray8
              code=true
            />;
          },
          {
            let bandPriceInBTC = financial.btcPrice;
            let usd24HrChange = financial.usd24HrChange;
            <div className={Styles.withWidth(170)}>
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
                  <Text value="BTC" color=Colors.gray7 weight=Text.Thin spacing={Text.Em(0.01)} />
                </div>
                <Text
                  value={usd24HrChange->Format.fPercent}
                  color={usd24HrChange >= 0. ? Colors.green4 : Colors.red5}
                  weight=Text.Semibold
                  code=true
                />
              </div>
            </div>;
          },
        )
        |> Sub.resolve;
      }
    />
    <HighlightCard
      label="MARKET CAP"
      valueAndExtraComponentSub={
        let%Sub (_, {financial}, _) = allSub;
        (
          {
            let marketcap = "$" ++ financial.usdMarketCap->Format.fPretty;
            <Text
              value=marketcap
              size=Text.Xxxl
              weight=Text.Semibold
              color=Colors.gray8
              code=true
            />;
          },
          {
            let marketcap = financial.circulatingSupply;
            <div className={Styles.withWidth(170)}>
              <div className=Styles.vFlex>
                <Text value={marketcap->Format.fPretty} code=true weight=Text.Thin />
                <HSpacing size=Spacing.xs />
                <Text value="BAND" color=Colors.gray7 weight=Text.Thin spacing={Text.Em(0.01)} />
              </div>
            </div>;
          },
        )
        |> Sub.resolve;
      }
    />
    <HighlightCard
      label="LATEST BLOCK"
      valueAndExtraComponentSub={
        let%Sub ({height, validator: {moniker}}, _, _) = allSub;
        (
          <TypeID.Block id=height position=TypeID.Landing />,
          <div className={Styles.withWidth(170)}>
            <Text value=moniker nowrap=true ellipsis=true block=true />
          </div>,
        )
        |> Sub.resolve;
      }
    />
    <HighlightCard
      label="ACTIVE VALIDATORS"
      valueAndExtraComponentSub={
        let%Sub (_, _, validatorCount) = allSub;
        (
          {
            let activeValidators = validatorCount->Format.iPretty ++ " Nodes";
            <Text value=activeValidators size=Text.Xxxl weight=Text.Semibold color=Colors.gray8 />;
          },
          <div className={Styles.withWidth(170)}>
            <div className=Styles.vFlex>
              // TODO: Replace this mock by the real value

                <Text value="A lot of" />
                <HSpacing size=Spacing.sm />
                <Text value=" BANDs Bonded" />
              </div>
          </div>,
        )
        |> Sub.resolve;
      }
    />
  </Row>;
};
