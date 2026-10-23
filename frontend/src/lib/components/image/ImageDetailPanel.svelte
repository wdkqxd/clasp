<script lang="ts">
	import type { ImageResult } from '../../types';
	import { getSearchSettings } from '../settings/searchSettings.svelte';
	import { ChevronLeft, ChevronRight, ArrowUpRight, X } from '@lucide/svelte';

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

	const MAX_URL_LEN = 40;

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
			const h = new URL(u).hostname.replace('www.', '');
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
				<ChevronLeft class="size-4" />
			</button>
			<button
				class="text-[#5c5c5c] hover:text-[#888] transition-colors duration-100 cursor-pointer disabled:opacity-30"
				disabled={index >= total - 1}
				onclick={onnext}
				aria-label="Next image"
			>
				<ChevronRight class="size-4" />
			</button>
			<button
				class="text-[#5c5c5c] hover:text-[#888] transition-colors duration-100 cursor-pointer"
				onclick={onclose}
				aria-label="Close detail panel"
			>
				<X class="size-4" />
			</button>
		</div>
	</div>

	<!-- Content -->
	<div class="flex flex-col gap-3 px-4 py-4 overflow-y-auto">
		<div class="flex justify-center bg-[#0a0a0a] rounded-lg">
			<img
				src={image.thumbnail}
				alt={image.title}
				class="w-full max-h-[70vh] object-contain rounded-md"
			/>
		</div>

		<div class="flex flex-col gap-1">
			<a
				href={image.url}
				target={settings.openInNewTab ? '_blank' : undefined}
				rel={settings.openInNewTab ? 'noopener noreferrer' : 'noreferrer'}
				class="text-[16px] text-[#c0c0c0] font-normal leading-snug no-underline hover:underline hover:text-[#fff] transition-colors"
			>{image.title}</a>

			<a
				href={image.url}
				target={settings.openInNewTab ? '_blank' : undefined}
				rel={settings.openInNewTab ? 'noopener noreferrer' : 'noreferrer'}
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
				<span>{displayURL(image.url)}</span>
			</a>

			{#if image.width && image.height}
				<p class="text-[11px] text-[#5c5c5c]">{image.width} x {image.height}</p>
			{/if}
		</div>

		<div class="pt-1">
			<a
				href={image.source_url}
				target={settings.openInNewTab ? '_blank' : undefined}
				rel={settings.openInNewTab ? 'noopener noreferrer' : 'noreferrer'}
				class="inline-flex cursor-pointer items-center gap-1.5 rounded-md border border-[#262626] bg-[#1e1e1e] px-3 py-1.5 text-[12px] text-[#d0d0d0] no-underline transition-colors duration-150 ease-out hover:bg-[#262626] hover:text-[#e0e0e0] active:scale-[0.97]"
			>
				view source
				<ArrowUpRight class="size-3.5" />
			</a>
		</div>
	</div>
</div>
