package zabbix

import (
	"encoding/json"
)


type (
	// InterfaceType different interface type
	InterfaceType int
)

const (
	// Differente type of zabbix interface
	// see "type" in https://www.zabbix.com/documentation/3.2/manual/api/reference/hostinterface/object

	// Agent type
	Agent InterfaceType = 1
	// SNMP type
	SNMP InterfaceType = 2
	// IPMI type
	IPMI InterfaceType = 3
	// JMX type
	JMX InterfaceType = 4
)

type InterfaceDetails struct {
	Version   int    `json:"version,string"`
	Community string `json:"community"`
}

// HostInterface represents zabbix host interface type
// https://www.zabbix.com/documentation/3.2/manual/api/reference/hostinterface/object
type HostInterface struct {
	InterfaceID string           `json:"interfaceid,omitempty"`
	DNS         string           `json:"dns"`
	IP          string           `json:"ip"`
	Main        int              `json:"main,string"`
	Port        string           `json:"port"`
	Type        InterfaceType    `json:"type,string"`
	UseIP       int              `json:"useip,string"`
	Details     InterfaceDetails `json:"details,omitempty"`
}

// HostInterfaces is an array of HostInterface
type HostInterfaces []HostInterface

// Custom unmarshal function
func (h *HostInterface) UnmarshalJSON(data []byte) error {
	type Alias HostInterface
	aux := &struct {
		Details json.RawMessage `json:"details"`
		*Alias
	}{
		Alias: (*Alias)(h),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle the details field based on its type
	if string(aux.Details) != "[]" {
		var details InterfaceDetails
		if err := json.Unmarshal(aux.Details, &details); err != nil {
			return err
		}
		h.Details = details
	} else {
		h.Details = InterfaceDetails{}
	}

	return nil
}
