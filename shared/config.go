package shared

type AppConfig struct {
	TradeDuration                 int
	PackageNames, TestExemptFuncs []string
	ActiveMarketThreshold         float64
}

func GetDefaultCfg() AppConfig {

	var activeMarketThreshold = 0.1

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
	}

	return *cfg
}
