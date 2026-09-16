<script lang="ts">
	import type { ImageResult } from '../../types';
	import { getSearchSettings } from '../settings/searchSettings.svelte';
	import chevronLeftIcon from 'pixelarticons/svg/chevron-left.svg?raw';
	import chevronRightIcon from 'pixelarticons/svg/chevron-right.svg?raw';
	import closeIcon from 'pixelarticons/svg/close.svg?raw';

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
		<div class="flex gap-3">
			<button
				class="text-[#5c5c5c] hover:text-[#888] transition-colors duration-100 cursor-pointer disabled:opacity-30"
				disabled={index <= 0}
				onclick={onprev}
				aria-label="Previous image"
			>
				<span class="icon">{@html chevronLeftIcon}</span>
			</button>
			<button
				class="text-[#5c5c5c] hover:text-[#888] transition-colors duration-100 cursor-pointer disabled:opacity-30"
				disabled={index >= total - 1}
				onclick={onnext}
				aria-label="Next image"
			>
				<span class="icon">{@html chevronRightIcon}</span>
			</button>
			<button
				class="text-[#5c5c5c] hover:text-[#888] transition-colors duration-100 cursor-pointer"
				onclick={onclose}
				aria-label="Close detail panel"
			>
				<span class="icon">{@html closeIcon}</span>
			</button>
		</div>
	</div>

	<!-- Content -->
	<div class="flex flex-col gap-3 px-4 py-4 overflow-y-auto">
		<div class="flex justify-center bg-[#0a0a0a] rounded-lg p-2">
			<img
				src={image.thumbnail}
				alt={image.title}
				class="w-full max-h-[70vh] object-contain rounded"
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
