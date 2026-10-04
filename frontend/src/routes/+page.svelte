<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { openBangOr } from '$lib/bangs';
	import SearchForm from '$lib/components/search/SearchForm.svelte';
	import SettingsPopup from '$lib/components/settings/SettingsPopup.svelte';
	import { getSearchSettings } from '$lib/components/settings/searchSettings.svelte';
	import logo from '$lib/assets/logo.svg';

	let settingsOpen = $state(false);
	const settings = getSearchSettings();

	onMount(() => {
		if (!page.url.searchParams.has('settings')) return;
		setTimeout(() => {
			settingsOpen = true;
		}, 400);
	});

	function onSearch(q: string) {
		openBangOr(q, `/search?q=${encodeURIComponent(q)}`, goto);
	}

</script>

<svelte:head>
	<title>clasp</title>
</svelte:head>

<div class="relative flex min-h-screen flex-col px-4">
	<div class="flex flex-1 items-center justify-center">
		<div class="w-full max-w-xl">
			<SettingsPopup bind:open={settingsOpen}>
				{#snippet trigger()}
					<img src={logo} alt="" class="block h-12 w-auto" />
				{/snippet}
			</SettingsPopup>
			<SearchForm {onSearch} />
		</div>
	</div>

	<footer class="flex items-center justify-center gap-5 pb-6 text-[13px] text-[#707070]">
		<a
			href="https://github.com/wdkqxd/clasp"
			target={settings.openInNewTab ? '_blank' : undefined}
			rel="noreferrer"
			class="transition-colors hover:text-[#c0c0c0]">github</a
		>
		<a href="/about/privacy" class="transition-colors hover:text-[#c0c0c0]">privacy</a>
		<a href="/about/general" class="transition-colors hover:text-[#c0c0c0]">about</a>
	</footer>
</div>
