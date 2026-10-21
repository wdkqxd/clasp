import { imageSearch } from '../../api';
import { BaseSearchState } from '../../baseSearchState.svelte';
import type { ImageResult, ImageSearchResponse } from '../../types';
import { getSearchSettings, type SearchEngine } from '../settings/searchSettings.svelte';

export const IMAGE_PAGE_SIZE = 90;

const TRANSPARENT_QUERY_RE = /\b(?:no\s+(?:bg|background)|transparent)\b/i;

export class ImageSearchState extends BaseSearchState<ImageResult, ImageSearchResponse> {
	constructor() {
		super(IMAGE_PAGE_SIZE, 'imageEngines');
	}

	async run(query: string): Promise<ImageSearchResponse | null> {
		getSearchSettings().transparentOnly = TRANSPARENT_QUERY_RE.test(query);
		return super.run(query);
	}

	protected async fetchResults(
		query: string,
		offset: number,
		limit: number,
		engines?: SearchEngine[]
	): Promise<ImageSearchResponse> {
		return imageSearch(query, offset, limit, engines);
	}

	protected getResults(response: ImageSearchResponse): ImageResult[] | null {
		return response.results;
	}

	protected getHasMore(response: ImageSearchResponse): boolean {
		return response.has_more;
	}

	protected getItemKey(item: ImageResult): string {
		return item.thumbnail || item.url;
	}
}
