export function formatDuration(nanoseconds: number): string {
	const ms = nanoseconds / 1_000_000;
	if (ms < 1) {
		return `${(nanoseconds / 1000).toFixed(0)}µs`;
	} else if (ms < 1000) {
		return `${ms.toFixed(0)}ms`;
	} else {
		return `${(ms / 1000).toFixed(1)}s`;
	}
}

export function formatDurationMs(ms: number): string {
	if (ms < 1) {
		return `${(ms * 1000).toFixed(0)} µs`;
	} else if (ms < 1000) {
		return `${ms.toFixed(0)} ms`;
	} else {
		return `${(ms / 1000).toFixed(2)} s`;
	}
}

export function getStatusColor(statusCode: number): string {
	if (statusCode >= 200 && statusCode < 300) return 'text-green-500';
	if (statusCode >= 300 && statusCode < 400) return 'text-blue-500';
	if (statusCode >= 400 && statusCode < 500) return 'text-yellow-500';
	return 'text-red-500';
}

export function truncateStackTrace(stackTrace: string, maxLength = 70): string {
	const firstLine = stackTrace.split('\n')[0];
	if (firstLine.length > maxLength) {
		return firstLine.slice(0, maxLength) + '...';
	}
	return firstLine;
}

export function formatRelativeTime(dateStr: string): string {
	const date = new Date(dateStr);
	const now = new Date();
	const diffMs = now.getTime() - date.getTime();
	const diffMins = Math.floor(diffMs / 60000);
	const diffHours = Math.floor(diffMs / 3600000);
	const diffDays = Math.floor(diffMs / 86400000);

	if (diffMins < 1) return 'just now';
	if (diffMins < 60) return `${diffMins}m`;
	if (diffHours < 24) return `${diffHours}h`;
	return `${diffDays}d`;
}
