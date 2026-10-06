const source =
	'https://cdn.jsdelivr.net/gh/laylavish/uBlockOrigin-HUGE-AI-Blocklist@main/list_uBlacklist.txt';

const storeKey = 'clasp.blocklist_v2';
const storeTs = 'clasp.blocklist_ts_v2';
const maxAge = 86_400_000;

const triggers = [
	'ai generated',
	'ai-generated',
	'ai art',
	'stable diffusion',
	'midjourney',
	'dall-e',
	'dalle',
	'generated with ai',
	'images using ai',
	'generated from ai',
	'ai image generator',
	'ai image',
	'ai digital art',
	'image creator',
	'canva ai',
];

const triggerWords = triggers.map((t) => t.split(/[\s-]+/).filter(Boolean));

type Rule = { kind: 're' | 'glob' | 'plain'; rx: RegExp | string };

let rules: Rule[] | undefined;

function compileRule(raw: string): Rule | null {
	const s = raw.trim();
	if (!s) return null;

	if (s.startsWith('/') && s.endsWith('/')) {
		try {
			return { kind: 're', rx: new RegExp(s.slice(1, -1), 'i') };
		} catch {
			return null;
		}
	}

	if (s.includes('*')) {
		const escaped = s
			.replace(/\*:\/\//g, 'https?://')
			.replace(/\*/g, '.*')
			.replace(/\./g, '\\.');
		try {
			return { kind: 'glob', rx: new RegExp(escaped, 'i') };
		} catch {
			return null;
		}
	}

	return { kind: 'plain', rx: s.toLowerCase() };
}

function parse(raw: string): Rule[] {
	const out: Rule[] = [];
	for (const line of raw.split('\n')) {
		if (line.startsWith('#')) continue;
		const r = compileRule(line);
		if (r) out.push(r);
	}
	return out;
}

function loadFromStore(): Rule[] | null {
	try {
		const data = sessionStorage.getItem(storeKey);
		const ts = sessionStorage.getItem(storeTs);
		if (data && ts && Date.now() - Number(ts) < maxAge) {
			return parse(data);
		}
	} catch {}
	return null;
}

function saveToStore(raw: string) {
	try {
		sessionStorage.setItem(storeKey, raw);
		sessionStorage.setItem(storeTs, String(Date.now()));
	} catch {}
}

async function fetchRules(): Promise<Rule[]> {
	const fresh = loadFromStore();
	if (fresh) return fresh;

	try {
		const res = await fetch(source);
		if (!res.ok) throw 0;
		const text = await res.text();
		saveToStore(text);
		return parse(text);
	} catch {
		return loadFromStore() ?? [];
	}
}

function hasAITitle(title: string): boolean {
	const words = title.toLowerCase().replaceAll('-', ' ').split(/\s+/).filter(Boolean);
	for (const kw of triggerWords) {
		if (kw.length === 1) {
			if (words.includes(kw[0])) return true;
			continue;
		}
		for (let i = 0; i <= words.length - kw.length; i++) {
			const window = words.slice(i, i + kw.length);
			if (kw.every((w) => window.includes(w))) return true;
		}
	}
	return false;
}

function matchesRules(url: string, hostname: string): boolean {
	const lowerUrl = url.toLowerCase();
	const lowerHost = hostname.toLowerCase();
	for (const { kind, rx } of rules ?? []) {
		if (kind === 'plain') {
			const t = rx as string;
			if (lowerUrl.includes(t) || lowerHost.includes(t)) return true;
		} else {
			const re = rx as RegExp;
			if (re.test(url) || re.test(hostname)) return true;
		}
	}
	return false;
}

export async function initBlocklist(): Promise<void> {
	rules = await fetchRules();
}

export function shouldBlockImage(img: {
	url?: string;
	title?: string;
	hostname?: string;
}): boolean {
	if (!rules) return false;

	const hostname = img.hostname || '';
	const url = img.url || '';

	if (hostname.endsWith('.ai')) return true;
	if (hasAITitle(img.title || '')) return true;
	if (matchesRules(url, hostname)) return true;

	return false;
}
