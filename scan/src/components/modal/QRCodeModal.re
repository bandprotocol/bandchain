module Styles = {
  open Css;
  let container =
    style([
      width(`percent(100.)),
      maxWidth(`px(468)),
      minHeight(`px(360)),
      padding(`px(40)),
      Media.mobile([maxWidth(`px(300))]),
    ]);
};

[@react.component]
let make = (~address) => {
  <div className=Styles.container>
    <Heading size=Heading.H3 value="QR Code" marginBottom=8 align=Heading.Center />
    <AddressRender address position=AddressRender.Subtitle clickable=false />
    <div
      className={Css.merge([
        CssHelper.flexBox(~justify=`center, ()),
        CssHelper.mt(~size=40, ()),
      ])}>
      <QRCode value={address |> Address.toBech32} size=200 />
    </div>
  </div>;
};
