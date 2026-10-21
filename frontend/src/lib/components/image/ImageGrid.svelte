<script lang="ts">
	import { fade } from 'svelte/transition';
	import type { ImageResult } from '../../types';
	import { getSearchSettings } from '../settings/searchSettings.svelte';
	import { shouldBlockImage } from './aiBlocklist';
	import ImageThumb from './ImageThumb.svelte';

	const settings = getSearchSettings();

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

	const MAX_URL_LEN = 32;
	let transparency = $state<Record<string, boolean | null>>({});

	function imageKey(image: ImageResult): string {
		return image.thumbnail || image.source_url || image.url;
	}

	function isVisible(image: ImageResult): boolean {
		if (
			settings.hideAIImages &&
			shouldBlockImage({
				url: image.url,
				title: image.title,
				hostname: (() => {
					try {
						return new URL(image.url).hostname;
					} catch {
						return '';
					}
				})()
			})
		) {
			return false;
		}
		return true;
	}

	let visibleItems = $derived(
		images.map((image, index) => ({ image, index })).filter(({ image }) => isVisible(image))
	);

	let transparencyPending = $derived(
		settings.transparentOnly &&
			visibleItems.some(({ image }) => transparency[imageKey(image)] === undefined)
	);

	let hasTransparentItems = $derived(
		visibleItems.some(({ image }) => transparency[imageKey(image)] === true)
	);

	function setTransparency(image: ImageResult, value: boolean | null) {
		const key = imageKey(image);
		if (key) transparency[key] = value;
	}

	function transparencyState(image: ImageResult): boolean | null | undefined {
		return transparency[imageKey(image)];
	}

	function span(image: ImageResult): string {
		const w = image.width || 400;
		const h = image.height || 300;
		const aspect = w / h;
		if (aspect > 1.5) return 'col-span-2';
		return 'col-span-1';
	}

	function displayURL(u: string): string {
		let out: string;
		try {
			const url = new URL(u);
			out = url.hostname.replace(/^www\./, '');
			const path = url.pathname;
			if (path && path !== '/') {
				out += path;
			}
		} catch {
			out = u;
		}
		return out.length > MAX_URL_LEN ? `${out.slice(0, MAX_URL_LEN)}…` : out;
	}

	function faviconURL(u: string): string {
		try {
			const hostname = new URL(u).hostname.replace(/^www\./, '');
			return `https://icons.duckduckgo.com/ip3/${hostname}.ico`;
		} catch {
			return '';
		}
	}

	function openSidebar(e: MouseEvent, index: number) {
		e.preventDefault();
		onselect?.(index);
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
{:else if visibleItems.length === 0}
	<p class="text-[#555] text-sm text-center py-12">no images found</p>
{:else if settings.transparentOnly && !transparencyPending && !hasTransparentItems}
	<p class="text-[#555] text-sm text-center py-12">no images found</p>
{:else}
	<div in:fade={{ duration: 300 }} class="grid grid-cols-[repeat(auto-fill,minmax(180px,1fr))] gap-3 pb-5" style="grid-auto-flow: dense; grid-auto-rows: 220px;">
		{#each visibleItems as { image, index } (imageKey(image) || index)}
			{@const transState = transparencyState(image)}
			<div
				class="flex flex-col text-left rounded-lg overflow-hidden image-card"
				class:col-span-2={span(image) === 'col-span-2'}
				class:invisible={settings.transparentOnly && transState === undefined}
				class:hidden={settings.transparentOnly && (transState === false || transState === null)}
			>
				<a
					href={image.url}
					class="overflow-hidden rounded-lg bg-[#1a1a1a] w-full h-full relative block cursor-pointer"
					class:selected={selectedIndex === index}
					onclick={(e) => openSidebar(e, index)}
				>
					<ImageThumb
						src={image.thumbnail}
						alt={image.title}
						selected={selectedIndex === index}
						detectTransparency={settings.transparentOnly && transState === undefined}
						ontransparency={(value) => setTransparency(image, value)}
					/>
				</a>
				<div class="flex flex-col gap-0.5 mt-2 px-0.5 image-info">
					<a
						href={image.url}
						class="text-xs text-[#ccc] leading-[1.35] line-clamp-1 no-underline"
						target={settings.openInNewTab ? '_blank' : undefined}
						rel={settings.openInNewTab ? 'noopener noreferrer' : undefined}
					>{image.title}</a>
					<a
						href={image.url}
						class="flex items-center gap-1.5 no-underline"
						target={settings.openInNewTab ? '_blank' : undefined}
						rel={settings.openInNewTab ? 'noopener noreferrer' : undefined}
					>
						{#if image.url}
							<img
								src={faviconURL(image.url)}
								alt=""
								class="w-3.5 h-3.5 rounded-sm flex-shrink-0 opacity-75"
								onerror={(e) => (e.currentTarget as HTMLElement).style.display = 'none'}
							/>
						{/if}
						<span class="text-[11px] text-[#5c5c5c] truncate">{displayURL(image.url)}</span>
					</a>
				</div>
			</div>
		{/each}
	</div>
{/if}

<style>
	.image-info:hover a {
		text-decoration: underline;
		text-underline-offset: 3px;
	}
	.image-card:hover :global(img:not(.selected-img)) {
		transform: scale(1.02);
	}
</style>
