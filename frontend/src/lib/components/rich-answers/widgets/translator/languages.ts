export interface Language {
	code: string;
	name: string;
	rtl?: boolean;
}

export const AUTO: Language = { code: 'auto', name: 'detect language' };

const RTL = new Set(['ar', 'fa', 'he', 'ur']);

const LANG_LIST: [string, string][] = [
	['af', 'afrikaans'],
	['sq', 'albanian'],
	['am', 'amharic'],
	['ar', 'arabic'],
	['hy', 'armenian'],
	['az', 'azerbaijani'],
	['be', 'belarusian'],
	['bn', 'bengali'],
	['bs', 'bosnian'],
	['bg', 'bulgarian'],
	['zh-CN', 'chinese (simplified)'],
	['zh-TW', 'chinese (traditional)'],
	['hr', 'croatian'],
	['cz', 'czech'],
	['dk', 'danish'],
	['nl', 'dutch'],
	['en', 'english'],
	['et', 'estonian'],
	['fi', 'finnish'],
	['fr', 'french'],
	['ka', 'georgian'],
	['de', 'german'],
	['el', 'greek'],
	['gu', 'gujarati'],
	['he', 'hebrew'],
	['hi', 'hindi'],
	['hu', 'hungarian'],
	['is', 'icelandic'],
	['id', 'indonesian'],
	['it', 'italian'],
	['ja', 'japanese'],
	['kn', 'kannada'],
	['kk', 'kazakh'],
	['km', 'khmer'],
	['ko', 'korean'],
	['ky', 'kyrgyz'],
	['lo', 'lao'],
	['lv', 'latvian'],
	['lt', 'lithuanian'],
	['mk', 'macedonian'],
	['ms', 'malay'],
	['ml', 'malayalam'],
	['mt', 'maltese'],
	['mr', 'marathi'],
	['mn', 'mongolian'],
	['my', 'burmese'],
	['ne', 'nepali'],
	['no', 'norwegian'],
	['ps', 'pashto'],
	['fa', 'persian'],
	['pl', 'polish'],
	['pt', 'portuguese'],
	['pa', 'punjabi'],
	['ro', 'romanian'],
	['ru', 'russian'],
	['sr', 'serbian'],
	['si', 'sinhala'],
	['sk', 'slovak'],
	['sl', 'slovenian'],
	['es', 'spanish'],
	['sw', 'swahili'],
	['sv', 'swedish'],
	['tl', 'filipino'],
	['ta', 'tamil'],
	['te', 'telugu'],
	['th', 'thai'],
	['tr', 'turkish'],
	['uk', 'ukrainian'],
	['ur', 'urdu'],
	['uz', 'uzbek'],
	['vi', 'vietnamese'],
];

export const LANGUAGES: Language[] = LANG_LIST.map(([code, name]) =>
	RTL.has(code.slice(0, 2)) ? { code, name, rtl: true } : { code, name }
).sort((a, b) => a.name.localeCompare(b.name));

const byCode = new Map(LANGUAGES.map((l) => [l.code, l]));

export function langByCode(code: string): Language | null {
	return byCode.get(code) ?? null;
}

const FLAG_OVERRIDES: Record<string, string> = {
	af: 'za',
	sq: 'al',
	am: 'et',
	ar: 'sa',
	hy: 'am',
	be: 'by',
	bn: 'bd',
	el: 'gr',
	en: 'gb',
	et: 'ee',
	fa: 'ir',
	gu: 'in',
	he: 'il',
	hi: 'in',
	ja: 'jp',
	ka: 'ge',
	kk: 'kz',
	km: 'kh',
	kn: 'in',
	ko: 'kr',
	ky: 'kg',
	lo: 'la',
	ml: 'in',
	mr: 'in',
	ms: 'my',
	my: 'mm',
	ne: 'np',
	pa: 'in',
	si: 'lk',
	sl: 'si',
	sr: 'rs',
	sv: 'se',
	sw: 'tz',
	ta: 'in',
	te: 'in',
	tl: 'ph',
	uk: 'ua',
	ur: 'pk',
	vi: 'vn',
	'zh-CN': 'cn',
	'zh-TW': 'tw',
};

export function matchesLanguage(lang: Language, query: string): boolean {
	return lang.name.toLowerCase().includes(query) || lang.code.toLowerCase().includes(query);
}

export function flagUrl(code: string): string | null {
	if (code === 'auto') return null;
	const country = FLAG_OVERRIDES[code] ?? code.slice(0, 2).toLowerCase();
	const src = `https://flagcdn.com/w40/${country}.png`;
	return `https://external-content.duckduckgo.com/iu/?u=${src}`;
}

export function browserLanguage(): string {
	for (const tag of navigator.languages?.length ? navigator.languages : [navigator.language]) {
		if (!tag) continue;
		const base = tag.toLowerCase().split('-')[0];
		if (byCode.has(base)) return base;
	}
	return 'en';
}
