const TRIGGERS = new Set([
	// latin
	'translate',
	'translation',
	'translator',
	'traduire',
	'traduction',
	'traducteur',
	'tradurre',
	'traduzione',
	'traduttore',
	'traduzir',
	'tradução',
	'übersetzen',
	'übersetzung',
	'übersetzer',
	'vertalen',
	'vertaling',
	'vertaal',
	'översätta',
	'oversættelse',
	'oversette',
	'oversättning',
	'käännös',
	'käänttää',
	'prevesti',
	'překlad',
	'přeložit',
	'sumon',
	'terjemhely',
	'terjemény',
	'çeviri',
	'tercüme',
	'thdịch',
	'dịch',
	'preklad',
	'tłumacz',
	'şahal',
	// cyrillic
	'перевод',
	'перевести',
	'переведи',
	'переводчик',
	'переклад',
	'перекласти',
	'превод',
	'преводи',
	'преводач',
	'перекладання',
	'тәрҗемә',
	'көчмә',
	'тарҷума',
	'тарчума',
	// greek, arabic, hebrew, armenian, georgian
	'μετάφραση',
	'مترجم',
	'ترجمة',
	'תרגום',
	'թարգմանություն',
	'თარგმანი',
	// cjk
	'翻译',
	'翻譯',
	'通译',
	'翻訳',
	'번역',
	'翻訳者',
	'번역자',
]);

export function parseTranslateQuery(query: string): boolean {
	const q = query.trim();
	if (!q) return false;
	const bare = q
		.replace(/^(?:google|bing|yandex|deepl|kagi)\s+/i, '')
		.replace(/[^\p{L}\p{N}-]/gu, '')
		.toLowerCase();
	return TRIGGERS.has(bare);
}
