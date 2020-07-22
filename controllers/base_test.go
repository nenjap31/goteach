package controllers

var (
	userJSON            = `{"username":"goteach","password":"goteach"}`
	userJSONfail            = `{"username":"goteach","password":"goteach123"}`
	userJSONfail2            = `{"username":"goteach123","password":"goteach123"}`
	UserCount = []map[string]interface{}{
		{"count(*)": 1},
	}
	userResp2           = []map[string]interface{}{
		{"username": "goteach","password":"$2a$08$g41oY514c5JUJnU6XVEjVu5CQTTWMgnI0Gv6kB0gnjISwHUr/TF22.","is_active": true,"role_id":1},
	}
	userResp3           = []map[string]interface{}{
		{"username": "goteach","password":"goteach123","is_active": true,"role_id":1},
	}
	userResp4           = []map[string]interface{}{
		{"username": "goteach","password":"goteach123","is_active": false,"role_id":1},
	}
	userRespfail2           = []map[string]interface{}{}
	userRespfail           = []map[string]interface{}{
		{},
	}
	UserResp = []map[string]interface{}{{
		"id": 1,
		"created_at": "2019-11-19T10:44:50+07:00",
		"updated_at": "2019-11-19T10:44:50+07:00",
		"deleted_at": "",
		"name": "goteach",
		"username": "goteach",
		"password": "$2a$08$dDHn8gOuiW7t4k/OBnecQexO3mUZ1/1/GSau9AWD7aDbuOfKz.iAq",
		"email": "goteach@xxx.xxx",
		"role_id": "1",
		"is_active": true,
		"role":"",
	},
	}

	userJSONAdd             = `{"name":"goteach test","username":"goteach2","password":"lusca2","email":"goteach2@alterra.id","is_active":true,"role_id":1}`
	userJSONAddFail             = `{"name":"goteach test fail","username":"goteach3","password":"goteach3","email":"","is_active":true,"role_id":1}`
	userJSONupdate             = `{"name":"goteach test updated","username":"goteach2update","email":"goteach_update@alterra.id","is_active":false,"role_id":1}`
	userJSONupdate2             = `{"name":"goteach test updated2","username":"goteach2update","email":"goteach_update@alterra.id","password":"123456789","is_active":false,"role_id":1}`
	userJSONupdate3             = `{"name":"goteach test updated21","password":"123456789"}`
	userJSONupdatefail             = `{"name":"goteach test updated","username":"goteach2update","email":"","is_active":false,"role_id":1}`
	getprofileresp = map[string]interface{}{
		"id": 0,
		"created_at": "0001-01-01T00:00:00Z",
		"updated_at": "0001-01-01T00:00:00Z",
		"deleted_at": "",
		"name": "",
		"username": "",
		"email": "",
		"role": "",
		"role_id": 0,
		"is_active": false,
		"permissions": "",
	}
	roleJSON            = `{"name":"admin2","is_admin":true}`
	roleJSONu            = `{"name":"admin4","is_admin":true}`
	permissionresp = []map[string]interface{}{
		{"resources":"user","permission": "read"},
	}
	Roleresp           = []map[string]interface{}{
		{"id":"1","name": "admin","permissions":permissionresp,"is_admin": true},
	}
	Roleresp2          = []map[string]interface{}{
		{"id":"1","name": "admin3","permissions":permissionresp,"is_admin": true},
	}
	Roleresp3          = []map[string]interface{}{
		{"id":"1","name": "admin4","permissions":permissionresp,"is_admin": true},
	}
	Rolepermission    =[]map[string]interface{}{
		{"role_id":1,"permission_id": 1},}


	Rolerespfail           = []map[string]interface{}{}
	Rolerespfail2          = []map[string]interface{}{
		{},
	}
)