<script lang="ts">
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { register, USER_TYPE, type RegistrationRole } from '$lib/auth/store';

	let email = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let givenName = $state('');
	let familyName = $state('');
	let userType = $state<RegistrationRole | ''>('');
	let error = $state('');
	let loading = $state(false);
	let emailSent = $state(false);

	const lang = $derived($page.params.lang);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		error = '';

		if (!userType) {
			error = $t('auth.select_role');
			return;
		}

		if (password !== confirmPassword) {
			error = $t('auth.passwords_mismatch');
			return;
		}

		loading = true;

		try {
			const result = await register(email, password, givenName, familyName, userType);
			if (result.emailVerificationSent) {
				emailSent = true;
			}
		} catch (err: any) {
			error = $t('auth.register_failed');
		} finally {
			loading = false;
		}
	}
</script>

<div class="page-centered">
	<div class="form-card">
		{#if emailSent}
			<div class="verify-email">
				<h1 class="heading heading--lg heading--center">
					{$t('auth.check_email_title')}
				</h1>
				<p class="verify-email__message">
					{$t('auth.check_email_message')}
				</p>
				<a href="{base}/{lang}/login" class="btn btn--primary">
					{$t('auth.back_to_login')}
				</a>
			</div>
		{:else}
			<h1 class="heading heading--lg heading--center">
				{$t('auth.register_title')}
			</h1>

			<form onsubmit={handleSubmit} class="form-stack">
				<fieldset class="role-selector">
					<legend class="field__label">{$t('auth.select_role')}</legend>
					<div class="role-selector__options">
						<label class="role-option" class:role-option--selected={userType === USER_TYPE.RESEARCHER}>
							<input
								type="radio"
								name="user-type"
								value={USER_TYPE.RESEARCHER}
								bind:group={userType}
								class="role-option__input"
							/>
							<span class="role-option__label">{$t('auth.role_researcher')}</span>
							<span class="role-option__desc">{$t('auth.role_researcher_desc')}</span>
						</label>
						<label class="role-option" class:role-option--selected={userType === USER_TYPE.SCITIZEN}>
							<input
								type="radio"
								name="user-type"
								value={USER_TYPE.SCITIZEN}
								bind:group={userType}
								class="role-option__input"
							/>
							<span class="role-option__label">{$t('auth.role_scitizen')}</span>
							<span class="role-option__desc">{$t('auth.role_scitizen_desc')}</span>
						</label>
					</div>
				</fieldset>

				<div class="form-row">
					<div class="field">
						<label for="given-name" class="field__label">{$t('auth.given_name')}</label>
						<input
							id="given-name"
							type="text"
							bind:value={givenName}
							required
							autocomplete="given-name"
							class="field__input"
						/>
					</div>
					<div class="field">
						<label for="family-name" class="field__label">{$t('auth.family_name')}</label>
						<input
							id="family-name"
							type="text"
							bind:value={familyName}
							required
							autocomplete="family-name"
							class="field__input"
						/>
					</div>
				</div>

				<div class="field">
					<label for="email" class="field__label">{$t('auth.email')}</label>
					<input
						id="email"
						type="email"
						bind:value={email}
						required
						autocomplete="email"
						class="field__input"
					/>
				</div>

				<div class="field">
					<label for="password" class="field__label">{$t('auth.password')}</label>
					<input
						id="password"
						type="password"
						bind:value={password}
						required
						autocomplete="new-password"
						class="field__input"
					/>
				</div>

				<div class="field">
					<label for="confirm-password" class="field__label">{$t('auth.confirm_password')}</label>
					<input
						id="confirm-password"
						type="password"
						bind:value={confirmPassword}
						required
						autocomplete="new-password"
						class="field__input"
					/>
				</div>

				{#if error}
					<p class="form-error" role="alert">{error}</p>
				{/if}

				<button
					type="submit"
					disabled={loading}
					class="btn btn--primary mt-2"
				>
					{loading ? $t('common.loading') : $t('auth.register_button')}
				</button>
			</form>

			<p class="form-footer">
				{$t('auth.have_account')}
				<a href="{base}/{lang}/login" class="form-footer__link">
					{$t('auth.login_link')}
				</a>
			</p>
		{/if}
	</div>
</div>
