module SendMsg = {
  [@react.component]
  let make = (~toAddress, ~amount) => {
    <div className={Css.merge([CssHelper.flexBox(~wrap=`nowrap, ()), CssHelper.overflowHidden])}>
      <AmountRender coins=amount />
      <HSpacing size=Spacing.sm />
      <Text value={j| to |j} size=Text.Md nowrap=true block=true />
      <HSpacing size=Spacing.sm />
      <AddressRender address=toAddress />
    </div>;
  };
};

module ReceiveMsg = {
  [@react.component]
  let make = (~fromAddress, ~amount) => {
    <div className={Css.merge([CssHelper.flexBox(~wrap=`nowrap, ()), CssHelper.overflowHidden])}>
      <AmountRender coins=amount />
      <HSpacing size=Spacing.sm />
      <Text value={j| from |j} size=Text.Md nowrap=true block=true />
      <HSpacing size=Spacing.sm />
      <AddressRender address=fromAddress />
    </div>;
  };
};
