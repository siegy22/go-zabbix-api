package zabbix_test

import (
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	zapi "github.com/claranet/go-zabbix-api"
)

func testCreateHost(group *zapi.HostGroup, t *testing.T) *zapi.Host {
	name := fmt.Sprintf("%s-%d", testGetHost(), rand.Int())
	iface := zapi.HostInterface{DNS: name, Port: "42", Type: zapi.Agent, UseIP: 0, Main: 1}
	hosts := zapi.Hosts{{
		Host:       name,
		Name:       "Name for " + name,
		GroupIds:   zapi.HostGroupIDs{{group.GroupID}},
		Interfaces: zapi.HostInterfaces{iface},
	}}

	err := testGetAPI(t).HostsCreate(hosts)
	if err != nil {
		t.Fatal(err)
	}
	return &hosts[0]
}

func testDeleteHost(host *zapi.Host, t *testing.T) {
	err := testGetAPI(t).HostsDelete(zapi.Hosts{*host})
	if err != nil {
		t.Fatal(err)
	}
}

func TestHosts(t *testing.T) {
	api := testGetAPI(t)

	group := testCreateHostGroup(t)
	defer testDeleteHostGroup(group, t)

	hosts, err := api.HostsGetByHostGroups(zapi.HostGroups{*group})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 0 {
		t.Errorf("Bad hosts: %#v", hosts)
	}

	host := testCreateHost(group, t)
	if host.HostID == "" || host.Host == "" {
		t.Errorf("Something is empty: %#v", host)
	}
	iface := host.Interfaces[0]
	host.GroupIds = nil
	host.Interfaces = nil

	newName := fmt.Sprintf("%s-%d", testGetHost(), rand.Int())
	host.Host = newName
	err = api.HostsUpdate(zapi.Hosts{*host})
	if err != nil {
		t.Fatal(err)
	}

	host2, err := api.HostGetByHost(host.Host)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(host, host2) {
		t.Errorf("Hosts are not equal:\n%#v\n%#v", host, host2)
	}

	host2, err = api.HostGetByID(host.HostID)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(host, host2) {
		t.Errorf("Hosts are not equal:\n%#v\n%#v", host, host2)
	}

	hosts, err = api.HostsGetByHostGroups(zapi.HostGroups{*group})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 1 {
		t.Errorf("Bad hosts: %#v", hosts)
	}

	hosts, err = api.HostsGet(zapi.Params{
		"hostids":               host.HostID,
		"selectInterfaces":      "extend",
		"selectMacros":          "extend",
		"selectParentTemplates": "extend",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 1 {
		t.Errorf("Bad hosts: %#v", hosts)
	}
	if len(hosts[0].Interfaces) != 1 {
		t.Errorf("Bad interfaces: %#v", hosts[0].Interfaces)
	}
	iface2 := hosts[0].Interfaces[0]
	iface.InterfaceID = iface2.InterfaceID
	if !reflect.DeepEqual(iface, iface2) {
		t.Errorf("Interfaces are not equal:\n%#v\n%#v", iface, iface2)
	}
	if len(hosts[0].Templates) != 0 {
		t.Errorf("Templates is not empty: %#v", hosts[0].Templates)
	}
	if len(hosts[0].UserMacros) != 0 {
		t.Errorf("UserMacros is not empty: %#v", hosts[0].UserMacros)
	}

	template := testCreateTemplate(group, t)
	defer testDeleteTemplate(template, t)

	iface = zapi.HostInterface{IP: "127.0.0.1", Port: "10050", Type: zapi.Agent, UseIP: 1, Main: 1}
	host.Interfaces = zapi.HostInterfaces{iface}
	host.TemplateIDs = zapi.TemplateIDs{{template.TemplateID}}
	macro := zapi.Macro{MacroName: "{$NAME}", Value: "Value"}
	host.UserMacros = zapi.Macros{macro}

	err = api.HostsUpdate(zapi.Hosts{*host})
	if err != nil {
		t.Fatal(err)
	}

	hosts, err = api.HostsGet(zapi.Params{
		"hostids":               host.HostID,
		"selectInterfaces":      "extend",
		"selectMacros":          "extend",
		"selectParentTemplates": "extend",
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 1 {
		t.Errorf("Bad hosts: %#v", hosts)
	}
	if len(hosts[0].Interfaces) != 1 {
		t.Errorf("Bad interfaces: %#v", hosts[0].Interfaces)
	}
	iface2 = hosts[0].Interfaces[0]
	iface.InterfaceID = iface2.InterfaceID
	if !reflect.DeepEqual(iface, iface2) {
		t.Errorf("Interfaces are not equal:\n%#v\n%#v", iface, iface2)
	}
	if len(hosts[0].Templates) != 1 {
		t.Errorf("Bad templates: %#v", hosts[0].Templates)
	}
	template2 := hosts[0].Templates[0]
	if template.Host != template2.Host {
		t.Errorf("Templates are not equal:\n%#v\n%#v", template.Host, template2.Host)
	}
	if len(hosts[0].UserMacros) != 1 {
		t.Errorf("Bad userMacros: %#v", hosts[0].UserMacros)
	}
	macro2 := hosts[0].UserMacros[0]
	macro.HostID = hosts[0].HostID
	if !reflect.DeepEqual(macro, macro2) {
		t.Errorf("UserMacros are not equal:\n%#v\n%#v", macro, macro2)
	}

	testDeleteHost(host, t)

	hosts, err = api.HostsGetByHostGroups(zapi.HostGroups{*group})
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 0 {
		t.Errorf("Bad hosts: %#v", hosts)
	}
}
