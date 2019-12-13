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
import { Project, ModelRef } from "../models";

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
          " "
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
    step: number;
    done: boolean;
    requestId: number;
    txHash: string;
  }
> {
  constructor(props: any) {
    super(props);
    this.state = {
      deploying: false,
      step: -1,
      done: false,
      requestId: null,
      txHash: null
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
      const bot1URL = "http://134.209.106.94:5000";
      const explorerURL = "http://134.209.106.94:12000";
      const abciURL = "http://134.209.106.94:26657";

      this.setState({
        step: 0,
        done: true,
        deploying: true
      });
      await new Promise(r => setTimeout(r, 500));

      const code = this.getWASMFileContent();

      this.setState({
        step: 1,
        done: false,
        deploying: true
      });

      const requestResult = await axios.post(bot1URL + "/request", {
        code,
        delay: 5
      });

      const { txHash, reqID } = requestResult.data;

      this.setState({
        step: 1,
        done: true,
        txHash: txHash,
        requestId: reqID
      });

      await new Promise(r => setTimeout(r, 1000));

      this.setState({
        step: 2,
        done: false
      });

      await new Promise(r => setTimeout(r, 1000));

      this.setState({
        step: 3,
        done: false
      });

      const {
        data: {
          result: { reportEnd }
        }
      } = await axios.get(bot1URL + "/status?reqID=" + reqID);

      const getHeight = async () =>
        parseInt(
          (await axios.get(abciURL + "/status")).data.result.sync_info
            .latest_block_height
        );

      while (true) {
        const height = await getHeight();

        if (height >= reportEnd) break;

        this.setState({
          step: 4 - (reportEnd - height) / 5,
          done: false
        });

        await new Promise(r => setTimeout(r, 200));
      }

      this.setState({
        step: 4,
        done: false
      });

      while ((await getHeight()) < reportEnd + 1) {
        await new Promise(r => setTimeout(r, 200));
      }

      this.setState({
        step: 5,
        done: true
      });

      // Reset
      this.setState({
        deploying: false
      });
    } catch (e) {
      alert("Deployment failed :( Please contact dev@bandprotocol.com");
      console.error("Deployment failed", e);
      this.setState({
        deploying: false
      });
    }
  }

  render() {
    const steps = [
      { svg: "/svg/deployment/step1.svg", label: "Validating\nOWASM script" },
      {
        svg: "/svg/deployment/step2.svg",
        label: "Sending txn\nto D3N"
      },
      {
        svg: "/svg/deployment/step3.svg",
        label: "Transaction\nConfirmed",
        getLink: this.state.txHash
          ? () => ({
              label: "Explorer",
              href: `http://134.209.106.94:12000/transactions/${this.state.txHash}`
            })
          : undefined
      },
      { svg: "/svg/deployment/step4.svg", label: "Waiting for\ndata queries" },
      { svg: "/svg/deployment/step5.svg", label: "Executing\nOWASM script" },
      {
        svg: "/svg/deployment/step6.svg",
        label: "Data ready\nto use",
        getLink:
          this.state.step === 5
            ? () => ({
                label: "Data & Proof",
                href: `http://134.209.106.94:5000/proof?reqID=${this.state.requestId}`
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
                    this.state.step === 5
                      ? 100
                      : Math.max((this.state.step + 0.5) * (100 / 6), 0)
                  }%`
                }}
              ></div>
            </div>
          </div>
          <div style={{ borderTop: "solid 1px #303030" }}>
            <Button
              icon={<GoX />}
              label="Cancel"
              title="Cancel"
              onClick={() => {
                this.props.onCancel();
              }}
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
