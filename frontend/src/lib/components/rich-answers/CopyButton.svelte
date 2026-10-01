<script lang="ts">
	import { Check, Copy } from '@lucide/svelte';

	let { active, size = 'h-4 w-4' }: { active: boolean; size?: string } = $props();
</script>

<span class="icon shrink-0 text-[#5c5c5c] opacity-0 transition-opacity group-hover:opacity-70 {size}">
	{#if active}
		<Check />
	{:else}
		<Copy />
	{/if}
</span>