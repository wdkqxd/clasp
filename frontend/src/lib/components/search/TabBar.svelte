<script lang="ts">
	import { page } from '$app/state';
	import { Clapperboard, Globe, Image, Settings2 } from '@lucide/svelte';
	import { fade } from 'svelte/transition';
	import {
		ENGINE_API_NAMES,
		getSearchSettings,
		SEARCH_ENGINES,
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
	const settings = getSearchSettings();

	const icons = {
		Web: Globe,
		Images: Image,
		Video: Clapperboard
	} as const;

	function tabHref(tab: string): string {
		const q = page.url.searchParams.get('q') ?? '';
		const params = new URLSearchParams({ q });
		if (tab === 'Images') params.set('type', 'images');
		return `/search?${params.toString()}`;
	}

	function toggleEngine(engine: SearchEngine) {
		if (settings.engines.includes(engine) && settings.engines.length === 1) return;
		settings.engines = settings.engines.includes(engine)
			? settings.engines.filter((item) => item !== engine)
			: [...settings.engines, engine];
		saveSearchSettings();
	}

	function engineStatusFor(engine: SearchEngine): EngineStatus | undefined {
		return engineStatus[ENGINE_API_NAMES[engine]];
	}

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
	<div class="flex items-center gap-5">
		{#each tabs as tab}
			{@const Icon = icons[tab as keyof typeof icons] ?? Globe}
			<a
				href={tab === 'Video' ? '#' : tabHref(tab)}
				class="inline-flex items-center gap-1.5 {tab.toLowerCase() === active.toLowerCase()
					? 'text-[#8a8a8a]'
					: 'text-[#4a4a4a] hover:text-[#6a6a6a]'} transition-colors"
				onclick={tab === 'Video' ? (e) => e.preventDefault() : undefined}
			>
				<Icon class="size-4" />
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
				class="absolute right-0 top-7 z-20 w-56 origin-top-right rounded-md border border-[#292929] bg-[#121212] p-3 shadow-2xl shadow-black/30 transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] {settingsOpen
					? 'visible translate-y-0 scale-100 opacity-100 blur-0'
					: 'pointer-events-none invisible -translate-y-1 scale-95 opacity-0 blur-sm'}"
			>
				<div class="space-y-3">
					<div>
						<div class="space-y-2">
							{#each SEARCH_ENGINES as engine}
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
									<span class="relative h-3.5 w-6 shrink-0 rounded-full transition-colors {settings.engines.includes(engine) ? 'bg-[#777]' : 'bg-[#292929]'}">
										<span class="absolute top-0.5 h-2.5 w-2.5 rounded-full bg-[#d0d0d0] transition-transform {settings.engines.includes(engine) ? 'translate-x-3' : 'translate-x-0.5'}"></span>
									</span>
								</button>
							{/each}
						</div>
					</div>

					{#if active === 'Images'}
						<div class="border-t border-[#292929] pt-3">
							<button
								type="button"
								class="flex w-full items-center justify-between text-left text-[12px] text-[#9b9b9b]"
								onclick={() => {
									settings.hideAIImages = !settings.hideAIImages;
									saveSearchSettings();
								}}
							>
								<span>hide ai images</span>
								<span class="relative h-3.5 w-6 rounded-full transition-colors {settings.hideAIImages ? 'bg-[#777]' : 'bg-[#292929]'}">
									<span class="absolute top-0.5 h-2.5 w-2.5 rounded-full bg-[#d0d0d0] transition-transform {settings.hideAIImages ? 'translate-x-3' : 'translate-x-0.5'}"></span>
								</span>
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>
	</div>
</nav>
