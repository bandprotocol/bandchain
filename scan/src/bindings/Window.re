[@bs.val] external alert: string => unit = "alert";

[@bs.val] external prompt: (string, string) => Js.Nullable.t(string) = "prompt";
