package kotsclient

import (
	"time"

	"github.com/pkg/errors"
	"github.com/replicatedhq/replicated/pkg/graphql"
	"github.com/replicatedhq/replicated/pkg/types"
)

const kotsCreateCustomer = `
	mutation createKotsCustomer($name: String!, $channelId: String!, $expiresAt: String) {
		createKotsCustomer(name: $name, channelId: $channelId, isKots: true, airgap: false, type: "dev", entitlements: [], expiresAt: $expiresAt, isGitopsSupported: false) {
			id
			name 
			type
			expiresAt
			# cant fetch channels here 
        }
	}`

type GraphQLResponseCreateCustomer struct {
	Data   *CreateCustomerData `json:"data,omitempty"`
	Errors []graphql.GQLError  `json:"errors,omitempty"`
}

type CreateCustomerData struct {
	Customer *Customer `json:"createKotsCustomer"`
}

func (c *GraphQLClient) CreateCustomer(name, channel string, expiresIn time.Duration) (*types.Customer, error) {

	response := GraphQLResponseCreateCustomer{}

	request := graphql.Request{
		Query: kotsCreateCustomer,
		Variables: map[string]interface{}{
			"name":      name,
			"channelId": channel,
		},
	}
	if expiresIn > 0 {
		request.Variables["expiresAt"] = (time.Now().Add(expiresIn)).Format(time.RFC3339)
	}

	if err := c.ExecuteRequest(request, &response); err != nil {
		return nil, errors.Wrap(err, "execute gql request")
	}

	customer := types.Customer{
		ID:      response.Data.Customer.ID,
		Name:    response.Data.Customer.Name,
		Type:    response.Data.Customer.Type,
		Expires: &response.Data.Customer.ExpiresAt,
	}

	return &customer, nil
}
