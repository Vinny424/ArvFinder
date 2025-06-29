<script lang="ts">
	import { apiService, type ArvCalculationResult } from '$lib/api';
	
	let purchasePrice = 0;
	let rehabCost = 0;
	let arv = 0;
	let holdingCosts = 0;
	let closingCosts = 0;
	let financingCosts = 0;
	let sellingCosts = 0;
	
	let calculationResult: ArvCalculationResult | null = null;
	let isCalculating = false;
	let error = '';
	
	// Local calculations for immediate feedback (while user types)
	$: maxOffer = arv * 0.7 - rehabCost;
	$: totalInvestment = purchasePrice + rehabCost + holdingCosts + closingCosts;
	$: potentialProfit = arv - totalInvestment;
	$: profitMargin = totalInvestment > 0 ? (potentialProfit / totalInvestment) * 100 : 0;
	$: is70RuleGood = purchasePrice <= maxOffer;
	
	async function calculateArv() {
		if (!purchasePrice || !arv) {
			error = 'Please enter purchase price and ARV';
			return;
		}
		
		isCalculating = true;
		error = '';
		
		try {
			calculationResult = await apiService.calculateArv({
				purchase_price: purchasePrice,
				rehab_cost: rehabCost,
				holding_costs: holdingCosts,
				closing_costs: closingCosts,
				arv: arv,
				financing_costs: financingCosts,
				selling_costs: sellingCosts
			});
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to calculate ARV';
			console.error('ARV calculation error:', err);
		} finally {
			isCalculating = false;
		}
	}
</script>

<div class="bg-white rounded-lg shadow-lg p-6">
	<h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
		<svg width="24" height="24" viewBox="0 0 24 24" class="mr-2 text-blue-600">
			<path fill="currentColor" d="M19 3H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2zm-5 14H7v-2h7v2zm3-4H7v-2h10v2zm0-4H7V7h10v2z"/>
		</svg>
		ARV Calculator
	</h2>
	
	<form on:submit|preventDefault={calculateArv} class="space-y-6">
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div>
				<label for="purchasePrice" class="block text-sm font-medium text-gray-700 mb-2">
					Purchase Price
				</label>
				<input
					id="purchasePrice"
					type="number"
					bind:value={purchasePrice}
					placeholder="150000"
					class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				/>
			</div>
			
			<div>
				<label for="rehabCost" class="block text-sm font-medium text-gray-700 mb-2">
					Rehab Cost
				</label>
				<input
					id="rehabCost"
					type="number"
					bind:value={rehabCost}
					placeholder="25000"
					class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				/>
			</div>
			
			<div>
				<label for="arv" class="block text-sm font-medium text-gray-700 mb-2">
					After Repair Value (ARV)
				</label>
				<input
					id="arv"
					type="number"
					bind:value={arv}
					placeholder="250000"
					class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				/>
			</div>
			
			<div>
				<label for="holdingCosts" class="block text-sm font-medium text-gray-700 mb-2">
					Holding Costs
				</label>
				<input
					id="holdingCosts"
					type="number"
					bind:value={holdingCosts}
					placeholder="5000"
					class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				/>
			</div>
		</div>
		
		<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
			<div>
				<label for="closingCosts" class="block text-sm font-medium text-gray-700 mb-2">
					Closing Costs
				</label>
				<input
					id="closingCosts"
					type="number"
					bind:value={closingCosts}
					placeholder="3000"
					class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				/>
			</div>
			
			<div>
				<label for="financingCosts" class="block text-sm font-medium text-gray-700 mb-2">
					Financing Costs (Optional)
				</label>
				<input
					id="financingCosts"
					type="number"
					bind:value={financingCosts}
					placeholder="2000"
					class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
				/>
			</div>
		</div>
		
		<div>
			<label for="sellingCosts" class="block text-sm font-medium text-gray-700 mb-2">
				Selling Costs (Optional)
			</label>
			<input
				id="sellingCosts"
				type="number"
				bind:value={sellingCosts}
				placeholder="15000"
				class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
			/>
		</div>
		
		{#if error}
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg">
				{error}
			</div>
		{/if}
		
		<button
			type="submit"
			disabled={isCalculating}
			class="w-full bg-blue-600 text-white py-3 rounded-lg hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
		>
			{isCalculating ? 'Calculating...' : 'Calculate ARV'}
		</button>
	</form>
	
	<!-- Live Preview (shows immediate calculations) -->
	{#if arv > 0 && !calculationResult}
		<div class="mt-8 p-6 bg-gray-50 rounded-lg">
			<h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Preview</h3>
			
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div class="bg-white p-4 rounded-lg">
					<div class="text-sm text-gray-600">70% Rule Max Offer</div>
					<div class="text-xl font-bold text-blue-600">${maxOffer.toLocaleString()}</div>
					<div class="text-xs {is70RuleGood ? 'text-green-600' : 'text-red-600'}">
						{is70RuleGood ? '✓ Meets 70% Rule' : '✗ Exceeds 70% Rule'}
					</div>
				</div>
				
				<div class="bg-white p-4 rounded-lg">
					<div class="text-sm text-gray-600">Potential Profit</div>
					<div class="text-xl font-bold {potentialProfit >= 0 ? 'text-green-600' : 'text-red-600'}">
						${potentialProfit.toLocaleString()}
					</div>
					<div class="text-xs text-gray-500">
						{profitMargin.toFixed(1)}% margin
					</div>
				</div>
			</div>
			
			<div class="mt-4 text-center">
				<p class="text-sm text-gray-600">Click "Calculate ARV" for comprehensive analysis</p>
			</div>
		</div>
	{/if}
	
	<!-- Detailed Results (from API) -->
	{#if calculationResult}
		<div class="mt-8 p-6 bg-gray-50 rounded-lg">
			<h3 class="text-lg font-semibold text-gray-900 mb-4">Comprehensive Analysis</h3>
			
			<!-- Key Metrics -->
			<div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
				<div class="bg-white p-4 rounded-lg">
					<div class="text-sm text-gray-600">70% Rule Max Offer</div>
					<div class="text-xl font-bold text-blue-600">${calculationResult.max_offer_70.toLocaleString()}</div>
					<div class="text-xs {calculationResult.is_70_rule_good ? 'text-green-600' : 'text-red-600'}">
						{calculationResult.is_70_rule_good ? '✓ Meets 70% Rule' : '✗ Exceeds 70% Rule'}
					</div>
				</div>
				
				<div class="bg-white p-4 rounded-lg">
					<div class="text-sm text-gray-600">Potential Profit</div>
					<div class="text-xl font-bold {calculationResult.potential_profit >= 0 ? 'text-green-600' : 'text-red-600'}">
						${calculationResult.potential_profit.toLocaleString()}
					</div>
					<div class="text-xs text-gray-500">
						{calculationResult.profit_margin.toFixed(1)}% margin
					</div>
				</div>
				
				<div class="bg-white p-4 rounded-lg">
					<div class="text-sm text-gray-600">ROI</div>
					<div class="text-xl font-bold {calculationResult.roi >= 0 ? 'text-green-600' : 'text-red-600'}">
						{calculationResult.roi.toFixed(1)}%
					</div>
					<div class="text-xs text-gray-500">
						Risk: {calculationResult.risk_level}
					</div>
				</div>
			</div>
			
			<!-- Investment Breakdown -->
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
				<div class="bg-white p-4 rounded-lg">
					<h4 class="font-semibold text-gray-900 mb-3">Investment Breakdown</h4>
					<div class="space-y-2 text-sm">
						<div class="flex justify-between">
							<span>Purchase Price:</span>
							<span>${calculationResult.purchase_price.toLocaleString()}</span>
						</div>
						<div class="flex justify-between">
							<span>Rehab Cost:</span>
							<span>${calculationResult.rehab_cost.toLocaleString()}</span>
						</div>
						<div class="flex justify-between">
							<span>Holding Costs:</span>
							<span>${calculationResult.holding_costs.toLocaleString()}</span>
						</div>
						<div class="flex justify-between">
							<span>Closing Costs:</span>
							<span>${calculationResult.closing_costs.toLocaleString()}</span>
						</div>
						{#if calculationResult.financing_costs > 0}
							<div class="flex justify-between">
								<span>Financing Costs:</span>
								<span>${calculationResult.financing_costs.toLocaleString()}</span>
							</div>
						{/if}
						<div class="border-t pt-2 flex justify-between font-semibold">
							<span>Total Investment:</span>
							<span>${calculationResult.total_investment.toLocaleString()}</span>
						</div>
					</div>
				</div>
				
				<div class="bg-white p-4 rounded-lg">
					<h4 class="font-semibold text-gray-900 mb-3">BRRRR Strategy</h4>
					<div class="space-y-2 text-sm">
						<div class="flex justify-between">
							<span>BRRRR Max Offer (75%):</span>
							<span>${calculationResult.brrrr_max_offer.toLocaleString()}</span>
						</div>
						<div class="flex justify-between">
							<span>BRRRR Profit:</span>
							<span class="{calculationResult.brrrr_profit >= 0 ? 'text-green-600' : 'text-red-600'}">${calculationResult.brrrr_profit.toLocaleString()}</span>
						</div>
					</div>
					<div class="mt-3 p-3 bg-blue-50 rounded text-xs text-blue-700">
						BRRRR strategy allows for 75% refinancing, potentially recovering more capital.
					</div>
				</div>
			</div>
			
			<!-- Recommendations -->
			{#if calculationResult.recommendations.length > 0}
				<div class="bg-white p-4 rounded-lg">
					<h4 class="font-semibold text-gray-900 mb-3">Recommendations</h4>
					<ul class="space-y-2">
						{#each calculationResult.recommendations as recommendation}
							<li class="flex items-start text-sm">
								<svg class="w-4 h-4 text-blue-500 mr-2 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
									<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
								</svg>
								{recommendation}
							</li>
						{/each}
					</ul>
				</div>
			{/if}
		</div>
	{/if}
</div>