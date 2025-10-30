package service

import (
	"github.com/Bilal-Cplusoft/sunready/internal/repo"
	"math"
	"fmt"
)

type QuoteService struct {
	quoteRepo *repo.QuoteRepo
}

type QuoteInput struct {
	SystemSizeKW          float64
	AnnualProductionKWh   float64
	MonthlyElectricBill   float64
	ElectricalOffsetPct   float64
	PanelCount            int
	State                 string
	CostPerWatt           *float64
	UtilityRatePerKWh     *float64
	AnnualUtilityIncrease *float64
	FederalTaxCredit      *float64
	LoanInterestRate      *float64
	LoanTermYears         *int
}

type QuoteResult struct {
	SystemCostBeforeIncentives float64 `json:"system_cost_before_incentives"`
	FederalTaxCredit           float64 `json:"federal_tax_credit"`
	SystemCostAfterIncentives  float64 `json:"system_cost_after_incentives"`
	EstimatedMonthlyPayment    float64 `json:"estimated_monthly_payment"`
	CurrentMonthlyBill         float64 `json:"current_monthly_bill"`
	EstimatedNewMonthlyBill    float64 `json:"estimated_new_monthly_bill"`
	MonthlySavings             float64 `json:"monthly_savings"`
	FirstYearSavings           float64 `json:"first_year_savings"`
	TwentyFiveYearSavings      float64 `json:"twenty_five_year_savings"`
	SystemSizeKW               float64 `json:"system_size_kw"`
	AnnualProductionKWh        float64 `json:"annual_production_kwh"`
	PanelCount                 int     `json:"panel_count"`
	ElectricalOffset           float64 `json:"electrical_offset_pct"`
	CostPerWatt                float64 `json:"cost_per_watt"`
	SimplePaybackYears         float64 `json:"simple_payback_years"`
	BreakEvenYear              int     `json:"break_even_year"`
	Summary                    string  `json:"summary"`
}


func NewQuoteService(quoteRepo *repo.QuoteRepo) *QuoteService {
	return &QuoteService{quoteRepo: quoteRepo}
}


func (s *QuoteService) CalculateQuote(input QuoteInput) (*QuoteResult, error) {
	if input.SystemSizeKW <= 0 {
		return nil, fmt.Errorf("system size must be greater than 0")
	}
	if input.MonthlyElectricBill <= 0 {
		return nil, fmt.Errorf("monthly electric bill must be greater than 0")
	}
	costPerWatt := 3.00
	utilityRate := 0.13
	annualIncrease := 0.03
	taxCredit := 0.30
	interestRate := 0.0699
	loanTermYears := 25
	sunHoursPerDay := 5.0
	electricalOffsetPct := 95.0
	annualProductionKWh := input.SystemSizeKW * sunHoursPerDay * 365 * 0.75
	panelCount := input.PanelCount
	systemSizeWatts := input.SystemSizeKW * 1000
	systemCostBeforeIncentives := systemSizeWatts * costPerWatt
	federalTaxCreditAmount := systemCostBeforeIncentives * taxCredit
	systemCostAfterIncentives := systemCostBeforeIncentives - federalTaxCreditAmount
	monthlyRate := interestRate / 12
	numPayments := float64(loanTermYears * 12)
	monthlyPayment := 0.0
	if interestRate > 0 {
		monthlyPayment = systemCostBeforeIncentives *
			(monthlyRate * math.Pow(1+monthlyRate, numPayments)) /
			(math.Pow(1+monthlyRate, numPayments) - 1)
	} else {
		monthlyPayment = systemCostBeforeIncentives / numPayments
	}

	annualCurrentBill := input.MonthlyElectricBill * 12
	offsetRatio := electricalOffsetPct / 100.0
	annualSolarSavings := annualProductionKWh * utilityRate
	remainingUsagePct := math.Max(0, 1.0-offsetRatio)
	newMonthlyBill := input.MonthlyElectricBill * remainingUsagePct
	monthlySavingsFromSolar := input.MonthlyElectricBill - newMonthlyBill
	netMonthlySavings := monthlySavingsFromSolar - monthlyPayment
	annualSavingsFromReducedBill := (input.MonthlyElectricBill - newMonthlyBill) * 12
	firstYearSavings := annualSavingsFromReducedBill - (monthlyPayment * 12)


	totalUtilityCostWithoutSolar := 0.0
	totalCostWithSolar := 0.0
	for year := 1; year <= 25; year++ {
		yearlyBill := annualCurrentBill * math.Pow(1+annualIncrease, float64(year-1))
		totalUtilityCostWithoutSolar += yearlyBill

		yearlyNewBill := newMonthlyBill * 12 * math.Pow(1+annualIncrease, float64(year-1))
		if year <= loanTermYears {
			totalCostWithSolar += (monthlyPayment * 12) + yearlyNewBill
		} else {
			totalCostWithSolar += yearlyNewBill
		}
	}
	twentyFiveYearSavings := totalUtilityCostWithoutSolar - totalCostWithSolar
	simplePayback := systemCostAfterIncentives / annualSolarSavings
	breakEvenYear := 0
	cumulativeSavings := 0.0
	for year := 1; year <= 30; year++ {
		yearlyBill := annualCurrentBill * math.Pow(1+annualIncrease, float64(year-1))
		yearlySolarSavings := yearlyBill * offsetRatio

		if year <= loanTermYears {
			cumulativeSavings += yearlySolarSavings - (monthlyPayment * 12)
		} else {
			cumulativeSavings += yearlySolarSavings
		}

		if cumulativeSavings >= systemCostBeforeIncentives && breakEvenYear == 0 {
			breakEvenYear = year
			break
		}
	}

	summary := fmt.Sprintf(
		"This %.2f kW solar system with %d panels will produce approximately %.0f kWh annually, "+
			"offsetting %.0f%% of your electricity usage. "+
			"The system costs $%.2f before incentives ($%.2f after federal tax credit). "+
			"Your estimated monthly payment is $%.2f, and you'll save approximately $%.2f in the first year. "+
			"Over 25 years, your total savings are estimated at $%.2f.",
		input.SystemSizeKW,
		panelCount,
		annualProductionKWh,
		electricalOffsetPct,
		systemCostBeforeIncentives,
		systemCostAfterIncentives,
		monthlyPayment,
		firstYearSavings,
		twentyFiveYearSavings,
	)

	return &QuoteResult{
		SystemCostBeforeIncentives: math.Round(systemCostBeforeIncentives*100) / 100,
		FederalTaxCredit:           math.Round(federalTaxCreditAmount*100) / 100,
		SystemCostAfterIncentives:  math.Round(systemCostAfterIncentives*100) / 100,
		EstimatedMonthlyPayment:    math.Round(monthlyPayment*100) / 100,
		CurrentMonthlyBill:         input.MonthlyElectricBill,
		EstimatedNewMonthlyBill:    math.Round(newMonthlyBill*100) / 100,
		MonthlySavings:             math.Round(netMonthlySavings*100) / 100,
		FirstYearSavings:           math.Round(firstYearSavings*100) / 100,
		TwentyFiveYearSavings:      math.Round(twentyFiveYearSavings*100) / 100,
		SystemSizeKW:               input.SystemSizeKW,
		AnnualProductionKWh:        math.Round(annualProductionKWh*100) / 100,
		PanelCount:                 panelCount,
		ElectricalOffset:           electricalOffsetPct,
		CostPerWatt:                costPerWatt,
		SimplePaybackYears:         math.Round(simplePayback*100) / 100,
		BreakEvenYear:              breakEvenYear,
		Summary:                    summary,
	}, nil
}
