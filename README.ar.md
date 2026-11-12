<div align="center">

<h1>Claude Cage</h1>

**مساحة تطوير معزولة**

<img src="assets/claude-cage-character.png" alt="Claude Cage" width="320" />

[![Release](https://img.shields.io/github/v/release/Lee-xeo/Claude-Cage?include_prereleases)](https://github.com/Lee-xeo/Claude-Cage/releases) [![CI](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/Lee-xeo/Claude-Cage/actions/workflows/ci.yml) [![MIT](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

<a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/daily?language=Go" alt="Trendshift — Go Repository of the Day" width="250" height="55" /></a> <a href="https://trendshift.io/repositories/339051#achievements"><img src="https://trendshift.io/api/badge/trendshift/repositories/339051/weekly?language=Go" alt="Trendshift — Go Repository of the Week" width="250" height="55" /></a>

**[موقع المشروع](https://claude-cage.com/) · [تنزيل التطبيق](https://claude-cage.com/ar/#download) · [البدء](docs/getting-started.md) · [الأسئلة الشائعة](docs/troubleshooting.md) · [دليل الاستخدام (بالصينية)](https://claude-cage.com/docs/)**

[简体中文](README.md) · [繁體中文](README.zh-TW.md) · [Русский](README.ru.md) · [العربية](README.ar.md) · [فارسی](README.fa.md) · [한국어](README.ko.md) · [Español](README.es.md) · [မြန်မာဘာသာ](README.my.md) · [Беларуская](README.be.md) · [پښتو](README.ps.md) · [دری](README.prs.md) · [አማርኛ](README.am.md) · [Soomaali](README.so.md) · [ትግርኛ](README.ti.md) · [Français](README.fr.md)

</div>

## مساحة تطوير معزولة

يجمع Claude Cage الطرفية والمتصفح والمحرر وتطبيق Claude Desktop الرسمي في مساحة عمل واحدة معزولة. استخدم جهاز macOS افتراضيًا على Mac، أو Windows على حاسوب متوافق، أو بيئة Linux خفيفة. شارك المشاريع التي تختارها، وأدر الشبكة بشكل مستقل، واتصل بخادمك عبر SSH، وتابع مهامك من الهاتف.

**الواجهة الرئيسية للقفص:**

![Cage](www/public/workbench-2026-11.jpg)

## ما يمكنك فعله مع كيج

- **العمل في نافذة واحدة.** قم بتشغيل Claude Code في الجهاز، وراجع التغييرات في المحرر، واستخدم متصفح Cage لمعاينة الخدمات وتسجيل الدخول.
- **طوّر على نظامك المألوف.** اختر نظامًا ضيفًا macOS أو Windows أو Linux بأدواته، أو استنسخ نظام المضيف لإنشاء بيئة مستقلة ونظيفة. [اختيار النظام الضيف](docs/guide/guest-systems.md)
- **السيطرة على ما تشاركه.** تتم مشاركة مساحة العمل الافتراضية ومجلدات المشروع المفتوحة بشكل صريح فقط. لا يتم استيراد حالة تسجيل دخول Claude للمضيف وملفات تعريف المتصفح ومفاتيح SSH تلقائيًا.
- **انظر حالة الشبكة الخاصة بك.** تتيح لك إدارة وكيل الشبكة المضمنة استيراد الاشتراكات وتحديد نقاط النهاية والتحقق من عنوان IP الخاص بالخروج. تظل البيئة غير متصلة بالإنترنت حتى تمر عمليات التحقق المهمة. إذا كان وكيل الشبكة غير متوفر، فلا يمكن الرجوع إلى الاتصال المباشر.
- **استمر في المرة القادمة.** تستمر الأدوات والملفات وحالة تسجيل الدخول على القرص الظاهري. يمكن أن تستمر مهام الخلفية في العمل عند إغلاق النافذة.
- **إدارة عدة حسابات Claude Code.** اعرض الحصص وخصص الحسابات للجلسات الجديدة، وبدّل الحساب لمتابعة السجل بعد الحصول على إذن صريح. [Claude Code](docs/guide/accounts.md)
- **تابع من هاتفك.** اربط Android أو iPhone بالكمبيوتر باستخدام رمز QR أو رمز الاقتران. تابع محادثات Claude Code نفسها، وأرسل الرسائل ووافق على الأدوات، واستخدم المتصفح البعيد أو نافذة Claude Desktop على الكمبيوتر.
- **تشغيل Cage في السحابة.** اتصل عبر SSH ببيئة Cage المنشورة على VPS متوافق يعمل بنظام Linux x86_64، دون الحاجة إلى KVM. تبقى البيانات ومخرج الشبكة على الخادم.
- **اتصل عبر شبكة خاصة.** سجّل الدخول إلى Tailscale داخل Cage للاتصال عبر SSH بين المضيف والبيئة في الاتجاهين. [Tailscale / SSH](docs/guide/tailscale.md)

## بداية سريعة

أرسل هذا الطلب إلى وكيل الذكاء الاصطناعي لديك (codex، dsh، Workbuddy - ربما ليس Claude Code 😁):

```markdown
ساعدني في التثبيت https://github.com/Lee-xeo/Claude-Cage وقم بتكوين وكيل شبكة Cage لاستخدام نقطة النهاية xxx في تثبيت Clash المحلي الخاص بي.
```

## التنزيل والتثبيت

يدعم Cage أنظمة macOS وWindows وLinux: أجهزة Apple Silicon مع macOS 14+، وWindows 10 22H2 / 11 x64، وDebian 12 / Ubuntu 22.04 / 24.04 x86_64. Intel Mac أو Windows ARM64 أو Linux ARM64. [macOS](docs/getting-started.md) · [Windows](docs/windows.md) · [Linux](docs/linux.md).

1. تنزيل التطبيق: [macOS / Windows / Linux](https://claude-cage.com/ar/#download).
2. macOS: افتح حزمة التثبيت واسحب **Cage** في **التطبيقات**; ثم قم بتشغيله من التطبيقات. Windows: [دليل التثبيت](docs/windows.md#安装并开始工作) · Linux: [دليل التثبيت](docs/linux.md#安装并开始工作).
3. اتبع الإرشادات الموجودة في التطبيق لتكوين وكيل الشبكة الخاص بك، وبدء تشغيل البيئة، وفتح المشروع.

تتضمن حزمة التثبيت مساحة العمل ومحرك الوكيل. لا تحتاج إلى تثبيت Go أو Node.js أو Docker أو Clash مسبقًا. عند إنشاء البيئة اختر النظام الضيف ومكان التخزين، ثم جهّز الصورة الرسمية وأدوات التطوير. لا يوفر Cage ترخيص نظام تشغيل أو خدمة وكيل أو حساب Claude.

الكود المصدري وحزم التثبيت متاحة للعامة. تجد ملاحظات الإصدارات وملفات التحقق في [Releases](https://github.com/Lee-xeo/Claude-Cage/releases).

تعتمد الميزات على البناء المستخدم؛ راجع ملاحظات الإصدار لقدرات المثبّت العام ونطاق التحقق منه.

تُحدَّث macOS وإصدارات Windows الموقعة عبر الوكيل المحدد. تُحدَّث معاينة Windows غير الموقعة وLinux يدويًا باستخدام المثبّت. احفظ المهام وأوقفها قبل التحديث؛ تبقى البيانات الشخصية محفوظة. [التحديث وإلغاء التثبيت](docs/guide/maintenance.md)

## ابدأ في ثلاث خطوات

### 1. حدد وكيل الشبكة

في الصفحة الرئيسية، انقر فوق "تكوين الشبكة والوكيل"، ثم "إضافة مصدر". قم باستيراد اشتراك Clash/Mihomo YAML أو ملف YAML محلي، وحدد نقطة نهاية ضمن "العقد". في حالة عدم توفر وكيل شبكة، يمكنك تحديد "تنزيل هذا الاشتراك من خلال الشبكة المضيفة" بشكل صريح، أو استيراد ملف محلي مباشرةً. لاحقًا، قم بإدارة نقاط النهاية ضمن "الإعدادات → الشبكة والوكيل".

### 2. ابدأ البيئة

عند الإنشاء الأول، اختر النظام الضيف أو استنساخ نظام المضيف وموقع التخزين. شغّل بيئة Cage واتبع خطوات التثبيت والأذونات. يعرض قسم تفاصيل العملية التقدم والأخطاء. تستخدم عمليات التشغيل التالية البيئة نفسها.

### 3. افتح مشروعًا وقم بتشغيل كلود

انقر فوق "فتح المشروع" لتحديد مجلد على جهاز الكمبيوتر الخاص بك، وافتح محطة جديدة، ثم أدخل `claude`; ثم قم بتفويض حسابك الخاص في متصفح Cage. يمكنك أيضًا البدء بمساحة العمل الافتراضية دون إعداد مشروع. للحصول على الخطوات الكاملة، راجع [البدء](docs/getting-started.md).

## افهم هذه الحدود قبل الاستخدام

- تؤثر عمليات التحرير والحذف في المشاريع المشتركة بشكل مباشر على الملفات الموجودة على جهاز الكمبيوتر الخاص بك. استخدم التحكم في الإصدار أو النسخ الاحتياطية. دليل `.git` في المشروع قابل للقراءة والكتابة، ويمكن تجهيز التغييرات وإنشاء الالتزامات وتبديل الفروع داخل Cage.
- تشترك جميع المشاريع في بيئة Cage واحدة، بما في ذلك أدواتها وحساباتها ونظام ملفات macOS / Windows / Debian / Ubuntu 22.04. المشاريع ليست معزولة عن بعضها البعض.
- يدعم الاتصال العام الصادر افتراضيًا IPv4 TCP ويمر DNS عبر الوكيل نفسه. تُقيَّد UDP وIPv6 وICMP العامة والاتصالات المباشرة بالمضيف والشبكة المحلية. عند تفعيل Tailscale صراحةً، يخضع الوصول إلى الأجهزة لأذونات الشبكة الخاصة.
- تُظهر عمليات فحص الشبكة التكوين الحالي والأدلة التي تم أخذ عينات منها. ولا يمكنها ضمان أن الحساب سوف يتجنب قيود النظام الأساسي، ولا تحل محل تدقيق أمان المحاكاة الافتراضية الكامل.

لمزيد من التفاصيل، انظر [حدود الشبكة والبيانات](docs/getting-started.md#网络与数据边界).

## التوثيق والمساهمات

للحصول على الفهرس الكامل، انظر [فهرس التوثيق](docs/README.md).

| مهمتك | ابدأ هنا |
|---|---|
| macOS / Windows / Debian / Ubuntu 22.04 | [اختيار النظام الضيف](https://claude-cage.com/docs/guest-systems/) |
| قم بتثبيت وكيل الشبكة وتكوينه وافتح مشروعك الأول | [البدء](docs/getting-started.md) |
| استكشاف أخطاء التنزيلات وبدء التشغيل والاتصال وتسجيل الدخول وإصلاحها | [الأسئلة الشائعة والاسترداد](docs/troubleshooting.md) |
| تعرف على عناصر التحكم في لوحة المفاتيح والمظهر واستمرارية البيانات وسلوك الخروج | [الاستخدام اليومي](docs/getting-started.md#日常使用) |
| الإبلاغ عن مشكلة قابلة للتكرار | [GitHub Issues](https://github.com/Lee-xeo/Claude-Cage/issues) |

المشروع يتطور بسرعة. انضم إلى مجموعة المجتمع لمشاركة التعليقات؛ مشرفو المجتمع مرحب بهم.

<p align="center">
  <a href="assets/community-qq.png"><img src="assets/community-qq.png" alt="QQ #2 · 1127457578" width="240" /></a>
  <br />QQ #2 · 1127457578
</p>

## الأصول والترخيص

القفص مشتق من [agent-vm](https://github.com/sylvinus/agent-vm).

الترخيص الرئيسي للمشروع هو [MIT](LICENSE). تحتفظ المكونات المجمعة بالتراخيص الخاصة بها. يتم توزيع Mihomo كملف قابل للتنفيذ منفصل تحت GPL-3.0. انظر أيضا [NOTICE](NOTICE) · [تراخيص الطرف الثالث](THIRD_PARTY_LICENSES) · [noVNC](ui/vendor/novnc/LICENSE.txt)。

Cage هو مشروع مجتمعي مستقل ولا ينتمي إلى Anthropic.

## المراجع والاعترافات

بفضل المشاريع التالية للقدرات التأسيسية ومراجع التصميم:

- [agent-vm](https://github.com/sylvinus/agent-vm): الأساس المنبع للقفص.
- [Lima](https://github.com/lima-vm/lima): إدارة الجهاز الظاهري لينكس.
- [gvisor-tap-vsock](https://github.com/containers/gvisor-tap-vsock): شبكات الآلة الافتراضية.
- [sshocker](https://github.com/lima-vm/sshocker) و [SFTP](https://github.com/pkg/sftp): مشاركة الملفات بين الجهاز الظاهري والمضيف.
- [noVNC](https://github.com/novnc/noVNC): العرض عن بعد والتفاعل في المتصفح.
- [Mihomo](https://github.com/MetaCubeX/mihomo): جوهر وكيل الشبكة المدمج.
- [code-server](https://github.com/coder/code-server): محرر التعليمات البرمجية القائم على المتصفح.
- [AIO Sandbox](https://github.com/agent-infra/sandbox): مراجع لبيئة العمل وتصميم التفاعل.


<a href="https://www.star-history.com/?repos=Lee-xeo%2FClaude-Cage&amp;type=date">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date&amp;theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" />
    <img alt="Claude Cage GitHub Star History" src="https://api.star-history.com/svg?repos=Lee-xeo/Claude-Cage&amp;type=Date" width="800" />
  </picture>
</a>
