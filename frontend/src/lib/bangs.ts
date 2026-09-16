let cache: Record<string, string> | null = null;
let loading: Promise<Record<string, string>> | null = null;

async function load(): Promise<Record<string, string>> {
	if (cache) return cache;
	if (!loading) {
		loading = fetch('/api/bangs')
			.then((r) => (r.ok ? r.json() : {}))
			.then((m) => (cache = m && typeof m === 'object' ? (m as Record<string, string>) : {}))
			.catch(() => (cache = {}))
			.finally(() => (loading = null));
	}
	return loading;
}

void load();

export async function checkBang(q: string): Promise<string | undefined> {
	const m = /^!(\S+)\s*(.*)$/.exec(q.trim());
	if (!m) return;
	const url = (await load())[m[1].toLowerCase()];
	if (!url) return;
	return url.replace('%s', encodeURIComponent(m[2].trim()));
}
