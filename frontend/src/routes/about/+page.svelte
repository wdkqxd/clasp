<script>
	import chevronLeftIcon from 'pixelarticons/svg/chevron-left.svg?raw';
</script>

<svelte:head>
	<title>about — clasp</title>
</svelte:head>

<div class="mx-auto w-full max-w-2xl px-4 py-16 sm:py-24">
	<a
		href="/"
		class="inline-flex items-center gap-1.5 text-[11px] text-[#3a3a3a] transition-colors hover:text-[#8a8a8a]"
	>
		<span class="inline-block size-[12px] shrink-0 align-middle [&>svg]:block [&>svg]:size-full"
			>{@html chevronLeftIcon}</span
		>
		back
	</a>

	<h1 class="mt-10 text-[15px] font-normal text-[#d0d0d0]">about</h1>

	<div class="mt-6 space-y-5 text-[13px] leading-relaxed text-[#707070]">
		<p>
			clasp is a small self-hosted metasearch engine. It queries Brave, DuckDuckGo and Startpage at the
			same time, merges what comes back, and re-ranks the results so nothing dominates the page.
		</p>
		<p>
			There is no database and no user profile. The whole thing is a single Go binary that proxies and
			parses search results in memory, which makes it cheap to run on a tiny machine and easy to audit.
		</p>
		<p>
			Open source and unmaintained by anyone in particular. Fork it, change it, host it yourself. Issues
			and pull requests are welcome on <a
				href="https://github.com/wdkqxd/clasp"
				target="_blank"
				rel="noreferrer"
				class="text-[#8a8a8a] underline underline-offset-4 transition-colors hover:text-[#c0c0c0]">github</a
			>.
		</p>
	</div>
</div>
