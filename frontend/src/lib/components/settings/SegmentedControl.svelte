<script lang="ts" generics="T extends string | boolean">
	import { tick } from 'svelte';

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
	let activeLeft = $state(0);
	let activeWidth = $state(0);
	let ready = $state(false);

	function updateIndicator() {
		if (!containerEl) return;
		const activeBtn = containerEl.querySelector<HTMLButtonElement>('[data-active="true"]');
		if (!activeBtn) return;

		activeLeft = activeBtn.offsetLeft;
		activeWidth = activeBtn.offsetWidth;
		if (!ready && activeWidth > 0) {
			ready = true;
		}
	}

	$effect(() => {
		const _ = value;
		const __ = options;
		tick().then(() => {
			updateIndicator();
		});
	});

	$effect(() => {
		if (!containerEl) return;
		const ro = new ResizeObserver(() => {
			updateIndicator();
		});
		ro.observe(containerEl);
		return () => ro.disconnect();
	});

	function select(val: T) {
		if (value === val) return;
		value = val;
		onchange?.(val);
	}
</script>

<div
	bind:this={containerEl}
	class="relative inline-flex max-w-full items-center rounded-[5px] bg-[#141414] p-[3px]"
>
	<span
		class="pointer-events-none absolute left-0 top-[3px] bottom-[3px] rounded-[3px] border border-[#2f2f2f] bg-[#222222] shadow-xs {ready
			? 'transition-[transform,width] duration-200 ease-[cubic-bezier(0.22,1,0.36,1)]'
			: ''}"
		style:transform="translateX({activeLeft}px)"
		style:width="{activeWidth}px"
		style:opacity={ready && activeWidth > 0 ? 1 : 0}
	></span>

	{#each options as opt (opt.value)}
		{@const isActive = opt.value === value}
		<button
			type="button"
			data-active={isActive}
			onclick={() => select(opt.value)}
			class="relative z-10 min-w-[38px] cursor-pointer rounded-[3px] px-2.5 py-1 text-center text-[12px] leading-tight select-none transition-colors duration-150 {isActive
				? 'text-[#e0e0e0]'
				: 'text-[#6e6e6e] hover:text-[#a0a0a0]'}"
		>
			{opt.label}
		</button>
	{/each}
</div>
