<script lang="ts">
	import type { RichAnswer } from '../../types';
	import Sparkline from './Sparkline.svelte';
	import { sym, price, price2 } from './format';

	let { answer }: { answer: RichAnswer } = $props();

	const v = $derived.by(() => {
		const c = answer.cryptocurrency;
		if (c?.quote) {
			return {
				name: c.quote.name,
				symbol: c.quote.symbol,
				badge: '',
				currency: c.vs_currency,
				price: c.quote.current_price,
				changePercent: c.quote.price_change_percentage_24h,
				change: null as number | null,
				low: c.quote.low_24h,
				high: c.quote.high_24h,
				rangeLabel: '24h',
				points: c.timeseries?.ts_price
			};
		}
		const st = answer.stock;
		if (st?.quote) {
			const q = st.quote;
			return {
				name: q.company_name || st.asset_info?.name || '',
				symbol: q.symbol,
				badge: q.primary_exchange || '',
				currency: st.asset_info?.currency || q.currency,
				price: q.latest_price,
				changePercent: q.change_percent,
				change: q.change ?? null,
				low: q.low,
				high: q.high,
				rangeLabel: 'Day',
				points: st.timeseries?.ts_close
			};
		}
		return null;
	});

	const s = $derived(sym(v?.currency ?? ''));
	const up = $derived((v?.changePercent ?? 0) >= 0);
</script>

{#if v}
	<div class="overflow-hidden rounded-md border border-[#292929] bg-[#121212] shadow-2xl shadow-black/30">
		<div class="p-4 pb-3">
			<div class="flex items-center gap-1.5">
				<span class="text-[12px] tracking-wide text-[#5c5c5c]">{v.name}</span>
				<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 text-[10.5px] tracking-wide text-[#5c5c5c]">
					{v.symbol?.toUpperCase()}
				</span>
				{#if v.badge}
					<span class="text-[10.5px] tracking-wide text-[#4a4a4a]">{v.badge}</span>
				{/if}
			</div>
			<div class="mt-3 flex items-baseline gap-3">
				<span class="text-[30px] leading-none text-[#e0e0e0]">{s}{price(v.price)}</span>
				{#if v.changePercent != null}
					<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 text-[10.5px] {up ? 'text-[#a6e3a1]' : 'text-[#f38ba8]'}">
						{v.change != null ? `${v.change >= 0 ? '+' : ''}${v.change.toFixed(2)} (` : ''}{v.changePercent >= 0 ? '+' : ''}{v.changePercent.toFixed(2)}%{v.change != null ? ')' : ''}
					</span>
				{/if}
			</div>
			{#if v.high && v.low}
				<div class="mt-1 text-[11px] text-[#5c5c5c]">{v.rangeLabel} {s}{price2(v.low)} – {s}{price2(v.high)}</div>
			{/if}
		</div>
		{#if v.points && v.points.length > 1}
			<div><Sparkline points={v.points} height={44} /></div>
		{/if}
	</div>
{/if}