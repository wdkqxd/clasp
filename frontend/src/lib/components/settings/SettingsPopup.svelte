<script lang="ts">
	import type { Snippet } from 'svelte';
	import settingsIcon from 'pixelarticons/svg/settings-cog.svg?raw';
	import chevronDownIcon from 'pixelarticons/svg/chevron-down.svg?raw';
	import { getSearchSettings, saveSearchSettings, type FontChoice } from './searchSettings.svelte';

	let {
		trigger = undefined,
		open = $bindable(false)
	}: { trigger?: Snippet; open?: boolean } = $props();

	const settings = getSearchSettings();

	const FONTS: { value: FontChoice; label: string }[] = [
		{ value: 'jetbrains', label: 'JetBrains Mono' },
		{ value: 'system', label: 'system mono' }
	];

	$effect(() => {
		if (!open) return;
		const closeOnOutside = (event: MouseEvent) => {
			if (event.target instanceof Element && !event.target.closest('[data-settings-popup]')) {
				open = false;
			}
		};
		const closeOnEscape = (event: KeyboardEvent) => {
			if (event.key === 'Escape') open = false;
		};
		document.addEventListener('click', closeOnOutside);
		document.addEventListener('keydown', closeOnEscape);
		return () => {
			document.removeEventListener('click', closeOnOutside);
			document.removeEventListener('keydown', closeOnEscape);
		};
	});
</script>

<div class="relative {trigger ? 'mx-auto w-fit mb-8' : ''}" data-settings-popup>
	{#if trigger}
		<button
			type="button"
			aria-label="Search settings"
			aria-expanded={open}
			onclick={() => (open = !open)}
			class="block cursor-pointer transition-transform duration-150 ease-[cubic-bezier(0.23,1,0.32,1)] active:scale-90 active:duration-100"
		>
			{@render trigger()}
		</button>
	{:else}
		<button
			type="button"
			aria-label="Search settings"
			aria-expanded={open}
			onclick={() => (open = !open)}
			class="inline-flex items-center gap-1.5 transition-colors text-[#5c5c5c] hover:text-[#8a8a8a] {open
				? 'text-[#9a9a9a]'
				: ''}"
		>
			<span class="icon">{@html settingsIcon}</span>
		</button>
	{/if}

	<div
		class="absolute right-0 top-full z-20 mt-2 w-64 origin-top-right rounded-md border border-[#292929] bg-[#121212] p-3 shadow-2xl shadow-black/30 transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] {open
			? 'visible translate-y-0 scale-100 opacity-100 blur-0'
			: 'pointer-events-none invisible -translate-y-1 scale-95 opacity-0 blur-sm'}"
	>
		<div class="space-y-2">
			<div class="flex items-center justify-between gap-3">
				<span class="text-[12px] text-[#9b9b9b]">font</span>
				<div class="relative flex items-center">
					<select
						value={settings.font}
						onchange={(event) => {
							settings.font = (event.currentTarget as HTMLSelectElement).value as FontChoice;
							saveSearchSettings();
						}}
						class="cursor-pointer appearance-none rounded bg-transparent py-0 pr-4 text-[12px] text-[#9b9b9b] outline-none transition-colors hover:text-[#d0d0d0]"
					>
						{#each FONTS as f}
							<option value={f.value}>{f.label}</option>
						{/each}
					</select>
					<span class="pointer-events-none absolute right-0 icon icon-sm text-[#5c5c5c]">{@html chevronDownIcon}</span>
				</div>
			</div>

			<button
				type="button"
				class="flex w-full items-center justify-between text-left text-[12px] text-[#9b9b9b]"
				onclick={() => {
					settings.showEngineBadges = !settings.showEngineBadges;
					saveSearchSettings();
				}}
			>
				<span>engine badges</span>
				<span class="relative h-3.5 w-6 rounded-full transition-colors {settings.showEngineBadges ? 'bg-[#777]' : 'bg-[#292929]'}">
					<span class="absolute top-0.5 h-2.5 w-2.5 rounded-full bg-[#d0d0d0] transition-transform duration-150 {settings.showEngineBadges ? 'translate-x-3' : 'translate-x-0.5'}"></span>
				</span>
			</button>

			<button
				type="button"
				class="flex w-full items-center justify-between text-left text-[12px] text-[#9b9b9b]"
				onclick={() => {
					settings.openInNewTab = !settings.openInNewTab;
					saveSearchSettings();
				}}
			>
				<span>open links in new tab</span>
				<span class="relative h-3.5 w-6 rounded-full transition-colors {settings.openInNewTab ? 'bg-[#777]' : 'bg-[#292929]'}">
					<span class="absolute top-0.5 h-2.5 w-2.5 rounded-full bg-[#d0d0d0] transition-transform duration-150 {settings.openInNewTab ? 'translate-x-3' : 'translate-x-0.5'}"></span>
				</span>
			</button>
		</div>
	</div>
</div>
