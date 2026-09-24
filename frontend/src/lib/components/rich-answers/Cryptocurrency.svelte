<script lang="ts">
	import type { RichCryptocurrency } from '../../types';
	import Sparkline from './Sparkline.svelte';
	import { sym, price, price2 } from './format';

	let { crypto }: { crypto: RichCryptocurrency } = $props();

	const q = crypto.quote;
	const s = sym(crypto.vs_currency ?? 'usd');
</script>

{#if q}
	<div class="overflow-hidden rounded-md border border-[#292929] bg-[#121212] shadow-2xl shadow-black/30">
		<div class="p-4 pb-3">
			<div class="flex items-center gap-1.5">
				<span class="text-[12px] tracking-wide text-[#5c5c5c]">{q.name}</span>
				<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 text-[10.5px] tracking-wide text-[#5c5c5c]">
					{q.symbol?.toUpperCase()}
				</span>
			</div>
			<div class="mt-3 flex items-baseline gap-3">
				<span class="text-[30px] leading-none text-[#e0e0e0]">{s}{price(q.current_price)}</span>
				{#if q.price_change_percentage_24h != null}
					<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 text-[10.5px] {q.price_change_percentage_24h >= 0 ? 'text-[#a6e3a1]' : 'text-[#f38ba8]'}">
						{q.price_change_percentage_24h >= 0 ? '+' : ''}{q.price_change_percentage_24h.toFixed(2)}%
					</span>
				{/if}
			</div>
			{#if q.high_24h && q.low_24h}
				<div class="mt-1 text-[11px] text-[#5c5c5c]">24h {s}{price2(q.low_24h)} – {s}{price2(q.high_24h)}</div>
			{/if}
		</div>
		{#if crypto.timeseries?.ts_price && crypto.timeseries.ts_price.length > 1}
			<div><Sparkline points={crypto.timeseries.ts_price} height={44} /></div>
		{/if}
	</div>
{/if}