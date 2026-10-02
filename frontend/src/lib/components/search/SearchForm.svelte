<script lang="ts">
	import { ArrowRight, Search, X } from '@lucide/svelte';
	import { cubicIn, cubicOut } from 'svelte/easing';
	import { fly } from 'svelte/transition';
	import { suggest } from '../../api';
	import { bangIcon, parseBang, resolveBang, searchBangs, type BangMatch } from '../../bangs';
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
	let bangResults = $state<BangMatch[]>([]);
	let activeIndex = $state(-1);
	let suggestionsOpen = $state(false);
	let brokenIcons = $state<Record<string, true>>({});
	let requestId = 0;

	let bang = $state<BangMatch | null>(null);
	let parsed = $derived(parseBang(value));
	let spaced = $derived(/^!\S+\s/.test(value.trimStart()));
	let chip = $derived(
		bang && spaced && parsed && parsed.trigger === bang.trigger ? bang : null
	);
	let rest = $derived(chip && parsed ? parsed.rest : value);
	let bangMode = $derived(!chip && value.trim().startsWith('!'));
	let itemCount = $derived(bangResults.length || suggestions.length);

	let showSuggestions = $derived(focused && suggestionsOpen && itemCount > 0);
	let showBangs = $derived(showSuggestions && bangResults.length > 0);

	$effect(() => {
		const query = value.trim();
		let cancelled = false;
		void resolveBang(query).then((found) => {
			if (!cancelled) bang = found;
		});
		return () => {
			cancelled = true;
		};
	});

	$effect(() => {
		const query = value.trim();
		const engine = settings.autocomplete;
		if (!focused) return;

		if (bangMode) {
			const term = query.slice(1).split(/\s/)[0];
			if (!term) {
				closeSuggestions();
				return;
			}
			const timer = setTimeout(() => {
				void loadBangs(term);
			}, DEBOUNCE_MS);
			return () => clearTimeout(timer);
		}

		const text = chip ? rest.trim() : query;
		if (text.length < MIN_QUERY_LENGTH) {
			closeSuggestions();
			return;
		}

		const timer = setTimeout(() => {
			void loadSuggestions(text, engine);
		}, DEBOUNCE_MS);
		return () => clearTimeout(timer);
	});

	async function loadSuggestions(query: string, engine: string) {
		const id = ++requestId;
		try {
			const res = await suggest(query, engine);
			if (id !== requestId) return;
			bangResults = [];
			suggestions = res.suggestions.filter((s) => s.toLowerCase() !== query.toLowerCase());
			activeIndex = -1;
			suggestionsOpen = suggestions.length > 0;
		} catch {
			if (id !== requestId) return;
			closeSuggestions();
		}
	}

	async function loadBangs(term: string) {
		const id = ++requestId;
		try {
			const res = await searchBangs(term);
			if (id !== requestId) return;
			suggestions = [];
			bangResults = res;
			activeIndex = -1;
			suggestionsOpen = res.length > 0;
		} catch {
			if (id !== requestId) return;
			closeSuggestions();
		}
	}

	function closeSuggestions() {
		requestId += 1;
		suggestions = [];
		bangResults = [];
		activeIndex = -1;
		suggestionsOpen = false;
	}

	function onInput(event: Event) {
		const text = (event.currentTarget as HTMLInputElement).value;
		value = chip && text ? `!${chip.trigger} ${text}` : chip ? `!${chip.trigger}` : text;
	}

	function clearBang() {
		if (!chip) return;
		value = `!${chip.trigger}`;
		closeSuggestions();
		input?.focus();
	}

	function pick(suggestion: string) {
		const q = chip ? `!${chip.trigger} ${suggestion}` : suggestion;
		value = q;
		closeSuggestions();
		input?.focus();
		onSearch(q);
	}

	function pickBang(target: BangMatch) {
		const tail = value.trim().slice(1).replace(/^\S*/, '').trim();
		value = `!${target.trigger} ${tail}`;
		closeSuggestions();
		input?.focus();
	}

	function matchParts(suggestion: string): [string, string] {
		const query = chip ? rest.trim() : value.trim();
		if (!query) return ['', suggestion];
		const q = query.toLowerCase();
		const s = suggestion.toLowerCase();
		const max = Math.min(q.length, s.length);
		let len = 0;
		while (len < max && q[len] === s[len]) len++;
		if (len === 0) return ['', suggestion];
		return [suggestion.slice(0, len), suggestion.slice(len)];
	}

	function moveActive(step: number) {
		if (!showSuggestions) return;
		activeIndex = (activeIndex + step + itemCount) % itemCount;
	}

	function onKeydown(event: KeyboardEvent) {
		if (event.key === 'Backspace' && chip && input?.selectionStart === 0) {
			event.preventDefault();
			clearBang();
			return;
		}
		if (event.key === 'ArrowDown' || event.key === 'ArrowUp') {
			if (!showSuggestions) return;
			event.preventDefault();
			event.stopPropagation();
			moveActive(event.key === 'ArrowDown' ? 1 : -1);
			return;
		}
		if (event.key === 'Enter' && showSuggestions && activeIndex >= 0) {
			event.preventDefault();
			if (showBangs) pickBang(bangResults[activeIndex]);
			else pick(suggestions[activeIndex]);
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
		{#if chip}
			<span
				class="flex shrink-0 items-center gap-1.5 rounded-sm bg-[#1e1e1e] py-1 pl-1 pr-1.5 text-[12px] leading-none text-[#c0c0c0]"
			>
				<img
					src={bangIcon(chip.host)}
					alt=""
					width="14"
					height="14"
					loading="lazy"
					class="h-3.5 w-3.5 rounded-sm object-contain opacity-75"
					onerror={() => (brokenIcons[chip.host] = true)}
				/>
				<span class="max-w-[120px] truncate">{chip.name}</span>
				<button
					type="button"
					aria-label="Remove bang"
					onmousedown={(event) => event.preventDefault()}
					onclick={clearBang}
					class="shrink-0 rounded-sm p-0.5 text-[#5c5c5c] transition-colors hover:bg-[#262626] hover:text-[#c0c0c0]"
				>
					<X size={11} />
				</button>
			</span>
		{/if}
		<input
			bind:this={input}
			value={rest}
			oninput={onInput}
			type="search"
			aria-label={placeholder}
			autocomplete="off"
			spellcheck="false"
			enterkeyhint="search"
			placeholder={chip ? 'search' : placeholder}
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
		<div
			data-suggestions
			in:fly={{ y: -4, duration: 140, easing: cubicOut }}
			out:fly={{ y: -4, duration: 110, easing: cubicIn }}
			class="absolute inset-x-0 top-full z-30 mt-2 overflow-hidden rounded-md border border-[#242424] bg-[#121212] p-1.5 shadow-2xl shadow-black/40"
		>
			<ul
				role="listbox"
				aria-label={showBangs ? 'bang suggestions' : 'search suggestions'}
				class={showBangs ? 'grid grid-cols-2 gap-1' : ''}
			>
			{#if showBangs}
				{#each bangResults as bang, i (bang.trigger + bang.url)}
					<li>
						<button
							type="button"
							role="option"
							aria-selected={i === activeIndex}
							onmousedown={(event) => {
								event.preventDefault();
							}}
							onclick={() => pickBang(bang)}
							onmouseenter={() => (activeIndex = i)}
							onmouseleave={() => (activeIndex = -1)}
							class="flex w-full items-center gap-2.5 rounded-md px-2.5 py-2 text-left text-[13px] transition-all duration-100 ease-out active:scale-[0.99] {i ===
							activeIndex
								? 'bg-[#1e1e1e] text-[#d0d0d0]'
								: 'text-[#8a8a8a]'}"
						>
							{#if !brokenIcons[bang.host]}
								<img
									src={bangIcon(bang.host)}
									alt=""
									width="14"
									height="14"
									loading="lazy"
									class="h-3.5 w-3.5 shrink-0 rounded-sm object-contain opacity-75"
									onerror={() => (brokenIcons[bang.host] = true)}
								/>
							{/if}
							<span class="min-w-0 flex-1 truncate text-[#c0c0c0]">{bang.name}</span>
							<span class="shrink-0 truncate text-[#8a8a8a]">!{bang.trigger}</span>
						</button>
					</li>
				{/each}
			{:else}
				{#each suggestions as suggestion, i (suggestion)}
					{@const [match, rest] = matchParts(suggestion)}
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
							<span class="truncate">{match}{#if rest}<span class={i === activeIndex
										? 'text-[#e8e8e8]'
										: 'text-[#c0c0c0]'}>{rest}</span
									>{/if}</span>
						</button>
					</li>
				{/each}
			{/if}
			</ul>
			{#if showBangs}
				<div class="mt-1.5 border-t border-[#1e1e1e]">
					<a
						href="https://helium.computer/bangs"
						target="_blank"
						rel="noreferrer"
						onmousedown={(event) => {
							event.preventDefault();
						}}
						class="flex items-center justify-center gap-1.5 px-2.5 pt-2 text-[11px] leading-none text-[#4a4a4a] transition-colors hover:text-[#8a8a8a]"
					>
						<span>bangs by</span>
						<svg viewBox="0 0 256 256" aria-hidden="true" class="size-3 shrink-0">
							<path
								d="M152.62 237.85L128 256L103.38 237.85L114.568 151.248L45.05 204.22L17 192L20.425 161.635L101.128 128L20.425 94.365L17 64L45.05 51.78L114.568 104.752L103.38 18.15L128 0L152.62 18.15L141.432 104.752L210.95 51.78L239 64L235.575 94.365L154.872 128L235.575 161.635L239 192L210.95 204.22L141.432 151.248L152.62 237.85Z"
								fill="currentColor"
							/>
						</svg>
						<span>helium</span>
					</a>
				</div>
			{/if}
		</div>
	{/if}
</form>
