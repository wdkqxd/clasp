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

	function onSearch(q: string) {
		openBangOr(q, `/search?q=${encodeURIComponent(q)}`, goto);
	}
</script>

<svelte:head>
	<title>clasp</title>
</svelte:head>

<div class="relative flex min-h-screen flex-col overflow-hidden px-4">
	<div class="flex flex-1 items-center justify-center">
		<div class="w-full max-w-xl">
			<div class="mb-8 flex justify-center">
				<img src={logo} alt="" class="block h-12 w-auto" />
			</div>
			<SearchForm bind:value={query} {onSearch} />
		</div>
	</div>

	<footer
		class="flex items-start justify-center gap-1 pb-3 transition-transform duration-[160ms] [transition-timing-function:cubic-bezier(0.23,1,0.32,1)] {query.trim()
			? 'pointer-events-none translate-y-full'
			: 'translate-y-0'}"
		aria-hidden={query.trim().length > 0}
	>
		{#each links as { href, label, icon: Icon } (href)}
			<a
				{href}
				tabindex={query.trim() ? -1 : 0}
				class="group inline-flex w-[76px] flex-col items-center gap-1 rounded-md px-3 py-2 text-[13px] text-[#707070] no-underline transition-all duration-150 ease-out hover:bg-[#1a1a1a] hover:text-[#c0c0c0] active:scale-[0.94] active:duration-75"
			>
				<span class="pointer-events-none">
					<Icon class="size-4" strokeWidth={1.75} />
				</span>
				{label}
			</a>
		{/each}
	</footer>
</div>
