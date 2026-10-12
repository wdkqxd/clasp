<script lang="ts">
	import { getCurrencyByCode, parseMoney } from 'numis';
	import type { RichCurrency } from '../../../../types';
	import Sparkline from '../../shared/Sparkline.svelte';
	import { price, price2 } from '../../shared/format';
	import CopyButton from '../../shared/CopyButton.svelte';

	let { query = '' }: { query?: string } = $props();

	const FRANKFURTER = 'https://api.frankfurter.dev/v2';
	const PAIR_RE = /^(.+?)\s+(?:to|in|into)\s+(.+)$/i;
	const TARGET_TAIL_RE = /\s+(?:exchange\s+rate|rate|chart|price|value|worth)$/i;
	const DEFAULT_CURRENCY = 'EUR';
	const FALLBACK_CURRENCY = 'USD';

	interface FrankfurterRate {
		date: string;
		base: string;
		quote: string;
		rate: number;
	}

	let resolved = $state<RichCurrency | null>(null);
	let requestId = 0;

	$effect(() => {
		const q = query.trim();
		if (!q) return;
		const current = ++requestId;
		void convert(q).then((res) => {
			if (current === requestId) resolved = res;
		});
	});

	const data = $derived(resolved);

	let copied = $state(false);
	function onCopy() {
		navigator.clipboard?.writeText(String(data?.result));
		copied = true;
		setTimeout(() => (copied = false), 1200);
	}

	async function convert(q: string): Promise<RichCurrency | null> {
		const intent = currencyIntent(q);
		if (!intent) return null;
		const { amount, from, to } = intent;
		if (from === to) return null;

		const rates = await getJSON<FrankfurterRate[]>(
			`${FRANKFURTER}/rates?base=${from}&quotes=${to}`
		);
		const latest = rates?.at(-1);
		if (!latest || !Number.isFinite(latest.rate)) return null;

		const out: RichCurrency = {
			amount,
			from_currency_code: from,
			to_currency_code: to,
			from_currency_name: getCurrencyByCode(from)?.currency,
			to_currency_name: getCurrencyByCode(to)?.currency,
			result: amount * latest.rate,
			exchange_rate: latest.rate,
			date: latest.date
		};

		const series = await exchangeSeries(from, to);
		if (series) out.timeseries = { time_range: '1M', ts_exchange_rate: series };
		return out;
	}

	function currencyIntent(q: string): { amount: number; from: string; to: string } | null {
		const pair = q.match(PAIR_RE);
		if (pair) {
			const from = money(pair[1]);
			const to = currencyCode(pair[2].replace(TARGET_TAIL_RE, ''));
			if (!from || !to) return null;
			return { amount: from.amount ?? 1, from: from.currency, to };
		}
		const bare = money(q);
		if (!bare?.amount) return null;
		return {
			amount: bare.amount,
			from: bare.currency,
			to: bare.currency === DEFAULT_CURRENCY ? FALLBACK_CURRENCY : DEFAULT_CURRENCY
		};
	}

	function money(text: string): { amount?: number; currency: string } | null {
		try {
			const parsed = parseMoney(text);
			if (!parsed?.currency || !parsed.original) return null;
			const info = getCurrencyByCode(parsed.currency.toUpperCase());
			if (!info) return null;
			if (compact(parsed.original) !== compact(text)) return null;
			return { amount: parsed.amount, currency: info.code };
		} catch {
			return null;
		}
	}

	function currencyCode(text: string): string | null {
		return money(text)?.currency ?? null;
	}

	function compact(s: string): string {
		return s.replace(/\s+/g, '').toLowerCase();
	}

	async function exchangeSeries(from: string, to: string): Promise<[string, number][] | null> {
		const toDate = new Date();
		const fromDate = new Date(toDate);
		fromDate.setDate(fromDate.getDate() - 30);
		const rates = await getJSON<FrankfurterRate[]>(
			`${FRANKFURTER}/rates?base=${from}&quotes=${to}` +
				`&from=${isoDate(fromDate)}&to=${isoDate(toDate)}`
		);
		if (!rates || rates.length < 2) return null;
		return rates.map((r) => [r.date, r.rate]);
	}

	function isoDate(d: Date): string {
		return d.toISOString().slice(0, 10);
	}

	async function getJSON<T>(url: string): Promise<T | null> {
		try {
			const res = await fetch(url);
			if (!res.ok) return null;
			return (await res.json()) as T;
		} catch {
			return null;
		}
	}
</script>

{#if data}
	{@const c = data}
	<div class="overflow-hidden rounded-md border border-[#292929] bg-[#121212] shadow-2xl shadow-black/30">
		<div class="p-4 pb-3">
			<p class="text-[16px] text-[#c0c0c0]">
				<span class="text-[#707070]">{price(c.amount)} {c.from_currency_name || c.from_currency_code} =</span>
				<button
					type="button"
					class="group inline-flex cursor-pointer items-center gap-3 text-[#e0e0e0]"
					title="copy result"
					onclick={onCopy}
				>
					{price2(c.result)} {c.to_currency_name || c.to_currency_code}
					<CopyButton active={copied} />
				</button>
			</p>
			{#if c.exchange_rate || c.date}
				<div class="mt-2 text-[11px] text-[#5c5c5c]">
					{#if c.exchange_rate}<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 tracking-wide">1 {c.from_currency_code} = {c.exchange_rate} {c.to_currency_code}</span>{/if}
					{#if c.exchange_rate && c.date}<span class="mx-1">·</span>{/if}
					{#if c.date}{c.date}{/if}
				</div>
			{/if}
		</div>
		{#if c.timeseries?.ts_exchange_rate?.length}
			{@const pts = c.timeseries.ts_exchange_rate.map((p) => [new Date(p[0]).getTime(), p[1]] as [number, number])}
			{#if pts.length > 1}
				<div class="-mx-0"><Sparkline points={pts} height={44} /></div>
			{/if}
		{/if}
	</div>
{/if}
