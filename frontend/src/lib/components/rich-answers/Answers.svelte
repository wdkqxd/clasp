<script lang="ts">
	import type { RichAnswer } from '../../types';
	import Weather from './widgets/weather/Weather.svelte';
	import Timer from './widgets/timer/Timer.svelte';
	import Stopwatch from './widgets/stopwatch/Stopwatch.svelte';
	import UnitConversion from './widgets/unit-conversion/UnitConversion.svelte';
	import Timezones from './widgets/timezones/Timezones.svelte';
	import Asset from './widgets/asset/Asset.svelte';
	import Currency from './widgets/currency/Currency.svelte';
	import Unixtimestamp from './widgets/unixtimestamp/Unixtimestamp.svelte';
	import Expression from './widgets/expression/Expression.svelte';
	import Translator from './widgets/translator/Translator.svelte';

	let { items, query = '' }: { items: RichAnswer[]; query?: string } = $props();

	const rest = $derived(items.filter((a) => a.type === 'cryptocurrency' || a.type === 'stock'));
</script>

{#if items.length > 0 || query}
	<section class="mb-0 space-y-5 border-[#262626] has-[>div]:mb-8 has-[>section]:mb-8">
		<Currency {query} />
		<Weather {query} />
		<Timer {query} />
		<Stopwatch {query} />
		<UnitConversion {query} />
		<Timezones {query} />
		<Unixtimestamp {query} />
		<Expression {query} />
		<Translator {query} />
		{#each rest as a, i (i)}
			{#if (a.type === 'cryptocurrency' || a.type === 'stock') && (a.cryptocurrency || a.stock)}
				<Asset answer={a} />
			{/if}
		{/each}
	</section>
{/if}
