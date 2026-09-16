<script lang="ts">
	import { market } from '../../api';
	import type { MarketQuote } from '../../types';

	let { quote: initialQuote }: { quote: MarketQuote } = $props();
	let override = $state<MarketQuote | null>(null);
	let quote = $derived(override ?? initialQuote);
	let selectedPeriod = $state('7d');
	let loading = $state(false);
	let chartElement = $state<SVGSVGElement | null>(null);
	let hoveredIndex = $state<number | null>(null);
	const periods = ['1d', '7d', '1m', '1y'];

	const numberFormat = new Intl.NumberFormat(undefined, { maximumFractionDigits: 6 });

	function format(value: number): string {
		return numberFormat.format(value);
	}

	function path(points: MarketQuote['points']): string {
		if (points.length < 2) return '';
		const values = points.map((point) => point.value);
		const min = Math.min(...values);
		const max = Math.max(...values);
		const range = max - min || Math.max(Math.abs(max) * 0.01, 1);
		return points
			.map((point, index) => {
				const x = (index / (points.length - 1)) * 320;
				const y = 8 + (1 - (point.value - min) / range) * 94;
				return `${x.toFixed(2)},${y.toFixed(2)}`;
			})
			.join(' ');
	}

	function hover(event: PointerEvent) {
		if (!chartElement || quote.points.length < 2) return;
		const bounds = chartElement.getBoundingClientRect();
		const ratio = Math.max(0, Math.min(1, (event.clientX - bounds.left) / bounds.width));
		hoveredIndex = Math.round(ratio * (quote.points.length - 1));
	}

	function pointX(index: number): number {
		return (index / Math.max(quote.points.length - 1, 1)) * 320;
	}

	function pointY(index: number): number {
		const values = quote.points.map((point) => point.value);
		const min = Math.min(...values);
		const max = Math.max(...values);
		const range = max - min || Math.max(Math.abs(max) * 0.01, 1);
		return 8 + (1 - (quote.points[index].value - min) / range) * 94;
	}

	function pointDate(time: number): string {
		return new Intl.DateTimeFormat(undefined, { month: 'short', day: 'numeric', hour: 'numeric', minute: '2-digit' }).format(new Date(time));
	}

	async function changePeriod(period: string) {
		if (period === selectedPeriod || loading) return;
		selectedPeriod = period;
		loading = true;
		try {
			const next = await market(quote.query, period);
			if (next) {
				override = next;
				hoveredIndex = null;
			}
		} finally {
			loading = false;
		}
	}
</script>

<div>
	<div class="flex items-start justify-between gap-4">
		<div class="min-w-0">
			<div class="truncate text-[11px] text-[#666]">{quote.name} · {quote.symbol}</div>
			<div class="mt-1 flex items-baseline gap-2">
				<span class="text-[21px] leading-tight text-[#d0d0d0]">{format(quote.current)}</span>
				<span class="text-[11px] text-[#666]">{quote.currency}</span>
			</div>
		</div>
		<div class:text-[#b5d7bd]={quote.change >= 0} class:text-[#d7a4a4]={quote.change < 0} class="pt-1 text-right text-[11px]">
			<div>{quote.change >= 0 ? '+' : ''}{format(quote.change)}</div>
			<div>{quote.change_percent >= 0 ? '+' : ''}{quote.change_percent.toFixed(2)}%</div>
		</div>
	</div>

	<div class:opacity-50={loading} class="relative mt-4 transition-opacity duration-200">
		<svg
			bind:this={chartElement}
			class="h-28 w-full cursor-crosshair overflow-visible"
			viewBox="0 0 320 110"
			preserveAspectRatio="none"
			aria-label="price chart"
			role="img"
			onpointermove={hover}
			onpointerleave={() => (hoveredIndex = null)}
		>
			<line x1="0" y1="103" x2="320" y2="103" stroke="#262626" stroke-width="1" />
			<polyline points={path(quote.points)} fill="none" stroke={quote.change >= 0 ? '#9abca2' : '#c49b9b'} stroke-width="2" stroke-linecap="round" stroke-linejoin="round" vector-effect="non-scaling-stroke" />
			{#if hoveredIndex !== null}
				<line x1={pointX(hoveredIndex)} y1="4" x2={pointX(hoveredIndex)} y2="106" stroke="#444" stroke-width="1" stroke-dasharray="3 3" vector-effect="non-scaling-stroke" />
				<circle cx={pointX(hoveredIndex)} cy={pointY(hoveredIndex)} r="4" fill="#101010" stroke={quote.change >= 0 ? '#b5d7bd' : '#d7a4a4'} stroke-width="2" vector-effect="non-scaling-stroke" />
			{/if}
		</svg>
		{#if hoveredIndex !== null}
			{@const hoveredPoint = quote.points[hoveredIndex]}
			<div
				class="pointer-events-none absolute top-0 z-10 -translate-x-1/2 rounded border border-[#333] bg-[#151515] px-2 py-1 text-[10px] shadow-lg shadow-black/20"
				style:left={`${(pointX(hoveredIndex) / 320) * 100}%`}
			>
				<div class="text-[#aaa]">{pointDate(hoveredPoint.time)}</div>
				<div class="mt-0.5 text-[#d0d0d0]">{format(hoveredPoint.value)} {quote.currency}</div>
			</div>
		{/if}
	</div>

	<div class="mt-2 flex items-center justify-between">
		<div class="flex gap-1">
			{#each periods as period}
				<button
					type="button"
					class:text-[#c8c8c8]={selectedPeriod === period}
					class:text-[#555]={selectedPeriod !== period}
					class="rounded px-1.5 py-1 text-[10px] transition-colors hover:text-[#aaa]"
					onclick={() => changePeriod(period)}
				>
					{period}
				</button>
			{/each}
		</div>
		<span class="text-[10px] text-[#555]">{quote.change >= 0 ? 'up' : 'down'}</span>
	</div>
</div>
