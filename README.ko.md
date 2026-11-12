<div align="center">

<h1>Claude Cage</h1>

**격리된 개발 작업 공간**

<img src="assets/claude-cage-character.png" alt="Claude Cage" width="320" />

[![Release](https://img.shields.io/github/v/release/Lee-xeo/Claude-Cage?include_prereleases)](https://github.com/Lee-xeo/Claude-Cage/releases) [![CI](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml) [![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

<a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/daily?language=Go" alt="Trendshift — Go Repository of the Day" width="250" height="55" /></a> <a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/weekly?language=Go" alt="Trendshift — Go Repository of the Week" width="250" height="55" /></a>

**[프로젝트 웹사이트](https://claude-cage.com/) · [앱 다운로드](https://claude-cage.com/ko/#download) · [시작하기](docs/getting-started.md) · [FAQ](docs/troubleshooting.md) · [사용 문서 (중국어)](https://claude-cage.com/docs/)**

[简体中文](README.md) · [繁體中文](README.zh-TW.md) · [Русский](README.ru.md) · [العربية](README.ar.md) · [فارسی](README.fa.md) · [한국어](README.ko.md) · [Español](README.es.md) · [မြန်မာဘာသာ](README.my.md) · [Беларуская](README.be.md) · [پښتو](README.ps.md) · [دری](README.prs.md) · [አማርኛ](README.am.md) · [Soomaali](README.so.md) · [ትግርኛ](README.ti.md) · [Français](README.fr.md)

</div>

## 격리된 개발 작업 공간

Claude Cage는 터미널, 브라우저, 편집기와 공식 Claude Desktop을 하나의 격리된 작업 공간에 모읍니다. Mac에서는 macOS 가상 머신을, 지원되는 Windows PC에서는 Windows 가상 머신을 사용하거나 가벼운 Linux 환경을 선택하세요. 필요한 프로젝트만 공유하고 네트워크를 따로 관리하며, SSH로 서버에 연결하거나 휴대폰에서 작업을 이어갈 수 있습니다.

**Cage의 기본 인터페이스:**

![Cage](www/public/workbench-2026-11.jpg)

## 케이지로 할 수 있는 일

- **단일 창에서 작업합니다.** 터미널에서 Claude Code를 실행하고, 편집기에서 변경 사항을 검토하고, Cage 브라우저를 사용하여 서비스를 미리 보고 로그인하세요.
- **익숙한 시스템에서 개발하세요.** macOS, Windows 또는 Linux 게스트에서 해당 플랫폼의 도구를 사용하거나, 호스트 시스템을 복제해 깨끗하고 독립적인 환경을 만드세요. [게스트 시스템 선택](docs/guide/guest-systems.md)
- **공유하는 내용을 제어하세요.** 기본 작업공간과 명시적으로 열린 프로젝트 폴더만 공유됩니다. 호스트의 Claude 로그인 상태, 브라우저 프로필 및 SSH 키는 자동으로 가져오지 않습니다.
- **네트워크 상태를 확인하세요.** 내장된 네트워크 프록시 관리를 통해 구독을 가져오고, 엔드포인트를 선택하고, 종료 IP를 확인할 수 있습니다. 중요한 검사가 통과될 때까지 환경은 오프라인 상태로 유지됩니다. 네트워크 프록시를 사용할 수 없는 경우 직접 연결로 대체할 수 없습니다.
- **다음에 계속하세요.** 도구, 파일 및 로그인 상태는 가상 디스크에 유지됩니다. 창을 닫아도 백그라운드 작업이 계속 실행될 수 있습니다.
- **Claude Code 계정을 여러 개 관리하세요.** 사용 한도를 확인하고 새 세션에 계정을 배정하며, 명시적으로 허용한 뒤 계정을 바꾸어 대화 기록을 이어가세요. [Claude Code](docs/guide/accounts.md)
- **휴대폰에서 이어서 작업하세요.** QR 코드나 페어링 코드로 Android·iPhone을 컴퓨터에 연결하세요. 같은 Claude Code 대화를 이어가고 메시지를 보내거나 도구 실행을 승인하며, 원격 브라우저와 컴퓨터의 Claude Desktop 창을 조작할 수 있습니다.
- **클라우드에서 Cage를 실행하세요.** 요구 조건을 충족하는 Linux x86_64 VPS에 Cage를 배포하고 SSH로 연결합니다. KVM은 필요 없으며 데이터와 네트워크 출구는 서버에 있습니다.
- **사설망으로 연결하세요.** Cage 안에서 Tailscale에 로그인해 호스트와 환경 사이에 양방향 SSH 연결을 사용하세요. [Tailscale / SSH](docs/guide/tailscale.md)

## 빠른 시작

이 요청을 AI 에이전트(codex, dsh, workbuddy — Claude Code는 아닐 수도 있음 😁) 에 보냅니다.

```markdown
설치 도와주세요 https://github.com/Lee-xeo/Claude-Cage 로컬 Clash 설치에서 xxx 엔드포인트를 사용하도록 Cage의 네트워크 프록시를 구성합니다.
```

## 다운로드 및 설치

Cage는 macOS, Windows, Linux를 지원합니다. 대상은 macOS 14+의 Apple Silicon Mac, Windows 10 22H2 / 11 x64, Debian 12 / Ubuntu 22.04 / 24.04 x86_64. Intel Mac, Windows ARM64, Linux ARM64는 아직 지원하지 않습니다. [macOS](docs/getting-started.md) · [Windows](docs/windows.md) · [Linux](docs/linux.md).

1. 앱 다운로드: [macOS / Windows / Linux](https://claude-cage.com/ko/#download).
2. macOS: 설치 패키지를 열고 드래그하세요. **Cage** 으로 **응용**; 그런 다음 응용 프로그램에서 실행하십시오. Windows: [설치 안내](docs/windows.md#安装并开始工作) · Linux: [설치 안내](docs/linux.md#安装并开始工作).
3. 앱 내 지침에 따라 자신의 네트워크 프록시를 구성하고, 환경을 시작하고, 프로젝트를 엽니다.

설치 패키지에는 작업 공간과 프록시 엔진이 포함됩니다. Go, Node.js, Docker, Clash를 미리 설치할 필요가 없습니다. 환경을 만들 때 게스트 시스템과 저장 위치를 선택하고 공식 이미지와 개발 도구를 준비하세요. Cage는 운영체제 라이선스, 프록시 서비스나 Claude 계정을 제공하지 않습니다.

소스 코드와 설치 패키지는 공개되어 있습니다. 버전 안내와 체크섬은 [Releases](https://github.com/Lee-xeo/Claude-Cage/releases)에서 확인하세요.

기능은 사용 중인 빌드에 따라 다르며, 공개 설치 파일의 기능과 검증 범위는 릴리스 설명을 기준으로 합니다.

macOS와 서명된 Windows 버전은 선택한 프록시를 통해 업데이트합니다. 서명되지 않은 Windows 미리 보기와 Linux는 설치 파일로 수동 업데이트합니다. 먼저 작업을 저장하고 중지하세요. 개인 데이터는 유지됩니다. [업데이트 및 제거](docs/guide/maintenance.md)

## 3단계로 시작하세요

### 1. 네트워크 프록시를 선택하세요

홈페이지에서 '네트워크 및 프록시 구성'을 클릭한 다음 '소스 추가'를 클릭하세요. Clash/Mihomo YAML 구독 또는 로컬 YAML 파일을 가져오고 '노드'에서 엔드포인트를 선택합니다. 네트워크 프록시를 사용할 수 없는 경우 '호스트 네트워크를 통해 이 구독 다운로드'를 명시적으로 선택하거나 로컬 파일을 직접 가져올 수 있습니다. 나중에 '설정 → 네트워크 및 프록시'에서 엔드포인트를 관리하세요.

### 2. 환경 시작

처음 만들 때 게스트 시스템 또는 호스트 시스템 복제와 저장 위치를 선택하세요. Cage 환경을 시작하고 설치와 권한 안내를 따르세요. 작업 상세에서 진행 상태와 오류를 볼 수 있습니다. 이후에는 같은 환경을 다시 사용합니다.

### 3. 프로젝트를 열고 Claude를 실행하세요.

컴퓨터에서 폴더를 선택하려면 '프로젝트 열기'를 클릭하고 새 터미널을 열고 다음을 입력하세요. `claude`; 그런 다음 Cage 브라우저에서 자신의 계정을 인증하세요. 프로젝트를 준비하지 않고 기본 작업공간으로 시작할 수도 있습니다. 전체 단계는 다음을 참조하세요. [시작하기](docs/getting-started.md).

## 사용하기 전에 이러한 제한 사항을 이해하십시오.

- 공유 프로젝트의 편집 및 삭제는 컴퓨터의 파일에 직접적인 영향을 미칩니다. 버전 제어 또는 백업을 사용하십시오. 프로젝트의 `.git`은 읽기와 쓰기가 가능하며 Cage에서 변경 사항을 스테이징하고 커밋하거나 브랜치를 전환할 수 있습니다.
- 모든 프로젝트는 도구, 계정 및 macOS / Windows / Debian / Ubuntu 22.04 파일 시스템을 포함하여 하나의 Cage 환경을 공유합니다. 프로젝트는 서로 격리되지 않습니다.
- 기본 인터넷 아웃바운드는 IPv4 TCP를 지원하며 DNS는 같은 프록시를 사용합니다. 인터넷 UDP, IPv6, ICMP와 호스트·LAN 직접 연결은 제한됩니다. Tailscale을 직접 활성화하면 사설망 권한에 따라 기기에 연결할 수 있습니다.
- 네트워크 검사에는 현재 구성과 샘플링된 증거가 표시됩니다. 계정이 플랫폼 제한을 피할 것이라고 보장할 수 없으며 완전한 가상화 보안 감사를 대체하지 않습니다.

자세한 내용은 다음을 참조하세요. [네트워크 및 데이터 경계](docs/getting-started.md#网络与数据边界).

## 문서화 및 기여

전체 색인을 보려면 다음을 참조하세요. [문서 색인](docs/README.md).

| 귀하의 작업 | 여기서 시작하세요 |
|---|---|
| macOS / Windows / Debian / Ubuntu 22.04 | [게스트 시스템 선택](https://claude-cage.com/docs/guest-systems/) |
| 네트워크 프록시 설치, 구성 및 첫 번째 프로젝트 열기 | [시작하기](docs/getting-started.md) |
| 다운로드, 시작, 연결 및 로그인 문제 해결 | [FAQ 및 복구](docs/troubleshooting.md) |
| 키보드 컨트롤, 모양, 데이터 지속성 및 종료 동작 알아보기 | [매일 사용](docs/getting-started.md#日常使用) |
| 재현 가능한 문제 보고 | [GitHub Issues](https://github.com/Lee-xeo/Claude-Cage/issues) |

이 프로젝트는 빠르게 발전하고 있습니다. 피드백을 공유하려면 커뮤니티 그룹에 가입하세요. 커뮤니티 관리자를 환영합니다.

<p align="center">
  <a href="assets/community-qq.png"><img src="assets/community-qq.png" alt="QQ #2 · 1127457578" width="240" /></a>
  <br />QQ #2 · 1127457578
</p>

## 원산지 및 라이선스

케이지는 다음에서 파생됩니다. [agent-vm](https://github.com/sylvinus/agent-vm).

프로젝트의 기본 라이센스는 다음과 같습니다. [MIT](LICENSE). 번들 구성 요소는 해당 라이센스를 유지합니다. Mihomo는 GPL-3.0에 따라 별도의 실행 파일로 배포됩니다. 또한 참조하십시오 [NOTICE](NOTICE) · [타사 라이선스](THIRD_PARTY_LICENSES) · [noVNC](ui/vendor/novnc/LICENSE.txt)。

Cage는 독립적인 커뮤니티 프로젝트이며 Anthropic과 관련이 없습니다.

## 참고자료 및 감사의 말

기본 기능 및 설계 참조에 대한 다음 프로젝트에 감사드립니다.

- [agent-vm](https://github.com/sylvinus/agent-vm): Cage의 업스트림 기반.
- [Lima](https://github.com/lima-vm/lima): Linux 가상 머신 관리.
- [gvisor-tap-vsock](https://github.com/containers/gvisor-tap-vsock): 가상 머신 네트워킹.
- [sshocker](https://github.com/lima-vm/sshocker) 그리고 [SFTP](https://github.com/pkg/sftp): 가상 머신과 호스트 간의 파일 공유.
- [noVNC](https://github.com/novnc/noVNC): 브라우저에서의 원격 표시 및 상호 작용.
- [Mihomo](https://github.com/MetaCubeX/mihomo): 내장된 네트워크 프록시 코어입니다.
- [code-server](https://github.com/coder/code-server): 브라우저 기반 코드 편집기입니다.
- [AIO Sandbox](https://github.com/agent-infra/sandbox): 작업 환경 및 인터랙션 디자인에 대한 참고 자료입니다.


<a href="https://www.star-history.com/?repos=Lee-xeo%2FClaude-Cage&amp;type=date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date&amp;theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" />
    <img alt="Claude Cage GitHub Star History" src="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" width="800" />
  </picture>
</a>
