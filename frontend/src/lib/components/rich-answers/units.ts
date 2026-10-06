import { create, UnitDependencies } from 'mathjs';

const math = create({ UnitDependencies });

// spellings mathjs does not know on its own
const UNIT_ALIASES: Record<string, string> = {
	d: 'day',
	days: 'day',
	wk: 'week',
	wks: 'week',
	mph: 'mi/h',
	kph: 'km/h',
	c: 'degC',
	celsius: 'degC',
	centigrade: 'degC',
	'°c': 'degC',
	f: 'degF',
	fahrenheit: 'degF',
	'°f': 'degF',
	k: 'K',
	kelvin: 'K',
};

const UNIT_LABELS: Record<string, string> = {
	degC: 'celsius',
	degF: 'fahrenheit',
	K: 'kelvin',
	'mi/h': 'mph',
	'km/h': 'km/h',
	day: 'day',
	week: 'week',
};

function unitName(unit: string): string | null {
	const raw = unit.trim();
	const alias = UNIT_ALIASES[raw.toLowerCase()] ?? raw;
	try {
		new math.Unit(1, alias);
		return alias;
	} catch {
		return null;
	}
}

const CONVERSION_RE = /^(-?[\d.,]+)\s*(.+?)\s+(?:to|in|into|as)\s+(.+?)$/i;

export function parseConversion(query: string): {
	amount: number;
	from_unit: string;
	to_unit: string;
	from_label: string;
	to_label: string;
} | null {
	const match = query.trim().match(CONVERSION_RE);
	if (!match) return null;

	const amount = Number(match[1].replace(/,(?=\d{3}\b)/g, '').replace(',', '.'));
	if (!Number.isFinite(amount)) return null;

	const from = unitName(match[2]);
	const to = unitName(match[3]);
	if (!from || !to) return null;

	const fromUnit = new math.Unit(amount, from);
	if (!fromUnit.equalBase(new math.Unit(1, to))) return null;

	return {
		amount,
		from_unit: from,
		to_unit: to,
		from_label: UNIT_LABELS[from] ?? match[2].trim(),
		to_label: UNIT_LABELS[to] ?? match[3].trim(),
	};
}

export function convertUnit(amount: number, from: string, to: string): string {
	const result = new math.Unit(amount, from).toNumber(to);
	if (!Number.isFinite(result)) return String(amount);
	if (Math.abs(result) >= 1000)
		return result.toLocaleString(undefined, { maximumFractionDigits: 2 });
	if (Math.abs(result) >= 1) return result.toFixed(4).replace(/\.?0+$/, '');
	return String(Number.parseFloat(result.toPrecision(4)));
}
