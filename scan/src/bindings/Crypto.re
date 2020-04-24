type t;

[@bs.module "crypto"] [@bs.val]
external createHash: string => t = "createHash";

[@bs.send] external update: (t, string) => t = "update";

[@bs.send] external digest: (t, string) => string = "digest";
