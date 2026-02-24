<script lang="ts">
	// Graph node: 0x40 (ScitizenRegisterPage)
	// Implements: FR-011 (Registration), FR-080 (ToS Acceptance)
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { t } from '$lib/i18n';
	import { userService } from '$lib/api/clients';

	const lang = $derived($page.params.lang);

	let email = $state('');
	let password = $state('');
	let givenName = $state('');
	let familyName = $state('');
	let tosAccepted = $state(false);
	let loading = $state(false);
	let error = $state('');
	let success = $state(false);

	async function handleSubmit() {
		if (!tosAccepted) {
			error = 'You must accept the Terms of Service';
			return;
		}

		loading = true;
		error = '';
		try {
			await userService.registerResearcher({
				email,
				password,
				givenName,
				familyName,
			});
			success = true;
		} catch (e: any) {
			error = e.message || 'Registration failed';
		} finally {
			loading = false;
		}
	}
</script>

<div class="register-page">
	{#if success}
		<div class="success-state">
			<h1 class="heading heading--xl">Check your email</h1>
			<p>We've sent a verification link to <strong>{email}</strong>.</p>
			<a href="{base}/{lang}/login" class="btn btn--primary mt-4">Go to Login</a>
		</div>
	{:else}
		<h1 class="heading heading--xl">Join as Citizen Scientist</h1>
		<p class="text-secondary mb-4">Register to contribute environmental data.</p>

		<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="register-form">
			<div class="form-group">
				<label for="given-name">First Name</label>
				<input id="given-name" type="text" bind:value={givenName} required class="input" autocomplete="given-name" />
			</div>
			<div class="form-group">
				<label for="family-name">Last Name</label>
				<input id="family-name" type="text" bind:value={familyName} required class="input" autocomplete="family-name" />
			</div>
			<div class="form-group">
				<label for="email">Email</label>
				<input id="email" type="email" bind:value={email} required class="input" autocomplete="email" />
			</div>
			<div class="form-group">
				<label for="password">Password</label>
				<input id="password" type="password" bind:value={password} required class="input" minlength="8" autocomplete="new-password" />
			</div>
			<div class="form-group form-group--checkbox">
				<input id="tos" type="checkbox" bind:checked={tosAccepted} />
				<label for="tos">I accept the Terms of Service and Privacy Policy</label>
			</div>

			{#if error}
				<p class="form-error" role="alert">{error}</p>
			{/if}

			<button type="submit" class="btn btn--primary" disabled={loading}>
				{loading ? 'Registering...' : 'Register'}
			</button>

			<p class="mt-4 text-secondary">
				Already have an account? <a href="{base}/{lang}/login">Log in</a>
			</p>
		</form>
	{/if}
</div>

<style>
	.register-page {
		max-width: 480px;
		margin: 2rem auto;
		padding: 0 1rem;
	}
	.register-form {
		display: flex;
		flex-direction: column;
		gap: 1rem;
	}
	.form-group {
		display: flex;
		flex-direction: column;
		gap: 0.25rem;
	}
	.form-group--checkbox {
		flex-direction: row;
		align-items: center;
		gap: 0.5rem;
	}
	.success-state {
		text-align: center;
		padding: 3rem 1rem;
	}
</style>
