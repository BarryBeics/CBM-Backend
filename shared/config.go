package shared

type AppConfig struct {
	TradeDuration                 int
	PackageNames, TestExemptFuncs []string
}

func GetDefaultCfg() AppConfig {

	// VALUES THAT WILL COME FROM A UI AT SOME POINT
	tradeDuration := 5

	packageNames := []string{
		"filters", "goBot", "priceData"}
	testExemptFuncs := []string{
		"GetDefaultCfg",
		"PrettyPrint",
		"SetupLogger",
		"WriteDataToJSON",
		"PrintFunctionsWithoutTestCoverage",
		"SavePricesToDB"}

	//REMEMBER TO BRING PERMUTATION INTO HERE

	cfg := &AppConfig{
		TradeDuration:   tradeDuration,
		PackageNames:    packageNames,
		TestExemptFuncs: testExemptFuncs,
	}

	return *cfg
}
