import { market, myIP } from '../../api';
import type { MarketQuote } from '../../types';

export type QuickAnswer = {
	type: 'calculator' | 'conversion' | 'ip' | 'time' | 'color' | 'market';
	title: string;
	value: string;
	detail?: string;
	fromCurrency?: string;
	toCurrency?: string;
	fromUnit?: string;
	toUnit?: string;
	amount?: string;
	market?: MarketQuote;
};

const currencyNames: Record<string, boolean> = {
	USD: true, EUR: true, GBP: true, RUB: true, UAH: true, JPY: true, CNY: true,
	CHF: true, CAD: true, AUD: true, PLN: true, TRY: true
};

function formatNumber(value: number): string {
	return new Intl.NumberFormat(undefined, { maximumFractionDigits: 6 }).format(value);
}

function calculate(expression: string): number | null {
	const normalized = expression.replace(/,/g, '.').replace(/\s+/g, '');
	const tokens = normalized.match(/\d+(?:\.\d+)?|[()+\-*/%]/g);
	if (!tokens || tokens.join('') !== normalized) return null;
	const values: number[] = [];
	const ops: string[] = [];
	const precedence: Record<string, number> = { '+': 1, '-': 1, '*': 2, '/': 2, '%': 2 };
	const apply = () => {
		const op = ops.pop(); const b = values.pop(); const a = values.pop();
		if (!op || a === undefined || b === undefined || (op === '/' && b === 0)) throw new Error();
		values.push(op === '+' ? a + b : op === '-' ? a - b : op === '*' ? a * b : op === '/' ? a / b : a % b);
	};
	try {
		for (const token of tokens) {
			if (/^\d/.test(token)) values.push(Number(token));
			else if (token === '(') ops.push(token);
			else if (token === ')') {
				while (ops.at(-1) !== '(') apply();
				if (ops.pop() !== '(') throw new Error();
			} else {
				while (ops.length && ops.at(-1) !== '(' && precedence[ops.at(-1) ?? ''] >= precedence[token]) apply();
				ops.push(token);
			}
		}
		while (ops.length) { if (ops.at(-1) === '(') throw new Error(); apply(); }
		return values.length === 1 && Number.isFinite(values[0]) ? values[0] : null;
	} catch { return null; }
}

function unitAnswer(query: string): QuickAnswer | null {
	const match = query.trim().match(/^(-?\d+(?:[.,]\d+)?)\s*(millimeters?|mm|centimeters?|cm|meters?|m|kilometers?|km|inch(?:es)?|in|feet|ft|yards?|yd|miles?|mi|milligrams?|mg|grams?|g|kilograms?|kg|tons?|t|ounces?|oz|pounds?|lb|milliliters?|ml|centiliters?|cl|deciliters?|dl|liters?|l|gallons?|gal|m\/s|km\/h|mph|c|°c|f|°f|k)\s+(?:in|to|в)\s*(millimeters?|mm|centimeters?|cm|meters?|m|kilometers?|km|inch(?:es)?|in|feet|ft|yards?|yd|miles?|mi|milligrams?|mg|grams?|g|kilograms?|kg|tons?|t|ounces?|oz|pounds?|lb|milliliters?|ml|centiliters?|cl|deciliters?|dl|liters?|l|gallons?|gal|m\/s|km\/h|mph|c|°c|f|°f|k)$/i);
	if (!match) return null;
	const amount = Number(match[1].replace(',', '.'));
	const normalizeUnit = (unit: string) => {
		const value = unit.toLowerCase().replace('°', '');
		const aliases: Record<string, string> = {
			millimeter: 'mm', millimeters: 'mm', centimeter: 'cm', centimeters: 'cm', meter: 'm', meters: 'm',
			kilometer: 'km', kilometers: 'km', inch: 'in', inches: 'in', foot: 'ft', feet: 'ft', yard: 'yd', yards: 'yd',
			mile: 'mi', miles: 'mi', milligram: 'mg', milligrams: 'mg', gram: 'g', grams: 'g', kilogram: 'kg', kilograms: 'kg',
			tons: 't', ton: 't', ounce: 'oz', ounces: 'oz', pound: 'lb', pounds: 'lb', milliliter: 'ml', milliliters: 'ml',
			centiliter: 'cl', centiliters: 'cl', deciliter: 'dl', deciliters: 'dl', liter: 'l', liters: 'l', gallon: 'gal', gallons: 'gal'
		};
		return aliases[value] ?? value;
	};
	const from = normalizeUnit(match[2]); const to = normalizeUnit(match[3]);
	let result: number;
	if (from === 'c' && to === 'f') result = amount * 9 / 5 + 32;
	else if (from === 'f' && to === 'c') result = (amount - 32) * 5 / 9;
	else if (from === 'c' && to === 'k') result = amount + 273.15;
	else if (from === 'k' && to === 'c') result = amount - 273.15;
	else if (from === 'f' && to === 'k') result = (amount - 32) * 5 / 9 + 273.15;
	else if (from === 'k' && to === 'f') result = (amount - 273.15) * 9 / 5 + 32;
	else {
		const factors: Record<string, number> = {
			mm: 0.001, cm: 0.01, m: 1, km: 1000, in: 0.0254, ft: 0.3048, yd: 0.9144, mi: 1609.344,
			mg: 0.000001, g: 0.001, kg: 1, t: 1000, oz: 0.0283495231, lb: 0.45359237,
			ml: 0.001, cl: 0.01, dl: 0.1, l: 1, gal: 3.785411784,
			'm/s': 1, 'km/h': 0.277777778, mph: 0.44704
		};
		if (!(from in factors) || !(to in factors)) return null;
		result = amount * factors[from] / factors[to];
	}
	return { type: 'conversion', title: 'conversion', value: `${formatNumber(result)} ${to}`, fromUnit: from, toUnit: to, amount: formatNumber(amount) };
}

export async function getQuickAnswer(query: string): Promise<QuickAnswer | null> {
	const q = query.trim();
	if (/^(my ip|what is my ip)$/i.test(q)) return { type: 'ip', title: 'your ip', value: await myIP() };
	const expression = q.replace(/^what is\s+/i, '').trim();
	if (/^[\d\s()+\-*/%,.]+$/.test(expression) && /[+*/%\-]/.test(expression)) {
		const result = calculate(expression.replace(/%/g, '/100'));
		if (result !== null) return { type: 'calculator', title: 'calculator', value: formatNumber(result), detail: q };
	}
	const units = unitAnswer(q); if (units) return units;
	const color = q.match(/^(?:color|colour)\s+(#[0-9a-f]{3,8})$/i);
	if (color) return { type: 'color', title: 'color', value: color[1].toUpperCase(), detail: color[1] };
	const time = q.match(/^(?:time)\s+(?:in|в)\s+([a-z_ /-]+)$/i);
	if (time) {
		try {
			const zone = time[1].trim().replace(/\s+/g, '_');
			return { type: 'time', title: `time in ${time[1].trim()}`, value: new Intl.DateTimeFormat(undefined, { timeZone: zone, dateStyle: 'medium', timeStyle: 'short' }).format(new Date()) };
		} catch { /* not a valid timezone */ }
	}
	const currency = q.match(/^(-?\d+(?:[.,]\d+)?)\s*([A-Za-z]{3})\s+(?:to|in|в)\s*([A-Za-z]{3})$/i);
	if (currency) {
		const amount = Number(currency[1].replace(',', '.')); const from = currency[2].toUpperCase(); const to = currency[3].toUpperCase();
		if (currencyNames[from] && currencyNames[to]) try {
			const response = await fetch(`https://api.frankfurter.dev/v2/rate/${from}/${to}`); if (!response.ok) return null;
			const data = (await response.json()) as { rate: number; date?: string };
			return {
				type: 'conversion', title: 'currency conversion', value: `${formatNumber(amount * data.rate)} ${to}`,
				detail: `rate ${formatNumber(data.rate)}${data.date ? ` · ${data.date}` : ''}`,
				fromCurrency: from, toCurrency: to, amount: formatNumber(amount)
			};
		} catch { return null; }
	}
	try {
		const quote = await market(q);
		if (quote) return { type: 'market', title: 'market', value: String(quote.current), market: quote };
	} catch { /* market data is optional */ }
	return null;
}
