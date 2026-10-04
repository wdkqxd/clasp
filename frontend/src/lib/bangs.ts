let cache: Record<string, string> | null = null;
let loading: Promise<Record<string, string>> | null = null;

async function load(): Promise<Record<string, string>> {
	if (cache) return cache;
	if (!loading) {
		loading = fetch('/api/bangs')
			.then((r) => (r.ok ? r.json() : null))
			.then((m) => {
				if (m && typeof m === 'object' && Object.keys(m).length > 0) {
					cache = m as Record<string, string>;
				}
				return cache ?? {};
			})
			.catch(() => cache ?? {})
			.finally(() => (loading = null));
	}
	return loading;
}

void load();

export async function checkBang(q: string): Promise<string | undefined> {
	const parsed = parseBang(q);
	if (!parsed) return;
	const url = (await load())[parsed.trigger];
	if (!url) return;
	return url.replace('%s', encodeURIComponent(parsed.rest));
}

export function parseBang(q: string): { trigger: string; rest: string } | null {
	const m = /^!(\S+)\s*(.*)$/.exec(q.trim());
	if (!m) return null;
	return { trigger: m[1].toLowerCase(), rest: m[2] };
}

export type BangMatch = {
	trigger: string;
	name: string;
	url: string;
	host: string;
};

export function bangIcon(host: string): string {
	return `https://icons.duckduckgo.com/ip3/${host}.ico`;
}

type BangEntry = {
	url: string;
	host: string;
	name: string;
	triggers: string[];
};

let index: BangEntry[] | null = null;
let triggerIndex: Map<string, BangEntry> | null = null;
let indexing: Promise<BangEntry[]> | null = null;

function hostOf(url: string): string {
	try {
		return new URL(url).hostname.replace(/^www\./, '');
	} catch {
		return '';
	}
}

function domainOf(host: string): string {
	const parts = host.split('.').filter(Boolean);
	if (parts.length < 2) return host;
	const tld = parts[parts.length - 1];
	return tld.length === 2 && parts.length > 2
		? parts.slice(-3).join('.')
		: parts.slice(-2).join('.');
}

function nameOf(host: string): string {
	const label = domainOf(host).split('.')[0].replace(/^www[-_]/, '');
	if (!label) return 'bang';
	return label.charAt(0).toUpperCase() + label.slice(1).replace(/[-_]+/g, ' ');
}

async function buildIndex(): Promise<BangEntry[]> {
	if (index) return index;
	if (!indexing) {
		indexing = load()
			.then((map) => {
				if (Object.keys(map).length === 0) return [] as BangEntry[];
				const byDomain = new Map<string, BangEntry>();
				const byTrigger = new Map<string, BangEntry>();
				for (const [trigger, url] of Object.entries(map)) {
					const key = trigger.toLowerCase();
					const host = hostOf(url);
					if (!host) continue;
					const domain = domainOf(host);
					const existing = byDomain.get(domain);
					if (existing) {
						existing.triggers.push(key);
						byTrigger.set(key, existing);
						continue;
					}
					const entry: BangEntry = {
						url,
						host,
						name: nameOf(host),
						triggers: [key]
					};
					byDomain.set(domain, entry);
					byTrigger.set(key, entry);
				}
				triggerIndex = byTrigger;
				index = [...byDomain.values()];
				return index;
			})
			.finally(() => (indexing = null));
	}
	return indexing;
}

export async function resolveBang(q: string): Promise<BangMatch | null> {
	const parsed = parseBang(q);
	if (!parsed) return null;
	await buildIndex();
	const entry = triggerIndex?.get(parsed.trigger);
	if (!entry) return null;
	return {
		trigger: parsed.trigger,
		name: entry.name,
		url: entry.url,
		host: entry.host
	};
}

const FEATURED_TRIGGERS = [
	'g',
	'yt',
	'gh',
	'gpt',
	'w',
	'ddg',
	'so',
	'r',
	'imdb',
	'mdn',
	'claude',
	'chatgpt',
	'x',
	'wa',
	'scholar',
	'arxiv',
	'spotify',
	'netflix',
	'twitch',
	'reddit',
	'pin',
	'ste',
	'lb'
];

export async function featuredBangs(limit = 8): Promise<BangMatch[]> {
	await buildIndex();
	const triggers = [...FEATURED_TRIGGERS];
	for (let i = triggers.length - 1; i > 0; i--) {
		const j = Math.floor(Math.random() * (i + 1));
		[triggers[i], triggers[j]] = [triggers[j], triggers[i]];
	}
	const seen = new Set<string>();
	const out: BangMatch[] = [];
	for (const trigger of triggers) {
		const entry = triggerIndex?.get(trigger);
		if (!entry) continue;
		const name = entry.name.toLowerCase();
		if (seen.has(name)) continue;
		seen.add(name);
		out.push({ trigger, name: entry.name, url: entry.url, host: entry.host });
		if (out.length >= limit) break;
	}
	return out;
}

function fuzzy(query: string, target: string): number {
	let t = 0;
	let g = 0;
	let n = 0;
	for (const ch of target) {
		if (t >= query.length) break;
		if (ch === query[t]) {
			g++;
			t++;
		} else {
			n = Math.max(g, n);
			g = 0;
		}
	}
	return t + n + 3 * g;
}

export async function searchBangs(rawQuery: string, limit = 8): Promise<BangMatch[]> {
	const query = rawQuery.trim().toLowerCase().replace(/\s+/g, '');
	if (!query) return featuredBangs(limit);

	const entries = await buildIndex();
	const scored: { trigger: string; entry: BangEntry; score: number }[] = [];

	for (const entry of entries) {
		let best = 0;
		let bestTrigger = entry.triggers[0];
		for (const trigger of entry.triggers) {
			let score: number;
			if (trigger === query) score = 40;
			else if (trigger.startsWith(query)) score = 20 + query.length;
			else score = 0.8 * fuzzy(query, trigger);
			if (score > best) {
				best = score;
				bestTrigger = trigger;
			}
		}
		const name = entry.name.toLowerCase();
		if (name.startsWith(query)) best = Math.max(best, 14 + query.length);
		else if (name.includes(query)) best = Math.max(best, 10);
		if (best <= 0) continue;
		scored.push({ trigger: bestTrigger, entry, score: best });
	}

	scored.sort(
		(a, b) =>
			b.score - a.score ||
			a.entry.name.length - b.entry.name.length ||
			a.trigger.length - b.trigger.length
	);

	const seen = new Set<string>();
	const out: BangMatch[] = [];
	for (const { trigger, entry } of scored) {
		const name = entry.name.toLowerCase();
		if (seen.has(name)) continue;
		seen.add(name);
		out.push({ trigger, name: entry.name, url: entry.url, host: entry.host });
		if (out.length >= limit) break;
	}
	return out;
}
