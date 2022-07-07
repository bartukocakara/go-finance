package v1

type ActionCreated struct {
	Created bool `json:"created"`
}

type ActionDeleted struct {
	Deleted bool `json:"deleted"`
}

type ActionUpdated struct {
	Updated bool `json:"updated"`
}
