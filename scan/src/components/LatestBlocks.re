type block = {
  id: int,
  proposer: string,
};

module Styles = {
  open Css;

  let block =
    style([
      backgroundColor(white),
      padding(Spacing.lg),
      marginBottom(Spacing.md),
      boxShadow(Shadow.box(~x=`px(0), ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      width(`px(120)),
      cursor(`pointer),
      transition(~duration=100, "transform"),
      hover([transform(translateY(`px(-3)))]),
    ]);
};

let renderBlock = ((b, moniker): (BlockHook.Block.t, string)) =>
  <div key={b.height |> string_of_int} className=Styles.block>
    <Text value="# " color=Colors.pink weight=Text.Semibold size=Text.Lg />
    <Text value={b.height->Format.iPretty} weight=Text.Semibold size=Text.Lg />
    <VSpacing size=Spacing.md />
    <Text value="PROPOSED BY" block=true size=Text.Xs color=Colors.grayText />
    <VSpacing size=Spacing.xs />
    <Text block=true value=moniker weight=Text.Semibold ellipsis=true />
  </div>;

[@react.component]
let make = () =>
  {
    let%Opt info = React.useContext(GlobalContext.context);
    let blocks = info.latestBlocks;
    let validators = info.validators;
    let blocksWithMonikers =
      blocks
      ->Belt_List.map(block => {
          Js.Console.log2(block, validators |> Belt_List.toArray);
          BlockHook.Block.getProposerMoniker(block, validators);
        })
      ->Belt_List.zip(blocks, _);

    Some(
      <Row alignItems=`initial>
        <Col>
          {blocksWithMonikers
           ->Belt.List.keepWithIndex((_b, i) => i mod 2 == 0)
           ->Belt.List.map(renderBlock)
           ->Array.of_list
           ->React.array}
        </Col>
        <HSpacing size=Spacing.sm />
        <Col>
          <VSpacing size=Spacing.xl />
          {blocksWithMonikers
           ->Belt.List.keepWithIndex((_b, i) => i mod 2 == 1)
           ->Belt.List.map(renderBlock)
           ->Array.of_list
           ->React.array}
        </Col>
      </Row>,
    );
  }
  ->Belt.Option.getWithDefault(React.null);
