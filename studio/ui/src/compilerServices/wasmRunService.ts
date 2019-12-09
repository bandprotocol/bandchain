/* Copyright 2018 Mozilla Foundation
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
import { CompilerService, ServiceInput, ServiceOutput } from "./types";

import { sendRequestJSON, ServiceTypes } from "./sendRequest";
import { base64EncodeBytes } from "../util";

export class WasmRunService implements CompilerService {
  async compile(input: ServiceInput): Promise<ServiceOutput> {
    let result;
    const options = input.options;

    const files = Object.values(input.files);
    if (files.length !== 1) {
      throw new Error(
        `Supporting execution of a single WASM file, but ${files.length} file(s) found`
      );
    }

    const code = base64EncodeBytes(
      new Uint8Array(files[0].content as Iterable<number>)
    );
    result = await sendRequestJSON({ code, options }, ServiceTypes.WasmRun);

    return {
      success: result.success,
      console: result.message,
      items: {}
    };
  }
}
