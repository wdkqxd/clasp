<script lang="ts">
	import { ArrowRight, Search } from '@lucide/svelte';
	import { cubicIn, cubicOut } from 'svelte/easing';
	import { fly } from 'svelte/transition';
	import { suggest } from '../../api';
	import { getSearchSettings } from '../settings/searchSettings.svelte';

	type ArrowTransition = { duration?: number };

	function arrowIn(node: Element, { duration = 190 }: ArrowTransition = {}) {
		return {
			duration,
			easing: cubicOut,
			css: (t: number) =>
				`opacity: ${t}; filter: blur(${(1 - t) * 5}px); transform: scale(${0.1 + 0.9 * t});`
		};
	}

	function arrowOut(node: Element, { duration = 110 }: ArrowTransition = {}) {
		return {
			duration,
			easing: cubicOut,
			css: (t: number) =>
				`opacity: ${t}; filter: blur(${(1 - t) * 4}px); transform: scale(${0.35 + 0.65 * t});`
		};
	}

	let {
		value = $bindable(''),
		placeholder = 'search the web',
		onSearch
	}: {
		value?: string;
		placeholder?: string;
		onSearch: (q: string) => void;
	} = $props();

	const MIN_QUERY_LENGTH = 2;
	const DEBOUNCE_MS = 140;

	const settings = getSearchSettings();

	let input = $state<HTMLInputElement | null>(null);
	let focused = $state(false);
	let hasQuery = $derived(value.trim().length > 0);
	let showSubmit = $derived(focused && hasQuery);

	let suggestions = $state<string[]>([]);
	let activeIndex = $state(-1);
	let suggestionsOpen = $state(false);
	let requestId = 0;

	let showSuggestions = $derived(focused && suggestionsOpen && suggestions.length > 0);

	$effect(() => {
		const query = value.trim();
		const engine = settings.autocomplete;
		if (!focused || query.length < MIN_QUERY_LENGTH) return;

		const timer = setTimeout(() => {
			void loadSuggestions(query, engine);
		}, DEBOUNCE_MS);
		return () => clearTimeout(timer);
	});

	$effect(() => {
		if (value.trim().length < MIN_QUERY_LENGTH) closeSuggestions();
	});

	async function loadSuggestions(query: string, engine: string) {
		const id = ++requestId;
		try {
			const res = await suggest(query, engine);
			if (id !== requestId) return;
			suggestions = res.suggestions.filter((s) => s.toLowerCase() !== query.toLowerCase());
			activeIndex = -1;
			suggestionsOpen = suggestions.length > 0;
		} catch {
			if (id !== requestId) return;
			closeSuggestions();
		}
	}

	function closeSuggestions() {
		requestId += 1;
		suggestions = [];
		activeIndex = -1;
		suggestionsOpen = false;
	}

	function pick(suggestion: string) {
		value = suggestion;
		closeSuggestions();
		input?.focus();
		onSearch(suggestion);
	}

	function moveActive(step: number) {
		if (!showSuggestions) return;
		activeIndex = (activeIndex + step + suggestions.length) % suggestions.length;
	}

	function onKeydown(event: KeyboardEvent) {
		if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
			if (!showSuggestions) return;
			event.preventDefault();
			event.stopPropagation();
			moveActive(event.key === 'ArrowDown' ? 1 : -1);
			return;
		}
		if (event.key === 'Enter' && showSuggestions && activeIndex >= 0) {
			event.preventDefault();
			pick(suggestions[activeIndex]);
			return;
		}
		if (event.key !== 'Escape') return;
		if (showSuggestions) {
			event.preventDefault();
			closeSuggestions();
			return;
		}
		if (hasQuery) {
			event.preventDefault();
			clear();
		} else {
			input?.blur();
		}
	}

	function submit(event: SubmitEvent) {
		event.preventDefault();
		const q = value.trim();
		if (!q) return;
		closeSuggestions();
		onSearch(q);
	}

	function clear() {
		value = '';
		input?.focus();
	}

	$effect(() => {
		const focusFromShortcut = (event: KeyboardEvent) => {
			if (event.key !== '/' || event.ctrlKey || event.metaKey || event.altKey) return;
			const target = event.target as HTMLElement | null;
			if (
				target &&
				(target.tagName === 'INPUT' ||
					target.tagName === 'TEXTAREA' ||
					target.isContentEditable)
			) {
				return;
			}
			event.preventDefault();
			input?.focus();
			input?.select();
		};
		window.addEventListener('keydown', focusFromShortcut);
		return () => window.removeEventListener('keydown', focusFromShortcut);
	});
</script>

<form onsubmit={submit} class="relative w-full">
	<div
		class="flex items-center gap-2.5 rounded-md border border-[#242424] bg-[#121212] pl-3 transition-colors duration-200 ease-out focus-within:border-[#3d3d3d] {hasQuery
			? 'pr-1.5'
			: 'pr-3'}"
	>
		<Search size={16} class="shrink-0 text-[#454545]" />
		<input
			bind:this={input}
			bind:value
			type="search"
			aria-label={placeholder}
			autocomplete="off"
			spellcheck="false"
			enterkeyhint="search"
			{placeholder}
			onkeydown={onKeydown}
			onfocus={() => (focused = true)}
			onblur={() => {
				focused = false;
				closeSuggestions();
			}}
			class="min-w-0 flex-1 bg-transparent py-2 text-[15px] text-[#c0c0c0] caret-[#8a8a8a] transition-colors placeholder:text-[#454545] focus:outline-none [&::-webkit-search-cancel-button]:appearance-none"
		/>
		{#if showSubmit}
			<button
				type="submit"
				aria-label="Search"
				onmousedown={(event) => event.preventDefault()}
				in:arrowIn
				out:arrowOut
				class="shrink-0 rounded-sm bg-[#242424] p-1.5 text-[#c8c8c8] transition-colors duration-150 hover:bg-[#2e2e2e] hover:text-[#e8e8e8] active:scale-90"
			>
				<ArrowRight size={13} />
			</button>
		{/if}
	</div>

	{#if showSuggestions}
		<ul
			data-suggestions
			role="listbox"
			aria-label="search suggestions"
			in:fly={{ y: -4, duration: 140, easing: cubicOut }}
			out:fly={{ y: -4, duration: 110, easing: cubicIn }}
			class="absolute inset-x-0 top-full z-30 mt-2 overflow-hidden rounded-md border border-[#242424] bg-[#121212] p-1.5 shadow-2xl shadow-black/40"
		>
			{#each suggestions as suggestion, i (suggestion)}
				<li>
					<button
						type="button"
						role="option"
						aria-selected={i === activeIndex}
						onmousedown={(event) => {
							event.preventDefault();
						}}
						onclick={() => pick(suggestion)}
						onmouseenter={() => (activeIndex = i)}
						onmouseleave={() => (activeIndex = -1)}
						class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-2 text-left text-[13px] transition-all duration-100 ease-out active:scale-[0.99] {i ===
						activeIndex
							? 'bg-[#1e1e1e] text-[#d0d0d0]'
							: 'text-[#8a8a8a]'}"
					>
						<Search size={13} class="shrink-0 text-[#4a4a4a]" />
						<span class="truncate">{suggestion}</span>
					</button>
				</li>
			{/each}
		</ul>
	{/if}
</form>
