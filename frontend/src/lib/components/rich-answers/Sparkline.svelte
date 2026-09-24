<script lang="ts">
	let {
		points,
		width = 260,
		height = 56
	}: {
		points: [number, number][];
		width?: number;
		height?: number;
	} = $props();

	const up = $derived(points.length > 1 && points[points.length - 1][1] >= points[0][1]);
	const lineColor = $derived(up ? '#a6e3a1' : '#f38ba8');
	const fillColor = $derived(up ? 'rgba(166, 227, 161, 0.15)' : 'rgba(243, 139, 168, 0.15)');

	const path = $derived.by(() => {
		if (points.length < 2) return '';
		const vs = points.map((p) => p[1]);
		const ts = points.map((p) => p[0]);
		let minV = Math.min(...vs);
		let maxV = Math.max(...vs);
		if (maxV === minV) {
			maxV += 1;
			minV -= 1;
		}
		const minT = ts[0];
		const maxT = ts[ts.length - 1] ?? ts[0];
		const spanT = maxT - minT || 1;
		const x = (i: number) => (i / (points.length - 1)) * width;
		const y = (v: number) => height - ((v - minV) / (maxV - minV)) * height;
		let d = `M ${x(0)} ${y(vs[0])}`;
		for (let i = 1; i < points.length; i++) {
			d += ` L ${x(i)} ${y(vs[i])}`;
		}
		return d;
	});

	let fill = $derived(`${path} L ${width} ${height} L 0 ${height} Z`);
</script>

{#if path}
	<svg
		width={width}
		height={height}
		viewBox={`0 0 ${width} ${height}`}
		class="block w-full"
		preserveAspectRatio="none"
		aria-hidden="true"
	>
		<path d={fill} fill={fillColor} stroke="none" />
		<path d={path} fill="none" stroke={lineColor} stroke-width="1.5" vector-effect="non-scaling-stroke" />
	</svg>
{/if}