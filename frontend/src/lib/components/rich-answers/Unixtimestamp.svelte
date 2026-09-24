<script lang="ts">
	import CopyButton from './CopyButton.svelte';

	let unixNow = $state(Math.floor(Date.now() / 1000));
	$effect(() => {
		const id = setInterval(() => (unixNow = Math.floor(Date.now() / 1000)), 1000);
		return () => clearInterval(id);
	});

	let copied = $state(false);
	function onCopy() {
		navigator.clipboard?.writeText(String(unixNow));
		copied = true;
		setTimeout(() => (copied = false), 1200);
	}

	let unixHuman = $derived(
		new Date(unixNow * 1000).toLocaleString().replace(/\s?([AP])M\b/, (_, x) => ` ${x.toLowerCase()}m`)
	);
</script>

<div class="rounded-md border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
	<button
		type="button"
		class="group inline-flex cursor-pointer items-center gap-3 font-mono text-[32px] leading-none text-[#e0e0e0]"
		title="copy timestamp"
		onclick={onCopy}
	>
		{unixNow}
		<CopyButton active={copied} />
	</button>
	<div class="mt-2 text-[13px] text-[#707070]">{unixHuman}</div>
</div>
