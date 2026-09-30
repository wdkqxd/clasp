import type { EngineStatus } from './types';
import { getSearchSettings, type SearchEngine } from './components/settings/searchSettings.svelte';

interface BaseResponse {
	took_ms: number;
	engines: Record<string, EngineStatus>;
}

export abstract class BaseSearchState<TItem, TResponse extends BaseResponse> {
	query = $state('');
	items = $state<TItem[]>([]);
	error = $state('');
	loading = $state(false);
	loadingMore = $state(false);
	hasMore = $state(true);
	page = $state(0);
	tookMs = $state<number | null>(null);
	engineStatus = $state<Record<string, EngineStatus>>({});

	#seen = new Set<string>();

	protected readonly pageSize: number;

	constructor(pageSize: number) {
		this.pageSize = pageSize;
	}

	protected abstract fetchResults(
		query: string,
		offset: number,
		limit: number,
		engines?: SearchEngine[]
	): Promise<TResponse>;

	protected abstract getResults(response: TResponse): TItem[] | null;

	protected abstract getHasMore(response: TResponse): boolean;

	protected abstract getItemKey(item: TItem): string;

	async run(q: string): Promise<TResponse | null> {
		this.query = q;
		this.items = [];
		this.error = '';
		this.page = 0;
		this.tookMs = null;
		this.engineStatus = {};
		this.hasMore = true;
		this.loading = true;
		this.#seen = new Set();
		try {
			const settings = getSearchSettings();
			const res = await this.fetchResults(q, 0, this.pageSize, settings.engines);
			this.tookMs = res.took_ms;
			this.engineStatus = res.engines ?? {};
			this.#append(this.getResults(res), this.getHasMore(res));
			return res;
		} catch (e) {
			this.error = e instanceof Error ? e.message : String(e);
			this.hasMore = false;
			return null;
		} finally {
			this.loading = false;
		}
	}

	async loadMore(): Promise<void> {
		if (this.loading || this.loadingMore || !this.hasMore || this.query === '') return;
		this.loadingMore = true;
		const next = this.page + 1;
		try {
			const settings = getSearchSettings();
			const res = await this.fetchResults(this.query, next, this.pageSize, settings.engines);
			this.tookMs = res.took_ms;
			this.engineStatus = res.engines ?? {};
			this.page = next;
			this.#append(this.getResults(res), this.getHasMore(res));
		} catch (e) {
			this.error = e instanceof Error ? e.message : String(e);
			this.hasMore = false;
		} finally {
			this.loadingMore = false;
		}
	}

	#append(results: TItem[] | null = [], hasMore = true): void {
		const list = results ?? [];
		const fresh = list.filter((r) => !this.#seen.has(this.getItemKey(r)));
		for (const r of fresh) this.#seen.add(this.getItemKey(r));
		this.items = [...this.items, ...fresh];
		this.hasMore = hasMore && fresh.length > 0;
	}

	clear(query = ''): void {
		this.query = query;
		this.items = [];
		this.error = '';
		this.loading = false;
		this.loadingMore = false;
		this.hasMore = true;
		this.page = 0;
		this.#seen = new Set();
	}
}
