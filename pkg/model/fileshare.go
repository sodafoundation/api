package model


type FileShareSpec struct {
	*BaseModel

	// The uuid of the project that the fileshare belongs to.
	TenantId string `json:"tenantId,omitempty"`

	// The uuid of the user that the fileshare belongs to.
	// +optional
	UserId string `json:"userId,omitempty"`

	// The name of the fileshare.
	Name string `json:"name,omitempty"`

	// The description of the fileshare.
	// +optional
	Description string `json:"description,omitempty"`

	// The size of the fileshare requested by the user.
	// Default unit of fileshare Size is GB.
	Size int64 `json:"size,omitempty"`

	// The locality that fileshare belongs to.
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// The status of the fileshare.
	// One of: "available", "error", "in-use", etc.
	Status string `json:"status,omitempty"`

	// The uuid of the pool which the fileshare belongs to.
	// +readOnly
	PoolId string `json:"poolId,omitempty"`

	// The uuid of the profile which the fileshare belongs to.
	ProfileId string `json:"profileId,omitempty"`

	// Metadata should be kept until the scemantics between opensds fileshare
	// and backend storage resouce description are clear.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`

}
