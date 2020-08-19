module Styles = {
  open Css;

  let container =
    style([
      display(`flex),
      flexDirection(`row),
      width(`percent(100.)),
      justifyContent(`center),
      minHeight(`px(30)),
      padding2(~v=`px(24), ~h=`zero),
    ]);

  let innerContainer =
    style([
      display(`flex),
      alignItems(`center),
      Media.mobile([
        width(`percent(100.)),
        justifyContent(`spaceBetween),
        padding2(~v=`zero, ~h=`px(5)),
      ]),
    ]);

  let angle = isFlip =>
    style([
      width(`px(12)),
      cursor(`pointer),
      transform(`rotateZ(`deg(isFlip ? 0. : 180.))),
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
let make = (~currentPage, ~pageCount, ~onPageChange: int => unit) =>
  if (pageCount > 1) {
    <div className=Styles.container>
      <div className=Styles.innerContainer>
        <ClickableText isFirst=true onClick={_ => onPageChange(1)} active={currentPage != 1} />
        <HSpacing size=Spacing.lg />
        <ClickableSymbol
          isPrevious=true
          active={currentPage != 1}
          onClick={_ => onPageChange(currentPage < 1 ? 1 : currentPage - 1)}
        />
        <div className=Styles.flex>
          <HSpacing size=Spacing.md />
          <Text value="Page" spacing={Text.Em(0.03)} />
          <HSpacing size=Spacing.sm />
          <Text
            value={currentPage |> Format.iPretty}
            spacing={Text.Em(0.03)}
            weight=Text.Semibold
          />
          <HSpacing size=Spacing.sm />
          <Text value="of" spacing={Text.Em(0.03)} />
          <HSpacing size=Spacing.sm />
          <Text
            value={pageCount |> Format.iPretty}
            spacing={Text.Em(0.03)}
            weight=Text.Semibold
          />
          <HSpacing size=Spacing.md />
        </div>
        <ClickableSymbol
          isPrevious=false
          active={currentPage != pageCount}
          onClick={_ => onPageChange(currentPage > pageCount ? pageCount : currentPage + 1)}
        />
        <HSpacing size=Spacing.lg />
        <ClickableText
          isFirst=false
          active={currentPage != pageCount}
          onClick={_ => onPageChange(pageCount)}
        />
      </div>
    </div>;
  } else {
    React.null;
  };
