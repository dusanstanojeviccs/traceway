<script lang="ts">
    import { onMount } from 'svelte';
    import { page } from '$app/state';
    import { api } from '$lib/api';
    import { LoadingCircle } from "$lib/components/ui/loading-circle";
    import { ErrorDisplay } from "$lib/components/ui/error-display";
    import { projectsState } from '$lib/state/projects.svelte';
    import { StackTraceCard, EventCard, EventsTable, PageHeader } from '$lib/components/issues';
    import type { ExceptionGroup, ExceptionOccurrence, LinkedTransaction } from '$lib/types/exceptions';

    let group = $state<ExceptionGroup | null>(null);
    let occurrences = $state<ExceptionOccurrence[]>([]);
    let loading = $state(true);
    let error = $state('');
    let notFound = $state(false);
    let total = $state(0);
    let linkedTransaction = $state<LinkedTransaction | null>(null);

    const exceptionHash = $derived(page.params.exceptionHash ?? '');
    const latestOccurrence = $derived(occurrences[0]);
    const isMessage = $derived(latestOccurrence?.isMessage ?? false);
    const hasMoreOccurrences = $derived(total > 10);
    const firstLineOfStackTrace = $derived(group?.stackTrace.split('\n')[0] || 'Exception');

    async function loadData() {
        loading = true;
        error = '';
        notFound = false;
        linkedTransaction = null;

        try {
            const exceptionHash = page.params.exceptionHash;
            const response = await api.post(`/exception-stack-traces/${exceptionHash}`, {
                pagination: {
                    page: 1,
                    pageSize: 10
                }
            }, { projectId: projectsState.currentProjectId ?? undefined });

            group = response.group;
            occurrences = response.occurrences || [];
            total = response.pagination.total;

            // Load linked transaction if the latest occurrence has a transactionId
            const firstOccurrence = occurrences[0];
            if (firstOccurrence?.transactionId) {
                try {
                    const txResponse = await api.post(
                        `/transactions/${firstOccurrence.transactionId}`,
                        {},
                        { projectId: projectsState.currentProjectId ?? undefined }
                    );
                    if (txResponse.transaction) {
                        linkedTransaction = {
                            id: txResponse.transaction.id,
                            endpoint: txResponse.transaction.endpoint,
                            duration: txResponse.transaction.duration,
                            statusCode: txResponse.transaction.statusCode,
                            recordedAt: txResponse.transaction.recordedAt
                        };
                    }
                } catch (txError) {
                    console.warn('Could not load linked transaction:', txError);
                }
            }
        } catch (e: any) {
            console.error(e);
            if (e.status === 404) {
                notFound = true;
            } else {
                error = e.message || 'Failed to load exception details';
            }
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        loadData();
    });
</script>

<div class="space-y-6">
    <PageHeader
        title={firstLineOfStackTrace}
        subtitle="Exception Hash: {exceptionHash}"
        backHref="/issues"
    />

    {#if loading && !group}
        <div class="flex items-center justify-center py-20">
            <LoadingCircle size="xlg" />
        </div>
    {:else if notFound}
        <ErrorDisplay
            status={404}
            title="Exception Not Found"
            description="The exception you're looking for doesn't exist or may have been removed. It's possible the data has expired or the link is incorrect."
            backHref="/issues"
            backLabel="Back to Issues"
            onRetry={() => loadData()}
            identifier={exceptionHash}
        />
    {:else if error}
        <ErrorDisplay
            status={400}
            title="Something Went Wrong"
            description={error}
            backHref="/issues"
            backLabel="Back to Issues"
            onRetry={() => loadData()}
        />
    {:else if group}
        <StackTraceCard
            stackTrace={group.stackTrace}
            {isMessage}
            firstSeen={group.firstSeen}
            lastSeen={group.lastSeen}
            totalCount={group.count}
        />

        {#if latestOccurrence}
            <EventCard
                occurrence={latestOccurrence}
                {linkedTransaction}
                title="Last Event"
                description="Details from the most recent occurrence of this exception"
            />
        {/if}

        <EventsTable
            {occurrences}
            {exceptionHash}
            {total}
            hasMore={hasMoreOccurrences}
            showViewAll={true}
        />
    {/if}
</div>
