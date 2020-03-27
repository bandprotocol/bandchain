type field_t =
  | Key(string)
  | Value(string)
  | DataSource(ID.DataSource.t, string);

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

let renderField = field =>
  switch (field) {
  | Key(x) =>
    <Text
      value=x
      size=Text.Sm
      weight=Text.Medium
      height={Text.Px(18)}
      nowrap=true
      ellipsis=true
      block=true
      code=true
    />
  | Value(x) =>
    <div className=Styles.valueContainer>
      <Text
        value=x
        size=Text.Sm
        weight=Text.Medium
        height={Text.Px(18)}
        nowrap=true
        ellipsis=true
        block=true
        code=true
      />
    </div>
  | DataSource(id, name) =>
    <div className=Styles.valueContainer>
      <TypeID.DataSource id />
      <HSpacing size=Spacing.sm />
      <Text value=name weight=Text.Regular spacing={Text.Em(0.02)} code=true />
    </div>
  };

[@react.component]
let make = (~headers=["KEY", "VALUE"], ~rows) => {
  <>
    <div className=Styles.thead>
      <Row>
        {headers
         ->Belt_List.map(header => {
             <Col key=header size=1.>
               <Text
                 value=header
                 size=Text.Xs
                 weight=Text.Semibold
                 spacing={Text.Em(0.05)}
                 height={Text.Px(18)}
                 color=Colors.gray6
               />
             </Col>
           })
         ->Belt_List.toArray
         ->React.array}
      </Row>
    </div>
    {rows
     ->Belt.List.mapWithIndex((i, fields) => {
         <div className=Styles.tbody key={i |> string_of_int}>
           <Row>
             {fields
              ->Belt_List.mapWithIndex((j, field) => {
                  <Col key={j |> string_of_int} size=1.> {renderField(field)} </Col>
                })
              ->Belt_List.toArray
              ->React.array}
           </Row>
         </div>
       })
     ->Belt.List.toArray
     ->React.array}
  </>;
};
