<script lang="ts">
	import { page } from '$app/state';
	import { Search, Settings2 } from '@lucide/svelte';
	import { tick } from 'svelte';
	import { fade } from 'svelte/transition';
	import {
		ENGINE_API_NAMES,
		getSearchSettings,
		IMAGE_ENGINES,
		WEB_ENGINES,
		saveSearchSettings,
		type SearchEngine
	} from '../settings/searchSettings.svelte';
	import type { EngineStatus } from '../../types';

	let {
		tabs,
		active,
		tookMs = null,
		engineStatus = {}
	}: {
		tabs: string[];
		active: string;
		tookMs?: number | null;
		engineStatus?: Record<string, EngineStatus>;
	} = $props();

	let settingsOpen = $state(false);
	let tabsEl = $state<HTMLDivElement | null>(null);
	let indicatorLeft = $state(0);
	let indicatorWidth = $state(0);
	let indicatorReady = $state(false);
	const settings = getSearchSettings();
	const engines = $derived(active === 'Images' ? settings.imageEngines : settings.engines);
	const engineOptions = $derived(active === 'Images' ? IMAGE_ENGINES : WEB_ENGINES);

	function tabHref(tab: string): string {
		const q = page.url.searchParams.get('q') ?? '';
		const params = new URLSearchParams({ q });
		if (tab === 'Images') params.set('type', 'images');
		return `/search?${params.toString()}`;
	}

	function toggleEngine(engine: SearchEngine) {
		const current = engines;
		if (current.includes(engine) && current.length === 1) return;
		const next = current.includes(engine)
			? current.filter((item) => item !== engine)
			: [...current, engine];
		if (active === 'Images') settings.imageEngines = next;
		else settings.engines = next;
		saveSearchSettings();
	}

	function engineStatusFor(engine: SearchEngine): EngineStatus | undefined {
		return engineStatus[ENGINE_API_NAMES[engine]];
	}

	function updateIndicator() {
		if (!tabsEl) return;
		const activeEl = tabsEl.querySelector<HTMLElement>('[data-active="true"]');
		if (!activeEl) return;

		indicatorLeft = activeEl.offsetLeft;
		indicatorWidth = activeEl.offsetWidth;
		if (!indicatorReady && indicatorWidth > 0) indicatorReady = true;
	}

	$effect(() => {
		const _ = active;
		const __ = tabs;
		tick().then(() => updateIndicator());
	});

	$effect(() => {
		if (!tabsEl) return;
		const ro = new ResizeObserver(() => updateIndicator());
		ro.observe(tabsEl);
		return () => ro.disconnect();
	});

	$effect(() => {
		if (!settingsOpen) return;
		const closeOnOutside = (event: MouseEvent) => {
			if (event.target instanceof Element && !event.target.closest('[data-search-settings]')) {
				settingsOpen = false;
			}
		};
		const closeOnEscape = (event: KeyboardEvent) => {
			if (event.key === 'Escape') settingsOpen = false;
		};
		document.addEventListener('click', closeOnOutside);
		document.addEventListener('keydown', closeOnEscape);
		return () => {
			document.removeEventListener('click', closeOnOutside);
			document.removeEventListener('keydown', closeOnEscape);
		};
	});
</script>

<nav class="relative flex items-center gap-5 pt-3 text-[13px] font-medium" data-search-settings>
	<div
		bind:this={tabsEl}
		class="relative inline-flex max-w-full items-center rounded-[5px] bg-[#141414] p-[3px]"
	>
		<span
			class="pointer-events-none absolute left-0 top-[3px] bottom-[3px] rounded-[3px] border border-[#2f2f2f] bg-[#222222] {indicatorReady
				? 'transition-[transform,width] duration-200 ease-[cubic-bezier(0.22,1,0.36,1)]'
				: ''}"
			style:transform="translateX({indicatorLeft}px)"
			style:width="{indicatorWidth}px"
			style:opacity={indicatorReady && indicatorWidth > 0 ? 1 : 0}
		></span>

		{#each tabs as tab}
			{@const isActive = tab.toLowerCase() === active.toLowerCase()}
			<a
				href={tab === 'Video' ? '#' : tabHref(tab)}
				data-active={isActive}
				class="relative z-10 inline-flex cursor-pointer items-center gap-1.5 rounded-[3px] px-2.5 py-1 text-[12px] leading-tight select-none transition-colors duration-150 {isActive
					? 'text-[#c0c0c0]'
					: 'text-[#6a6a6a] hover:text-[#9a9a9a]'}"
				onclick={tab === 'Video' ? (e) => e.preventDefault() : undefined}
			>
				{#if tab === 'All'}
					<Search class="size-3.5" />
				{/if}
				{tab}
			</a>
		{/each}
	</div>

	<div class="ml-auto flex items-center gap-3">
		{#if tookMs !== null}
			<span transition:fade={{ duration: 300 }} class="text-[11px] text-[#4a4a4a]">{tookMs} ms</span>
		{/if}

		<div class="relative">
		<button
			type="button"
			aria-label="Search settings"
			aria-expanded={settingsOpen}
			onclick={() => (settingsOpen = !settingsOpen)}
			class="inline-flex items-center text-[#4a4a4a] hover:text-[#9a9a9a] transition-colors {settingsOpen ? 'text-[#9a9a9a]' : ''}"
		>
			<Settings2 class="size-3.5 translate-y-[2px]" />
		</button>

			<div
				class="absolute right-0 top-7 z-20 w-56 origin-top-right rounded-[5px] border border-[#292929] bg-[#121212] p-3 shadow-2xl shadow-black/30 transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] {settingsOpen
					? 'visible translate-y-0 scale-100 opacity-100 blur-0'
					: 'pointer-events-none invisible -translate-y-1 scale-95 opacity-0 blur-sm'}"
			>
				<div class="space-y-3">
					<div>
						<div class="space-y-2">
							{#each engineOptions as engine}
								{@const status = engineStatusFor(engine)}
								<button
									type="button"
									class="flex w-full items-center justify-between text-left text-[12px] text-[#9b9b9b]"
									onclick={() => toggleEngine(engine)}
								>
									<span class="flex min-w-0 items-center gap-1.5">
										<span class="truncate">{engine.toLowerCase()}</span>
										{#if status?.status === 'error'}
											<span class="group/err relative shrink-0 cursor-help text-[11px] text-[#d05a5a]">
												error
												<span
													class="pointer-events-none absolute left-0 top-full z-30 mt-1.5 w-max max-w-[190px] rounded border border-[#292929] bg-[#121212] px-2 py-1.5 text-[11px] leading-snug font-normal text-[#8a8a8a] opacity-0 shadow-2xl shadow-black/40 transition-opacity duration-200 ease-out group-hover/err:opacity-100"
												>
													{status?.error || 'unknown error'}
												</span>
											</span>
										{:else if typeof status?.took_ms === 'number'}
											<span class="shrink-0 text-[11px] text-[#5c5c5c]">{status.took_ms} ms</span>
										{/if}
									</span>
									<span class="relative h-3.5 w-6 shrink-0 rounded-full transition-colors {engines.includes(engine) ? 'bg-[#777]' : 'bg-[#292929]'}">
										<span class="absolute top-0.5 h-2.5 w-2.5 rounded-full bg-[#d0d0d0] transition-transform {engines.includes(engine) ? 'translate-x-3' : 'translate-x-0.5'}"></span>
									</span>
								</button>
							{/each}
						</div>
					</div>

					{#if active === 'Images'}
						<div class="space-y-2 border-t border-[#292929] pt-3">
							<button
								type="button"
								class="flex w-full items-center justify-between text-left text-[12px] text-[#9b9b9b]"
								onclick={() => {
									settings.hideAIImages = !settings.hideAIImages;
									saveSearchSettings();
								}}
							>
								<span>hide ai</span>
								<span class="relative h-3.5 w-6 rounded-full transition-colors {settings.hideAIImages ? 'bg-[#777]' : 'bg-[#292929]'}">
									<span class="absolute top-0.5 h-2.5 w-2.5 rounded-full bg-[#d0d0d0] transition-transform {settings.hideAIImages ? 'translate-x-3' : 'translate-x-0.5'}"></span>
								</span>
							</button>
							<button
								type="button"
								class="flex w-full items-center justify-between text-left text-[12px] text-[#9b9b9b]"
								onclick={() => {
									settings.transparentOnly = !settings.transparentOnly;
									saveSearchSettings();
								}}
							>
								<span>transparent only</span>
								<span class="relative h-3.5 w-6 rounded-full transition-colors {settings.transparentOnly ? 'bg-[#777]' : 'bg-[#292929]'}">
									<span class="absolute top-0.5 h-2.5 w-2.5 rounded-full bg-[#d0d0d0] transition-transform {settings.transparentOnly ? 'translate-x-3' : 'translate-x-0.5'}"></span>
								</span>
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
</nav>
