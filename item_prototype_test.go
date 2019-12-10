package zabbix_test

import (
	"testing"

	dd "github.com/claranet/go-zabbix-api"
)

func createItemPrototype(template *dd.Template, lldRule *dd.LLDRule, t *testing.T) *dd.ItemPrototype {
	items := dd.ItemPrototypes{{
		RuleID:      lldRule.ItemID,
		Delay:       "30",
		HostID:      template.TemplateID,
		InterfaceID: "0",
		Key:         "key.lala.lolo",
		Type:        dd.ZabbixAgent,
		Name:        "item prototype test",
		ValueType:   dd.Unsigned,
	}}
	err := getAPI(t).ItemPrototypesCreate(items)
	if err != nil {
		t.Fatal(err)
	}
	return &items[0]
}

func deleteItemPrototype(item *dd.ItemPrototype, t *testing.T) {
	err := getAPI(t).ItemPrototypesDelete(dd.ItemPrototypes{*item})
	if err != nil {
		t.Fatal(err)
	}
}

func testItemPrototype(t *testing.T) {
	api := getAPI(t)

	hostGroup := CreateHostGroup(t)
	defer DeleteHostGroup(hostGroup, t)

	template := CreateTemplate(hostGroup, t)
	defer DeleteTemplate(template, t)

	lldRule := CreateLLDRule(template, t)
	defer DeleteLLDRule(lldRule, t)

	itemPrototype := createItemPrototype(template, lldRule, t)

	getItemPrototype, err := api.ItemPrototypesGet(dd.Params{"itemids": itemPrototype.ItemID})
	if err != nil {
		t.Error(err)
	}
	if len(getItemPrototype) != 1 {
		t.Errorf("Expecting one result and got : %d", len(getItemPrototype))
	}

	itemPrototype.Name = "update_item_prototype_name"
	err = api.ItemPrototypesUpdate(dd.ItemPrototypes{*itemPrototype})
	if err != nil {
		t.Error(err)
	}

	getByIdItemPrototype, err := api.ItemPrototypeGetByID(itemPrototype.ItemID)
	if err != nil {
		t.Error(err)
	}
	if getByIdItemPrototype.Name != itemPrototype.Name {
		t.Errorf("Item prototype name is %q and should be %q", getByIdItemPrototype.Name, itemPrototype.Name)
	}

	deleteItemPrototype(itemPrototype, t)
}
