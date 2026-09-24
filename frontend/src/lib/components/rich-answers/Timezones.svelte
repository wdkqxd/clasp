<script lang="ts">
	import type { RichTimezones } from '../../types';

	let { timezones }: { timezones: RichTimezones } = $props();

	function tzTime(strftime?: string): { time: string; ampm: string } {
		const raw = (strftime ?? '').trim();
		const mer = raw.match(/\s+([AP]M)$/i);
		if (mer) return { time: raw.slice(0, mer.index), ampm: mer[1].toLowerCase() };
		return { time: raw, ampm: '' };
	}
</script>

<div class="rounded-md border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
	{#if timezones.timezones?.length}
		<div class="flex flex-col gap-1.5">
			{#each timezones.timezones.slice(0, 6) as tz, i (i)}
				{@const t = tzTime(tz.converted_time?.strftime)}
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
	{/if}
</div>