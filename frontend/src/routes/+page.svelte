<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { checkBang } from '$lib/bangs';
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
		checkBang(q).then((bang) => {
			if (bang) {
				window.location.href = bang;
				return;
			}
			goto(`/search?q=${encodeURIComponent(q)}`);
		});
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
			<SearchForm {onSearch} placeholder="search..." />
		</div>
	</div>

	<footer class="flex items-center justify-center gap-6 pb-6 text-[12px] text-[#3a3a3a]">
		<a
			href="https://github.com/wdkqxd/clasp"
			target={settings.openInNewTab ? '_blank' : undefined}
			rel="noreferrer"
			class="transition-colors hover:text-[#8a8a8a]">github</a
		>
		<a href="/privacy" class="transition-colors hover:text-[#8a8a8a]">privacy</a>
		<a href="/about" class="transition-colors hover:text-[#8a8a8a]">about</a>
	</footer>
</div>
