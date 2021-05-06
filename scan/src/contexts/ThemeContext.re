type t = {
  isDarkMode: bool,
  theme: Theme.t,
};

let keyword = "theme";
let context = React.createContext(ContextHelper.default);

let getThemeMode = () => {
  let localOpt = LocalStorage.getItem(keyword);
  {
    let%Opt local = localOpt;
    local == "dark" ? Some(Theme.Dark) : Some(Day);
  }
  ->Belt.Option.getWithDefault(Day);
};

let setThemeMode =
  fun
  | Theme.Day => LocalStorage.setItem(keyword, "day")
  | Dark => LocalStorage.setItem(keyword, "dark");

[@react.component]
let make = (~children) => {
  let (mode, setMode) = React.useState(_ => getThemeMode());

  let toggle = () =>
    setMode(prevMode => {
      switch (prevMode) {
      | Day =>
        setThemeMode(Dark);
        Dark;
      | Dark =>
        setThemeMode(Day);
        Day;
      }
    });

  let theme = React.useMemo1(() => Theme.get(mode), [|mode|]);
  let data = {isDarkMode: mode == Dark, theme};

  React.createElement(
    React.Context.provider(context),
    {"value": (data, toggle), "children": children},
  );
};
