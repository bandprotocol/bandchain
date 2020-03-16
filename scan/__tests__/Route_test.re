open Jest;
open Route;
open Expect;

describe("Expect Search Functionality to work correctly", () => {
  test("test block route", () =>
    expect("123" |> search |> toString) |> toEqual("/block/123")
  );
  test("test block route", () =>
    expect(
      "22638794cb5f306ef929b90c58b27d26cb35a77ca5c5c624cf2025a98528c323" |> search |> toString,
    )
    |> toEqual("/tx/22638794cb5f306ef929b90c58b27d26cb35a77ca5c5c624cf2025a98528c323")
  );
  test("test block prefix is B", () =>
    expect("B123" |> search |> toString) |> toEqual("/block/123")
  );
  test("test block prefix is b", () =>
    expect("b123" |> search |> toString) |> toEqual("/block/123")
  );
  test("test data soure route prefix is D", () =>
    expect("D123" |> search |> toString) |> toEqual("/data-source/123")
  );
  test("test data soure route prefix is d", () =>
    expect("d123" |> search |> toString) |> toEqual("/data-source/123")
  );
  test("test request route prefix is R", () =>
    expect("R123" |> search |> toString) |> toEqual("/request/123")
  );
  test("test request route prefix is r", () =>
    expect("r123" |> search |> toString) |> toEqual("/request/123")
  );
  test("test oracle script route prefix is O", () =>
    expect("O123" |> search |> toString) |> toEqual("/script/123")
  );
  test("test oracle script route prefix is o", () =>
    expect("O123" |> search |> toString) |> toEqual("/script/123")
  );
  test("test validator route", () =>
    expect("bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec" |> search |> toString)
    |> toEqual("/validator/bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec#delegators")
  );
  test("test oracle script route", () =>
    expect("band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun" |> search |> toString)
    |> toEqual("/account/band1p40yh3zkmhcv0ecqp3mcazy83sa57rgjp07dun")
  );
  test("test page not found", () =>
    expect("D123DD" |> search |> toString) |> toEqual("/")
  );
});
