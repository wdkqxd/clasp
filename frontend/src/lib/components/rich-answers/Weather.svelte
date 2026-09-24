<script lang="ts">
	import type { RichWeather } from '../../types';

	let { weather }: { weather: RichWeather } = $props();

	const loc = [weather.location?.name, weather.location?.state, weather.location?.country].filter(Boolean).join(', ');
</script>

<div class="rounded-md border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
	<div class="flex items-center justify-between gap-3">
		<span class="text-[12px] tracking-wide text-[#5c5c5c]">{loc}</span>
		{#if weather.current_weather?.weather?.description}
			<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 text-[10.5px] tracking-wide text-[#9b9b9b]">
				{weather.current_weather.weather.description}
			</span>
		{/if}
	</div>
	{#if weather.current_weather}
		{@const cw = weather.current_weather}
		<div class="mt-3 flex items-baseline gap-3">
			<span class="text-[34px] leading-none text-[#e0e0e0]">{Math.round(cw.temp)}°</span>
			{#if cw.feels_like != null}
				<span class="text-[12px] text-[#5c5c5c]">feels {Math.round(cw.feels_like)}°</span>
			{/if}
			<span class="ml-auto text-[12px] text-[#707070]">
				{#if cw.humidity != null}{Math.round(cw.humidity)}% humidity{/if}
				{#if cw.humidity != null && (cw.wind_speed != null || cw.uvi != null || cw.pressure != null)} · {/if}
				{#if cw.wind_speed != null}{Math.round(cw.wind_speed)} m/s wind{/if}
				{#if cw.wind_speed != null && (cw.uvi != null || cw.pressure != null)} · {/if}
				{#if cw.uvi != null}uv {Math.round(cw.uvi)}{/if}
				{#if cw.uvi != null && cw.pressure != null} · {/if}
				{#if cw.pressure != null}{Math.round(cw.pressure)} hPa{/if}
			</span>
		</div>
	{/if}
</div>