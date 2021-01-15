module Styles = {
  open Css;
  let container =
    style([
      width(`percent(100.)),
      maxWidth(`px(468)),
      minHeight(`px(200)),
      padding(`px(40)),
      Media.mobile([maxWidth(`px(300))]),
    ]);
};

[@react.component]
let make = () => {
  <div className=Styles.container>
    <Heading
      size=Heading.H3
      value="Cosmoscan is syncing the database state."
      marginBottom=24
      align=Heading.Center
    />
    <Text
      value="The database has suddenly updated some changes for fixing some bugs. It means the old block and transaction information will not show on Cosmoscan for now. Please wait until the state is up to date."
      size=Text.Lg
    />
  </div>;
};
