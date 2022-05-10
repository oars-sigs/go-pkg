package idaas

import "testing"

func TestPermissionInitData(t *testing.T) {
	c := New(cfg)
	err := c.PermissionInitData(&InitPermissionData{
		Group:   "oars_filebase",
		Name:    "oars_filebase_init",
		Version: "0.0.5",
		Resources: []string{
			"namespace_p,云盘项目空间",
			"namespace_a,应用文件空间",
		},
		Actions: []string{
			"namespace_p,create,创建",
			"namespace_a,create,创建",
		},
		Roles: []string{
			"所有者,owner,rr,namespace,*",
			"读写,rw,rr,namespace,*",
			"只读,r,rr,namespace,*",
		},
		Rules: []string{
			"owner,rr,namespace,*,namespace,*,*,a",
			"owner,rr,namespace,*,file,*,*,a",
			"rw,rr,namespace,*,file,*,r,a",
			"rw,rr,namespace,*,file,*,w,a",
			"r,rr,namespace,*,file,*,r,a",
			"djh@gzsunrun.cn,u,namespace/file,*,namespace,test1,owner,rr",
			"djh@gzsunrun.cn,u,namespace/file,*,namespace,test2,rw,rr",
			"djh@gzsunrun.cn,u,namespace/file,*,namespace,test3,r,rr",
		},
	})
	if err != nil {
		t.Error(err)
	}
}

// func TestPermissionEnforce(t *testing.T) {
// 	c := New(cfg)
// 	res, err := c.PermissionEnforce(EnforceParam{
// 		Group:        "oars_filebase",
// 		Resource:     "file",
// 		ResourceName: "test2/1",
// 		Action:       "w",
// 		UserId:       "djh@gzsunrun.cn",
// 		DeptPath:     "1",
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(res)
// }

// func TestPermissionResources(t *testing.T) {
// 	c := New(cfg)
// 	res, err := c.PermissionResources(EnforceParam{
// 		Group:        "oars_filebase",
// 		Resource:     "file",
// 		ResourceName: "test2/*",
// 		Action:       "w",
// 		UserId:       "djh@gzsunrun.cn",
// 		DeptPath:     "1",
// 	})
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	t.Log(res)
// }

func TestPermissionPutRule(t *testing.T) {
	c := New(cfg)
	param := &PermissionRulePutParam{
		Filter: &PermissionRolebindings{
			Group:             "oars_filebase",
			MatchResource:     "namespace/file",
			MatchResourceName: "*",
			Resource:          "namespace",
			ResourceName:      "test4",
			ActionKind:        ActionKindResourceRole,
		},
		Data: []PermissionRolebindings{
			PermissionRolebindings{
				Group:             "oars_filebase",
				User:              "24541fa8-6de3-4536-aca8-9e28ba3bdaac",
				UserKind:          "u",
				MatchResource:     "namespace/file",
				MatchResourceName: "*",
				Resource:          "namespace",
				ResourceName:      "test4",
				ActionKind:        ActionKindResourceRole,
				Action:            "owner",
			},
			PermissionRolebindings{
				Group:             "oars_filebase",
				User:              "30bad16d-de58-4f28-93cf-b08228e08bf3",
				UserKind:          "d",
				MatchResource:     "namespace/file",
				MatchResourceName: "*",
				Resource:          "namespace",
				ResourceName:      "test4",
				ActionKind:        ActionKindResourceRole,
				Action:            "owner",
			},
		},
	}
	err := c.PermissionPutRule(param)
	if err != nil {
		t.Error(err)
	}
}

func TestPermissionListRule(t *testing.T) {
	c := New(cfg)
	res, err := c.PermissionListRule(&PermissionRolebindings{
		Group:             "oars_filebase",
		MatchResource:     "namespace/file",
		MatchResourceName: "*",
		Resource:          "namespace",
		ResourceName:      "test4",
		ActionKind:        ActionKindResourceRole,
	})
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
