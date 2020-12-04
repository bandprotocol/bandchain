module Styles = {
  open Css;
  let container =
    style([
      width(`px(468)),
      minHeight(`px(500)),
      padding3(~top=`px(32), ~h=`px(24), ~bottom=`px(24)),
      borderRadius(`px(4)),
    ]);

  let modalTitle = style([paddingBottom(`px(24))]);

  let resultContainer = style([minHeight(`px(400)), width(`percent(100.))]);

  let btn = style([width(`percent(100.))]);

  let btnBack =
    style([
      display(`table),
      cursor(`pointer),
      margin3(~top=`px(24), ~h=`auto, ~bottom=`zero),
    ]);

  let jsonDisplay =
    style([
      resize(`none),
      fontSize(`px(12)),
      backgroundColor(Colors.bg),
      border(`px(1), `solid, Colors.gray9),
      borderRadius(`px(4)),
      width(`percent(100.)),
      height(`px(300)),
      overflowY(`scroll),
      marginBottom(`px(16)),
      fontFamilies([
        `custom("IBM Plex Mono"),
        `custom("cousine"),
        `custom("sfmono-regular"),
        `custom("Consolas"),
        `custom("Menlo"),
        `custom("liberation mono"),
        `custom("ubuntu mono"),
        `custom("Courier"),
        `monospace,
      ]),
    ]);

  let resultIcon = style([width(`px(48)), marginBottom(`px(16))]);

  let txhashContainer = style([cursor(`pointer)]);
};

type state_t =
  | Nothing
  | Signing
  | Broadcasting
  | Success(Hash.t)
  | Error(string);

[@react.component]
let make = (~rawTx, ~onBack, ~account: AccountContext.t) => {
  let (_, dispatchModal) = React.useContext(ModalContext.context);
  let (state, setState) = React.useState(_ => Nothing);
  let jsonTx = TxCreator.sortAndStringify(rawTx);

  <div className=Styles.container>
    <div className=Styles.modalTitle>
      <Text value="Confirm Transaction" weight=Text.Medium size=Text.Xl />
    </div>
    {switch (state) {
     | Nothing =>
       <div>
         <div className={CssHelper.mb(~size=16, ())}>
           <Text
             value="Please verify the transaction details below before proceeding"
             size=Text.Lg
           />
         </div>
         <textarea
           className=Styles.jsonDisplay
           disabled=true
           defaultValue={rawTx |> TxCreator.stringifyWithSpaces}
         />
         <div id="broadcastButtonContainer">
           <Button
             py=10
             style=Styles.btn
             onClick={_ => {
               dispatchModal(DisableExit);
               setState(_ => Signing);
               let _ =
                 Wallet.sign(jsonTx, account.wallet)
                 |> Js.Promise.then_(signature => {
                      setState(_ => Broadcasting);
                      let signedTx =
                        TxCreator.createSignedTx(
                          ~network=Env.network,
                          ~signature=signature |> JsBuffer.toBase64,
                          ~pubKey=account.pubKey,
                          ~tx=rawTx,
                          ~mode="sync",
                          (),
                        );
                      ignore(
                        TxCreator.broadcast(signedTx)
                        |> Js.Promise.then_(res =>
                             switch (res) {
                             | TxCreator.Tx(txResponse) =>
                               txResponse.success
                                 ? {
                                   setState(_ => Success(txResponse.txHash));
                                 }
                                 : {
                                   Js.Console.error(txResponse);
                                   setState(_ => Error(txResponse.code |> TxResError.parse));
                                 };
                               dispatchModal(EnableExit);
                               Js.Promise.resolve();
                             | _ =>
                               setState(_ => Error("Fail to broadcast"));
                               dispatchModal(EnableExit);
                               Js.Promise.resolve();
                             }
                           )
                        |> Js.Promise.catch(err => {
                             switch (Js.Json.stringifyAny(err)) {
                             | Some(errorValue) => setState(_ => Error(errorValue))
                             | None => setState(_ => Error("Can not stringify error"))
                             };
                             dispatchModal(EnableExit);
                             Js.Promise.resolve();
                           }),
                      );

                      Promise.ret();
                    })
                 |> Js.Promise.catch(_ => {
                      setState(_ => Error("Failed to sign message"));
                      dispatchModal(EnableExit);
                      Promise.ret();
                    });
               ();
             }}>
             <Text value="Broadcast" weight=Text.Semibold size=Text.Lg color=Colors.white />
           </Button>
         </div>
         <a className=Styles.btnBack onClick=onBack>
           <Text value="Back" weight=Text.Semibold size=Text.Lg color=Colors.gray7 />
         </a>
       </div>
     | Success(txHash) =>
       <div
         className={Css.merge([
           CssHelper.flexBox(~direction=`column, ~justify=`center, ()),
           Styles.resultContainer,
         ])}>
         <img src=Images.success className=Styles.resultIcon />
         <div id="successMsgContainer" className={CssHelper.mb(~size=16, ())}>
           <Text value="Broadcast transaction success" size=Text.Lg block=true align=Text.Center />
         </div>
         <Link className=Styles.txhashContainer route={Route.TxIndexPage(txHash)}>
           <Button py=8 px=13 variant=Button.Outline onClick={_ => {dispatchModal(CloseModal)}}>
             <Text
               block=true
               value="View Details"
               weight=Text.Semibold
               ellipsis=true
               color=Colors.bandBlue
             />
           </Button>
         </Link>
       </div>
     | Signing =>
       <div
         className={Css.merge([
           CssHelper.flexBox(~direction=`column, ~justify=`center, ()),
           Styles.resultContainer,
         ])}>
         <div className={CssHelper.mb(~size=16, ())}>
           <Icon name="fad fa-spinner-third fa-spin" color=Colors.bandBlue size=48 />
         </div>
         <Text value="Waiting for signing transaction" size=Text.Lg block=true align=Text.Center />
       </div>
     | Broadcasting =>
       <div
         className={Css.merge([
           CssHelper.flexBox(~direction=`column, ~justify=`center, ()),
           Styles.resultContainer,
         ])}>
         <div className={CssHelper.mb(~size=16, ())}>
           <Icon name="fad fa-spinner-third fa-spin" color=Colors.bandBlue size=48 />
         </div>
         <Text
           value="Waiting for broadcasting transaction"
           size=Text.Lg
           block=true
           align=Text.Center
         />
       </div>
     | Error(err) =>
       <div
         className={Css.merge([
           CssHelper.flexBox(~direction=`column, ~justify=`center, ()),
           Styles.resultContainer,
         ])}>
         <img src=Images.fail className=Styles.resultIcon />
         <div className={CssHelper.mb()}>
           <Text value="Broadcast transaction fail" size=Text.Lg block=true align=Text.Center />
         </div>
         <Text value=err color=Colors.red3 align=Text.Center breakAll=true />
       </div>
     }}
  </div>;
};
