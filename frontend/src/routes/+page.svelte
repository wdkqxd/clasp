<script lang="ts">
	import { goto } from '$app/navigation';
	import { Info, Settings2, ShieldCheck } from '@lucide/svelte';
	import { openBangOr } from '$lib/bangs';
	import SearchForm from '$lib/components/search/SearchForm.svelte';
	import logo from '$lib/assets/logo.svg';

	const links = [
		{ href: '/settings', label: 'settings', icon: Settings2 },
		{ href: '/about/privacy', label: 'privacy', icon: ShieldCheck },
		{ href: '/about/general', label: 'about', icon: Info }
	];

	let query = $state('');
	let menuOpen = $state(false);
	let panelHeight = $state(0);
	let footerEl = $state<HTMLElement | null>(null);
	let closeTimer: ReturnType<typeof setTimeout> | undefined;

	function onSearch(q: string) {
		openBangOr(q, `/search?q=${encodeURIComponent(q)}`, goto);
	}

	function toggleMenu() {
		menuOpen = !menuOpen;
		cancelClose();
	}

	function cancelClose() {
		if (closeTimer) {
			clearTimeout(closeTimer);
			closeTimer = undefined;
		}
	}

	function scheduleClose() {
		if (!menuOpen || closeTimer) return;
		closeTimer = setTimeout(() => {
			menuOpen = false;
			closeTimer = undefined;
		}, 500);
	}

	$effect(() => () => cancelClose());
</script>

<svelte:head>
	<title>clasp</title>
</svelte:head>

<svelte:window
	onkeydown={(e) => {
		if (e.key === 'Escape') menuOpen = false;
	}}
	onmousemove={(e) => {
		if (!footerEl || !menuOpen) return;
		if (e.clientY < footerEl.getBoundingClientRect().top) scheduleClose();
		else cancelClose();
	}}
/>

<div class="relative flex min-h-screen flex-col overflow-hidden px-4">
	<div class="flex flex-1 items-center justify-center">
		<div class="w-full max-w-xl">
			<div class="mb-8 flex justify-center">
				<img src={logo} alt="" class="block h-12 w-auto" />
			</div>
			<SearchForm bind:value={query} {onSearch} />
		</div>
	</div>

	<div
		bind:this={footerEl}
		class="relative shrink-0 transition-transform duration-[340ms] ease-[cubic-bezier(0.22,1,0.36,1)]"
		style:transform={menuOpen ? `translateY(-${panelHeight}px)` : 'translateY(0)'}
	>
		<button
			type="button"
			aria-label="Toggle footer menu"
			aria-expanded={menuOpen}
			onclick={toggleMenu}
			class="group/handle block w-full cursor-pointer px-6 pt-1 pb-2"
		>
			<span
				class="mx-auto block h-[3px] w-10 rounded-full bg-[#3d3d3d] transition-colors duration-200 group-hover/handle:bg-[#707070]"
			></span>
		</button>

		<footer
			bind:clientHeight={panelHeight}
			class="absolute top-full left-0 flex w-full items-start justify-center gap-1 px-1 pt-2 pb-2"
		>
			{#each links as { href, label, icon: Icon } (href)}
				<a
					{href}
					tabindex={menuOpen ? 0 : -1}
					class="inline-flex w-[76px] flex-col items-center gap-1 rounded-[5px] px-3 py-2 text-[13px] text-[#707070] no-underline transition-all duration-150 ease-out hover:bg-[#1a1a1a] hover:text-[#c0c0c0] active:scale-[0.94] active:duration-75"
				>
					<span class="pointer-events-none">
						<Icon class="size-4" strokeWidth={1.75} />
					</span>
					{label}
				</a>
			{/each}
		</footer>
	</div>
</div>
