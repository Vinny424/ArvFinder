// API configuration and service functions
const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// API response types
export interface ArvCalculationRequest {
	purchase_price: number;
	rehab_cost: number;
	holding_costs: number;
	closing_costs: number;
	arv: number;
	financing_costs?: number;
	selling_costs?: number;
}

export interface ArvCalculationResult {
	purchase_price: number;
	rehab_cost: number;
	holding_costs: number;
	closing_costs: number;
	arv: number;
	financing_costs: number;
	selling_costs: number;
	max_offer_70: number;
	is_70_rule_good: boolean;
	total_investment: number;
	potential_profit: number;
	profit_margin: number;
	roi: number;
	brrrr_max_offer: number;
	brrrr_profit: number;
	risk_level: string;
	recommendations: string[];
}

export interface ApiResponse<T> {
	success: boolean;
	data: T;
	error?: string;
}

// API service class
class ApiService {
	private async fetchApi<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
		const url = `${API_BASE_URL}${endpoint}`;
		
		const response = await fetch(url, {
			headers: {
				'Content-Type': 'application/json',
				...options.headers,
			},
			...options,
		});

		if (!response.ok) {
			throw new Error(`API error: ${response.status} ${response.statusText}`);
		}

		const result = await response.json();
		
		if (!result.success) {
			throw new Error(result.error || 'API request failed');
		}

		return result.data;
	}

	// ARV Calculation endpoints
	async calculateArv(request: ArvCalculationRequest): Promise<ArvCalculationResult> {
		return this.fetchApi<ArvCalculationResult>('/arv/calculate', {
			method: 'POST',
			body: JSON.stringify(request),
		});
	}

	async calculate70Rule(arv: number, rehabCost: number): Promise<{
		arv: number;
		rehab_cost: number;
		max_offer: number;
		rule: string;
	}> {
		return this.fetchApi('/arv/70-rule', {
			method: 'POST',
			body: JSON.stringify({
				arv,
				rehab_cost: rehabCost,
			}),
		});
	}

	async calculateROI(profit: number, investment: number): Promise<{
		profit: number;
		investment: number;
		roi: number;
		roi_formatted: number;
	}> {
		return this.fetchApi('/arv/roi', {
			method: 'POST',
			body: JSON.stringify({
				profit,
				investment,
			}),
		});
	}

	async calculateCashOnCash(annualCashFlow: number, totalCashInvested: number): Promise<{
		annual_cash_flow: number;
		total_cash_invested: number;
		cash_on_cash_return: number;
	}> {
		return this.fetchApi('/arv/cash-on-cash', {
			method: 'POST',
			body: JSON.stringify({
				annual_cash_flow: annualCashFlow,
				total_cash_invested: totalCashInvested,
			}),
		});
	}

	async calculateCapRate(netOperatingIncome: number, propertyValue: number): Promise<{
		net_operating_income: number;
		property_value: number;
		cap_rate: number;
	}> {
		return this.fetchApi('/arv/cap-rate', {
			method: 'POST',
			body: JSON.stringify({
				net_operating_income: netOperatingIncome,
				property_value: propertyValue,
			}),
		});
	}

	// Property endpoints
	async getProperties(): Promise<any[]> {
		return this.fetchApi<any[]>('/properties/');
	}

	// Health check
	async healthCheck(): Promise<{ status: string; service: string }> {
		const response = await fetch(`${API_BASE_URL.replace('/api/v1', '')}/health`);
		return response.json();
	}
}

// Export singleton instance
export const apiService = new ApiService();