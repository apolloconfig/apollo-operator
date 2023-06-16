package reconcile

import "apollo.io/apollo-operator/pkg/reconcile/apolloportal"

var (
	apolloPortal      ApolloObject
	apolloEnvironment ApolloObject
	apolloAllInOne    ApolloObject
)

func init() {
	apolloPortal = apolloportal.NewApolloPortal()
	apolloEnvironment = apolloportal.NewApolloPortal()
	apolloAllInOne = apolloportal.NewApolloPortal()
}

func ApolloPortal() ApolloObject {
	return apolloPortal
}

func ApolloEnvironment() ApolloObject {
	return apolloEnvironment
}

func ApolloAllInOne() ApolloObject {
	return apolloAllInOne
}
