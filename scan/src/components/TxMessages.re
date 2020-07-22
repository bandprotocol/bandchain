module Styles = {
  open Css;
  let msgContainer = overflowed =>
    style([
      position(`relative),
      height(overflowed ? `px(60) : `auto),
      overflow(overflowed ? `hidden : `visible),
    ]);
  let showButton =
    style([
      display(`flex),
      backgroundColor(Colors.gray3),
      borderRadius(`px(30)),
      width(`px(60)),
      alignItems(`center),
      justifyContent(`center),
      fontSize(`px(10)),
      cursor(`pointer),
      color(Colors.gray8),
      height(`px(20)),
    ]);
  let showContainer = style([display(`flex), marginTop(`px(10))]);
};

[@react.component]
let make = (~txHash: Hash.t, ~messages, ~width: int, ~success: bool, ~errMsg: string) => {
  let (overflowed, setOverflowed) = React.useState(_ => false);
  let (expanded, setExpanded) = React.useState(_ => false);

  let msgEl = React.useRef(Js.Nullable.null);

  React.useEffect0(_ => {
    msgEl
    ->React.Ref.current
    ->Js.Nullable.toOption
    ->Belt_Option.map(msgRef => {
        let divHeight = ReactDOMRe.domElementToObj(msgRef)##clientHeight;
        divHeight > 60 ? setOverflowed(_ => true) : ();
      })
    ->ignore;
    None;
  });
  <>
    <div ref={ReactDOMRe.Ref.domRef(msgEl)} className={Styles.msgContainer(overflowed)}>
      {messages
       ->Belt_List.toArray
       ->Belt_Array.mapWithIndex((i, msg) =>
           <React.Fragment key={(txHash |> Hash.toHex) ++ (i |> string_of_int)}>
             <VSpacing size=Spacing.sm />
             <Msg msg width />
             <VSpacing size=Spacing.sm />
           </React.Fragment>
         )
       ->React.array}
      {success ? React.null : <TxError.Mini msg=errMsg />}
    </div>
    {overflowed || expanded
       ? <div>
           <div
             className=Styles.showContainer
             onClick={_ => {
                setOverflowed(_ => !overflowed);
                setExpanded(_ => !expanded);
             }}>
             {expanded
                ? <div className=Styles.showButton> {"show less" |> React.string} </div>
                : Media.isMobile() ? <Link route=Route.TxIndexPage(txHash) > "Show More" |> React.string </Link>
                    : <div className=Styles.showButton> {"show more" |> React.string} </div> }

           </div>
         </div>
       : React.null}
  </>;
};
