<script>
	import { goto } from '$app/navigation';
	import { apiURL } from '$lib/utils';

	let username = '';
	let password = '';

	let result = null;
	let errorMessage = '';
	let error = false;
	async function login() {
		const res = await fetch(apiURL + 'users/login', {
			method: 'POST',
			body: JSON.stringify({
				username,
				password
			})
		});

		const json = await res.json();

		if (!res.ok) {
			error = true;
			errorMessage = json.error;
		} else {
			localStorage.token = json.token;
			goto('/');
		}
	}
</script>

<div id="login-prompt-container">
	<div id="login-prompt">
		<form>
			<div class="login-field-container">
				<input class="login-field" bind:value={username} placeholder="username" />
			</div>
			<div class="login-field-container">
				<input class="login-field" bind:value={password} type="password" placeholder="password" />
			</div>
			<div id="login-button-container">
				<button id="login-button" on:click={login}>Login</button>
			</div>
			<div id="login-error-message-container" style={`padding: ${error ? '8px' : '0px'}`}>
				{errorMessage}
			</div>
			<div id="register-message">
				Don't have an account? <a href="/register">Register</a> instead.
			</div>
		</form>
	</div>
</div>
