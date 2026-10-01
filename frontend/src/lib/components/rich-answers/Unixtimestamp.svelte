<script lang="ts">
	import CopyButton from './CopyButton.svelte';

	let { query = '' }: { query?: string } = $props();

	let unixNow = $state(Math.floor(Date.now() / 1000));
	$effect(() => {
		const id = setInterval(() => (unixNow = Math.floor(Date.now() / 1000)), 1000);
		return () => clearInterval(id);
	});

	const NOW_RE =
		/^(?:current\s+)?(?:unix\s+)?timestamp(?:\s+(?:now|in\s+seconds))?$|^(?:what(?:'s| is)\s+)?(?:the\s+)?unix\s+(?:time\s+)?now$/i;
	const EPOCH_PREFIX_RE = /^(?:convert\s+)?(?:unix\s+)?timestamp(?:\s+(?:to\s+\w+\s+)?(?:for|of|in|at))?\s+(\d{1,13})$/i;
	const EPOCH_SUFFIX_RE = /^(?:convert\s+)?(\d{9,13})\s*(?:unix(?:\s+time(?:stamp)?)?|epoch(?:\s+time(?:stamp)?)?|to\s+(?:a\s+)?date|time)$/i;
	const DATE_RE =
		/^(?:convert\s+)?(\d{4}-\d{2}-\d{2}(?:[ t]\d{1,2}:\d{2}(?::\d{2})?)?)\s*(?:to|into|in)\s+(?:an?\s+)?unix(?:\s+time(?:stamp)?)?$/i;

	const parsed = $derived.by(() => {
		const q = query.trim().replace(/\s+/g, ' ');
		if (!q) return null;

		if (NOW_RE.test(q)) return { epoch: null, date: null };

		let m = q.match(EPOCH_PREFIX_RE);
		if (m) return fromEpoch(Number(m[1]));

		m = q.match(EPOCH_SUFFIX_RE);
		if (m) return fromEpoch(Number(m[1]));

		m = q.match(DATE_RE);
		if (m) {
			const ms = Date.parse(m[1].replace(' ', 'T'));
			if (!Number.isNaN(ms)) {
				const epoch = Math.floor(ms / 1000);
				return { epoch, date: formatDate(epoch) };
			}
			return null;
		}
		return null;
	});

	function fromEpoch(value: number) {
		if (!Number.isFinite(value) || value < 0) return null;
		return { epoch: value, date: formatDate(value) };
	}

	function formatDate(epoch: number) {
		return new Date(epoch * 1000).toLocaleString(undefined, {
			weekday: 'short',
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit',
			second: '2-digit',
			timeZoneName: 'short'
		});
	}

	const epoch = $derived(parsed?.epoch ?? unixNow);
	const human = $derived(parsed?.date ?? formatDate(unixNow));
	const isNow = $derived(parsed?.epoch === null);

	let copied = $state(false);
	function onCopy() {
		navigator.clipboard?.writeText(String(epoch));
		copied = true;
		setTimeout(() => (copied = false), 1200);
	}
</script>

{#if parsed}
	<div class="rounded-md border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
		<button
			type="button"
			class="group inline-flex cursor-pointer items-center gap-3 font-mono text-[32px] leading-none text-[#e0e0e0]"
			title="copy timestamp"
			onclick={onCopy}
		>
			{epoch}
			<CopyButton active={copied} />
		</button>
		<div class="mt-2 text-[13px] text-[#707070]">{human}</div>
		{#if isNow}
			<div class="mt-1 text-[12px] text-[#505050]">current unix time</div>
		{/if}
	</div>
{/if}