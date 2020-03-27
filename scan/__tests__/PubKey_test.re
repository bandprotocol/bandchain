open Jest;
open PubKey;
open Expect;

describe("Expect PubKey to work correctly", () => {
  test("should be able to create PubKey from hex", () =>
    expect(
      "eb5ae9872103a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33" |> fromHex,
    )
    |> toEqual(
         PubKey("eb5ae9872103a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33"),
       )
  );

  test("should be able to create PubKey from hex with 0x prefix", () =>
    expect(
      "0xeb5ae9872103a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33" |> fromHex,
    )
    |> toEqual(
         PubKey("eb5ae9872103a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33"),
       )
  );

  test("should be able to get hexString with 0x prefix", () =>
    expect(
      fromHex("eb5ae9872103a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33")
      |> toHex(~with0x=true),
    )
    |> toBe("0xeb5ae9872103a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba33")
  );

  test("should be able to create PubKey from fromBech32", () =>
    expect(
      "bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj"
      |> fromBech32,
    )
    |> toEqual(
         PubKey("eb5ae9872103d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f"),
       )
  );

  test("should be able to convert self to hex", () =>
    expect("AjW2GKsPDp9Isa8y9Wt42VXEMieYk3FHN6k3A1AkuDxY" |> fromBase64)
    |> toEqual(PubKey("0235b618ab0f0e9f48b1af32f56b78d955c432279893714737a937035024b83c58"))
  );

  test("should be able to convert self to address", () =>
    expect(
      [|
        PubKey.fromBech32(
          "bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj",
        ),
        PubKey.fromBech32(
          "bandvalconspub1addwnpepq06h7wvh5n5pmrejr6t3pyn7ytpwd5c0kmv0wjdfujs847em8dusjl96sxg",
        ),
        PubKey.fromBech32(
          "bandvalconspub1addwnpepqwj5l74gfj8j77v8st0gh932s3uyu2yys7n50qf6pptjgwnqu2arxkkn82m",
        ),
        PubKey.fromBech32(
          "bandvalconspub1addwnpepqfey4c5ul6m5juz36z0dlk8gyg6jcnyrvxm4werkgkmcerx8fn5g2gj9q6w",
        ),
      |]
      |> Array.map(toAddress),
    )
    |> toEqual([|
         Address.Address("f0c23921727d869745c4f9703cf33996b1d2b715"),
         Address.Address("f23391b5dbf982e37fb7dadea64aae21cae4c172"),
         Address.Address("bdb6a0728c8dfe2124536f16f2ba428fe767a8f9"),
         Address.Address("5179b0bb203248e03d2a1342896133b5c58e1e44"),
       |])
  );

  test("should be able to convert self to hex", () =>
    expect(
      PubKey("eb5ae9872103d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f")
      |> toHex,
    )
    |> toEqual("eb5ae9872103d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f")
  );

  test("should be able to convert self to hex", () =>
    expect(
      PubKey("eb5ae9872103d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f")
      |> toPubKeyHexOnly,
    )
    |> toEqual("03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f")
  );

  test("should be able to convert self to hex with 0x", () =>
    expect(
      PubKey("eb5ae9872103d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f")
      |> toPubKeyHexOnly(~with0x=true),
    )
    |> toEqual("0x03d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f")
  );

  test("should be able to convert self to bech32", () =>
    expect(
      PubKey("eb5ae9872103d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f")
      |> toBech32,
    )
    |> toEqual(
         "bandvalconspub1addwnpepq0grwz83v8g4s06fusnq5s4jkzxnhgvx67qr5g7v8tx39ur5m8tk7rg2nxj",
       )
  );
});
