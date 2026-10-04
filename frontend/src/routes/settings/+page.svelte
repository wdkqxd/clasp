<script lang="ts">
	import { ChevronLeft } from '@lucide/svelte';
	import SettingsItems from '$lib/components/settings/SettingsItems.svelte';
</script>

<svelte:head>
	<title>settings / clasp</title>
</svelte:head>

<div class="mx-auto w-full max-w-2xl px-4 py-10 sm:py-16">
	<a
		href="/"
		class="inline-flex items-center gap-1.5 text-[13px] text-[#707070] transition-colors hover:text-[#c0c0c0]"
	>
		<span class="inline-block size-[14px] shrink-0 align-middle [&>svg]:block [&>svg]:size-full"
			><ChevronLeft /></span
		>
		back
	</a>

	<h1 class="mt-10 text-[20px] leading-tight font-normal text-[#d0d0d0]">settings</h1>

	<div class="mt-8 max-w-xl">
		<SettingsItems />
	</div>
</div>