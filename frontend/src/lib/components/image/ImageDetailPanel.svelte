<script lang="ts">
	import type { ImageResult } from '../../types';
	import { getSearchSettings } from '../settings/searchSettings.svelte';

	const settings = getSearchSettings();

	let {
		image,
		index,
		total,
		onclose,
		onprev,
		onnext
	}: {
		image: ImageResult;
		index: number;
		total: number;
		onclose: () => void;
		onprev: () => void;
		onnext: () => void;
	} = $props();

	function hostname(url: string): string {
		try {
			return new URL(url).hostname.replace('www.', '');
		} catch {
			return '';
		}
	}

	function faviconURL(url: string): string {
		try {
			const h = new URL(url).hostname.replace('www.', '');
			return `https://icons.duckduckgo.com/ip3/${h}.ico`;
		} catch {
			return '';
		}
	}
</script>

<div class="flex flex-col h-full">
	<!-- Header: nav + close -->
	<div class="flex items-center justify-between px-4 py-3 border-b border-[#ffffff0f] flex-shrink-0">
		<span class="text-[13px] text-[#5c5c5c]">{index + 1} / {total}</span>
		<div class="flex gap-1.5">
			<button
				class="w-[30px] h-[30px] rounded-md bg-[#1a1a1a] text-[#5c5c5c] flex items-center justify-center hover:bg-[#252525] hover:text-[#888] transition-colors duration-100 cursor-pointer disabled:opacity-30"
				disabled={index <= 0}
				onclick={onprev}
				aria-label="Previous image"
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M15 18l-6-6 6-6"/></svg>
			</button>
			<button
				class="w-[30px] h-[30px] rounded-md bg-[#1a1a1a] text-[#5c5c5c] flex items-center justify-center hover:bg-[#252525] hover:text-[#888] transition-colors duration-100 cursor-pointer disabled:opacity-30"
				disabled={index >= total - 1}
				onclick={onnext}
				aria-label="Next image"
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M9 18l6-6-6-6"/></svg>
			</button>
			<button
				class="w-[30px] h-[30px] rounded-md bg-[#1a1a1a] text-[#5c5c5c] flex items-center justify-center hover:bg-[#252525] hover:text-[#888] transition-colors duration-100 cursor-pointer"
				onclick={onclose}
				aria-label="Close detail panel"
			>
				<svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M18 6l-12 12"/><path d="M6 6l12 12"/></svg>
			</button>
		</div>
	</div>

	<!-- Content -->
	<div class="flex flex-col gap-3 px-4 py-4 overflow-y-auto">
		<div class="flex justify-center bg-[#0a0a0a] rounded-lg p-2">
			<img
				src={image.thumbnail}
				alt={image.title}
				class="w-full max-h-[45vh] object-contain rounded"
			/>
		</div>

		<div class="flex flex-col gap-1">
			<h3 class="text-[14px] text-[#c0c0c0] font-normal leading-snug">{image.title}</h3>

			<a
				href={image.url}
				target={settings.openInNewTab ? '_blank' : undefined}
				rel={settings.openInNewTab ? 'noopener noreferrer' : undefined}
				class="flex items-center gap-1.5 text-[12px] text-[#5c5c5c] hover:text-[#888] no-underline transition-colors"
			>
				{#if image.url}
					<img
						src={faviconURL(image.url)}
						alt=""
						class="w-3.5 h-3.5 rounded-sm opacity-75"
						onerror={(e) => (e.currentTarget as HTMLElement).style.display = 'none'}
					/>
				{/if}
				<span>{hostname(image.url)}</span>
			</a>

			{#if image.width && image.height}
				<p class="text-[11px] text-[#5c5c5c]">{image.width} x {image.height}</p>
			{/if}
		</div>

		<div class="pt-1">
			<a
				href={image.url}
				target={settings.openInNewTab ? '_blank' : undefined}
				rel={settings.openInNewTab ? 'noopener noreferrer' : undefined}
				class="inline-block px-4 py-2 bg-[#1e1e1e] text-[#d0d0d0] text-[12px] rounded-md no-underline hover:bg-[#262626] transition-colors duration-100"
			>
				View source
			</a>
		</div>
	</div>
</div>
