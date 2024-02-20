<script>
	import { isLoggedIn } from '$lib/utils';
	import { goto } from '$app/navigation';
	import { Search, HouseFill, Grid3x3GapFill, EyeFill } from 'svelte-bootstrap-icons';

	function redirectPFP() {
		if (isLoggedIn()) {
			goto('/user');
		} else {
			goto('/login');
		}
	}
</script>

<div id="sidebar">
	<div id="user-pfp" on:click={redirectPFP}></div>
	<div id="sidebar-entries">
		<div class="sidebar-entry">
			<div class="sidebar-entry-icon"><Search class="sidebar-entry-icon-svg" /></div>
			<div id="sidebar-search-box" class="sidebar-entry-text">Search...</div>
		</div>
		<div
			class="sidebar-entry"
			on:click={() => {
				goto('/');
			}}
		>
			<div class="sidebar-entry-icon"><HouseFill class="sidebar-entry-icon-svg" /></div>
			<div class="sidebar-entry-text">Home</div>
		</div>
		<div
			class="sidebar-entry"
			on:click={() => {
				goto('/observations');
			}}
		>
			<div class="sidebar-entry-icon"><EyeFill class="sidebar-entry-icon-svg" /></div>
			<div class="sidebar-entry-text">Observations</div>
		</div>
		<div
			class="sidebar-entry"
			on:click={() => {
				goto('/explore');
			}}
		>
			<div class="sidebar-entry-icon"><Grid3x3GapFill class="sidebar-entry-icon-svg" /></div>
			<div class="sidebar-entry-text">Explore</div>
		</div>
	</div>
</div>
