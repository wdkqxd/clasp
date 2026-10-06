<script lang="ts">
	import { Check, Copy } from '@lucide/svelte';

	let { active, size = 16 }: { active: boolean; size?: number } = $props();

	const cls = 'shrink-0 text-[#5c5c5c] opacity-0 transition-opacity group-hover:opacity-70';
</script>

{#if active}
	<Check {size} class={cls} />
{:else}
	<Copy {size} class={cls} />
{/if}