type block = {
  id: int,
  proposer: string,
};

module Styles = {
  open Css;

  let block =
    style([
      backgroundColor(white),
      padding4(~top=`px(14), ~left=`px(10), ~right=`px(18), ~bottom=`px(16)),
      marginBottom(`px(3)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      width(`px(120)),
      cursor(`pointer),
      transition(~duration=100, "transform"),
      hover([transform(translateY(`px(-3)))]),
    ]);

  let rightCol = style([marginLeft(`px(-3))]);

  let topicBar =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`row),
      justifyContent(`spaceBetween),
    ]);

  let seeAll = style([display(`flex), flexDirection(`row), cursor(`pointer)]);
  let cFlex = style([display(`flex), flexDirection(`column)]);
  let amount =
    style([fontSize(`px(20)), lineHeight(`px(24)), color(Colors.gray8), fontWeight(`bold)]);
  let rightArrow = style([width(`px(25)), marginTop(`px(17)), marginLeft(`px(16))]);
};

let renderBlock = (b: BlockSub.t) =>
  <div
    key={b.height |> ID.Block.toString}
    className=Styles.block
    onClick={_ => Route.redirect(b.height |> ID.Block.getRoute)}>
    <TypeID.Block id={b.height} />
    <VSpacing size=Spacing.md />
    <Text value="PROPOSED BY" block=true size=Text.Xs color=Colors.gray7 spacing={Text.Em(0.1)} />
    <VSpacing size={`px(1)} />
    <Text
      block=true
      value={b.validator.moniker}
      weight=Text.Bold
      ellipsis=true
      height={Text.Px(15)}
      spacing={Text.Em(0.02)}
    />
  </div>;

[@react.component]
let make = () =>
  {
    let%Opt info = React.useContext(GlobalContext.context);
    let blocks = info.latestBlocks;

    let blocksCount = BlockSub.count()->Sub.default(0);

    Some(
      <>
        <div className=Styles.topicBar>
          <Text
            value="Latest Blocks"
            size=Text.Xxl
            weight=Text.Bold
            block=true
            color=Colors.gray8
          />
        </div>
        <VSpacing size=Spacing.lg />
        <VSpacing size=Spacing.sm />
        <div className=Styles.seeAll onClick={_ => Route.redirect(Route.BlockHomePage)}>
          <div className=Styles.cFlex>
            <span className=Styles.amount> {blocksCount |> Format.iPretty |> React.string} </span>
            <VSpacing size=Spacing.xs />
            <Text
              value="ALL BLOCKS"
              size=Text.Sm
              color=Colors.bandBlue
              spacing={Text.Em(0.05)}
              weight=Text.Medium
            />
          </div>
          <img src=Images.rightArrow className=Styles.rightArrow />
        </div>
        <VSpacing size=Spacing.lg />
        <Row alignItems=`initial>
          <Col>
            {blocks
             ->Belt.List.keepWithIndex((_b, i) => i mod 2 == 0)
             ->Belt.List.map(renderBlock)
             ->Array.of_list
             ->React.array}
          </Col>
          <Col>
            <div className=Styles.rightCol>
              <VSpacing size=Spacing.xl />
              {blocks
               ->Belt.List.keepWithIndex((_b, i) => i mod 2 == 1)
               ->Belt.List.map(renderBlock)
               ->Array.of_list
               ->React.array}
            </div>
          </Col>
        </Row>
      </>,
    );
  }
  ->Belt.Option.getWithDefault(React.null);
