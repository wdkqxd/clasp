<script lang="ts">
	import { goto } from '$app/navigation';
	import { X } from '@lucide/svelte';

	function closePage() {
		void goto('/');
	}
</script>

<div class="flex w-full shrink-0 items-center bg-[#0c0c0c] px-3 py-1.5">
	<button
		type="button"
		onclick={closePage}
		aria-label="close"
		title="close"
		class="inline-flex size-7 cursor-pointer items-center justify-center rounded-[5px] text-[#707070] transition-all duration-150 ease-out hover:bg-[#1a1a1a] hover:text-[#c0c0c0] active:scale-[0.92] active:duration-75 focus-visible:outline-1 focus-visible:outline-offset-1 focus-visible:outline-[#707070]"
	>
		<X size={16} strokeWidth={2} />
	</button>
</div>
