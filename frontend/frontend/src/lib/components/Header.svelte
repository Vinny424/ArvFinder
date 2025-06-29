<script lang="ts">
	export let currentPage = '';
</script>

<header class="bg-gradient-to-r from-blue-600 to-blue-700 text-white shadow-xl">
	<div class="container mx-auto px-6 py-4">
		<div class="flex items-center justify-between">
			<!-- Logo Section -->
			<div class="flex items-center space-x-3">
				<div class="bg-white/10 p-2 rounded-lg backdrop-blur">
					<svg width="28" height="28" viewBox="0 0 32 32" class="text-white">
						<path fill="currentColor" d="M16 2L3 14h3v14h8v-8h4v8h8V14h3L16 2zm0 4.83L24 15v11h-4v-8H12v8H8V15l8-8.17z"/>
					</svg>
				</div>
				<div>
					<h1 class="text-2xl font-bold tracking-tight">ArvFinder</h1>
					<p class="text-xs text-blue-200 -mt-1">Property Investment Analysis</p>
				</div>
			</div>
			
			<!-- Navigation -->
			<nav class="hidden lg:flex items-center space-x-1">
				<a href="/" class="nav-link {currentPage === 'dashboard' ? 'active' : ''}">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"/>
					</svg>
					Dashboard
				</a>
				<a href="/properties" class="nav-link {currentPage === 'properties' ? 'active' : ''}">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5"/>
					</svg>
					Properties
				</a>
				<a href="/reports" class="nav-link {currentPage === 'reports' ? 'active' : ''}">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
					</svg>
					Reports
				</a>
				<a href="/pricing" class="nav-link {currentPage === 'pricing' ? 'active' : ''}">
					<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"/>
					</svg>
					Pricing
				</a>
			</nav>
			
			<!-- Actions -->
			<div class="flex items-center space-x-3">
				<div class="hidden md:flex items-center text-sm text-blue-200">
					<div class="w-2 h-2 bg-green-400 rounded-full mr-2"></div>
					Starter Plan
				</div>
				<button class="bg-white/10 hover:bg-white/20 px-4 py-2 rounded-lg transition-all duration-200 backdrop-blur border border-white/20">
					<svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
						<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
					</svg>
					Account
				</button>
			</div>
		</div>
	</div>
</header>

<style>
	.nav-link {
		@apply flex items-center space-x-2 px-3 py-2 rounded-lg text-blue-100 hover:text-white hover:bg-white/10 transition-all duration-200;
	}
	
	.nav-link.active {
		@apply bg-white/20 text-white;
	}
</style>