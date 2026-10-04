<script lang="ts">
	import SearchForm from './SearchForm.svelte';
	import TabBar from './TabBar.svelte';
	import type { EngineStatus } from '../../types';

	type HeaderFeed = {
		query: string;
		tookMs: number | null;
		engineStatus: Record<string, EngineStatus>;
	};

	let {
		feed,
		active,
		onSearch,
		placeholder
	}: {
		feed: HeaderFeed;
		active: string;
		onSearch: (q: string) => void;
		placeholder?: string;
	} = $props();
</script>

<header class="mb-5 min-h-[86px]">
	<SearchForm bind:value={feed.query} {onSearch} {placeholder} />
	<TabBar
		tabs={['Web', 'Images', 'Video']}
		{active}
		tookMs={feed.tookMs}
		engineStatus={feed.engineStatus}
	/>
</header>