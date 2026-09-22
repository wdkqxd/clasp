<script lang="ts">
	import loadingIcon from 'pixelarticons/svg/loader.svg?raw';
	import ResultCard from './ResultCard.svelte';
	import SkeletonResults from './SkeletonResults.svelte';
	import Discussions from './Discussions.svelte';
	import { PAGE_SIZE } from './search.svelte.ts';
	import type { SearchState } from './search.svelte.ts';
	import type { Discussion } from '../../types';

	let { feed, baseIndex = 0, discussions = [] }: { feed: SearchState; baseIndex?: number; discussions?: Discussion[] } = $props();

	let sentinel = $state<HTMLDivElement | null>(null);

	$effect(() => {
		const el = sentinel;
		if (!el) return;
		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0]?.isIntersecting) feed.loadMore();
			},
			{ rootMargin: '400px' }
		);
		observer.observe(el);
		return () => observer.disconnect();
	});
</script>

{#if feed.loading}
	<div class="fade-out">
		<SkeletonResults />
	</div>
{:else if feed.error && feed.items.length === 0}
	<p class="text-[13px] text-[#5c5c5c]">error: {feed.error}</p>
{:else if feed.items.length === 0}
	<p class="text-[13px] text-[#5c5c5c]">no results</p>
{:else}
	<div class="space-y-10">
		{#each feed.items as r, i (r.url)}
			{#if discussions.length > 0 && i === PAGE_SIZE - 1}
				<Discussions items={discussions} />
			{/if}
			<div class="result-item" style="--i:{Math.min(baseIndex + i, 10)}">
				<ResultCard result={r} />
			</div>
		{/each}
		{#if discussions.length > 0 && feed.items.length < PAGE_SIZE}
			<Discussions items={discussions} />
		{/if}
	</div>
{/if}

{#if feed.loadingMore}
	<div class="flex justify-center py-4">
		<span class="inline-block w-7 h-7 spin text-[#5c5c5c]">{@html loadingIcon}</span>
	</div>
{/if}
<div bind:this={sentinel} class="h-5"></div>

<style>
	.result-item {
		animation: resultSlideIn 0.26s cubic-bezier(0.23, 1, 0.32, 1) backwards;
		animation-delay: calc(var(--i) * 28ms);
	}

	@keyframes resultSlideIn {
		from {
			opacity: 0;
			transform: translateY(6px);
		}
	}

	.fade-out {
		animation: fadeOutIn 0.2s ease-out;
	}

	@keyframes fadeOutIn {
		from { opacity: 0; }
		to { opacity: 1; }
	}
</style>
