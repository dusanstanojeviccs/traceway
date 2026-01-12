import { DateTime } from 'luxon';

export const timezoneState = $state({
	timezone: ''
});

export function initTimezone() {
	if (typeof window !== 'undefined') {
		const stored = localStorage.getItem('timezone');
		timezoneState.timezone = stored || DateTime.local().zoneName || 'UTC';
	}
}

export function getTimezone(): string {
	return timezoneState.timezone || DateTime.local().zoneName || 'UTC';
}

export function setTimezone(tz: string) {
	timezoneState.timezone = tz;
	localStorage.setItem('timezone', tz);
}
