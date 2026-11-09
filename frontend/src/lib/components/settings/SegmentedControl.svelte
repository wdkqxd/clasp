<script lang="ts" generics="T extends string | number | boolean">
	import SlidingIndicator from '$lib/components/ui/SlidingIndicator.svelte';
	import { pillGroupClass } from '$lib/components/ui/pill';

	let {
		options,
		value = $bindable(),
		onchange
	}: {
		options: { value: T; label: string }[];
		value: T;
		onchange?: (value: T) => void;
	} = $props();

	let containerEl = $state<HTMLDivElement | null>(null);

	function select(val: T) {
		if (value === val) return;
		value = val;
		onchange?.(val);
	}
</script>

<div bind:this={containerEl} class={pillGroupClass(options.length)}>
	<SlidingIndicator
		container={containerEl}
		sync={`${value}|${options.map((option) => option.label).join(',')}`}
		extraClass="shadow-xs"
	/>

	{#each options as opt (opt.value)}
		{@const isActive = opt.value === value}
		<button
			type="button"
			data-active={isActive}
			onclick={() => select(opt.value)}
			class="relative z-10 min-w-[40px] cursor-pointer rounded-[3px] px-2.5 py-1 text-center text-[13px] leading-tight select-none transition-colors duration-150 ease-out {isActive
				? 'text-[#c0c0c0]'
				: 'text-[#6a6a6a] hover:text-[#9a9a9a]'}"
		>
			{opt.label}
		</button>
	{/each}
</div>
