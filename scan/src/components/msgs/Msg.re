module Styles = {
  open Css;
  let withWidth = (w: int) => style([width(`px(w)), display(`flex), flexDirection(`row)]);
  let withBg = (color: Types.Color.t) =>
    style([height(`px(16)), backgroundColor(color), borderRadius(`px(100))]);
};

[@react.component]
let make = (~msg: TxHook.Msg.t, ~width: int) => {
  switch (msg.action) {
  | Send({fromAddress, toAddress, amount}) =>
    <div className={Styles.withWidth(width)}>
      <div className={Styles.withWidth(width / 2 - 18)}>
        <Text value="band" weight=Text.Semibold height={Text.Px(16)} code=true />
        <Text
          value={fromAddress |> Address.toBech32 |> Js.String.sliceToEnd(~from=4)}
          weight=Text.Regular
          height={Text.Px(16)}
          ellipsis=true
          nowrap=true
          block=true
          code=true
        />
      </div>
      <div className={Styles.withBg(Colors.blue1)}>
        <Text value="SEND" color=Colors.blue7 />
      </div>
      <div className={Styles.withWidth(width / 2 - 18)}>
        <Text value="band" weight=Text.Semibold height={Text.Px(16)} code=true />
        <Text
          value={toAddress |> Address.toBech32 |> Js.String.sliceToEnd(~from=4)}
          weight=Text.Regular
          height={Text.Px(16)}
          ellipsis=true
          nowrap=true
          block=true
          code=true
        />
      </div>
    </div>
  | CreateDataSource(cds) => React.null
  | EditDataSource(eds) => React.null
  | CreateOracleScript(cos) => React.null
  | EditOracleScript(eos) => React.null
  | Request(req) => React.null
  | Report(rep) => React.null
  | Unknown => React.null
  };
};
