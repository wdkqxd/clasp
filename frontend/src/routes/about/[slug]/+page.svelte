<script lang="ts">
	import { ChevronLeft } from '@lucide/svelte';

	let { data } = $props();
</script>

<svelte:head>
	<title>{data.doc.title} / clasp</title>
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

	<h1 class="mt-10 text-[20px] leading-tight font-normal text-[#d0d0d0]">{data.doc.title}</h1>

	<div class="prose mt-6">
		{@html data.doc.html}
	</div>
</div>

<style>
	.prose {
		font-size: 13px;
		line-height: 1.625;
		color: #707070;
	}
	.prose :global(h2) {
		margin-top: 2rem;
		margin-bottom: 0.375rem;
		font-size: 15px;
		line-height: 1.375;
		font-weight: 400;
		color: #c0c0c0;
	}
	.prose :global(h2:first-child) {
		margin-top: 0;
	}
	.prose :global(a) {
		color: #8a8a8a;
		text-decoration: underline;
		text-underline-offset: 4px;
		transition: color 0.15s;
	}
	.prose :global(a:hover) {
		color: #c0c0c0;
	}
</style>