// Code generated by "libovsdb.modelgen"
// DO NOT EDIT.

package ovs

const SFlowTable = "sFlow"

// SFlow defines an object in sFlow table
type SFlow struct {
	UUID        string            `ovsdb:"_uuid"`
	Agent       *string           `ovsdb:"agent"`
	ExternalIDs map[string]string `ovsdb:"external_ids"`
	Header      *int              `ovsdb:"header"`
	Polling     *int              `ovsdb:"polling"`
	Sampling    *int              `ovsdb:"sampling"`
	Targets     []string          `ovsdb:"targets"`
}