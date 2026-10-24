<script lang="ts">
	import type { ImageResult } from '../../types';
	import { getSearchSettings } from '../settings/searchSettings.svelte';
	import { ChevronLeft, ChevronRight, ArrowUpRight, Check, Copy, Download, X } from '@lucide/svelte';

	const settings = getSearchSettings();
	let copied = $state(false);

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
	const DOWNLOAD_PROXY = 'https://wsrv.nl/?url=';

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

	function imageLink(): string {
		return image.source_url || image.thumbnail || image.url;
	}

	function copyLink() {
		navigator.clipboard?.writeText(imageLink());
		copied = true;
		setTimeout(() => (copied = false), 1200);
	}

	function downloadName(url: string): string {
		try {
			const name = new URL(url).pathname.split('/').pop() ?? '';
			if (/\.[a-z0-9]{2,5}$/i.test(name)) return name;
		} catch {
			// ignore malformed urls
		}
		const base = image.title
			.toLowerCase()
			.replace(/[^a-z0-9]+/g, '-')
			.replace(/^-+|-+$/g, '')
			.slice(0, 48);
		return base || 'image';
	}

	async function downloadImage() {
		const url = imageLink();
		if (!url) return;
		try {
			const res = await fetch(`${DOWNLOAD_PROXY}${encodeURIComponent(url)}`);
			if (!res.ok) throw new Error('failed to download image');
			const blob = await res.blob();
			const ext = blob.type.split('/')[1]?.split('+')[0] ?? '';
			let name = downloadName(url);
			if (ext && !/\.[a-z0-9]{2,5}$/i.test(name)) {
				name += `.${ext === 'jpeg' ? 'jpg' : ext}`;
			}
			const objectURL = URL.createObjectURL(blob);
			const a = document.createElement('a');
			a.href = objectURL;
			a.download = name;
			a.click();
			setTimeout(() => URL.revokeObjectURL(objectURL), 1000);
		} catch {
			window.open(url, '_blank', 'noopener,noreferrer');
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
		<div class="flex justify-center bg-[#0a0a0a] rounded-md">
			<img
				src={image.thumbnail}
				alt={image.title}
				class="w-full max-h-[70vh] object-contain rounded-[5px]"
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

		<div class="flex items-center gap-1.5 pt-1">
			<a
				href={image.source_url}
				target={settings.openInNewTab ? '_blank' : undefined}
				rel={settings.openInNewTab ? 'noopener noreferrer' : 'noreferrer'}
				class="inline-flex cursor-pointer items-center gap-1.5 rounded-[4px] border border-[#262626] bg-[#1e1e1e] px-3 py-1.5 text-[12px] text-[#d0d0d0] no-underline transition-colors duration-150 ease-out hover:bg-[#262626] hover:text-[#e0e0e0] active:scale-[0.97]"
			>
				view source
				<ArrowUpRight class="size-3.5" />
			</a>
			<button
				type="button"
				title="copy image link"
				aria-label="Copy image link"
				onclick={copyLink}
				class="inline-flex size-8 cursor-pointer items-center justify-center rounded-[4px] border border-[#262626] bg-[#1e1e1e] text-[#d0d0d0] transition-colors duration-150 ease-out hover:bg-[#262626] hover:text-[#e0e0e0] active:scale-[0.97]"
			>
				{#if copied}
					<Check class="size-3.5" />
				{:else}
					<Copy class="size-3.5" />
				{/if}
			</button>
			<button
				type="button"
				title="download image"
				aria-label="Download image"
				onclick={downloadImage}
				class="inline-flex size-8 cursor-pointer items-center justify-center rounded-[4px] border border-[#262626] bg-[#1e1e1e] text-[#d0d0d0] transition-colors duration-150 ease-out hover:bg-[#262626] hover:text-[#e0e0e0] active:scale-[0.97]"
			>
				<Download class="size-3.5" />
			</button>
		</div>
	</div>
</div>
