module Styles = {
  open Css;
  let msgContainer = overflowed =>
    style([
      position(`relative),
      height(overflowed ? `px(60) : `auto),
      overflow(overflowed ? `hidden : `visible),
      selector("> div + div", [marginTop(`px(10))]),
      Media.mobile([height(overflowed ? `px(50) : `auto)]),
    ]);
  let showButton =
    style([
      display(`flex),
      backgroundColor(Colors.gray3),
      borderRadius(`px(30)),
      width(`px(65)),
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

  let isMobile = Media.isMobile();

  let msgCount = isMobile ? 1 : 2;

  React.useEffect0(_ => {
    let msgLength = Belt.List.length(messages);
    msgLength > msgCount ? setOverflowed(_ => true) : ();
    None;
  });
  <>
    <div ref={ReactDOMRe.Ref.domRef(msgEl)} className={Styles.msgContainer(overflowed)}>
      {messages
       ->Belt_List.toArray
       ->Belt_Array.mapWithIndex((i, msg) =>
           <React.Fragment key={(txHash |> Hash.toHex) ++ (i |> string_of_int)}>
             <Msg msg width />
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
                : isMobile
                    ? <Link className=Styles.showButton route={Route.TxIndexPage(txHash)}>
                        {"show more" |> React.string}
                      </Link>
                    : <div className=Styles.showButton> {"show more" |> React.string} </div>}
           </div>
         </div>
       : React.null}
  </>;
};
