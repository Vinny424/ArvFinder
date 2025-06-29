<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import ArvCalculator from '$lib/components/ArvCalculator.svelte';
	
	let activeTab = 'arv-calculator';
	let searchAddress = '';
	let isSearching = false;
	let searchResults: any = null;
	
	const tabs = [
		{
			id: 'arv-calculator',
			name: 'ARV Calculator',
			icon: 'calculator',
			description: 'Calculate After Repair Value'
		},
		{
			id: 'seventy-rule',
			name: '70% Rule',
			icon: 'percentage',
			description: 'Quick investment screening'
		},
		{
			id: 'brrrr-analysis',
			name: 'BRRRR Analysis',
			icon: 'refresh',
			description: 'Buy, Rehab, Rent, Refinance, Repeat'
		},
		{
			id: 'cash-flow',
			name: 'Cash Flow',
			icon: 'trending-up',
			description: 'Monthly cash flow calculator'
		}
	];
	
	async function searchProperty() {
		if (!searchAddress.trim()) return;
		
		isSearching = true;
		searchResults = null;
		
		try {
			// This would integrate with a property data API like RentSpree, Zillow, etc.
			// For now, we'll simulate the API call
			await new Promise(resolve => setTimeout(resolve, 2000));
			
			// Simulated API response
			searchResults = {
				address: searchAddress,
				price: 225000,
				bedrooms: 3,
				bathrooms: 2,
				sqft: 1200,
				yearBuilt: 1985,
				zestimate: 280000,
				rentEstimate: 1800,
				neighborhood: 'Downtown',
				comps: [
					{ address: '789 Pine St', price: 215000, sqft: 1150 },
					{ address: '321 Elm Rd', price: 235000, sqft: 1280 },
					{ address: '654 Birch Ave', price: 220000, sqft: 1200 }
				]
			};
		} catch (error) {
			console.error('Error searching property:', error);
		} finally {
			isSearching = false;
		}
	}
	
	function handleKeyPress(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			searchProperty();
		}
	}
</script>

<svelte:head>
	<title>ArvFinder - Property Investment Analysis</title>
</svelte:head>

<Header currentPage="dashboard" />

<main class="min-h-screen bg-gradient-to-br from-gray-50 to-blue-50">
	<!-- Hero Section with Address Search -->
	<section class="relative overflow-hidden">
		<div class="absolute inset-0 bg-gradient-to-r from-blue-600/5 to-purple-600/5"></div>
		<div class="relative container mx-auto px-6 py-16">
			<div class="text-center mb-12">
				<h1 class="text-5xl font-bold text-gray-900 mb-6 tracking-tight">
					Professional Property Analysis
				</h1>
				<p class="text-xl text-gray-600 max-w-3xl mx-auto mb-12 leading-relaxed">
					Enter any property address to get instant ARV calculations, investment analysis, and comprehensive market data.
				</p>
				
				<!-- Address Search Bar -->
				<div class="max-w-2xl mx-auto">
					<div class="relative">
						<input 
							type="text" 
							bind:value={searchAddress}
							on:keypress={handleKeyPress}
							placeholder="Enter property address (e.g., 123 Main St, Denver, CO)"
							class="w-full px-6 py-4 text-lg border-2 border-gray-200 rounded-2xl shadow-lg focus:border-blue-500 focus:ring-4 focus:ring-blue-500/20 transition-all duration-200 bg-white/90 backdrop-blur"
						/>
						<button 
							on:click={searchProperty}
							disabled={isSearching || !searchAddress.trim()}
							class="absolute right-3 top-1/2 -translate-y-1/2 bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white px-6 py-2 rounded-xl transition-all duration-200 font-medium shadow-lg"
						>
							{#if isSearching}
								<svg class="animate-spin w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
								</svg>
							{:else}
								<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
								</svg>
							{/if}
						</button>
					</div>
				</div>
			</div>
			
			<!-- Search Results -->
			{#if searchResults}
				<div class="max-w-4xl mx-auto mb-12">
					<div class="bg-white/80 backdrop-blur rounded-3xl shadow-xl border border-white/50 p-8">
						<div class="flex items-start justify-between mb-6">
							<div>
								<h3 class="text-2xl font-bold text-gray-900 mb-2">{searchResults.address}</h3>
								<div class="flex items-center space-x-6 text-gray-600">
									<span class="flex items-center">
										<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"/>
										</svg>
										{searchResults.bedrooms} bed
									</span>
									<span class="flex items-center">
										<svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16"/>
										</svg>
										{searchResults.bathrooms} bath
									</span>
									<span>{searchResults.sqft.toLocaleString()} sqft</span>
									<span>Built {searchResults.yearBuilt}</span>
								</div>
							</div>
							<div class="text-right">
								<div class="text-3xl font-bold text-gray-900">${searchResults.price.toLocaleString()}</div>
								<div class="text-sm text-gray-500">List Price</div>
							</div>
						</div>
						
						<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
							<div class="bg-gradient-to-br from-blue-50 to-blue-100 rounded-2xl p-6 border border-blue-200">
								<div class="text-2xl font-bold text-blue-900">${searchResults.zestimate.toLocaleString()}</div>
								<div class="text-sm text-blue-700">Estimated ARV</div>
							</div>
							<div class="bg-gradient-to-br from-green-50 to-green-100 rounded-2xl p-6 border border-green-200">
								<div class="text-2xl font-bold text-green-900">${searchResults.rentEstimate.toLocaleString()}</div>
								<div class="text-sm text-green-700">Monthly Rent Estimate</div>
							</div>
							<div class="bg-gradient-to-br from-purple-50 to-purple-100 rounded-2xl p-6 border border-purple-200">
								<div class="text-2xl font-bold text-purple-900">{searchResults.neighborhood}</div>
								<div class="text-sm text-purple-700">Neighborhood</div>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</section>

	<!-- Tab Navigation -->
	<section class="container mx-auto px-6 py-8">
		<div class="bg-white/80 backdrop-blur rounded-3xl shadow-xl border border-white/50 overflow-hidden">
			<!-- Tab Headers -->
			<div class="border-b border-gray-200 bg-gradient-to-r from-gray-50 to-white">
				<nav class="flex overflow-x-auto scrollbar-hide">
					{#each tabs as tab}
						<button
							on:click={() => activeTab = tab.id}
							class="flex-shrink-0 px-8 py-6 text-left border-b-2 transition-all duration-200 hover:bg-white/50 {activeTab === tab.id ? 'border-blue-500 bg-white text-blue-600' : 'border-transparent text-gray-600 hover:text-gray-900'}"
						>
							<div class="flex items-center space-x-3">
								<div class="w-8 h-8 rounded-lg bg-gradient-to-br {activeTab === tab.id ? 'from-blue-500 to-blue-600' : 'from-gray-400 to-gray-500'} flex items-center justify-center text-white">
									{#if tab.icon === 'calculator'}
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
										</svg>
									{:else if tab.icon === 'percentage'}
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
										</svg>
									{:else if tab.icon === 'refresh'}
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
										</svg>
									{:else if tab.icon === 'trending-up'}
										<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
										</svg>
									{/if}
								</div>
								<div>
									<div class="font-semibold text-sm">{tab.name}</div>
									<div class="text-xs text-gray-500">{tab.description}</div>
								</div>
							</div>
						</button>
					{/each}
				</nav>
			</div>

			<!-- Tab Content -->
			<div class="p-8">
				{#if activeTab === 'arv-calculator'}
					<ArvCalculator />
				{:else if activeTab === 'seventy-rule'}
					<div class="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-2xl p-8 border border-blue-200">
						<h3 class="text-2xl font-bold text-gray-900 mb-6">70% Rule Calculator</h3>
						<p class="text-gray-600 mb-6">Quick screening tool: Maximum purchase price should not exceed 70% of ARV minus repair costs.</p>
						<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">After Repair Value (ARV)</label>
								<input type="number" placeholder="250000" class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">Estimated Repair Costs</label>
								<input type="number" placeholder="25000" class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
							</div>
							<div>
								<label class="block text-sm font-medium text-gray-700 mb-2">Maximum Purchase Price</label>
								<input type="number" readonly placeholder="150000" class="w-full px-4 py-3 bg-gray-100 border border-gray-300 rounded-xl" />
							</div>
						</div>
					</div>
				{:else if activeTab === 'brrrr-analysis'}
					<div class="bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl p-8 border border-green-200">
						<h3 class="text-2xl font-bold text-gray-900 mb-6">BRRRR Strategy Analysis</h3>
						<p class="text-gray-600 mb-6">Analyze the complete Buy, Rehab, Rent, Refinance, Repeat investment strategy.</p>
						<div class="text-center text-gray-500 py-12">
							<svg class="w-16 h-16 mx-auto mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z"/>
							</svg>
							<p>BRRRR calculator coming soon...</p>
						</div>
					</div>
				{:else if activeTab === 'cash-flow'}
					<div class="bg-gradient-to-br from-purple-50 to-pink-50 rounded-2xl p-8 border border-purple-200">
						<h3 class="text-2xl font-bold text-gray-900 mb-6">Cash Flow Calculator</h3>
						<p class="text-gray-600 mb-6">Calculate monthly cash flow and return on investment for rental properties.</p>
						<div class="text-center text-gray-500 py-12">
							<svg class="w-16 h-16 mx-auto mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"/>
							</svg>
							<p>Cash flow calculator coming soon...</p>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</section>
</main>

<style>
	:global(body) {
		margin: 0;
		font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
	}
	
	.scrollbar-hide {
		-ms-overflow-style: none;
		scrollbar-width: none;
	}
	
	.scrollbar-hide::-webkit-scrollbar {
		display: none;
	}
</style>
