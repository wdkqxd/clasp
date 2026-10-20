export const SEARCH_SETTINGS_KEY = 'clasp.search.settings';

export const WEB_ENGINES = ['Brave', 'DuckDuckGo'] as const;
export const IMAGE_ENGINES = ['Brave', 'DuckDuckGo', 'Startpage'] as const;
export type SearchEngine = (typeof IMAGE_ENGINES)[number];

export const ENGINE_API_NAMES: Record<SearchEngine, string> = {
	Brave: 'brave',
	DuckDuckGo: 'ddg',
	Startpage: 'startpage',
};

export const AUTOCOMPLETE_PROVIDERS = [
	{ value: 'kagi', label: 'kagi' },
	{ value: 'ddg', label: 'duckduckgo' },
	{ value: 'brave', label: 'brave' },
	{ value: 'startpage', label: 'startpage' },
] as const;

export type AutocompleteProvider = (typeof AUTOCOMPLETE_PROVIDERS)[number]['value'];

export type FontChoice = 'jetbrains' | 'system';

export const IMAGE_TYPE_FILTERS = ['any', 'gif', 'svg'] as const;
export type ImageTypeFilter = (typeof IMAGE_TYPE_FILTERS)[number];

export type SearchSettings = {
	engines: SearchEngine[];
	imageEngines: SearchEngine[];
	showEngineBadges: boolean;
	openInNewTab: boolean;
	font: FontChoice;
	autocomplete: AutocompleteProvider;
	hideAIImages: boolean;
	imageType: ImageTypeFilter;
	transparentOnly: boolean;
};

const defaults: SearchSettings = {
	engines: [...WEB_ENGINES],
	imageEngines: [...IMAGE_ENGINES],
	showEngineBadges: false,
	openInNewTab: false,
	font: 'jetbrains',
	autocomplete: 'kagi',
	hideAIImages: false,
	imageType: 'any',
	transparentOnly: false,
};

function load(): SearchSettings {
	if (typeof localStorage === 'undefined')
		return {
			...defaults,
			engines: [...defaults.engines],
			imageEngines: [...defaults.imageEngines],
		};
	try {
		const parsed = JSON.parse(localStorage.getItem(SEARCH_SETTINGS_KEY) ?? 'null') ?? {};
		const engines = Array.isArray(parsed.engines)
			? WEB_ENGINES.filter((engine) => parsed.engines.includes(engine))
			: [...defaults.engines];
		const imageEngines = Array.isArray(parsed.imageEngines)
			? IMAGE_ENGINES.filter((engine) => parsed.imageEngines.includes(engine))
			: [...defaults.imageEngines];
		return {
			engines: engines.length > 0 ? engines : [...defaults.engines],
			imageEngines: imageEngines.length > 0 ? imageEngines : [...defaults.imageEngines],
			showEngineBadges:
				typeof parsed.showEngineBadges === 'boolean'
					? parsed.showEngineBadges
					: defaults.showEngineBadges,
			openInNewTab:
				typeof parsed.openInNewTab === 'boolean' ? parsed.openInNewTab : defaults.openInNewTab,
			font: parsed.font === 'system' ? 'system' : 'jetbrains',
			autocomplete: AUTOCOMPLETE_PROVIDERS.some((p) => p.value === parsed.autocomplete)
				? parsed.autocomplete
				: defaults.autocomplete,
			hideAIImages:
				typeof parsed.hideAIImages === 'boolean' ? parsed.hideAIImages : defaults.hideAIImages,
			imageType: defaults.imageType,
			transparentOnly: defaults.transparentOnly,
		};
	} catch {
		return {
			...defaults,
			engines: [...defaults.engines],
			imageEngines: [...defaults.imageEngines],
		};
	}
}

export const searchSettings = $state<SearchSettings>(load());

export function getSearchSettings(): SearchSettings {
	return searchSettings;
}

export function saveSearchSettings(): void {
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem(
			SEARCH_SETTINGS_KEY,
			JSON.stringify({
				engines: searchSettings.engines,
				imageEngines: searchSettings.imageEngines,
				showEngineBadges: searchSettings.showEngineBadges,
				openInNewTab: searchSettings.openInNewTab,
				font: searchSettings.font,
				autocomplete: searchSettings.autocomplete,
				hideAIImages: searchSettings.hideAIImages,
			})
		);
	}
}
