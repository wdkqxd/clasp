<script lang="ts">
	import { page } from '$app/state';
	import HashHighlight from '$lib/components/navigation/HashHighlight.svelte';
	import PageTopBar from '$lib/components/navigation/PageTopBar.svelte';
	import { isOfficialHost } from '$lib/instance';

	let { data } = $props();

	const isOfficialInstance = $derived(isOfficialHost(page.url.hostname));
</script>

<svelte:head>
	<title>{data.doc.title} / clasp</title>
</svelte:head>

<div class="flex h-dvh w-full flex-col overflow-hidden bg-[#0c0c0c]">
	<PageTopBar />

	<div
		class="hide-scrollbar min-h-0 flex-1 overflow-y-auto rounded-t-[14px] border-t border-[#242424] bg-[#101010]"
	>
		<div
			class="prose mx-auto w-full max-w-2xl px-4 py-10 sm:py-12"
			class:official-instance={isOfficialInstance}
		>
			{@html data.doc.html}
		</div>
	</div>
</div>

<HashHighlight />

<style>
	.prose {
		font-size: 14px;
		line-height: 1.625;
		color: #c0c0c0;
	}
	.prose :global(h2) {
		margin-top: 2rem;
		margin-bottom: 0.375rem;
		font-size: 16px;
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
	.prose.official-instance :global([data-instance-only]) {
		display: none;
	}
	.prose :global(a) {
		color: inherit;
		text-decoration: none;
		border-bottom: 1px solid color-mix(in srgb, currentColor 40%, transparent);
		border-radius: 1px;
		box-decoration-break: clone;
		-webkit-box-decoration-break: clone;
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
		border-bottom: none;
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