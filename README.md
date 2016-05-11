# revel_gorp

## conf/db.conf - configurations
```
connection_pool = 5
encoding = utf8
driver = mysql

[dev]
user = username
pass = password
host = 127.0.0.1
db   = dbname

[prod]

```

## app/init.go
```
func init() {
  //...
  revel.OnAppStart(InitDB)

  //...
}

func InitDB() {
	var err error
	DBMap, err = revel_gorp.InitDatabase()
	if (err != nil) {
		revel.ERROR.Printf("LoadDatabase error: %s", err.Error())
	}

  // add models
	DBMap.AddTableWithName(models.User{}, "users").SetKeys(true, "Id")

  // create tables if not exists
	err = DBMap.CreateTablesIfNotExists()
	if err != nil {
		revel.ERROR.Printf("Create tables error: %s", err.Error())
	}
}
```

Feel free to let me know on your feedback.
Give me a star if you are using this :)
