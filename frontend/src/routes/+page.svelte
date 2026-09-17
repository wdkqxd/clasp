<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';
	import { checkBang } from '$lib/bangs';
	import SearchForm from '$lib/components/search/SearchForm.svelte';
	import SettingsPopup from '$lib/components/settings/SettingsPopup.svelte';
	import logo from '$lib/assets/logo.svg';

	let settingsOpen = $state(false);

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

<div class="min-h-screen relative flex items-center justify-center px-4">
	<div class="w-full max-w-xl">
		<SettingsPopup bind:open={settingsOpen}>
			{#snippet trigger()}
				<img src={logo} alt="" class="block h-12 w-auto" />
			{/snippet}
		</SettingsPopup>
		<SearchForm {onSearch} placeholder="Search..." />
	</div>
</div>
