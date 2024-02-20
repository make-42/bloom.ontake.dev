<script>
	import { goto } from '$app/navigation';
	import { apiURL } from '$lib/utils';

	let username = '';
	let email = '';
	let lat = 0;
	let long = 0;
	let password = '';
	let repeatPassword = '';
	let permissions = 4;

	let result = null;
	let errorMessage = '';
	let error = false;
	async function register() {
		if (password != repeatPassword) {
			error = true;
			errorMessage = "the passwords don't match";
		} else {
			const res = await fetch(apiURL + 'users/user', {
				method: 'POST',
				body: JSON.stringify({
					username,
					password,
					email,
					lat,
					long,
					permissions
				})
			});

			const json = await res.json();

			if (!res.ok) {
				error = true;
				errorMessage = json.error;
			} else {
				localStorage.token = json.token;
				goto('/login');
			}
		}
	}
</script>

<div id="login-prompt-container">
	<div id="login-prompt">
		<form>
			<div class="login-field-container">
				<input class="login-field" bind:value={username} type="text" placeholder="username" />
			</div>
			<div class="login-field-container">
				<input class="login-field" bind:value={email} type="email" placeholder="email" />
			</div>
			<div class="login-field-container">
				<input class="login-field" bind:value={password} type="password" placeholder="password" />
			</div>
			<div class="login-field-container">
				<input
					class="login-field"
					bind:value={repeatPassword}
					type="password"
					placeholder="repeat your password"
				/>
			</div>
			<div class="login-field-container">
				<label>Latitude</label><br />
				<input
					class="login-field"
					bind:value={lat}
					type="number"
					min="-90"
					max="90"
					step="0.001"
					placeholder="latitude"
				/>
			</div>
			<div class="login-field-container">
				<label>Longitude</label><br />
				<input
					class="login-field"
					bind:value={long}
					type="number"
					min="-180"
					max="180"
					step="0.001"
					placeholder="longitude"
				/>
			</div>
			<div id="login-button-container">
				<button id="login-button" on:click={register}>Register</button>
			</div>
			<div id="login-error-message-container" style={`padding: ${error ? '8px' : '0px'}`}>
				{errorMessage}
			</div>
			<div id="register-message">
				Have an account already? <a href="/login">Login</a> instead.
			</div>
		</form>
	</div>
</div>
