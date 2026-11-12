<div align="center">

<h1>Claude Cage</h1>

**محیط توسعهٔ ایزوله**

<img src="assets/claude-cage-character.png" alt="Claude Cage" width="320" />

[![Release](https://img.shields.io/github/v/release/Lee-xeo/Claude-Cage?include_prereleases)](https://github.com/Lee-xeo/Claude-Cage/releases) [![CI](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml) [![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

<a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/daily?language=Go" alt="Trendshift — Go Repository of the Day" width="250" height="55" /></a> <a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/weekly?language=Go" alt="Trendshift — Go Repository of the Week" width="250" height="55" /></a>

**[وب سایت پروژه](https://claude-cage.com/) · [دانلود برنامه](https://claude-cage.com/fa/#download) · [شروع کردن](docs/getting-started.md) · [سوالات متداول](docs/troubleshooting.md) · [راهنمای استفاده (چینی)](https://claude-cage.com/docs/)**

[简体中文](README.md) · [繁體中文](README.zh-TW.md) · [Русский](README.ru.md) · [العربية](README.ar.md) · [فارسی](README.fa.md) · [한국어](README.ko.md) · [Español](README.es.md) · [မြန်မာဘာသာ](README.my.md) · [Беларуская](README.be.md) · [پښتو](README.ps.md) · [دری](README.prs.md) · [አማርኛ](README.am.md) · [Soomaali](README.so.md) · [ትግርኛ](README.ti.md) · [Français](README.fr.md)

</div>

## محیط توسعهٔ ایزوله

Claude Cage ترمینال، مرورگر، ویرایشگر و برنامهٔ رسمی Claude Desktop را در یک فضای کاری جدا گرد هم می‌آورد. روی Mac از ماشین مجازی macOS، روی رایانهٔ Windows سازگار از Windows، یا از محیط سبک Linux استفاده کنید. پروژه‌های دلخواه را به اشتراک بگذارید، شبکه را جداگانه مدیریت کنید، با SSH به سرور خود وصل شوید و کارها را از تلفن ادامه دهید.

**رابط اصلی Cage:**

![Cage](www/public/workbench-2026-11.jpg)

## کاری که می توانید با Cage انجام دهید

- **در یک پنجره واحد کار کنید.** Claude Code را در ترمینال اجرا کنید، تغییرات ویرایشگر را بررسی کنید و از مرورگر Cage برای پیش نمایش خدمات و ورود به سیستم استفاده کنید.
- **در سیستم آشنا توسعه دهید.** سیستم مهمان macOS، Windows یا Linux را با ابزارهای همان پلتفرم انتخاب کنید، یا سیستم میزبان را برای ساخت محیطی مستقل و پاک شبیه‌سازی کنید. [انتخاب سیستم مهمان](docs/guide/guest-systems.md)
- **آنچه را که به اشتراک می گذارید کنترل کنید.** فقط فضای کاری پیش‌فرض و پوشه‌های پروژه به‌صراحت باز شده به اشتراک گذاشته می‌شوند. وضعیت ورود به سیستم Claude میزبان، نمایه‌های مرورگر و کلیدهای SSH به‌طور خودکار وارد نمی‌شوند.
- **وضعیت شبکه خود را ببینید.** مدیریت پروکسی شبکه داخلی به شما امکان می‌دهد اشتراک‌ها را وارد کنید، نقاط پایانی را انتخاب کنید و IP خروجی را بررسی کنید. محیط آفلاین می‌ماند تا زمانی که بررسی‌های حیاتی با موفقیت انجام شوند. اگر پراکسی شبکه در دسترس نباشد، هیچ بازگشتی به اتصال مستقیم وجود ندارد.
- **دفعه بعد ادامه بده** ابزارها، فایل‌ها و وضعیت ورود به سیستم روی دیسک مجازی باقی می‌مانند. وقتی پنجره را می‌بندید، کارهای پس‌زمینه می‌توانند همچنان اجرا شوند.
- **چند حساب Claude Code را مدیریت کنید.** سهمیه‌ها را ببینید، به نشست‌های جدید حساب اختصاص دهید و با اجازهٔ صریح، حساب را برای ادامهٔ تاریخچه تغییر دهید. [Claude Code](docs/guide/accounts.md)
- **از تلفن ادامه دهید.** Android یا iPhone را با کد QR یا کد جفت‌سازی به رایانه وصل کنید. همان گفت‌وگوهای Claude Code را ادامه دهید، پیام بفرستید، ابزارها را تأیید کنید و مرورگر راه دور یا پنجرهٔ Claude Desktop رایانه را به کار ببرید.
- **Cage را در فضای ابری اجرا کنید.** از طریق SSH به Cage مستقر روی VPS سازگار با Linux x86_64 وصل شوید؛ KVM لازم نیست. داده‌ها و خروجی شبکه روی سرور قرار دارند.
- **از شبکهٔ خصوصی وصل شوید.** در Cage وارد Tailscale شوید تا میزبان و محیط از هر دو جهت با SSH به هم وصل شوند. [Tailscale / SSH](docs/guide/tailscale.md)

## شروع سریع

این درخواست را برای عامل هوش مصنوعی خود (codex، dsh، workbuddy - احتمالاً Claude Code) ارسال کنید:

```markdown
کمکم کن نصب کنم https://github.com/Lee-xeo/Claude-Cage و پروکسی شبکه Cage را برای استفاده از نقطه پایانی xxx در نصب محلی Clash من پیکربندی کنید.
```

## دانلود و نصب کنید

Cage از macOS، Windows و Linux پشتیبانی می‌کند: Apple Silicon با macOS 14+، Windows 10 22H2 / 11 x64 و Debian 12 / Ubuntu 22.04 / 24.04 x86_64. Intel Mac، Windows ARM64 و Linux ARM64 فعلاً پشتیبانی نمی‌شوند. [macOS](docs/getting-started.md) · [Windows](docs/windows.md) · [Linux](docs/linux.md).

1. دانلود برنامه: [macOS / Windows / Linux](https://claude-cage.com/fa/#download).
2. macOS: بسته نصب را باز کرده و بکشید **Cage** به **برنامه های کاربردی**; سپس آن را از Applications اجرا کنید. Windows: [راهنمای نصب](docs/windows.md#安装并开始工作) · Linux: [راهنمای نصب](docs/linux.md#安装并开始工作).
3. دستورالعمل های درون برنامه را دنبال کنید تا پروکسی شبکه خود را پیکربندی کنید، محیط را راه اندازی کنید و یک پروژه را باز کنید.

بستهٔ نصب شامل فضای کاری و موتور پراکسی است. نیازی به نصب قبلی Go، Node.js، Docker یا Clash نیست. هنگام ساخت محیط، سیستم مهمان و محل ذخیره را انتخاب و تصویر رسمی و ابزارها را آماده کنید. Cage مجوز سیستم‌عامل، سرویس پراکسی یا حساب Claude ارائه نمی‌کند.

کد منبع و بسته‌های نصب عمومی هستند. یادداشت نسخه‌ها و فایل‌های بررسی صحت در [Releases](https://github.com/Lee-xeo/Claude-Cage/releases) قرار دارند.

قابلیت‌ها به ساخت مورد استفاده بستگی دارند؛ قابلیت‌ها و دامنهٔ بررسی نصب‌کنندهٔ عمومی در یادداشت انتشار آمده‌اند.

macOS و نسخه‌های امضاشدهٔ Windows از پراکسی انتخابی به‌روز می‌شوند. پیش‌نمایش بدون امضای Windows و Linux با نصب‌کننده به‌صورت دستی به‌روز می‌شوند. ابتدا کارها را ذخیره و متوقف کنید؛ داده‌های شخصی حفظ می‌شوند. [به‌روزرسانی و حذف](docs/guide/maintenance.md)

## در سه مرحله شروع کنید

### 1. یک پروکسی شبکه را انتخاب کنید

در صفحه اصلی، روی «پیکربندی شبکه و پروکسی» و سپس «افزودن منبع» کلیک کنید. یک اشتراک Clash/Mihomo YAML یا یک فایل YAML محلی را وارد کنید و یک نقطه پایانی را در زیر «Nodes» انتخاب کنید. اگر هیچ پروکسی شبکه در دسترس نیست، می‌توانید صراحتاً «دانلود این اشتراک از طریق شبکه میزبان» را انتخاب کنید یا یک فایل محلی را مستقیماً وارد کنید. بعداً، نقاط پایانی را در قسمت «تنظیمات → شبکه و پروکسی» مدیریت کنید.

### 2. محیط را شروع کنید

در ساخت نخست، سیستم مهمان یا تکثیر سیستم میزبان و محل ذخیره را انتخاب کنید. محیط Cage را راه‌اندازی کرده و مراحل نصب و مجوزها را دنبال کنید. پیشرفت و خطاها در جزئیات عملیات دیده می‌شوند. راه‌اندازی‌های بعدی همان محیط را به کار می‌برند.

### 3. یک پروژه را باز کنید و Claude را اجرا کنید

برای انتخاب یک پوشه در رایانه، روی «باز کردن پروژه» کلیک کنید، یک ترمینال جدید باز کنید و وارد شوید `claude`; سپس حساب خود را در مرورگر Cage مجاز کنید. همچنین می توانید بدون تهیه پروژه با فضای کاری پیش فرض شروع کنید. برای مشاهده کامل مراحل، نگاه کنید [شروع کردن](docs/getting-started.md).

## قبل از استفاده این محدودیت ها را درک کنید

- ویرایش‌ها و حذف‌ها در پروژه‌های مشترک مستقیماً روی فایل‌های رایانه شما تأثیر می‌گذارد. از کنترل نسخه یا نسخه پشتیبان استفاده کنید. پوشهٔ `.git` پروژه خواندنی و نوشتنی است؛ می‌توانید در Cage تغییرات را آماده و ثبت کنید و شاخه را تغییر دهید.
- همه پروژه ها یک محیط Cage شامل ابزارها، حساب ها و سیستم فایل لینوکس را به اشتراک می گذارند. پروژه ها از یکدیگر جدا نیستند.
- خروجی عمومی به‌طور پیش‌فرض IPv4 TCP را پشتیبانی می‌کند و DNS از همان پراکسی می‌گذرد. UDP، IPv6 و ICMP عمومی و اتصال مستقیم به میزبان یا شبکهٔ محلی محدودند. با فعال‌سازی صریح Tailscale، دسترسی به دستگاه‌ها تابع مجوزهای شبکهٔ خصوصی است.
- بررسی های شبکه پیکربندی فعلی و شواهد نمونه برداری شده را نشان می دهد. آن‌ها نمی‌توانند تضمین کنند که یک حساب از محدودیت‌های پلتفرم جلوگیری می‌کند، و جایگزین ممیزی امنیتی کامل مجازی‌سازی نمی‌شوند.

برای جزئیات، نگاه کنید [مرزهای شبکه و داده](docs/getting-started.md#网络与数据边界).

## اسناد و مشارکت

برای نمایه کامل، نگاه کنید [نمایه اسناد و مدارک](docs/README.md).

| وظیفه شماست | از اینجا شروع کنید |
|---|---|
| macOS / Windows / Debian / Ubuntu 22.04 | [انتخاب سیستم مهمان](https://claude-cage.com/docs/guest-systems/) |
| یک پروکسی شبکه را نصب کنید، پیکربندی کنید و اولین پروژه خود را باز کنید | [شروع کردن](docs/getting-started.md) |
| عیب‌یابی دانلود، راه‌اندازی، اتصال و ورود به سیستم | [سوالات متداول و بازیابی](docs/troubleshooting.md) |
| کنترل های صفحه کلید، ظاهر، ماندگاری داده ها و رفتار خروج را بیاموزید | [استفاده روزانه](docs/getting-started.md#日常使用) |
| یک مشکل قابل تکرار را گزارش کنید | [GitHub Issues](https://github.com/Lee-xeo/Claude-Cage/issues) |

این پروژه به سرعت در حال پیشرفت است. برای به اشتراک گذاشتن بازخورد به گروه جامعه بپیوندید. از نگهبانان جامعه استقبال می شود.

<p align="center">
  <a href="assets/community-qq.png"><img src="assets/community-qq.png" alt="QQ #2 · 1127457578" width="240" /></a>
  <br />QQ #2 · 1127457578
</p>

## مبدا و مجوز

Cage مشتق شده از [agent-vm](https://github.com/sylvinus/agent-vm).

مجوز اصلی پروژه می باشد [MIT](LICENSE). اجزای همراه مجوز مربوطه خود را حفظ می کنند. Mihomo به عنوان یک فایل اجرایی جداگانه تحت GPL-3.0 توزیع شده است. همچنین ببینید [NOTICE](NOTICE) · [مجوزهای شخص ثالث](THIRD_PARTY_LICENSES) · [noVNC](ui/vendor/novnc/LICENSE.txt)。

Cage یک پروژه اجتماعی مستقل است و به آنتروپیک وابسته نیست.

## مراجع و قدردانی

با تشکر از پروژه های زیر برای قابلیت های اساسی و مراجع طراحی:

- [agent-vm](https://github.com/sylvinus/agent-vm): پایه بالادست Cage.
- [Lima](https://github.com/lima-vm/lima): مدیریت ماشین مجازی لینوکس.
- [gvisor-tap-vsock](https://github.com/containers/gvisor-tap-vsock): شبکه ماشین مجازی.
- [sshocker](https://github.com/lima-vm/sshocker) و [SFTP](https://github.com/pkg/sftp): به اشتراک گذاری فایل بین ماشین مجازی و میزبان.
- [noVNC](https://github.com/novnc/noVNC): نمایش از راه دور و تعامل در مرورگر.
- [Mihomo](https://github.com/MetaCubeX/mihomo): هسته پروکسی شبکه داخلی.
- [code-server](https://github.com/coder/code-server): ویرایشگر کد مبتنی بر مرورگر.
- [AIO Sandbox](https://github.com/agent-infra/sandbox): مراجع برای طراحی محیط کار و تعامل.


<a href="https://www.star-history.com/?repos=Lee-xeo%2FClaude-Cage&amp;type=date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date&amp;theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" />
    <img alt="Claude Cage GitHub Star History" src="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" width="800" />
  </picture>
</a>
