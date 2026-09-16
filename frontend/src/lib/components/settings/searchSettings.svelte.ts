export const SEARCH_SETTINGS_KEY = 'metasearch.search.settings';

export const SEARCH_ENGINES = ['Brave', 'DuckDuckGo', 'Startpage'] as const;
export type SearchEngine = (typeof SEARCH_ENGINES)[number];

export type FontChoice = 'jetbrains' | 'system';

export type SearchSettings = {
	engines: SearchEngine[];
	showEngineBadges: boolean;
	openInNewTab: boolean;
	font: FontChoice;
};

const defaults: SearchSettings = {
	engines: [...SEARCH_ENGINES],
	showEngineBadges: false,
	openInNewTab: false,
	font: 'jetbrains',
};

function load(): SearchSettings {
	if (typeof localStorage === 'undefined') return { ...defaults, engines: [...defaults.engines] };
	try {
		const parsed = JSON.parse(localStorage.getItem(SEARCH_SETTINGS_KEY) ?? 'null') ?? {};
		const engines = Array.isArray(parsed.engines)
			? SEARCH_ENGINES.filter((engine) => parsed.engines.includes(engine))
			: [...defaults.engines];
		return {
			engines: engines.length > 0 ? engines : [...defaults.engines],
			showEngineBadges:
				typeof parsed.showEngineBadges === 'boolean'
					? parsed.showEngineBadges
					: defaults.showEngineBadges,
			openInNewTab:
				typeof parsed.openInNewTab === 'boolean' ? parsed.openInNewTab : defaults.openInNewTab,
			font: parsed.font === 'system' ? 'system' : 'jetbrains'
		};
	} catch {
		return { ...defaults, engines: [...defaults.engines] };
	}
}

export const searchSettings = $state<SearchSettings>(load());

export function getSearchSettings(): SearchSettings {
	return searchSettings;
}

export function saveSearchSettings(): void {
	if (typeof localStorage !== 'undefined') {
		localStorage.setItem(SEARCH_SETTINGS_KEY, JSON.stringify(searchSettings));
	}
}