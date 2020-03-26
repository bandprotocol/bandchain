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
      backgroundColor(Colors.gray7),
    ]);

  let textContainer = style([paddingLeft(Spacing.lg), display(`flex)]);

  let logo = style([width(`px(50)), marginRight(`px(10))]);

  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);

  let fullWidth = style([width(`percent(100.0)), display(`flex)]);

  let feeContainer = style([display(`flex), justifyContent(`flexEnd), maxWidth(`px(150))]);

  let loadingContainer =
    style([
      display(`flex),
      justifyContent(`center),
      alignItems(`center),
      height(`px(200)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(white),
    ]);
};

[@react.component]
let make = () =>
  {
    let (page, setPage) = React.useState(_ => 1);
    let pageSize = 10;

    let dataSourcesCountSub = DataSourceSub.count();
    let dataSourcesSub = DataSourceSub.getList(~pageSize, ~page, ());

    let%Sub dataSourcesCount = dataSourcesCountSub;
    let%Sub dataSources = dataSourcesSub;

    let pageCount = Page.getPageCount(dataSourcesCount, pageSize);

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
              color=Colors.gray7
              block=true
            />
            <div className=Styles.seperatedLine />
            <Text
              value={dataSourcesCount->string_of_int ++ " In total"}
              size=Text.Md
              weight=Text.Thin
              spacing={Text.Em(0.06)}
              color=Colors.gray7
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
                  color=Colors.gray5
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
                color=Colors.gray5
                spacing={Text.Em(0.1)}
              />
            </Col>
            <Col size=1.>
              <Text
                block=true
                value="OWNER"
                size=Text.Sm
                weight=Text.Semibold
                color=Colors.gray5
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
                  color=Colors.gray5
                  spacing={Text.Em(0.1)}
                />
              </div>
            </Col>
            <Col> <HSpacing size=Spacing.xl /> </Col>
          </Row>
        </THead>
        {dataSources
         ->Belt_Array.map(({id, name, timestamp, owner, fee}) => {
             <TBody key=name>
               <div className=Styles.fullWidth>
                 <Row>
                   <Col> <HSpacing size=Spacing.xl /> </Col>
                   <Col size=0.5> <TElement elementType={TElement.DataSource(id, name)} /> </Col>
                   <Col size=0.5> <TElement elementType={timestamp->TElement.Timestamp} /> </Col>
                   <Col size=1.> <TElement elementType={owner->TElement.Address} /> </Col>
                   <Col size=0.4>
                     <TElement
                       elementType={fee->TxHook.Coin.getBandAmountFromCoins->TElement.Fee}
                     />
                   </Col>
                   <Col> <HSpacing size=Spacing.xl /> </Col>
                 </Row>
               </div>
             </TBody>
           })
         ->React.array}
      </>
      <VSpacing size=Spacing.lg />
      <Pagination currentPage=page pageCount onPageChange={newPage => setPage(_ => newPage)} />
      <VSpacing size=Spacing.lg />
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(_, React.null);
