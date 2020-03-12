module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`row), alignItems(`center)]);

  let pageContainer = style([paddingTop(`px(50)), minHeight(`px(500))]);

  let seperatedLine =
    style([
      width(`px(13)),
      height(`px(1)),
      marginLeft(`px(10)),
      marginRight(`px(10)),
      backgroundColor(Colors.grayHeader),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);

  let fullWidth = style([width(`percent(100.0)), display(`flex)]);

  let feeContainer = style([display(`flex), justifyContent(`flexEnd), maxWidth(`px(150))]);
};

[@react.component]
let make = () => {
  <div className=Styles.pageContainer>
    <Row>
      <Col>
        <div className=Styles.vFlex>
          <img src=Images.dataSourceLogo className=Styles.logo />
          <Text
            value="ALL SOURCES"
            weight=Text.Medium
            size=Text.Md
            spacing={Text.Em(0.06)}
            height={Text.Px(15)}
            nowrap=true
            color=Colors.grayHeader
            block=true
          />
          <div className=Styles.seperatedLine />
          <Text
            value="20 In total"
            size=Text.Md
            weight=Text.Thin
            spacing={Text.Em(0.06)}
            color=Colors.grayHeader
            nowrap=true
          />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.xl />
    <>
      <THead>
        <Row>
          <Col> <HSpacing size=Spacing.xl /> </Col>
          <Col size=0.5>
            <div className=TElement.Styles.hashContainer>
              <Text
                block=true
                value="NAME"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.grayText
                spacing={Text.Em(0.1)}
              />
            </div>
          </Col>
          <Col size=0.5>
            <Text
              block=true
              value="AGE"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.grayText
              spacing={Text.Em(0.1)}
            />
          </Col>
          <Col size=1.>
            <Text
              block=true
              value="OWNER"
              size=Text.Sm
              weight=Text.Semibold
              color=Colors.grayText
              spacing={Text.Em(0.1)}
            />
          </Col>
          <Col size=0.4>
            <div className=Styles.feeContainer>
              <Text
                block=true
                value="REQUEST FEE (BAND)"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.grayText
                spacing={Text.Em(0.1)}
              />
            </div>
          </Col>
          <Col> <HSpacing size=Spacing.xl /> </Col>
        </Row>
      </THead>
      {[
         (
           1,
           "CoinGecko V.2",
           MomentRe.momentNow(),
           "band17rprjgtj0krfw3wyl9creueej6ca9dc4dgxv6e" |> Address.fromBech32,
           123.3,
         ),
         (
           2,
           "Koo",
           MomentRe.momentNow()
           |> MomentRe.Moment.subtract(~duration=MomentRe.duration(2., `hours)),
           "band17rprjgtj0krfw3wyl9creueej6ca9dc4dgxv6e" |> Address.fromBech32,
           123.3,
         ),
         (
           3,
           "Binance",
           MomentRe.momentNow()
           |> MomentRe.Moment.subtract(~duration=MomentRe.duration(10., `hours)),
           "band17rprjgtj0krfw3wyl9creueej6ca9dc4dgxv6e" |> Address.fromBech32,
           123.3,
         ),
         (
           4,
           "CMC",
           MomentRe.momentNow()
           |> MomentRe.Moment.subtract(~duration=MomentRe.duration(22., `hours)),
           "band17rprjgtj0krfw3wyl9creueej6ca9dc4dgxv6e" |> Address.fromBech32,
           123.3,
         ),
       ]
       ->Belt.List.map(((id, name, timestamp, owner, fee)) => {
           <TBody key=name>
             <div className=Styles.fullWidth>
               <Row>
                 <Col> <HSpacing size=Spacing.xl /> </Col>
                 <Col size=0.5>
                   <TElement elementType={TElement.DataSource(ID.DataSource.ID(id), name)} />
                 </Col>
                 <Col size=0.5> <TElement elementType={timestamp->TElement.Timestamp} /> </Col>
                 <Col size=1.> <TElement elementType={owner->TElement.Address} /> </Col>
                 <Col size=0.4> <TElement elementType={fee->TElement.Fee} /> </Col>
                 <Col> <HSpacing size=Spacing.xl /> </Col>
               </Row>
             </div>
           </TBody>
         })
       ->Array.of_list
       ->React.array}
    </>
    <VSpacing size=Spacing.xl />
  </div>;
};
