module Styles = {
  open Css;

  let typeContainer = w => style([marginRight(`px(20)), width(w)]);

  let resolveIcon = style([width(`px(20)), height(`px(20)), marginLeft(Spacing.sm)]);

  let hashContainer = style([maxWidth(`px(220))]);
  let feeContainer = style([display(`flex), justifyContent(`flexEnd)]);
  let timeContainer = style([display(`flex), alignItems(`center), maxWidth(`px(150))]);
  let textContainer = style([display(`flex)]);
  let countContainer = style([maxWidth(`px(80))]);
  let proposerBox = style([maxWidth(`px(270)), display(`flex), flexDirection(`column)]);
  let idContainer = style([display(`flex), maxWidth(`px(200))]);
  let dataSourcesContainer = style([display(`flex)]);
  let dataSourceContainer = style([display(`flex), width(`px(170))]);
  let oracleScriptContainer = style([display(`flex), width(`px(170))]);
  let resolveStatusContainer =
    style([display(`flex), alignItems(`center), justifyContent(`flexEnd)]);
};

let renderText = (text, weight) =>
  <div className={Styles.typeContainer(`px(150))}>
    <Text value=text size=Text.Lg weight block=true ellipsis=true />
  </div>;

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
  <div className=Styles.timeContainer> <Timestamp time size=Text.Md /> </div>;
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
  <div className=Styles.dataSourceContainer>
    <TypeID.DataSource id position=TypeID.Text />
    <HSpacing size=Spacing.xs />
    <Text
      value=name
      block=true
      height={Text.Px(16)}
      spacing={Text.Em(0.02)}
      nowrap=true
      ellipsis=true
    />
  </div>;
};

let renderOracleScript = (id, name) => {
  <div className=Styles.oracleScriptContainer>
    <TypeID.OracleScript id position=TypeID.Text />
    <HSpacing size=Spacing.xs />
    <Text
      value=name
      block=true
      height={Text.Px(16)}
      spacing={Text.Em(0.02)}
      ellipsis=true
      nowrap=true
    />
  </div>;
};

let renderRelatedDataSources = ids => {
  switch (ids |> Belt_List.length) {
  | 0 => <Text value="TBD" size=Text.Md spacing={Text.Em(0.02)} />
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

let renderRequest = id => {
  <div className=Styles.idContainer> <TypeID.Request id position=TypeID.Text /> </div>;
};

let renderRequestStatus = status => {
  <div className=Styles.resolveStatusContainer>
    <Text
      block=true
      size=Text.Md
      weight=Text.Medium
      align=Text.Right
      value={
        switch (status) {
        | RequestStatus.Success => "Success"
        | Failure => "Fail"
        | Pending => "Pending"
        | Expired => "Expired"
        | Unknown => "???"
        }
      }
    />
    <img
      src={
        switch (status) {
        | Success => Images.success
        | Failure => Images.fail
        | Pending => Images.pending
        | Expired => Images.expired
        | Unknown => Images.unknown
        }
      }
      className=Styles.resolveIcon
    />
  </div>;
};

type t =
  | Height(int)
  | HeightWithTime(int, MomentRe.Moment.t)
  | Name(string)
  | Timestamp(MomentRe.Moment.t)
  | TxHash(Hash.t, MomentRe.Moment.t)
  | Detail(string)
  | Status(string)
  | Count(int)
  | Fee(float)
  | Hash(Hash.t)
  | Address(Address.t)
  | Value(Js.Json.t)
  | Proposer(string, string)
  | DataSource(ID.DataSource.t, string)
  | OracleScript(ID.OracleScript.t, string)
  | RelatedDataSources(list(ID.DataSource.t))
  | Request(ID.Request.t)
  | RequestStatus(RequestStatus.t);

[@react.component]
let make = (~elementType) => {
  switch (elementType) {
  | Height(height) => renderHeight(height)
  | HeightWithTime(height, timestamp) => renderHeightWithTime(height, timestamp)
  | Name(name) => renderName(name)
  | Timestamp(time) => renderTime(time)
  | TxHash(hash, timestamp) => renderTxHash(hash, timestamp)
  | Detail(detail) => renderText(detail, Text.Semibold)
  | Status(status) => renderText(status, Text.Semibold)
  | Count(count) => renderCount(count)
  | Fee(fee) => renderFee(fee)
  | Hash(hash) => renderHash(hash)
  | Address(address) => renderAddress(address)
  | Value(value) => renderText(value->Js.Json.stringify, Text.Regular)
  | Proposer(moniker, proposer) => renderProposer(moniker, proposer)
  | DataSource(id, name) => renderDataSource(id, name)
  | OracleScript(id, name) => renderOracleScript(id, name)
  | RelatedDataSources(ids) => renderRelatedDataSources(ids)
  | Request(id) => renderRequest(id)
  | RequestStatus(status) => renderRequestStatus(status)
  };
};
