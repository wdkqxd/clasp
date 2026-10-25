<script lang="ts">
	import { page } from '$app/state';

	let top = $state(0);
	let left = $state(0);
	let width = $state(0);
	let height = $state(0);
	let active = $state(false);

	$effect(() => {
		const hash = page.url.hash;
		if (!hash) {
			active = false;
			return;
		}

		let showTimer: ReturnType<typeof setTimeout> | undefined;
		let hideTimer: ReturnType<typeof setTimeout> | undefined;
		let frame = 0;

		active = false;

		showTimer = setTimeout(() => {
			const id = decodeURIComponent(hash.slice(1));
			const heading = id ? document.getElementById(id) : null;
			if (!heading) return;

			heading.scrollIntoView();

			frame = requestAnimationFrame(() => {
				const target = heading.closest('section') ?? heading;
				const rect = target.getBoundingClientRect();
				const pad = 10;

				top = rect.top - pad;
				left = rect.left - pad;
				width = rect.width + pad * 2;
				height = rect.height + pad * 2;
				active = true;

				hideTimer = setTimeout(() => {
					active = false;
				}, 1000);
			});
		}, 80);

		return () => {
			clearTimeout(showTimer);
			if (hideTimer) clearTimeout(hideTimer);
			cancelAnimationFrame(frame);
			active = false;
		};
	});
</script>

<div
	aria-hidden="true"
	class="pointer-events-none fixed z-40 rounded-[6px] shadow-[0_0_0_100vmax_rgba(16,16,16,0.72)] transition-opacity duration-300 ease-out motion-reduce:transition-none"
	style:top="{top}px"
	style:left="{left}px"
	style:width="{width}px"
	style:height="{height}px"
	style:opacity={active ? 1 : 0}
></div>
