module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      flexDirection(`row),
      width(`percent(100.)),
      justifyContent(`center),
      minHeight(`px(30)),
    ]);

  let innerContainer = style([display(`inlineFlex), alignItems(`center)]);

  let angle = isflip =>
    style([
      width(`px(12)),
      cursor(`pointer),
      transform(`rotateZ(`deg(isflip ? 0. : 180.))),
    ]);

  let clickable = active =>
    style([
      cursor(`pointer),
      pointerEvents(active ? `auto : `none),
      opacity(active ? 1. : 0.5),
    ]);

  let flex = style([display(`flex)]);
};

module ClickableText = {
  [@react.component]
  let make = (~isFirst, ~active, ~onClick) => {
    <div className={Css.merge([Styles.flex, Styles.clickable(active)])} onClick>
      <Text value={isFirst ? "First" : "Last"} spacing={Text.Em(0.03)} />
    </div>;
  };
};

module ClickableSymbol = {
  [@react.component]
  let make = (~isPrevious, ~active, ~onClick) => {
    <img
      src=Images.leftAngle
      className={Css.merge([Styles.angle(isPrevious), Styles.clickable(active)])}
      onClick
    />;
  };
};

[@react.component]
let make = (~currentPage, ~pageCount, ~onChangePage: int => unit) => {
  <div className=Styles.container>
    <div className=Styles.innerContainer>
      <ClickableText isFirst=true onClick={_ => onChangePage(1)} active={currentPage != 1} />
      <HSpacing size=Spacing.lg />
      <ClickableSymbol
        isPrevious=true
        active={currentPage != 1}
        onClick={_ => onChangePage(currentPage < 1 ? 1 : currentPage - 1)}
      />
      <HSpacing size=Spacing.md />
      <Text value="Page" spacing={Text.Em(0.03)} />
      <HSpacing size=Spacing.sm />
      <Text value={currentPage |> Format.iPretty} spacing={Text.Em(0.03)} weight=Text.Semibold />
      <HSpacing size=Spacing.sm />
      <Text value="of" spacing={Text.Em(0.03)} />
      <HSpacing size=Spacing.sm />
      <Text value={pageCount |> Format.iPretty} spacing={Text.Em(0.03)} weight=Text.Semibold />
      <HSpacing size=Spacing.md />
      <ClickableSymbol
        isPrevious=false
        active={currentPage != pageCount}
        onClick={_ => onChangePage(currentPage > pageCount ? pageCount : currentPage + 1)}
      />
      <HSpacing size=Spacing.lg />
      <ClickableText
        isFirst=false
        active={currentPage != pageCount}
        onClick={_ => onChangePage(pageCount)}
      />
    </div>
  </div>;
};
