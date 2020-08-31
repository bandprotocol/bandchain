[@react.component]
let make = (~value) => {
  value->MarkedJS.marked->MarkedJS.parse;
};
