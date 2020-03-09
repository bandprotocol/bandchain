open Jest;
open Parser;

open Expect;

describe("Expect Parser to work correctly", () => {
  test("should be able to create Hash from hex", () =>
    expect("benz=$1" |> checker) |> toEqual(true)
  );

  test("should be able to create Hash from hex", () =>
    expect(
      {f|#!/bin/sh

symbol=$1

# Cryptocurrency price endpoint: https://www.coingecko.com/api/documentations/v3
url=\"https://api.coingecko.com/api/v3/simple/price?ids=$symbol&vs_currencies=usd\"

# Performs data fetching and parses the result
curl -s -X GET $url -H \"accept: application/json\" | jq -er ".[\"$symbol\"].usd\"
|f}
      |> getVariable,
    )
    |> toEqual(["symbol"])
  );

  test("should be able to create Hash from hex", () =>
    expect({f|"#!/bin/bash

calldata=$1

echo $calldata
"|f} |> getVariable)
    |> toEqual(["calldata"])
  );

  test("should be able to create Hash from hex", () =>
    expect({f|"#!/bin/bash

calldata=$1
symbol=$1

echo $calldata
"|f} |> getVariable)
    |> toEqual(["calldata", "symbol"])
  );

  test("should be able to create Hash from hex", () =>
    expect(
      "IyEvYmluL3NoCgpzeW1ib2w9JDEKCiMgQ3J5cHRvY3VycmVuY3kgcHJpY2UgZW5kcG9pbnQ6IGh0dHBzOi8vd3d3LmNvaW5nZWNrby5jb20vYXBpL2RvY3VtZW50YXRpb25zL3YzCnVybD0iaHR0cHM6Ly9hcGkuY29pbmdlY2tvLmNvbS9hcGkvdjMvc2ltcGxlL3ByaWNlP2lkcz0kc3ltYm9sJnZzX2N1cnJlbmNpZXM9dXNkIgoKIyBQZXJmb3JtcyBkYXRhIGZldGNoaW5nIGFuZCBwYXJzZXMgdGhlIHJlc3VsdApjdXJsIC1zIC1YIEdFVCAkdXJsIC1IICJhY2NlcHQ6IGFwcGxpY2F0aW9uL2pzb24iIHwganEgLWVyICIuW1wiJHN5bWJvbFwiXS51c2QiCg"
      |> JsBuffer.fromBase64
      |> parseExecutableScript,
    )
    |> toEqual(["symbol"])
  );
});
// IyEvYmluL3NoCgpzeW1ib2w9JDEKCiMgQ3J5cHRvY3VycmVuY3kgcHJpY2UgZW5kcG9pbnQ6IGh0dHBzOi8vd3d3LmNvaW5nZWNrby5jb20vYXBpL2RvY3VtZW50YXRpb25zL3YzCnVybD0iaHR0cHM6Ly9hcGkuY29pbmdlY2tvLmNvbS9hcGkvdjMvc2ltcGxlL3ByaWNlP2lkcz0kc3ltYm9sJnZzX2N1cnJlbmNpZXM9dXNkIgoKIyBQZXJmb3JtcyBkYXRhIGZldGNoaW5nIGFuZCBwYXJzZXMgdGhlIHJlc3VsdApjdXJsIC1zIC1YIEdFVCAkdXJsIC1IICJhY2NlcHQ6IGFwcGxpY2F0aW9uL2pzb24iIHwganEgLWVyICIuW1wiJHN5bWJvbFwiXS51c2QiCg
