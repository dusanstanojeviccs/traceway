<script lang="ts">
	import { TriangleAlert } from 'lucide-svelte';

	type ImpactLevel = 'critical' | 'high' | 'medium' | null;

	let {
		level,
		showIcon
	}: {
		level: ImpactLevel;
		showIcon?: boolean;
	} = $props();

	// Default showIcon to true for critical and high levels
	const shouldShowIcon = $derived(showIcon ?? (level === 'critical' || level === 'high'));

	const config = $derived(() => {
		switch (level) {
			case 'critical':
				return {
					bg: 'bg-red-500/15',
					text: 'text-red-600 dark:text-red-400',
					label: 'Critical'
				};
			case 'high':
				return {
					bg: 'bg-orange-500/15',
					text: 'text-orange-600 dark:text-orange-400',
					label: 'High'
				};
			case 'medium':
				return {
					bg: 'bg-yellow-500/15',
					text: 'text-yellow-600 dark:text-yellow-500',
					label: 'Medium'
				};
			default:
				return null;
		}
	});
</script>

{#if config()}
	<span
		class="inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium {config()
			?.bg} {config()?.text}"
	>
		{#if shouldShowIcon}
			<TriangleAlert class="h-3 w-3" />
		{/if}
		{config()?.label}
	</span>
{/if}
