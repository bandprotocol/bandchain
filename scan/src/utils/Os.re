let checkHID: unit => bool = [%bs.raw {|
function() {
  return !!navigator.hid;
}
|}];

let isWindows: unit => bool = [%bs.raw
  {|
function() {
  let x = navigator.userAgent;
  if (x) {
    return !!x.match(/NT/)
  } else {
    return false
  }
}
  |}
];
