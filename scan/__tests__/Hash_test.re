open Jest;
open Hash;

describe("Expect Hash to work correctly", () => {
  Expect.(
    test("should be able create Hash from hex", () =>
      expect("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff"->fromHex)
      |> toEqual(Hash("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff"))
    )
  );

  Expect.(
    test("should be able create Hash from hex with 0x", () =>
      expect("0x28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff"->fromHex)
      |> toEqual(Hash("28913c89fa628136fffce7ded99d65a4e3f5c211f82639fed4adca30d53b8dff"))
    )
  );

  Expect.(
    test("should be able create Hash from base64", () =>
      expect("AupJR1NqB6pTA+sUrqv4HxcR3U6SjLp3gGCngA+KhB+H"->fromBase64)
      |> toEqual(Hash("02ea4947536a07aa5303eb14aeabf81f1711dd4e928cba778060a7800f8a841f87"))
    )
  );

  Expect.(
    test("should be able get base64 string from Hash", () =>
      expect(Hash.fromBase64("AupJR1NqB6pTA+sUrqv4HxcR3U6SjLp3gGCngA+KhB+H")->toBase64)
      |> toBe("AupJR1NqB6pTA+sUrqv4HxcR3U6SjLp3gGCngA+KhB+H")
    )
  );

  Expect.(
    test("should be able get hex string from Hash", () =>
      expect(Hash.fromBase64("AupJR1NqB6pTA+sUrqv4HxcR3U6SjLp3gGCngA+KhB+H")->toHex)
      |> toBe("02ea4947536a07aa5303eb14aeabf81f1711dd4e928cba778060a7800f8a841f87")
    )
  );
});
