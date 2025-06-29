<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import ArvCalculator from '$lib/components/ArvCalculator.svelte';
	import PropertyCard from '$lib/components/PropertyCard.svelte';
	
	// Sample data for demonstration
	const sampleProperties = [
		{
			id: '1',
			address: '123 Main St, Denver, CO',
			price: 180000,
			arv: 250000,
			roi: 15.8
		},
		{
			id: '2', 
			address: '456 Oak Ave, Boulder, CO',
			price: 220000,
			arv: 300000,
			roi: 12.4
		}
	];
</script>

<svelte:head>
	<title>ArvFinder - Property Investment Analysis</title>
</svelte:head>

<Header />

<main class="min-h-screen bg-gray-50">
	<div class="container mx-auto px-4 py-8">
		<div class="text-center mb-12">
			<h1 class="text-4xl font-bold text-gray-900 mb-4">
				Professional ARV Analysis
			</h1>
			<p class="text-xl text-gray-600 max-w-2xl mx-auto">
				Calculate accurate After Repair Values and analyze investment opportunities with our advanced property analysis tools.
			</p>
		</div>
		
		<div class="mb-12">
			<ArvCalculator />
		</div>
		
		<section class="mb-12">
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Recent Properties</h2>
			<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
				{#each sampleProperties as property}
					<PropertyCard {property} />
				{/each}
			</div>
		</section>
	</div>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
	}
</style>
