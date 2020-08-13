module Styles = {
  open Css;

  // Sort Dropdown

  let sortDrowdownContainer =
    style([position(`relative), zIndex(2), flexBasis(`percent(40.))]);
  let sortDrowdownPanel = show => {
    style([
      display(
        {
          show ? `block : `none;
        },
      ),
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      backgroundColor(Colors.white),
      position(`absolute),
      right(`zero),
      top(`percent(100.)),
      width(`px(165)),
    ]);
  };
  let sortDropdownItem = isActive => {
    style([
      backgroundColor(
        {
          isActive ? Colors.blue1 : Colors.white;
        },
      ),
      cursor(`pointer),
      display(`flex),
      alignItems(`center),
      padding2(~v=`px(8), ~h=`px(10)),
      selector("> img", [marginRight(`px(5))]),
    ]);
  };
  let sortDropdownTextItem = {
    style([
      paddingRight(`px(15)),
      cursor(`pointer),
      after([
        contentRule(`text("")),
        backgroundImage(`url(Images.sortDown)),
        width(`px(8)),
        height(`px(8)),
        backgroundRepeat(`noRepeat),
        backgroundSize(`contain),
        display(`block),
        position(`absolute),
        top(`percent(50.)),
        right(`zero),
        transform(translateY(`percent(-50.))),
      ]),
    ]);
  };
};

[@react.component]
let make = (~sortedBy, ~setSortedBy, ~sortList) => {
  let (show, setShow) = React.useState(_ => false);
  <div className=Styles.sortDrowdownContainer>
    <div className=Styles.sortDropdownTextItem onClick={_ => setShow(prev => !prev)}>
      <Text
        block=true
        value="Sort By"
        size=Text.Md
        weight=Text.Regular
        color=Colors.gray6
        align=Text.Right
      />
    </div>
    <div className={Styles.sortDrowdownPanel(show)}>
      {sortList
       ->Belt.List.mapWithIndex((i, (value, name)) => {
           let isActive = sortedBy == value;
           <div
             key={i |> string_of_int}
             className={Styles.sortDropdownItem(isActive)}
             onClick={_ => {
               setSortedBy(_ => value);
               setShow(_ => false);
             }}>
             <Text
               block=true
               value=name
               size=Text.Md
               weight=Text.Regular
               color={isActive ? Colors.blue7 : Colors.gray6}
             />
           </div>;
         })
       ->Array.of_list
       ->React.array}
    </div>
  </div>;
};
