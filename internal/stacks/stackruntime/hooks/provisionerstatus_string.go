// Code generated by "stringer -type=ProvisionerStatus resource_instance.go"; DO NOT EDIT.

package hooks

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ProvisionerStatusInvalid-0]
	_ = x[ProvisionerProvisioning-112]
	_ = x[ProvisionerProvisioned-80]
	_ = x[ProvisionerErrored-69]
}

const (
	_ProvisionerStatus_name_0 = "ProvisionerStatusInvalid"
	_ProvisionerStatus_name_1 = "ProvisionerErrored"
	_ProvisionerStatus_name_2 = "ProvisionerProvisioned"
	_ProvisionerStatus_name_3 = "ProvisionerProvisioning"
)

func (i ProvisionerStatus) String() string {
	switch {
	case i == 0:
		return _ProvisionerStatus_name_0
	case i == 69:
		return _ProvisionerStatus_name_1
	case i == 80:
		return _ProvisionerStatus_name_2
	case i == 112:
		return _ProvisionerStatus_name_3
	default:
		return "ProvisionerStatus(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
