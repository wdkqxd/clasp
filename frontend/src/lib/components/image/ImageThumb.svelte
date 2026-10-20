<script lang="ts">
	let {
		src,
		alt,
		selected = false,
		classNames = '',
		detectTransparency = false,
		ontransparency
	}: {
		src: string;
		alt: string;
		selected?: boolean;
		classNames?: string;
		detectTransparency?: boolean;
		ontransparency?: (value: boolean | null) => void;
	} = $props();

	let ref = $state<HTMLDivElement | null>(null);
	let imageEl = $state<HTMLImageElement | null>(null);
	let inView = $state(false);
	let loaded = $state(false);
	let detectedSrc = '';

	let observer: IntersectionObserver | undefined;

	function detectAlpha(img: HTMLImageElement): boolean | null {
		try {
			const canvas = document.createElement('canvas');
			canvas.width = 48;
			canvas.height = 48;
			const ctx = canvas.getContext('2d', { willReadFrequently: true });
			if (!ctx) return null;
			ctx.clearRect(0, 0, 48, 48);
			ctx.drawImage(img, 0, 0, 48, 48);
			const { data } = ctx.getImageData(0, 0, 48, 48);
			let transparent = 0;
			for (let i = 3; i < data.length; i += 4) {
				if (data[i] < 250) transparent++;
			}
			return transparent / (data.length / 4) > 0.01;
		} catch {
			return null;
		}
	}

	$effect(() => {
		if (!detectTransparency || !loaded || !imageEl || detectedSrc === src) return;
		detectedSrc = src;
		ontransparency?.(detectAlpha(imageEl));
	});

	$effect(() => {
		if (!ref) return;
		observer = new IntersectionObserver(
			(entries) => {
				for (const entry of entries) {
					if (entry.isIntersecting) {
						inView = true;
						observer?.disconnect();
					}
				}
			},
			{ rootMargin: '400px' }
		);
		observer.observe(ref);
		return () => observer?.disconnect();
	});
</script>

<div bind:this={ref} class="w-full h-full">
	{#if inView}
		<img
			bind:this={imageEl}
			src={src}
			alt={alt}
			crossorigin={src.startsWith('https://imgs.search.brave.com/') ? 'anonymous' : undefined}
			class="w-full h-full object-cover transition-[opacity,transform] duration-300 ease-out {classNames}"
			class:img-loaded={loaded}
			class:selected-img={selected}
			onload={() => (loaded = true)}
		/>
	{/if}
</div>

<style>
	div :global(img) {
		opacity: 0;
		transform: scale(1.04);
	}
	div :global(img.img-loaded) {
		opacity: 1;
		transform: scale(1);
	}
</style>