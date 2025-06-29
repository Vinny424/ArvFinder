<script lang="ts">
	import Header from '$lib/components/Header.svelte';
	import ArvCalculator from '$lib/components/ArvCalculator.svelte';
</script>

<svelte:head>
	<title>ARV Calculator - ArvFinder</title>
</svelte:head>

<Header />

<main class="min-h-screen bg-gray-50">
	<div class="container mx-auto px-4 py-8">
		<div class="text-center mb-8">
			<h1 class="text-3xl font-bold text-gray-900 mb-4">
				ARV Calculator
			</h1>
			<p class="text-xl text-gray-600 max-w-2xl mx-auto">
				Calculate After Repair Values, analyze investment opportunities, and make informed decisions with our comprehensive property analysis tools.
			</p>
		</div>
		
		<div class="max-w-4xl mx-auto">
			<ArvCalculator />
		</div>
		
		<!-- Additional Calculator Tools -->
		<div class="max-w-4xl mx-auto mt-12">
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Quick Calculators</h2>
			
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
				<div class="bg-white rounded-lg shadow p-6">
					<h3 class="text-lg font-semibold text-gray-900 mb-3">70% Rule</h3>
					<p class="text-gray-600 text-sm mb-4">
						Quick calculation to determine maximum offer based on ARV and rehab costs.
					</p>
					<button class="w-full bg-blue-100 text-blue-700 py-2 rounded-lg hover:bg-blue-200 transition-colors">
						Coming Soon
					</button>
				</div>
				
				<div class="bg-white rounded-lg shadow p-6">
					<h3 class="text-lg font-semibold text-gray-900 mb-3">Cash-on-Cash Return</h3>
					<p class="text-gray-600 text-sm mb-4">
						Calculate your annual return on actual cash invested in the property.
					</p>
					<button class="w-full bg-blue-100 text-blue-700 py-2 rounded-lg hover:bg-blue-200 transition-colors">
						Coming Soon
					</button>
				</div>
				
				<div class="bg-white rounded-lg shadow p-6">
					<h3 class="text-lg font-semibold text-gray-900 mb-3">Cap Rate</h3>
					<p class="text-gray-600 text-sm mb-4">
						Determine capitalization rate based on net operating income and property value.
					</p>
					<button class="w-full bg-blue-100 text-blue-700 py-2 rounded-lg hover:bg-blue-200 transition-colors">
						Coming Soon
					</button>
				</div>
			</div>
		</div>
		
		<!-- Educational Content -->
		<div class="max-w-4xl mx-auto mt-12 bg-blue-50 rounded-lg p-8">
			<h2 class="text-2xl font-bold text-gray-900 mb-4">Understanding ARV Calculations</h2>
			
			<div class="grid grid-cols-1 md:grid-cols-2 gap-8">
				<div>
					<h3 class="text-lg font-semibold text-gray-900 mb-3">The 70% Rule</h3>
					<p class="text-gray-700 mb-3">
						The 70% rule is a guideline that suggests investors should pay no more than 70% of a property's ARV minus repair costs.
					</p>
					<div class="bg-white p-4 rounded border-l-4 border-blue-500">
						<code class="text-sm text-blue-600">
							Max Offer = (ARV × 0.70) - Repair Costs
						</code>
					</div>
				</div>
				
				<div>
					<h3 class="text-lg font-semibold text-gray-900 mb-3">BRRRR Strategy</h3>
					<p class="text-gray-700 mb-3">
						Buy, Rehab, Rent, Refinance, Repeat - a strategy that allows you to recycle capital through refinancing.
					</p>
					<div class="bg-white p-4 rounded border-l-4 border-green-500">
						<code class="text-sm text-green-600">
							Refinance = ARV × 0.75
						</code>
					</div>
				</div>
			</div>
		</div>
	</div>
</main>