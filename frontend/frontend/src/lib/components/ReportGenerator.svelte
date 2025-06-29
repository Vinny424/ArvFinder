<script lang="ts">
	import { onMount } from 'svelte';
	
	export let propertyData: any = null;
	export let userTier: string = 'starter'; // starter, professional, enterprise
	
	let loading = false;
	let error = '';
	let reportPrice = 999; // $9.99 in cents
	let freeReports = false;
	let stripe: any = null;
	let elements: any = null;
	let cardElement: any = null;
	let paymentProcessing = false;
	
	// Stripe publishable key
	const STRIPE_PUBLISHABLE_KEY = 'pk_test_51Rf9L600n2nnxa7p51rzPLaOKBP8neatshbEpcfQuqmrORetSHfAidHUZeW9dtz0FL3NWbzLRDcmw4GrtGvUXKCL00HEIAOFka';
	
	onMount(async () => {
		// Check if user gets free reports
		freeReports = userTier === 'professional' || userTier === 'enterprise';
		
		if (!freeReports) {
			// Load Stripe for paid users
			await loadStripe();
		}
	});
	
	async function loadStripe() {
		if (typeof window !== 'undefined') {
			// Load Stripe.js
			const script = document.createElement('script');
			script.src = 'https://js.stripe.com/v3/';
			script.onload = initializeStripe;
			document.head.appendChild(script);
		}
	}
	
	function initializeStripe() {
		// @ts-ignore
		stripe = Stripe(STRIPE_PUBLISHABLE_KEY);
		elements = stripe.elements();
		
		// Create card element
		cardElement = elements.create('card', {
			style: {
				base: {
					fontSize: '16px',
					color: '#424770',
					'::placeholder': {
						color: '#aab7c4',
					},
				},
			},
		});
		
		// Mount card element
		const cardContainer = document.getElementById('card-element');
		if (cardContainer) {
			cardElement.mount('#card-element');
		}
	}
	
	async function generateReport() {
		if (!propertyData) {
			error = 'No property data available for report generation';
			return;
		}
		
		loading = true;
		error = '';
		
		try {
			if (freeReports) {
				// Generate report immediately for paid users
				await generateReportDirectly();
			} else {
				// Process payment first for starter users
				await processPaymentAndGenerateReport();
			}
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to generate report';
		} finally {
			loading = false;
			paymentProcessing = false;
		}
	}
	
	async function generateReportDirectly() {
		// For professional/enterprise users with free reports
		const response = await fetch('http://localhost:8080/api/v1/payments/create-report-payment', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				customer_email: 'user@example.com', // This would come from auth
				customer_name: 'Current User',
				property_id: propertyData.id || 'demo_property',
				user_tier: userTier,
			}),
		});
		
		const result = await response.json();
		
		if (result.success && result.data.free_report) {
			// Simulate report generation
			alert('Report generated successfully! (This would download a PDF in production)');
		} else {
			throw new Error('Failed to generate free report');
		}
	}
	
	async function processPaymentAndGenerateReport() {
		if (!stripe || !cardElement) {
			error = 'Payment system not loaded. Please refresh the page.';
			return;
		}
		
		paymentProcessing = true;
		
		// Create payment intent
		const response = await fetch('http://localhost:8080/api/v1/payments/create-report-payment', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				customer_email: 'user@example.com', // This would come from auth
				customer_name: 'Current User',
				property_id: propertyData.id || 'demo_property',
				user_tier: userTier,
			}),
		});
		
		const result = await response.json();
		
		if (!result.success) {
			throw new Error(result.error || 'Failed to create payment');
		}
		
		// Confirm payment
		const { error: stripeError } = await stripe.confirmCardPayment(result.data.client_secret, {
			payment_method: {
				card: cardElement,
				billing_details: {
					name: 'Current User',
					email: 'user@example.com',
				},
			},
		});
		
		if (stripeError) {
			throw new Error(stripeError.message);
		}
		
		// Payment successful, generate report
		alert('Payment successful! Report generated. (This would download a PDF in production)');
	}
</script>

<div class="bg-white rounded-lg shadow-lg p-6">
	<h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center">
		<svg class="w-6 h-6 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
			<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
		</svg>
		Generate Professional Report
	</h2>
	
	{#if !propertyData}
		<div class="bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded-lg mb-6">
			<p>Please complete an ARV calculation first to generate a report.</p>
		</div>
	{:else}
		<div class="mb-6">
			<h3 class="text-lg font-semibold text-gray-900 mb-2">Report Contents</h3>
			<ul class="space-y-2 text-gray-600">
				<li class="flex items-center">
					<svg class="w-4 h-4 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
					</svg>
					Comprehensive ARV analysis with detailed calculations
				</li>
				<li class="flex items-center">
					<svg class="w-4 h-4 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
					</svg>
					70% Rule assessment and investment recommendations
				</li>
				<li class="flex items-center">
					<svg class="w-4 h-4 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
					</svg>
					BRRRR strategy analysis and profit projections
				</li>
				<li class="flex items-center">
					<svg class="w-4 h-4 text-green-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
					</svg>
					Professional formatting suitable for lenders and partners
				</li>
			</ul>
		</div>
		
		{#if freeReports}
			<div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-lg mb-6">
				<p class="flex items-center">
					<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
					</svg>
					Report generation is included in your {userTier} subscription!
				</p>
			</div>
		{:else}
			<div class="bg-blue-50 border border-blue-200 text-blue-700 px-4 py-3 rounded-lg mb-6">
				<p class="font-medium">One-time payment: $9.99</p>
				<p class="text-sm mt-1">Upgrade to Professional ($29/month) for unlimited free reports!</p>
			</div>
			
			{#if !paymentProcessing}
				<div class="mb-6">
					<label class="block text-sm font-medium text-gray-700 mb-2">
						Payment Information
					</label>
					<div id="card-element" class="p-3 border border-gray-300 rounded-lg">
						<!-- Stripe Elements will create form elements here -->
					</div>
				</div>
			{/if}
		{/if}
		
		{#if error}
			<div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg mb-6">
				{error}
			</div>
		{/if}
		
		<button
			on:click={generateReport}
			disabled={loading || paymentProcessing}
			class="w-full bg-blue-600 text-white py-3 px-6 rounded-lg hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
		>
			{#if loading || paymentProcessing}
				<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white inline" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
					<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
					<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
				</svg>
				{paymentProcessing ? 'Processing Payment...' : 'Generating Report...'}
			{:else}
				{freeReports ? 'Generate Report (Free)' : 'Pay $9.99 & Generate Report'}
			{/if}
		</button>
		
		{#if !freeReports}
			<p class="text-xs text-gray-500 mt-3 text-center">
				Secure payment processing by Stripe. Your payment information is encrypted and secure.
			</p>
		{/if}
	{/if}
</div>