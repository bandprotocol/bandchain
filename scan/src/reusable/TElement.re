module Styles = {
  open Css;

  let typeContainer = w => style([marginRight(`px(20)), width(w)]);

  let msgIcon =
    style([
      width(`px(30)),
      height(`px(30)),
      marginTop(`px(5)),
      marginLeft(Spacing.xl),
      marginRight(Spacing.xl),
    ]);

  let hashContainer = style([maxWidth(`px(220))]);
  let feeContainer = style([display(`flex), justifyContent(`flexEnd)]);
  let timeContainer = style([display(`flex), alignItems(`center), maxWidth(`px(150))]);
  let textContainer = style([display(`flex)]);
  let countContainer = style([maxWidth(`px(80))]);
  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
  let idContainer = style([display(`flex)]);
  let dataSourcesContainer = style([display(`flex)]);
};

let renderText = (text, weight) =>
  <div className={Styles.typeContainer(`px(150))}>
    <Text value=text size=Text.Lg weight block=true ellipsis=true />
  </div>;

let renderTxTypeWithDetail = (msgs: list(TxHook.Msg.t)) => {
  <div className={Styles.typeContainer(`px(150))}>
    <MsgBadge msgs />
    <VSpacing size=Spacing.xs />
    <Text
      value={msgs->Belt.List.getExn(0)->TxHook.Msg.getDescription}
      size=Text.Lg
      weight=Text.Semibold
      block=true
      ellipsis=true
    />
  </div>;
};

let renderTxHash = (hash, time) => {
  <div className=Styles.hashContainer>
    <TimeAgos time />
    <VSpacing size={`px(6)} />
    <Text
      block=true
      code=true
      value={hash |> Hash.toHex(~with0x=true)}
      size=Text.Lg
      weight=Text.Bold
      ellipsis=true
    />
  </div>;
};

let renderHash = hash => {
  <div className=Styles.hashContainer>
    <Text
      block=true
      code=true
      value={hash |> Hash.toHex}
      size=Text.Lg
      weight=Text.Bold
      ellipsis=true
    />
  </div>;
};

let renderHashWithLink = hash => {
  <div className=Styles.hashContainer onClick={_ => hash->Route.TxIndexPage->Route.redirect}>
    <Text
      block=true
      code=true
      value={hash |> Hash.toHex}
      size=Text.Lg
      weight=Text.Bold
      ellipsis=true
      color=Colors.purple3
    />
  </div>;
};

let renderAddress = address => {
  <AddressRender address />;
};

let renderFee = fee => {
  <div className=Styles.feeContainer>
    <Text value={fee->Format.fPretty} color=Colors.gray7 code=true block=true nowrap=true />
  </div>;
};

let renderHeight = blockHeight => {
  <div className=Styles.textContainer> <TypeID.Block id={ID.Block.ID(blockHeight)} /> </div>;
};

let renderHeightWithTime = (height, time) => {
  <> <TimeAgos time /> <VSpacing size={`px(6)} /> {renderHeight(height)} </>;
};

let renderName = name => {
  <div className=Styles.hashContainer>
    <Text block=true code=true value=name size=Text.Lg weight=Text.Bold ellipsis=true />
  </div>;
};

let renderTime = time => {
  <div className=Styles.timeContainer> <TimeAgos time size=Text.Md /> </div>;
};

let renderCount = count => {
  <div className=Styles.countContainer>
    <Text value={count |> string_of_int} size=Text.Md weight=Text.Semibold />
  </div>;
};

let renderProposer = (moniker, proposer) => {
  <div className=Styles.proposerBox>
    <Text block=true value=moniker size=Text.Sm weight=Text.Regular color=Colors.gray7 />
    <VSpacing size=Spacing.sm />
    <Text
      block=true
      value=proposer
      size=Text.Md
      weight=Text.Bold
      code=true
      ellipsis=true
      color=Colors.gray8
    />
  </div>;
};

let renderDataSource = (id, name) => {
  <div className=Styles.idContainer>
    <TypeID.DataSource id position=TypeID.Text />
    <HSpacing size=Spacing.xs />
    <Text value=name block=true height={Text.Px(16)} spacing={Text.Em(0.02)} />
  </div>;
};

let renderOracleScript = (id, name) => {
  <div className=Styles.idContainer>
    <TypeID.OracleScript id position=TypeID.Text />
    <HSpacing size=Spacing.xs />
    <Text value=name block=true height={Text.Px(16)} spacing={Text.Em(0.02)} />
  </div>;
};

let renderRelatedDataSources = ids => {
  switch (ids |> Belt_List.length) {
  | 0 => <Text value="Undetermined" size=Text.Md spacing={Text.Em(0.02)} />
  | _ =>
    <div className=Styles.dataSourcesContainer>
      {ids
       ->Belt_List.map(id => {
           <div className=Styles.idContainer key={id |> ID.DataSource.toString}>
             <TypeID.DataSource id position=TypeID.Text />
             <HSpacing size=Spacing.sm />
           </div>
         })
       ->Array.of_list
       ->React.array}
    </div>
  };
};

let msgIcon =
  fun
  | TxHook.Msg.CreateDataSource(_) => Images.newScript
  | EditDataSource(_) => Images.newScript
  | CreateOracleScript(_) => Images.newScript
  | EditOracleScript(_) => Images.newScript
  | Send(_) => Images.sendCoin
  | Request(_) => Images.dataRequest
  | Report(_) => Images.report
  | Unknown => Images.checkIcon;

type t =
  | Icon(TxHook.Msg.t)
  | Height(int)
  | HeightWithTime(int, MomentRe.Moment.t)
  | Name(string)
  | Timestamp(MomentRe.Moment.t)
  | TxHash(Hash.t, MomentRe.Moment.t)
  | TxTypeWithDetail(list(TxHook.Msg.t))
  | Detail(string)
  | Status(string)
  | Count(int)
  | Fee(float)
  | Hash(Hash.t)
  | HashWithLink(Hash.t)
  | Address(Address.t)
  | Value(Js.Json.t)
  | Proposer(string, string)
  | DataSource(ID.DataSource.t, string)
  | OracleScript(ID.OracleScript.t, string)
  | RelatedDataSources(list(ID.DataSource.t));

[@react.component]
let make = (~elementType) => {
  switch (elementType) {
  | Icon({action, _}) => <img src={action->msgIcon} className=Styles.msgIcon />
  | Height(height) => renderHeight(height)
  | HeightWithTime(height, timestamp) => renderHeightWithTime(height, timestamp)
  | Name(name) => renderName(name)
  | Timestamp(time) => renderTime(time)
  | TxHash(hash, timestamp) => renderTxHash(hash, timestamp)
  | TxTypeWithDetail(msgs) => renderTxTypeWithDetail(msgs)
  | Detail(detail) => renderText(detail, Text.Semibold)
  | Status(status) => renderText(status, Text.Semibold)
  | Count(count) => renderCount(count)
  | Fee(fee) => renderFee(fee)
  | Hash(hash) => renderHash(hash)
  | HashWithLink(hash) => renderHashWithLink(hash)
  | Address(address) => renderAddress(address)
  | Value(value) => renderText(value->Js.Json.stringify, Text.Regular)
  | Proposer(moniker, proposer) => renderProposer(moniker, proposer)
  | DataSource(id, name) => renderDataSource(id, name)
  | OracleScript(id, name) => renderOracleScript(id, name)
  | RelatedDataSources(ids) => renderRelatedDataSources(ids)
  };
};
