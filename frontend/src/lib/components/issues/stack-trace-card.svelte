<script lang="ts">
    import * as Card from "$lib/components/ui/card";
    import { formatDateTime } from '$lib/utils/formatters';
    import { getTimezone } from '$lib/state/timezone.svelte';

    interface Props {
        stackTrace: string;
        isMessage?: boolean;
        firstSeen?: string;
        lastSeen?: string;
        totalCount?: number;
        timezone?: string;
    }

    let { stackTrace, isMessage = false, firstSeen, lastSeen, totalCount, timezone }: Props = $props();

    const tz = $derived(timezone ?? getTimezone());
    const showStats = $derived(firstSeen && lastSeen && totalCount !== undefined);
</script>

<Card.Root>
    <Card.Header>
        <div class="flex items-center gap-2">
            <Card.Title>Stack Trace</Card.Title>
            {#if isMessage}
                <span class="inline-flex items-center rounded-md bg-blue-50 dark:bg-blue-900/30 px-2 py-1 text-xs font-medium text-blue-700 dark:text-blue-300 ring-1 ring-inset ring-blue-700/10 dark:ring-blue-400/30">
                    Message
                </span>
            {/if}
        </div>
        {#if showStats}
            <Card.Description>
                First seen: {formatDateTime(firstSeen!, { timezone: tz })} ·
                Last seen: {formatDateTime(lastSeen!, { timezone: tz })} ·
                Total occurrences: {totalCount}
            </Card.Description>
        {/if}
    </Card.Header>
    <Card.Content>
        <div class="rounded-md bg-muted p-4 overflow-x-auto">
            <pre class="text-sm whitespace-pre-wrap font-mono">{stackTrace}</pre>
        </div>
    </Card.Content>
</Card.Root>
