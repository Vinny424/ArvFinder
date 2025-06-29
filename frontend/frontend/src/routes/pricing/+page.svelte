<script lang="ts">
	import { onMount } from 'svelte';
	import Header from '$lib/components/Header.svelte';
	
	interface SubscriptionPlan {
		name: string;
		price: number;
		price_id: string;
		features: string[];
		arv_limit: number;
		popular: boolean;
	}
	
	let plans: Record<string, SubscriptionPlan> = {};
	let loading = true;
	let error = '';
	
	// Stripe publishable key
	const STRIPE_PUBLISHABLE_KEY = 'pk_test_51Rf9L600n2nnxa7p51rzPLaOKBP8neatshbEpcfQuqmrORetSHfAidHUZeW9dtz0FL3NWbzLRDcmw4GrtGvUXKCL00HEIAOFka';
	
	onMount(async () => {
		await loadPlans();
		
		// Load Stripe script
		const script = document.createElement('script');
		script.src = 'https://js.stripe.com/v3/';
		document.head.appendChild(script);
	});
	
	async function loadPlans() {
		try {
			loading = true;
			const response = await fetch('http://localhost:8080/api/v1/payments/plans');
			const result = await response.json();
			
			if (result.success) {
				plans = result.data;
			} else {
				error = 'Failed to load pricing plans';
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load pricing plans';
		} finally {
			loading = false;
		}
	}
	
	async function selectPlan(planKey: string) {
		const plan = plans[planKey];
		
		if (planKey === 'starter') {
			// Free plan - just redirect to signup
			alert('Free plan selected! Sign up to get started.');
			return;
		}
		
		if (!plan.price_id) {
			alert('This plan is not yet available for purchase. Please contact support.');
			return;
		}
		
		// For now, show payment coming soon
		alert(`${plan.name} plan selected! Payment integration coming soon. This will cost $${plan.price / 100}/month.`);
		
		// TODO: Implement actual Stripe payment flow
		// 1. Create customer and subscription
		// 2. Redirect to Stripe checkout or show payment form
		// 3. Handle payment success/failure
	}
	
	function formatPrice(price: number): string {
		if (price === 0) return 'Free';
		return `$${price / 100}`;
	}
</script>

<svelte:head>
	<title>Pricing - ArvFinder</title>
</svelte:head>

<Header />

<main class="min-h-screen bg-gray-50">
	<div class="container mx-auto px-4 py-16">
		<!-- Header -->
		<div class="text-center mb-16">
			<h1 class="text-4xl font-bold text-gray-900 mb-4">
				Simple, Transparent Pricing
			</h1>
			<p class="text-xl text-gray-600 max-w-3xl mx-auto">
				Choose the plan that fits your investment goals. Start free and upgrade as you grow your portfolio.
			</p>
		</div>
		
		<!-- Loading State -->
		{#if loading}
			<div class="flex justify-center items-center py-12">
				<div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
			</div>
		{/if}
		
		<!-- Error State -->
		{#if error}
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6 max-w-md mx-auto">
				{error}
				<button 
					on:click={loadPlans}
					class="ml-4 text-red-800 underline hover:text-red-900"
				>
					Retry
				</button>
			</div>
		{/if}
		
		<!-- Pricing Cards -->
		{#if !loading && Object.keys(plans).length > 0}
			<div class="grid grid-cols-1 md:grid-cols-3 gap-8 max-w-6xl mx-auto">
				{#each Object.entries(plans) as [planKey, plan]}
					<div class="bg-white rounded-2xl shadow-lg overflow-hidden {plan.popular ? 'ring-2 ring-blue-500 relative' : ''}">
						{#if plan.popular}
							<div class="absolute top-0 left-1/2 transform -translate-x-1/2 bg-blue-500 text-white px-6 py-1 text-sm font-medium rounded-b-lg">
								Most Popular
							</div>
						{/if}
						
						<div class="p-8 {plan.popular ? 'pt-12' : ''}">
							<!-- Plan Header -->
							<div class="text-center mb-8">
								<h3 class="text-2xl font-bold text-gray-900 mb-2">{plan.name}</h3>
								<div class="mb-4">
									<span class="text-4xl font-bold text-gray-900">{formatPrice(plan.price)}</span>
									{#if plan.price > 0}
										<span class="text-gray-600">/month</span>
									{/if}
								</div>
								<p class="text-gray-600">
									{#if plan.arv_limit === -1}
										Unlimited ARV calculations
									{:else}
										{plan.arv_limit} ARV calculations per month
									{/if}
								</p>
							</div>
							
							<!-- Features -->
							<ul class="space-y-4 mb-8">
								{#each plan.features as feature}
									<li class="flex items-start">
										<svg class="w-5 h-5 text-green-500 mr-3 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
										</svg>
										<span class="text-gray-600">{feature}</span>
									</li>
								{/each}
							</ul>
							
							<!-- CTA Button -->
							<button
								on:click={() => selectPlan(planKey)}
								class="w-full py-3 px-6 rounded-lg font-medium transition-colors
									{plan.popular 
										? 'bg-blue-600 text-white hover:bg-blue-700' 
										: 'bg-gray-100 text-gray-900 hover:bg-gray-200'}"
							>
								{planKey === 'starter' ? 'Get Started Free' : `Choose ${plan.name}`}
							</button>
						</div>
					</div>
				{/each}
			</div>
		{/if}
		
		<!-- Additional Information -->
		<div class="mt-16 text-center">
			<h2 class="text-2xl font-bold text-gray-900 mb-8">Frequently Asked Questions</h2>
			
			<div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto text-left">
				<div class="bg-white rounded-lg p-6 shadow-sm">
					<h3 class="font-semibold text-gray-900 mb-3">Can I change plans anytime?</h3>
					<p class="text-gray-600">Yes, you can upgrade or downgrade your plan at any time. Changes take effect immediately and you'll be prorated accordingly.</p>
				</div>
				
				<div class="bg-white rounded-lg p-6 shadow-sm">
					<h3 class="font-semibold text-gray-900 mb-3">What payment methods do you accept?</h3>
					<p class="text-gray-600">We accept all major credit cards (Visa, MasterCard, American Express) and ACH bank transfers through our secure Stripe integration.</p>
				</div>
				
				<div class="bg-white rounded-lg p-6 shadow-sm">
					<h3 class="font-semibold text-gray-900 mb-3">Is there a free trial?</h3>
					<p class="text-gray-600">Yes! Our Starter plan is completely free and includes 10 ARV calculations per month. No credit card required.</p>
				</div>
				
				<div class="bg-white rounded-lg p-6 shadow-sm">
					<h3 class="font-semibold text-gray-900 mb-3">Do you offer annual discounts?</h3>
					<p class="text-gray-600">Yes, save 20% when you pay annually. Annual plans are available for Professional ($290/year) and Enterprise ($590/year).</p>
				</div>
			</div>
		</div>
		
		<!-- Contact Section -->
		<div class="mt-16 bg-blue-50 rounded-2xl p-8 text-center">
			<h3 class="text-2xl font-bold text-gray-900 mb-4">Need a Custom Solution?</h3>
			<p class="text-gray-600 mb-6">
				For teams of 10+ or custom integrations, we offer enterprise solutions tailored to your needs.
			</p>
			<button class="bg-blue-600 text-white px-8 py-3 rounded-lg hover:bg-blue-700 transition-colors">
				Contact Sales
			</button>
		</div>
	</div>
</main>