<script lang="ts">
	import type { SearchResult } from '../../types';
	import { getSearchSettings } from '../settings/searchSettings.svelte';

	let { result }: { result: SearchResult } = $props();

	const settings = getSearchSettings();

	const MAX_URL_LEN = 32;

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

	const ENGINE_LABELS: Record<string, string> = {
		brave: 'brave',
		ddg: 'duckduckgo'
	};

	function engineLabel(e: string): string {
		return ENGINE_LABELS[e] ?? e;
	}

	function resultEngines(): string[] {
		return result.engines?.length ? result.engines : [result.engine];
	}

function cleanSnippet(s: string): string {
		return s
			.replace(/<\/?(?:b|strong|bold|em|i)>/gi, '')
			.replace(/<[^>]+>/g, '')
			.slice(0, 140)
			.replace(/\s+/g, ' ')
			.trim();
	}
</script>

<article class="group">
	<a
		href={result.url}
		target={settings.openInNewTab ? '_blank' : undefined}
		rel={settings.openInNewTab ? 'noopener noreferrer' : 'noreferrer'}
		class="block visited:[&_h2]:text-[#83709f]"
	>
		<div class="mb-1.5 flex items-center gap-1.5 text-[12px] tracking-wide text-[#5c5c5c]">
			{#if faviconURL(result.url)}
				<img
					src={faviconURL(result.url)}
					alt=""
					width="14"
					height="14"
					loading="lazy"
					class="h-3.5 w-3.5 rounded-sm object-contain opacity-75"
				/>
			{/if}
			{displayURL(result.url)}
		</div>
		<h2
			class="truncate text-[16.5px] text-[#c0c0c0] font-normal leading-snug mb-1.5 underline-offset-[3px] transition-colors group-hover:underline"
			title={result.title}
		>
			{result.title}
		</h2>
		<p class="text-[13px] text-[#707070] leading-relaxed">{cleanSnippet(result.snippet)}</p>
	</a>
	<div class="mt-2 flex flex-wrap gap-1.5" class:hidden={!settings.showEngineBadges}>
		{#each resultEngines() as e}
			<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 text-[10.5px] tracking-wide text-[#5c5c5c]">
				{engineLabel(e)}
			</span>
		{/each}
	</div>
	{#if result.sitelinks?.length}
		<div class="grid grid-cols-2 gap-x-12 gap-y-2 mt-4 pl-0">
			{#each result.sitelinks as link}
				<a
					href={link.url}
					target={settings.openInNewTab ? '_blank' : undefined}
					rel={settings.openInNewTab ? 'noopener noreferrer' : 'noreferrer'}
					class="text-[12px] text-[#5c5c5c] underline-offset-[3px] transition-colors hover:underline visited:text-[#83709f]"
				>
					{link.title}
				</a>
			{/each}
		</div>
	{/if}
</article>
