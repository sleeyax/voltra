package utils

import "math"

// RoundStepSize rounds a given quantity to a specific step size.
//
// Parameters:
//   - quantity: The quantity to be rounded.
//   - stepSize: The step size to round to.
//
// Returns:
//
//	The rounded quantity.
func RoundStepSize(quantity, stepSize float64) float64 {
	// Calculate the precision required for rounding.
	precision := int(math.Round(-math.Log10(stepSize)))

	// Round the quantity to the calculated precision.
	return math.Round(quantity*math.Pow(10, float64(precision))) / math.Pow(10, float64(precision))
}
