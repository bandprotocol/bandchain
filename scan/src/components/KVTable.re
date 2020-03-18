module Styles = {
  open Css;

  let thead =
    style([
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(2), Css.rgba(0, 0, 0, 0.05))),
      backgroundColor(Colors.gray3),
      marginBottom(`px(1)),
      display(`flex),
      alignItems(`center),
      height(`px(20)),
      paddingLeft(`px(7)),
      paddingRight(`px(7)),
    ]);

  let tbody =
    style([
      boxShadow(Shadow.box(~x=`zero, ~y=`px(2), ~blur=`px(4), Css.rgba(0, 0, 0, 0.08))),
      backgroundColor(Colors.gray1),
      marginBottom(`px(1)),
      display(`flex),
      alignItems(`center),
      height(`px(20)),
      paddingLeft(`px(7)),
      paddingRight(`px(7)),
    ]);

  let valueContainer = style([maxWidth(`px(230))]);
};

[@react.component]
let make = (~header=["KEY", "VALUE"], ~kv) => {
  <>
    <div className=Styles.thead>
      <Row>
        <Col size=1.>
          <Text
            value={header |> Belt_List.getExn(_, 0)}
            size=Text.Xs
            weight=Text.Semibold
            spacing={Text.Em(0.05)}
            height={Text.Px(18)}
            color=Colors.gray6
          />
        </Col>
        <Col size=1.>
          <Text
            value={header |> Belt_List.getExn(_, 1)}
            size=Text.Xs
            weight=Text.Semibold
            spacing={Text.Em(0.05)}
            height={Text.Px(18)}
            color=Colors.gray6
          />
        </Col>
      </Row>
    </div>
    {kv
     ->Belt.List.map(((key, value)) => {
         <div className=Styles.tbody>
           <Row>
             <Col size=1.>
               <Text
                 value=key
                 size=Text.Sm
                 weight=Text.Medium
                 height={Text.Px(18)}
                 nowrap=true
                 ellipsis=true
                 block=true
                 code=true
               />
             </Col>
             <Col size=1.>
               <div className=Styles.valueContainer>
                 <Text
                   value
                   size=Text.Sm
                   weight=Text.Medium
                   height={Text.Px(18)}
                   nowrap=true
                   ellipsis=true
                   block=true
                   code=true
                 />
               </div>
             </Col>
           </Row>
         </div>
       })
     ->Belt.List.toArray
     ->React.array}
  </>;
};
