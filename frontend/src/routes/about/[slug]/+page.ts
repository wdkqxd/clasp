import { error } from '@sveltejs/kit';
import { aboutDoc, aboutDocs } from '$lib/content/about';
import type { PageLoad } from './$types';

export const load: PageLoad = ({ params }) => {
	const doc = aboutDoc(params.slug);
	if (!doc) error(404, 'not found');
	return { doc, docs: aboutDocs() };
};
