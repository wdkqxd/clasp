import { marked, type Tokens } from 'marked';
import { headingAnchor, slugifyHeading } from '$lib/content/headings';

const headingSlugs = new Map<string, number>();

marked.use({
	renderer: {
		heading({ tokens, depth }: Tokens.Heading) {
			const inner = this.parser.parseInline(tokens);
			const base = slugifyHeading(inner.replace(/<[^>]*>/g, ''));
			const count = (headingSlugs.get(base) ?? 0) + 1;
			headingSlugs.set(base, count);
			const id = count === 1 ? base : `${base}-${count}`;

			return `<h${depth} id="${id}">${inner}${headingAnchor(id)}</h${depth}>\n`;
		},
	},
});

const sources = import.meta.glob('./about/*.md', {
	query: '?raw',
	import: 'default',
	eager: true,
}) as Record<string, string>;

export type AboutDoc = {
	slug: string;
	title: string;
	html: string;
};

function groupIntoSections(html: string): string {
	return html
		.split(/(?=<h2\s)/)
		.map((part) =>
			part.startsWith('<h2 ') ? `<section class="doc-section">${part}</section>` : part
		)
		.join('');
}

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

	headingSlugs.clear();
	const html = marked.parse(body, { async: false }) as string;
	return { slug, title, html: groupIntoSections(html) };
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
