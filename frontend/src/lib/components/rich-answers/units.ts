const conversions: Record<string, Record<string, number>> = {
	mass: {
		gram: 1,
		kilogram: 1000,
		pound: 453.592,
		poundmass: 453.592,
		ounce: 28.3495,
		milligram: 0.001,
		ton: 907185,
		tonne: 1000000,
		stone: 6350.29,
	},
	length: {
		meter: 1,
		metre: 1,
		kilometer: 1000,
		kilometre: 1000,
		mile: 1609.34,
		foot: 0.3048,
		feet: 0.3048,
		inch: 0.0254,
		centimeter: 0.01,
		centimetre: 0.01,
		yard: 0.9144,
		millimeter: 0.001,
		millimetre: 0.001,
	},
	volume: {
		liter: 1,
		litre: 1,
		milliliter: 0.001,
		millilitre: 0.001,
		gallon: 3.78541,
		quart: 0.946353,
		pint: 0.473176,
		cup: 0.24,
		fluidounce: 0.0295735,
		tablespoon: 0.0147868,
		teaspoon: 0.00492892,
	},
	area: {
		squaremeter: 1,
		squarefoot: 0.092903,
		squareyard: 0.836127,
		acre: 4046.86,
		hectare: 10000,
		squarekilometer: 1000000,
		squaremile: 2590000,
	},
	time: {
		second: 1,
		minute: 60,
		hour: 3600,
		day: 86400,
		week: 604800,
		month: 2628000,
		year: 31536000,
	},
	speed: {
		meterpersecond: 1,
		kilometerperhour: 0.277778,
		mileperhour: 0.44704,
		knot: 0.514444,
	},
	digital: {
		byte: 1,
		kilobyte: 1024,
		megabyte: 1048576,
		gigabyte: 1073741824,
		terabyte: 1099511627776,
		bit: 0.125,
	},
};

export function convertUnit(amount: number, from: string, to: string, dim: string): string {
	if (dim === 'temperature') {
		const fromLower = from.toLowerCase();
		const toLower = to.toLowerCase();
		if (fromLower === 'celsius' && toLower === 'fahrenheit') return ((amount * 9) / 5 + 32).toFixed(2);
		if (fromLower === 'fahrenheit' && toLower === 'celsius') return (((amount - 32) * 5) / 9).toFixed(2);
		if (fromLower === 'celsius' && toLower === 'kelvin') return (amount + 273.15).toFixed(2);
		if (fromLower === 'kelvin' && toLower === 'celsius') return (amount - 273.15).toFixed(2);
		if (fromLower === 'fahrenheit' && toLower === 'kelvin') return (((amount - 32) * 5) / 9 + 273.15).toFixed(2);
		if (fromLower === 'kelvin' && toLower === 'fahrenheit') return (((amount - 273.15) * 9) / 5 + 32).toFixed(2);
		return String(amount);
	}

	const table = conversions[dim];
	if (!table) return String(amount);

	const fromKey = from.toLowerCase().replace(/[^a-z]/g, '');
	const toKey = to.toLowerCase().replace(/[^a-z]/g, '');
	const fromFactor = table[fromKey] || 1;
	const toFactor = table[toKey] || 1;
	const result = (amount * fromFactor) / toFactor;

	if (Math.abs(result) >= 1000) return result.toLocaleString(undefined, { maximumFractionDigits: 2 });
	if (Math.abs(result) >= 1) return result.toFixed(4).replace(/\.?0+$/, '');
	return result.toPrecision(4);
}

export function formatUnit(unit: string): string {
	const mapped: Record<string, string> = {
		poundmass: 'pound',
		squaremeter: 'square meter',
		squarefoot: 'square foot',
		squareyard: 'square yard',
		squarekilometer: 'square kilometer',
		squaremile: 'square mile',
		meterpersecond: 'meter/second',
		kilometerperhour: 'km/h',
		mileperhour: 'mph',
		fluidounce: 'fluid ounce',
	};
	const lower = unit.toLowerCase().replace(/[^a-z]/g, '');
	return (
		mapped[lower] ??
		unit
			.replace(/([a-z])([A-Z])/g, '$1 $2')
			.replace(/_/g, ' ')
			.toLowerCase()
	);
}