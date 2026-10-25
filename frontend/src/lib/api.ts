import type { SearchEngine, SearchTimeout } from './components/settings/searchSettings.svelte';
import type {
	HealthResponse,
	ImageSearchResponse,
	SearchResponse,
	SuggestResponse,
	WikipediaArticle,
} from './types';

export async function health(): Promise<HealthResponse | null> {
	try {
		const res = await fetch('/api/healthz');
		if (!res.ok) return null;
		return res.json();
	} catch {
		return null;
	}
}

export async function search(
	query: string,
	offset = 0,
	limit = 10,
	engines?: SearchEngine[],
	timeout?: SearchTimeout
): Promise<SearchResponse> {
	const engineParam = engines?.length ? `&engines=${encodeURIComponent(engines.join(','))}` : '';
	const timeoutParam = timeout ? `&timeout=${timeout}` : '';
	const res = await fetch(
		`/api/search?q=${encodeURIComponent(query)}&limit=${limit}&offset=${offset}${engineParam}${timeoutParam}`
	);
	if (!res.ok) throw new Error(`HTTP ${res.status}`);
	return res.json();
}

export async function suggest(query: string, engine: string): Promise<SuggestResponse> {
	const res = await fetch(
		`/api/suggest?q=${encodeURIComponent(query)}&engine=${encodeURIComponent(engine)}`
	);
	if (!res.ok) throw new Error(`HTTP ${res.status}`);
	return res.json();
}

export async function wikipedia(
	query: string,
	language = 'auto'
): Promise<WikipediaArticle | null> {
	const res = await fetch(
		`/api/wikipedia?q=${encodeURIComponent(query)}&lang=${encodeURIComponent(language)}`
	);
	if (res.status === 404) return null;
	if (!res.ok) throw new Error(`HTTP ${res.status}`);
	return res.json();
}

export async function imageSearch(
	query: string,
	offset = 0,
	limit = 20,
	engines?: SearchEngine[],
	timeout?: SearchTimeout
): Promise<ImageSearchResponse> {
	const engineParam = engines?.length ? `&engines=${encodeURIComponent(engines.join(','))}` : '';
	const timeoutParam = timeout ? `&timeout=${timeout}` : '';
	const res = await fetch(
		`/api/images?q=${encodeURIComponent(query)}&limit=${limit}&offset=${offset}${engineParam}${timeoutParam}`
	);
	if (!res.ok) throw new Error(`HTTP ${res.status}`);
	return res.json();
}
