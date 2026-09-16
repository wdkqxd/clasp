<script lang="ts">
	import arrowsIcon from 'pixelarticons/svg/arrows-horizontal.svg?raw';
	import MarketCard from '../market/MarketCard.svelte';
	import type { QuickAnswer } from './quickAnswer';
	let { answer }: { answer: QuickAnswer } = $props();
	let swapped = $state(false);

	function numericValue(value: string, unit?: string): string {
		return unit ? value.replace(new RegExp(`\\s${unit}$`, 'i'), '') : value;
	}

	function unitLabel(unit?: string): string {
		const labels: Record<string, string> = {
			mm: 'millimeters', cm: 'centimeters', m: 'meters', km: 'kilometers', in: 'inches', ft: 'feet',
			yd: 'yards', mi: 'miles', mg: 'milligrams', g: 'grams', kg: 'kilograms', t: 'tons', oz: 'ounces',
			lb: 'pounds', ml: 'milliliters', cl: 'centiliters', dl: 'deciliters', l: 'liters', gal: 'gallons',
			'm/s': 'meters per second', 'km/h': 'kilometers per hour', mph: 'miles per hour',
			c: 'degrees Celsius', f: 'degrees Fahrenheit', k: 'kelvin'
		};
		return unit ? labels[unit] ?? unit : '';
	}
</script>

<section class="rounded-md border border-[#262626] px-4 py-3">
	{#if answer.type === 'market' && answer.market}
		<MarketCard quote={answer.market} />
	{:else if (answer.fromCurrency && answer.toCurrency) || (answer.fromUnit && answer.toUnit)}
		<div class="flex items-center gap-2">
			<div class="min-w-0 flex-1 text-center">
				<div class="text-[18px] text-[#c8c8c8]">{swapped ? numericValue(answer.value, answer.toUnit) : answer.fromUnit ? answer.amount : `${answer.amount} ${answer.fromCurrency}`}</div>
				<div class="mt-1 text-[11px] text-[#5c5c5c]">{swapped ? answer.toCurrency ?? unitLabel(answer.toUnit) : answer.fromCurrency ?? unitLabel(answer.fromUnit)}</div>
			</div>
			<button
				type="button"
				aria-label="Swap currencies"
				onclick={() => (swapped = !swapped)}
				class="icon !h-5 !w-5 shrink-0 p-0 text-[#5c5c5c] transition-all duration-100 hover:text-[#a0a0a0] active:scale-90"
			>
				{@html arrowsIcon}
			</button>
			<div class="min-w-0 flex-1 text-center">
				<div class="text-[18px] text-[#c8c8c8]">{swapped ? (answer.fromUnit ? answer.amount : `${answer.amount} ${answer.fromCurrency}`) : answer.fromUnit ? numericValue(answer.value, answer.toUnit) : answer.value}</div>
				<div class="mt-1 text-[11px] text-[#5c5c5c]">{swapped ? answer.fromCurrency ?? unitLabel(answer.fromUnit) : answer.toCurrency ?? unitLabel(answer.toUnit)}</div>
			</div>
		</div>
	{:else if answer.type === 'color'}
		<div class="flex items-center gap-3"><span class="h-8 w-8 rounded border border-[#333]" style:background={answer.detail}></span><span class="text-[18px] text-[#c8c8c8]">{answer.value}</span></div>
	{:else}
		<div class="mb-2 text-[10px] lowercase tracking-wide text-[#5c5c5c]">{answer.title}</div>
		{#if answer.detail}<div class="text-[12px] text-[#5c5c5c]">{answer.detail}</div>{/if}
		<div class="text-[21px] leading-tight text-[#d0d0d0]">{answer.value}</div>
	{/if}
</section>
