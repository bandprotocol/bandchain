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

let letterSpacing =
  fun
  | Title
  | Subtitle
  | Text => Text.Em(0.02);

module Styles = {
  open Css;

  let link = style([cursor(`pointer)]);

  let pointerEvents =
    fun
    | Title => style([pointerEvents(`none)])
    | Subtitle
    | Text => style([pointerEvents(`auto)]);
};

[@react.component]
let make = (~id, ~position=Text) => {
  <div
    className={Css.merge([Styles.link, Styles.pointerEvents(position)])}
    onClick={_ => Route.redirect(ID.getRoute(id))}>
    <Text
      value={id |> ID.toString}
      size={position |> fontSize}
      weight=Text.Semibold
      height={position |> lineHeight}
      color={id |> ID.getColor}
      spacing={position |> letterSpacing}
      nowrap=true
      code=true
      block=true
    />
  </div>;
};
