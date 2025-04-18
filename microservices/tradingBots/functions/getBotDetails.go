package functions

import (
	"context"
	"encoding/json"

	"cryptobotmanager.com/cbm-backend/resolvers/graph/model"
	"cryptobotmanager.com/cbm-backend/shared/graph"
	"github.com/Khan/genqlient/graphql"
)

// This function takes in two file paths for strategy and exit parameters, loads the
// parameters from the respective files, and returns the loaded strategy and exit parameters.
//
// The loadStrategies and loadExits functions are called to load the parameters from
// the files specified in the strategyParams and exitParams arguments, respectively. If
// an error occurs during the loading of the strategy or exit parameters, the function logs an
// error message, exits the program with status code 1, and returns nil for both
// strategyDetails and exitDetails.
//
// If both sets of parameters are loaded successfully, strategyDetails and
// exitDetails are returned along with nil for the err return value.
func GetParameters(ctx context.Context, client graphql.Client) ([]model.StrategyInput, error) {
	response, err := graph.GetAllStrategies(ctx, client)
	if err != nil {
		return nil, err
	}

	var strategyDetails []model.StrategyInput

	// Access the "data" key and then the "getAllStrategies" key
	for _, obj := range response.GetAllStrategies {
		// Check if the strategy is not tested
		if !obj.Tested {
			objJSON, err := json.Marshal(obj)
			if err != nil {
				return nil, err
			}

			var details model.StrategyInput
			err = json.Unmarshal(objJSON, &details)
			if err != nil {
				return nil, err
			}

			// Append the converted object to the slice
			strategyDetails = append(strategyDetails, details)
		}
	}

	return strategyDetails, nil
}
