let normalizeHexString = hexstr => {
  let len = hexstr->String.length;
  let prefix = len >= 2 ? hexstr->String.sub(0, 2) : "";

  switch (prefix) {
  | "0x"
  | "0X" => hexstr->String.lowercase_ascii->String.sub(2, len - 2)
  | _ => hexstr->String.lowercase_ascii
  };
};
