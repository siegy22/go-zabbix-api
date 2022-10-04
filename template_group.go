package zabbix

// TemplateGroup represent Zabbix template group object, new in v6.2
// https://www.zabbix.com/documentation/6.2/en/manual/api/reference/templategroup/object
type TemplateGroup struct {
	GroupID string `json:"groupid,omitempty"`
	Name    string `json:"name"`
}

// TemplateGroups is an array of TemplateGroup
type TemplateGroups []TemplateGroup

// TemplateGroupsGet Wrapper for templategroup.get
// https://www.zabbix.com/documentation/6.2/en/manual/api/reference/templategroup/get
func (api *API) TemplateGroupsGet(params Params) (res TemplateGroups, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("templategroup.get", params, &res)
	return
}

// TemplateGroupGetByID Gets host group by Id only if there is exactly 1 matching host group.
// https://www.zabbix.com/documentation/6.2/en/manual/api/reference/templategroup/get
func (api *API) TemplateGroupGetByID(id string) (res *TemplateGroup, err error) {
	groups, err := api.TemplateGroupsGet(Params{"groupids": id})
	if err != nil {
		return
	}

	if len(groups) == 1 {
		res = &groups[0]
	} else {
		e := ExpectedOneResult(len(groups))
		err = &e
	}
	return
}

// TemplateGroupsCreate Wrapper for templategroup.create
// https://www.zabbix.com/documentation/6.2/en/manual/api/reference/templategroup/create
func (api *API) TemplateGroupsCreate(TemplateGroups TemplateGroups) (err error) {
	response, err := api.CallWithError("templategroup.create", TemplateGroups)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	groupids := result["groupids"].([]interface{})
	for i, id := range groupids {
		TemplateGroups[i].GroupID = id.(string)
	}
	return
}

// TemplateGroupsUpdate Wrapper for templategroup.update
// https://www.zabbix.com/documentation/6.2/en/manual/api/reference/templategroup/update
func (api *API) TemplateGroupsUpdate(TemplateGroups TemplateGroups) (err error) {
	_, err = api.CallWithError("templategroup.update", TemplateGroups)
	return
}

// TemplateGroupsDelete Wrapper for templategroup.delete
// Cleans GroupId in all TemplateGroups elements if call succeed.
// https://www.zabbix.com/documentation/6.2/en/manual/api/reference/templategroup/delete
func (api *API) TemplateGroupsDelete(TemplateGroups TemplateGroups) (err error) {
	ids := make([]string, len(TemplateGroups))
	for i, group := range TemplateGroups {
		ids[i] = group.GroupID
	}

	err = api.TemplateGroupsDeleteByIds(ids)
	if err == nil {
		for i := range TemplateGroups {
			TemplateGroups[i].GroupID = ""
		}
	}
	return
}

// TemplateGroupsDeleteByIds Wrapper for templategroup.delete
// https://www.zabbix.com/documentation/6.2/en/manual/api/reference/templategroup/delete
func (api *API) TemplateGroupsDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("templategroup.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	groupids := result["groupids"].([]interface{})
	if len(ids) != len(groupids) {
		err = &ExpectedMore{len(ids), len(groupids)}
	}
	return
}
