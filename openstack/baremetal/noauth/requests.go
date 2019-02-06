package noauth

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
)

// EndpointOpts specifies a "noauth" Ironic Endpoint.
type EndpointOpts struct {
	// IronicEndpoint [required] is currently only used with "noauth" Ironic.
	// An Ironic endpoint with "auth_strategy=noauth" is necessary, for example:
	// http://ironic.example.com:8776/v1.
	IronicEndpoint string
}

func initClientOpts(client *gophercloud.ProviderClient, eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	sc := new(gophercloud.ServiceClient)
	if eo.IronicEndpoint == "" {
		return nil, fmt.Errorf("IronicEndpoint is required")
	}

	endpoint := fmt.Sprintf("%s%s", gophercloud.NormalizeURL(eo.IronicEndpoint), client.TokenID)
	sc.Endpoint = gophercloud.NormalizeURL(endpoint)
	sc.ProviderClient = client
	return sc, nil
}

// NewBlockStorageNoAuth creates a ServiceClient that may be used to access a
// "noauth" bare metal service.
func NewBaremetalNoAuth(eo EndpointOpts) (*gophercloud.ServiceClient, error) {
	return initClientOpts(&gophercloud.ProviderClient{
		TokenID: "fake-token",
	}, eo)
}
