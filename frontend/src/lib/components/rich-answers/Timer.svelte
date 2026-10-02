<script lang="ts">
	let { query = '' }: { query?: string } = $props();

	const TIMER_RE = /^(?:set\s+(?:a|up)\s+)?(?:a\s+)?timer(?:\s+for)?\s*(.*)$/i;
	const DURATION_RE = /^(?:(\d+(?:\.\d+)?)\s*(?:h|hr|hrs|hour|hours))?\s*(?:(\d+(?:\.\d+)?)\s*(?:m|min|mins|minute|minutes))?\s*(?:(\d+(?:\.\d+)?)\s*(?:s|sec|secs|second|seconds))?$/i;
	const CLOCK_RE = /^(\d{1,2}):([0-5]?\d)(?::([0-5]?\d))?$/;

	let data = $derived(parseTimer(query.trim()));

	type TimerData = { duration: number; start_on_load: boolean };

	function parseTimer(q: string): TimerData | null {
		const match = q.match(TIMER_RE);
		if (!match) return null;
		const rest = match[1].trim();
		if (!rest) return { duration: 300, start_on_load: false };

		const clock = rest.match(CLOCK_RE);
		if (clock) {
			const [, a, b, c] = clock;
			const duration = Number(a) * 3600 + Number(b) * 60 + Number(c ? Number(c) : 0);
			return duration > 0 ? { duration, start_on_load: false } : null;
		}

		const parts = rest.match(DURATION_RE);
		if (!parts) return null;
		const [, h, m, s] = parts;
		if (!h && !m && !s) return null;
		const duration = Number(h ?? 0) * 3600 + Number(m ?? 0) * 60 + Number(s ?? 0);
		return duration > 0 ? { duration, start_on_load: false } : null;
	}

	let timerTotal = $state(0);
	let timerLeft = $state(0);
	let timerRunning = $state(false);
	let timerId: ReturnType<typeof setInterval> | null = null;

	$effect(() => {
		const duration = data?.duration;
		if (!duration || duration === timerTotal) return;
		timerTotal = duration;
		timerLeft = duration;
		timerRunning = false;
		if (timerId) clearInterval(timerId);
	});

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

{#if data}
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
{/if}