package engine

import (
	"reflect"
	"testing"
)

func TestParseBraveDiscussions(t *testing.T) {
	doc := `<div class="snippet svelte-jmfu5f" id="discussions"><header>Discussions</header>` +
		`<hr class="snippet-divider"/><div class="discussions-item svelte-1gf5iup"><div class="accordion-item"><details><summary title="What is linux?" class="svelte-1rwoi5h"><div class="title svelte-1rwoi5h">What is linux?</div></summary>` +
		`<div class="content desktop-default-regular t-secondary svelte-1gf5iup">Just like windows and macos, just different cuz it's a different operating system <a class="source svelte-1gf5iup" href="https://www.reddit.com/r/linux4noobs/comments/1ljteoi/what_is_linux/">More on reddit.com</a></div></details></div></div>` +
		`<div class="discussions-item svelte-1gf5iup"><div class="accordion-item"><details><summary title="Linux+" class="svelte-1rwoi5h"><div class="title svelte-1rwoi5h">Linux+</div></summary>` +
		`<div class="content t-secondary svelte-1gf5iup">Some compTIA talk <a class="source svelte-1gf5iup" href="https://www.reddit.com/r/CompTIA/comments/15bjymi/linux/">More on reddit.com</a></div></details></div></div>` +
		`<div class="discussions-item svelte-1gf5iup"><div class="accordion-item"><details><summary title="Linux+" class="svelte-1rwoi5h"><div class="title svelte-1rwoi5h">Linux+</div></summary>` +
		`<div class="content t-secondary svelte-1gf5iup">Same talk <a class="source svelte-1gf5iup" href="https://www.reddit.com/r/CompTIA/comments/15bjymi/linux/">More on reddit.com</a></div></details></div></div>` +
		`</div>`

	got := parseBraveDiscussions(doc)
	want := []Discussion{
		{Title: "What is linux?", URL: "https://www.reddit.com/r/linux4noobs/comments/1ljteoi/what_is_linux/", Snippet: "Just like windows and macos, just different cuz it's a different operating system"},
		{Title: "Linux+", URL: "https://www.reddit.com/r/CompTIA/comments/15bjymi/linux/", Snippet: "Some compTIA talk"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("parseBraveDiscussions = %#v, want %#v", got, want)
	}
}

func TestParseBraveDiscussionsEmpty(t *testing.T) {
	if got := parseBraveDiscussions(`<html><body>no discussions</body></html>`); len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}