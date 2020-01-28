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

export class DeploymentDialog extends React.Component<
  {
    isOpen: boolean;
    templatesName: string;
    onCancel: () => void;
    project: ModelRef<Project>;
  },
  {
    deploying: boolean;
    name: string;
    nameEditable: boolean;
    codeUrl: string;
    step: number;
    done: boolean;
    codeHash: string;
    explorerLink: string;
  }
> {
  constructor(props: any) {
    super(props);
    this.state = {
      deploying: false,
      name: "",
      nameEditable: true,
      codeUrl: "",
      step: -1,
      done: false,
      codeHash: null,
      explorerLink: null
    };
  }

  async componentDidMount() {
    // Load WASM file
    const wasmFile = this.props.project.getModel().getFile("out.wasm");
    if (!wasmFile) {
      this.props.onCancel();
      alert("Please build before deploying. File not found: out.wasm");
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

  async deploy() {
    try {
      const { bandsvUrl, explorerUrl, studiosvUrl } = await getConfig();

      this.setState({
        deploying: true,
        nameEditable: false,
        step: 0,
        done: true,
        codeHash: null,
        explorerLink: null
      });
      await new Promise(r => setTimeout(r, 500));

      const code = this.getWASMFileContent();

      function getAllFiles(
        currentFile: Directory
      ): Array<{ name: string; content: string | ArrayBuffer }> {
        if (currentFile.type !== "directory") {
          return [
            {
              name: currentFile.getPath(),
              content: currentFile.getData()
            }
          ];
        }

        return currentFile.children.flatMap(getAllFiles);
      }

      const project = this.props.project.getModel();

      const files = [
        {
          name: project.getFile("Cargo.toml").getPath(),
          content: project.getFile("Cargo.toml").getData()
        },
        ...getAllFiles(project.getFile("src") as Directory)
      ];

      await new Promise(r => setTimeout(r, 1000));

      const {
        data: { txHash, codeHash }
      } = await axios.post(bandsvUrl + "/store", {
        code,
        name: this.state.name
      });

      this.setState({
        step: 1,
        done: false,
        deploying: true,
        codeHash,
        explorerLink: explorerUrl + "/script/" + codeHash
      });

      const { data: codeUrl } = await axios.post(studiosvUrl + "/upload", {
        wasm: code,
        name: this.state.name,
        code: JSON.stringify(files)
      });

      this.setState({
        step: 2,
        done: true,
        codeUrl
      });

      await new Promise(r => setTimeout(r, 1000));

      this.setState({
        step: 2,
        done: false
      });

      await new Promise(r => setTimeout(r, 1000));

      this.setState({
        step: 3,
        done: true
      });
    } catch (e) {
      alert("Deployment failed :( Please contact dev@bandprotocol.com");
      console.error("Deployment failed", e);
    } finally {
      // Reset states
      this.setState({
        deploying: false,
        nameEditable: true
      });
    }
  }

  render() {
    const steps = [
      {
        svg: "/svg/deployment/step2.svg",
        label: "Create Data Request\nScript on D3N"
      },
      {
        svg: "/svg/deployment/step1.svg",
        label: "Uploading\nSource Code",
        getLink: this.state.codeUrl
          ? () => ({
              label: "Code",
              href: this.state.codeUrl
            })
          : undefined
      },
      {
        svg: "/svg/deployment/step6.svg",
        label: "Script ready\nfor query",
        getLink:
          this.state.step === 3
            ? () => ({
                label: "Block Explorer",
                href: `https://d3n-scan.onrender.com/script/${this.state.codeHash}`
              })
            : undefined
      }
    ];

    return (
      <ReactModal
        isOpen={this.props.isOpen}
        contentLabel="Deploy on D3N Test Network"
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
          <div className="modal-title-bar">Deploy on D3N Test Network</div>
          <div
            style={{
              position: "relative",
              flex: 1,
              display: "flex",
              background: "#1f1f1f",
              padding: "30px 15px"
            }}
          >
            {steps.map((step, idx) => (
              <Step
                key={idx}
                svg={step.svg}
                label={step.label}
                active={this.state.step >= idx}
                getLink={step.getLink}
                done={
                  !Math.floor(this.state.step - idx) ? this.state.done : true
                }
              />
            ))}
            <div
              style={{
                position: "absolute",
                background: "#484848",
                height: 10,
                left: 20,
                right: 20,
                top: 112,
                borderRadius: 5
              }}
            >
              <div
                style={{
                  position: "relative",
                  backgroundImage:
                    "linear-gradient(270deg, #FACC8F 0%, #E682B5 100%)",
                  height: 10,
                  borderRadius: 5,
                  transition: "width 0.5s ease-out",
                  width: `${
                    this.state.step === 3
                      ? 100
                      : Math.max((this.state.step + 0.5) * (100 / 3), 0)
                  }%`
                }}
              ></div>
            </div>
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
            <input
              style={{
                display: "inline-block",
                height: 40,
                lineHeight: 40,
                border: 0,
                backgroundColor: "var(--grey-900)",
                padding: "0 8px",
                width: 180,
                borderBottom: "2px solid var(--title)",
                color: "2px solid var(--title)",
                fontFamily: '"Roboto", sans-serif',
                fontSize: 14
              }}
              placeholder="Build name (optional)"
              value={this.state.name}
              disabled={!this.state.nameEditable}
              onChange={e => this.setState({ name: e.target.value })}
            />
            <Button
              icon={<GoFile />}
              label="Deploy"
              title="Deploy"
              isDisabled={this.state.deploying}
              onClick={() => {
                this.deploy();
              }}
            />
          </div>
        </div>
      </ReactModal>
    );
  }
}
