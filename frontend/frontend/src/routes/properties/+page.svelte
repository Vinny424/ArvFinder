<script lang="ts">
	import { onMount } from 'svelte';
	import Header from '$lib/components/Header.svelte';
	import PropertyCard from '$lib/components/PropertyCard.svelte';
	import { apiService } from '$lib/api';
	
	let properties: any[] = [];
	let loading = true;
	let error = '';
	let showAddForm = false;
	
	// Summary statistics
	$: totalProperties = properties.length;
	$: totalValue = properties.reduce((sum, p) => sum + (p.arv || 0), 0);
	$: totalInvestment = properties.reduce((sum, p) => sum + (p.price || 0), 0);
	$: averageROI = totalProperties > 0 
		? properties.reduce((sum, p) => sum + (p.roi || 0), 0) / totalProperties 
		: 0;
	
	onMount(async () => {
		await loadProperties();
	});
	
	async function loadProperties() {
		try {
			loading = true;
			properties = await apiService.getProperties();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load properties';
		} finally {
			loading = false;
		}
	}
	
	function addProperty() {
		showAddForm = true;
	}
</script>

<svelte:head>
	<title>Properties - ArvFinder</title>
</svelte:head>

<Header />

<main class="min-h-screen bg-gray-50">
	<div class="container mx-auto px-4 py-8">
		<!-- Page Header -->
		<div class="flex justify-between items-center mb-8">
			<div>
				<h1 class="text-3xl font-bold text-gray-900">Property Portfolio</h1>
				<p class="text-gray-600 mt-2">Manage and analyze your investment properties</p>
			</div>
			<button
				on:click={addProperty}
				class="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors flex items-center"
			>
				<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
				</svg>
				Add Property
			</button>
		</div>
		
		<!-- Portfolio Summary -->
		<div class="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
			<div class="bg-white rounded-lg shadow p-6">
				<div class="flex items-center">
					<svg class="w-8 h-8 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
					</svg>
					<div class="ml-4">
						<p class="text-sm font-medium text-gray-600">Total Properties</p>
						<p class="text-2xl font-bold text-gray-900">{totalProperties}</p>
					</div>
				</div>
			</div>
			
			<div class="bg-white rounded-lg shadow p-6">
				<div class="flex items-center">
					<svg class="w-8 h-8 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"/>
					</svg>
					<div class="ml-4">
						<p class="text-sm font-medium text-gray-600">Total ARV</p>
						<p class="text-2xl font-bold text-gray-900">${totalValue.toLocaleString()}</p>
					</div>
				</div>
			</div>
			
			<div class="bg-white rounded-lg shadow p-6">
				<div class="flex items-center">
					<svg class="w-8 h-8 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
					</svg>
					<div class="ml-4">
						<p class="text-sm font-medium text-gray-600">Total Investment</p>
						<p class="text-2xl font-bold text-gray-900">${totalInvestment.toLocaleString()}</p>
					</div>
				</div>
			</div>
			
			<div class="bg-white rounded-lg shadow p-6">
				<div class="flex items-center">
					<svg class="w-8 h-8 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
					</svg>
					<div class="ml-4">
						<p class="text-sm font-medium text-gray-600">Average ROI</p>
						<p class="text-2xl font-bold text-gray-900">{averageROI.toFixed(1)}%</p>
					</div>
				</div>
			</div>
		</div>
		
		<!-- Loading State -->
		{#if loading}
			<div class="flex justify-center items-center py-12">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
			</div>
		{/if}
		
		<!-- Error State -->
		{#if error}
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
				{error}
				<button 
					on:click={loadProperties}
					class="ml-4 text-red-800 underline hover:text-red-900"
				>
					Retry
				</button>
			</div>
		{/if}
		
		<!-- Properties Grid -->
		{#if !loading && properties.length > 0}
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each properties as property (property.id)}
					<PropertyCard {property} />
				{/each}
			</div>
		{:else if !loading && !error}
			<!-- Empty State -->
			<div class="text-center py-12">
				<svg class="mx-auto h-24 w-24 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
				</svg>
				<h3 class="text-xl font-semibold text-gray-900 mt-4">No properties yet</h3>
				<p class="text-gray-600 mt-2 max-w-md mx-auto">
					Start building your portfolio by adding your first investment property.
				</p>
				<button
					on:click={addProperty}
					class="mt-6 bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors"
				>
					Add Your First Property
				</button>
			</div>
		{/if}
	</div>
</main>

<!-- Add Property Modal -->
{#if showAddForm}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
		<div class="bg-white rounded-lg max-w-md w-full p-6">
			<h3 class="text-lg font-semibold text-gray-900 mb-4">Add New Property</h3>
			<p class="text-gray-600 mb-4">
				Property management coming soon! For now, use the ARV Calculator to analyze potential investments.
			</p>
			<div class="flex justify-end space-x-3">
				<button
					on:click={() => showAddForm = false}
					class="px-4 py-2 text-gray-600 border border-gray-300 rounded-lg hover:bg-gray-50"
				>
					Close
				</button>
				<a
					href="/"
					class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
				>
					Go to Calculator
				</a>
			</div>
		</div>
	</div>
{/if}