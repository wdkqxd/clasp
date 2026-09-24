<script lang="ts">
	let swMs = $state(0);
	let swRunning = $state(false);
	let swLaps = $state<number[]>([]);
	let swId: ReturnType<typeof setInterval> | null = null;

	function onSw() {
		if (swRunning) {
			swRunning = false;
			if (swId) clearInterval(swId);
			return;
		}
		swRunning = true;
		swId = setInterval(() => (swMs += 10), 10);
	}
	function onSwLap() {
		if (swRunning) swLaps = [swMs, ...swLaps];
		else {
			swMs = 0;
			swLaps = [];
		}
	}
	function fmtSw(ms: number): string {
		const s = Math.floor(ms / 1000);
		const c = Math.floor((ms % 1000) / 10);
		return `${Math.floor(s / 60)}:${String(s % 60).padStart(2, '0')}.${String(c).padStart(2, '0')}`;
	}
</script>

<div class="rounded-md border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
	<div class="font-mono text-[32px] leading-none text-[#e0e0e0]">{fmtSw(swMs)}</div>
	<div class="mt-3 flex items-center gap-2 text-[12px]">
		<button
			type="button"
			class="cursor-pointer rounded bg-[#1c1c1c] px-3 py-1 text-[#9b9b9b] transition-colors hover:bg-[#262626] hover:text-[#d0d0d0]"
			onclick={onSw}
		>
			{swRunning ? 'pause' : swMs > 0 ? 'resume' : 'start'}
		</button>
		<button
			type="button"
			class="cursor-pointer rounded bg-[#1c1c1c] px-3 py-1 text-[#707070] transition-colors hover:bg-[#262626] hover:text-[#9b9b9b]"
			onclick={onSwLap}
		>
			{swRunning ? 'lap' : 'reset'}
		</button>
	</div>
	{#if swLaps.length}
		<div class="mt-3 flex flex-col gap-1 text-[12px] text-[#707070]">
			{#each swLaps as lap, i (i)}
				<span>lap {swLaps.length - i} · {fmtSw(lap)}</span>
			{/each}
		</div>
	{/if}
</div>