<script lang="ts">
	import { goto } from '$app/navigation';
	import { openBangOr } from '$lib/bangs';
	import SearchForm from '$lib/components/search/SearchForm.svelte';
	import logo from '$lib/assets/logo.svg';

	function onSearch(q: string) {
		openBangOr(q, `/search?q=${encodeURIComponent(q)}`, goto);
	}
</script>

<svelte:head>
	<title>clasp</title>
</svelte:head>

<div class="relative flex min-h-screen flex-col px-4">
	<div class="flex flex-1 items-center justify-center">
		<div class="w-full max-w-xl">
			<div class="mb-8 flex justify-center">
				<img src={logo} alt="" class="block h-12 w-auto" />
			</div>
			<SearchForm {onSearch} />
		</div>
	</div>

	<footer class="flex items-center justify-center gap-5 pb-6 text-[13px] text-[#707070]">
		<a href="/settings" class="transition-colors hover:text-[#c0c0c0]">settings</a>
		<a href="/about/privacy" class="transition-colors hover:text-[#c0c0c0]">privacy</a>
		<a href="/about/general" class="transition-colors hover:text-[#c0c0c0]">about</a>
	</footer>
</div>
