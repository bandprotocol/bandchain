[@bs.deriving jsConverter]
type style_t('a, 'b) = {
  control: 'a => 'a,
  option: 'b => 'b,
};

[@bs.obj]
external makeProps:
  (~value: 'a, ~onChange: 'a => unit, ~options: array('a), ~styles: style_t('b, 'c), unit) => _ =
  "";

[@bs.module "react-select"]
external make:
  React.component({
    .
    "value": 'a,
    "onChange": 'a => unit,
    "options": array('a),
    "styles": style_t('b, 'c),
  }) =
  "default";
