<script lang="ts">
	import type { RichAnswer } from '../../types';
	import Colorpicker from './Colorpicker.svelte';
	import Weather from './Weather.svelte';
	import Timer from './Timer.svelte';
	import Stopwatch from './Stopwatch.svelte';
	import UnitConversion from './UnitConversion.svelte';
	import Timezones from './Timezones.svelte';
	import Cryptocurrency from './Cryptocurrency.svelte';
	import Currency from './Currency.svelte';
	import Unixtimestamp from './Unixtimestamp.svelte';
	import TopAnswer from './TopAnswer.svelte';

	let { items }: { items: RichAnswer[] } = $props();
</script>

{#if items.length > 0}
	<section class="space-y-10 border-[#262626]">
		{#each items as a, i (i)}
			{#if a.type === 'topAnswer' && a.topAnswer}
				<TopAnswer topAnswer={a.topAnswer} />
			{:else if a.type === 'colorpicker' && a.colorpicker}
				<Colorpicker colorpicker={a.colorpicker} />
			{:else if a.type === 'weather' && a.weather}
				<Weather weather={a.weather} />
			{:else if a.type === 'timer' && a.timer}
				<Timer timer={a.timer} />
			{:else if a.type === 'stopwatch'}
				<Stopwatch />
			{:else if a.type === 'unitConversion' && a.unitConversion}
				<UnitConversion unit={a.unitConversion} />
			{:else if a.type === 'timezones' && a.timezones?.timezones?.length}
				<Timezones timezones={a.timezones} />
			{:else if a.type === 'cryptocurrency' && a.cryptocurrency}
				<Cryptocurrency crypto={a.cryptocurrency} />
			{:else if a.type === 'currency' && a.currency}
				<Currency currency={a.currency} />
			{:else if a.type === 'unixtimestamp'}
				<Unixtimestamp />
			{/if}
		{/each}
	</section>
{/if}
