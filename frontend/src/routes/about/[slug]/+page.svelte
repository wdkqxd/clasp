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

	<div class="prose mt-10">
		{@html data.doc.html}
	</div>
</div>

<style>
	.prose {
		font-size: 13px;
		line-height: 1.625;
		color: #c0c0c0;
	}
	.prose :global(h2) {
		margin-top: 2rem;
		margin-bottom: 0.375rem;
		font-size: 15px;
		line-height: 1.375;
		font-weight: 600;
		color: #c0c0c0;
		width: fit-content;
		scroll-margin-top: 1.5rem;
	}
	.prose :global(.doc-section:first-child h2:first-child) {
		margin-top: 0;
	}
	.prose :global(p) {
		margin: 0;
	}
	.prose :global(p + p) {
		margin-top: 1.25rem;
	}
	.prose :global(a) {
		color: inherit;
		text-decoration: underline;
		text-underline-offset: 4px;
		transition: color 0.15s;
	}
	.prose :global(a:hover) {
		color: #fff;
	}
	.prose :global(a.heading-anchor) {
		display: inline-flex;
		margin-left: 0.375rem;
		vertical-align: middle;
		color: #707070;
		opacity: 0;
		text-decoration: none;
		transition:
			color 0.15s,
			opacity 0.15s;
	}
	.prose :global(.doc-section:hover .heading-anchor),
	.prose :global(.doc-section:focus-within .heading-anchor) {
		opacity: 1;
	}
	.prose :global(a.heading-anchor:hover) {
		color: #c0c0c0;
	}
</style>