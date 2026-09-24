<script lang="ts">
	let {
		src,
		alt,
		selected = false,
		classNames = ''
	}: {
		src: string;
		alt: string;
		selected?: boolean;
		classNames?: string;
	} = $props();

	let ref = $state<HTMLDivElement | null>(null);
	let inView = $state(false);
	let loaded = $state(false);

	let observer: IntersectionObserver | undefined;

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
			src={src}
			alt={alt}
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