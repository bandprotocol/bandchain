let let_ = (a, b) =>
  switch (a) {
  | None => None
  | Some(x) => b(x)
  };

let ret = x => Some(x);
