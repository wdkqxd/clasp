import { search as apiSearch } from '../../api';
import { BaseSearchState } from '../../baseSearchState.svelte';
import type { SearchEngine } from '../settings/searchSettings.svelte';
import type { SearchResponse, SearchResult } from '../../types';

export const PAGE_SIZE = 10;

export class SearchState extends BaseSearchState<SearchResult, SearchResponse> {
	constructor() {
		super(PAGE_SIZE);
	}

	protected async fetchResults(
		query: string,
		offset: number,
		limit: number,
		engines?: SearchEngine[]
	): Promise<SearchResponse> {
		return apiSearch(query, offset, limit, engines);
	}

	protected getResults(response: SearchResponse): SearchResult[] | null {
		return response.results;
	}

	protected getHasMore(response: SearchResponse): boolean {
		return response.has_more;
	}

	protected getItemKey(item: SearchResult): string {
		return item.url;
	}
}
