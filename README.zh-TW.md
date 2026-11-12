<div align="center">

<h1>Claude Cage</h1>

**隔離開發工作台**

<img src="assets/claude-cage-character.png" alt="Claude Cage" width="320" />

[![Release](https://img.shields.io/github/v/release/Lee-xeo/Claude-Cage?include_prereleases)](https://github.com/Lee-xeo/Claude-Cage/releases) [![CI](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml) [![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

<a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/daily?language=Go" alt="Trendshift — Go Repository of the Day" width="250" height="55" /></a> <a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/weekly?language=Go" alt="Trendshift — Go Repository of the Week" width="250" height="55" /></a>

**[專案主頁](https://claude-cage.com/) · [下載應用程式](https://claude-cage.com/zh-TW/#download) · [開始使用](docs/getting-started.md) · [常見問題](docs/troubleshooting.md) · [使用文件（簡體中文）](https://claude-cage.com/docs/)**

[简体中文](README.md) · [繁體中文](README.zh-TW.md) · [Русский](README.ru.md) · [العربية](README.ar.md) · [فارسی](README.fa.md) · [한국어](README.ko.md) · [Español](README.es.md) · [မြန်မာဘာသာ](README.my.md) · [Беларуская](README.be.md) · [پښتو](README.ps.md) · [دری](README.prs.md) · [አማርኛ](README.am.md) · [Soomaali](README.so.md) · [ትግርኛ](README.ti.md) · [Français](README.fr.md)

</div>

## 隔離開發工作台

Claude Cage 讓你在熟悉的系統裡與 Claude 一起工作：Mac 上使用 macOS 虛擬機器，Windows 上使用 Windows 虛擬機器，也可以選擇輕量的 Linux 環境。終端、瀏覽器、編輯器和官方 Claude Desktop 集中在同一個工作臺，專案按需共享，網路獨立管理；還能連接自己的伺服器，或用手機接續電腦上的工作。

**Cage 執行主介面：**

![Cage](www/public/workbench-2026-11.jpg)

## 你可以用它做什麼

- **在一個視窗裡完成工作。** 用終端執行 Claude Code，用編輯器檢視改動，用 Cage 瀏覽器預覽服務和完成登入。
- **用熟悉的系統開發。** 選擇 macOS、Windows 或 Linux 客機，使用對應平台的工具；也可自動匹配本機系統版本，建立乾淨的獨立環境。 [選擇客機系統](docs/guide/guest-systems.md)
- **控制共享範圍。** 只共享預設工作區和明確開啟的專案資料夾；主機的 Claude 登入、瀏覽器配置和 SSH 金鑰不會自動匯入。
- **看清網路狀態。** 內建代理管理支援匯入訂閱、選擇節點和檢查出口。關鍵檢查未透過時，環境保持離線；代理不可用時沒有直連回退。
- **下次繼續使用。** 工具、檔案和登入狀態儲存在虛擬磁碟裡，關閉視窗可以保留後臺任務。
- **管理多個 Claude Code 帳號。** 查看額度、為新會話分配帳號，並在明確授權後切換帳號接續歷史。 [Claude Code](docs/guide/accounts.md)
- **用手機接續工作。** Android / iPhone 透過 QR Code 或配對碼連接電腦，接續同一份 Claude Code 對話、傳送訊息和處理工具審批，也可操作遠端瀏覽器和電腦上的 Claude Desktop。
- **在雲端執行環境。** 透過 SSH 連接部署在符合條件的 Linux x86_64 VPS 上的 Cage，無需 KVM；資料與網路出口位於伺服器。
- **透過私網連接。** 在 Cage 內登入 Tailscale，讓本機與環境透過 SSH 雙向連接。 [Tailscale / SSH](docs/guide/tailscale.md)

## 快速開始

把這句話發給你的 agent（codex，dsh，workbuddy，反正不會是 claude code😁）

```markdown
幫我安裝 https://github.com/Lee-xeo/Claude-Cage ，並把 Cage 的代理配置為我本機 Clash 的 xxx 節點。
```

## 下載與安裝

Cage 支援 macOS、Windows 和 Linux：macOS 14+ 的 Apple Silicon Mac、Windows 10 22H2 / 11 x64，以及 Debian 12 / Ubuntu 22.04 / 24.04 x86_64。 Intel Mac、Windows ARM64 和 Linux ARM64 尚不受支援。 [macOS](docs/getting-started.md) · [Windows](docs/windows.md) · [Linux](docs/linux.md)。

1. 開啟[專案主頁](https://claude-cage.com/zh-TW/#download)，點選「下載應用程式」並選擇平台。
2. macOS：開啟安裝包，將 **Cage** 拖入 **Applications（應用程式）**，再從“應用程式”啟動。Windows：[安裝說明](docs/windows.md#安装并开始工作) · Linux：[安裝說明](docs/linux.md#安装并开始工作)。
3. 按應用內引導配置自己的代理，啟動環境後開啟專案。

安裝包隨附工作臺和代理核心，無需先安裝 Go、Node.js、Docker 或 Clash。首次建立環境時選擇客機系統與儲存位置，再準備對應的官方系統映像和開發工具。Cage 不提供作業系統授權、代理服務或 Claude 帳戶。

原始碼與安裝包公開，版本說明和校驗檔案見 [Releases](https://github.com/Lee-xeo/Claude-Cage/releases)。

功能以所用建置為準，公開安裝包的能力與驗證範圍以發行說明為準。

macOS 和已簽章 Windows 版本支援透過所選代理更新；Windows 未簽章預覽版與 Linux 使用安裝包手動更新。更新前儲存並停止工作，個人資料繼續保留。 [升級與解除安裝](docs/guide/maintenance.md)

## 三步開始

### 1. 選擇代理

點選首頁“配置網路與代理”，再點選“新增來源”，匯入 Clash/Mihomo YAML 訂閱或本地 YAML 檔案，在“節點”中選擇一個節點。沒有可用代理時，可以明確勾選“透過主機網路下載這份訂閱”，或直接匯入本地檔案。後續可從“設定 → 網路與代理”管理節點。

### 2. 啟動環境

首次選擇客機系統或「自動匹配本機系統版本」，再選擇儲存位置。點擊「啟動 Cage 環境」，按引導完成安裝與授權；進度和錯誤可在「操作詳情」查看。以後啟動沿用原環境。

### 3. 開啟專案，執行 Claude

點選“開啟專案”選擇 本機 上的資料夾，再新建終端，輸入 `claude`，在 Cage 瀏覽器中完成自己的賬戶授權。也可以先使用預設工作區，無需提前準備專案。完整步驟見[開始使用](docs/getting-started.md)。

## 使用前瞭解這些邊界

- 共享專案中的修改和刪除會直接影響 本機 上的檔案。請使用版本控制或備份；專案內的 `.git` 可讀寫，可直接在 Cage 中暫存、提交與切換分支。
- 所有專案共用一個 Cage 環境，也共用其中的工具、賬戶與 macOS / Windows / Debian / Ubuntu 22.04 檔案系統。專案之間沒有獨立隔離。
- 預設公網出站支援 IPv4 TCP，DNS 透過同一代理；公網 UDP、IPv6、ICMP 及主機和區域網路直連受限。主動啟用 Tailscale 後，可依私網權限連接裝置。
- 網路檢查展示當前配置和取樣證據，無法保證賬戶不會受到平臺限制，也不能代替完整的虛擬化安全審計。

詳細機制見[網路與資料邊界](docs/getting-started.md#网络与数据边界)。

## 文件與貢獻

完整入口見[文件導航](docs/README.md)。

| 你要做什麼 | 從這裡開始 |
|---|---|
| macOS / Windows / Debian / Ubuntu 22.04 | [選擇客機系統](https://claude-cage.com/docs/guest-systems/) |
| 第一次安裝、配置代理、開啟專案 | [開始使用](docs/getting-started.md) |
| 處理下載、啟動、聯網和登入問題 | [常見問題與恢復](docs/troubleshooting.md) |
| 瞭解按鍵、外觀、資料儲存和退出行為 | [日常使用](docs/getting-started.md#日常使用) |
| 報告可復現問題 | [GitHub Issues](https://github.com/Lee-xeo/Claude-Cage/issues) |

本專案仍在快速迭代中，歡迎加入交流群提意見，歡迎社群大佬一起維護

<p align="center">
  <a href="assets/community-qq.png"><img src="assets/community-qq.png" alt="QQ 二群 · 1127457578" width="240" /></a>
  <br />QQ 二群 · 1127457578
</p>

## 來源與許可

Cage 基於 [agent-vm](https://github.com/sylvinus/agent-vm) 二次開發。

專案主許可為 [MIT](LICENSE)。隨附元件保留各自許可，Mihomo 作為獨立可執行檔案按 GPL-3.0 分發；請同時參閱 [NOTICE](NOTICE) · [第三方許可](THIRD_PARTY_LICENSES) · [noVNC](ui/vendor/novnc/LICENSE.txt)。

Cage 是獨立社群專案，與 Anthropic 無隸屬關係。

## 參考與致謝

感謝以下專案為 Cage 提供基礎能力或設計參考：

- [agent-vm](https://github.com/sylvinus/agent-vm)：Cage 的二次開發底座。
- [Lima](https://github.com/lima-vm/lima)：Linux 虛擬機器管理。
- [gvisor-tap-vsock](https://github.com/containers/gvisor-tap-vsock)：虛擬機器網路能力。
- [sshocker](https://github.com/lima-vm/sshocker) 和 [SFTP](https://github.com/pkg/sftp)：虛擬機器與主機之間的檔案共享能力。
- [noVNC](https://github.com/novnc/noVNC)：瀏覽器中的遠端畫面與互動。
- [Mihomo](https://github.com/MetaCubeX/mihomo)：內建代理核心。
- [code-server](https://github.com/coder/code-server)：瀏覽器中的程式碼編輯器。
- [AIO Sandbox](https://github.com/agent-infra/sandbox)：工作環境和互動體驗參考。


<a href="https://www.star-history.com/?repos=Lee-xeo%2FClaude-Cage&amp;type=date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date&amp;theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" />
    <img alt="Claude Cage 的 GitHub Star 成長曲線" src="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" width="800" />
  </picture>
</a>
