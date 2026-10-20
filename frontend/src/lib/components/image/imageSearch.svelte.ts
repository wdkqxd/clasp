import { imageSearch } from '../../api';
import { BaseSearchState } from '../../baseSearchState.svelte';
import type { ImageResult, ImageSearchResponse } from '../../types';
import {
	getSearchSettings,
	type ImageTypeFilter,
	type SearchEngine,
} from '../settings/searchSettings.svelte';

export const IMAGE_PAGE_SIZE = 90;

const TRANSPARENT_QUERY_RE = /\b(?:no\s+(?:bg|background)|transparent)\b/i;
const IMAGE_TYPE_QUERY_RE = /\b(gif|svg)\b/gi;

function filtersFromQuery(query: string): {
	imageType: ImageTypeFilter;
	transparentOnly: boolean;
} {
	const rest = query.replace(TRANSPARENT_QUERY_RE, ' ').replace(IMAGE_TYPE_QUERY_RE, ' ').trim();
	if (!rest) return { imageType: 'any', transparentOnly: false };
	const type = query.match(/\b(gif|svg)\b/i)?.[0].toLowerCase() as ImageTypeFilter | undefined;
	return {
		imageType: type ?? 'any',
		transparentOnly: TRANSPARENT_QUERY_RE.test(query),
	};
}

export class ImageSearchState extends BaseSearchState<ImageResult, ImageSearchResponse> {
	constructor() {
		super(IMAGE_PAGE_SIZE, 'imageEngines');
	}

	async run(query: string): Promise<ImageSearchResponse | null> {
		const settings = getSearchSettings();
		const filters = filtersFromQuery(query);
		settings.imageType = filters.imageType;
		settings.transparentOnly = filters.transparentOnly;
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
