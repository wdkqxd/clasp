<script lang="ts">
	import '../app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { onNavigate } from '$app/navigation';
	import { page } from '$app/state';
	import { getSearchSettings, type FontChoice } from '$lib/components/settings/searchSettings.svelte';

	let { children } = $props();

	const settings = getSearchSettings();

	const FONT_STACKS: Record<FontChoice, string> = {
		jetbrains: "'JetBrains Mono', ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace",
		system:
			"ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Cascadia Mono', monospace"
	};

	$effect(() => {
		document.documentElement.style.setProperty('--font-mono', FONT_STACKS[settings.font]);
	});

	onNavigate((navigation) => {
		if (!document.startViewTransition) return;

		const isSearchTabSwitch =
			page.url.pathname === '/search' && navigation.to?.url.pathname === '/search';
		if (isSearchTabSwitch) return;

		return new Promise((resolve) => {
			document.startViewTransition(async () => {
				resolve();
				await navigation.complete;
			});
		});
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<div class="min-h-screen bg-base-950 text-[#d0d0d0] antialiased font-mono">
	{@render children()}
</div>
