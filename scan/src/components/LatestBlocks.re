type block = {
  id: int,
  proposer: string,
};

module Styles = {
  open Css;

  let block = (i, ID.Block.ID(bh)) =>
    style([
      position(`absolute),
      left(`zero),
      backgroundColor(white),
      padding(`px(12)),
      width(`calc((`sub, `percent(50.), `px(4)))),
      height(`px(75)),
      borderRadius(`px(4)),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.08))),
      ),
      opacity(i == 0 || i == 11 ? 0. : 1.),
      pointerEvents(i == 0 || i == 11 ? `none : `auto),
      transition(~duration=800, "all"),
      transform(
        `translate((
          bh mod 2 == 1 ? `calc((`add, `percent(100.), `px(8))) : `zero,
          `px(i mod 2 == 1 ? i / 2 * 85 : (i + 1) / 2 * 85 - 45),
        )),
      ),
      Media.mobile([
        opacity(i == 0 || i == 7 ? 0. : 1.),
        height(`px(75)),
        width(`calc((`sub, `percent(50.), `px(5)))),
        transform(
          `translate((
            bh mod 2 == 1 ? `calc((`add, `percent(100.), `px(10))) : `zero,
            `px(i mod 2 == 1 ? i / 2 * 80 : (i + 1) / 2 * 80 - 70),
          )),
        ),
      ]),
    ]);

  let blocksWrapper =
    style([
      overflow(`hidden),
      minWidth(`px(245)),
      minHeight(`px(500)),
      position(`relative),
      Media.mobile([minHeight(`px(300)), overflow(`hidden), width(`percent(100.))]),
    ]);
};

let renderMiniBlock =
    (index: int, blockHeight: ID.Block.t, validatorOpt: option(ValidatorSub.Mini.t)) => {
  <div
    key={blockHeight |> ID.Block.toString}
    className={Css.merge([
      CssHelper.flexBox(~direction=`column, ~justify=`spaceBetween, ~align=`flexStart, ()),
      Styles.block(index, blockHeight),
    ])}>
    <TypeID.Block id=blockHeight />
    {switch (validatorOpt) {
     | Some(validator) =>
       <ValidatorMonikerLink
         size=Text.Md
         validatorAddress={validator.operatorAddress}
         width={`percent(100.)}
         moniker={validator.moniker}
         identity={validator.identity}
         avatarWidth=20
       />
     | None => React.null
     }}
  </div>;
};

let renderMiniBlockLoading = (imaginaryIndex: int, imaginaryBlockHeight: ID.Block.t) =>
  <div
    key={imaginaryBlockHeight |> ID.Block.toString}
    className={Styles.block(imaginaryIndex, imaginaryBlockHeight)}>
    <LoadingCensorBar width=60 height=15 />
    <VSpacing size=Spacing.lg />
    <VSpacing size={`px(1)} />
    <LoadingCensorBar width=100 height=15 />
  </div>;

let renderMobile = (~blocksSub) => {
  <div className={Css.merge([CssHelper.flexBox(~direction=`column, ()), Styles.blocksWrapper])}>
    {switch (blocksSub) {
     | ApolloHooks.Subscription.Data(blocks) =>
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
  </div>;
};

let renderDesktop = (~blocksSub) => {
  <div className={Css.merge([CssHelper.flexBox(~direction=`column, ()), Styles.blocksWrapper])}>
    {switch (blocksSub) {
     | ApolloHooks.Subscription.Data(blocks) =>
       let {BlockSub.height: ID.Block.ID(blocksCount)} = blocks->Belt_Array.getExn(0);
       [|renderMiniBlock(0, ID.Block.ID(blocksCount + 1), None)|]
       ->Belt_Array.concat(
           blocks->Belt_Array.mapWithIndex((i, {height, validator}) => {
             renderMiniBlock(i + 1, height, Some(validator))
           }),
         )
       ->React.array;
     | _ => Belt_Array.makeBy(11, i => renderMiniBlockLoading(i, ID.Block.ID(i)))->React.array
     }}
  </div>;
};

[@react.component]
let make = (~blocksSub) => {
  <>
    <div
      className={CssHelper.flexBox(~justify=`spaceBetween, ~align=`flexEnd, ())}
      id="latestBlocksSectionHeader">
      <div>
        <Text
          value="Latest Blocks"
          size=Text.Lg
          block=true
          color=Colors.gray7
          weight=Text.Medium
        />
        <VSpacing size={`px(4)} />
        {switch (blocksSub) {
         | ApolloHooks.Subscription.Data(blocks) =>
           <Text
             value={
               //  TODO: hack for request count, let's fix it later
               blocks
               ->Belt.Array.get(0)
               ->Belt.Option.mapWithDefault(
                   0,
                   ({BlockSub.height}) => {
                     let ID.Block.ID(height_) = height;
                     height_;
                   },
                 )
               ->Format.iPretty
             }
             size=Text.Lg
             color=Colors.gray7
             weight=Text.Medium
           />
         | _ => <LoadingCensorBar width=90 height=18 />
         }}
      </div>
      <Link className={CssHelper.flexBox(~align=`center, ())} route=Route.BlockHomePage>
        <Text value="All Blocks" color=Colors.bandBlue weight=Text.Medium />
        <HSpacing size=Spacing.md />
        <Icon name="fal fa-angle-right" color=Colors.bandBlue />
      </Link>
    </div>
    <VSpacing size={`px(16)} />
    {Media.isMobile() ? renderMobile(blocksSub) : renderDesktop(blocksSub)}
  </>;
};
