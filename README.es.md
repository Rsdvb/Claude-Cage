<div align="center">

<h1>Claude Cage</h1>

**Entorno de desarrollo aislado**

<img src="assets/claude-cage-character.png" alt="Claude Cage" width="320" />

[![Release](https://img.shields.io/github/v/release/Lee-xeo/Claude-Cage?include_prereleases)](https://github.com/Lee-xeo/Claude-Cage/releases) [![CI](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml) [![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

<a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/daily?language=Go" alt="Trendshift — Go Repository of the Day" width="250" height="55" /></a> <a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/weekly?language=Go" alt="Trendshift — Go Repository of the Week" width="250" height="55" /></a>

**[Sitio web del proyecto](https://claude-cage.com/) · [Descargar aplicación](https://claude-cage.com/es/#download) · [Empezando](docs/getting-started.md) · [Preguntas frecuentes](docs/troubleshooting.md) · [Guía de uso (chino)](https://claude-cage.com/docs/)**

[简体中文](README.md) · [繁體中文](README.zh-TW.md) · [Русский](README.ru.md) · [العربية](README.ar.md) · [فارسی](README.fa.md) · [한국어](README.ko.md) · [Español](README.es.md) · [မြန်မာဘာသာ](README.my.md) · [Беларуская](README.be.md) · [پښتو](README.ps.md) · [دری](README.prs.md) · [አማርኛ](README.am.md) · [Soomaali](README.so.md) · [ትግርኛ](README.ti.md) · [Français](README.fr.md)

</div>

## Entorno de desarrollo aislado

Claude Cage reúne terminal, navegador, editor y Claude Desktop oficial en un mismo espacio de trabajo aislado. Usa una máquina virtual macOS en Mac, Windows en un PC compatible o un entorno Linux ligero. Comparte solo los proyectos que elijas, gestiona la red por separado, conecta tu servidor por SSH y continúa las tareas desde el teléfono.

**Interfaz principal de Cage:**

![Cage](www/public/workbench-2026-11.jpg)

## ¿Qué puedes hacer con la jaula?

- **Trabaja en una sola ventana.** Ejecute Claude Code en la terminal, revise los cambios en el editor y use el navegador Cage para obtener una vista previa de los servicios e iniciar sesión.
- **Desarrolla en tu sistema habitual.** Elige macOS, Windows o Linux como sistema invitado con sus herramientas, o clona el sistema del equipo para crear un entorno independiente y limpio. [Elegir el sistema invitado](docs/guide/guest-systems.md)
- **Controla lo que compartes.** Solo se comparten el espacio de trabajo predeterminado y las carpetas de proyectos abiertas explícitamente. El estado de inicio de sesión de Claude, los perfiles del navegador y las claves SSH del host no se importan automáticamente.
- **Vea el estado de su red.** La administración de proxy de red incorporada le permite importar suscripciones, seleccionar puntos finales y verificar la IP de salida. El entorno permanece fuera de línea hasta que se superan las comprobaciones críticas. Si el proxy de red no está disponible, no se recurre a una conexión directa como alternativa.
- **Continúe la próxima vez.** Las herramientas, los archivos y el estado de inicio de sesión persisten en el disco virtual. Las tareas en segundo plano pueden seguir ejecutándose cuando cierras la ventana.
- **Gestione varias cuentas de Claude Code.** Consulta cuotas, asigna cuentas a nuevas sesiones y, con autorización explícita, cambia de cuenta para continuar el historial. [Claude Code](docs/guide/accounts.md)
- **Continúa desde el teléfono.** Vincula Android o iPhone al ordenador mediante QR o código de emparejamiento. Retoma las mismas conversaciones de Claude Code, envía mensajes, aprueba herramientas y usa el navegador remoto o la ventana de Claude Desktop del ordenador.
- **Ejecute Cage en la nube.** Conéctese por SSH a Cage desplegado en un VPS Linux x86_64 compatible, sin KVM. Los datos y la salida de red están en el servidor.
- **Conecta por una red privada.** Inicia sesión en Tailscale dentro de Cage para conectar el equipo y el entorno por SSH en ambos sentidos. [Tailscale / SSH](docs/guide/tailscale.md)

## Inicio rápido

Envíe esta solicitud a su agente de IA (codex, dsh, workbuddy; probablemente no Claude Code 😁):

```markdown
Ayúdame a instalar https://github.com/Lee-xeo/Claude-Cage y configurar el proxy de red de Cage para usar el punto final xxx en mi instalación local de Clash.
```

## Descarga e instalación

Cage admite macOS, Windows y Linux: Apple Silicon con macOS 14+, Windows 10 22H2 / 11 x64 y Debian 12 / Ubuntu 22.04 / 24.04 x86_64. Intel Mac, Windows ARM64 y Linux ARM64 aún no son compatibles. [macOS](docs/getting-started.md) · [Windows](docs/windows.md) · [Linux](docs/linux.md).

1. Descargar aplicación: [macOS / Windows / Linux](https://claude-cage.com/es/#download).
2. macOS: Abra el paquete de instalación y arrastre **Cage** en **Aplicaciones**; luego ejecútelo desde Aplicaciones. Windows: [Guía de instalación](docs/windows.md#安装并开始工作) · Linux: [Guía de instalación](docs/linux.md#安装并开始工作).
3. Siga las instrucciones de la aplicación para configurar su propio proxy de red, iniciar el entorno y abrir un proyecto.

El instalador incluye el espacio de trabajo y el motor del proxy. No necesitas instalar antes Go, Node.js, Docker ni Clash. Al crear el entorno, elige el sistema invitado y su ubicación y prepara la imagen oficial y las herramientas. Cage no proporciona licencias del sistema, servicio de proxy ni cuentas Claude.

El código fuente y los instaladores son públicos. Las notas y sumas de comprobación están en [Releases](https://github.com/Lee-xeo/Claude-Cage/releases).

Las funciones dependen de la compilación; las capacidades y verificaciones del instalador público figuran en las notas de la versión.

macOS y las versiones firmadas de Windows se actualizan mediante el proxy elegido. La vista previa de Windows sin firma y Linux se actualizan manualmente con el instalador. Guarda y detén las tareas antes de actualizar; los datos personales se conservan. [Actualizar y desinstalar](docs/guide/maintenance.md)

## Comience en tres pasos

### 1. Seleccione un proxy de red

En la página de inicio, haga clic en "Configurar red y proxy" y luego en "Agregar fuente". Importe una suscripción YAML de Clash/Mihomo o un archivo YAML local y seleccione un punto final en "Nodos". Si no hay ningún proxy de red disponible, puede seleccionar explícitamente "Descargar esta suscripción a través de la red de host" o importar un archivo local directamente. Luego, administre los puntos finales en 'Configuración → Red y proxy'.

### 2. Iniciar el entorno

Al crear el entorno, elige el sistema invitado o la clonación del sistema del equipo y la ubicación de almacenamiento. Inicia el entorno Cage y sigue la instalación y las autorizaciones. El progreso y los errores aparecen en los detalles de la operación. Los siguientes inicios reutilizan el mismo entorno.

### 3. Abra un proyecto y ejecute Claude.

Haga clic en 'Abrir proyecto' para seleccionar una carpeta en su equipo, abra una nueva terminal e ingrese `claude`; luego autorice su propia cuenta en el navegador Cage. También puedes comenzar con el espacio de trabajo predeterminado sin preparar un proyecto. Para conocer los pasos completos, consulte [Empezando](docs/getting-started.md).

## Comprenda estos límites antes de usarlo

- Las ediciones y eliminaciones en proyectos compartidos afectan directamente a los archivos de su equipo. Utilice control de versiones o copias de seguridad. El directorio `.git` del proyecto permite lectura y escritura; puede preparar cambios, crear commits y cambiar de rama en Cage.
- Todos los proyectos comparten un entorno Cage, incluidas sus herramientas, cuentas y sistema de archivos macOS / Windows / Debian / Ubuntu 22.04. Los proyectos no están aislados unos de otros.
- Por defecto, las conexiones públicas salientes admiten IPv4 TCP y DNS usa el mismo proxy. UDP, IPv6, ICMP públicos y el acceso directo al equipo o a la red local están limitados. Al activar Tailscale, el acceso a los dispositivos sigue los permisos de la red privada.
- Las comprobaciones de red muestran la configuración actual y pruebas de muestra. No pueden garantizar que una cuenta evite las restricciones de la plataforma y no reemplazan una auditoría completa de seguridad de virtualización.

Para más detalles, consulte [Límites de red y datos](docs/getting-started.md#网络与数据边界).

## Documentación y aportaciones

Para consultar el índice completo, consulte [Índice de documentación](docs/README.md).

| tu tarea | Empieza aquí |
|---|---|
| macOS / Windows / Debian / Ubuntu 22.04 | [Elegir el sistema invitado](https://claude-cage.com/docs/guest-systems/) |
| Instale, configure un proxy de red y abra su primer proyecto | [Empezando](docs/getting-started.md) |
| Solucionar problemas de descargas, inicio, conectividad e inicio de sesión | [Preguntas frecuentes y recuperación](docs/troubleshooting.md) |
| Aprenda los controles del teclado, la apariencia, la persistencia de datos y el comportamiento de salida | [Uso diario](docs/getting-started.md#日常使用) |
| Informar un problema reproducible | [GitHub Issues](https://github.com/Lee-xeo/Claude-Cage/issues) |

El proyecto está evolucionando rápidamente. Únase al grupo comunitario para compartir comentarios; Los mantenedores de la comunidad son bienvenidos.

<p align="center">
  <a href="assets/community-qq.png"><img src="assets/community-qq.png" alt="QQ #2 · 1127457578" width="240" /></a>
  <br />QQ #2 · 1127457578
</p>

## Orígenes y licencias

La jaula se deriva de [agent-vm](https://github.com/sylvinus/agent-vm).

La licencia principal del proyecto es [MIT](LICENSE). Los componentes incluidos conservan sus respectivas licencias. Mihomo se distribuye como un ejecutable independiente bajo GPL-3.0. Ver también [NOTICE](NOTICE) · [Licencias de terceros](THIRD_PARTY_LICENSES) · [noVNC](ui/vendor/novnc/LICENSE.txt)。

Cage es un proyecto comunitario independiente y no está afiliado a Anthropic.

## Referencias y agradecimientos

Gracias a los siguientes proyectos por capacidades fundamentales y referencias de diseño:

- [agent-vm](https://github.com/sylvinus/agent-vm): la base aguas arriba de Cage.
- [Lima](https://github.com/lima-vm/lima): Gestión de máquinas virtuales Linux.
- [gvisor-tap-vsock](https://github.com/containers/gvisor-tap-vsock): redes de máquinas virtuales.
- [sshocker](https://github.com/lima-vm/sshocker) y [SFTP](https://github.com/pkg/sftp): intercambio de archivos entre la máquina virtual y el host.
- [noVNC](https://github.com/novnc/noVNC): visualización remota e interacción en el navegador.
- [Mihomo](https://github.com/MetaCubeX/mihomo): el núcleo de proxy de red integrado.
- [code-server](https://github.com/coder/code-server): el editor de código basado en navegador.
- [AIO Sandbox](https://github.com/agent-infra/sandbox): referencias para el ambiente de trabajo y el diseño de interacción.


<a href="https://www.star-history.com/?repos=Lee-xeo%2FClaude-Cage&amp;type=date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date&amp;theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" />
    <img alt="Claude Cage GitHub Star History" src="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" width="800" />
  </picture>
</a>
