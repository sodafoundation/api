package nimble

const (
	DriverName = "HPE Nimble Storage Driver"
)

const (
	ThickLuntype         = 0
	ThinLuntype          = 1
	MaxNameLength        = 31
	MaxDescriptionLength = 170
	PortNumPerContr      = 2
	PwdExpired           = 3
	PwdReset             = 4
)

// Error Code
const (
	ErrorUnauthorizedToServer = "SM_http_unauthorized"
	ErrorSmVolSizeDecreased   = "SM_vol_size_decreased"
	ErrorSmHttpConflict       = "SM_http_conflict"
)
