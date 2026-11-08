<script lang="ts">
	import HashHighlight from '$lib/components/navigation/HashHighlight.svelte';
	import PageTopBar from '$lib/components/navigation/PageTopBar.svelte';
	import SettingsItems from '$lib/components/settings/SettingsItems.svelte';
</script>

<svelte:head>
	<title>settings / clasp</title>
</svelte:head>

<div class="flex h-dvh w-full flex-col overflow-hidden bg-[#0c0c0c]">
	<PageTopBar />

	<div
		class="hide-scrollbar min-h-0 flex-1 scroll-smooth overflow-y-auto rounded-t-[14px] border-t border-[#242424] bg-[#101010]"
	>
		<div class="mx-auto w-full max-w-2xl px-4 py-10 sm:py-12" data-hash-root>
			<SettingsItems />
		</div>
	</div>
</div>

<HashHighlight />
