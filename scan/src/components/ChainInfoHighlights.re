module Styles = {
  open Css;

  let card =
    style([
      position(`relative),
      backgroundColor(Colors.white),
      borderRadius(`px(4)),
      boxShadow(Shadow.box(~x=`zero, ~y=`zero, ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08)))),
      Media.smallMobile([margin2(~v=`zero, ~h=`px(-5))]),
    ]);

  let innerCard =
    style([
      position(`relative),
      zIndex(2),
      height(`px(130)),
      padding2(~v=`px(16), ~h=`px(24)),
      Media.mobile([padding2(~v=`px(10), ~h=`px(12)), height(`px(120))]),
    ]);

  let bgCard = (url: string) =>
    style([
      backgroundImage(`url(url)),
      backgroundPosition(`center),
      top(`px(12)),
      backgroundSize(`contain),
      backgroundRepeat(`noRepeat),
      width(`percent(100.)),
      height(`percent(100.)),
      position(`absolute),
      zIndex(1),
      opacity(0.4),
    ]);

  let fullWidth = style([width(`percent(100.))]);

  let creditContaier =
    style([paddingTop(`px(8)), Media.mobile([padding2(~v=`px(8), ~h=`zero)])]);
};

module HighlightCard = {
  [@react.component]
  let make = (~label, ~valueAndExtraComponentSub, ~bgUrl=?) => {
    <div className=Styles.card>
      {switch (bgUrl, valueAndExtraComponentSub) {
       | (Some(url), ApolloHooks.Subscription.Data(_)) => <div className={Styles.bgCard(url)} />
       | _ => React.null
       }}
      <div
        id={"highlight-" ++ label}
        className={Css.merge([
          Styles.innerCard,
          CssHelper.flexBox(~direction=`column, ~justify=`spaceBetween, ~align=`flexStart, ()),
        ])}>
        {switch (valueAndExtraComponentSub) {
         | Data((valueComponent, extraComponent)) =>
           <>
             <Text value=label size=Text.Lg color=Colors.gray7 weight=Text.Semibold block=true />
             valueComponent
             extraComponent
           </>
         | _ =>
           <>
             <LoadingCensorBar width=90 height=18 />
             <LoadingCensorBar width=120 height=20 />
             <LoadingCensorBar width=75 height=15 />
           </>
         }}
      </div>
    </div>;
  };
};

[@react.component]
let make = (~latestBlockSub: Sub.t(BlockSub.t)) => {
  let infoSub = React.useContext(GlobalContext.context);
  let activeValidatorCountSub = ValidatorSub.countByActive(true);
  let bondedTokenCountSub = ValidatorSub.getTotalBondedAmount();

  let validatorInfoSub = Sub.all2(activeValidatorCountSub, bondedTokenCountSub);
  let allSub = Sub.all3(latestBlockSub, infoSub, validatorInfoSub);

  <Row justify=Row.Between>
    <Col.Grid col=Col.Three colSm=Col.Six>
      <HighlightCard
        label="Band Price"
        bgUrl=Images.graphBG
        valueAndExtraComponentSub={
          let%Sub (_, {financial}, _) = allSub;
          (
            {
              let bandPriceInUSD = "$" ++ (financial.usdPrice |> Format.fPretty(~digits=2));
              <Text value=bandPriceInUSD size=Text.Xxxl color=Colors.gray7 code=true />;
            },
            {
              let bandPriceInBTC = financial.btcPrice;
              let usd24HrChange = financial.usd24HrChange;

              <div
                className={Css.merge([
                  CssHelper.flexBox(~justify=`spaceBetween, ()),
                  Styles.fullWidth,
                ])}>
                <Text
                  value={bandPriceInBTC->Format.fPretty ++ " BTC"}
                  color=Colors.gray6
                  code=true
                  spacing={Text.Em(0.01)}
                />
                <Text
                  value={usd24HrChange->Format.fPercentChange}
                  color={usd24HrChange >= 0. ? Colors.green4 : Colors.red5}
                  code=true
                />
              </div>;
            },
          )
          |> Sub.resolve;
        }
      />
      <div className={Css.merge([CssHelper.flexBox(), Styles.creditContaier])}>
        <Text value="Powered by" size=Text.Sm color=Colors.gray7 />
        <HSpacing size=Spacing.xs />
        // TODO: make it to link later
        <Text value="Band Oracle" size=Text.Sm color=Colors.bandBlue weight=Text.Medium />
      </div>
    </Col.Grid>
    <Col.Grid col=Col.Three colSm=Col.Six>
      <HighlightCard
        label="Market Cap"
        valueAndExtraComponentSub={
          let%Sub (_, {financial}, _) = allSub;
          (
            {
              <Text
                value={"$" ++ (financial.usdMarketCap |> Format.fCurrency)}
                size=Text.Xxxl
                color=Colors.gray7
                code=true
              />;
            },
            {
              let marketcap = financial.btcMarketCap;

              <Text
                value={(marketcap |> Format.fPretty) ++ " BTC"}
                code=true
                color=Colors.gray6
              />;
            },
          )
          |> Sub.resolve;
        }
      />
    </Col.Grid>
    <Col.Grid col=Col.Three colSm=Col.Six>
      <HighlightCard
        label="Latest Block"
        valueAndExtraComponentSub={
          let%Sub ({height, validator: {moniker, identity, operatorAddress}}, _, _) = allSub;
          (
            <TypeID.Block id=height position=TypeID.Landing />,
            <ValidatorMonikerLink
              validatorAddress=operatorAddress
              moniker
              identity
              width={`percent(100.)}
              avatarWidth=20
            />,
          )
          |> Sub.resolve;
        }
      />
    </Col.Grid>
    <Col.Grid col=Col.Three colSm=Col.Six>
      <HighlightCard
        label="Active Validators"
        valueAndExtraComponentSub={
          let%Sub (_, _, (activeValidatorCount, bondedTokenCount)) = allSub;
          (
            {
              let activeValidators = activeValidatorCount->Format.iPretty ++ " Nodes";
              <Text value=activeValidators size=Text.Xxxl color=Colors.gray7 />;
            },
            <Text
              value={
                (bondedTokenCount |> Coin.getBandAmountFromCoin |> Format.fPretty)
                ++ " BAND Bonded"
              }
              code=true
              color=Colors.gray6
            />,
          )
          |> Sub.resolve;
        }
      />
    </Col.Grid>
  </Row>;
};
