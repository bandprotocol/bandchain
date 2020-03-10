type pos_t =
  | Title
  | Subtitle
  | Text;

let fontSize =
  fun
  | Title => Text.Xxl
  | Subtitle => Text.Lg
  | Text => Text.Md;

let lineHeight =
  fun
  | Title => Text.Px(23)
  | Subtitle => Text.Px(18)
  | Text => Text.Px(16);

module Styles = {
  open Css;

  let link = style([cursor(`pointer)]);

  let pointerEvents =
    fun
    | Title => style([pointerEvents(`none)])
    | Subtitle
    | Text => style([pointerEvents(`auto)]);
};

module ComponentCreator = (RawID: ID.IDSig) => {
  [@react.component]
  let make = (~id, ~position=Text) =>
    <div className={Css.merge([Styles.link, Styles.pointerEvents(position)])}>
      <Text
        value={id |> RawID.toString}
        size={position |> fontSize}
        weight=Text.Semibold
        height={position |> lineHeight}
        color=RawID.color
        spacing={Text.Em(0.02)}
        nowrap=true
        code=true
        block=true
      />
    </div>;
};

module DataSource = ComponentCreator(ID.DataSource);
module OracleScript = ComponentCreator(ID.OracleScript);
module Request = ComponentCreator(ID.Request);
module Block = ComponentCreator(ID.Block);
