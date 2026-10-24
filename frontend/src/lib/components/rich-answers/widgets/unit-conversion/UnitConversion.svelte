<script lang="ts">
	import { convertUnit, parseConversion } from './units';
	import CopyButton from '../../shared/CopyButton.svelte';

	let { query = '' }: { query?: string } = $props();

	const data = $derived(parseConversion(query.trim()));

	const out = $derived(data ? convertUnit(data.amount, data.from_unit, data.to_unit) : '');
	const outN = $derived(Number.parseFloat(out.replace(/,/g, '')) || 0);

	let copied = $state(false);
	function onCopy() {
		navigator.clipboard?.writeText(out);
		copied = true;
		setTimeout(() => (copied = false), 1200);
	}
	const SYMBOLS = new Set([
		'm', 'km', 'cm', 'mm', 'mi', 'ft', 'in', 'yd', 'l', 'ml', 'g', 'kg', 'mg', 't',
		'lb', 'oz', 'h', 'min', 's', 'd', 'b', 'kb', 'mb', 'gb', 'tb', 'mph', 'kph',
		'km/h', 'mi/h', '°c', '°f'
	]);

	function plural(u: string, n: number): string {
		if (Math.abs(n) === 1 || /s$/i.test(u) || SYMBOLS.has(u.toLowerCase())) return u;
		return `${u}s`;
	}
</script>

{#if data}
	{@const u = data}
	<div class="rounded-[5px] border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
		<p class="text-[16px] text-[#c0c0c0]">
			<span class="text-[#707070]">{u.amount} {plural(u.from_label, u.amount)} =</span>
			<button
				type="button"
				class="group inline-flex cursor-pointer items-center gap-3 text-[#e0e0e0]"
				title="copy result"
				onclick={onCopy}
			>
				{out} {plural(u.to_label, outN)}
				<CopyButton active={copied} />
			</button>
		</p>
	</div>
{/if}