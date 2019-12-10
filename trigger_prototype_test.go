package zabbix_test

import (
	"fmt"
	"testing"

	"github.com/claranet/go-zabbix-api"
	dd "github.com/claranet/go-zabbix-api"
)

func createTriggerPrototype(template *dd.Template, item *dd.ItemPrototype, t *testing.T) *dd.TriggerPrototype {
	triggers := dd.TriggerPrototypes{{
		Description: "test trigger prototype",
		Expression:  fmt.Sprintf("{%s:%s.last()}=312", template.Host, item.Key),
		Priority:    zabbix.Warning,
	}}
	err := getAPI(t).TriggerPrototypesCreate(triggers)
	if err != nil {
		t.Fatal(err)
	}
	return &triggers[0]
}

func deleteTriggerPrototype(trigger *dd.TriggerPrototype, t *testing.T) {
	err := getAPI(t).TriggerPrototypesDelete(dd.TriggerPrototypes{*trigger})
	if err != nil {
		t.Fatal(err)
	}
}

func testTriggerPrototype(t *testing.T) {
	api := getAPI(t)

	hostGroup := CreateHostGroup(t)
	defer DeleteHostGroup(hostGroup, t)

	template := CreateTemplate(hostGroup, t)
	defer DeleteTemplate(template, t)

	lldRule := CreateLLDRule(template, t)
	defer DeleteLLDRule(lldRule, t)

	itemPrototype := createItemPrototype(template, lldRule, t)
	defer deleteItemPrototype(itemPrototype, t)

	triggerPrototype := createTriggerPrototype(template, itemPrototype, t)

	getTriggerPrototype, err := api.TriggerPrototypesGet(dd.Params{"triggerids": triggerPrototype.TriggerID})
	if err != nil {
		t.Error(err)
	}
	if len(getTriggerPrototype) != 1 {
		t.Errorf("Expecting one result and got : %d", len(getTriggerPrototype))
	}

	triggerPrototype.Description = "update_trigger_prototype_name"
	err = api.TriggerPrototypesUpdate(dd.TriggerPrototypes{*triggerPrototype})
	if err != nil {
		t.Error(err)
	}

	getByIdTriggerPrototype, err := api.TriggerPrototypeGetByID(triggerPrototype.TriggerID)
	if err != nil {
		t.Error(err)
	}
	if getByIdTriggerPrototype.Description != triggerPrototype.Description {
		t.Errorf("Trigger prototype description is %q and should be %q", getByIdTriggerPrototype.Description, triggerPrototype.Description)
	}

	deleteTriggerPrototype(triggerPrototype, t)
}
