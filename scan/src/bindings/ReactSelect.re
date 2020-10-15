type react_select_option_t = {
  value: string,
  label: string,
};

[@bs.deriving jsConverter]
type style_t('a, 'b, 'c, 'd, 'e) = {
  control: 'a => 'a,
  option: 'b => 'b,
  container: 'c => 'c,
  singleValue: 'd => 'd,
  indicatorSeparator: 'e => 'e,
};

[@bs.obj]
external makeProps:
  (
    ~value: react_select_option_t,
    ~onChange: 'a => unit,
    ~options: array('a),
    ~styles: style_t('b, 'c, 'd, 'e, 'f),
    unit
  ) =>
  _ =
  "";

[@bs.module "react-select"]
external make:
  React.component({
    .
    "value": react_select_option_t,
    "onChange": 'a => unit,
    "options": array('a),
    "styles": style_t('b, 'c, 'd, 'e, 'f),
  }) =
  "default";
