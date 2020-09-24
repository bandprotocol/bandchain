module Styles = {
  open Css;

  let container =
    style([
      maxWidth(`px(600)),
      width(`percent(100.)),
      height(`percent(100.)),
      position(`relative),
      marginTop(Spacing.xs),
      Media.mobile([margin(`zero), display(`flex), maxWidth(`percent(100.))]),
    ]);
  let searchIcon =
    style([
      position(`absolute),
      top(`px(10)),
      left(`px(15)),
      width(`px(15)),
      height(`px(15)),
    ]);
  let search =
    style([
      width(`percent(100.)),
      background(white),
      borderRadius(`px(4)),
      padding4(~left=`px(15), ~right=Spacing.md, ~top=`px(10), ~bottom=`px(10)),
      boxShadows([
        Shadow.box(~x=`zero, ~y=`px(1), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.07))),
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(12), Css.rgba(0, 0, 0, `num(0.02))),
      ]),
      fontSize(`px(12)),
      outline(`px(1), `none, white),
      border(`px(1), `solid, white),
      placeholder([color(Colors.blueGray3)]),
    ]);

  let button =
    style([
      position(`absolute),
      right(`zero),
      width(`px(45)),
      height(`percent(100.)),
      backgroundColor(Colors.blue1),
      borderTopRightRadius(`px(4)),
      borderBottomRightRadius(`px(4)),
      fontSize(`px(14)),
      fontWeight(`medium),
      color(rgba(51, 51, 51, `num(0.54))),
      cursor(`pointer),
      border(`zero, `solid, white),
    ]);
};

module SearchResults = {
  module Styles = {
    open Css;
    let container =
      style([
        position(`absolute),
        left(`zero),
        right(`px(110)),
        top(`percent(90.)),
        backgroundColor(white),
        borderRadius(`px(4)),
        boxShadows([
          Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, `num(0.07))),
          Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(12), Css.rgba(0, 0, 0, `num(0.02))),
        ]),
      ]);

    let result = style([padding(Spacing.sm), paddingLeft(`px(38)), cursor(`pointer)]);

    let lastResult =
      style([
        borderBottomLeftRadius(`px(4)),
        borderBottomRightRadius(`px(4)),
        paddingBottom(Spacing.md),
        paddingTop(`px(9)),
      ]);

    let resultFocused = style([backgroundColor(blueviolet), color(white)]);
  };

  let isValidAddress = a =>
    a->String.sub(0, min(a->String.length, 2)) == "0x" && a->String.length > 2;

  let isValidTx = isValidAddress;

  [@react.component]
  let make = (~searchTerm, ~focusIndex, ~onHover) => {
    let results =
      [|
        searchTerm->isValidAddress
          ? <>
              <VSpacing size={`px(-2)} />
              <Text value="ADDRESS" size=Text.Xs color=Colors.gray5 weight=Text.Semibold />
              <VSpacing size=Spacing.xs />
              <Text value={searchTerm ++ "1f2bce"} weight=Text.Bold size=Text.Lg block=true />
              <VSpacing size=Spacing.sm />
            </>
          : React.null,
        searchTerm->isValidTx
          ? <>
              <VSpacing size={`px(-2)} />
              <Text value="TRANSACTION" size=Text.Xs color=Colors.gray5 weight=Text.Semibold />
              <VSpacing size=Spacing.xs />
              <Text value={searchTerm ++ "dd92b"} weight=Text.Bold size=Text.Lg block=true />
              <VSpacing size=Spacing.sm />
            </>
          : React.null,
        <> <Text value="Show all results for " /> <Text value=searchTerm weight=Text.Bold /> </>,
      |]
      ->Belt.Array.keep(r => r != React.null);

    <div className=Styles.container>
      {results
       ->Belt.Array.mapWithIndex((i, result) =>
           <div
             onMouseOver={_evt => onHover(i)}
             key={i |> string_of_int}
             className={Css.merge([
               Styles.result,
               i == results->Array.length - 1 ? Styles.lastResult : "",
               i == focusIndex mod results->Array.length ? Styles.resultFocused : "",
             ])}>
             result
           </div>
         )
       ->React.array}
    </div>;
  };
};

type resultState =
  | Hidden
  | ShowAndFocus(int);

type validArrowDirection =
  | Up
  | Down;

type state = {
  searchTerm: string,
  resultState,
};

type action =
  | ChangeSearchTerm(string)
  | ArrowPressed(validArrowDirection)
  | StartTyping
  | StopTyping
  | HoverResultAt(int);

let reducer = state =>
  fun
  | ChangeSearchTerm(newTerm) => {...state, searchTerm: newTerm}
  | ArrowPressed(direction) =>
    switch (state.resultState) {
    | Hidden => state
    | ShowAndFocus(focusIndex) => {
        ...state,
        resultState:
          ShowAndFocus(
            switch (direction) {
            | Up => focusIndex - 1
            | Down => focusIndex + 1
            },
          ),
      }
    }
  | StartTyping => {...state, resultState: ShowAndFocus(0)}
  | StopTyping => {...state, resultState: Hidden}
  | HoverResultAt(resultIndex) => {...state, resultState: ShowAndFocus(resultIndex)};

[@react.component]
let make = () => {
  let ({searchTerm, resultState}, dispatch) =
    React.useReducer(reducer, {searchTerm: "", resultState: Hidden});

  <div className=Styles.container>
    <input
      onFocus={_evt => dispatch(StartTyping)}
      onBlur={_evt => dispatch(StopTyping)}
      onChange={evt => dispatch(ChangeSearchTerm(ReactEvent.Form.target(evt)##value))}
      onKeyDown={event =>
        switch (ReactEvent.Keyboard.key(event)) {
        | "ArrowUp" =>
          dispatch(ArrowPressed(Up));
          ReactEvent.Keyboard.preventDefault(event);
        | "ArrowDown" =>
          dispatch(ArrowPressed(Down));
          ReactEvent.Keyboard.preventDefault(event);
        | "Enter" =>
          dispatch(ChangeSearchTerm(""));
          ReactEvent.Keyboard.preventDefault(event);
          Route.redirect(searchTerm |> Route.search);
        | _ => ()
        }
      }
      value=searchTerm
      className=Styles.search
      placeholder="Search Address / TXN Hash / Block"
    />
    {switch (resultState) {
     | ShowAndFocus(_)
     | Hidden => React.null
     }}
    <button
      className=Styles.button
      onClick={_ => {
        Route.redirect(searchTerm |> Route.search);
        dispatch(ChangeSearchTerm(""));
      }}>
      <img src=Images.searchIcon className=Styles.searchIcon />
    </button>
  </div>;
};
