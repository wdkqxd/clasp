<script lang="ts">
	import type { RichTopAnswer } from '../../types';
	import { getSearchSettings } from '../settings/searchSettings.svelte';
	import { ChevronUp } from '@lucide/svelte';
	import CopyButton from './CopyButton.svelte';

	let { topAnswer }: { topAnswer: RichTopAnswer } = $props();

	const settings = getSearchSettings();

	const blocks = $derived(topAnswer.blocks ?? []);

	let copied = $state(-1);
	function onCopy(index: number, value: string) {
		navigator.clipboard?.writeText(value);
		copied = index;
		setTimeout(() => {
			if (copied === index) copied = -1;
		}, 1200);
	}

	const host = $derived(
		topAnswer.meta_url?.hostname ||
			(() => {
				try {
					return new URL(topAnswer.url).hostname.replace(/^www\./, '');
				} catch {
					return '';
				}
			})()
	);
	const favicon = $derived(host ? `https://icons.duckduckgo.com/ip3/${host}.ico` : '');
</script>

<section class="rounded-md border border-[#262626] bg-[#161616] px-4 py-3.5">
	<div class="flex items-start gap-2.5">
		<h2 class="text-[15px] leading-snug font-medium text-[#e2e2e2] line-clamp-3">{topAnswer.question}</h2>
	</div>

	{#if blocks.length > 0}
		<div class="mt-3 space-y-2.5 border-l border-[#2a2a2a] pl-3.5">
			{#each blocks as b, i (i)}
				{#if b.kind === 'code' && b.text}
					<div class="group relative">
						<pre class="overflow-x-auto rounded bg-[#0e0e0e] py-2 pr-9 pl-3"><code class="text-[12.5px] leading-relaxed text-[#d0d0d0]">{b.text}</code></pre>
						<button
							type="button"
							title="copy code"
							aria-label="copy code"
							onclick={() => onCopy(i, b.text ?? '')}
							class="absolute top-1.5 right-1.5 flex cursor-pointer items-center rounded p-1 {copied === i
								? '[&>span]:opacity-100'
								: ''}"
						>
							<CopyButton active={copied === i} size="h-3.5 w-3.5" />
						</button>
					</div>
				{:else if b.kind === 'list' && b.items?.length}
					<ul class="list-disc space-y-1 pl-5">
						{#each b.items as item, j (j)}
							<li class="text-[13.5px] leading-relaxed text-[#b0b0b0]">{item}</li>
						{/each}
					</ul>
				{:else if b.kind === 'quote' && b.text}
					<blockquote class="border-l-2 border-[#333333] pl-3 text-[13.5px] leading-relaxed text-[#8f8f8f] italic">
						{b.text}
					</blockquote>
				{:else if b.text}
					<p class="max-w-[68ch] text-[13.5px] leading-relaxed text-[#b0b0b0]">{b.text}</p>
				{/if}
			{/each}
		</div>
	{/if}

	<div class="mt-3.5 flex items-center gap-2 text-[12px] text-[#5c5c5c]">
		{#if topAnswer.author}
			<span class="min-w-0 truncate">{topAnswer.author}</span>
		{/if}
		{#if topAnswer.upvotes}
			<span class="flex shrink-0 items-center gap-1">
				<span class="icon text-[#5c5c5c]"><ChevronUp /></span>
				<span>{topAnswer.upvotes.toLocaleString()}</span>
			</span>
		{/if}
		{#if host}
			<a
				href={topAnswer.url}
				target={settings.openInNewTab ? '_blank' : undefined}
				rel={settings.openInNewTab ? 'noopener noreferrer' : undefined}
				class="sitelink ml-auto flex min-w-0 shrink items-center gap-1.5 hover:text-[#8f8f8f]"
			>
				{#if favicon}
					<img
						src={favicon}
						alt=""
						width="14"
						height="14"
						loading="lazy"
						class="h-3.5 w-3.5 rounded-sm object-contain opacity-75"
					/>
				{/if}
				<span class="truncate">{host}</span>
			</a>
		{/if}
	</div>
</section>
