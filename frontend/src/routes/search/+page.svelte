<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { checkBang } from '$lib/bangs';
	import Infobox from '$lib/components/infobox/Infobox.svelte';
	import { wikipedia } from '$lib/api';
	import ResultsList from '$lib/components/search/ResultsList.svelte';
	import Answers from '$lib/components/rich-answers/Answers.svelte';
	import ImageGrid from '$lib/components/image/ImageGrid.svelte';
	import ImageDetailPanel from '$lib/components/image/ImageDetailPanel.svelte';
	import SearchForm from '$lib/components/search/SearchForm.svelte';
	import TabBar from '$lib/components/search/TabBar.svelte';
	import { SearchState } from '$lib/components/search/search.svelte';
	import { ImageSearchState } from '$lib/components/image/imageSearch.svelte';
	import { initBlocklist } from '$lib/components/image/aiBlocklist';
	import { Loader } from '@lucide/svelte';

	const feed = new SearchState();
	const imageFeed = new ImageSearchState();
	initBlocklist();
	let wikipediaAnswer = $state<import('$lib/types').WikipediaArticle | null>(null);
	let discussions = $state<import('$lib/types').Discussion[]>([]);
	let rich = $state<import('$lib/types').RichAnswer[]>([]);
	let startedQuery = '';
	let requestId = 0;
	let imageSentinel = $state<HTMLElement | null>(null);
	let selectedImageIndex = $state(-1);

	let activeTab = $derived(page.url.searchParams.get('type') === 'images' ? 'Images' : 'Web');

	$effect(() => {
		if (activeTab !== 'Images' || !imageSentinel) return;
		const observer = new IntersectionObserver(
			(entries) => {
				if (entries[0].isIntersecting && imageFeed.hasMore && !imageFeed.loadingMore) {
					void imageFeed.loadMore();
				}
			},
			{ rootMargin: '400px' }
		);
		observer.observe(imageSentinel);
		return () => observer.disconnect();
	});

	function search(q: string) {
		const isImages = page.url.searchParams.get('type') === 'images';
		if (isImages) {
			selectedImageIndex = -1;
			void imageFeed.run(q);
		} else {
			discussions = [];
			rich = [];
			void feed.run(q).then((res) => {
				if (q === feed.query) {
					discussions = res?.discussions ?? [];
					rich = res?.rich ?? [];
				}
				if (!res?.infobox) {
					wikipediaAnswer = null;
					return;
				}
				wikipedia(q)
					.then((article) => {
						wikipediaAnswer = article;
					})
					.catch(() => {
						wikipediaAnswer = null;
					});
			});
		}
	}

	$effect(() => {
		const q = page.url.searchParams.get('q') ?? '';
		const isImages = page.url.searchParams.get('type') === 'images';
		const tabKey = isImages ? 'images' : 'web';
		if (q === '') {
			requestId += 1;
			startedQuery = '';
			feed.clear();
			imageFeed.clear();
			wikipediaAnswer = null;
			discussions = [];
			rich = [];
			selectedImageIndex = -1;
			return;
		}
		const stateKey = `${tabKey}:${q}`;
		if (stateKey === startedQuery) return;
		startedQuery = stateKey;
		feed.query = q;
		selectedImageIndex = -1;
		discussions = [];
		rich = [];
		const currentRequest = ++requestId;
		checkBang(q).then((bang) => {
			if (currentRequest !== requestId) return;
			if (bang) {
				window.location.href = bang;
				return;
			}
			search(q);
		});
	});

	function onSearch(q: string) {
		checkBang(q).then((bang) => {
			if (bang) {
				window.location.href = bang;
				return;
			}
			const isImages = page.url.searchParams.get('type') === 'images';
			const params = new URLSearchParams({ q });
			if (isImages) params.set('type', 'images');
			goto(`/search?${params.toString()}`);
		});
	}

	function selectImage(index: number) {
		selectedImageIndex = index;
	}

	function closeDetail() {
		selectedImageIndex = -1;
	}

	function prevImage() {
		if (selectedImageIndex > 0) selectedImageIndex--;
	}

	function nextImage() {
		if (selectedImageIndex < imageFeed.items.length - 1) selectedImageIndex++;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (activeTab !== 'Images') return;
		if (e.key === 'Escape') closeDetail();
		else if (e.key === 'ArrowLeft') { e.preventDefault(); prevImage(); }
		else if (e.key === 'ArrowRight') { e.preventDefault(); nextImage(); }
	}

	let sidebarOpen = $derived(activeTab === 'Images' && selectedImageIndex >= 0 && !!imageFeed.items[selectedImageIndex]);
</script>

<svelte:head>
	<title>{feed.query ? `${feed.query}` : 'clasp'}</title>
</svelte:head>

<svelte:window onkeydown={handleKeydown} />

{#if activeTab === 'Images'}
	<div class="min-h-screen" class:pr-[420px]={sidebarOpen}>
		<div class="w-full max-w-3xl pl-4 sm:pl-8 md:pl-16 lg:pl-28 xl:pl-44 pr-4 sm:pr-6 pt-8 sm:pt-10">
			<header class="mb-8 min-h-[86px]">
				<SearchForm bind:value={feed.query} {onSearch} />
				<TabBar tabs={['Web', 'Images', 'Video']} active={activeTab} tookMs={imageFeed.tookMs} engineStatus={imageFeed.engineStatus} />
			</header>
		</div>

		<div class="w-full px-4 sm:px-6">
			<ImageGrid
				images={imageFeed.items}
				loading={imageFeed.loading}
				selectedIndex={selectedImageIndex}
				onselect={selectImage}
			/>

			{#if imageFeed.loadingMore}
				<div class="flex justify-center py-4">
					<span class="inline-block w-7 h-7 spin text-[#5c5c5c]"><Loader /></span>
				</div>
			{/if}
			<div bind:this={imageSentinel}></div>
		</div>
	</div>

	<div
		class="fixed top-0 right-0 bottom-0 w-[420px] max-w-full bg-[#0e0e10] border-l border-[#ffffff0f] z-50 overflow-y-auto transition-transform duration-250 ease-out {sidebarOpen
			? 'translate-x-0'
			: 'translate-x-full'}"
	>
		{#if sidebarOpen}
			<ImageDetailPanel
				image={imageFeed.items[selectedImageIndex]}
				index={selectedImageIndex}
				total={imageFeed.items.length}
				onclose={closeDetail}
				onprev={prevImage}
				onnext={nextImage}
			/>
		{/if}
	</div>
{:else}
	<div
		class="w-full max-w-3xl pl-4 sm:pl-8 md:pl-16 lg:pl-28 xl:pl-44 pr-4 sm:pr-6 pt-8 sm:pt-10"
	>
		<header class="mb-8 min-h-[86px]">
			<SearchForm bind:value={feed.query} {onSearch} />
			<TabBar tabs={['Web', 'Images', 'Video']} active={activeTab} tookMs={feed.tookMs} engineStatus={feed.engineStatus} />
		</header>

		<main class="space-y-12">
			<div class="relative flex flex-col md:block">
				<section class="order-last min-w-0 md:order-none">
					{#if rich.length > 0}
						<div class="mb-10">
							<Answers items={rich} />
						</div>
					{/if}
					<ResultsList feed={feed} baseIndex={2} discussions={discussions} />
				</section>
				{#if wikipediaAnswer && !feed.loading}
					<aside class="result-item order-first mt-0 mb-12 w-full min-w-0 md:absolute md:left-[calc(100%+3rem)] md:top-0 md:order-none md:mb-0 md:mt-0 md:w-[320px]" style="--i:1">
						<Infobox article={wikipediaAnswer} />
					</aside>
				{/if}
			</div>
		</main>
	</div>
{/if}

<style>
	.result-item {
		animation: resultSlideIn 0.26s cubic-bezier(0.23, 1, 0.32, 1) backwards;
		animation-delay: calc(var(--i) * 28ms);
	}

	@keyframes resultSlideIn {
		from {
			opacity: 0;
			transform: translateY(6px);
		}
	}
</style>
