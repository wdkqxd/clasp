import { marked } from 'marked';

const sources = import.meta.glob('./about/*.md', {
	query: '?raw',
	import: 'default',
	eager: true
}) as Record<string, string>;

export type AboutDoc = {
	slug: string;
	title: string;
	html: string;
};

function parseDoc(path: string, source: string): AboutDoc {
	const slug = (path.split('/').pop() ?? '').replace(/\.md$/, '');
	const fm = /^---\r?\n([\s\S]*?)\r?\n---\r?\n?/.exec(source);

	let title = slug;
	let body = source;

	if (fm) {
		body = source.slice(fm[0].length);
		const t = /^title:\s*(.+)$/m.exec(fm[1]);
		if (t) title = t[1].trim().replace(/^["']|["']$/g, '');
	}

	return { slug, title, html: marked.parse(body, { async: false }) as string };
}

const docs: AboutDoc[] = Object.entries(sources)
	.map(([path, source]) => parseDoc(path, source))
	.sort((a, b) => a.slug.localeCompare(b.slug));

export function aboutDocs(): AboutDoc[] {
	return docs;
}

export function aboutDoc(slug: string): AboutDoc | undefined {
	return docs.find((d) => d.slug === slug);
}
