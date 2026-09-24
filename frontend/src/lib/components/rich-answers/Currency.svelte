<script lang="ts">
	import type { RichCurrency } from '../../types';
	import Sparkline from './Sparkline.svelte';
	import { price, price2 } from './format';
	import CopyButton from './CopyButton.svelte';

	let { currency }: { currency: RichCurrency } = $props();

	let copied = $state(false);
	function onCopy() {
		navigator.clipboard?.writeText(String(currency.result));
		copied = true;
		setTimeout(() => (copied = false), 1200);
	}
</script>

<div class="overflow-hidden rounded-md border border-[#292929] bg-[#121212] shadow-2xl shadow-black/30">
	<div class="p-4 pb-3">
		<p class="text-[16px] text-[#c0c0c0]">
			<span class="text-[#707070]">{price(currency.amount)} {currency.from_currency_name || currency.from_currency_code} =</span>
			<button
				type="button"
				class="group inline-flex cursor-pointer items-center gap-3 text-[#e0e0e0]"
				title="copy result"
				onclick={onCopy}
			>
				{price2(currency.result)} {currency.to_currency_name || currency.to_currency_code}
				<CopyButton active={copied} />
			</button>
		</p>
		{#if currency.exchange_rate || currency.date}
			<div class="mt-2 text-[11px] text-[#5c5c5c]">
				{#if currency.exchange_rate}<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 tracking-wide">1 {currency.from_currency_code} = {currency.exchange_rate} {currency.to_currency_code}</span>{/if}
				{#if currency.exchange_rate && currency.date}<span class="mx-1">·</span>{/if}
				{#if currency.date}{currency.date}{/if}
			</div>
		{/if}
	</div>
	{#if currency.timeseries?.ts_exchange_rate?.length}
		{@const pts = currency.timeseries.ts_exchange_rate.map((p) => [new Date(p[0]).getTime(), p[1]] as [number, number])}
		{#if pts.length > 1}
			<div class="-mx-0"><Sparkline points={pts} height={44} /></div>
		{/if}
	{/if}
</div>