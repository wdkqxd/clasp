<script lang="ts">
	let { query = '' }: { query?: string } = $props();

	type Timezone = {
		abbreviation: string;
		generic_name: string;
		utc_offset: string;
		converted_time?: { strftime?: string; location?: string; city?: string; utc_diff?: string };
	};
	type TimezonesData = { timezones?: Timezone[] };

	const GEOCODING = 'https://geocoding-api.open-meteo.com/v1/search';

	const TIME_RE = /^(?:what(?:'s| is)\s+)?(?:the\s+)?(?:current\s+)?time(?:\s+zone)?\s*(?:in|at|for)\s+(.+?)\s*$/i;

	interface GeoResult {
		name: string;
		country?: string;
		admin1?: string;
		timezone: string;
	}

	let resolved = $state<TimezonesData | null>(null);
	let requestId = 0;

	$effect(() => {
		const q = query.trim();
		if (!q) return;
		const current = ++requestId;
		void load(q).then((res) => {
			if (current === requestId) resolved = res;
		});
	});

	const data = $derived(resolved);

	let tick = $state(0);
	$effect(() => {
		if (!data) return;
		const id = setInterval(() => (tick += 1), 30000);
		return () => clearInterval(id);
	});

	async function load(q: string): Promise<TimezonesData | null> {
		const match = q.match(TIME_RE);
		if (!match) return null;
		const place = match[1];
		if (place.length > 60) return null;

		const geo = await getJSON<{ results?: GeoResult[] }>(
			`${GEOCODING}?name=${encodeURIComponent(place)}&count=1&language=en&format=json`
		);
		const hit = geo?.results?.[0];
		if (!hit?.timezone) return null;

		const city = [hit.name, hit.admin1].filter(Boolean).join(', ');
		return {
			timezones: [
				{
					abbreviation: tzAbbreviation(hit.timezone),
					generic_name: hit.timezone.replace(/_/g, ' '),
					utc_offset: tzOffset(hit.timezone),
					converted_time: {
						strftime: tzTime(hit.timezone),
						location: city,
						city,
						utc_diff: tzOffset(hit.timezone)
					}
				}
			]
		};
	}

	function tzTime(zone: string): string {
		return new Intl.DateTimeFormat(undefined, {
			timeZone: zone,
			hour: 'numeric',
			minute: '2-digit'
		}).format(new Date());
	}

	function tzOffset(zone: string): string {
		const parts = new Intl.DateTimeFormat('en-US', {
			timeZone: zone,
			timeZoneName: 'shortOffset'
		}).formatToParts(new Date());
		return parts.find((p) => p.type === 'timeZoneName')?.value ?? '';
	}

	function tzAbbreviation(zone: string): string {
		const parts = new Intl.DateTimeFormat('en-US', {
			timeZone: zone,
			timeZoneName: 'short'
		}).formatToParts(new Date());
		return parts.find((p) => p.type === 'timeZoneName')?.value ?? '';
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

	function tzTimeOf(t: Timezone): { time: string; ampm: string } {
		const raw = (t.converted_time?.strftime ?? '').trim();
		const mer = raw.match(/\s+([AP]M)$/i);
		if (mer) return { time: raw.slice(0, mer.index), ampm: mer[1].toLowerCase() };
		return { time: raw, ampm: '' };
	}
</script>

{#if data?.timezones?.length}
	{@const list = data.timezones.slice(0, 6)}
	<div class="rounded-[5px] border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
		<div class="flex flex-col gap-1.5">
			{#each list as tz, i (i)}
				{@const t = tzTimeOf(tz)}
				<div class="flex items-baseline gap-4 text-[13px]">
					<span class="w-24 shrink-0 font-mono text-[#e0e0e0]">
						{t.time}{#if t.ampm}<span class="text-[#5c5c5c]"> {t.ampm}</span>{/if}
					</span>
					<span class="min-w-0 flex-1 truncate text-[#9a9a9a]">
						{tz.converted_time?.city || tz.converted_time?.location || tz.abbreviation}
					</span>
					<span class="shrink-0 text-[#5c5c5c]">
						{tz.converted_time?.utc_diff || tz.utc_offset || ''}
					</span>
				</div>
			{/each}
		</div>
	</div>
{/if}
