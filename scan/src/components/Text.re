type size =
  | Xs
  | Sm
  | Md
  | Lg
  | Xl
  | Xxl;

type weight =
  | Regular
  | Semibold
  | Bold;

module Styles = {
  open Css;
  open Belt.Option;

  let fontSize =
    mapWithDefault(
      _,
      style([fontSize(px(12))]),
      fun
      | Xs => style([fontSize(px(8)), letterSpacing(`em(0.07))])
      | Sm => style([fontSize(px(10)), letterSpacing(`em(0.05))])
      | Md => style([fontSize(px(12))])
      | Lg => style([fontSize(px(14)), lineHeight(`px(18))])
      | Xl => style([fontSize(px(16))])
      | Xxl => style([fontSize(px(24))]),
    );

  let fontWeight =
    mapWithDefault(
      _,
      style([]),
      fun
      | Regular => style([fontWeight(`normal)])
      | Semibold => style([fontWeight(`medium)])
      | Bold => style([fontWeight(`bold)]),
    );

  let noWrap = style([whiteSpace(`nowrap)]);
  let block = style([display(`block)]);

  let code =
    style([
      fontFamily(
        "Fira Code, cousine, sfmono-regular,Consolas,Menlo,liberation mono,ubuntu mono,Courier,monospace",
      ),
    ]);
};

[@react.component]
let make =
    (
      ~size=?,
      ~weight=?,
      ~nowrap=?,
      ~color=?,
      ~block=?,
      ~code=?,
      ~value: string,
    ) => {
  <span
    className={Cn.make([
      Styles.fontSize(size),
      Styles.fontWeight(weight),
      Styles.noWrap->Cn.ifTrue(nowrap->Belt.Option.getWithDefault(false)),
      Styles.block->Cn.ifTrue(block->Belt.Option.getWithDefault(false)),
      Styles.code->Cn.ifTrue(code->Belt.Option.getWithDefault(false)),
      color->Cn.mapSome(c => Css.style([Css.color(c)])),
    ])}>
    {React.string(value)}
  </span>;
};