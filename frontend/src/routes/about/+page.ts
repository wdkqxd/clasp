import { redirect } from '@sveltejs/kit';
import { aboutDocs } from '$lib/content/about';
import type { PageLoad } from './$types';

export const load: PageLoad = () => {
	const first = aboutDocs()[0];
	if (first) redirect(307, `/about/${first.slug}`);
	return {};
};