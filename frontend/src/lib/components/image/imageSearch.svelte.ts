import { imageSearch } from '../../api';
import { BaseSearchState } from '../../baseSearchState.svelte';
import type { SearchEngine } from '../settings/searchSettings.svelte';
import type { ImageSearchResponse, ImageResult } from '../../types';

export const IMAGE_PAGE_SIZE = 90;

export class ImageSearchState extends BaseSearchState<ImageResult, ImageSearchResponse> {
	constructor() {
		super(IMAGE_PAGE_SIZE);
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
