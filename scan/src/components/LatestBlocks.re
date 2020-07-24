type block = {
  id: int,
  proposer: string,
};

module Styles = {
  open Css;

  let desktopContainer = style([marginRight(`px(18))]);

  let block = (i, ID.Block.ID(bh)) =>
    style([
      position(`absolute),
      backgroundColor(white),
      padding4(~top=`px(14), ~left=`px(10), ~right=`px(18), ~bottom=`px(16)),
      marginBottom(`px(3)),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      width(`px(120)),
      height(`px(80)),
      pointerEvents(i == 0 || i == 11 ? `none : `auto),
      cursor(`pointer),
      opacity(i == 0 || i == 11 ? 0. : 1.),
      transform(
        `translate((
          `px(bh mod 2 == 1 ? 125 : 0),
          `px(i mod 2 == 1 ? i / 2 * 85 : (i + 1) / 2 * 85 - 42),
        )),
      ),
      transition(~duration=800, "all"),
      Media.mobile([
        opacity(i == 0 || i == 7 ? 0. : 1.),
        height(`px(88)),
        width(`px(168)),
        transform(
          `translate((
            `px(bh mod 2 == 1 ? 180 : 0),
            `px(i mod 2 == 1 ? i / 2 * 100 : (i + 1) / 2 * 100 - 90),
          )),
        ),
      ]),
    ]);

  let rightCol = style([marginLeft(`px(-3))]);

  let topicBar =
    style([
      width(`percent(100.)),
      display(`flex),
      flexDirection(`row),
      justifyContent(`spaceBetween),
    ]);

  let blocksWrapper =
    style([
      display(`flex),
      flexDirection(`column),
      minWidth(`px(245)),
      minHeight(`px(500)),
      position(`relative),
      Media.mobile([minHeight(`px(370)), overflow(`hidden), width(`percent(100.))]),
    ]);

  let seeAll = style([display(`flex), flexDirection(`row), cursor(`pointer)]);
  let cFlex = style([display(`flex), flexDirection(`column)]);
  let amount =
    style([fontSize(`px(20)), lineHeight(`px(24)), color(Colors.gray8), fontWeight(`bold)]);
  let rightArrow =
    style([
      width(`px(25)),
      marginTop(`px(17)),
      marginLeft(`px(16)),
      Media.mobile([marginTop(`zero)]),
    ]);

  let mobileContainer = style([width(`percent(100.))]);

  let allBlocks = style([display(`flex), flexDirection(`column), alignItems(`flexEnd)]);
};

let renderMiniBlock =
    (index: int, blockHeight: ID.Block.t, validatorOpt: option(ValidatorSub.Mini.t)) => {
  <div key={blockHeight |> ID.Block.toString} className={Styles.block(index, blockHeight)}>
    <TypeID.Block id=blockHeight />
    <VSpacing size=Spacing.sm />
    <Text value="PROPOSED BY" block=true size=Text.Xs color=Colors.gray7 spacing={Text.Em(0.1)} />
    <VSpacing size={`px(5)} />
    {switch (validatorOpt) {
     | Some(validator) =>
       <ValidatorMonikerLink
         size=Text.Sm
         validatorAddress={validator.operatorAddress}
         width={`px(100)}
         moniker={validator.moniker}
         identity={validator.identity}
       />
     | None =>
       <ValidatorMonikerLink
         size=Text.Sm
         validatorAddress={"" |> Address.fromHex}
         width={`px(10)}
         moniker=""
       />
     }}
  </div>;
};

let renderMiniBlockLoading = (imaginaryIndex: int, imaginaryBlockHeight: ID.Block.t) =>
  <div
    key={imaginaryBlockHeight |> ID.Block.toString}
    className={Styles.block(imaginaryIndex, imaginaryBlockHeight)}>
    <LoadingCensorBar width=90 height=16 />
    <VSpacing size=Spacing.lg />
    <VSpacing size={`px(1)} />
    <LoadingCensorBar width=45 height=16 />
  </div>;

let renderMobile = (~blocksSub) => {
  <div className=Styles.mobileContainer>
    <Row>
      <Col size=0.9>
        <div className=Styles.topicBar>
          <Text
            value="Latest Blocks"
            size=Text.Xxl
            weight=Text.Bold
            block=true
            color=Colors.gray8
          />
        </div>
        {switch (blocksSub) {
         | ApolloHooks.Subscription.Data(blocks) =>
           let {BlockSub.height: ID.Block.ID(blocksCount)} = blocks->Belt_Array.getExn(0);
           <Text
             value={blocksCount |> Format.iPretty}
             size=Text.Xxl
             height={Text.Px(24)}
             color=Colors.gray8
             block=true
             weight=Text.Bold
           />;
         | _ => <LoadingCensorBar width=90 height=24 />
         }}
      </Col>
      <Col size=0.3 justifyContent=Col.End>
        <div className=Styles.allBlocks>
          <img src=Images.rightArrow className=Styles.rightArrow />
          <VSpacing size=Spacing.md />
          <Text
            value="ALL BLOCKS"
            size=Text.Sm
            color=Colors.bandBlue
            spacing={Text.Em(0.05)}
            weight=Text.Medium
          />
        </div>
      </Col>
    </Row>
    <VSpacing size=Spacing.md />
    <Row alignItems=`initial>
      <div className=Styles.blocksWrapper>
        {switch (blocksSub) {
         | Data(blocks) =>
           let {BlockSub.height: ID.Block.ID(blocksCount)} = blocks->Belt_Array.getExn(0);
           [|renderMiniBlock(0, ID.Block.ID(blocksCount + 1), None)|]
           ->Belt_Array.concat(
               blocks->Belt_Array.mapWithIndex((i, {height, validator}) =>
                 renderMiniBlock(i + 1, height, Some(validator))
               ),
             )
           ->React.array;
         | _ => Belt_Array.makeBy(7, i => renderMiniBlockLoading(i, ID.Block.ID(i)))->React.array
         }}
      </div>
    </Row>
  </div>;
};

let renderDesktop = (~blocksSub: ApolloHooks.Subscription.variant(_)) => {
  <div className=Styles.desktopContainer>
    <div className=Styles.topicBar>
      <Text value="Latest Blocks" size=Text.Xxl weight=Text.Bold block=true color=Colors.gray8 />
    </div>
    <VSpacing size=Spacing.lg />
    <VSpacing size=Spacing.sm />
    <Link className=Styles.seeAll route=Route.BlockHomePage>
      <div className=Styles.cFlex>
        {switch (blocksSub) {
         | Data(blocks) =>
           let {BlockSub.height: ID.Block.ID(blocksCount)} = blocks->Belt_Array.getExn(0);
           <Text
             value={blocksCount |> Format.iPretty}
             size=Text.Xxl
             height={Text.Px(24)}
             color=Colors.gray8
             block=true
             weight=Text.Bold
           />;
         | _ => <LoadingCensorBar width=90 height=24 />
         }}
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
    </Link>
    <VSpacing size=Spacing.lg />
    <Row alignItems=`initial>
      <div className=Styles.blocksWrapper>
        {switch (blocksSub) {
         | Data(blocks) =>
           let {BlockSub.height: ID.Block.ID(blocksCount)} = blocks->Belt_Array.getExn(0);
           [|renderMiniBlock(0, ID.Block.ID(blocksCount + 1), None)|]
           ->Belt_Array.concat(
               blocks->Belt_Array.mapWithIndex((i, {height, validator}) => {
                 renderMiniBlock(i + 1, height, Some(validator))
               }),
             )
           ->React.array;
         | _ =>
           Belt_Array.makeBy(11, i => renderMiniBlockLoading(i, ID.Block.ID(i)))->React.array
         }}
      </div>
    </Row>
  </div>;
};

[@react.component]
let make = (~blocksSub) => {
  Media.isMobile() ? renderMobile(~blocksSub) : renderDesktop(~blocksSub);
};
