<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { Button } from '$lib/components/ui/button';
	import { Skeleton } from '$lib/components/ui/skeleton';
	import { RefreshCw } from 'lucide-svelte';
	import MetricCard from '$lib/components/dashboard/metric-card.svelte';
	import type { DashboardData, DashboardMetric } from '$lib/types/dashboard';
	import { Card, CardContent } from '$lib/components/ui/card';
	import { api } from '$lib/api';
	import { ErrorDisplay } from '$lib/components/ui/error-display';
	import { projectsState } from '$lib/state/projects.svelte';
	import { TimeRangePicker } from '$lib/components/ui/time-range-picker';
	import { CalendarDate } from '@internationalized/date';

	let dashboardData = $state<DashboardData | null>(null);
	let loading = $state(true);
	let error = $state('');
	let errorStatus = $state<number>(0);

	// Preset definitions (must match TimeRangePicker)
	const presetMinutes: Record<string, number> = {
		'30m': 30,
		'60m': 60,
		'3h': 180,
		'6h': 360,
		'12h': 720,
		'24h': 1440,
		'3d': 4320,
		'7d': 10080,
		'1M': 43200,
		'3M': 129600,
	};

	// Calculate time range from preset
	function getTimeRangeFromPreset(presetValue: string): { from: Date; to: Date } {
		const minutes = presetMinutes[presetValue] || 360; // Default to 6h
		const now = new Date();
		const from = new Date(now.getTime() - minutes * 60 * 1000);
		return { from, to: now };
	}

	// Parse URL params
	function parseUrlParams(): { preset: string | null; from: Date | null; to: Date | null } {
		if (!browser) return { preset: '6h', from: null, to: null };
		const params = new URLSearchParams(window.location.search);
		const presetParam = params.get('preset');
		const fromParam = params.get('from');
		const toParam = params.get('to');

		// If preset is specified, use it
		if (presetParam && presetMinutes[presetParam]) {
			return { preset: presetParam, from: null, to: null };
		}

		// If custom from/to specified
		if (fromParam && toParam) {
			const from = new Date(fromParam);
			const to = new Date(toParam);
			if (!isNaN(from.getTime()) && !isNaN(to.getTime())) {
				return { preset: null, from, to };
			}
		}

		// Default to 6h preset
		return { preset: '6h', from: null, to: null };
	}

	function dateToCalendarDate(date: Date): CalendarDate {
		return new CalendarDate(date.getFullYear(), date.getMonth() + 1, date.getDate());
	}

	function dateToTimeString(date: Date): string {
		return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
	}

	// Initialize from URL or defaults
	const initialUrlParams = parseUrlParams();
	const initialRange = initialUrlParams.preset
		? getTimeRangeFromPreset(initialUrlParams.preset)
		: { from: initialUrlParams.from!, to: initialUrlParams.to! };

	let selectedPreset = $state<string | null>(initialUrlParams.preset);
	let fromDate = $state<CalendarDate>(dateToCalendarDate(initialRange.from));
	let toDate = $state<CalendarDate>(dateToCalendarDate(initialRange.to));
	let fromTime = $state(dateToTimeString(initialRange.from));
	let toTime = $state(dateToTimeString(initialRange.to));
	let sharedTimeDomain = $state<[Date, Date] | null>(null);

	// Update URL with current time range
	function updateUrl(pushState = true) {
		if (!browser) return;

		const params = new URLSearchParams();

		if (selectedPreset) {
			// Store preset in URL
			params.set('preset', selectedPreset);
		} else {
			// Store custom from/to in URL
			const fromDateTime = new Date(getFromDateTime());
			const toDateTime = new Date(getToDateTime());
			params.set('from', fromDateTime.toISOString());
			params.set('to', toDateTime.toISOString());
		}

		const newUrl = `${window.location.pathname}?${params.toString()}`;

		if (pushState) {
			window.history.pushState({}, '', newUrl);
		} else {
			window.history.replaceState({}, '', newUrl);
		}
	}

	// Handle browser back/forward navigation
	function handlePopState() {
		const urlParams = parseUrlParams();

		if (urlParams.preset) {
			selectedPreset = urlParams.preset;
			const range = getTimeRangeFromPreset(urlParams.preset);
			fromDate = dateToCalendarDate(range.from);
			fromTime = dateToTimeString(range.from);
			toDate = dateToCalendarDate(range.to);
			toTime = dateToTimeString(range.to);
		} else if (urlParams.from && urlParams.to) {
			selectedPreset = null;
			fromDate = dateToCalendarDate(urlParams.from);
			fromTime = dateToTimeString(urlParams.from);
			toDate = dateToCalendarDate(urlParams.to);
			toTime = dateToTimeString(urlParams.to);
		}

		loadDashboard(false); // Don't push to history on popstate
	}

	// Combine date and time into ISO datetime string
	function getFromDateTime(): string {
		const dateStr = `${fromDate.year}-${String(fromDate.month).padStart(2, '0')}-${String(fromDate.day).padStart(2, '0')}`;
		return `${dateStr}T${fromTime || '00:00'}`;
	}

	function getToDateTime(): string {
		const dateStr = `${toDate.year}-${String(toDate.month).padStart(2, '0')}-${String(toDate.day).padStart(2, '0')}`;
		return `${dateStr}T${toTime || '23:59'}`;
	}

	function handleTimeRangeChange(from: { date: CalendarDate; time: string }, to: { date: CalendarDate; time: string }, preset: string | null) {
		fromDate = from.date;
		fromTime = from.time;
		toDate = to.date;
		toTime = to.time;
		selectedPreset = preset;
		loadDashboard();
	}

	// Handle drag-to-zoom selection from chart overlay
	function handleChartRangeSelect(from: Date, to: Date) {
		selectedPreset = null; // Chart selection is always custom
		fromDate = new CalendarDate(from.getFullYear(), from.getMonth() + 1, from.getDate());
		fromTime = `${String(from.getHours()).padStart(2, '0')}:${String(from.getMinutes()).padStart(2, '0')}`;
		toDate = new CalendarDate(to.getFullYear(), to.getMonth() + 1, to.getDate());
		toTime = `${String(to.getHours()).padStart(2, '0')}:${String(to.getMinutes()).padStart(2, '0')}`;
		loadDashboard();
	}

	async function loadDashboard(pushToHistory = true) {
		loading = true;
		error = '';
		errorStatus = 0;

		// Update URL with current time range
		updateUrl(pushToHistory);

		try {
			const fromDateTime = new Date(getFromDateTime());
			const toDateTime = new Date(getToDateTime());

			// Store shared time domain for charts
			sharedTimeDomain = [fromDateTime, toDateTime];

			// Build query params for date range
			const dateParams = `fromDate=${fromDateTime.toISOString()}&toDate=${toDateTime.toISOString()}`;
			const response = await api.get(`/dashboard?${dateParams}`, {
				projectId: projectsState.currentProjectId ?? undefined
			});

			// Transform the API response to match the DashboardData type
			// Convert timestamp strings to Date objects
			const metrics: DashboardMetric[] = response.metrics.map((m: any) => ({
				id: m.id,
				name: m.name,
				value: m.value,
				unit: m.unit,
				trend: m.trend.map((t: any) => ({
					timestamp: new Date(t.timestamp),
					value: t.value
				})),
				change24h: m.change24h,
				status: m.status
			}));

			dashboardData = {
				metrics,
				lastUpdated: new Date(response.lastUpdated)
			};
		} catch (e: any) {
			errorStatus = e.status || 0;
			error = e.message || 'Failed to load dashboard data';
			console.error(e);
		} finally {
			loading = false;
		}
	}

	onMount(() => {
		// On initial load, use replaceState to set URL without creating history entry
		loadDashboard(false);

		// Listen for browser back/forward navigation
		window.addEventListener('popstate', handlePopState);

		return () => {
			window.removeEventListener('popstate', handlePopState);
		};
	});

	const lastUpdatedFormatted = $derived(
		dashboardData?.lastUpdated
			? dashboardData.lastUpdated.toLocaleString('en-US', {
					hour: '2-digit',
					minute: '2-digit',
					second: '2-digit',
					hour12: false
				})
			: ''
	);
</script>

<div class="space-y-4">
	<!-- Header -->
	<div class="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
		<div>
			<h2 class="text-3xl font-bold tracking-tight">Metrics</h2>
			{#if dashboardData}
				<p class="mt-1 text-sm text-muted-foreground">Last updated: {lastUpdatedFormatted}</p>
			{/if}
		</div>
		<div class="flex items-center gap-2">
			<TimeRangePicker
				bind:fromDate
				bind:toDate
				bind:fromTime
				bind:toTime
				bind:preset={selectedPreset}
				onApply={handleTimeRangeChange}
			/>
			<Button variant="outline" size="sm" onclick={() => loadDashboard()} disabled={loading}>
				<RefreshCw class="h-4 w-4 {loading ? 'animate-spin' : ''}" />
			</Button>
		</div>
	</div>

	{#if error && !loading}
		<ErrorDisplay
			status={errorStatus === 404 ? 404 : errorStatus === 400 ? 400 : errorStatus === 422 ? 422 : 400}
			title="Failed to Load Metrics"
			description={error}
			onRetry={() => loadDashboard()}
		/>
	{/if}

	<!-- Metrics Grid -->
	{#if !error}
	<div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
		{#if loading}
			{#each Array(11) as _}
				<Card>
					<CardContent class="flex gap-3 flex-col">
						<div class="flex items-center justify-between">
							<Skeleton class="h-4 w-24" />
							<Skeleton class="h-2 w-2 rounded-full" />
						</div>
						<Skeleton class="h-8 w-32" />
						<Skeleton class="h-16 w-full" />
						<Skeleton class="h-3 w-28" />
					</CardContent>
				</Card>
			{/each}
		{:else if dashboardData}
			{#each dashboardData.metrics as metric (metric.id)}
				<MetricCard {metric} timeDomain={sharedTimeDomain} onRangeSelect={handleChartRangeSelect} />
			{/each}
		{/if}
	</div>
	{/if}
</div>
