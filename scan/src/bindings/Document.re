[@bs.scope "document"] [@bs.val]
external addKeyboardEventListener: (string, ReactEvent.Keyboard.t => unit) => unit =
  "addEventListener";

[@bs.scope "document"] [@bs.val]
external removeKeyboardEventListener: (string, ReactEvent.Keyboard.t => unit) => unit =
  "removeEventListener";
