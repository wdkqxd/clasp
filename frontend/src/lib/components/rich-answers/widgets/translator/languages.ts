import iso from '@cospired/i18n-iso-languages';
import en from '@cospired/i18n-iso-languages/langs/en.json';

iso.registerLocale(en);

export interface Language {
	code: string;
	name: string;
	rtl?: boolean;
}

export const AUTO: Language = { code: 'auto', name: 'detect language' };

const RTL = new Set(['ar', 'fa', 'he', 'ur', 'ps', 'sd', 'ug', 'yi']);

export const LANGUAGES: Language[] = (() => {
	const written: Record<string, string> = {
		'zh-CN': 'chinese (simplified)',
		'zh-TW': 'chinese (traditional)',
	};
	const shadowed = new Set(Object.keys(written).map((c) => c.split('-')[0]));
	const codes = [
		...Object.keys(iso.getAlpha2Codes()).filter((c) => !shadowed.has(c)),
		...Object.keys(written),
	].sort();

	return codes.map((code) => {
		const name = (written[code] ?? iso.getName(code, 'en') ?? code).toLowerCase();
		return RTL.has(code.slice(0, 2)) ? { code, name, rtl: true } : { code, name };
	});
})();

const byCode = new Map(LANGUAGES.map((l) => [l.code, l]));

export function langByCode(code: string): Language | null {
	return byCode.get(code) ?? null;
}

const FLAG_OVERRIDES: Record<string, string> = {
	af: 'za',
	am: 'et',
	ar: 'sa',
	bn: 'bd',
	ca: 'es',
	el: 'gr',
	en: 'gb',
	et: 'ee',
	fa: 'ir',
	ga: 'ie',
	gl: 'es',
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
	lb: 'lu',
	lo: 'la',
	ml: 'in',
	mr: 'in',
	ms: 'my',
	my: 'mm',
	ne: 'np',
	pa: 'in',
	sr: 'rs',
	si: 'lk',
	sl: 'si',
	sv: 'se',
	sw: 'tz',
	ta: 'in',
	te: 'in',
	tl: 'ph',
	uk: 'ua',
	ur: 'pk',
	vi: 'vn',
	cy: 'gb',
	'zh-CN': 'cn',
	'zh-TW': 'tw',
};

export function matchesLanguage(lang: Language, query: string): boolean {
	return lang.name.toLowerCase().includes(query) || lang.code.toLowerCase().includes(query);
}

export function flagUrl(code: string): string | null {
	if (code === 'auto') return null;
	const country = FLAG_OVERRIDES[code] ?? code.slice(0, 2).toLowerCase();
	const src = `https://flagcdn.com/w80/${country}.png`;
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
