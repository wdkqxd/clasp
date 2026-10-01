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

	function host(u: string): string {
		try {
			return new URL(u).hostname.replace(/^www\./, '');
		} catch {
			return '';
		}
	}

	function favicon(u: string): string {
		const h = host(u);
		return h ? `https://icons.duckduckgo.com/ip3/${h}.ico` : '';
	}
</script>

<section>
	<div>
		{#each items as d (d.url)}
			<div class="disc-item rounded transition-colors duration-300 ease-out hover:bg-[#1c1c1c]" class:open={isOpen(d.url)}>
				<button
					type="button"
					onclick={() => toggle(d.url)}
					aria-expanded={isOpen(d.url)}
					class="flex w-full cursor-pointer items-center justify-between gap-3 px-2 py-2 text-left outline-none"
				>
					<span class="min-w-0 flex-1 truncate text-[14px] leading-snug text-[#c0c0c0]">{d.title}</span>
					<span class="icon shrink-0 text-[#5c5c5c] transition-transform duration-200 disc-chevron"><ChevronDown /></span>
				</button>

				<div class="disc-body">
					<div class="disc-body-inner">
						<div class="px-2 pb-3">
							{#if d.snippet}
								<p class="max-w-[52ch] text-[13px] leading-relaxed text-[#707070] line-clamp-3">{d.snippet}</p>
							{/if}

							<div class="mt-2.5">
								<a
									href={d.url}
									target={settings.openInNewTab ? '_blank' : undefined}
									rel={settings.openInNewTab ? 'noopener noreferrer' : undefined}
									class="sitelink inline-flex items-center text-[12px] lowercase text-[#5c5c5c]"
								>
									<span>more on</span>
								{#if favicon(d.url)}
									<img
										src={favicon(d.url)}
										alt=""
										width="14"
										height="14"
										loading="lazy"
										class="ml-2 h-3.5 w-3.5 rounded-sm object-contain opacity-75"
									/>
								{/if}
									<span class="ml-1">{host(d.url)}</span>
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
