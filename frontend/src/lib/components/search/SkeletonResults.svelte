<script lang="ts">
	let { count = 6 }: { count?: number } = $props();
</script>

<div class="space-y-10">
	{#each Array(count) as _, i (i)}
		<article class="result group">
			<div class="mb-1.5"><div class="skeleton" style="width:140px;height:12px"></div></div>
			<div class="mb-1.5"><div class="skeleton" style="width:70%;height:15px"></div></div>
			<div class="mb-1"><div class="skeleton" style="height:13px"></div></div>
			<div class="mb-1"><div class="skeleton" style="width:92%;height:13px"></div></div>
			<div><div class="skeleton" style="width:60%;height:13px"></div></div>
		</article>
	{/each}
</div>
