<script lang="ts">
	import type { RichColorpicker } from '../../types';
	import CopyButton from './CopyButton.svelte';

	let { colorpicker }: { colorpicker: RichColorpicker } = $props();

	let pickColor = $state(colorpicker.color);
	let copied = $state('');
	function onCopy(label: string, value: string) {
		navigator.clipboard?.writeText(value);
		copied = label;
		setTimeout(() => {
			if (copied === label) copied = '';
		}, 1200);
	}

	function parseColor(hex: string) {
		const m = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex);
		return m ? { r: parseInt(m[1], 16), g: parseInt(m[2], 16), b: parseInt(m[3], 16) } : null;
	}
	function hslOf(rgb: { r: number; g: number; b: number }): string {
		let r = rgb.r / 255;
		let g = rgb.g / 255;
		let b = rgb.b / 255;
		const max = Math.max(r, g, b);
		const min = Math.min(r, g, b);
		let h = 0;
		let s = 0;
		const l = (max + min) / 2;
		if (max !== min) {
			const d = max - min;
			s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
			switch (max) {
				case r: h = (g - b) / d + (g < b ? 6 : 0); break;
				case g: h = (b - r) / d + 2; break;
				default: h = (r - g) / d + 4; break;
			}
			h /= 6;
		}
		return `${Math.round(h * 360)}, ${Math.round(s * 100)}%, ${Math.round(l * 100)}%`;
	}
	function colorFormats(hex: string) {
		const rgb = parseColor(hex);
		if (!rgb) return [];
		return [
			{ label: 'HEX', value: hex.toUpperCase() },
			{ label: 'RGB', value: `rgb(${rgb.r}, ${rgb.g}, ${rgb.b})` },
			{ label: 'HSL', value: `hsl(${hslOf(rgb)})` }
		];
	}
</script>

<div class="flex items-center gap-4 rounded-md border border-[#292929] bg-[#121212] p-4 shadow-2xl shadow-black/30">
	<label
		class="block h-10 w-10 shrink-0 cursor-pointer rounded-md border border-[#292929]"
		style="background:{pickColor}"
		title="pick color"
	>
		<input
			type="color"
			value={pickColor}
			oninput={(e) => (pickColor = e.currentTarget.value)}
			class="absolute inset-0 h-0 w-0 opacity-0"
			aria-label="pick color"
		/>
	</label>
	<div class="min-w-0">
		{#each colorFormats(pickColor) as f (f.label)}
			<div class="flex items-baseline gap-3">
				<span class="w-8 shrink-0 text-[11px] tracking-wide text-[#5c5c5c]">{f.label}</span>
				<button
					type="button"
					class="group inline-flex cursor-pointer items-center gap-2 text-[13px] text-[#c0c0c0]"
					title="copy {f.label}"
					onclick={() => onCopy(f.label, f.value)}
				>
					{f.value}
					<CopyButton active={copied === f.label} size="h-3.5 w-3.5" />
				</button>
			</div>
		{/each}
	</div>
</div>