<script lang="ts">
    import { goto } from '$app/navigation';
    import { authState } from '$lib/state/auth.svelte';
    
    $effect(() => {
            if (authState.token) {
            goto('/issues');
        } else {
            goto('/login');
        }
    });
</script>
