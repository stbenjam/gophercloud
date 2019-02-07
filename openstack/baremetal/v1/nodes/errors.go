package nodes

import "github.com/gophercloud/gophercloud"

// When specifying a config_drive for a host, only a path to a file on disk, or a raw string value may be specified.
type ErrConfigDriveMustBeEitherPathOrValueNotBoth struct{ gophercloud.ErrInvalidInput }

func (e ErrConfigDriveMustBeEitherPathOrValueNotBoth) Error() string {
	return "Either a path to a file, or a string value may be specified for a config drive. Not both."
}
