<div align="center">

<h1>Claude Cage</h1>

**Espace de développement isolé**

<img src="assets/claude-cage-character.png" alt="Claude Cage" width="320" />

[![Release](https://img.shields.io/github/v/release/Lee-xeo/Claude-Cage?include_prereleases)](https://github.com/Lee-xeo/Claude-Cage/releases) [![CI](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml) [![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

<a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/daily?language=Go" alt="Trendshift — Go Repository of the Day" width="250" height="55" /></a> <a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/weekly?language=Go" alt="Trendshift — Go Repository of the Week" width="250" height="55" /></a>

**[Site Web du projet](https://claude-cage.com/) · [Télécharger l’application](https://claude-cage.com/fr/#download) · [Commencer](docs/getting-started.md) · [FAQ](docs/troubleshooting.md) · [Guide d’utilisation (chinois)](https://claude-cage.com/docs/)**

[简体中文](README.md) · [繁體中文](README.zh-TW.md) · [Русский](README.ru.md) · [العربية](README.ar.md) · [فارسی](README.fa.md) · [한국어](README.ko.md) · [Español](README.es.md) · [မြန်မာဘာသာ](README.my.md) · [Беларуская](README.be.md) · [پښتو](README.ps.md) · [دری](README.prs.md) · [አማርኛ](README.am.md) · [Soomaali](README.so.md) · [ትግርኛ](README.ti.md) · [Français](README.fr.md)

</div>

## Espace de développement isolé

Claude Cage réunit terminal, navigateur, éditeur et Claude Desktop officiel dans un même espace de travail isolé. Utilisez une machine virtuelle macOS sur Mac, Windows sur un PC compatible, ou un environnement Linux léger. Partagez les projets de votre choix, gérez le réseau séparément, connectez votre serveur par SSH et poursuivez vos tâches depuis votre téléphone.

**L'interface principale de Cage :**

![Cage](www/public/workbench-2026-11.jpg)

## Ce que vous pouvez faire avec Cage

- **Travaillez dans une seule fenêtre.** Exécutez Claude Code dans le terminal, examinez les modifications dans l'éditeur et utilisez le navigateur Cage pour prévisualiser les services et vous connecter.
- **Développez dans votre système habituel.** Choisissez un système invité macOS, Windows ou Linux avec ses outils, ou clonez le système hôte pour créer un environnement indépendant et propre. [Choisir le système invité](docs/guide/guest-systems.md)
- **Contrôlez ce que vous partagez.** Seuls l'espace de travail par défaut et les dossiers de projet explicitement ouverts sont partagés. L'état de connexion Claude de l'hôte, les profils de navigateur et les clés SSH ne sont pas importés automatiquement.
- **Consultez l'état de votre réseau.** La gestion intégrée du proxy réseau vous permet d'importer des abonnements, de sélectionner des points de terminaison et de vérifier l'adresse IP de sortie. L'environnement reste hors ligne jusqu'à ce que les vérifications critiques soient réussies. Si le proxy réseau n'est pas disponible, il n'est pas possible de recourir à une connexion directe.
- **Continuez la prochaine fois.** Les outils, les fichiers et l'état de connexion persistent sur le disque virtuel. Les tâches en arrière-plan peuvent continuer à s'exécuter lorsque vous fermez la fenêtre.
- **Gérez plusieurs comptes Claude Code.** Consultez les quotas, attribuez les comptes aux nouvelles sessions et, avec une autorisation explicite, changez de compte en conservant l’historique. [Claude Code](docs/guide/accounts.md)
- **Poursuivez depuis votre téléphone.** Associez Android ou iPhone à l’ordinateur par QR code ou code d’association. Retrouvez les mêmes conversations Claude Code, envoyez des messages, approuvez les outils et utilisez le navigateur distant ou la fenêtre Claude Desktop de l’ordinateur.
- **Exécutez Cage dans le cloud.** Connectez-vous par SSH à Cage déployé sur un VPS Linux x86_64 compatible, sans KVM. Les données et la sortie réseau restent sur le serveur.
- **Connectez-vous par un réseau privé.** Connectez Tailscale dans Cage pour accéder par SSH à l’hôte et à l’environnement dans les deux sens. [Tailscale / SSH](docs/guide/tailscale.md)

## Démarrage rapide

Envoyez cette demande à votre agent IA (codex, dsh, workbuddy — probablement pas Claude Code 😁) :

```markdown
Aide-moi à installer https://github.com/Lee-xeo/Claude-Cage et configurez le proxy réseau de Cage pour utiliser le point de terminaison xxx dans mon installation Clash locale.
```

## Téléchargement et installation

Cage prend en charge macOS, Windows et Linux : Apple Silicon avec macOS 14+, Windows 10 22H2 / 11 x64 et Debian 12 / Ubuntu 22.04 / 24.04 x86_64. Intel Mac, Windows ARM64 et Linux ARM64 ne sont pas encore pris en charge. [macOS](docs/getting-started.md) · [Windows](docs/windows.md) · [Linux](docs/linux.md).

1. Télécharger l’application: [macOS / Windows / Linux](https://claude-cage.com/fr/#download).
2. macOS: Ouvrez le package d'installation et faites glisser **Cage** dans **Applications**; puis lancez-le depuis Applications. Windows: [Guide d’installation](docs/windows.md#安装并开始工作) · Linux: [Guide d’installation](docs/linux.md#安装并开始工作).
3. Suivez les instructions de l'application pour configurer votre propre proxy réseau, démarrer l'environnement et ouvrir un projet.

L’installateur comprend l’espace de travail et le moteur de proxy. Inutile d’installer Go, Node.js, Docker ou Clash au préalable. À la création, choisissez le système invité et son emplacement, puis préparez l’image officielle et les outils. Cage ne fournit ni licence de système, ni service proxy, ni compte Claude.

Le code source et les installateurs sont publics. Consultez les notes et les sommes de contrôle dans [Releases](https://github.com/Lee-xeo/Claude-Cage/releases).

Les fonctions dépendent de la compilation utilisée ; les fonctions et les vérifications de l’installateur public sont précisées dans les notes de version.

macOS et les versions Windows signées se mettent à jour via le proxy choisi. La préversion Windows non signée et Linux utilisent l’installateur manuellement. Enregistrez et arrêtez les tâches avant la mise à jour ; les données personnelles sont conservées. [Mise à jour et désinstallation](docs/guide/maintenance.md)

## Commencez en trois étapes

### 1. Sélectionnez un proxy réseau

Sur la page d'accueil, cliquez sur « Configurer le réseau et le proxy », puis sur « Ajouter une source ». Importez un abonnement YAML Clash/Mihomo ou un fichier YAML local et sélectionnez un point de terminaison sous « Nœuds ». Si aucun proxy réseau n'est disponible, vous pouvez sélectionner explicitement « Télécharger cet abonnement via le réseau hôte » ou importer directement un fichier local. Plus tard, gérez les points de terminaison sous « Paramètres → Réseau et proxy ».

### 2. Démarrez l'environnement

À la première création, choisissez le système invité ou le clonage du système hôte, puis le stockage. Démarrez l’environnement Cage et suivez l’installation et les autorisations. La progression et les erreurs figurent dans les détails de l’opération. Les lancements suivants réutilisent le même environnement.

### 3. Ouvrez un projet et exécutez Claude

Cliquez sur « Ouvrir le projet » pour sélectionner un dossier sur votre ordinateur, ouvrez un nouveau terminal et entrez `claude`; puis autorisez votre propre compte dans le navigateur Cage. Vous pouvez également commencer avec l'espace de travail par défaut sans préparer de projet. Pour les étapes complètes, voir [Commencer](docs/getting-started.md).

## Comprendre ces limites avant utilisation

- Les modifications et suppressions dans les projets partagés affectent directement les fichiers sur votre ordinateur. Utilisez le contrôle de version ou les sauvegardes. Le dossier `.git` du projet est accessible en lecture et en écriture ; vous pouvez préparer des modifications, créer des commits et changer de branche dans Cage.
- Tous les projets partagent un environnement Cage, y compris ses outils, ses comptes et son système de fichiers macOS / Windows / Debian / Ubuntu 22.04. Les projets ne sont pas isolés les uns des autres.
- Par défaut, les sorties publiques prennent en charge IPv4 TCP et DNS passe par le même proxy. UDP, IPv6, ICMP publics et les accès directs à l’hôte ou au réseau local sont limités. En activant Tailscale, les accès aux appareils suivent les règles du réseau privé.
- Les vérifications du réseau montrent la configuration actuelle et des échantillons de preuves. Ils ne peuvent pas garantir qu’un compte évitera les restrictions de la plateforme et ne remplacent pas un audit complet de sécurité de la virtualisation.

Pour plus de détails, voir [Limites du réseau et des données](docs/getting-started.md#网络与数据边界).

## Documentation et contributions

Pour l'index complet, voir [Index des documents](docs/README.md).

| Votre tâche | Commencez ici |
|---|---|
| macOS / Windows / Debian / Ubuntu 22.04 | [Choisir le système invité](https://claude-cage.com/docs/guest-systems/) |
| Installez, configurez un proxy réseau et ouvrez votre premier projet | [Commencer](docs/getting-started.md) |
| Dépanner les téléchargements, le démarrage, la connectivité et la connexion | [FAQ et récupération](docs/troubleshooting.md) |
| Apprenez les commandes du clavier, l'apparence, la persistance des données et le comportement de sortie | [Utilisation quotidienne](docs/getting-started.md#日常使用) |
| Signaler un problème reproductible | [GitHub Issues](https://github.com/Lee-xeo/Claude-Cage/issues) |

Le projet évolue rapidement. Rejoignez le groupe communautaire pour partager vos commentaires ; les responsables de la communauté sont les bienvenus.

<p align="center">
  <a href="assets/community-qq.png"><img src="assets/community-qq.png" alt="QQ #2 · 1127457578" width="240" /></a>
  <br />QQ #2 · 1127457578
</p>

## Origines et licences

La cage est dérivée de [agent-vm](https://github.com/sylvinus/agent-vm).

La licence principale du projet est [MIT](LICENSE). Les composants groupés conservent leurs licences respectives. Mihomo est distribué sous forme d'exécutable distinct sous GPL-3.0. Voir aussi [NOTICE](NOTICE) · [Licences tierces](THIRD_PARTY_LICENSES) · [noVNC](ui/vendor/novnc/LICENSE.txt)。

Cage est un projet communautaire indépendant et n'est pas affilié à Anthropic.

## Références et remerciements

Merci aux projets suivants pour leurs capacités fondamentales et leurs références de conception :

- [agent-vm](https://github.com/sylvinus/agent-vm): la fondation amont de Cage.
- [Lima](https://github.com/lima-vm/lima): Gestion des machines virtuelles Linux.
- [gvisor-tap-vsock](https://github.com/containers/gvisor-tap-vsock): mise en réseau de machines virtuelles.
- [sshocker](https://github.com/lima-vm/sshocker) et [SFTP](https://github.com/pkg/sftp): partage de fichiers entre la machine virtuelle et l'hôte.
- [noVNC](https://github.com/novnc/noVNC): affichage à distance et interaction dans le navigateur.
- [Mihomo](https://github.com/MetaCubeX/mihomo): le noyau proxy réseau intégré.
- [code-server](https://github.com/coder/code-server): l'éditeur de code basé sur un navigateur.
- [AIO Sandbox](https://github.com/agent-infra/sandbox): références pour l'environnement de travail et le design d'interaction.


<a href="https://www.star-history.com/?repos=Lee-xeo%2FClaude-Cage&amp;type=date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date&amp;theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" />
    <img alt="Claude Cage GitHub Star History" src="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" width="800" />
  </picture>
</a>
