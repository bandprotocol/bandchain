type field_t =
  | Value(string)
  | Values(list(string))
  | DataSource(ID.DataSource.t, string)
  | Block(ID.Block.t)
  | TxHash(Hash.t)
  | Validator(ValidatorSub.Mini.t);

module Styles = {
  open Css;
  let tabletContainer =
    style([
      padding2(~v=`px(8), ~h=`px(24)),
      backgroundColor(Colors.profileBG),
      Media.mobile([padding2(~v=`px(8), ~h=`px(12))]),
    ]);

  let tableSpacing =
    style([padding2(~v=`px(8), ~h=`zero), Media.mobile([padding2(~v=`px(4), ~h=`zero)])]);

  let valueContainer = mw =>
    style([
      maxWidth(`px(mw)),
      minHeight(`px(20)),
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
    ]);
};

let renderField = (field, maxWidth) => {
  switch (field) {
  | Value(v) =>
    <div className={Styles.valueContainer(maxWidth)}>
      <Text value=v nowrap=true ellipsis=true block=true />
    </div>
  | Values(vals) =>
    <div className={CssHelper.flexBox(~direction=`column, ())}>
      {vals
       ->Belt_List.mapWithIndex((i, v) =>
           <div key={i->string_of_int ++ v} className={Styles.valueContainer(maxWidth)}>
             <Text value=v nowrap=true ellipsis=true block=true align=Text.Right />
           </div>
         )
       ->Belt_List.toArray
       ->React.array}
    </div>
  | DataSource(id, name) =>
    <div className={Styles.valueContainer(maxWidth)}>
      <TypeID.DataSource id position=TypeID.Mini />
      <HSpacing size=Spacing.sm />
      <Text
        value=name
        weight=Text.Regular
        spacing={Text.Em(0.02)}
        size=Text.Sm
        height={Text.Px(16)}
      />
    </div>
  | Block(id) =>
    <div className={Styles.valueContainer(maxWidth)}>
      <TypeID.Block id position=TypeID.Mini />
    </div>
  | TxHash(txHash) =>
    <div className={Styles.valueContainer(maxWidth)}>
      <TxLink txHash width=maxWidth size=Text.Sm />
    </div>
  | Validator(validator) =>
    <div className={Styles.valueContainer(maxWidth)}>
      <ValidatorMonikerLink
        size=Text.Sm
        validatorAddress={validator.operatorAddress}
        width={`px(maxWidth)}
        moniker={validator.moniker}
        identity={validator.identity}
      />
    </div>
  };
};

[@react.component]
let make = (~headers=["Key", "Value"], ~rows) => {
  let columnSize = headers |> Belt_List.length > 2 ? Col.Four : Col.Six;
  let valueWidth = Media.isMobile() ? 70 : 480;
  <div className=Styles.tabletContainer>
    <div className=Styles.tableSpacing>
      <Row>
        {headers
         ->Belt_List.mapWithIndex((i, header) => {
             <Col.Grid key={header ++ (i |> string_of_int)} col=columnSize colSm=columnSize>
               <Text value=header weight=Text.Semibold height={Text.Px(18)} color=Colors.gray7 />
             </Col.Grid>
           })
         ->Belt_List.toArray
         ->React.array}
      </Row>
    </div>
    {rows
     ->Belt.List.mapWithIndex((i, row) => {
         <div
           key={"outerRow" ++ (i |> string_of_int)} className={Css.merge([Styles.tableSpacing])}>
           <Row>
             {row
              ->Belt_List.mapWithIndex((j, value) => {
                  <Col.Grid
                    key={"innerRow" ++ (j |> string_of_int)} col=columnSize colSm=columnSize>
                    {renderField(value, valueWidth)}
                  </Col.Grid>
                })
              ->Belt_List.toArray
              ->React.array}
           </Row>
         </div>
       })
     ->Belt.List.toArray
     ->React.array}
  </div>;
};
