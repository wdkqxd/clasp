export function pillGroupClass(optionCount: number): string {
	const pressScale = optionCount > 2 ? 'active:scale-[0.985]' : 'active:scale-[0.97]';
	return `relative inline-flex max-w-full items-center rounded-[5px] bg-[#141414] p-[3px] transition-[scale] duration-150 ease-out ${pressScale} active:duration-75`;
}
