<script lang="ts">
	import { Link as LinkIcon } from '@lucide/svelte';
	import type { Snippet } from 'svelte';

	let {
		id,
		children
	}: {
		id: string;
		children: Snippet;
	} = $props();
</script>

<h2
	{id}
	class="flex w-fit scroll-mt-6 items-center gap-1.5 text-[16px] leading-[1.375] font-semibold text-[#c0c0c0]"
>
	{@render children()}
	<a
		href="#{id}"
		aria-label="link to this section"
		class="inline-flex text-[#707070] opacity-0 transition-[color,opacity] duration-150 hover:text-[#c0c0c0] focus-visible:opacity-100 group-hover:opacity-100"
	>
		<LinkIcon size={14} />
	</a>
</h2>
