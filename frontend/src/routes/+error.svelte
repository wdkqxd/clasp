<script lang="ts">
	import { page } from '$app/state';

	const TITLES: Record<number, string> = {
		404: 'not found',
		403: 'forbidden',
		500: 'internal error'
	};

	const title = $derived(TITLES[page.status] ?? 'something went wrong');
</script>

<svelte:head>
	<title>{page.status} / clasp</title>
</svelte:head>

<div class="flex min-h-screen flex-col px-4">
	<div class="flex flex-1 flex-col items-center justify-center gap-8">

		<div class="text-center">
			<h1 class="mt-3 text-[22px] leading-tight font-normal text-[#d0d0d0]">{title}</h1>
			{#if page.error?.message && page.status !== 404}
				<p class="mt-3 text-[13px] break-words text-[#707070]">{page.error.message}</p>
			{/if}
			<p class="mt-3 text-[13px] text-[#707070]">
				{page.status === 404
					? 'this page does not exist, or it moved somewhere else.'
					: 'something broke on our side. try again in a moment.'}
			</p>
		</div>
	</div>
</div>
