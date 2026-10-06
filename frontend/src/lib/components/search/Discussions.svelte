<script lang="ts">
	import type { Discussion } from '../../types';
	import { getSearchSettings } from '../settings/searchSettings.svelte';
	import { ChevronDown } from '@lucide/svelte';

	let { items }: { items: Discussion[] } = $props();

	const settings = getSearchSettings();

	let openSet = $state(new Set<string>());

	function isOpen(u: string): boolean {
		return openSet.has(u);
	}

	function toggle(u: string) {
		const next = new Set(openSet);
		if (next.has(u)) {
			next.delete(u);
		} else {
			next.add(u);
		}
		openSet = next;
	}

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

	function favicon(u: string): string {
		try {
			const h = new URL(u).hostname.replace(/^www\./, '');
			return h ? `https://icons.duckduckgo.com/ip3/${h}.ico` : '';
		} catch {
			return '';
		}
	}
</script>

<section>
	<div class="space-y-2.5">
		{#each items as d (d.url)}
			<div
				class="disc-item rounded-md border border-[#262626] bg-[#121212] transition-colors duration-300 ease-out hover:bg-[#181818]"
				class:open={isOpen(d.url)}
			>
				<button
					type="button"
					onclick={() => toggle(d.url)}
					aria-expanded={isOpen(d.url)}
					class="flex w-full cursor-pointer items-center justify-between gap-3 px-4 py-3 text-left outline-none"
				>
					<span class="min-w-0 flex-1 text-[15px] leading-snug text-[#c0c0c0]">{d.title}</span>
					<span class="shrink-0 transition-transform duration-200 disc-chevron"
						><ChevronDown class="size-4 text-[#5c5c5c]" /></span
					>
				</button>

				<div class="disc-body">
					<div class="disc-body-inner">
						<div class="px-4 pb-3.5">
							{#if d.snippet}
								<p class="max-w-[68ch] text-[13px] leading-relaxed text-[#707070] line-clamp-3">{d.snippet}</p>
							{/if}

							<div class="mt-2.5">
								<a
									href={d.url}
									target={settings.openInNewTab ? '_blank' : undefined}
									rel={settings.openInNewTab ? 'noopener noreferrer' : undefined}
								class="flex items-center gap-1.5 text-[12px] text-[#5c5c5c] transition-colors hover:text-[#888]"
							>
								{#if favicon(d.url)}
									<img
										src={favicon(d.url)}
										alt=""
										width="14"
										height="14"
										loading="lazy"
										class="h-3.5 w-3.5 rounded-sm object-contain opacity-75"
									/>
								{/if}
								<span class="truncate">{displayURL(d.url)}</span>
							</a>
							</div>
						</div>
					</div>
				</div>
			</div>
		{/each}
	</div>
</section>

<style>
	.disc-body {
		display: grid;
		grid-template-rows: 0fr;
		transition:
			grid-template-rows 0.28s ease,
			opacity 0.22s ease;
		opacity: 0;
	}
	.disc-item.open .disc-body {
		grid-template-rows: 1fr;
		opacity: 1;
	}
	.disc-body-inner {
		overflow: hidden;
		min-height: 0;
	}
	.disc-item.open .disc-chevron {
		transform: rotate(180deg);
	}
</style>
