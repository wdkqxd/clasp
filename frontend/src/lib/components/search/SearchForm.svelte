<script lang="ts">
	let {
		value = $bindable(''),
		placeholder = 'Search...',
		onSearch
	}: {
		value?: string;
		placeholder?: string;
		onSearch: (q: string) => void;
	} = $props();
</script>

<form
	onsubmit={(e) => {
		e.preventDefault();
		const q = value.trim();
		if (q) onSearch(q);
	}}
>
	<div class="relative">
		<input
			type="search"
			bind:value
			autocomplete="off"
			spellcheck="false"
			placeholder={placeholder}
			class="w-full bg-transparent border-0 border-b border-[#363636] rounded-none py-2.5 px-0 text-[15px] text-[#d8d8d8] focus:outline-none focus:border-[#8a8a8a] placeholder:text-[#4a4a4a] transition-colors"
		/>
	</div>
</form>
