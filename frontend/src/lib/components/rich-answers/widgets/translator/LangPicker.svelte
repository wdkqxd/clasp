<script lang="ts">
	import { Check, Globe } from '@lucide/svelte';
	import { AUTO, LANGUAGES, flagUrl, langByCode, matchesLanguage } from './languages';

	let {
		value = $bindable('auto'),
		detect = false,
		detected = null,
		align = 'start',
		onchange
	}: {
		value?: string;
		detect?: boolean;
		detected?: string | null;
		align?: 'start' | 'end';
		onchange?: (code: string) => void;
	} = $props();

	const entries = $derived(detect ? [AUTO, ...LANGUAGES] : LANGUAGES);

	let open = $state(false);
	let query = $state('');
	let active = $state(0);
	let root = $state<HTMLElement | null>(null);
	let searchEl = $state<HTMLInputElement | null>(null);
	let listEl = $state<HTMLElement | null>(null);

	const filtered = $derived.by(() => {
		const q = query.trim().toLowerCase();
		if (!q) return entries;
		return entries.filter((l) => matchesLanguage(l, q));
	});

	const shown = $derived(detect && value === 'auto' && detected ? langByCode(detected) : null);

	function openPop() {
		open = true;
		query = '';
		active = Math.max(
			0,
			filtered.findIndex((l) => l.code === value)
		);
		queueMicrotask(() => {
			searchEl?.focus();
			scrollActive();
		});
	}

	function close() {
		open = false;
		query = '';
	}

	function pick(code: string) {
		value = code;
		close();
		onchange?.(code);
	}

	function scrollActive() {
		listEl?.children[active]?.scrollIntoView({ block: 'nearest' });
	}

	function onListKey(e: KeyboardEvent) {
		if (e.key === 'Escape') {
			e.stopPropagation();
			close();
		} else if (e.key === 'ArrowDown') {
			e.preventDefault();
			active = Math.min(active + 1, filtered.length - 1);
			scrollActive();
		} else if (e.key === 'ArrowUp') {
			e.preventDefault();
			active = Math.max(active - 1, 0);
			scrollActive();
		} else if (e.key === 'Enter') {
			e.preventDefault();
			const hit = filtered[active];
			if (hit) pick(hit.code);
		}
	}

	function onWindowPointerDown(e: PointerEvent) {
		if (open && root && !root.contains(e.target as Node)) close();
	}
</script>

<svelte:window onpointerdown={onWindowPointerDown} />

<div bind:this={root} class="relative" class:open>
	<button
		type="button"
		class="flex min-h-10 w-full max-w-full cursor-pointer items-center gap-2 rounded-[5px] border border-transparent bg-[#1a1a1a] px-2.5 py-2 text-[13px] text-[#d0d0d0] transition-[background-color,border-color] duration-100 ease-out hover:bg-[#212121] {open
			? 'border-[#333] bg-[#212121]'
			: ''}"
		aria-haspopup="listbox"
		aria-expanded={open}
		onclick={() => (open ? close() : openPop())}
	>
		<span class="flex min-w-0 items-center gap-2 overflow-hidden">
			{#if shown}
				{@const shownFlag = flagUrl(shown.code)}
				{#if shownFlag}
					<img class="size-[17px] shrink-0 rounded-full object-cover opacity-75" src={shownFlag} alt="" width="17" height="17" loading="lazy" />
				{:else}
					<span class="inline-flex shrink-0 text-[#707070]"><Globe size={16} /></span>
				{/if}
			{:else if value === 'auto'}
				<span class="inline-flex shrink-0 text-[#707070]"><Globe size={16} /></span>
			{:else if flagUrl(value)}
				<img
					class="size-[17px] shrink-0 rounded-full object-cover opacity-75"
					src={flagUrl(value)}
					alt=""
					width="20"
					height="20"
					loading="lazy"
				/>
			{/if}
			<span class="truncate">
				{#if value === 'auto'}
					{shown ? shown.name : AUTO.name}
					{#if shown}<span class="ml-0.5 text-[10px] text-[#5c5c5c]">· detected</span>{/if}
				{:else}
					{langByCode(value)?.name ?? value}
				{/if}
			</span>
		</span>
		<svg
			class="ml-auto inline-flex shrink-0 text-[#5c5c5c] transition-transform duration-200 ease-out {open ? 'rotate-180' : ''}"
			width="14"
			height="14"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<path d="M6 9l6 6l6 -6" />
		</svg>
	</button>

	<div
		class="absolute top-full left-0 z-30 mt-1.5 w-[min(280px,calc(100vw-40px))] origin-top-left rounded-[5px] border border-[#292929] bg-[#121212] shadow-2xl shadow-black/30 transition-all duration-300 ease-[cubic-bezier(0.22,1,0.36,1)] {align ===
		'end'
			? 'right-0 left-auto origin-top-right'
			: ''} {open
			? 'visible translate-y-0 scale-100 opacity-100 blur-0'
			: 'invisible pointer-events-none -translate-y-1 scale-95 opacity-0 blur-sm'}"
		role="listbox"
	>
			<input
				bind:this={searchEl}
				bind:value={query}
				type="text"
				placeholder="search languages…"
				autocomplete="off"
				spellcheck="false"
				class="w-full border-b border-[#262626] bg-transparent px-3 py-2.5 text-[13px] text-[#d0d0d0] outline-none placeholder:text-[#5c5c5c]"
				oninput={() => (active = 0)}
				onkeydown={onListKey}
			/>
			<div
				bind:this={listEl}
				class="relative flex max-h-[260px] flex-col gap-0.5 overflow-y-auto p-1.5 lang-list [scrollbar-width:none]"
			>
				{#each filtered as lang, i (lang.code)}
					<button
						type="button"
						role="option"
						aria-selected={lang.code === value}
						class="flex w-full shrink-0 cursor-pointer items-center gap-2.5 rounded-[5px] border-none bg-transparent px-2.5 py-1.5 text-left text-[13px] text-[#c0c0c0] transition-colors duration-100 hover:bg-[#1e1e1e] {i === active
							? 'bg-[#1e1e1e]'
							: ''}"
						onclick={() => pick(lang.code)}
						onmouseenter={() => (active = i)}
					>
						{#if lang.code === 'auto'}
							<span class="inline-flex shrink-0 text-[#707070]"><Globe size={16} /></span>
						{:else}
							{@const itemFlag = flagUrl(lang.code)}
							{#if itemFlag}
								<img class="size-[17px] shrink-0 rounded-full object-cover opacity-75" src={itemFlag} alt="" width="17" height="17" loading="lazy" />
							{:else}
								<span class="inline-flex shrink-0 text-[#707070]"><Globe size={16} /></span>
							{/if}
						{/if}
						<span class="min-w-0 flex-1 truncate">{lang.name}</span>
						{#if lang.code === value}
							<Check size={13} class="shrink-0 text-[#5c5c5c]" />
						{/if}
					</button>
				{/each}
				{#if filtered.length === 0}
					<p class="px-2.5 py-2 text-[12px] text-[#5c5c5c]">no matches</p>
				{/if}
			</div>
		</div>
</div>

<style>
	.lang-list::-webkit-scrollbar {
		display: none;
	}
</style>
