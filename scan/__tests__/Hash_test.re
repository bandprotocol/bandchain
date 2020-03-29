open Jest;
open Hash;
open Expect;

describe("Expect Hash to work correctly", () => {
  test("should be able to create Hash from hex", () =>
    expect("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff" |> fromHex)
    |> toEqual(Hash("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff"))
  );

  test("should be able to create Hash from hex with 0x", () =>
    expect("0x28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff" |> fromHex)
    |> toEqual(Hash("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff"))
  );

  test("should be able to create Hash from hex with 0X", () =>
    expect("0X28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff" |> fromHex)
    |> toEqual(Hash("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff"))
  );

  test({j|should be able to create Hash from hex with \\x|j}, () =>
    expect({j|\\x28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff|j} |> fromHex)
    |> toEqual(Hash("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff"))
  );

  test("should be able to get hexString with 0x prefix", () =>
    expect(
      fromHex("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff")
      |> toHex(~with0x=true),
    )
    |> toBe("0x28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff")
  );

  test("should be able to create Hash from base64", () =>
    expect("AupJR1NqB6pTA+sUrqv4HxcR3U6SjLp3gGCngA+KhB+H" |> fromBase64)
    |> toEqual(Hash("02ea4947536a07aa5303eb14aeabf81f1711dd4e928cba778060a7800f8a841f87"))
  );

  test("should be able get base64 string from Hash", () =>
    expect(Hash.fromBase64("AupJR1NqB6pTA+sUrqv4HxcR3U6SjLp3gGCngA+KhB+H") |> toBase64)
    |> toBe("AupJR1NqB6pTA+sUrqv4HxcR3U6SjLp3gGCngA+KhB+H")
  );

  test("should be able get hex string from Hash", () =>
    expect(Hash.fromBase64("AupJR1NqB6pTA+sUrqv4HxcR3U6SjLp3gGCngA+KhB+H") |> toHex)
    |> toBe("02ea4947536a07aa5303eb14aeabf81f1711dd4e928cba778060a7800f8a841f87")
  );

  test("should be able to get hexString in uppercase", () =>
    expect(
      fromHex("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff")
      |> toHex(~upper=true),
    )
    |> toBe("28913C89FA628136FFFCE7DED99D65A4E3F5C211F82639FED4ADCA30D53B8DFF")
  );

  test("should be able to get hexString with 0x prefix in uppercase", () =>
    expect(
      fromHex("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff")
      |> toHex(~with0x=true, ~upper=true),
    )
    |> toBe("0X28913C89FA628136FFFCE7DED99D65A4E3F5C211F82639FED4ADCA30D53B8DFF")
  );
});
