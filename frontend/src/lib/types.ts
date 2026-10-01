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

export interface Discussion {
	title: string;
	url: string;
	snippet?: string;
}

export interface RichAnswer {
	type: string;
	calculator?: RichCalculator;
	colorpicker?: RichColorpicker;
	weather?: RichWeather;
	timer?: RichTimer;
	unitConversion?: RichUnitConversion;
	timezones?: RichTimezones;
	cryptocurrency?: RichCryptocurrency;
	unixtimestamp?: RichUnixtimestamp;
	currency?: RichCurrency;
	news?: RichNews;
	topAnswer?: RichTopAnswer;
	stopwatch?: boolean;
}

export interface RichAnswerBlock {
	kind: string;
	text?: string;
	items?: string[];
}

export interface RichTopAnswer {
	question: string;
	blocks?: RichAnswerBlock[];
	author?: string;
	upvotes?: number;
	title?: string;
	url: string;
	meta_url?: RichMetaURL;
}

export interface RichMetaURL {
	hostname?: string;
	favicon?: string;
}

export interface RichCalculator {
	expression: string;
	answer: number;
}

export interface RichColorpicker {
	color: string;
	format?: string;
}

export interface RichWeatherCurrent {
	temp: number;
	feels_like?: number;
	humidity?: number;
	wind_speed?: number;
	uvi?: number;
	pressure?: number;
	weather?: RichWeatherCond;
}

export interface RichWeatherCond {
	main: string;
	description?: string;
	icon?: string;
}

export interface RichWeatherDay {
	ts?: number;
	date_i18n?: string;
	temperature?: { max?: number; min?: number; day?: number };
	weather?: RichWeatherCond;
}

export interface RichWeather {
	location?: { name: string; country?: string; state?: string };
	current_weather?: RichWeatherCurrent;
	daily?: RichWeatherDay[];
	alerts?: RichWeatherAlert[];
}

export interface RichWeatherAlert {
	event?: string;
	description?: string;
	start_relative_i18n?: string;
}

export interface RichTimer {
	duration: number;
	start_on_load: boolean;
}

export interface RichUnitConversion {
	amount: number;
	from_unit: string;
	to_unit: string;
	dimensionality: string;
}

export interface RichTimezone {
	abbreviation?: string;
	generic_name?: string;
	utc_offset?: string;
	converted_time?: { strftime?: string; location?: string; utc_diff?: string; city?: string };
}

export interface RichTimezones {
	type?: string;
	timezones?: RichTimezone[];
}

export interface RichCryptoQuote {
	id?: string;
	symbol: string;
	name: string;
	image?: string;
	current_price: number;
	market_cap?: number;
	high_24h?: number;
	low_24h?: number;
	price_change_24h?: number;
	price_change_percentage_24h?: number;
}

export interface RichCryptocurrency {
	intent_type?: string;
	vs_currency?: string;
	quote?: RichCryptoQuote;
	timeseries?: { time_range?: string; ts_price?: [number, number][] };
}

export interface RichUnixtimestamp {
	conversion?: Record<string, unknown>;
}

export interface RichCurrency {
	amount: number;
	from_currency_code: string;
	to_currency_code: string;
	from_currency_name?: string;
	to_currency_name?: string;
	result: number;
	exchange_rate?: number;
	date?: string;
	timeseries?: { time_range?: string; ts_exchange_rate?: [string, number][] };
}

export interface RichNewsArticle {
	title: string;
	url: string;
	description?: string;
	img?: string;
	publish_time?: number;
	publisher_name?: string;
	page_age?: string;
	meta_url?: { hostname?: string; favicon?: string };
}

export interface RichNews {
	is_topical?: boolean;
	topics?: string[];
	articles?: RichNewsArticle[];
}

export interface SearchResponse {
	query: string;
	took_ms: number;
	results: SearchResult[];
	engines: Record<string, EngineStatus>;
	has_more: boolean;
	infobox?: SearchInfobox;
	discussions?: Discussion[];
	rich?: RichAnswer[];
}

export interface SuggestResponse {
	query: string;
	engine: string;
	suggestions: string[];
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
