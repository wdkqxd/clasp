<script>
	import { ChevronLeft } from '@lucide/svelte';
</script>

<svelte:head>
	<title>about / clasp</title>
</svelte:head>

<div class="mx-auto w-full max-w-2xl px-4 py-16 sm:py-24">
	<a
		href="/"
		class="inline-flex items-center gap-1.5 text-[11px] text-[#3a3a3a] transition-colors hover:text-[#8a8a8a]"
	>
		<span class="inline-block size-[12px] shrink-0 align-middle [&>svg]:block [&>svg]:size-full"
			><ChevronLeft /></span
		>
		back
	</a>

	<h1 class="mt-10 text-[20px] leading-tight font-normal text-[#d0d0d0]">about</h1>

	<div class="mt-6 space-y-8 text-[13px] leading-relaxed text-[#707070]">
		<section>
			<h2 class="mb-1.5 text-[15px] leading-snug font-normal text-[#c0c0c0]">what it does</h2>
			<p>
				clasp is a small open-source metasearch engine. it asks brave, duckduckgo and startpage at the
				same time, merges what comes back and re-ranks it, so no single engine takes over the page.
			</p>
		</section>

		<section>
			<h2 class="mb-1.5 text-[15px] leading-snug font-normal text-[#c0c0c0]">how it works</h2>
			<p>
				there is no database and no user account. the whole thing is one go binary that fetches and
				parses results in memory, which makes it cheap to run on a tiny machine and easy to read.
			</p>
		</section>

		<section>
			<h2 class="mb-1.5 text-[15px] leading-snug font-normal text-[#c0c0c0]">source</h2>
			<p>
				the code is open and nobody in particular maintains it. fork it, change it, run it yourself.
				issues and pull requests are welcome on <a
					href="https://github.com/wdkqxd/clasp"
					target="_blank"
					rel="noreferrer"
					class="text-[#8a8a8a] underline underline-offset-4 transition-colors hover:text-[#c0c0c0]"
					>github</a
				>.
			</p>
		</section>
	</div>
</div>
