<script>
	import { ChevronLeft } from '@lucide/svelte';
</script>

<svelte:head>
	<title>privacy / clasp</title>
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

	<h1 class="mt-10 text-[20px] leading-tight font-normal text-[#d0d0d0]">privacy</h1>

	<div class="mt-6 space-y-8 text-[13px] leading-relaxed text-[#707070]">
		<section>
			<h2 class="mb-1.5 text-[15px] leading-snug font-normal text-[#c0c0c0]">general terms</h2>
			<p>
			we believe privacy is important, so we don't collect, store, or sell any data related to you.
			</p>
		</section>

		<section>
			<h2 class="mb-1.5 text-[15px] leading-snug font-normal text-[#c0c0c0]">analytics</h2>
			<p>
				we look at anonymous, aggregated numbers through described <a href="https://analytics.wdkq.cc" target="_blank" rel="noreferrer" class="text-[#8a8a8a] underline underline-offset-4 transition-colors hover:text-[#c0c0c0]"
				>our own plausible instance</a> it sets no cookies
				and builds no fingerprint, so a single visit cannot be traced back to you. 

				<br /><br />if you would rather we
				counted nothing at all, you can switch it off in the settings at any time. how it works is
				described <a
					href="https://plausible.io/privacy-focused-web-analytics"
					target="_blank"
					rel="noreferrer"
					class="text-[#8a8a8a] underline underline-offset-4 transition-colors hover:text-[#c0c0c0]"
					>here</a
				>.
			</p>
		</section>

		<section>
			<h2 class="mb-1.5 text-[15px] leading-snug font-normal text-[#c0c0c0]">queries</h2>
			<p>
				your query travels from your browser to this server, and from here to the search engines you
				turned on. it is never written to disk and never handed to anyone else. 

				<br /><br />all queries are processed in RAM and sent to search engines in anonymized form, so searches cannot be traced back to anyone.
			</p>
		</section>

		<section>
			<h2 class="mb-1.5 text-[15px] leading-snug font-normal text-[#c0c0c0]">self-hosted instances</h2>
			<p>
				we do our best to keep this instance private. anyone can run their own copy of clasp though, and
				we can neither see nor control what they do with it, so we cannot answer for it. <br /><br />read the
				source before you trust one. it is all on <a
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
