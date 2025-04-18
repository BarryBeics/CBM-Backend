package shared

type AppConfig struct {
	TopAverages                   []int
	TradeDuration                 int
	PackageNames, TestExemptFuncs []string
	ActiveMarketThreshold         float64
}

func GetDefaultCfg() AppConfig {

	var activeMarketThreshold = 0.1

	// Select how many top coin to return
	topAverages := []int{3, 5, 10}

	// VALUES THAT WILL COME FROM A UI AT SOME POINT
	tradeDuration := 5

	packageNames := []string{
		"filters", "backTesting", "resolvers"}
	testExemptFuncs := []string{
		"GetDefaultCfg",
		"PrettyPrint",
		"SetupLogger",
		"WriteDataToJSON",
		"PrintFunctionsWithoutTestCoverage",
		"SavePricesToDB"}

	//REMEMBER TO BRING PERMUTATION INTO HERE

	cfg := &AppConfig{
		ActiveMarketThreshold: activeMarketThreshold,
		TradeDuration:         tradeDuration,
		PackageNames:          packageNames,
		TestExemptFuncs:       testExemptFuncs,
		TopAverages:           topAverages,
	}

	return *cfg
}
