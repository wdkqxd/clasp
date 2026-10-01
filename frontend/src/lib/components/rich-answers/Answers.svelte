<script lang="ts">
	import type { RichAnswer } from '../../types';
	import Weather from './Weather.svelte';
	import Timer from './Timer.svelte';
	import Stopwatch from './Stopwatch.svelte';
	import UnitConversion from './UnitConversion.svelte';
	import Timezones from './Timezones.svelte';
	import Cryptocurrency from './Cryptocurrency.svelte';
	import Currency from './Currency.svelte';
	import Unixtimestamp from './Unixtimestamp.svelte';
	import Expression from './Expression.svelte';
	import TopAnswer from './TopAnswer.svelte';

	let { items, query = '' }: { items: RichAnswer[]; query?: string } = $props();

	const rest = $derived(items.filter((a) => a.type === 'topAnswer' || a.type === 'cryptocurrency'));
</script>

{#if items.length > 0 || query}
	<section class="mb-0 space-y-5 border-[#262626] has-[>div]:mb-10 has-[>section]:mb-10">
		<Currency {query} />
		<Weather {query} />
		<Timer {query} />
		<Stopwatch {query} />
		<UnitConversion {query} />
		<Timezones {query} />
		<Unixtimestamp {query} />
		<Expression {query} />
		{#each rest as a, i (i)}
			{#if a.type === 'topAnswer' && a.topAnswer}
				<TopAnswer topAnswer={a.topAnswer} />
			{:else if a.type === 'cryptocurrency' && a.cryptocurrency}
				<Cryptocurrency crypto={a.cryptocurrency} />
			{/if}
		{/each}
	</section>
{/if}
