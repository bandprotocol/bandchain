module Styles = {
  open Css;
  let cardContainer =
    style([
      position(`relative),
      padding(`px(10)),
      backgroundColor(Colors.white),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      selector("+ div", [marginTop(`px(10))]),
    ]);
  let cardItem = (alignItems_, isOneColumn) =>
    style([display(isOneColumn ? `block : `flex), alignItems(alignItems_), padding(`px(5))]);
  let cardItemHeading =
    style([
      display(`flex),
      flexDirection(`column),
      flexGrow(0.),
      flexShrink(0.),
      flexBasis(`percent(25.)),
    ]);
  let logo = style([width(`px(20)), position(`absolute), top(`px(5)), right(`px(8))]);
  let cardItemHeadingLg = style([padding2(~v=`px(10), ~h=`zero)]);
  let infoContainer = isOneColumn =>
    style([
      width(`percent(100.)),
      marginTop(isOneColumn ? `px(10) : `zero),
      overflow(`hidden),
    ]);
};

[@react.component]
let make = (~values, ~idx, ~status=?, ~requestStatus=?) => {
  <div className=Styles.cardContainer>
    {switch (status) {
     | Some(success) => <img src={success ? Images.success : Images.fail} className=Styles.logo />
     | None => React.null
     }}
    {switch (requestStatus) {
     | Some(resolveStatus) =>
       switch (resolveStatus) {
       | RequestSub.Success => <img src=Images.success className=Styles.logo />
       | Failure => <img src=Images.fail className=Styles.logo />
       | Pending => <img src=Images.pending className=Styles.logo />
       | Expired => <img src=Images.expired className=Styles.logo />
       | Unknown => <img src=Images.unknown className=Styles.logo />
       }
     | None => React.null
     }}
    {values
     ->Belt_List.mapWithIndex((index, (heading, value)) => {
         let alignItem =
           switch (value) {
           | InfoMobileCard.Messages(_)
           | PubKey(_) => `baseline
           | _ => `center
           };
         let isOneColumn =
           switch (value) {
           | InfoMobileCard.KVTableReport(_)
           | KVTableRequest(_) => true
           | _ => false
           };
         <div
           className={Styles.cardItem(alignItem, isOneColumn)}
           key={idx ++ (index |> string_of_int)}>
           <div className=Styles.cardItemHeading>
             {heading
              ->Js.String2.split("\n")
              ->Belt.Array.map(each => {
                  switch (value) {
                  | InfoMobileCard.Nothing =>
                    <div className=Styles.cardItemHeadingLg>
                      <Text key=each value=each size=Text.Sm weight=Text.Bold color=Colors.gray6 />
                    </div>
                  | _ =>
                    <Text
                      key=each
                      value=each
                      size=Text.Sm
                      weight=Text.Semibold
                      color=Colors.gray6
                    />
                  }
                })
              ->React.array}
           </div>
           <div className={Styles.infoContainer(isOneColumn)}>
             <InfoMobileCard info=value />
           </div>
         </div>;
       })
     ->Belt.List.toArray
     ->React.array}
  </div>;
};
