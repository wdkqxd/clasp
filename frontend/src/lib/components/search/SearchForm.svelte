<script lang="ts">
	import { ArrowRight, Search } from '@lucide/svelte';
	import { cubicOut } from 'svelte/easing';

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

	let input = $state<HTMLInputElement | null>(null);
	let focused = $state(false);
	let hasQuery = $derived(value.trim().length > 0);
	let showSubmit = $derived(focused && hasQuery);

	function submit(event: SubmitEvent) {
		event.preventDefault();
		const q = value.trim();
		if (q) onSearch(q);
	}

	function clear() {
		value = '';
		input?.focus();
	}

	function onKeydown(event: KeyboardEvent) {
		if (event.key !== 'Escape') return;
		if (hasQuery) {
			event.preventDefault();
			clear();
		} else {
			input?.blur();
		}
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

<form onsubmit={submit} class="w-full">
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
			onblur={() => (focused = false)}
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
</form>
