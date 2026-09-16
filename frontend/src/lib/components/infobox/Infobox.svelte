<script lang="ts">
	import type { WikipediaArticle } from '../../types';

	let { article }: { article: WikipediaArticle } = $props();
	let expanded = $state(false);
	const extractLimit = 200;
	const COLLAPSED_H = 260;
	const hasLongExtract = $derived((article.extract?.length ?? 0) > extractLimit);

	let bodyRef: HTMLDivElement | undefined = $state();
	let canExpand = $state(false);

	function preview(text: string): string {
		return hasLongExtract ? `${text.slice(0, extractLimit).trimEnd()}…` : text;
	}

	$effect(() => {
		const el = bodyRef;
		if (!el) return;
		canExpand = el.scrollHeight > COLLAPSED_H;
	});

	const resetHeight = () => {
		const el = bodyRef;
		if (el) el.style.height = '';
	};

	function toggle() {
		const el = bodyRef;
		if (!el) return;
		if (expanded) {
			el.style.height = `${el.offsetHeight}px`;
			expanded = false;
			requestAnimationFrame(() => {
				requestAnimationFrame(() => {
					el.style.height = `${COLLAPSED_H}px`;
				});
			});
		} else {
			el.style.height = `${el.offsetHeight}px`;
			expanded = true;
			requestAnimationFrame(() => {
				requestAnimationFrame(() => {
					el.style.height = `${el.scrollHeight}px`;
				});
			});
		}
		el.addEventListener('transitionend', resetHeight, { once: true });
		el.addEventListener('transitioncancel', resetHeight, { once: true });
	}
</script>

<section class="infobox">
	<div class="infobox-header">
		<div class="infobox-header-left">
			<h2 class="infobox-title">
				<a href={article.url} target="_blank" rel="noreferrer">{article.title}</a>
			</h2>
			{#if article.description}
				<p class="infobox-desc">{article.description}</p>
			{/if}
		</div>
		{#if article.thumbnail}
			<div class="infobox-header-right">
				<a class="infobox-img-wrap" href={article.url} target="_blank" rel="noreferrer">
					<img class="infobox-img" src={article.thumbnail} alt="" />
				</a>
			</div>
		{/if}
	</div>

	<div
		class="infobox-body"
		bind:this={bodyRef}
		class:expanded
		class:collapsible={canExpand}
	>
		{#if article.extract}
			<p class="infobox-long-desc">
				{preview(article.extract)}
				<a class="infobox-link" href={article.url} target="_blank" rel="noreferrer">Wikipedia</a>
			</p>
		{/if}
		{#if article.facts?.length}
			<dl class="infobox-attrs">
				{#each article.facts as fact}
					<div class="infobox-attr-row">
						<dt>{fact.label}</dt>
						<dd>{fact.value}</dd>
					</div>
				{/each}
			</dl>
		{/if}
		{#if canExpand && !expanded}
			<div class="infobox-gradient"></div>
		{/if}
	</div>

	{#if canExpand}
		<button type="button" class="infobox-toggle" onclick={toggle}>
			{expanded ? 'show less' : 'show more'}
		</button>
	{/if}
</section>

<style>
	.infobox {
		padding-top: 8px;
	}

	.infobox-header {
		display: flex;
		gap: 16px;
	}

	.infobox-header-left {
		flex: 1;
		min-width: 0;
	}

	.infobox-header-right {
		display: flex;
		flex-direction: column;
		gap: 8px;
		flex-shrink: 0;
	}

	.infobox-img-wrap {
		width: 80px;
		height: 80px;
		border-radius: 4px;
		overflow: hidden;
		position: relative;
		flex-shrink: 0;
		display: block;
		line-height: 0;
	}

	.infobox-img {
		width: 100%;
		height: 100%;
		object-fit: cover;
		display: block;
		transition: opacity 0.15s;
	}

	.infobox-img-wrap:hover .infobox-img {
		opacity: 0.85;
	}

	.infobox-title {
		font-size: 1.4rem;
		margin: 0 0 6px;
		color: #d0d0d0;
		line-height: 1.3;
		font-weight: 400;
	}

	.infobox-title a {
		color: inherit;
		text-decoration: none;
	}

	.infobox-title a:hover {
		text-decoration: underline;
		text-underline-offset: 3px;
	}

	.infobox-desc {
		font-size: 0.85rem;
		color: #707070;
		line-height: 1.4;
		margin: 0 0 8px;
	}

	.infobox-body {
		position: relative;
		overflow: hidden;
		margin-top: 16px;
		padding-top: 16px;
		border-top: 1px solid #262626;
		transition: height 0.35s cubic-bezier(0.23, 1, 0.32, 1);
	}

	.infobox-body.collapsible {
		height: 260px;
	}

	.infobox-body.collapsible.expanded {
		height: auto;
	}

	.infobox-long-desc {
		font-size: 0.85rem;
		color: #a0a0a0;
		line-height: 1.5;
		margin: 0;
	}

	.infobox-link {
		color: #707070;
		text-decoration: none;
		margin-left: 4px;
	}

	.infobox-link:hover {
		text-decoration: underline;
		text-underline-offset: 3px;
		color: #a0a0a0;
	}

	.infobox-attrs {
		margin: 14px 0 0;
		padding-top: 14px;
		border-top: 1px solid #262626;
		font-size: 0.85rem;
		line-height: 1.5;
		display: grid;
		grid-template-columns: fit-content(40%) 1fr;
		column-gap: 16px;
		row-gap: 7px;
		align-items: baseline;
	}

	.infobox-attr-row {
		display: contents;
	}

	.infobox-attr-row dt {
		color: #a0a0a0;
		font-weight: 500;
	}

	.infobox-attr-row dd {
		margin: 0;
		color: #707070;
	}

	.infobox-gradient {
		position: absolute;
		bottom: 0;
		left: 0;
		right: 0;
		height: 3em;
		background: linear-gradient(to bottom, transparent, #101010);
		pointer-events: none;
	}

	.infobox-toggle {
		background: transparent;
		font-size: 0.8rem;
		padding: 6px 12px;
		cursor: pointer;
		margin-top: 12px;
		font-family: inherit;
		color: #a0a0a0;
		width: 100%;
		border: 1px solid #363636;
		border-radius: 3px;
		transition:
			border-color 0.15s,
			color 0.15s,
			transform 0.15s cubic-bezier(0.23, 1, 0.32, 1);
	}

	.infobox-toggle:hover {
		border-color: #4a4a4a;
		color: #d0d0d0;
	}

	.infobox-toggle:active {
		transform: scale(0.98);
		transition-duration: 0.1s;
	}
</style>