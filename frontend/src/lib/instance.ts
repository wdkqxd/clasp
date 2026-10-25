export const OFFICIAL_HOST = 'clasp.wdkq.cc';

export function isOfficialHost(hostname: string): boolean {
	return hostname === OFFICIAL_HOST;
}
