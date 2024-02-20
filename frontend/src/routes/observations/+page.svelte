<script lang="ts">
	import Sidebar from '$lib/components/sidebar/sidebar.svelte';
	import { apiURL, toDateString } from '$lib/utils';
	import { Plus } from 'svelte-bootstrap-icons';
	async function fetchData() {
		const res = await fetch(apiURL + 'observations/observation', {
			headers: {
				Authorization: 'Bearer ' + localStorage.getItem('token')
				// 'Content-Type': 'application/x-www-form-urlencoded',
			}
		});
		const data = await res.json();

		if (res.ok) {
			return data;
		} else {
			throw new Error(data);
		}
	}

	async function search(searchedName: String) {
		const res = await fetch(apiURL + 'taxon/search?q=' + searchedName);
		const data = await res.json();

		if (res.ok) {
			return data;
		} else {
			throw new Error(data);
		}
	}

	async function getSpeciesName(ID: String) {
		const res = await fetch(apiURL + 'taxon/taxon?id=' + ID);
		const data = await res.json();

		if (res.ok) {
			return data;
		} else {
			throw new Error(data);
		}
	}

	async function submit() {
		const res = await fetch(apiURL + 'observations/observation', {
			method: createMode ? 'POST' : 'PATCH',
			body: createMode
				? JSON.stringify({
						taxonID,
						bloomStartDate: Math.round(new Date(bloomStartDate).getTime() / 1000),
						bloomPeakDate: Math.round(new Date(bloomPeakDate).getTime() / 1000),
						bloomEndDate: Math.round(new Date(bloomEndDate).getTime() / 1000),
						lat,
						long
					})
				: JSON.stringify({
						ID,
						taxonID,
						bloomStartDate: Math.round(new Date(bloomStartDate).getTime() / 1000),
						bloomPeakDate: Math.round(new Date(bloomPeakDate).getTime() / 1000),
						bloomEndDate: Math.round(new Date(bloomEndDate).getTime() / 1000),
						lat,
						long
					}),
			headers: {
				Authorization: 'Bearer ' + localStorage.getItem('token')
				// 'Content-Type': 'application/x-www-form-urlencoded',
			}
		});

		const json = await res.json();

		if (!res.ok) {
			error = true;
			errorMessage = json.error;
		} else {
			error = false;
			errorMessage = '';
			showObservation = false;
			fetchedData = fetchData();
		}
	}

	async function deleteItem() {
		const res = await fetch(apiURL + 'observations/observation?id=' + ID, {
			method: 'DELETE',
			headers: {
				Authorization: 'Bearer ' + localStorage.getItem('token')
				// 'Content-Type': 'application/x-www-form-urlencoded',
			}
		});

		const json = await res.json();

		if (!res.ok) {
			error = true;
			errorMessage = json.error;
		} else {
			error = false;
			errorMessage = '';
			showObservation = false;
			fetchedData = fetchData();
		}
	}

	let createMode = false;
	let showObservation = false;

	let ID = '';
	let taxonIDSearch = '';
	let taxonID = '';
	let bloomStartDate = '';
	let bloomPeakDate = '';
	let bloomEndDate = '';
	let lat = 0;
	let long = 0;

	let error = false;
	let errorMessage = '';

	let fetchedData = fetchData();
</script>

<Sidebar></Sidebar>
<div id="content">
	<div
		id="add-observation-button"
		on:click={() => {
			createMode = true;
			showObservation = true;
		}}
	>
		<Plus id="add-observation-icon-svg" />
	</div>
	{#await fetchedData then items}
		<div class="observation-item-container">
			{#each items.data as item}
				{#await getSpeciesName(item.TaxonID) then data}
					<div
						class="observation-item"
						onclick={() => {
							ID = item.id;
							taxonIDSearch = data.data.ScientificName;
							taxonID = item.TaxonID;
							bloomStartDate = toDateString(item.BloomStartDate);
							bloomPeakDate = toDateString(item.BloomPeakDate);
							bloomEndDate = toDateString(item.BloomEndDate);
							createMode = false;
							showObservation = true;
						}}
					>
						<div>
							{data.data.ScientificName}
						</div>
						<div class="observation-mod-date">
							{new Date(item.DateModified * 1000).toLocaleDateString('en-US')}
						</div>
					</div>
				{/await}
			{/each}
		</div>
	{:catch error}
		<p style="color: red">{error.message}</p>
	{/await}
</div>
{#if showObservation}
	<div id="observation-container-container">
		<div id="observation-container">
			<form>
				<div class="login-field-container">
					<input
						class="login-field"
						bind:value={taxonIDSearch}
						type="text"
						placeholder="Search taxonomy..."
					/>
					{#await search(taxonIDSearch)}
						<br /><select class="login-field" bind:value={taxonID}> </select>
					{:then data}
						<br /><select class="login-field" bind:value={taxonID}>
							{#each data.data as item}
								<option value={item.id}>{item.ScientificName}</option>
							{/each}
						</select>
					{:catch error}
						<br /><select class="login-field" bind:value={taxonID}> </select>
					{/await}
				</div>
				<div class="login-field-container">
					<label>Date of bloom start</label><br />
					<input
						class="login-field"
						bind:value={bloomStartDate}
						type="date"
						placeholder="Search taxonomy..."
					/>
				</div>
				<div class="login-field-container">
					<label>Date of bloom peak</label><br />
					<input
						class="login-field"
						bind:value={bloomPeakDate}
						type="date"
						placeholder="Search taxonomy..."
					/>
				</div>
				<div class="login-field-container">
					<label>Date of bloom end</label><br />
					<input
						class="login-field"
						bind:value={bloomEndDate}
						type="date"
						placeholder="Search taxonomy..."
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
				<div id="login-error-message-container" style={`padding: ${error ? '8px' : '0px'}`}>
					{errorMessage}
				</div>
				<div id="login-button-container">
					<button id="login-button" on:click={submit}>Submit</button>

					{#if !createMode}
						<button
							class="delete-button"
							on:click={() => {
								var answer = window.confirm('Are you sure you want to delete this observation?');
								if (answer) {
									deleteItem();
								}
							}}>Delete</button
						>
					{/if}
					<button
						class="warning-button"
						on:click={() => {
							showObservation = false;
						}}>Close</button
					>
				</div>
			</form>
		</div>
	</div>
{/if}
