package main

import (
	"fmt"
	"os"
	"runtime"
)

var (
	userAgent          string
	introspectionQuery string
)

func getIntrospectionQuery() string {
	if len(introspectionQuery) > 0 {
		return introspectionQuery
	} else {
		return `
                query IntrospectionQuery {
                  __schema {
                   s
                    queryType { name }
                    mutationType { name }
                    subscriptionType { name }
                    types {
                      ...FullType
                    }
                    directives {
                      name
                      description
                      locations
                     s
                      args {
                        ...InputValue
                      }
                    }
                  }
                }
                fragment FullType on __Type {
                  kind
                  name
                  description
                 s
                 s
                  fields(includeDeprecated: true) {
                    name
                    description
                    args {
                      ...InputValue
                    }
                    type {
                      ...TypeRef
                    }
                    isDeprecated
                    deprecationReason
                  }
                  inputFields {
                    ...InputValue
                  }
                  interfaces {
                    ...TypeRef
                  }
                  enumValues(includeDeprecated: true) {
                    name
                    description
                    isDeprecated
                    deprecationReason
                  }
                  possibleTypes {
                    ...TypeRef
                  }
                }
                fragment InputValue on __InputValue {
                  name
                  description
                  type { ...TypeRef }
                  defaultValue
                 s
                 s
                }
                fragment TypeRef on __Type {
                  kind
                  name
                  ofType {
                    kind
                    name
                    ofType {
                      kind
                      name
                      ofType {
                        kind
                        name
                        ofType {
                          kind
                          name
                          ofType {
                            kind
                            name
                            ofType {
                              kind
                              name
                              ofType {
                                kind
                                name
                              }
                            }
                          }
                        }
                      }
                    }
                  }
                }
`
	}
}

func getUserAgent() string {
	if len(userAgent) > 0 {
		return userAgent
	}

	hostname, _ := os.Hostname()
	userAgent = fmt.Sprintf(
		"Webex Go SDK(v%s) - OS(%s) - hostname(%s) - Go Version(%s)",
		VERSION,
		runtime.GOOS,
		hostname,
		runtime.Version())

	return userAgent
}
