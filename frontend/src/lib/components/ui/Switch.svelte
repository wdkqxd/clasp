<script lang="ts">
	let { checked }: { checked: boolean } = $props();
</script>

<span
	class="relative h-3.5 w-6 shrink-0 rounded-full transition-colors {checked
		? 'bg-[#777]'
		: 'bg-[#292929]'}"
>
	<span
		class="absolute top-0.5 h-2.5 w-2.5 rounded-full bg-[#d0d0d0] transition-transform {checked
			? 'translate-x-3'
			: 'translate-x-0.5'}"
	></span>
</span>
