export function sym(code: string): string {
	const map: Record<string, string> = { usd: '$', eur: '€', gbp: '£', rub: '₽', jpy: '¥', cny: '¥' };
	return map[code.toLowerCase()] ?? '';
}

export function price(n: number): string {
	return n.toLocaleString(undefined, { maximumFractionDigits: 2 });
}

export function price2(n: number): string {
	return n.toLocaleString(undefined, { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}