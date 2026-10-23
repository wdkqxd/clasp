<script lang="ts">
	import SegmentedControl from './SegmentedControl.svelte';
	import {
		AUTOCOMPLETE_PROVIDERS,
		getSearchSettings,
		saveSearchSettings,
		type FontChoice
	} from './searchSettings.svelte';

	const settings = getSearchSettings();

	const FONTS: { value: FontChoice; label: string }[] = [
		{ value: 'jetbrains', label: 'jetbrains mono' },
		{ value: 'system', label: 'system mono' }
	];

	const BOOLEAN_OPTIONS = [
		{ value: false, label: 'off' },
		{ value: true, label: 'on' }
	];
</script>

<div>
	<section>
		<h2 class="text-[15px] font-normal leading-[1.375] text-[#c0c0c0]">font</h2>
		<p class="mt-1 text-[13px] leading-[1.625] text-[#707070]">
			typeface used across the interface
		</p>
		<div class="mt-2.5">
			<SegmentedControl
				options={FONTS}
				bind:value={settings.font}
				onchange={() => saveSearchSettings()}
			/>
		</div>
	</section>

	<section class="mt-7">
		<h2 class="text-[15px] font-normal leading-[1.375] text-[#c0c0c0]">autocomplete</h2>
		<p class="mt-1 text-[13px] leading-[1.625] text-[#707070]">
			which provider answers search suggestions
		</p>
		<div class="mt-2.5">
			<SegmentedControl
				options={AUTOCOMPLETE_PROVIDERS.map((p) => ({ value: p.value, label: p.label }))}
				bind:value={settings.autocomplete}
				onchange={() => saveSearchSettings()}
			/>
		</div>
	</section>

	<section class="mt-7">
		<h2 class="text-[15px] font-normal leading-[1.375] text-[#c0c0c0]">engine badges</h2>
		<p class="mt-1 text-[13px] leading-[1.625] text-[#707070]">
			show which engine returned each result
		</p>
		<div class="mt-2.5">
			<SegmentedControl
				options={BOOLEAN_OPTIONS}
				bind:value={settings.showEngineBadges}
				onchange={() => saveSearchSettings()}
			/>
		</div>
	</section>

	<section class="mt-7">
		<h2 class="text-[15px] font-normal leading-[1.375] text-[#c0c0c0]">open links in new tab</h2>
		<p class="mt-1 text-[13px] leading-[1.625] text-[#707070]">
			open results and pages in a separate tab
		</p>
		<div class="mt-2.5">
			<SegmentedControl
				options={BOOLEAN_OPTIONS}
				bind:value={settings.openInNewTab}
				onchange={() => saveSearchSettings()}
			/>
		</div>
	</section>
</div>
