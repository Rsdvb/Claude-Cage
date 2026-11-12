import React, { useState } from "react";
import {
  ShieldCheck,
  ShieldAlert,
  Check,
  CircleHelp,
  X,
  LoaderCircle,
  Copy,
  Settings,
  ChevronDown,
} from "lucide-react";
import { api, writeClipboard } from "./api";
import { Modal } from "./views";
import { networkState, reportForSharing } from "./safety-state.mjs";

export async function bindDetectedMacOSExit(ip, onConfirm) {
  if (!(await onConfirm({ title: "绑定固定出口？", body: `客机终端、浏览器和代理均测得 ${ip}。保存后请重启客机并重新检测。`, confirm: "绑定固定出口" }))) return null;
  return api("/macos/exit-policy", { ip });
}

export function SafetyPanel({
  status,
  onClose,
  onSettings,
  onConfirm,
}) {
  const [error, setError] = useState(""),
    [notice, setNotice] = useState("");
  const network = networkState(status),
    safety = status.safety || {},
    report = safety.report;
  const busy = Boolean(status.job?.running),
    online = network.mode === "open";
  const passed =
    report?.checks?.filter((c) => c.critical && c.state === "passed").length ||
    0;
  const total = report?.checks?.filter((c) => c.critical).length || 0;
  const run = async () => {
    setError("");
    if (
      status.vm_running &&
      !safety.gate &&
      !(await onConfirm({
        title: "重启环境并加载检测组件？",
        body: "当前终端会话将结束。请先保存文件，浏览器登录与项目文件会保留。",
        confirm: "重启并检查",
      }))
    )
      return;
    try {
      await api(status.ready && safety.gate ? "/safety/check" : "/start", {});
    } catch (e) {
      setError(e.message);
    }
  };
  const pause = async () => {
    setError("");
    try {
      await api("/safety/close", {});
    } catch (e) {
      setError(e.message);
    }
  };
  const copy = async () => {
    try {
      await writeClipboard(JSON.stringify(reportForSharing(report), null, 2));
      setNotice("已复制报告，不含代理凭据、IP 地址和本机路径");
    } catch {
      setError("无法写入本机剪贴板，请检查浏览器的剪贴板权限。");
    }
  };
  return (
    <Modal
      title="网络与隔离检查"
      wide
      onClose={onClose}
      footer={
        <>
          <span className="save-note">
            {busy ? status.job.message : "外网关闭时，离线终端和文件仍可使用"}
          </span>
          {(online || network.mode === "checking") && (
            <button onClick={pause}>
              {network.mode === "checking" ? "取消检测并关闭外网" : "暂停外网"}
            </button>
          )}
          <button className="primary" disabled={busy} onClick={run}>
            {busy ? (
              <LoaderCircle size={15} className="spin" />
            ) : (
              <ShieldCheck size={15} />
            )}
            {busy
              ? "检查进行中"
              : online
                ? "暂停外网并重检"
                : status.ready && safety.gate
                  ? "重新检测"
                  : status.vm_running
                    ? "重启并检查"
                    : "启动并检查"}
          </button>
        </>
      }
    >
      <div className="modal-body safety-body">
        {safety.macos_exit_candidate?.ip && <button disabled={busy} onClick={async () => {
          try { if (await bindDetectedMacOSExit(safety.macos_exit_candidate.ip, onConfirm)) setNotice("固定出口已提交，请在保存完成后重启客机并重新检测。"); }
          catch (e) { setError(e.message); }
        }}>绑定检测到的固定出口 {safety.macos_exit_candidate.ip}</button>}
        <div className={`safety-summary ${network.mode}`}>
          {online ? (
            <ShieldCheck size={27} />
          ) : network.mode === "checking" ? (
            <LoaderCircle size={27} className="spin" />
          ) : (
            <ShieldAlert size={27} />
          )}
          <div>
            <h3>{network.label}</h3>
            <p>
              {busy
                ? status.job.message
                : online
                  ? status.network_error || "关键检查已通过，后台继续核对出口规则与环境配置。"
                  : status.network_error || safety.error ||
                    safety.gate?.reason ||
                    "完成关键检查后，才允许连接外部网站。"}
            </p>
          </div>
          <button className="text-button" onClick={() => onSettings("network")}>
            <Settings size={15} />
            出口规则
          </button>
        </div>
        {(error || notice) && (
          <div
            role="status"
            className={error ? "error-message" : "inline-info"}
          >
            {error || notice}
          </div>
        )}
        <div className="safety-section-heading">
          <h4>关键检查</h4>
          <span>{total ? `${passed} / ${total} 通过` : "尚无检测记录"}</span>
          {report && (
            <button className="text-button" onClick={copy}>
              <Copy size={14} />
              复制报告
            </button>
          )}
        </div>
        {report ? (
          <>
            <p className="safety-timestamp">
              检测于 {new Date(report.created_at).toLocaleString("zh-CN")} ·{" "}
              {online ? "当前运行实例" : "报告不代表当前已放行"}
            </p>
            <div className="safety-checks">
              {report.checks
                .filter((check) => check.critical)
                .map((check) => {
                  const Icon =
                    check.state === "passed"
                      ? Check
                      : check.state === "failed"
                        ? X
                        : CircleHelp;
                  return (
                    <details
                      className={`safety-check ${check.state}`}
                      key={check.id}
                    >
                      <summary>
                        <Icon size={16} />
                        <strong>{check.title}</strong>
                        <span>
                          {check.state === "passed"
                            ? "通过"
                            : check.state === "failed"
                              ? "未通过"
                              : busy && check.checked_at === report.created_at
                                ? "待检查"
                                : "未知"}
                        </span>
                        <ChevronDown size={14} />
                      </summary>
                      <div>
                        <p>{check.evidence}</p>
                        {check.state !== "passed" && check.remedy && (
                          <p className="check-remedy">{check.remedy}</p>
                        )}
                      </div>
                    </details>
                  );
                })}
            </div>
          </>
        ) : (
          <p className="safety-empty">
            检查涵盖出站规则、代理故障关闭、DNS、共享目录、宿主管理接口，以及系统和浏览器的实际配置。关键项失败或未知时保持外网关闭。
          </p>
        )}
        {report?.checks?.some((check) => !check.critical) && (
          <div className="safety-advisories">
            {report.checks
              .filter((check) => !check.critical)
              .map((check) => (
                <p key={check.id}>
                  <CircleHelp size={15} />
                  <span>
                    <strong>{check.title}：</strong>
                    {check.evidence}
                  </span>
                </p>
              ))}
          </div>
        )}
        <details className="safety-scope">
          <summary>检查范围与来源</summary>
          <p>
            当前版本验证 IPv4 出口，客机公网 IPv6 与 UDP 关闭。出口监测约每 30
            秒执行；位置来自 IP
            数据库，城市判断具有误差。代理自身按目标分流的行为仍需在代理侧约束。
          </p>
          <p>
            报告包含运行时策略和采样证据，不能代替完整虚拟化审计。全面检查有效期为
            24 小时；控制服务失联后，联网授权最多保留 90 秒。
          </p>
        </details>
      </div>
    </Modal>
  );
}
