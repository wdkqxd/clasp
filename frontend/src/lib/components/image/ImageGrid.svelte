<script lang="ts">
	import { fade } from 'svelte/transition';
	import type { ImageResult } from '../../types';

	let {
		images,
		loading = false,
		selectedIndex = -1,
		onselect
	}: {
		images: ImageResult[];
		loading?: boolean;
		selectedIndex?: number;
		onselect?: (index: number) => void;
	} = $props();

	function span(image: ImageResult): string {
		const w = image.width || 400;
		const h = image.height || 300;
		const aspect = w / h;
		if (aspect > 1.5) return 'col-span-2';
		return 'col-span-1';
	}

	function hostname(url: string): string {
		try {
			return new URL(url).hostname.replace('www.', '');
		} catch {
			return '';
		}
	}
</script>

{#if loading && images.length === 0}
	<div out:fade={{ duration: 200 }} class="grid grid-cols-[repeat(auto-fill,minmax(180px,1fr))] gap-3 pb-5" style="grid-auto-flow: dense; grid-auto-rows: 220px;">
		{#each Array(20) as _, i (i)}
			<div class="flex flex-col rounded-lg overflow-hidden">
				<div class="skeleton flex-1 rounded-lg"></div>
				<div class="flex flex-col gap-1 mt-2 px-0.5">
					<div class="skeleton h-3 rounded" style="width: {60 + (i % 3) * 15}%"></div>
					<div class="skeleton h-2.5 rounded" style="width: {40 + (i % 2) * 20}%"></div>
				</div>
			</div>
		{/each}
	</div>
{:else if images.length === 0}
	<p class="text-[#555] text-sm text-center py-12">no images found</p>
{:else}
	<div in:fade={{ duration: 300 }} class="grid grid-cols-[repeat(auto-fill,minmax(180px,1fr))] gap-3 pb-5" style="grid-auto-flow: dense; grid-auto-rows: 220px;">
		{#each images as image, i}
			<button
				class="group flex flex-col text-left cursor-pointer outline-none rounded-lg overflow-hidden"
				class:selected={selectedIndex === i}
				class:col-span-2={span(image) === 'col-span-2'}
				onclick={() => onselect?.(i)}
			>
				<div class="overflow-hidden rounded-lg bg-[#1a1a1a] w-full h-full relative">
					<img
						src={image.thumbnail}
						alt={image.title}
						loading="lazy"
						class="w-full h-full object-cover rounded-lg transition-[transform,opacity] duration-200 ease-out"
						class:selected-img={selectedIndex === i}
					/>
				</div>
				<div class="flex flex-col gap-0.5 mt-2 px-0.5">
					<p class="text-xs text-[#ccc] leading-[1.35] line-clamp-1">{image.title}</p>
					<div class="flex items-center gap-1.5">
						{#if image.url}
							<img
								src="https://icons.duckduckgo.com/ip3/{hostname(image.url)}.ico"
								alt=""
								class="w-4 h-4 rounded-full flex-shrink-0"
								onerror={(e) => (e.currentTarget as HTMLElement).style.display = 'none'}
							/>
						{/if}
						<span class="text-[11px] text-[#777] truncate">{hostname(image.url)}</span>
					</div>
				</div>
			</button>
		{/each}
	</div>
{/if}

<style>
	.selected-img {
		filter: brightness(0.5);
	}
	button:hover img:not(.selected-img) {
		transform: scale(1.02);
	}
</style>
