<script lang="ts">
	import type { RichUnitConversion } from '../../types';
	import { convertUnit, formatUnit } from './units';
	import CopyButton from './CopyButton.svelte';

	let { unit }: { unit: RichUnitConversion } = $props();

	const out = convertUnit(unit.amount, unit.from_unit, unit.to_unit, unit.dimensionality);
	const outN = parseFloat(out.replace(/,/g, ''));

	let copied = $state(false);
	function onCopy() {
		navigator.clipboard?.writeText(out);
		copied = true;
		setTimeout(() => (copied = false), 1200);
	}
	function plural(u: string, n: number): string {
		return Math.abs(n) === 1 ? u : `${u}s`;
	}
</script>

<div class="rounded-md border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
	<p class="text-[16px] text-[#c0c0c0]">
		<span class="text-[#707070]">{unit.amount} {plural(formatUnit(unit.from_unit), unit.amount)} =</span>
		<button
			type="button"
			class="group inline-flex cursor-pointer items-center gap-3 text-[#e0e0e0]"
			title="copy result"
			onclick={onCopy}
		>
			{out} {plural(formatUnit(unit.to_unit), outN)}
			<CopyButton active={copied} />
		</button>
	</p>
</div>