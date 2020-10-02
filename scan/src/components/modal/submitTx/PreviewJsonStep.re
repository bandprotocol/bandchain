module Styles = {
  open Css;
  let container =
    style([
      display(`flex),
      flexDirection(`column),
      width(`px(640)),
      minHeight(`px(480)),
      padding4(~top=`px(50), ~bottom=`px(34), ~left=`px(50), ~right=`px(50)),
      backgroundColor(rgb(249, 249, 251)),
      borderRadius(`px(5)),
    ]);

  let modalTitle = style([display(`flex), justifyContent(`center)]);

  let rowContainer =
    style([
      display(`flex),
      alignItems(`center),
      justifyContent(`spaceBetween),
      paddingTop(`px(20)),
      minHeight(`px(70)),
      maxHeight(`px(70)),
    ]);

  let resultContainer =
    style([
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
      justifyContent(`center),
      paddingTop(`px(15)),
      minHeight(`px(70)),
      maxHeight(`px(70)),
    ]);

  let errorContainer =
    style([
      display(`flex),
      flexDirection(`column),
      alignItems(`center),
      justifyContent(`center),
      marginTop(`px(15)),
    ]);

  let errMsg = style([maxHeight(`px(100)), overflowY(`auto)]);

  let rFlex =
    style([
      display(`flex),
      flexDirection(`row),
      alignItems(`center),
      justifyContent(`center),
    ]);

  let btn =
    style([
      width(`px(100)),
      height(`px(30)),
      display(`flex),
      justifySelf(`right),
      justifyContent(`center),
      alignItems(`center),
      backgroundImage(
        `linearGradient((
          `deg(90.),
          [(`percent(0.), Colors.blue7), (`percent(100.), Colors.bandBlue)],
        )),
      ),
      boxShadow(
        Shadow.box(~x=`zero, ~y=`px(4), ~blur=`px(8), Css.rgba(82, 105, 255, `num(0.25))),
      ),
      borderRadius(`px(4)),
      cursor(`pointer),
      alignSelf(`center),
    ]);

  let jsonDisplay =
    style([resize(`none), width(`percent(100.)), height(`px(300)), overflowY(`scroll)]);

  let loading = style([width(`px(100))]);

  let resultIcon = style([width(`px(30))]);

  let txhashContainer = style([cursor(`pointer)]);

  let warning =
    style([
      display(`flex),
      flexDirection(`row),
      padding2(~v=`px(5), ~h=`px(8)),
      color(Colors.yellow6),
      backgroundColor(Colors.yellow1),
      border(`px(1), `solid, Colors.yellow6),
      borderRadius(`px(4)),
    ]);
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
      <Text value="Confirm Transaction" weight=Text.Bold size=Text.Xxxl />
    </div>
    <VSpacing size=Spacing.lg />
    <div className=Styles.warning>
      <Text value="Please verify the transaction details below before proceeding" />
    </div>
    <VSpacing size=Spacing.md />
    <textarea
      className=Styles.jsonDisplay
      disabled=true
      defaultValue={rawTx |> TxCreator.stringifyWithSpaces}
    />
    {switch (state) {
     | Nothing =>
       <div className=Styles.rowContainer>
         <div className=Styles.btn onClick=onBack>
           <Text value="Back" weight=Text.Bold size=Text.Md color=Colors.white />
         </div>
         <div
           className=Styles.btn
           id="broadcastButton"
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
           <Text value="Broadcast" weight=Text.Bold size=Text.Md color=Colors.white />
         </div>
       </div>
     | Success(txHash) =>
       <div className=Styles.resultContainer>
         <div className=Styles.rFlex id="successMsgContainer">
           <img src=Images.success2 className=Styles.resultIcon />
           <HSpacing size=Spacing.md />
           <Text value="Broadcast Transaction Success" weight=Text.Semibold />
         </div>
         <VSpacing size=Spacing.md />
         <div onClick={_ => dispatchModal(CloseModal)}>
           <Link className=Styles.txhashContainer route={Route.TxIndexPage(txHash)}>
             <Text
               block=true
               code=true
               spacing={Text.Em(0.02)}
               value={txHash |> Hash.toHex(~upper=true)}
               weight=Text.Medium
               ellipsis=true
               size=Text.Sm
               color=Colors.gray7
             />
           </Link>
         </div>
       </div>
     | Signing =>
       <div className=Styles.resultContainer>
         <img src=Images.loadingCircles className=Styles.loading />
         <VSpacing size=Spacing.sm />
         <Text
           value="Waiting for signing transaction"
           spacing={Text.Em(0.03)}
           weight=Text.Medium
         />
       </div>
     | Broadcasting =>
       <div className=Styles.resultContainer>
         <img src=Images.loadingCircles className=Styles.loading />
         <VSpacing size=Spacing.sm />
         <Text
           value="Waiting for broadcasting transaction"
           spacing={Text.Em(0.03)}
           weight=Text.Medium
         />
       </div>
     | Error(err) =>
       <div className=Styles.errorContainer>
         <div className=Styles.rFlex>
           <img src=Images.fail2 className=Styles.resultIcon />
           <HSpacing size=Spacing.md />
           <Text value="Broadcast Transaction Failed" weight=Text.Semibold />
         </div>
         <VSpacing size=Spacing.md />
         <div className=Styles.errMsg>
           <Text value=err color=Colors.red3 align=Text.Center breakAll=true />
         </div>
       </div>
     }}
  </div>;
};
