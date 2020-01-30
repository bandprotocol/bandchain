/* Copyright 2019 Band Protocol
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

import * as React from "react";
import * as ReactModal from "react-modal";
import axios from "axios";
import { Button } from "./shared/Button";
import { GoGear, GoFile, GoX, Icon } from "./shared/Icons";
import { Project, ModelRef, Directory } from "../models";
import getConfig from "../config";

const Step: React.SFC<{
  svg: string;
  label: string;
  active: boolean;
  done: boolean;
  getLink?: () => { label: string; href: string };
}> = ({ svg, label, getLink, active, done }) => {
  return (
    <div
      className={done ? "" : "blink"}
      style={{
        flex: 1,
        display: "flex",
        flexDirection: "column",
        alignItems: "center",
        height: 150,
        transition: "opacity 0.5s",
        opacity: active ? 1 : 0.5
      }}
    >
      <img style={{ height: 36 }} src={svg} />
      <div
        style={{
          flex: 1,
          display: "flex",
          flexDirection: "column",
          alignItems: "center",
          justifyContent: "flex-end",
          height: 32
        }}
      >
        {getLink ? (
          <a
            style={{
              padding: "2px 8px",
              borderRadius: 2,
              background: "#4D4D4D",
              cursor: "pointer",
              fontSize: 10,
              color: "#ffffff"
            }}
            target="_blank"
            href={getLink().href}
          >
            {getLink().label}
          </a>
        ) : (
          ""
        )}
      </div>
      <div
        style={{
          textAlign: "center",
          whiteSpace: "pre",
          fontSize: 12,
          fontWeight: "bold",
          lineHeight: "18px",
          marginTop: 48
        }}
      >
        {label}
      </div>
    </div>
  );
};

export class RunDialog extends React.Component<
  {
    isOpen: boolean;
    templatesName: string;
    onCancel: () => void;
    project: ModelRef<Project>;
  },
  {
    loadingTemplate: boolean;
    running: boolean;
    paramsTemplate: Array<Array<string>>;
    params: { [k: string]: string };
    result: string;
  }
> {
  constructor(props: any) {
    super(props);
    this.state = {
      loadingTemplate: false,
      running: false,
      paramsTemplate: null,
      params: null,
      result: null
    };
  }

  async componentDidMount() {
    // Load WASM file
    const wasmFile = this.props.project.getModel().getFile("out.wasm");
    if (!wasmFile) {
      this.props.onCancel();
      alert("Please build before running. File not found: out.wasm");
    }

    this.setState({
      loadingTemplate: true
    });

    // Load template
    try {
      const { bandsvUrl } = await getConfig();
      const code = this.getWASMFileContent();

      const {
        data: { params }
      } = await axios.post(bandsvUrl + "/params-info", {
        code
      });

      this.setState({
        loadingTemplate: false,
        paramsTemplate: params,
        params: {}
      });
    } catch (e) {
      alert(
        "Failed to load parameter info :( Please contact dev@bandprotocol.com"
      );
      console.log(e);
      this.props.onCancel();
    } finally {
    }
  }

  getWASMFileContent() {
    const wasmFile = this.props.project.getModel().getFile("out.wasm");
    const wasmFileContent = Array.from(
      new Uint8Array(wasmFile.getData() as Iterable<number>)
    )
      .map(p => p.toString(16).padStart(2, "0"))
      .join("");
    return wasmFileContent;
  }

  setParam(k: string, v: string) {
    this.setState({
      params: {
        ...this.state.params,
        [k]: v
      }
    });
  }

  async run() {
    // Make sure all the parameters are present and formatted
    const formattedParams: { [key: string]: string | number } = {};
    this.state.paramsTemplate.forEach(([name, type]) => {
      if (!this.state.params[name])
        return alert(`Parameter ${name} must not be empty`);

      if (["coins::Coins", "String"].includes(type))
        formattedParams[name] = this.state.params[name];

      if ("Int" === type) {
        if (isNaN(parseInt(this.state.params[name])))
          return alert(`Parameter ${name} must be integer`);
        formattedParams[name] = parseInt(this.state.params[name]);
      }
      if ("Float" === type) {
        if (isNaN(parseFloat(this.state.params[name])))
          return alert(`Parameter ${name} must be float`);
        formattedParams[name] = parseFloat(this.state.params[name]);
      }
    });

    this.setState({ running: true });

    try {
      const { bandsvUrl } = await getConfig();
      const code = this.getWASMFileContent();

      const {
        data: { result }
      } = await axios.post(bandsvUrl + "/execute", {
        code,
        params: formattedParams
      });

      this.setState({ running: false, result });
    } catch (e) {
      alert("Failed run OWASM script :( Please contact dev@bandprotocol.com");
      this.setState({ running: false, result: "" });
    }
  }

  render() {
    return (
      <ReactModal
        isOpen={this.props.isOpen}
        contentLabel="Run OWASM Script"
        className="modal show-file-icons"
        overlayClassName="overlay"
        ariaHideApp={false}
      >
        <div
          style={{
            display: "flex",
            flexDirection: "column",
            height: "100%"
          }}
        >
          <div className="modal-title-bar">Run OWASM Script</div>
          <div
            style={{
              position: "relative",
              flex: 1,
              background: "#1f1f1f",
              padding: "30px 15px"
            }}
          >
            {this.state.loadingTemplate && (
              <div
                style={{
                  padding: 30,
                  textAlign: "center"
                }}
              >
                Loading parameter info ...
              </div>
            )}
            {!this.state.loadingTemplate &&
              this.state.paramsTemplate &&
              !this.state.paramsTemplate.length && (
                <div
                  style={{
                    padding: 30,
                    textAlign: "center"
                  }}
                >
                  No parameter requires
                </div>
              )}
            {this.state.paramsTemplate &&
              this.state.paramsTemplate.map(([name, type], idx) => (
                <div style={{ display: "flex", marginBottom: 10 }}>
                  <div
                    style={{
                      lineHeight: "40px",
                      color: "#ffffff",
                      width: 150,
                      paddingRight: 30,
                      textAlign: "right"
                    }}
                  >
                    {name}
                  </div>
                  {type === "coins::Coins" ? (
                    <select
                      style={{
                        display: "inline-block",
                        height: 40,
                        lineHeight: "40px",
                        border: 0,
                        backgroundColor: "var(--grey-800)",
                        padding: "0 8px",
                        width: 196,
                        borderBottom: "2px solid var(--title)",
                        color: "2px solid var(--title)",
                        fontFamily: '"Roboto", sans-serif',
                        fontSize: 14
                      }}
                      placeholder="Build name (optional)"
                      value={this.state.params[name] || ""}
                      onChange={e => this.setParam(name, e.target.value)}
                      required
                    >
                      <option value="" disabled>
                        Select token
                      </option>
                      <option value="ADA">ADA</option>
                      <option value="BAND">BAND</option>
                      <option value="BCH">BCH</option>
                      <option value="BNB">BNB</option>
                      <option value="BTC">BTC</option>
                      <option value="EOS">EOS</option>
                      <option value="ETH">ETH</option>
                      <option value="LTC">LTC</option>
                      <option value="ETC">ETC</option>
                      <option value="TRX">TRX</option>
                      <option value="XRP">XRP</option>
                    </select>
                  ) : (
                    <input
                      style={{
                        display: "inline-block",
                        height: 40,
                        lineHeight: 40,
                        border: 0,
                        backgroundColor: "var(--grey-800)",
                        padding: "0 8px",
                        width: 180,
                        borderBottom: "2px solid var(--title)",
                        color: "2px solid var(--title)",
                        fontFamily: '"Roboto", sans-serif',
                        fontSize: 14
                      }}
                      placeholder={type}
                      value={this.state.params[name] || ""}
                      onChange={e => this.setParam(name, e.target.value)}
                      required
                    />
                  )}
                </div>
              ))}

            {this.state.result && (
              <>
                <div
                  style={{
                    borderBottom: "solid 1px var(--grey-800)",
                    margin: 20
                  }}
                ></div>
                {Object.entries(this.state.result).map(([key, val]) => (
                  <div style={{ display: "flex", marginBottom: 10 }}>
                    <div
                      style={{
                        lineHeight: "40px",
                        color: "#ffffff",
                        width: 150,
                        paddingRight: 30,
                        textAlign: "right"
                      }}
                    >
                      {key}
                    </div>
                    <div
                      style={{
                        lineHeight: "40px",
                        color: "#ff93c9",
                        paddingLeft: 8
                      }}
                    >
                      {val}
                    </div>
                  </div>
                ))}
              </>
            )}
          </div>
          <div style={{ borderTop: "solid 1px #303030" }}>
            <span style={{ float: "right" }}>
              <Button
                icon={<GoX />}
                label="Cancel"
                title="Cancel"
                onClick={() => {
                  this.props.onCancel();
                }}
              />
            </span>
            <Button
              icon={<GoFile />}
              label="Run script"
              title="Run script"
              isDisabled={this.state.running}
              onClick={() => {
                this.run();
              }}
            />
            <span
              style={{
                height: 40,
                padding: "0 20px",
                fontSize: "14px",
                color: "#ffffff"
              }}
            >
              {this.state.running && "Processing ..."}
            </span>
          </div>
        </div>
      </ReactModal>
    );
  }
}
