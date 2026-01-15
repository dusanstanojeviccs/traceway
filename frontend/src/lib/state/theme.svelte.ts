export const themeState = $state({
    isDark: false
});

export function initTheme() {
    if (typeof document !== 'undefined') {
        // Check localStorage first
        const stored = localStorage.getItem('theme');

        if (stored === 'dark' || stored === 'light') {
            // Use stored preference
            themeState.isDark = stored === 'dark';
        } else {
            // Fall back to system preference
            themeState.isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
        }

        // Apply the theme
        document.documentElement.classList.toggle('dark', themeState.isDark);

        // Watch for system preference changes (only matters if no stored preference)
        const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
        const handleSystemChange = (e: MediaQueryListEvent) => {
            if (!localStorage.getItem('theme')) {
                themeState.isDark = e.matches;
                document.documentElement.classList.toggle('dark', themeState.isDark);
            }
        };
        mediaQuery.addEventListener('change', handleSystemChange);

        // Watch for external class changes
        const observer = new MutationObserver(() => {
            themeState.isDark = document.documentElement.classList.contains('dark');
        });
        observer.observe(document.documentElement, {
            attributes: true,
            attributeFilter: ['class']
        });

        return () => {
            observer.disconnect();
            mediaQuery.removeEventListener('change', handleSystemChange);
        };
    }
}

export function toggleTheme() {
    themeState.isDark = !themeState.isDark;
    document.documentElement.classList.toggle('dark', themeState.isDark);
    localStorage.setItem('theme', themeState.isDark ? 'dark' : 'light');
}
