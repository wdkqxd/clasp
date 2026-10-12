<script lang="ts">
	import type { RichAnswer } from '../../types';
	import Weather from './Weather.svelte';
	import Timer from './Timer.svelte';
	import Stopwatch from './Stopwatch.svelte';
	import UnitConversion from './UnitConversion.svelte';
	import Timezones from './Timezones.svelte';
	import Asset from './Asset.svelte';
	import Currency from './Currency.svelte';
	import Unixtimestamp from './Unixtimestamp.svelte';
	import Expression from './Expression.svelte';
	import Translator from './Translator.svelte';

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
