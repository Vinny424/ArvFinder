<script lang="ts">
	export let property: {
		id: string;
		address: string;
		price: number;
		arv: number;
		roi: number;
		image?: string;
	};
</script>

<div class="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-lg transition-shadow">
	<div class="aspect-video bg-gray-200 relative">
		{#if property.image}
			<img src={property.image} alt={property.address} class="w-full h-full object-cover" />
		{:else}
			<div class="w-full h-full flex items-center justify-center">
				<svg width="48" height="48" viewBox="0 0 48 48" class="text-gray-400">
					<path fill="currentColor" d="M24 4L6 20h4v20h12v-12h4v12h12V20h4L24 4zm0 6.5L36 22v16h-8v-12H20v12h-8V22l12-11.5z"/>
				</svg>
			</div>
		{/if}
		<div class="absolute top-2 right-2 bg-green-500 text-white px-2 py-1 rounded text-sm font-medium">
			{property.roi.toFixed(1)}% ROI
		</div>
	</div>
	
	<div class="p-4">
		<h3 class="font-semibold text-gray-900 mb-2">{property.address}</h3>
		
		<div class="grid grid-cols-2 gap-4 text-sm">
			<div>
				<span class="text-gray-600">Purchase Price</span>
				<div class="font-medium">${property.price.toLocaleString()}</div>
			</div>
			<div>
				<span class="text-gray-600">ARV</span>
				<div class="font-medium">${property.arv.toLocaleString()}</div>
			</div>
		</div>
		
		<button class="w-full mt-4 bg-blue-600 text-white py-2 rounded-lg hover:bg-blue-700 transition-colors">
			View Details
		</button>
	</div>
</div>