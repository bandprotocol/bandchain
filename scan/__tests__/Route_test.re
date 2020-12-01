open Jest;
open Route;
open Expect;

describe("Expect Search Functionality to work correctly", () => {
  test("test block route", () =>
    expect("123" |> search) |> toEqual(BlockIndexPage(123))
  );
  test("test transaction route", () =>
    expect("22638794cb5f306ef929b90c58b27d26cb35a77ca5c5c624cf2025a98528c323" |> search)
    |> toEqual(
         TxIndexPage(
           "22638794cb5f306ef929b90c58b27d26cb35a77ca5c5c624cf2025a98528c323" |> Hash.fromHex,
         ),
       )
  );
  test("test transaction route prefix is 0x", () =>
    expect("22638794cb5f306ef929b90c58b27d26cb35a77ca5c5c624cf2025a98528c323" |> search)
    |> toEqual(
         TxIndexPage(
           "0x22638794cb5f306ef929b90c58b27d26cb35a77ca5c5c624cf2025a98528c323" |> Hash.fromHex,
         ),
       )
  );
  test("test block prefix is B", () =>
    expect("B123" |> search) |> toEqual(BlockIndexPage(123))
  );
  test("test block prefix is b", () =>
    expect("b123" |> search) |> toEqual(BlockIndexPage(123))
  );
  test("test data soure route prefix is D", () =>
    expect("D123" |> search) |> toEqual(DataSourceIndexPage(123, DataSourceRequests))
  );
  test("test data soure route prefix is d", () =>
    expect("d123" |> search) |> toEqual(DataSourceIndexPage(123, DataSourceRequests))
  );
  test("test request route prefix is R", () =>
    expect("R123" |> search) |> toEqual(RequestIndexPage(123))
  );
  test("test request route prefix is r", () =>
    expect("r123" |> search) |> toEqual(RequestIndexPage(123))
  );
  test("test oracle script route prefix is O", () =>
    expect("O123" |> search) |> toEqual(OracleScriptIndexPage(123, OracleScriptRequests))
  );
  test("test oracle script route prefix is o", () =>
    expect("O123" |> search) |> toEqual(OracleScriptIndexPage(123, OracleScriptRequests))
  );
  test("test validator route", () =>
    expect("bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec" |> search)
    |> toEqual(
         ValidatorIndexPage(
           "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec" |> Address.fromBech32,
           Reports,
         ),
       )
  );
  test("test account route", () =>
    expect("band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun" |> search)
    |> toEqual(
         AccountIndexPage(
           "band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun" |> Address.fromBech32,
           AccountTransactions,
         ),
       )
  );
  test("test page not found", () =>
    expect("D123DD" |> search) |> toEqual(NotFound)
  );
});
