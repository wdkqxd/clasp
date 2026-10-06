<script lang="ts">
	import { ChevronDown } from '@lucide/svelte';
	import {
		AUTOCOMPLETE_PROVIDERS,
		getSearchSettings,
		saveSearchSettings,
		type AutocompleteProvider,
		type FontChoice
	} from './searchSettings.svelte';

	const settings = getSearchSettings();

	const FONTS: { value: FontChoice; label: string }[] = [
		{ value: 'jetbrains', label: 'JetBrains Mono' },
		{ value: 'system', label: 'system mono' }
	];
</script>

<div class="space-y-6">
	<div class="flex items-center justify-between gap-4">
		<div class="min-w-0">
			<div class="text-[14px] leading-snug text-[#c0c0c0]">font</div>
			<div class="mt-1 text-[12px] leading-relaxed text-[#5c5c5c]">
				typeface used across the interface
			</div>
		</div>
		<div class="relative flex shrink-0 items-center">
			<select
				value={settings.font}
				onchange={(event) => {
					settings.font = (event.currentTarget as HTMLSelectElement).value as FontChoice;
					saveSearchSettings();
				}}
				class="cursor-pointer appearance-none rounded bg-transparent py-0 pr-3.5 text-[13px] text-[#9b9b9b] outline-none transition-colors hover:text-[#d0d0d0]"
			>
				{#each FONTS as f}
					<option value={f.value}>{f.label}</option>
				{/each}
			</select>
			<ChevronDown class="pointer-events-none absolute right-0 size-3.5 text-[#5c5c5c]" />
		</div>
	</div>

	<div class="border-t border-[#1c1c1c]"></div>

	<div class="space-y-6">
		<div class="flex items-center justify-between gap-4">
			<div class="min-w-0">
				<div class="text-[14px] leading-snug text-[#c0c0c0]">autocomplete</div>
				<div class="mt-1 text-[12px] leading-relaxed text-[#5c5c5c]">
					which provider answers search suggestions
				</div>
			</div>
			<div class="relative flex shrink-0 items-center">
				<select
					value={settings.autocomplete}
					onchange={(event) => {
						settings.autocomplete = (event.currentTarget as HTMLSelectElement).value as AutocompleteProvider;
						saveSearchSettings();
					}}
					class="cursor-pointer appearance-none rounded bg-transparent py-0 pr-3.5 text-[13px] text-[#9b9b9b] outline-none transition-colors hover:text-[#d0d0d0]"
				>
					{#each AUTOCOMPLETE_PROVIDERS as provider}
						<option value={provider.value}>{provider.label}</option>
					{/each}
				</select>
				<ChevronDown class="pointer-events-none absolute right-0 size-3.5 text-[#5c5c5c]" />
			</div>
		</div>

		<button
			type="button"
			class="flex w-full cursor-pointer items-center justify-between gap-10 text-left"
			onclick={() => {
				settings.showEngineBadges = !settings.showEngineBadges;
				saveSearchSettings();
			}}
		>
			<span class="min-w-0">
				<span class="block text-[14px] leading-snug text-[#c0c0c0]">engine badges</span>
				<span class="mt-1 block text-[12px] leading-relaxed text-[#5c5c5c]">
					show which engine returned each result
				</span>
			</span>
			<span
				class="relative h-5 w-9 shrink-0 rounded-full transition-colors {settings.showEngineBadges
					? 'bg-[#777]'
					: 'bg-[#292929]'}"
			>
				<span
					class="absolute top-0.5 h-4 w-4 rounded-full bg-[#d0d0d0] transition-transform duration-150 {settings.showEngineBadges
						? 'translate-x-4.5'
						: 'translate-x-0.5'}"
				></span>
			</span>
		</button>
	</div>

	<div class="border-t border-[#1c1c1c]"></div>

	<div>
		<button
			type="button"
			class="flex w-full cursor-pointer items-center justify-between gap-10 text-left"
			onclick={() => {
				settings.openInNewTab = !settings.openInNewTab;
				saveSearchSettings();
			}}
		>
			<span class="min-w-0">
				<span class="block text-[14px] leading-snug text-[#c0c0c0]">open links in new tab</span>
				<span class="mt-1 block text-[12px] leading-relaxed text-[#5c5c5c]">
					open results and pages in a separate tab
				</span>
			</span>
			<span
				class="relative h-5 w-9 shrink-0 rounded-full transition-colors {settings.openInNewTab
					? 'bg-[#777]'
					: 'bg-[#292929]'}"
			>
				<span
					class="absolute top-0.5 h-4 w-4 rounded-full bg-[#d0d0d0] transition-transform duration-150 {settings.openInNewTab
						? 'translate-x-4.5'
						: 'translate-x-0.5'}"
				></span>
			</span>
		</button>
	</div>
</div>