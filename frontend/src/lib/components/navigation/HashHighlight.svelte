<script lang="ts">
	import { page } from '$app/state';

	$effect(() => {
		const hash = page.url.hash;
		if (!hash) return;

		let showTimer: ReturnType<typeof setTimeout> | undefined;
		let hideTimer: ReturnType<typeof setTimeout> | undefined;
		let dimmed: HTMLElement[] = [];

		const undim = (fade: boolean) => {
			for (const el of dimmed) {
				el.style.opacity = '1';
				if (fade) {
					setTimeout(() => {
						el.style.transition = '';
						el.style.opacity = '';
					}, 320);
				} else {
					el.style.transition = '';
					el.style.opacity = '';
				}
			}
			dimmed = [];
		};

		showTimer = setTimeout(() => {
			const id = decodeURIComponent(hash.slice(1));
			const heading = id ? document.getElementById(id) : null;
			if (!heading) return;

			const target = heading.closest('section') ?? heading;
			const root = target.parentElement;
			if (!root) return;

			for (const section of root.querySelectorAll('section')) {
				if (section === target) continue;
				section.style.transition = 'opacity 300ms ease';
				section.style.opacity = '0.35';
				dimmed.push(section);
			}

			if (dimmed.length > 0) {
				hideTimer = setTimeout(() => undim(true), 1000);
			}
		}, 80);

		return () => {
			clearTimeout(showTimer);
			if (hideTimer) clearTimeout(hideTimer);
			undim(false);
		};
	});
</script>
