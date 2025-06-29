<script lang="ts">
	// Header component for ArvFinder
</script>

<header class="bg-blue-600 text-white shadow-lg">
	<div class="container mx-auto px-4 py-4 flex items-center justify-between">
		<div class="flex items-center space-x-2">
			<svg width="32" height="32" viewBox="0 0 32 32" class="text-white">
				<path fill="currentColor" d="M16 2L3 14h3v14h8v-8h4v8h8V14h3L16 2zm0 4.83L24 15v11h-4v-8H12v8H8V15l8-8.17z"/>
			</svg>
			<h1 class="text-xl font-bold">ArvFinder</h1>
		</div>
		
		<nav class="hidden md:flex space-x-6">
			<a href="/" class="hover:text-blue-200 transition-colors">Dashboard</a>
			<a href="/properties" class="hover:text-blue-200 transition-colors">Properties</a>
			<a href="/calculator" class="hover:text-blue-200 transition-colors">ARV Calculator</a>
			<a href="/reports" class="hover:text-blue-200 transition-colors">Reports</a>
			<a href="/pricing" class="hover:text-blue-200 transition-colors">Pricing</a>
		</nav>
		
		<div class="flex items-center space-x-4">
			<button class="bg-blue-700 hover:bg-blue-800 px-4 py-2 rounded-lg transition-colors">
				Account
			</button>
		</div>
	</div>
</header>

<style>
	/* Add Tailwind CSS classes via app.css */
</style>