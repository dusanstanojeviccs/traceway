<script lang="ts">
	let {
		scope,
		sorted = true
	}: {
		scope: Record<string, string>;
		sorted?: boolean;
	} = $props();

	const entries = $derived(() => {
		const items = Object.entries(scope);
		if (sorted) {
			return items.sort((a, b) => a[0].localeCompare(b[0]));
		}
		return items;
	});
</script>

<div class="grid grid-cols-1 gap-3 sm:grid-cols-2 md:grid-cols-3">
	{#each entries() as [key, value]}
		<div class="flex flex-col gap-1 rounded-md bg-muted p-3">
			<span class="text-xs font-medium text-muted-foreground">{key}</span>
			<span class="font-mono text-sm break-all">{value}</span>
		</div>
	{/each}
</div>
