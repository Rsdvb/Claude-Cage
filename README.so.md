<div align="center">

<h1>Claude Cage</h1>

**Deegaan horumarineed oo go’doonsan**

<img src="assets/claude-cage-character.png" alt="Claude Cage" width="320" />

[![Release](https://img.shields.io/github/v/release/Lee-xeo/Claude-Cage?include_prereleases)](https://github.com/Lee-xeo/Claude-Cage/releases) [![CI](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml) [![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

<a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/daily?language=Go" alt="Trendshift — Go Repository of the Day" width="250" height="55" /></a> <a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/weekly?language=Go" alt="Trendshift — Go Repository of the Week" width="250" height="55" /></a>

**[Mareegta mashruuca](https://claude-cage.com/) · [Soo dejiso abka](https://claude-cage.com/so/#download) · [Bilaabidda](docs/getting-started.md) · [FAQ](docs/troubleshooting.md) · [Hagaha isticmaalka (Shiine)](https://claude-cage.com/docs/)**

[简体中文](README.md) · [繁體中文](README.zh-TW.md) · [Русский](README.ru.md) · [العربية](README.ar.md) · [فارسی](README.fa.md) · [한국어](README.ko.md) · [Español](README.es.md) · [မြန်မာဘာသာ](README.my.md) · [Беларуская](README.be.md) · [پښتو](README.ps.md) · [دری](README.prs.md) · [አማርኛ](README.am.md) · [Soomaali](README.so.md) · [ትግርኛ](README.ti.md) · [Français](README.fr.md)

</div>

## Deegaan horumarineed oo go’doonsan

Claude Cage wuxuu terminalka, biraawsarka, tafatiraha iyo Claude Desktop-ka rasmiga ah ku mideeyaa hal goob shaqo oo go’doonsan. Ku isticmaal macOS mashiin macmal ah oo Mac ah, Windows kombiyuutar la jaanqaadaya, ama deegaan Linux fudud. Wadaag mashaariicda aad doorato, shabakadda si gaar ah u maamul, serverkaaga SSH ku xir, hawlahana telefoonka ka sii wad.

**Interface ugu weyn ee Cage:**

![Cage](www/public/workbench-2026-11.jpg)

## Maxaad ku samayn kartaa Cage

- **Ka shaqee hal daaqad.** Orod Claude Code gudaha terminalka, dib u eeg isbeddelada tifaftiraha, oo isticmaal browserka Cage si aad u aragto adeegyada oo aad u soo gasho.
- **Ku horumari nidaamka aad taqaan.** Dooro macOS, Windows ama Linux marti ah iyo qalabkiisa, ama nuqul ka samee nidaamka martigeliyaha si aad u abuurto deegaan nadiif ah oo madax bannaan. [Doorashada nidaamka martida](docs/guide/guest-systems.md)
- **Xakamee waxa aad wadaagto.** Kaliya goobta shaqada ee caadiga ah iyo galka mashruuca ee si cad loo furay ayaa la wadaagaa. Xaaladda soo galitaanka Claude ee martida loo yahay, boggaga browserka, iyo furayaasha SSH si toos ah looma soo dejiyo.
- **Eeg heerka shabakadaada** Maaraynta wakiilnimada shabakada ee ku dhex jirta ayaa kuu ogolaanaysa inaad soo dejiso rukunsiyada, dooro meelaha dhamaadka, oo aad hubiso IP-ga ka baxa. Deegaanku internetka kuma xirmayo ilaa hubinnada muhiimka ahi ay si guul leh u gudbaan. Haddii wakiillada shabakadu maqan tahay, ma dhacayso xiriir toos ah.
- **Sii wad marka xigta** Qalabka, faylalka, iyo xaaladda gelitaanku waxay ku sii jiraan saxanka casriga ah. Hawlaha asalka ahi waxay sii socon karaan markaad daaqadda xidhid.
- **Maamul akoonno badan oo Claude Code ah.** Eeg qoondada, u qoondee akoonno kulammada cusub, kuna beddel akoonka oggolaansho cad si aad taariikhda u sii wadato. [Claude Code](docs/guide/accounts.md)
- **Telefoonka ka sii wad.** Android ama iPhone ku lammaani kombiyuutarka adigoo adeegsanaya QR ama koodhka lammaaninta. Sii wad isla wada sheekaysiyada Claude Code, dir farriimo, oggolow qalabka, isticmaalna biraawsarka fog ama daaqadda Claude Desktop ee kombiyuutarka.
- **Cage ku socodsii daruurta.** SSH kaga xirmo Cage lagu rakibay VPS ku habboon oo Linux x86_64 ah, adigoon u baahnayn KVM. Xogta iyo halka shabakadu ka baxdo waxay ku yaallaan serverka.
- **Ku xidh shabakad gaar ah.** Gudaha Cage ka gal Tailscale si martigeliyaha iyo deegaanku SSH ugu xidhmaan labada jiho. [Tailscale / SSH](docs/guide/tailscale.md)

## Bilaw degdeg ah

U dir codsigan wakiilkaaga AI (codex, dsh, workbuddy - malaha maaha Claude Code 😁):

```markdown
Iga caawi in aan rakibo https://github.com/Lee-xeo/Claude-Cage oo habee wakiilka shabakada Cage si aad u isticmaasho xxx dhamaadka rakibida xaafadayda.
```

## Soo deji oo rakib

Cage wuxuu taageeraa macOS, Windows iyo Linux: Apple Silicon oo leh macOS 14+, Windows 10 22H2 / 11 x64 iyo Debian 12 / Ubuntu 22.04 / 24.04 x86_64. Intel Mac, Windows ARM64 iyo Linux ARM64 weli lama taageero. [macOS](docs/getting-started.md) · [Windows](docs/windows.md) · [Linux](docs/linux.md).

1. Soo dejiso abka: [macOS / Windows / Linux](https://claude-cage.com/so/#download).
2. macOS: Fur xirmada rakibaadda oo jiid **Cage** galay **Codsiyada**; ka dibna ka bilow Applications. Windows: [Hagaha rakibidda](docs/windows.md#安装并开始工作) · Linux: [Hagaha rakibidda](docs/linux.md#安装并开始工作).
3. Raac tilmaamaha app-ka si aad u habayso wakiilkaaga shabakad, bilaw deegaanka, oo aad u furto mashruuc.

Rakibuhu wuxuu ka kooban yahay goobta shaqada iyo matoorka proxy-ga. Uma baahnid inaad horay u rakibto Go, Node.js, Docker ama Clash. Markaad deegaanka abuurayso, dooro nidaamka martida iyo goobta kaydinta, dabadeed diyaari sawirka rasmiga ah iyo qalabka. Cage ma bixiyo ruqsad nidaam, adeeg proxy ama akoon Claude.

Koodhka isha iyo rakibayaashu waa dadweyne. Qoraallada noocyada iyo checksum-yada ka eeg [Releases](https://github.com/Lee-xeo/Claude-Cage/releases).

Awooduhu waxay ku xiran yihiin build-ka la isticmaalayo; awoodaha iyo xadka tijaabada rakibaha dadweynaha waxay ku qoran yihiin qoraalka sii-deynta.

macOS iyo Windows saxeexan waxay ku cusboonaadaan proxy-ga la doortay. Windows preview aan saxeexnayn iyo Linux gacanta ayaa loogu cusboonaysiiyaa rakibaha. Marka hore kaydi oo jooji hawlaha; xogta shakhsiga ah waa la hayaa. [Cusboonaysiin iyo saarid](docs/guide/maintenance.md)

## Ku bilow saddex tallaabo

### 1. Dooro wakiil shabakad

Bogga guriga, dhagsii 'Configure network and proxy', ka dibna 'Ku dar isha'. Soo rar isdiiwaangelinta Clash/Mihomo YAML ama faylka YAML maxalli ah, oo dooro dhammaadka hoosta 'Nodes'. Haddii aysan jirin wakiil shabakadeed, waxaad si cad u dooran kartaa 'Soo deji rukuntan iyada oo loo marayo shabakadda martida loo yahay', ama si toos ah u soo deji faylka maxalliga ah. Ka dib, maamul meelaha dhamaadka ee hoos yimaada 'Settings → Network and proxy'.

### 2. Bilow deegaanka

Marka ugu horreysa dooro nidaamka martida ama nuqulka nidaamka martigeliyaha iyo goobta kaydinta. Bilow deegaanka Cage oo raac rakibidda iyo oggolaanshaha. Horumarka iyo khaladaadku waxay ka muuqdaan faahfaahinta hawsha. Bilowga dambe wuxuu isticmaalaa isla deegaanka.

### 3. Fur mashruuc oo maamul Claude

Guji 'Mashruuca Furan' si aad u dooratid galka kombiyuutarka kaaga, fur terminal cusub, oo gal `claude`; dabadeed ogolaw akoonkaaga browserka Cage. Waxaad sidoo kale ku bilaabi kartaa goobta shaqada ee caadiga ah adigoon diyaarin mashruuc. Tallaabooyinka buuxa, eeg [Bilaabidda](docs/getting-started.md).

## Fahmo xadkan ka hor intaanad isticmaalin

- Wax ka beddelka iyo tirtiridda mashaariicda la wadaago waxay si toos ah u saameeyaan faylasha kombiyuutarka kaaga. Isticmaal xakamaynta nooca ama kaydinta. Galka `.git` ee mashruuca waa la akhriyi karaa lana qori karaa; Cage gudaheeda waxaad ku diyaarin kartaa isbeddelada, samayn kartaa commits, kuna beddeli kartaa branches.
- Dhammaan mashruucyada waxay wadaagaan hal deegaan Cage, oo ay ku jiraan qalabkeeda, xisaabaadka, iyo nidaamka faylka macOS / Windows / Debian / Ubuntu 22.04. Mashruucyadu maaha kuwo midba midka kale ka soocan.
- Sida caadiga ah bixitaanka internetku wuxuu taageeraa IPv4 TCP, DNS-na isla proxy-ga ayuu maraa. UDP, IPv6 iyo ICMP dadweyne iyo xidhiidh toos ah oo martigeliye ama LAN ah waa xaddidan yihiin. Markaad si cad u hawlgeliso Tailscale, gelitaanka aaladuhu wuxuu raacaa oggolaanshaha shabakadda gaarka ah.
- Jeegaga shabakada waxay tusinaysaa qaabka hadda jira iyo cadaymaha la muunada ah. Ma dammaanad qaadi karaan in akoonku uu ka fogaanayo xannibaadaha goobta, mana beddelaan xisaab-xidhka amniga oo dhammaystiran.

Faahfaahinta, eeg [Xuduudaha shabakada iyo xogta](docs/getting-started.md#网络与数据边界).

## Dukumentiyada iyo wax ku biirinta

Si aad u hesho tusmada buuxda, eeg [Tusmada dhokumentiyada](docs/README.md).

| Hawshaada | Halkan ka bilow |
|---|---|
| macOS / Windows / Debian / Ubuntu 22.04 | [Doorashada nidaamka martida](https://claude-cage.com/docs/guest-systems/) |
| Ku rakib, habee wakiil shabakad, oo fur mashruucaaga ugu horreeya | [Bilaabidda](docs/getting-started.md) |
| Xallinta soo dejinta, bilaabista, isku xidhka, iyo soo gelista | [FAQ iyo soo kabashada](docs/troubleshooting.md) |
| Baro kontaroolada kiiboodhka, muuqaalka, adkaysiga xogta, iyo habdhaqanka bixista | [Isticmaalka maalinlaha ah](docs/getting-started.md#日常使用) |
| Ka warbixi arrin dib loo soo saari karo | [GitHub Issues](https://github.com/Lee-xeo/Claude-Cage/issues) |

Mashruucu si xawli ah ayuu u socdaa. Ku biir kooxda bulshada si aad ula wadaagto jawaab celinta; ilaaliyayaasha bulshada waa la soo dhawaynayaa.

<p align="center">
  <a href="assets/community-qq.png"><img src="assets/community-qq.png" alt="QQ #2 · 1127457578" width="240" /></a>
  <br />QQ #2 · 1127457578
</p>

## Asalka iyo shatiga

Qafis waxaa laga soo qaatay [agent-vm](https://github.com/sylvinus/agent-vm).

Shatiga ugu weyn ee mashruuca waa [MIT](LICENSE). Qaybaha xidhxidhan waxay haystaan ​​shatiyadooda. Mihomo waxaa loo qaybiyaa sidii mid gaar ah oo la fulin karo iyadoo hoos timaada GPL-3.0. Sidoo kale arag [NOTICE](NOTICE) · [Shatiyada xisbiga saddexaad](THIRD_PARTY_LICENSES) · [noVNC](ui/vendor/novnc/LICENSE.txt)。

Cage waa mashruuc bulsho oo madax banaan oo aan xidhiidh la lahayn Anthropic.

## Tix-raac iyo mahadnaq

Waad ku mahadsan tahay mashaariicda soo socota ee kartida aasaasiga ah iyo tixraacyada naqshadeynta:

- [agent-vm](https://github.com/sylvinus/agent-vm): aasaaska sare ee Cage.
- [Lima](https://github.com/lima-vm/lima): Maareynta mashiinka farsamada ee Linux.
- [gvisor-tap-vsock](https://github.com/containers/gvisor-tap-vsock): isku xirka mashiinka farsamada.
- [sshocker](https://github.com/lima-vm/sshocker) iyo [SFTP](https://github.com/pkg/sftp): wadaagista faylka u dhexeeya mashiinka farsamada iyo martigeliyaha.
- [noVNC](https://github.com/novnc/noVNC): soo bandhigida fog iyo isdhexgalka browserka.
- [Mihomo](https://github.com/MetaCubeX/mihomo):-ku-xidhan xudunta wakiilnimada.
- [code-server](https://github.com/coder/code-server): tifaftiraha koodhka ku salaysan browserka.
- [AIO Sandbox](https://github.com/agent-infra/sandbox): tixraacyada jawiga shaqada iyo naqshadeynta isdhexgalka.


<a href="https://www.star-history.com/?repos=Lee-xeo%2FClaude-Cage&amp;type=date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date&amp;theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" />
    <img alt="Claude Cage GitHub Star History" src="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" width="800" />
  </picture>
</a>
