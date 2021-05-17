module Styles = {
  open Css;

  let version =
    style([
      display(`flex),
      borderRadius(`px(10)),
      backgroundColor(Colors.blue1),
      padding2(~v=`pxFloat(5.), ~h=`px(10)),
      minWidth(`px(120)),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(Spacing.xs),
      marginTop(`px(1)),
      position(`relative),
      cursor(`pointer),
      zIndex(3),
      Media.mobile([padding2(~v=`pxFloat(5.), ~h=`px(10))]),
      Media.smallMobile([minWidth(`px(90))]),
    ]);

  let versionLoading =
    style([
      display(`flex),
      borderRadius(`px(10)),
      backgroundColor(Colors.blue1),
      overflow(`hidden),
      height(`px(16)),
      justifyContent(`center),
      alignItems(`center),
      marginLeft(Spacing.xs),
      marginTop(`px(1)),
    ]);

  let downIcon = show =>
    style([
      width(`px(6)),
      marginTop(`px(1)),
      transform(`rotate(`deg(show ? 180. : 0.))),
      Media.mobile([width(`px(8)), height(`px(6))]),
    ]);

  let dropdown = show =>
    style([
      display(`flex),
      borderRadius(`px(10)),
      flexDirection(`column),
      justifyContent(`center),
      position(`absolute),
      width(`percent(100.)),
      alignItems(`center),
      backgroundColor(Colors.blue1),
      transition(~duration=200, "all"),
      top(`px(25)),
      height(`auto),
      padding4(~top=`pxFloat(4.), ~bottom=`zero, ~left=`px(8), ~right=`px(8)),
      opacity(show ? 1. : 0.),
      pointerEvents(show ? `auto : `none),
      Media.mobile([top(`px(35))]),
    ]);

  let link = style([textDecoration(`none)]);
};

type chainID =
  | WenchangTestnet
  | WenchangMainnet
  | GuanYuDevnet
  | GuanYuTestnet
  | GuanYuPOA
  | GuanYuMainnet
  | LaoziTestnet
  | Unknown;

let parseChainID =
  fun
  | "band-wenchang-testnet3" => WenchangTestnet
  | "band-wenchang-mainnet" => WenchangMainnet
  | "band-guanyu-devnet5"
  | "band-guanyu-devnet6"
  | "band-guanyu-devnet7"
  | "band-guanyu-devnet8"
  | "bandchain" => GuanYuDevnet
  | "band-guanyu-testnet1"
  | "band-guanyu-testnet2"
  | "band-guanyu-testnet3"
  | "band-guanyu-testnet4" => GuanYuTestnet
  | "band-guanyu-poa" => GuanYuPOA
  | "band-guanyu-mainnet" => GuanYuMainnet
  | "band-laozi-internal"
  | "band-laozi-testnet1"
  | "band-laozi-testnet2"
  | "band-guanyu-laozi1" => LaoziTestnet
  | _ => Unknown;

let getLink =
  fun
  | WenchangTestnet => "https://wenchang-testnet3.cosmoscan.io/"
  | WenchangMainnet
  | GuanYuMainnet => "https://cosmoscan.io/"
  | GuanYuDevnet => "https://guanyu-devnet.cosmoscan.io/"
  | GuanYuTestnet => "https://guanyu-testnet3.cosmoscan.io/"
  | GuanYuPOA => "https://guanyu-poa.cosmoscan.io/"
  | LaoziTestnet => "https://laozi-testnet1.cosmoscan.io/"
  | Unknown => "";

let getName =
  fun
  | WenchangTestnet => "wenchang-testnet"
  | WenchangMainnet => "wenchang-mainnet"
  | GuanYuDevnet => "guanyu-devnet"
  | GuanYuTestnet => "guanyu-testnet"
  | GuanYuPOA => "guanyu-poa"
  | GuanYuMainnet => "guanyu-mainnet"
  | LaoziTestnet => "laozi-testnet"
  | Unknown => "unknown";

[@react.component]
let make = () =>
  {
    let (show, setShow) = React.useState(_ => false);
    let trackingSub = TrackingSub.use();
    let%Sub tracking = trackingSub;
    let currentChainID = tracking.chainID->parseChainID;

    <div
      className=Styles.version
      onClick={event => {
        setShow(oldVal => !oldVal);
        ReactEvent.Mouse.stopPropagation(event);
      }}>
      <Text
        value={currentChainID->getName}
        size=Text.Sm
        color=Colors.blue6
        nowrap=true
        weight=Text.Semibold
        spacing={Text.Em(0.03)}
      />
      <HSpacing size=Spacing.sm />
      <img src=Images.triangleDown className={Styles.downIcon(show)} />
      <div className={Styles.dropdown(show)}>
        {[|GuanYuMainnet, GuanYuTestnet|]
         ->Belt.Array.keep(chainID => chainID != currentChainID)
         ->Belt.Array.map(chainID => {
             let name = chainID->getName;
             <a
               href={getLink(chainID)}
               key=name
               className=Styles.link
               target="_blank"
               rel="noopener">
               <Text
                 value=name
                 size=Text.Sm
                 color=Colors.blue6
                 nowrap=true
                 weight=Text.Semibold
                 spacing={Text.Em(0.03)}
               />
               <VSpacing size={`px(8)} />
             </a>;
           })
         ->React.array}
      </div>
    </div>
    |> Sub.resolve;
  }
  |> Sub.default(
       _,
       {
         let width = Media.isSmallMobile() ? 80 : 110;
         <div className=Styles.versionLoading>
           <LoadingCensorBar width height=20 colorBase=Colors.blue1 colorLighter=Colors.white />
         </div>;
       },
     );
