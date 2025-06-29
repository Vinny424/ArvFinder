<script lang="ts">
	import Header from '$lib/components/Header.svelte';
</script>

<svelte:head>
	<title>Reports - ArvFinder</title>
</svelte:head>

<Header />

<main class="min-h-screen bg-gray-50">
	<div class="container mx-auto px-4 py-8">
		<div class="text-center mb-8">
			<h1 class="text-3xl font-bold text-gray-900 mb-4">
				Investment Reports
			</h1>
			<p class="text-xl text-gray-600 max-w-2xl mx-auto">
				Generate professional reports for your property analysis and share with partners, lenders, or team members.
			</p>
		</div>
		
		<!-- Report Types -->
		<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-12">
			<div class="bg-white rounded-lg shadow-lg p-6">
				<div class="flex items-center mb-4">
					<svg class="w-8 h-8 text-blue-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
					</svg>
					<h3 class="text-xl font-semibold text-gray-900">ARV Analysis Report</h3>
				</div>
				<p class="text-gray-600 mb-4">
					Comprehensive property analysis including ARV calculations, 70% rule assessment, and investment recommendations.
				</p>
				<button class="w-full bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition-colors">
					Generate Report
				</button>
			</div>
			
			<div class="bg-white rounded-lg shadow-lg p-6">
				<div class="flex items-center mb-4">
					<svg class="w-8 h-8 text-green-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"/>
					</svg>
					<h3 class="text-xl font-semibold text-gray-900">Portfolio Summary</h3>
				</div>
				<p class="text-gray-600 mb-4">
					Overview of your entire property portfolio with key metrics, performance data, and growth trends.
				</p>
				<button class="w-full bg-green-600 text-white py-2 rounded-lg hover:bg-green-700 transition-colors">
					Generate Report
				</button>
			</div>
			
			<div class="bg-white rounded-lg shadow-lg p-6">
				<div class="flex items-center mb-4">
					<svg class="w-8 h-8 text-purple-500 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
					</svg>
					<h3 class="text-xl font-semibold text-gray-900">Market Comparison</h3>
				</div>
				<p class="text-gray-600 mb-4">
					Compare your properties against market trends and similar investments in your area.
				</p>
				<button class="w-full bg-purple-600 text-white py-2 rounded-lg hover:bg-purple-700 transition-colors">
					Generate Report
				</button>
			</div>
		</div>
		
		<!-- Recent Reports -->
		<div class="bg-white rounded-lg shadow-lg p-6">
			<h2 class="text-2xl font-bold text-gray-900 mb-6">Recent Reports</h2>
			
			<div class="text-center py-12">
				<svg class="mx-auto h-24 w-24 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
				</svg>
				<h3 class="text-xl font-semibold text-gray-900 mt-4">No reports generated yet</h3>
				<p class="text-gray-600 mt-2 max-w-md mx-auto">
					Generate your first report to start tracking your investment analysis and performance.
				</p>
			</div>
		</div>
		
		<!-- Report Features -->
		<div class="mt-12 bg-blue-50 rounded-lg p-8">
			<h2 class="text-2xl font-bold text-gray-900 mb-6 text-center">Report Features</h2>
			
			<div class="grid grid-cols-1 md:grid-cols-2 gap-8">
				<div class="flex items-start">
					<svg class="w-6 h-6 text-blue-500 mr-3 mt-1 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
					</svg>
					<div>
						<h3 class="font-semibold text-gray-900 mb-2">Professional Formatting</h3>
						<p class="text-gray-700">Clean, professional PDF reports suitable for sharing with lenders, partners, and investors.</p>
					</div>
				</div>
				
				<div class="flex items-start">
					<svg class="w-6 h-6 text-blue-500 mr-3 mt-1 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"/>
					</svg>
					<div>
						<h3 class="font-semibold text-gray-900 mb-2">Detailed Analytics</h3>
						<p class="text-gray-700">Comprehensive analysis including comparable properties, market trends, and risk assessment.</p>
					</div>
				</div>
				
				<div class="flex items-start">
					<svg class="w-6 h-6 text-blue-500 mr-3 mt-1 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 4V2a1 1 0 011-1h4a1 1 0 011 1v2h4a1 1 0 011 1v2a1 1 0 01-1 1h-1v9a2 2 0 01-2 2H7a2 2 0 01-2-2V8H4a1 1 0 01-1-1V5a1 1 0 011-1h3z"/>
					</svg>
					<div>
						<h3 class="font-semibold text-gray-900 mb-2">Custom Branding</h3>
						<p class="text-gray-700">Add your company logo and branding to create professional reports for your business.</p>
					</div>
				</div>
				
				<div class="flex items-start">
					<svg class="w-6 h-6 text-blue-500 mr-3 mt-1 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.367 2.684 3 3 0 00-5.367-2.684z"/>
					</svg>
					<div>
						<h3 class="font-semibold text-gray-900 mb-2">Easy Sharing</h3>
						<p class="text-gray-700">Share reports via email, download as PDF, or generate shareable links for stakeholders.</p>
					</div>
				</div>
			</div>
		</div>
	</div>
</main>