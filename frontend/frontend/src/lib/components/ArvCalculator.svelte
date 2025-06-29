<script lang="ts">
	import { apiService, type ArvCalculationResult } from '$lib/api';
	
	let purchasePrice = 0;
	let rehabCost = 0;
	let arv = 0;
	let holdingCosts = 0;
	let closingCosts = 0;
	let financingCosts = 0;
	let sellingCosts = 0;
	let downPayment = 0;
	let loanType = 'conventional'; // 'conventional', 'rehab-loan', 'cash'
	
	let calculationResult: ArvCalculationResult | null = null;
	let isCalculating = false;
	let error = '';
	
	// Local calculations for immediate feedback (while user types)
	$: maxOffer = arv * 0.7 - rehabCost;
	$: totalInvestment = purchasePrice + rehabCost + holdingCosts + closingCosts;
	$: potentialProfit = arv - totalInvestment;
	$: profitMargin = totalInvestment > 0 ? (potentialProfit / totalInvestment) * 100 : 0;
	$: is70RuleGood = purchasePrice <= maxOffer;
	
	// Updated ROI calculation based on actual cash invested
	$: actualCashInvested = (() => {
		if (loanType === 'cash') {
			return totalInvestment; // All cash purchase
		} else if (loanType === 'rehab-loan') {
			return downPayment; // Down payment only (rehab included in loan)
		} else {
			return downPayment + rehabCost; // Down payment + rehab costs
		}
	})();
	$: actualROI = actualCashInvested > 0 ? (potentialProfit / actualCashInvested) * 100 : 0;
	
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
	
	function estimateCosts() {
		if (!purchasePrice || !arv) {
			error = 'Please enter purchase price and ARV first to estimate costs';
			return;
		}
		
		// Rough cost estimates based on industry standards
		const estimatedClosingCosts = Math.round(purchasePrice * 0.025); // 2.5% of purchase price
		const estimatedFinancingCosts = Math.round(purchasePrice * 0.015); // 1.5% for loan fees
		const estimatedSellingCosts = Math.round(arv * 0.08); // 8% of ARV (6% realtor + 2% other)
		
		// Update the values
		closingCosts = estimatedClosingCosts;
		financingCosts = estimatedFinancingCosts;
		sellingCosts = estimatedSellingCosts;
		
		// Clear any previous errors
		error = '';
	}
</script>

<div class="space-y-8">
	<div class="text-center">
		<h3 class="text-2xl font-bold text-gray-900 mb-2">ARV Investment Calculator</h3>
		<p class="text-gray-600">Enter property details for comprehensive investment analysis</p>
	</div>
	
	<form on:submit|preventDefault={calculateArv} class="space-y-8">
		<!-- Primary Investment Details -->
		<div class="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-2xl p-6 border border-blue-200">
			<h4 class="text-lg font-semibold text-gray-900 mb-4">Property Investment Details</h4>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div>
					<label for="purchasePrice" class="block text-sm font-medium text-gray-700 mb-2">
						Purchase Price
					</label>
					<div class="relative">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
						<input
							id="purchasePrice"
							type="number"
							bind:value={purchasePrice}
							placeholder="150,000"
							class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
						/>
					</div>
				</div>
				
				<div>
					<label for="arv" class="block text-sm font-medium text-gray-700 mb-2">
						After Repair Value (ARV)
					</label>
					<div class="relative">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
						<input
							id="arv"
							type="number"
							bind:value={arv}
							placeholder="250,000"
							class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
						/>
					</div>
				</div>
				
				<div>
					<label for="rehabCost" class="block text-sm font-medium text-gray-700 mb-2">
						Rehab Cost
					</label>
					<div class="relative">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
						<input
							id="rehabCost"
							type="number"
							bind:value={rehabCost}
							placeholder="25,000"
							class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
						/>
					</div>
				</div>
				
				<div>
					<label for="holdingCosts" class="block text-sm font-medium text-gray-700 mb-2">
						Holding Costs
					</label>
					<div class="relative">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
						<input
							id="holdingCosts"
							type="number"
							bind:value={holdingCosts}
							placeholder="5,000"
							class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
						/>
					</div>
				</div>
			</div>
		</div>
		
		<!-- Financing Details -->
		<div class="bg-gradient-to-r from-green-50 to-emerald-50 rounded-2xl p-6 border border-green-200">
			<h4 class="text-lg font-semibold text-gray-900 mb-4">Financing Details</h4>
			<div class="grid grid-cols-1 md:grid-cols-2 gap-6">
				<div>
					<label for="loanType" class="block text-sm font-medium text-gray-700 mb-2">
						Loan Type
					</label>
					<select
						id="loanType"
						bind:value={loanType}
						class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
					>
						<option value="conventional">Conventional (Purchase Only)</option>
						<option value="rehab-loan">Rehab Loan (Purchase + Repair)</option>
						<option value="cash">All Cash Purchase</option>
					</select>
					<p class="text-xs text-gray-500 mt-1">
						{#if loanType === 'conventional'}
							Down payment + rehab costs out of pocket
						{:else if loanType === 'rehab-loan'}
							Down payment only (rehab financed)
						{:else}
							No financing, all cash investment
						{/if}
					</p>
				</div>
				
				{#if loanType !== 'cash'}
					<div>
						<label for="downPayment" class="block text-sm font-medium text-gray-700 mb-2">
							Down Payment
						</label>
						<div class="relative">
							<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
							<input
								id="downPayment"
								type="number"
								bind:value={downPayment}
								placeholder={Math.round(purchasePrice * 0.2).toLocaleString()}
								class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
							/>
						</div>
						<p class="text-xs text-gray-500 mt-1">
							Suggested: 20% = ${Math.round(purchasePrice * 0.2).toLocaleString()}
						</p>
					</div>
				{/if}
			</div>
		</div>
		
		<!-- Additional Costs -->
		<div class="bg-gradient-to-r from-gray-50 to-gray-100 rounded-2xl p-6 border border-gray-200">
			<h4 class="text-lg font-semibold text-gray-900 mb-4">Additional Costs</h4>
			<div class="grid grid-cols-1 md:grid-cols-3 gap-6">
				<div>
					<label for="closingCosts" class="block text-sm font-medium text-gray-700 mb-2">
						Closing Costs
					</label>
					<div class="relative">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
						<input
							id="closingCosts"
							type="number"
							bind:value={closingCosts}
							placeholder="3,000"
							class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
						/>
					</div>
				</div>
				
				<div>
					<label for="financingCosts" class="block text-sm font-medium text-gray-700 mb-2">
						Financing Costs
					</label>
					<div class="relative">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
						<input
							id="financingCosts"
							type="number"
							bind:value={financingCosts}
							placeholder="2,000"
							class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
						/>
					</div>
				</div>
				
				<div>
					<label for="sellingCosts" class="block text-sm font-medium text-gray-700 mb-2">
						Selling Costs
					</label>
					<div class="relative">
						<span class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500">$</span>
						<input
							id="sellingCosts"
							type="number"
							bind:value={sellingCosts}
							placeholder="15,000"
							class="w-full pl-8 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-blue-500 focus:border-transparent bg-white"
						/>
					</div>
				</div>
			</div>
			
			<div class="mt-6 text-center">
				<button
					type="button"
					on:click={estimateCosts}
					class="bg-gradient-to-r from-amber-500 to-orange-500 hover:from-amber-600 hover:to-orange-600 text-white px-6 py-3 rounded-xl font-medium transition-all duration-200 shadow-lg transform hover:scale-105"
				>
					<div class="flex items-center justify-center">
						<svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
						</svg>
						Estimate Additional Costs
					</div>
				</button>
				<div class="mt-3 text-xs text-amber-700 bg-amber-50 rounded-lg p-3 border border-amber-200">
					<div class="flex items-start">
						<svg class="w-4 h-4 mr-2 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
							<path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
						</svg>
						<div>
							<p class="font-medium">Rough Estimates Only</p>
							<p>Closing: ~2.5% of purchase • Financing: ~1.5% of purchase • Selling: ~8% of ARV</p>
							<p>Always verify with professionals for accurate costs.</p>
						</div>
					</div>
				</div>
			</div>
		</div>
		
		{#if error}
			<div class="bg-red-50 border-l-4 border-red-400 p-4 rounded-lg">
				<div class="flex">
					<svg class="w-5 h-5 text-red-400 mr-2" fill="currentColor" viewBox="0 0 20 20">
						<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"/>
					</svg>
					<p class="text-red-700 font-medium">{error}</p>
				</div>
			</div>
		{/if}
		
		<div class="text-center">
			<button
				type="submit"
				disabled={isCalculating}
				class="bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white px-12 py-4 rounded-2xl font-semibold text-lg shadow-xl transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed transform hover:scale-105"
			>
				{#if isCalculating}
					<div class="flex items-center justify-center">
						<svg class="animate-spin w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
						</svg>
						Analyzing Property...
					</div>
				{:else}
					<div class="flex items-center justify-center">
						<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 7h6m0 10v-3m-3 3h.01M9 17h.01M9 14h.01M12 14h.01M15 11h.01M12 11h.01M9 11h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
						</svg>
						Run Complete Analysis
					</div>
				{/if}
			</button>
		</div>
	</form>
	
	<!-- Live Preview (shows immediate calculations) -->
	{#if arv > 0 && !calculationResult}
		<div class="bg-gradient-to-br from-amber-50 to-orange-50 rounded-2xl p-6 border border-amber-200">
			<div class="flex items-center mb-4">
				<div class="w-8 h-8 bg-gradient-to-br from-amber-500 to-orange-500 rounded-lg flex items-center justify-center mr-3">
					<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"/>
					</svg>
				</div>
				<h3 class="text-lg font-semibold text-gray-900">Quick Preview</h3>
			</div>
			
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
				<div class="bg-white/80 backdrop-blur p-4 rounded-xl border border-white/50">
					<div class="text-sm text-gray-600 mb-1">70% Rule Max Offer</div>
					<div class="text-xl font-bold text-blue-600">${maxOffer.toLocaleString()}</div>
					<div class="text-xs {is70RuleGood ? 'text-green-600' : 'text-red-600'} flex items-center mt-1">
						{#if is70RuleGood}
							<svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/>
							</svg>
							Meets 70% Rule
						{:else}
							<svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
								<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
							</svg>
							Exceeds 70% Rule
						{/if}
					</div>
				</div>
				
				<div class="bg-white/80 backdrop-blur p-4 rounded-xl border border-white/50">
					<div class="text-sm text-gray-600 mb-1">Cash-on-Cash ROI</div>
					<div class="text-xl font-bold {actualROI >= 0 ? 'text-green-600' : 'text-red-600'}">
						{actualROI.toFixed(1)}%
					</div>
					<div class="text-xs text-gray-500 mt-1">
						Cash invested: ${actualCashInvested.toLocaleString()}
					</div>
				</div>
			</div>
			
			<div class="text-center bg-white/60 rounded-xl p-3">
				<p class="text-sm text-gray-600 flex items-center justify-center">
					<svg class="w-4 h-4 mr-1 text-amber-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
					</svg>
					Run complete analysis for detailed investment recommendations
				</p>
			</div>
		</div>
	{/if}
	
	<!-- Detailed Results (from API) -->
	{#if calculationResult}
		<div class="bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl p-6 border border-green-200">
			<div class="flex items-center mb-6">
				<div class="w-8 h-8 bg-gradient-to-br from-green-500 to-emerald-500 rounded-lg flex items-center justify-center mr-3">
					<svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
					</svg>
				</div>
				<h3 class="text-xl font-bold text-gray-900">Investment Analysis Complete</h3>
			</div>
			
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