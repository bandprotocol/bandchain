type t =
  | DataSource(int)
  | OracleScript(int)
  | Request(int)
  | Block(int);

let getRoute =
  fun
  | DataSource(id) => Route.DataSourceIndexPage(id, Route.DataSourceExecute)
  | OracleScript(id) => Route.HomePage // TODO: change it later
  | Request(id) => Route.RequestIndexPage(id, Route.RequestReportStatus)
  | Block(height) => Route.BlockIndexPage(height);

let toString =
  fun
  | DataSource(id) => "#D" ++ (id |> string_of_int)
  | OracleScript(id) => "#O" ++ (id |> string_of_int)
  | Request(id) => "#R" ++ (id |> string_of_int)
  | Block(height) => "#B" ++ (height |> string_of_int);

let getColor =
  fun
  | DataSource(_) => Colors.brightOrange
  | OracleScript(_) => Colors.brightRed
  | Request(_) => Colors.orange
  | Block(_) => Colors.brightBlue;
