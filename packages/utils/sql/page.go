package sql

func (db *DCDB) GetPageMenus(pagePrefix, pageName string) (string, error) {
	return db.Single(`SELECT menu FROM "`+pagePrefix+`_pages" WHERE name = ?`, pageName).String()
}

func (db *DCDB) GetPage(tablePrefix, pageName string) (map[string]string, error) {
	return db.OneRow(`SELECT * FROM "`+tablePrefix+`_pages" WHERE name = ?`, pageName).String()
}

func (db *DCDB) GetInterfacePages(tableprefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tableprefix+`_pages" where menu!='0' order by name`, -1)
}

func (db *DCDB) GetInterfaceBlocks(tablePrefix string) ([]map[string]string, error) {
	return db.GetAll(`SELECT * FROM "`+tablePrefix+`_pages" where menu='0' order by name`, -1)
}

func (db *DCDB) GetValueFromPage(tableName string, page string) (string, error) {
	return db.Single(`SELECT value FROM`+tableName+`WHERE name =?`, page).String()
}

func (db *DCDB) CreatePagesTable(id string) error {
	return db.ExecSQL(`CREATE TABLE "` + id + `_pages" (
				"name" varchar(255)  NOT NULL DEFAULT '',
				"value" text  NOT NULL DEFAULT '',
				"menu" varchar(255)  NOT NULL DEFAULT '',
				"conditions" bytea  NOT NULL DEFAULT '',
				"rb_id" bigint NOT NULL DEFAULT '0'
				);
				ALTER TABLE ONLY "` + id + `_pages" ADD CONSTRAINT "` + id + `_pages_pkey" PRIMARY KEY (name);
				`)
}

func (db *DCDB) CreateFirstPagesRecords(id string, sid string) error {
	return db.ExecSQL(`INSERT INTO "`+id+`_pages" (name, value, menu, conditions) VALUES
		(?, ?, ?, ?),
		(?, ?, ?, ?)`,
		`dashboard_default`, `FullScreen(1)

If(StateVal(type_office))
Else:
Title : Basic Apps
Divs: col-md-4
		Divs: panel panel-default elastic
			Divs: panel-body text-center fill-area flexbox-item-grow
				Divs: flexbox-item-grow flex-center
					Divs: pv-lg
					Image("/static/img/apps/money.png", Basic, center-block img-responsive img-circle img-thumbnail thumb96 )
					DivsEnd:
					P(h4,Basic Apps)
					P(text-left,"Election and Assign, Polling, Messenger, Simple Money System")
				DivsEnd:
			DivsEnd:
			Divs: panel-footer
				Divs: clearfix
					Divs: pull-right
						BtnPage(app-basic, Install,'',btn btn-primary lang)
					DivsEnd:
				DivsEnd:
			DivsEnd:
		DivsEnd:
	DivsEnd:
IfEnd:
PageEnd:
`, `menu_default`, sid,

		`government`, `FullScreen(1)

If(StateVal(type_office))
Else:
Title : Basic Apps
Divs: col-md-4
		Divs: panel panel-default elastic
			Divs: panel-body text-center fill-area flexbox-item-grow
				Divs: flexbox-item-grow flex-center
					Divs: pv-lg
					Image("/static/img/apps/money.png", Basic, center-block img-responsive img-circle img-thumbnail thumb96 )
					DivsEnd:
					P(h4,Basic Apps)
					P(text-left,"Election and Assign, Polling, Messenger, Simple Money System")
				DivsEnd:
			DivsEnd:
			Divs: panel-footer
				Divs: clearfix
					Divs: pull-right
						BtnPage(app-basic, Install,'',btn btn-primary lang)
					DivsEnd:
				DivsEnd:
			DivsEnd:
		DivsEnd:
	DivsEnd:
IfEnd:
PageEnd:
`, `government`, sid,
	)
}
