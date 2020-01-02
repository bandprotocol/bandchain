module Styles = {
  open Css;
};

[@react.component]
let make = () => {
  <Table header={"Latest Transactions" |> React.string} body={"123" |> React.string} />;
};
