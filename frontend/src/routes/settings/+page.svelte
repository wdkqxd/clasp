<script lang="ts">
	import { ChevronLeft } from '@lucide/svelte';
	import HashHighlight from '$lib/components/navigation/HashHighlight.svelte';
	import SettingsItems from '$lib/components/settings/SettingsItems.svelte';
</script>

<svelte:head>
	<title>settings / clasp</title>
</svelte:head>

<div class="mx-auto w-full max-w-2xl px-4 py-10 sm:py-16 lg:max-w-4xl">
	<a
		href="/"
		class="-my-1 -ml-2 inline-flex items-center gap-1.5 rounded-[5px] px-2 py-1 text-[13px] text-[#707070] no-underline transition-all duration-150 ease-out hover:bg-[#1a1a1a] hover:text-[#c0c0c0] active:scale-[0.94] active:duration-75 lg:ml-3"
	>
		<span class="inline-block size-[14px] shrink-0 align-middle [&>svg]:block [&>svg]:size-full"
			><ChevronLeft /></span
		>
		back
	</a>

	<div class="mt-10">
		<SettingsItems />
	</div>
</div>

<HashHighlight />
