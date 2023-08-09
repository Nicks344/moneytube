package serverutils

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"

	"github.com/meandrewdev/graphqlws"

	"github.com/graphql-go/graphql"
)

var subscriptionManager graphqlws.SubscriptionManager
var schema graphql.Schema

type graphQLRequest struct {
	Query         string                 `json:"query" url:"query" schema:"query"`
	Variables     map[string]interface{} `json:"variables" url:"variables" schema:"variables"`
	OperationName string                 `json:"operationName" url:"operationName" schema:"operationName"`
}

func SetSchema(s graphql.Schema) {
	schema = s
}

func GetGQLHTTPHandler(authKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Content-Length, X-Requested-With, Key")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusOK)
			return
		}
		/*
			if r.Header["Key"][0] != authKey {
				w.Write([]byte("Auth error"))
				return
			}
		*/
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		/* the array of pre-defined struct */
		var requests []*graphQLRequest
		/* Parse JSON body */
		err = json.Unmarshal(body, &requests)
		if err != nil {
			return
		}

		results := make([]*graphql.Result, len(requests))

		ctx := r.Context()
		var wg sync.WaitGroup
		for reqOrder, r := range requests {
			wg.Add(1)

			go func(order int, req *graphQLRequest) {
				defer wg.Done()
				params := graphql.Params{
					Schema:         schema,
					RequestString:  req.Query,
					VariableValues: req.Variables,
					OperationName:  req.OperationName,
					Context:        ctx,
				}
				results[order] = graphql.Do(params)
			}(reqOrder, r)
		}
		/* this is going to wait for all query to be executed */
		wg.Wait()

		w.Header().Add("Content-Type", "application/json; charset=utf-8")

		var buff []byte
		w.WriteHeader(http.StatusOK)
		buff, _ = json.Marshal(results)

		w.Write(buff)
	})
}

func GetGQLWsHandler(authKey string) http.Handler {
	subscriptionManager = graphqlws.NewSubscriptionManager(&schema)
	// Create a WebSocket/HTTP handler
	return graphqlws.NewHandler(graphqlws.HandlerConfig{
		// Wire up the GraphqL WebSocket handler with the subscription manager
		SubscriptionManager: subscriptionManager,

		// Optional: Add a hook to resolve auth tokens into users that are
		// then stored on the GraphQL WS connections
		Authenticate: func(authToken string) (interface{}, error) {
			/*
				if authToken != authKey {
					return nil, errors.New("Auth error")
				}
			*/
			return "Auth success", nil
		},
	})
}

func OnGQLEvent(event string, id string, data interface{}) {
	if subscriptionManager == nil {
		return
	}
	for _, subsContainer := range subscriptionManager.Subscriptions() {
		if subsContainer == nil {
			return
		}
		for _, subscription := range subsContainer {
			if subscription.OperationName == event {
				ctx := context.Background()

				var hasID bool
				if subscription.Variables != nil && len(subscription.Variables) > 0 {
					if varID, ok := subscription.Variables["id"]; ok && varID.(string) == id {
						ctx = context.WithValue(ctx, id, data)
						hasID = true
					}
				}
				if !hasID {
					continue
				}
				// Re-execute the subscription query
				params := graphql.Params{
					Schema:         schema, // The GraphQL schema
					RequestString:  subscription.Query,
					VariableValues: subscription.Variables,
					OperationName:  subscription.OperationName,
					Context:        ctx,
				}
				result := graphql.Do(params)

				// Send query results back to the subscriber at any point
				sendData := graphqlws.DataMessagePayload{
					// Data can be anything (interface{})
					Data: result.Data,
					// Errors is optional ([]error)
					Errors: graphqlws.ErrorsFromGraphQLErrors(result.Errors),
				}
				subscription.SendData(&sendData)
			}
		}
	}
}
