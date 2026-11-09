<script lang="ts">
	import AnchorHeading from './AnchorHeading.svelte';
	import SegmentedControl from './SegmentedControl.svelte';
	import {
		AUTOCOMPLETE_PROVIDERS,
		TIMEOUT_OPTIONS,
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

<div class="flex flex-col gap-7">
	<div class="flex flex-col gap-7">
		<section class="group">
			<AnchorHeading id="font">font</AnchorHeading>
			<p class="mt-1 text-[14px] leading-relaxed text-[#c0c0c0]">
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

		<section class="group">
			<AnchorHeading id="engine-badges">engine badges</AnchorHeading>
			<p class="mt-1 text-[14px] leading-relaxed text-[#c0c0c0]">
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

		<section class="group">
			<AnchorHeading id="result-stats">result stats</AnchorHeading>
			<p class="mt-1 text-[14px] leading-relaxed text-[#c0c0c0]">
				show how many results were found and how long the search took
			</p>
			<div class="mt-2.5">
				<SegmentedControl
					options={BOOLEAN_OPTIONS}
					bind:value={settings.showResultMeta}
					onchange={() => saveSearchSettings()}
				/>
			</div>
		</section>
	</div>

	<hr class="border-0 border-t border-base-700" />

	<div class="flex flex-col gap-7">
		<section class="group">
			<AnchorHeading id="autocomplete">autocomplete</AnchorHeading>
			<p class="mt-1 text-[14px] leading-relaxed text-[#c0c0c0]">
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

		<section class="group">
			<AnchorHeading id="timeout">timeout</AnchorHeading>
			<p class="mt-1 text-[14px] leading-relaxed text-[#c0c0c0]">
				how long to wait for search engines
			</p>
			<div class="mt-2.5">
				<SegmentedControl
					options={TIMEOUT_OPTIONS}
					bind:value={settings.timeout}
					onchange={() => saveSearchSettings()}
				/>
			</div>
		</section>

		<section class="group">
			<AnchorHeading id="open-links-in-new-tab">open links in new tab</AnchorHeading>
			<p class="mt-1 text-[14px] leading-relaxed text-[#c0c0c0]">
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

	<hr class="border-0 border-t border-base-700" />

	<div class="flex flex-col gap-7">
		<section class="group">
			<AnchorHeading id="analytics">analytics</AnchorHeading>
			<p class="mt-1 text-[14px] leading-relaxed text-[#c0c0c0]">
				participate in anonymous usage analytics
			</p>
			<div class="mt-2.5">
				<SegmentedControl
					options={BOOLEAN_OPTIONS}
					bind:value={settings.analytics}
					onchange={() => saveSearchSettings()}
				/>
			</div>
		</section>
	</div>
</div>
