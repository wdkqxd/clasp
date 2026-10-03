<script lang="ts">
	import type { WikipediaArticle } from '../../types';

	let { article }: { article: WikipediaArticle } = $props();
	const extractLimit = 200;
	const hasLongExtract = $derived((article.extract?.length ?? 0) > extractLimit);

	function preview(text: string): string {
		return hasLongExtract ? `${text.slice(0, extractLimit).trimEnd()}…` : text;
	}
</script>

<section class="pt-2">
	<div class="flex gap-4">
		<div class="min-w-0 flex-1">
			<h2 class="m-0 mb-1.5 text-[1.4rem] leading-[1.3] font-normal text-[#d0d0d0]">
				<a
					href={article.url}
					target="_blank"
					rel="noreferrer"
					class="text-inherit no-underline hover:underline hover:underline-offset-[3px]"
				>
					{article.title}
				</a>
			</h2>
			{#if article.description}
				<p class="m-0 mb-2 text-[0.85rem] leading-[1.4] text-[#707070]">{article.description}</p>
			{/if}
		</div>
		{#if article.thumbnail}
			<div class="flex flex-shrink-0 flex-col gap-2">
				<a
					class="group relative block h-20 w-20 flex-shrink-0 overflow-hidden rounded-sm leading-[0]"
					href={article.url}
					target="_blank"
					rel="noreferrer"
				>
					<img
						class="block h-full w-full object-cover transition-opacity duration-150 group-hover:opacity-[0.85]"
						src={article.thumbnail}
						alt=""
					/>
				</a>
			</div>
		{/if}
	</div>

	{#if article.extract}
		<p class="m-0 mt-4 border-t border-[#262626] pt-4 text-[0.85rem] leading-[1.5] text-[#a0a0a0]">
			{preview(article.extract)}
			<a
				class="ml-1 text-[#707070] no-underline transition-colors duration-150 hover:text-[#a0a0a0] hover:underline hover:underline-offset-[3px]"
				href={article.url}
				target="_blank"
				rel="noreferrer"
			>Wikipedia</a>
		</p>
	{/if}
</section>
