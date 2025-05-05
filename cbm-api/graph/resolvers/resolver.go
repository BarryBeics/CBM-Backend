package resolvers

import "cryptobotmanager.com/cbm-backend/cbm-api/graph/generated"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app.

type Resolver struct{}

func (r *Resolver) Project() generated.ProjectResolver {
	return &projectResolver{r}
}

func (r *Resolver) Query() generated.QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Mutation() generated.MutationResolver {
	return &mutationResolver{r}
}

type projectResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
