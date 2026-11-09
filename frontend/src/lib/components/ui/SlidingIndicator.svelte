<script lang="ts">
	import { tick } from 'svelte';

	let {
		container,
		sync,
		extraClass = ''
	}: {
		container: HTMLElement | null;
		sync: unknown;
		extraClass?: string;
	} = $props();

	let left = $state(0);
	let width = $state(0);
	let ready = $state(false);

	function measure() {
		if (!container) return;
		const activeEl = container.querySelector<HTMLElement>('[data-active="true"]');
		if (!activeEl || activeEl.offsetWidth <= 0) return;

		left = activeEl.offsetLeft;
		width = activeEl.offsetWidth;

		if (!ready) requestAnimationFrame(() => (ready = true));
	}

	$effect(() => {
		void container;
		void sync;
		tick().then(measure);
	});

	$effect(() => {
		if (!container) return;
		const ro = new ResizeObserver(measure);
		ro.observe(container);
		return () => ro.disconnect();
	});
</script>

<span
	class="pointer-events-none absolute left-0 top-[3px] bottom-[3px] rounded-[3px] border border-[#2f2f2f] bg-[#222222] {extraClass} {ready
		? 'transition-[transform,width] duration-200 ease-[cubic-bezier(0.22,1,0.36,1)]'
		: ''}"
	style:transform="translateX({left}px)"
	style:width="{width}px"
	style:opacity={ready && width > 0 ? 1 : 0}
></span>
