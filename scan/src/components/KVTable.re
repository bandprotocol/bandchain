type field_t =
  | Value(string)
  | Values(list(string))
  | DataSource(ID.DataSource.t, string)
  | Block(ID.Block.t)
  | TxHash(Hash.t)
  | Validator(ValidatorSub.Mini.t);

type theme_t =
  | MessageMiniTable
  | RequestMiniTable;

type with_setting_t('a) = {
  mainElem: 'a,
  size: float,
  isRight: bool,
};

module Styles = {
  open Css;

  let vFlex = style([display(`flex), flexDirection(`column)]);
  let hFlex = style([display(`flex), alignItems(`center)]);

  let thead = theme =>
    style([
      boxShadow(
        Shadow.box(
          ~x=`zero,
          ~y=`px(2),
          ~blur=`px(2),
          switch (theme) {
          | MessageMiniTable => Css.rgba(0, 0, 0, 0.05)
          | RequestMiniTable => Css.rgba(11, 29, 142, 0.05)
          },
        ),
      ),
      backgroundColor(
        switch (theme) {
        | MessageMiniTable => Colors.gray3
        | RequestMiniTable => Colors.blue1
        },
      ),
      marginBottom(`px(1)),
      display(`flex),
      alignItems(`center),
      height(
        switch (theme) {
        | MessageMiniTable => `px(20)
        | RequestMiniTable => `px(25)
        },
      ),
      paddingLeft(`px(7)),
      paddingRight(`px(7)),
    ]);

  let tbody = theme =>
    style([
      boxShadow(
        Shadow.box(
          ~x=`zero,
          ~y=`px(2),
          ~blur=`px(4),
          switch (theme) {
          | MessageMiniTable => Css.rgba(0, 0, 0, 0.08)
          | RequestMiniTable => Css.rgba(11, 29, 142, 0.08)
          },
        ),
      ),
      backgroundColor(
        switch (theme) {
        | MessageMiniTable => Colors.gray1
        | RequestMiniTable => Colors.blueGray1
        },
      ),
      marginBottom(`px(1)),
      display(`flex),
      padding2(
        ~v=
          switch (theme) {
          | MessageMiniTable => `px(1)
          | RequestMiniTable => `px(5)
          },
        ~h=`px(7),
      ),
    ]);

  let valueContainer = mw =>
    style([
      maxWidth(`px(mw)),
      minHeight(`px(20)),
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
    ]);
  let fillRight = style([marginRight(`auto)]);
};

let renderField = (field, maxWidth, isRight) => {
  switch (field) {
  | Value(v) =>
    <div className={Styles.valueContainer(maxWidth)}>
      <Text
        value=v
        size=Text.Sm
        weight=Text.Medium
        height={Text.Px(18)}
        nowrap=true
        ellipsis=true
        block=true
        code=true
      />
    </div>
  | Values(vals) =>
    <div className=Styles.vFlex>
      {vals
       ->Belt_List.map(v =>
           <div key=v className={Styles.valueContainer(maxWidth)}>
             {isRight ? <div className=Styles.fillRight /> : React.null}
             <Text
               value=v
               size=Text.Sm
               weight=Text.Medium
               height={Text.Px(18)}
               nowrap=true
               ellipsis=true
               block=true
               code=true
             />
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
      />
    </div>
  };
};

let withSetting = (arr, sizes, isRights) =>
  arr->Belt_List.mapWithIndex((i, elem) =>
    {
      mainElem: elem,
      size: sizes->Belt_List.get(i)->Belt_Option.getWithDefault(1.),
      isRight: isRights->Belt_List.get(i)->Belt_Option.getWithDefault(false),
    }
  );

[@react.component]
let make =
    (
      ~tableWidth,
      ~headers=["KEY", "VALUE"],
      ~rows,
      ~sizes=[],
      ~isRights=[],
      ~theme=MessageMiniTable,
    ) => {
  let headersWithSetting = headers->withSetting(sizes, isRights);
  let rowsWithSetting = rows->Belt_List.map(fields => fields->withSetting(sizes, isRights));
  <>
    <div className={Styles.thead(theme)}>
      <Row>
        {headersWithSetting
         ->Belt_List.mapWithIndex((i, {mainElem, size, isRight}) => {
             <Col key={i |> string_of_int} size>
               <div className=Styles.hFlex>
                 {isRight ? <div className=Styles.fillRight /> : React.null}
                 <Text
                   value=mainElem
                   size={
                     switch (theme) {
                     | MessageMiniTable => Text.Xs
                     | RequestMiniTable => Text.Sm
                     }
                   }
                   weight=Text.Semibold
                   spacing={Text.Em(0.05)}
                   height={Text.Px(18)}
                   color={
                     switch (theme) {
                     | MessageMiniTable => Colors.gray6
                     | RequestMiniTable => Colors.bandBlue
                     }
                   }
                 />
               </div>
             </Col>
           })
         ->Belt_List.toArray
         ->React.array}
      </Row>
    </div>
    {let sumSizes =
       switch (sizes |> Belt_List.length) {
       | 0 => headers |> Belt_List.length |> float_of_int
       | _ => sizes->Belt_List.reduce(0., (+.))
       };
     rowsWithSetting
     ->Belt.List.mapWithIndex((i, rowWithSetting) => {
         <div key={i |> string_of_int} className={Styles.tbody(theme)}>
           {rowWithSetting
            ->Belt_List.mapWithIndex((j, {mainElem, size, isRight}) => {
                <Col key={j |> string_of_int} size>
                  <div className=Styles.hFlex>
                    {isRight ? <div className=Styles.fillRight /> : React.null}
                    {renderField(
                       mainElem,
                       sumSizes <= 0.
                         ? tableWidth
                         : (tableWidth |> float_of_int) *. size /. sumSizes |> int_of_float,
                       isRight,
                     )}
                  </div>
                </Col>
              })
            ->Belt_List.toArray
            ->React.array}
         </div>
       })
     ->Belt.List.toArray
     ->React.array}
  </>;
};
