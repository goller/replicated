package kotsclient

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedhq/replicated/pkg/graphql"
	"github.com/replicatedhq/replicated/pkg/types"
)

const kotsListCustomers = `
	query customers($appId: String!, $appType: String!) {
		customers(appId: $appId, appType: $appType) {
            customers {
		        id
		        name
				createdAt
				expiresAt
				isArchived
				type
				shipInstallStatus {
					status
				}
		        channels {
		            id
		            name
					description
		            currentVersion
					created
					updated
					releaseSequence
		        }
            }
        }
	}
	`

const kotsListCustomer = `
    query customers($appId: String, $appType: String, $customerName: String) {
		customers(appId: $appId, appType: $appType, customerName: $customerName) {
            customers {
		        id
		        name
				createdAt
				expiresAt
				isArchived
				type
				shipInstallStatus {
					status
				}
		        channels {
		            id
		            name
					description
		            currentVersion
					releaseSequence
		        }
            }
        }
	}
	`

type GraphQLResponseListCustomers struct {
	Data   *CustomerDataWrapper `json:"data,omitempty"`
	Errors []graphql.GQLError   `json:"errors,omitempty"`
}

type CustomerDataWrapper struct {
	Customers CustomerData `json:"customers"`
}

type CustomerData struct {
	Customers []Customer `json:"customers"`
}

type ShipInstallStatus struct {
	Status string `json:"status"`
}

type Channel struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	CurrentVersion  string    `json:"currentVersion"`
	Created         time.Time `json:"created"`
	Updated         time.Time `json:"updated"`
	ReleaseSequence int64     `json:"releaseSequence"`
}

type Customer struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Channels          []Channel         `json:"channels"`
	Type              string            `json:"type"`
	IsArchived        bool              `json:"isArchived"`
	CreatedAt         time.Time         `json:"createdAt"`
	ExpiresAt         time.Time         `json:"expiresAt"`
	ShipInstallStatus ShipInstallStatus `json:"shipInstallStatus"`
}

func (c *GraphQLClient) Customers(appID string) ([]Customer, error) {
	response := GraphQLResponseListCustomers{}

	request := graphql.Request{
		Query: kotsListCustomers,

		Variables: map[string]interface{}{
			"appId":   appID,
			"appType": "kots",
		},
	}

	if err := c.ExecuteRequest(request, &response); err != nil {
		return nil, errors.Wrap(err, "execute gql request")
	}

	return response.Data.Customers.Customers, nil
}

func (c *GraphQLClient) Customer(appID, name string) ([]Customer, error) {
	response := GraphQLResponseListCustomers{}

	request := graphql.Request{
		Query: kotsListCustomer,

		Variables: map[string]interface{}{
			"appId":        appID,
			"appType":      "kots",
			"customerName": name,
		},
	}

	if err := c.ExecuteRequest(request, &response); err != nil {
		return nil, errors.Wrap(err, "execute gql request")
	}

	return response.Data.Customers.Customers, nil
}

func (c *GraphQLClient) ListCustomers(appID string) ([]types.Customer, error) {
	cs, err := c.Customers(appID)
	if err != nil {
		return nil, err
	}

	customers := make([]types.Customer, 0, 0)
	for _, kotsCustomer := range cs {

		kotsChannels := make([]types.Channel, 0, 0)
		for _, kotsChannel := range kotsCustomer.Channels {
			channel := types.Channel{
				ID:              kotsChannel.ID,
				Name:            kotsChannel.Name,
				ReleaseLabel:    kotsChannel.CurrentVersion,
				ReleaseSequence: kotsChannel.ReleaseSequence,
			}
			kotsChannels = append(kotsChannels, channel)
		}
		customer := types.Customer{
			ID:       kotsCustomer.ID,
			Name:     kotsCustomer.Name,
			Type:     kotsCustomer.Type,
			Channels: kotsChannels,
			Expires:  &kotsCustomer.ExpiresAt,
		}

		customers = append(customers, customer)
	}

	return customers, nil
}
