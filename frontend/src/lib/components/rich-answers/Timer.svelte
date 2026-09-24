<script lang="ts">
	import type { RichTimer } from '../../types';

	let { timer }: { timer: RichTimer } = $props();

	let timerTotal = $state(timer.duration);
	let timerLeft = $state(timer.duration);
	let timerRunning = $state(false);
	let timerId: ReturnType<typeof setInterval> | null = null;

	function fmtTimer(s: number): string {
		const h = Math.floor(s / 3600);
		const m = Math.floor((s % 3600) / 60);
		const sec = Math.round(s % 60);
		if (h > 0) return `${h}:${String(m).padStart(2, '0')}:${String(sec).padStart(2, '0')}`;
		return `${m}:${String(sec).padStart(2, '0')}`;
	}
	function onTimer() {
		if (timerRunning) {
			timerRunning = false;
			if (timerId) clearInterval(timerId);
			return;
		}
		if (timerLeft <= 0) timerLeft = timerTotal;
		timerRunning = true;
		timerId = setInterval(() => {
			timerLeft = Math.max(0, timerLeft - 0.1);
			if (timerLeft <= 0) {
				timerRunning = false;
				if (timerId) clearInterval(timerId);
			}
		}, 100);
	}
	function resetTimer() {
		timerRunning = false;
		if (timerId) clearInterval(timerId);
		timerLeft = timerTotal;
	}
</script>

<div class="rounded-md border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
	<div class="font-mono text-[32px] leading-none text-[#e0e0e0]">{fmtTimer(timerLeft)}</div>
	<div class="mt-3 flex items-center gap-2 text-[12px]">
		<button
			type="button"
			class="cursor-pointer rounded bg-[#1c1c1c] px-3 py-1 text-[#9b9b9b] transition-colors hover:bg-[#262626] hover:text-[#d0d0d0]"
			onclick={onTimer}
		>
			{timerRunning ? 'pause' : timerLeft <= 0 || timerLeft === timerTotal ? 'start' : 'resume'}
		</button>
		<button
			type="button"
			class="cursor-pointer rounded bg-[#1c1c1c] px-3 py-1 text-[#707070] transition-colors hover:bg-[#262626] hover:text-[#9b9b9b]"
			onclick={resetTimer}
		>
			reset
		</button>
	</div>
</div>