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
import { decodeBinary } from "./utils";
import * as Tar from "tar-js";
import { base64EncodeBytes } from "../util";

export class RustTestService implements CompilerService {
  async compile(input: ServiceInput): Promise<ServiceOutput> {
    let result;
    const options = input.options;
    const cargo = options["cargo"];

    if (cargo) {
      const tarBuffer = new Tar();

      const files = input.files;
      Object.entries(files).forEach(([name, file]) => {
        tarBuffer.append(name, file.content, {});
      });

      const tar = base64EncodeBytes(tarBuffer.out);

      result = await sendRequestJSON({ tar, options }, ServiceTypes.CargoTest);

      return {
        success: result.success,
        console: result.message,
        items: {}
      };
    } else {
      throw new Error(`Supporting cargo test only`);
    }
  }
}
