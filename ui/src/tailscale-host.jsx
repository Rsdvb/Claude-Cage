import React, { useState } from "react";
import { ChevronDown, Copy, Network, Trash2 } from "lucide-react";
import { api } from "./api";

const nameOf = peer => peer.HostName || peer.DNSName || peer.TailscaleIPs?.[0] || peer.ID;

export function HostTailscaleServices({ onConfirm, onCopy }) {
  const [open, setOpen] = useState(false);
  const [peers, setPeers] = useState([]);
  const [forwards, setForwards] = useState([]);
  const [device, setDevice] = useState("");
  const [port, setPort] = useState("22");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const selected = peers.find(peer => peer.ID === device);
  const tcpPort = Number(port);
  const validPort = Number.isInteger(tcpPort) && tcpPort >= 1 && tcpPort <= 65535;

  async function request(action, body = {}) {
    setBusy(true);
    setError("");
    try {
      const result = await api(`/tailscale/host-${action}`, body);
      if (action === "peers") {
        const devices = result.peers || [];
        setPeers(devices);
        setDevice(current => devices.some(peer => peer.ID === current) ? current : "");
      }
      setForwards(result.forwards);
      setError(result.warning || "");
    } catch (e) {
      setError(e.message);
    } finally {
      setBusy(false);
    }
  }

  async function create() {
    if (!selected || !validPort) return;
    if (await onConfirm({
      title: "创建宿主机 Tailscale 连接",
      body: `允许当前客机中的应用通过宿主机访问 ${nameOf(selected)} 的 TCP 端口 ${tcpPort}。目标设备会识别为宿主机的 Tailscale 身份。`,
      confirm: "创建连接地址",
    })) await request("open-service", { device, port: tcpPort });
  }

  return <section className="tailscale-section">
    <button className="tailscale-disclosure" aria-expanded={open} aria-controls="tailnet-host-services"
      disabled={busy} onClick={() => { setOpen(!open); if (!open) void request("peers"); }}>
      <span><Network size={16} />通过宿主机 Tailscale 连接</span><ChevronDown size={16} />
    </button>
    {open && <div id="tailnet-host-services" className="tailscale-detail">
      <p className="field-help">在宿主机登录 Tailscale 后，选择一台设备和一个 TCP 端口。生成的地址供当前 Cage 客机使用。</p>
      <div className="tailscale-connect-form">
        <select aria-label="宿主机私网设备" value={device} disabled={busy} onChange={e => setDevice(e.target.value)}>
          <option value="">选择设备</option>
          {peers.map(peer => <option key={peer.ID} value={peer.ID}>{nameOf(peer)} · {peer.TailscaleIPs[0]}</option>)}
        </select>
        <input aria-label="宿主机私网服务端口" type="number" min="1" max="65535" value={port} onChange={e => setPort(e.target.value)} />
        <button className="quiet-button" disabled={busy || !selected || !validPort} onClick={create}>创建连接地址</button>
        <button className="quiet-button" disabled={busy} onClick={() => request("peers")}>刷新设备</button>
      </div>
      {error && <p className="error-message" role="alert">{error}</p>}
      {forwards.map(entry => <div className="tailscale-forward" key={`${entry.device}:${entry.port}`}>
        <span><strong>{nameOf(peers.find(peer => peer.ID === entry.device) || { HostName: entry.device })} · {entry.port}</strong><code>{entry.address}</code></span>
        <button className="icon-button" aria-label="复制宿主机连接地址" onClick={() => onCopy(entry.address)}><Copy size={14} /></button>
        <button className="icon-button" aria-label="关闭宿主机连接" disabled={busy} onClick={() => request("close-service", { device: entry.device, port: entry.port })}><Trash2 size={14} /></button>
      </div>)}
    </div>}
  </section>;
}
