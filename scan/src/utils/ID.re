module type RawIDSig = {
  type tab_t;
  let prefix: string;
  let color: Css.Types.Color.t;
  let route: (int, tab_t) => Route.t;
  let defaultTab: tab_t;
};

module RawDataSourceID = {
  type tab_t = Route.data_source_tab_t;
  let prefix = "#D";
  let color = Colors.yellow5;
  let route = (id, tab) => Route.DataSourceIndexPage(id, tab);
  let defaultTab = Route.DataSourceExecute;
};

module RawOracleScriptID = {
  type tab_t = Route.oracle_script_tab_t;
  let prefix = "#O";
  let color = Colors.pink5;
  let route = (id, tab) => Route.OracleScriptIndexPage(id, tab);
  let defaultTab = Route.OracleScriptExecute;
};

module RawRequestID = {
  type tab_t = Route.request_tab_t;
  let prefix = "#R";
  let color = Colors.orange5;
  let route = (id, tab) => Route.RequestIndexPage(id, tab);
  let defaultTab = Route.RequestReportStatus;
};

module RawBlock = {
  type tab_t = unit;
  let prefix = "#B";
  let color = Colors.bandBlue;
  let route = (height, _) => Route.BlockIndexPage(height);
  let defaultTab = ();
};

module type IDSig = {
  include RawIDSig;
  type t;
  let getRoute: t => Route.t;
  let toString: t => string;
};

module IDCreator = (RawID: RawIDSig) => {
  include RawID;

  type t =
    | ID(int);

  let getRoute =
    fun
    | ID(id) => RawID.route(id, RawID.defaultTab);

  let getRouteWithTab = (ID(id), tab) => RawID.route(id, tab);

  let toString =
    fun
    | ID(id) => RawID.prefix ++ string_of_int(id);

  let toInt =
    fun
    | ID(id) => id;

  let fromJson = json => ID(json |> Js.Json.decodeNumber |> Belt.Option.getExn |> int_of_float);

  let toJson =
    fun
    | ID(id) => id |> float_of_int |> Js.Json.number;
};

module DataSource = IDCreator(RawDataSourceID);
module OracleScript = IDCreator(RawOracleScriptID);
module Request = IDCreator(RawRequestID);
module Block = IDCreator(RawBlock);
