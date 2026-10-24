<script lang="ts">
	import type { RichWeather, RichWeatherCond } from '../../../../types';

	let { query = '' }: { query?: string } = $props();

	const GEOCODING = 'https://geocoding-api.open-meteo.com/v1/search';
	const FORECAST = 'https://api.open-meteo.com/v1/forecast';

	const WEATHER_RE =
		/^(?:what(?:'s| is)\s+the\s+)?(?:weather|forecast|temperature)\b\s*(?:in|for|at|near)?\s*(.+?)\s*$/i;
	const WEATHER_TAIL_RE = /\s+(?:today|now|right now|currently|at the moment|this week|tomorrow)$/i;

	interface GeoResult {
		name: string;
		country?: string;
		admin1?: string;
		latitude: number;
		longitude: number;
	}

	interface ForecastResponse {
		current?: {
			temperature_2m?: number;
			apparent_temperature?: number;
			relative_humidity_2m?: number;
			wind_speed_10m?: number;
			surface_pressure?: number;
			weather_code?: number;
		};
		daily?: {
			time?: string[];
			weather_code?: number[];
			temperature_2m_max?: number[];
			temperature_2m_min?: number[];
			uv_index_max?: number[];
		};
	}

	let resolved = $state<RichWeather | null>(null);
	let requestId = 0;

	$effect(() => {
		const q = query.trim();
		if (!q) return;
		const current = ++requestId;
		void forecast(q).then((res) => {
			if (current === requestId) resolved = res;
		});
	});

	const data = $derived(resolved);
	const loc = $derived(
		[data?.location?.name, data?.location?.state, data?.location?.country]
			.filter(Boolean)
			.join(', ')
	);

	async function forecast(q: string): Promise<RichWeather | null> {
		const match = q.match(WEATHER_RE);
		if (!match) return null;
		const place = match[1].replace(WEATHER_TAIL_RE, '').trim();
		if (!place || place.length > 60) return null;

		const geo = await getJSON<{ results?: GeoResult[] }>(
			`${GEOCODING}?name=${encodeURIComponent(place)}&count=1&language=en&format=json`
		);
		const hit = geo?.results?.[0];
		if (!hit) return null;

		const params = new URLSearchParams({
			latitude: String(hit.latitude),
			longitude: String(hit.longitude),
			current:
				'temperature_2m,relative_humidity_2m,apparent_temperature,weather_code,wind_speed_10m,surface_pressure',
			daily: 'weather_code,temperature_2m_max,temperature_2m_min,uv_index_max',
			wind_speed_unit: 'ms',
			timezone: 'auto',
			forecast_days: '5'
		});
		const data = await getJSON<ForecastResponse>(`${FORECAST}?${params}`);
		if (!data?.current || data.current.temperature_2m === undefined) return null;

		const daily = data.daily;
		const out: RichWeather = {
			location: { name: hit.name, country: hit.country ?? '', state: hit.admin1 ?? '' },
			current_weather: {
				temp: data.current.temperature_2m,
				feels_like: data.current.apparent_temperature,
				humidity: data.current.relative_humidity_2m,
				wind_speed: data.current.wind_speed_10m,
				pressure: data.current.surface_pressure,
				uvi: daily?.uv_index_max?.[0],
				weather: weatherCond(data.current.weather_code)
			}
		};

		if (daily?.time?.length) {
			out.daily = daily.time.map((date, i) => ({
				ts: Math.floor(Date.parse(`${date}T12:00:00`) / 1000),
				date_i18n: date,
				temperature: {
					max: daily.temperature_2m_max?.[i],
					min: daily.temperature_2m_min?.[i]
				},
				weather: weatherCond(daily.weather_code?.[i])
			}));
		}
		return out;
	}

	const WMO_CODES: Record<number, [string, string]> = {
		0: ['Clear', 'clear sky'],
		1: ['Clear', 'mainly clear'],
		2: ['Clouds', 'partly cloudy'],
		3: ['Clouds', 'overcast'],
		45: ['Fog', 'fog'],
		48: ['Fog', 'depositing rime fog'],
		51: ['Drizzle', 'light drizzle'],
		53: ['Drizzle', 'moderate drizzle'],
		55: ['Drizzle', 'dense drizzle'],
		56: ['Drizzle', 'light freezing drizzle'],
		57: ['Drizzle', 'dense freezing drizzle'],
		61: ['Rain', 'slight rain'],
		63: ['Rain', 'moderate rain'],
		65: ['Rain', 'heavy rain'],
		66: ['Rain', 'light freezing rain'],
		67: ['Rain', 'heavy freezing rain'],
		71: ['Snow', 'slight snowfall'],
		73: ['Snow', 'moderate snowfall'],
		75: ['Snow', 'heavy snowfall'],
		77: ['Snow', 'snow grains'],
		80: ['Rain', 'slight rain showers'],
		81: ['Rain', 'moderate rain showers'],
		82: ['Rain', 'violent rain showers'],
		85: ['Snow', 'slight snow showers'],
		86: ['Snow', 'heavy snow showers'],
		95: ['Thunderstorm', 'thunderstorm'],
		96: ['Thunderstorm', 'thunderstorm with slight hail'],
		99: ['Thunderstorm', 'thunderstorm with heavy hail']
	};

	function weatherCond(code?: number): RichWeatherCond | undefined {
		if (code === undefined) return undefined;
		const hit = WMO_CODES[code];
		if (!hit) return { main: 'Unknown', description: `weather code ${code}` };
		return { main: hit[0], description: hit[1] };
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
	{@const w = data}
	<div class="rounded-[5px] border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
		<div class="flex items-center justify-between gap-3">
			<span class="text-[12px] tracking-wide text-[#5c5c5c]">{loc}</span>
			{#if w.current_weather?.weather?.description}
				<span class="rounded bg-[#1c1c1c] px-1.5 py-0.5 text-[10.5px] tracking-wide text-[#9b9b9b]">
					{w.current_weather.weather.description}
				</span>
			{/if}
		</div>
		{#if w.current_weather}
			{@const cw = w.current_weather}
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
{/if}
