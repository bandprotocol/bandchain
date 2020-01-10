open Jest;
open HexUtils;
open Expect;

describe("Expect Hash to work correctly", () => {
  test("should be able to strip 0x out", () =>
    expect("0x"->normalizeHexString) |> toBe("")
  );

  test("should be able to strip 0x out and return the lower case of substring after 0x", () =>
    expect("0xaAA22E077492CbaD414098EBD98AA8dc1C7AE8D9"->normalizeHexString)
    |> toBe("aaa22e077492cbad414098ebd98aa8dc1c7ae8d9")
  );

  test("should be able to strip 0X out", () =>
    expect("0X"->normalizeHexString) |> toBe("")
  );

  test("should be able to strip 0X out and return the lower case of substring after 0X", () =>
    expect("0XaAA22E077492CbaD414098EBD98AA8dc1C7AE8D9"->normalizeHexString)
    |> toBe("aaa22e077492cbad414098ebd98aa8dc1c7ae8d9")
  );

  test("should return an empty string if the input is an empty string", () =>
    expect(""->normalizeHexString) |> toBe("")
  );

  test("should be able to to return lower case of the input", () =>
    expect("aAA22E077492CbaD414098EBD98AA8dc1C7AE8D9"->normalizeHexString)
    |> toBe("aaa22e077492cbad414098ebd98aa8dc1c7ae8d9")
  );
});
