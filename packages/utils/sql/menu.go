package sql

func (db *DCDB) GetValueFromMenu(pagePrefix, menuName string) (string, error) {
	return db.Single(`SELECT value FROM "`+pagePrefix+`_menu" WHERE name = ?`, menuName).String()
}

func (db *DCDB) GetMenu(pagePrefix, menuName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+pagePrefix+`_menu" WHERE name = ?`, menuName).String()
}

func (db *DCDB) GetAllMenus(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_menu" order by name`, -1)
}

func (db *DCDB) CreateMenuTable(id string) error {
	return db.ExecSQL(`CREATE TABLE "` + id + `_menu" (
				"name" varchar(255)  NOT NULL DEFAULT '',
				"value" text  NOT NULL DEFAULT '',
				"conditions" bytea  NOT NULL DEFAULT '',
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER TABLE ONLY "` + id + `_menu" ADD CONSTRAINT "` + id + `_menu_pkey" PRIMARY KEY (name);
				`)
}

func (db *DCDB) CreateFirstMenuRecord(id string, sid string) error {
	return db.ExecSQL(`INSERT INTO "`+id+`_menu" (name, value, conditions) VALUES
		(?, ?, ?),
		(?, ?, ?)`,
		`menu_default`, `MenuItem(Dashboard, dashboard_default)
 MenuItem(Government dashboard, government)`, sid,
		`government`, `MenuItem(Citizen dashboard, dashboard_default)
MenuItem(Government dashboard, government)
MenuGroup(Admin tools,admin)
MenuItem(Tables,sys-listOfTables)
MenuItem(Smart contracts, sys-contracts)
MenuItem(Interface, sys-interface)
MenuItem(App List, sys-app_catalog)
MenuItem(Export, sys-export_tpl)
MenuItem(Wallet,  sys-edit_wallet)
MenuItem(Languages, sys-languages)
MenuItem(Signatures, sys-signatures)
MenuItem(Gen Keys, sys-gen_keys)
MenuEnd:
MenuBack(Welcome)`, sid)
}
