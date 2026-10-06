<script lang="ts">
	import { Loader } from '@lucide/svelte';
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
	<div class="animate-fade-in">
		<SkeletonResults />
	</div>
{:else if feed.error && feed.items.length === 0}
	<p class="text-[13px] text-[#5c5c5c]">error: {feed.error}</p>
{:else if feed.items.length === 0}
	<p class="text-[13px] text-[#5c5c5c]">no results</p>
{:else}
	<div class="space-y-8">
		{#each feed.items as r, i (r.url)}
			{#if discussions.length > 0 && i === PAGE_SIZE - 1}
				<Discussions items={discussions} />
			{/if}
			<div class="animate-result-in" style="animation-delay:{Math.min(baseIndex + i, 10) * 28}ms">
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
		<span class="inline-block animate-spin text-[#5c5c5c]"><Loader size={28} /></span>
	</div>
{/if}
<div bind:this={sentinel} class="h-5"></div>
