[@react.component]
let make = (~value) => {
  value->MarkedJS.marked->HTMLParser.parse;
};
