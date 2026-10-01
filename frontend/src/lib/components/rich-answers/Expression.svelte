<script lang="ts">
	import { all, create } from 'mathjs';
	import CopyButton from './CopyButton.svelte';

	const math = create(all, {});

	let { query = '' }: { query?: string } = $props();

	const CHARS_RE = /^[\d\s+\-*/^%().,!a-z=]+$/i;
	const OPERATOR_RE =
		/[+\-*/^%!=]|\b(?:sqrt|abs|sign|round|floor|ceil|trunc|fix|exp|log2?|log10|hypot|sin|cos|tan|asin|acos|atan|atan2|sinh|cosh|tanh|min|max|pow|mod|nthRoot|factorial|gamma|unaryfactorial)\b/i;
	const ISO_DATE_RE = /^\d{4}-\d{1,2}-\d{1,2}(?:[ t]\d{1,2}:\d{2}(?::\d{2})?)?$/;
	const PHONE_RE = /^\d{1,4}(?:[-\s]\d{2,4}){2,}$/;

	function toNumber(source: string): number | null {
		let result;
		try {
			result = math.evaluate(source);
		} catch {
			return null;
		}
		let value: unknown = result;
		if (math.isMatrix(result)) {
			const size = result.size();
			if (size.length !== 1 || size[0] !== 1) return null;
			value = result.get([0]);
		}
		if (math.isUnit(value)) return null;
		const num = Number(value);
		return Number.isFinite(num) ? num : null;
	}

	function evaluate(input: string): { expression: string; value: number | boolean } | null {
		const cleaned = input.trim().replace(/\s+/g, ' ');
		if (!cleaned || cleaned.length > 80) return null;
		if (!/\d/.test(cleaned)) return null;
		if (ISO_DATE_RE.test(cleaned) || PHONE_RE.test(cleaned)) return null;
		if (!CHARS_RE.test(cleaned)) return null;
		if (!OPERATOR_RE.test(cleaned)) return null;

		const raw = cleaned.replace(/(\d),(\d)/g, '$1.$2');

		const equality = raw.match(/^(.+?)\s*=\s*(.+)$/);
		if (equality) {
			const left = toNumber(equality[1]);
			const right = toNumber(equality[2]);
			if (left === null || right === null) return null;
			const scale = Math.max(1, Math.abs(left), Math.abs(right));
			return { expression: raw, value: Math.abs(left - right) < scale * 1e-9 };
		}

		const ofMatch = raw.match(/^(\d+(?:\.\d+)?)\s*(?:%|percent)\s+of\s+([\d.]+)$/i);
		const source = ofMatch ? `(${ofMatch[1]} / 100) * ${ofMatch[2]}` : raw;

		const value = toNumber(source);
		if (value === null) return null;
		return { expression: raw, value };
	}

	const solved = $derived(evaluate(query));

	let copied = $state(false);
	function onCopy() {
		if (!solved) return;
		navigator.clipboard?.writeText(format(solved.value));
		copied = true;
		setTimeout(() => (copied = false), 1200);
	}

	function format(value: number | boolean) {
		if (typeof value === 'boolean') return value ? 'true' : 'false';
		if (Number.isInteger(value) && Math.abs(value) < 1e15) return String(value);
		const rounded = Number(value.toPrecision(12));
		if (Math.abs(rounded) >= 1e15 || (Math.abs(rounded) < 1e-6 && rounded !== 0)) {
			return rounded.toExponential(6).replace(/e([+-])(\d)$/, 'e$10$2');
		}
		return String(rounded);
	}
</script>

{#if solved}
	<div class="rounded-md border border-[#292929] bg-[#121212] px-5 py-4 shadow-2xl shadow-black/30">
		<div class="font-mono text-[13px] break-words text-[#707070]">{solved.expression}</div>
		<button
			type="button"
			class="group mt-1 inline-flex cursor-pointer items-center gap-3 font-mono text-[32px] leading-none text-[#e0e0e0]"
			title="copy result"
			onclick={onCopy}
		>
			{format(solved.value)}
			<CopyButton active={copied} />
		</button>
	</div>
{/if}
