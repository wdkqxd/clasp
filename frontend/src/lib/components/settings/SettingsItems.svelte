<script lang="ts">
	import { page } from '$app/state';
	import { isOfficialHost } from '$lib/instance';
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
	const isOfficialInstance = $derived(isOfficialHost(page.url.hostname));

	type Group = {
		id: string;
		label: string;
		sections: string[];
	};

	const groups: Group[] = $derived([
		{ id: 'appearance', label: 'appearance', sections: ['font', 'engine-badges'] },
		{
			id: 'search',
			label: 'search',
			sections: ['autocomplete', 'timeout', 'open-links-in-new-tab']
		},
		...(isOfficialInstance
			? []
			: [
					{
						id: 'privacy',
						label: 'privacy',
						sections: ['unofficial-instance-notice']
					}
				])
	]);

	let active = $state('appearance');

	$effect(() => {
		const hash = decodeURIComponent(page.url.hash.slice(1));
		const group = groups.find((g) => g.sections.includes(hash));
		if (group) active = group.id;
	});

	const FONTS: { value: FontChoice; label: string }[] = [
		{ value: 'jetbrains', label: 'jetbrains mono' },
		{ value: 'system', label: 'system mono' }
	];

	const BOOLEAN_OPTIONS = [
		{ value: false, label: 'off' },
		{ value: true, label: 'on' }
	];
</script>

<div class="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-center lg:gap-6">
	<div
		class="flex gap-1 overflow-x-auto pb-1 lg:w-44 lg:shrink-0 lg:flex-col lg:gap-0 lg:overflow-visible lg:pb-0"
		role="tablist"
		aria-label="settings sections"
	>
		{#each groups as group (group.id)}
			<button
				type="button"
				role="tab"
				id="tab-{group.id}"
				aria-selected={active === group.id}
				aria-controls="panel-{group.id}"
				onclick={() => (active = group.id)}
				class="cursor-pointer rounded-[5px] px-2.5 py-1.5 text-left whitespace-nowrap text-[13px] transition-all duration-150 ease-out select-none active:scale-[0.94] active:duration-75 {active ===
				group.id
					? 'bg-[#262626] text-[#e0e0e0]'
					: 'text-[#707070] hover:bg-[#1a1a1a] hover:text-[#c0c0c0]'}"
			>
				{group.label}
			</button>
		{/each}
	</div>

	<div class="lg:w-160 lg:shrink-0">
		<div
			id="panel-appearance"
			role="tabpanel"
			aria-labelledby="tab-appearance"
			class:hidden={active !== 'appearance'}
		>
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

			<section class="group mt-7">
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
		</div>

		<div
			id="panel-search"
			role="tabpanel"
			aria-labelledby="tab-search"
			class:hidden={active !== 'search'}
		>
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

			<section class="group mt-7">
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

			<section class="group mt-7">
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

		{#if !isOfficialInstance}
			<div
				id="panel-privacy"
				role="tabpanel"
				aria-labelledby="tab-privacy"
				class:hidden={active !== 'privacy'}
			>
				<section class="group">
					<AnchorHeading id="unofficial-instance-notice">unofficial instance notice</AnchorHeading>
					<p class="mt-1 text-[14px] leading-relaxed text-[#c0c0c0]">
						show an unofficial instance warning on the homepage
					</p>
					<div class="mt-2.5">
						<SegmentedControl
							options={BOOLEAN_OPTIONS}
							bind:value={settings.showInstanceNotice}
							onchange={() => saveSearchSettings()}
						/>
					</div>
				</section>
			</div>
		{/if}
	</div>
</div>
