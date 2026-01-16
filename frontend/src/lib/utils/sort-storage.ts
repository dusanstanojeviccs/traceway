import { browser } from '$app/environment';

export type SortDirection = 'asc' | 'desc';

export type SortState = {
	field: string;
	direction: SortDirection;
};

const STORAGE_PREFIX = 'traceway_sort_';

export function getSortState(pageKey: string, defaultState: SortState): SortState {
	if (!browser) return defaultState;

	try {
		const stored = localStorage.getItem(`${STORAGE_PREFIX}${pageKey}`);
		if (stored) {
			const parsed = JSON.parse(stored);
			if (parsed.field && (parsed.direction === 'asc' || parsed.direction === 'desc')) {
				return parsed;
			}
		}
	} catch {
		// Ignore parse errors, return default
	}

	return defaultState;
}

export function setSortState(pageKey: string, state: SortState): void {
	if (!browser) return;

	try {
		localStorage.setItem(`${STORAGE_PREFIX}${pageKey}`, JSON.stringify(state));
	} catch {
		// Ignore storage errors
	}
}

export function toggleSortDirection(current: SortDirection): SortDirection {
	return current === 'asc' ? 'desc' : 'asc';
}

export function handleSortClick(
	field: string,
	currentField: string,
	currentDirection: SortDirection,
	defaultDirection: SortDirection = 'desc'
): SortState {
	if (field === currentField) {
		return { field, direction: toggleSortDirection(currentDirection) };
	}
	return { field, direction: defaultDirection };
}
