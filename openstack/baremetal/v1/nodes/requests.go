package nodes

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to the
// List request.
type ListOptsBuilder interface {
	ToNodeListQuery() (string, error)
	ToNodeListDetailQuery() (string, error)
}

// Provision state reports the current provision state of the node, these are only used in filtering
type ProvisionState string

const (
	Enroll       ProvisionState = "enroll"
	Verifying                   = "verifying"
	Manageable                  = "manageable"
	Available                   = "available"
	Active                      = "active"
	DeployWait                  = "wait call-back"
	Deploying                   = "deploying"
	DeployFail                  = "deploy failed"
	DeployDone                  = "deploy complete"
	Deleting                    = "deleting"
	Deleted                     = "deleted"
	Cleaning                    = "cleaning"
	CleanWait                   = "clean wait"
	CleanFail                   = "clean failed"
	Error                       = "error"
	Rebuild                     = "rebuild"
	Inpsecting                  = "inspecting"
	InspectFail                 = "inspect failed"
	InspectWait                 = "inspect wait"
	Adopting                    = "adopting"
	AdoptFail                   = "adopt failed"
	Rescue                      = "rescue"
	RescueFail                  = "rescue failed"
	Rescuing                    = "rescuing"
	UnrescueFail                = "unrescue failed"
)

// ListOpts allows the filtering and sorting of paginated collections through
// the API. Filtering is achieved by passing in struct field values that map to
// the node attributes you want to see returned. Marker and Limit are used
// for pagination.
type ListOpts struct {
	// Filter the list by specific instance UUID
	InstanceUUID string `q:"instance_uuid"`

	// Filter the list by chassis UUID
	ChassisUUID string `q:"chassis_uuid"`

	// Filter the list by maintenance set to True or False
	Maintenance bool `q:"maintenance"`

	// Nodes which are, or are not, associated with an instance_uuid.
	Associated bool `q:"associated"`

	// Only return those with the specified provision_state.
	ProvisionState ProvisionState `q:"provision_state"`

	// Filter the list with the specified driver.
	Driver string `q:"driver"`

	// Filter the list with the specified resource class.
	ResourceClass string `q:"resource_class"`

	// Filter the list with the specified conductor_group.
	ConductorGroup string `q:"conductor_group"`

	// Filter the list with the specified fault.
	Fault string `q:"fault"`

	// One or more fields to be returned in the response.
	Fields []string `q:"fields"`

	// Requests a page size of items.
	Limit int `q:"limit"`

	// The ID of the last-seen item.
	Marker string `q:"marker"`

	// Sorts the response by the requested sort direction.
	SortDir string `q:"sort_dir"`

	// Sorts the response by the this attribute value.
	SortKey string `q:"sort_key"`

	// A string or UUID of the tenant who owns the baremetal node.
	Owner string `q:"owner"`
}

// ToNodeListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToNodeListQuery() (string, error) {
	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// List makes a request against the API to list nodes accessible to you.
func List(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToNodeListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return NodePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// ToNodeListDetailQuery formats a ListOpts into a query string for the list details API.
func (opts ListOpts) ToNodeListDetailQuery() (string, error) {
	// Detail endpoint can't filter by Fields
	if len(opts.Fields) > 0 {
		return "", fmt.Errorf("fields is not a valid option when getting a detailed listing of nodes")
	}

	q, err := gophercloud.BuildQueryString(opts)
	return q.String(), err
}

// Return a list of bare metal Nodes with complete details. Some filtering is possible by passing in flags in ListOpts,
// but you cannot limit by the fields returned.
func ListDetail(client *gophercloud.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	// This URL is deprecated. In the future, we should compare the microversion and if >= 1.43, hit the listURL
	// with ListOpts{Detail: true,}
	url := listDetailURL(client)
	if opts != nil {
		query, err := opts.ToNodeListDetailQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return NodePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get requests details on a single node, by ID.
func Get(client *gophercloud.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to the
// Create request.
type CreateOptsBuilder interface {
	ToNodeCreateMap() (map[string]interface{}, error)
}

// CreateOpts specifies node creation parameters.
type CreateOpts struct {
	// The boot interface for a Node, e.g. “pxe”.
	BootInterface string `json:"boot_interface,omitempty"`

	// The conductor group for a node. Case-insensitive string up to 255 characters, containing a-z, 0-9, _, -, and ..
	ConductorGroup string `json:"conductor_group,omitempty"`

	// The console interface for a node, e.g. “no-console”.
	ConsoleInterface string `json:"console_interface,omitempty"`

	// The deploy interface for a node, e.g. “iscsi”.
	DeployInterface string `json:"deploy_interface,omitempty"`

	// All the metadata required by the driver to manage this Node. List of fields varies between drivers, and can
	// be retrieved from the /v1/drivers/<DRIVER_NAME>/properties resource.
	DriverInfo map[string]interface{} `json:"driver_info,omitempty"`

	// name of the driver used to manage this Node.
	Driver string `json:"driver,omitempty"`

	// A set of one or more arbitrary metadata key and value pairs.
	Extra map[string]interface{} `json:"extra,omitempty"`

	// The interface used for node inspection, e.g. “no-inspect”.
	InspectInterface string `json:"inspect_interface,omitempty"`

	// Interface for out-of-band node management, e.g. “ipmitool”.
	ManagementInterface string `json:"management_interface,omitempty"`

	// Human-readable identifier for the Node resource. May be undefined. Certain words are reserved.
	Name string `json:"name,omitempty"`

	// Which Network Interface provider to use when plumbing the network connections for this Node.
	NetworkInterface string `json:"network_interface,omitempty"`

	// Interface used for performing power actions on the node, e.g. “ipmitool”.
	PowerInterface string `json:"power_interface,omitempty"`

	// Physical characteristics of this Node. Populated during inspection, if performed. Can be edited via the REST
	// API at any time.
	Properties map[string]interface{} `json:"properties,omitempty"`

	// Interface used for configuring RAID on this node, e.g. “no-raid”.
	RAIDInterface string `json:"raid_interface,omitempty"`

	// The interface used for node rescue, e.g. “no-rescue”.
	RescueInterface string `json:"rescue_interface,omitempty"`

	// A string which can be used by external schedulers to identify this Node as a unit of a specific type
	// of resource.
	ResourceClass string `json:"resource_class,omitempty"`

	// Interface used for attaching and detaching volumes on this node, e.g. “cinder”.
	StorageInterface string `json:"storage_interface,omitempty"`

	// The UUID for the resource.
	UUID string `json:"uuid,omitempty"`

	// Interface for vendor-specific functionality on this node, e.g. “no-vendor”.
	VendorInterface string `json:"vendor_interface,omitempty"`

	// A string or UUID of the tenant who owns the baremetal node.
	Owner string `json:"owner,omitempty"`
}

// ToNodeCreateMap assembles a request body based on the contents of a CreateOpts.
func (opts CreateOpts) ToNodeCreateMap() (map[string]interface{}, error) {
	body, err := gophercloud.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Create requests a node to be created
func Create(client *gophercloud.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToNodeCreateMap()
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Post(createURL(client), reqBody, &r.Body, nil)
	return
}

type Patch interface {
	ToNodeUpdateMap() map[string]interface{}
}

// UpdateOpts is a slice of Patches used to update a node
type UpdateOpts []Patch

type UpdateOp string

const (
	ReplaceOp UpdateOp = "replace"
	AddOp     UpdateOp = "add"
	RemoveOp  UpdateOp = "remove"
)

type UpdateOperation struct {
	Op    UpdateOp `json:"op,required"`
	Path  string   `json:"path,required"`
	Value string   `json:"value,omitempty"`
}

func (opts UpdateOperation) ToNodeUpdateMap() map[string]interface{} {
	return map[string]interface{}{
		"op":    opts.Op,
		"path":  opts.Path,
		"value": opts.Value,
	}
}

// Update requests that a node be updated
func Update(client *gophercloud.ServiceClient, id string, opts UpdateOpts) (r UpdateResult) {
	body := make([]map[string]interface{}, len(opts))
	for i, patch := range opts {
		body[i] = patch.ToNodeUpdateMap()
	}

	resp, err := client.Request("PATCH", updateURL(client, id), &gophercloud.RequestOpts{
		JSONBody: &body,
		OkCodes:  []int{200},
	})

	r.Body = resp.Body
	r.Header = resp.Header
	r.Err = err

	return
}

// Delete requests that a node be removed
func Delete(client *gophercloud.ServiceClient, id string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, id), nil)
	return
}

// Request that Ironic validate whether the Node’s driver has enough information to manage the Node. This polls each
// interface on the driver, and returns the status of that interface.
func Validate(client *gophercloud.ServiceClient, id string) (r ValidateResult) {
	_, r.Err = client.Get(validateURL(client, id), &r.Body, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	return
}
