<script lang="ts">
	import { ArrowLeftRight, Volume2, X } from '@lucide/svelte';
	import CopyButton from '../../shared/CopyButton.svelte';
	import LangPicker from './LangPicker.svelte';
	import { browserLanguage, langByCode } from './languages';
	import { parseTranslateQuery } from './translateIntent';

	let { query = '' }: { query?: string } = $props();

	const MAX = 1000;

	const show = $derived(parseTranslateQuery(query));

	let text = $state('');
	let from = $state('auto');
	let to = $state('en');

	$effect(() => {
		to = browserLanguage();
	});
	let swapping = $state(false);

	const sourceLang = $derived(from === 'auto' ? null : from);
	const sourceRtl = $derived(!!langByCode(sourceLang ?? '')?.rtl);
	const targetRtl = $derived(!!langByCode(to)?.rtl);

	function onSwap() {
		if (!sourceLang) return;
		const a = sourceLang;
		from = to;
		to = a;
		swapping = !swapping;
	}

	function onClear() {
		text = '';
	}
</script>

{#if show}
	<div class="@container rounded-md border border-[#292929] bg-[#121212]">
		<div class="flex items-center gap-1 border-b border-[#262626] px-2 py-2">
			<div class="min-w-0 flex-1">
				<LangPicker bind:value={from} detect align="start" />
			</div>
			<button
				type="button"
				class="flex size-10 shrink-0 cursor-pointer items-center justify-center text-[#5c5c5c] transition-colors duration-150 hover:text-[#c0c0c0] disabled:cursor-default disabled:opacity-40 disabled:hover:text-[#5c5c5c]"
				title="swap languages"
				aria-label="swap languages"
				disabled={!sourceLang}
				onclick={onSwap}
			>
				<ArrowLeftRight size={16} class="transition-transform duration-300 ease-out {swapping ? 'rotate-180' : ''}" />
			</button>
			<div class="min-w-0 flex-1">
				<LangPicker bind:value={to} align="end" />
			</div>
		</div>

		<div class="grid grid-cols-2 max-[560px]:grid-cols-1">
			<div class="relative flex min-h-[150px] flex-col rounded-bl-md max-[560px]:min-h-[112px] max-[560px]:border-b max-[560px]:border-[#1c1c1c] border-r border-[#1c1c1c]">
				<textarea
					class="block w-full flex-1 resize-none bg-transparent px-4 pt-3.5 pr-[52px] pb-2 text-[15px] leading-relaxed text-[#d0d0d0] outline-none placeholder:text-[#3d3d3d] {text.length <= 80
						? 'text-[18px]'
						: ''}"
					dir={sourceRtl ? 'rtl' : 'auto'}
					rows="4"
					maxlength={MAX}
					placeholder="enter text…"
					bind:value={text}
				></textarea>
				{#if text}
					<button
						type="button"
						class="absolute top-1.5 right-1.5 cursor-pointer rounded p-1.5 text-[#3d3d3d] transition-colors hover:bg-[#1e1e1e] hover:text-[#8a8a8a]"
						title="clear"
						aria-label="clear text"
						onclick={onClear}
					>
						<X size={14} />
					</button>
				{/if}
				<div class="mt-auto flex min-h-[46px] items-center gap-1 px-2 pb-1.5">
					{#if text}
						<button
							type="button"
							class="flex cursor-pointer items-center gap-1.5 rounded px-1.5 py-1 text-[#5c5c5c] transition-colors hover:bg-[#1e1e1e] hover:text-[#c0c0c0]"
							title="speak source text"
							onclick={() => {
								const u = new SpeechSynthesisUtterance(text.slice(0, 400));
								u.lang = sourceLang ?? to;
								speechSynthesis.cancel();
								speechSynthesis.speak(u);
							}}
						>
							<Volume2 size={14} />
						</button>
					{/if}
					<span class="ml-auto pr-2 text-[12px] tabular-nums text-[#5c5c5c]">
						{text.length.toLocaleString('en-US')} / {MAX.toLocaleString('en-US')}
					</span>
				</div>
			</div>

			<div class="flex min-h-[150px] flex-col rounded-br-md bg-[#151515] max-[560px]:min-h-[112px] max-[560px]:rounded-br-none max-[560px]:rounded-bl-md">
				<div
					class="max-h-[280px] flex-1 overflow-y-auto px-4 pt-3.5 pb-2 text-[15px] leading-relaxed break-words whitespace-pre-wrap text-[#d0d0d0] {text.length <=
					80 ? 'text-[18px]'
					: ''}"
					dir={targetRtl ? 'rtl' : 'auto'}
					aria-live="polite"
				>
					{#if text}
						<span class="text-[#4a4a4a] italic">translation</span>
					{:else}
						<span class="text-[#3d3d3d]">{langByCode(to)?.name ?? 'translation'}</span>
					{/if}
				</div>
				<div class="mt-auto flex min-h-[46px] items-center gap-1 px-2 pb-1.5">
					<span class="flex items-center gap-1.5 px-1.5 py-1 text-[12px] text-[#3d3d3d]">
						<CopyButton active={false} size={14} />
					</span>
					<span class="mr-auto pl-1 text-[12px] text-[#5c5c5c]"></span>
				</div>
			</div>
		</div>
	</div>
{/if}
