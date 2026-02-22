<script lang="ts">
	import { goto } from '$app/navigation';
	import { base } from '$app/paths';
	import { page } from '$app/stores';
	import { t } from '$lib/i18n';
	import { campaignService } from '$lib/api/clients';
	import { currentUser } from '$lib/auth/store';

	const lang = $derived($page.params.lang);

	/* -------------------------------------------------------
	   Wizard state
	   Steps: 1=basics, 2=parameters, 3=regions, 4=eligibility, 5=review
	   Graph refs: 0x1b (campaigns), 0x1e (parameters), 0x18 (regions),
	               0x12 (eligibility), 0x26 (draft state)
	   ------------------------------------------------------- */

	let currentStep = $state(1);
	const totalSteps = 5;

	const stepLabels = ['Basics', 'Parameters', 'Regions', 'Eligibility', 'Review'];

	// Step 1: Basics (graph 0x1b)
	let windowStart = $state('');
	let windowEnd = $state('');

	// Step 2: Parameters (graph 0x1e)
	interface Parameter {
		name: string;
		unit: string;
		minRange: string;
		maxRange: string;
		precision: string;
	}
	let parameters = $state<Parameter[]>([{ name: '', unit: '', minRange: '', maxRange: '', precision: '' }]);

	function addParameter() {
		parameters = [...parameters, { name: '', unit: '', minRange: '', maxRange: '', precision: '' }];
	}

	function removeParameter(index: number) {
		parameters = parameters.filter((_, i) => i !== index);
	}

	// Step 3: Regions (graph 0x18)
	interface Region {
		geoJson: string;
	}
	let regions = $state<Region[]>([{ geoJson: '' }]);

	function addRegion() {
		regions = [...regions, { geoJson: '' }];
	}

	function removeRegion(index: number) {
		regions = regions.filter((_, i) => i !== index);
	}

	// Step 4: Eligibility (graph 0x12)
	interface Eligibility {
		deviceClass: string;
		tier: string;
		requiredSensors: string;
		firmwareMin: string;
	}
	let eligibility = $state<Eligibility[]>([{ deviceClass: '', tier: '1', requiredSensors: '', firmwareMin: '' }]);

	function addEligibility() {
		eligibility = [...eligibility, { deviceClass: '', tier: '1', requiredSensors: '', firmwareMin: '' }];
	}

	function removeEligibility(index: number) {
		eligibility = eligibility.filter((_, i) => i !== index);
	}

	// Submission
	let submitting = $state(false);
	let error = $state('');

	async function handleSubmit() {
		submitting = true;
		error = '';

		try {
			await campaignService.createCampaign({
				orgId: '',
				createdBy: $currentUser?.id ?? '',
				windowStart: windowStart || undefined,
				windowEnd: windowEnd || undefined,
				parameters: parameters
					.filter((p) => p.name)
					.map((p) => ({
						name: p.name,
						unit: p.unit,
						minRange: p.minRange ? parseFloat(p.minRange) : undefined,
						maxRange: p.maxRange ? parseFloat(p.maxRange) : undefined,
						precision: p.precision ? parseInt(p.precision) : undefined,
					})),
				regions: regions
					.filter((r) => r.geoJson)
					.map((r) => ({ geoJson: r.geoJson })),
				eligibility: eligibility
					.filter((e) => e.deviceClass)
					.map((e) => ({
						deviceClass: e.deviceClass,
						tier: parseInt(e.tier) || 1,
						requiredSensors: e.requiredSensors
							.split(',')
							.map((s) => s.trim())
							.filter(Boolean),
						firmwareMin: e.firmwareMin,
					})),
			});

			await goto(`${base}/${lang}/researcher/`);
		} catch (e: any) {
			error = $t('campaign.create_error');
		} finally {
			submitting = false;
		}
	}

	function nextStep() {
		if (currentStep < totalSteps) currentStep++;
	}

	function prevStep() {
		if (currentStep > 1) currentStep--;
	}
</script>

<div class="wizard">
	<!-- Step indicator -->
	<div class="wizard__steps">
		{#each stepLabels as label, i}
			{#if i > 0}
				<div class="wizard__step-divider"></div>
			{/if}
			<div class="wizard__step" class:wizard__step--active={currentStep === i + 1} class:wizard__step--completed={currentStep > i + 1}>
				<span>{i + 1}.</span>
				<span>{$t(`campaign.${label.toLowerCase()}`) || label}</span>
			</div>
		{/each}
	</div>

	<!-- Step 1: Basics -->
	{#if currentStep === 1}
		<div class="wizard__section">
			<h2 class="wizard__section-title">{$t('campaign.basics')}</h2>
			<div class="form-stack">
				<div class="form-row">
					<div class="field">
						<label for="window-start" class="field__label">{$t('campaign.window_start')}</label>
						<input
							id="window-start"
							type="date"
							bind:value={windowStart}
							class="field__input"
						/>
					</div>
					<div class="field">
						<label for="window-end" class="field__label">{$t('campaign.window_end')}</label>
						<input
							id="window-end"
							type="date"
							bind:value={windowEnd}
							class="field__input"
						/>
					</div>
				</div>
			</div>
		</div>

	<!-- Step 2: Parameters -->
	{:else if currentStep === 2}
		<div class="wizard__section">
			<h2 class="wizard__section-title">{$t('campaign.parameters')}</h2>
			<div class="repeater">
				{#each parameters as param, i}
					<div class="repeater__item">
						<div class="repeater__item-header">
							<span class="field__label">{$t('campaign.parameters')} {i + 1}</span>
							{#if parameters.length > 1}
								<button type="button" class="repeater__remove" onclick={() => removeParameter(i)}>
									{$t('campaign.remove')}
								</button>
							{/if}
						</div>
						<div class="form-stack">
							<div class="form-row">
								<div class="field">
									<label for="param-name-{i}" class="field__label">{$t('campaign.param_name')}</label>
									<input id="param-name-{i}" type="text" bind:value={param.name} class="field__input" required />
								</div>
								<div class="field">
									<label for="param-unit-{i}" class="field__label">{$t('campaign.param_unit')}</label>
									<input id="param-unit-{i}" type="text" bind:value={param.unit} class="field__input" />
								</div>
							</div>
							<div class="form-row">
								<div class="field">
									<label for="param-min-{i}" class="field__label">{$t('campaign.param_min')}</label>
									<input id="param-min-{i}" type="number" step="any" bind:value={param.minRange} class="field__input" />
								</div>
								<div class="field">
									<label for="param-max-{i}" class="field__label">{$t('campaign.param_max')}</label>
									<input id="param-max-{i}" type="number" step="any" bind:value={param.maxRange} class="field__input" />
								</div>
								<div class="field">
									<label for="param-prec-{i}" class="field__label">{$t('campaign.param_precision')}</label>
									<input id="param-prec-{i}" type="number" min="0" bind:value={param.precision} class="field__input" />
								</div>
							</div>
						</div>
					</div>
				{/each}
				<button type="button" class="btn btn--secondary repeater__add" onclick={addParameter}>
					{$t('campaign.add_parameter')}
				</button>
			</div>
		</div>

	<!-- Step 3: Regions -->
	{:else if currentStep === 3}
		<div class="wizard__section">
			<h2 class="wizard__section-title">{$t('campaign.regions')}</h2>
			<div class="repeater">
				{#each regions as region, i}
					<div class="repeater__item">
						<div class="repeater__item-header">
							<span class="field__label">{$t('campaign.regions')} {i + 1}</span>
							{#if regions.length > 1}
								<button type="button" class="repeater__remove" onclick={() => removeRegion(i)}>
									{$t('campaign.remove')}
								</button>
							{/if}
						</div>
						<div class="field">
							<label for="region-json-{i}" class="field__label">{$t('campaign.region_geojson')}</label>
							<textarea
								id="region-json-{i}"
								bind:value={region.geoJson}
								class="field__input"
								rows="4"
								placeholder={'{"type": "Polygon", "coordinates": [...]}'}

							></textarea>
						</div>
					</div>
				{/each}
				<button type="button" class="btn btn--secondary repeater__add" onclick={addRegion}>
					{$t('campaign.add_region')}
				</button>
			</div>
		</div>

	<!-- Step 4: Eligibility -->
	{:else if currentStep === 4}
		<div class="wizard__section">
			<h2 class="wizard__section-title">{$t('campaign.eligibility')}</h2>
			<div class="repeater">
				{#each eligibility as rule, i}
					<div class="repeater__item">
						<div class="repeater__item-header">
							<span class="field__label">{$t('campaign.eligibility')} {i + 1}</span>
							{#if eligibility.length > 1}
								<button type="button" class="repeater__remove" onclick={() => removeEligibility(i)}>
									{$t('campaign.remove')}
								</button>
							{/if}
						</div>
						<div class="form-stack">
							<div class="form-row">
								<div class="field">
									<label for="elig-class-{i}" class="field__label">{$t('campaign.device_class')}</label>
									<input id="elig-class-{i}" type="text" bind:value={rule.deviceClass} class="field__input" required />
								</div>
								<div class="field">
									<label for="elig-tier-{i}" class="field__label">{$t('campaign.tier')}</label>
									<input id="elig-tier-{i}" type="number" min="1" bind:value={rule.tier} class="field__input" />
								</div>
							</div>
							<div class="form-row">
								<div class="field">
									<label for="elig-sensors-{i}" class="field__label">{$t('campaign.required_sensors')}</label>
									<input id="elig-sensors-{i}" type="text" bind:value={rule.requiredSensors} class="field__input" placeholder="temp, humidity, pressure" />
								</div>
								<div class="field">
									<label for="elig-fw-{i}" class="field__label">{$t('campaign.firmware_min')}</label>
									<input id="elig-fw-{i}" type="text" bind:value={rule.firmwareMin} class="field__input" placeholder="1.0.0" />
								</div>
							</div>
						</div>
					</div>
				{/each}
				<button type="button" class="btn btn--secondary repeater__add" onclick={addEligibility}>
					{$t('campaign.add_eligibility')}
				</button>
			</div>
		</div>

	<!-- Step 5: Review -->
	{:else if currentStep === 5}
		<div class="wizard__section">
			<h2 class="wizard__section-title">{$t('campaign.review')}</h2>

			<div class="form-stack">
				<div class="campaign-card">
					<div class="campaign-card__top">
						<span class="field__label">{$t('campaign.window')}</span>
						<span class="status-badge status-badge--draft">draft</span>
					</div>
					<div class="campaign-card__dates">
						{windowStart || '—'} — {windowEnd || '—'}
					</div>
				</div>

				{#if parameters.some(p => p.name)}
					<div class="campaign-card">
						<span class="field__label">{$t('campaign.parameters')} ({parameters.filter(p => p.name).length})</span>
						<div class="campaign-card__meta">
							{#each parameters.filter(p => p.name) as p}
								<span>{p.name} ({p.unit})</span>
							{/each}
						</div>
					</div>
				{/if}

				{#if regions.some(r => r.geoJson)}
					<div class="campaign-card">
						<span class="field__label">{$t('campaign.regions')} ({regions.filter(r => r.geoJson).length})</span>
					</div>
				{/if}

				{#if eligibility.some(e => e.deviceClass)}
					<div class="campaign-card">
						<span class="field__label">{$t('campaign.eligibility')} ({eligibility.filter(e => e.deviceClass).length})</span>
						<div class="campaign-card__meta">
							{#each eligibility.filter(e => e.deviceClass) as e}
								<span>{e.deviceClass} (tier {e.tier})</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>

			{#if error}
				<p class="form-error mt-4">{error}</p>
			{/if}
		</div>
	{/if}

	<!-- Navigation -->
	<div class="wizard__actions">
		{#if currentStep > 1}
			<button type="button" class="btn btn--secondary" onclick={prevStep}>
				{$t('campaign.back')}
			</button>
		{:else}
			<div></div>
		{/if}

		{#if currentStep < totalSteps}
			<button type="button" class="btn btn--primary" onclick={nextStep}>
				{$t('campaign.next')}
			</button>
		{:else}
			<button
				type="button"
				class="btn btn--primary"
				disabled={submitting}
				onclick={handleSubmit}
			>
				{submitting ? $t('common.loading') : $t('campaign.submit')}
			</button>
		{/if}
	</div>
</div>
