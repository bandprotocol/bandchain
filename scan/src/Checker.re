module SyncStaus = {
  [@react.component]
  let make = () => {
    let (syncing, setSyncing) = React.useState(_ => false);
    let (_, dispatchModal) = React.useContext(ModalContext.context);
    let trackingSub = TrackingSub.use();

    // If database is syncing the state (when replayOffset = -2).
    React.useEffect2(
      () => {
        switch (trackingSub) {
        | Data({replayOffset}) when replayOffset != (-2) && !syncing =>
          Syncing->OpenModal->dispatchModal;
          setSyncing(_ => true);
        | _ => ()
        };
        None;
      },
      (trackingSub, syncing),
    );

    React.null;
  };
};

let checkNetwork = () => {
  exception WrongNetwork(string);
  switch (Env.network) {
  | "WENCHANG"
  | "GUANYU38"
  | "GUANYU" => ()
  | _ => raise(WrongNetwork("Incorrect or unspecified NETWORK environment variable"))
  };
};

[@react.component]
let make = () => {
  checkNetwork();

  <SyncStaus />;
};
