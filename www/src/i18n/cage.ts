import { languages, type CageLocale } from './languages';
export { languages, type CageLocale } from './languages';
export const locales = languages.map(language => language.code);
const readmes = import.meta.glob<string>('../../../README*.md', { query: '?raw', import: 'default', eager: true });

function plain(text: string) {
  return text.replace(/<img\b[^>]*>/g, '').replace(/!\[[^\]]*\]\([^)]*\)/g, '')
    .replace(/\[([^\]]+)\]\([^)]*\)/g, '$1').replace(/<[^>]+>/g, '')
    .replace(/\*\*|`/g, '').trim();
}
function paragraphs(text: string) {
  return text.replace(/```[\s\S]*?```/g, '').split(/\n\s*\n/)
    .filter(block => !/^\s*(?:[-#|]|\d+\.)/.test(block))
    .map(plain).filter(Boolean);
}
function bullets(text: string) {
  return [...text.matchAll(/^- (.+)$/gm)].map(match => match[1]);
}
function pair(text: string): [string, string] {
  const bold = text.match(/^\*\*(.+?)\*\*\s*(.*)$/);
  return bold ? [plain(bold[1]), plain(bold[2])] : ['', plain(text)];
}
function readCopy(language: typeof languages[number]) {
  const file = language.code === 'zh' ? 'README.md' : `README.${language.code}.md`;
  const source = readmes[`../../../${file}`];
  if (!source) throw new Error(`Missing website language source: ${file}`);
  const allHeadings = [...source.matchAll(/^## (.+)$/gm)];
  const headings = allHeadings.filter(heading => !/^Star (?:History|增长曲线|成長曲線)$/.test(heading[1]));
  if (headings.length !== 9) throw new Error(`${file}: expected the nine shared README sections; update the website mapping when the document structure changes`);
  const sections = headings.map(heading => ({
    title: heading[1], body: source.slice(heading.index! + heading[0].length, allHeadings[allHeadings.indexOf(heading) + 1]?.index ?? source.length),
  }));
  const [reason, uses, quick, installation, steps, boundaries, contribution, license] = sections;
  const reasonParagraphs = paragraphs(reason.body);
  const installParagraphs = paragraphs(installation.body);
  const captions = [...reason.body.matchAll(/^\*\*(.+?)\*\*\s*$/gm)].map(match => plain(match[1]));
  const stepHeadings = [...steps.body.matchAll(/^### \d+\. (.+)$/gm)];
  if (stepHeadings.length !== 3 || bullets(boundaries.body).length !== 4) {
    throw new Error(`${file}: README setup or boundary structure has changed`);
  }
  const setup = stepHeadings.map((heading, i): [string,string] => [heading[1], plain(steps.body.slice(heading.index! + heading[0].length, stepHeadings[i+1]?.index ?? steps.body.length))]);
  const docs = [...contribution.body.matchAll(/^\|\s*([^|]+)\s*\|\s*\[([^\]]+)\]\(([^)]+)\)\s*\|/gm)].map(match => ({
    title: plain(match[1]), label: match[2], href: match[3].startsWith('https://') ? match[3] : `https://github.com/Lee-xeo/Claude-Cage/blob/main/${match[3]}`,
  }));
  const slogan = source.slice(0, headings[0].index).match(/^\*\*(.+?)\*\*$/m)?.[1] || 'Claude Cage';
  const isChinese = language.code === 'zh' || language.code === 'zh-TW';
  return {
    ...language, title: `Claude Cage — ${plain(slogan)}`, description: reasonParagraphs[0],
    home: 'Claude Cage', start: quick.title, boundary: boundaries.title,
    tagline: plain(slogan),
    intro: reasonParagraphs[0], requirements: 'macOS · Windows · Linux · Android · iOS', artAlt: 'Claude Cage',
    featureTitle: uses.title, features: bullets(uses.body).map(pair),
    workbench: captions[0] || 'Claude Cage',
    quickTitle: quick.title, quickText: paragraphs(quick.body)[0],
    prompt: quick.body.match(/```[^\n]*\n([\s\S]*?)```/)?.[1]?.trim() || '',
    quickNote: paragraphs(quick.body).slice(1).join('\n'),
    stepsTitle: steps.title, steps: setup,
    boundaryTitle: boundaries.title, boundaries: bullets(boundaries.body).map(pair),
    downloadTitle: installation.title, downloadText: installParagraphs[1],
    downloadNote: installParagraphs[0], updateNote: installParagraphs.at(-1) || '',
    communityTitle: contribution.title, communityText: paragraphs(contribution.body.replace(/<p\b[^>]*>[\s\S]*?<\/p>/g, '')).at(-1) || '',
    group: isChinese ? 'QQ 二群 · 1127457578' : 'QQ #2 · 1127457578', docsTitle: contribution.title, docs,
    licenseTitle: license.title, licenseText: paragraphs(license.body).slice(0,-1).join('\n\n'),
    independent: paragraphs(license.body).at(-1) || '',
  };
}
export const cageCopy = Object.fromEntries(languages.map(language => [language.code, readCopy(language)])) as Record<CageLocale, ReturnType<typeof readCopy>>;
