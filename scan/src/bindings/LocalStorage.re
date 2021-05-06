[@bs.scope "localStorage"] [@bs.val] external getItem: string => option(string) = "getItem";
[@bs.scope "localStorage"] [@bs.val] external setItem: (string, string) => unit = "setItem";
