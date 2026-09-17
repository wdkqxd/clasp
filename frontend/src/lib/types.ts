export interface Sitelink {
	title: string;
	url: string;
}

export interface SearchResult {
	title: string;
	url: string;
	snippet: string;
	engine: string;
	engines?: string[];
	sitelinks?: Sitelink[];
}

export interface EngineStatus {
	status: string;
	count?: number;
	took_ms?: number;
	error?: string;
}

export interface SearchInfobox {
	title: string;
	url: string;
}

export interface SearchResponse {
	query: string;
	took_ms: number;
	results: SearchResult[];
	engines: Record<string, EngineStatus>;
	has_more: boolean;
	infobox?: SearchInfobox;
}

export interface WikipediaArticle {
	title: string;
	description?: string;
	extract?: string;
	url: string;
	thumbnail?: string;
	facts?: WikipediaFact[];
}

export interface WikipediaFact {
	label: string;
	value: string;
}

export interface MarketPoint {
	time: number;
	value: number;
}

export interface MarketQuote {
	query: string;
	symbol: string;
	name: string;
	currency: string;
	current: number;
	change: number;
	change_percent: number;
	period: string;
	points: MarketPoint[];
}

export interface ImageResult {
	title: string;
	url: string;
	source_url?: string;
	thumbnail: string;
	width: number;
	height: number;
	engine: string;
	engines?: string[];
}

export interface ImageSearchResponse {
	query: string;
	took_ms: number;
	results: ImageResult[];
	engines: Record<string, EngineStatus>;
	has_more: boolean;
}
