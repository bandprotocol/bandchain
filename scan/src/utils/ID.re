module type RawIDSig = {
  let prefix: string;
  let color: Css.Types.Color.t;
  let route: int => Route.t;
};

module RawDataSourceID = {
  let prefix = "#D";
  let color = Colors.yellow5;
  let route = id => Route.DataSourceIndexPage(id, Route.DataSourceExecute);
};

module RawOracleScriptID = {
  let prefix = "#O";
  let color = Colors.pink5;
  let route = (id: int) => Route.HomePage;
};

module RawRequestID = {
  let prefix = "#R";
  let color = Colors.orange5;
  let route = id => Route.RequestIndexPage(id, Route.RequestReportStatus);
};

module RawBlock = {
  let prefix = "#B";
  let color = Colors.bandBlue;
  let route = height => Route.BlockIndexPage(height);
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
    | ID(id) => RawID.route(id);

  let toString =
    fun
    | ID(id) => RawID.prefix ++ string_of_int(id);
};

module DataSource = IDCreator(RawDataSourceID);
module OracleScript = IDCreator(RawOracleScriptID);
module Request = IDCreator(RawRequestID);
module Block = IDCreator(RawBlock);
