package target_provision_state

// target_provision_state.State is used when telling Ironic which provision state
type State string

const (
	Active   State = "active"
	Delete   State = "delete"
	Manage   State = "manage"
	Provide  State = "provide"
	Inspect  State = "inspect"
	Abort    State = "abort"
	Clean    State = "clean"
	Adopt    State = "adopt"
	Rescue   State = "rescue"
	Unrescue State = "unrescue"
)
