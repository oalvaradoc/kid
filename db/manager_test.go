package db

import (
	"git.multiverse.io/eventkit/kit/handler/config"
	"github.com/beego/beego/v2/adapter/orm"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

type Test struct {
	ID   int    `xorm:"not null pk INT(64) 'id'" orm:"column(id);pk" description:"ID number"`
	Name string `xorm:"VARCHAR(64) 'name'" orm:"column(name);size(100)" description:"test name"`
}

type DemoTable struct {
	ID        int    `xorm:"not null pk INT(64) 'id'" orm:"column(id);pk" description:"ID number"`
	FirstName string `xorm:"VARCHAR(64) 'first_name'" orm:"column(first_name);size(45)" description:"first name"`
	LastName  string `xorm:"VARCHAR(64) 'last_name'" orm:"column(last_name);size(45)" description:"last name"`
}

type Demo1Table struct {
	ID        int    `xorm:"not null pk INT(64) 'id'" orm:"column(id);pk" description:"ID number"`
	ClassName string `xorm:"VARCHAR(64) 'class_name'" orm:"column(class_name);size(45)" description:"class name"`
}

type Demo2Table struct {
	ID        int    `orm:"column(id);pk" description:"ID number"`
	ClassName string `orm:"column(class_name);size(45)" description:"class name"`
}

func TestMain(t *testing.M) {
	orm.RegisterModel(new(Test))
	orm.RegisterModel(new(DemoTable))
	orm.RegisterModel(new(Demo1Table))
	orm.RegisterModel(new(Demo2Table))
	t.Run()
}

func TestXormDBManagement(t *testing.T) {
	t.Skip()
	dbConfigs := map[string]config.Db{
		"su0001-db1": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0001",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
		"su0001-db2": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0002",
			Topics:   []string{"DAS00004", "DAS00005", "DAS00006"},
			Default:  false,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo1",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
	}
	InitDBManagerForXorm(dbConfigs)

	// select from db1
	eg, err := GetXormEngine("su0001", "DAS000011")
	if nil != err {
		t.Errorf("Failed to get ormer, error:%++v", err)
		return
	}
	session := eg.NewSession()
	td := &Test{ID: 0}
	session.ID(td.ID)
	_, rerr := session.Get(td)
	if nil != rerr {
		t.Errorf("Failed to read record from DB, error:%++v", rerr)
		return
	}

	t.Logf("[db1]The result of Test record:%++v", td)

	dmt := &DemoTable{ID: 0}
	session.ID(dmt.ID)
	_, rerr = session.Get(dmt)
	if nil != rerr {
		t.Errorf("Failed to read record from DB, error:%++v", rerr)
		return
	}

	t.Logf("[db1]The result of DemoTable record:%++v", dmt)

	// select from db2
	eg, err = GetXormEngine("su0002", "DAS00004")
	if nil != err {
		t.Errorf("Failed to get ormer, error:%++v", err)
		return
	}
	session = eg.NewSession()
	td = &Test{ID: 0}
	session.ID(td.ID)
	_, rerr = session.Get(td)
	if nil != rerr {
		t.Errorf("[db2]Failed to read Test record from DB, error:%++v", rerr)
		return
	}

	t.Logf("[db2]The result of Test record:%++v", td)

	d1t := &Demo1Table{ID: 0}
	session.ID(d1t.ID)

	_, rerr = session.Get(d1t)
	if nil != rerr {
		t.Errorf("[db2]Failed to read Demo1Table record from DB, error:%++v", rerr)
		return
	}
	t.Logf("[db2]The result of Demo1Table record:%++v", d1t)

	dbConfigs = map[string]config.Db{
		"su0001-db1": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0001",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
		"su0001-db2": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0002",
			Topics:   []string{"DAS00004", "DAS00005"},
			Default:  false,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo1",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
	}

	if e := Rotate(dbConfigs); nil != e {
		t.Errorf("Failed to rotate db configs, error:%++v", e)
		return
	}

	t.Log("End of beego DB management roate")
}

func TestBeegoDBManagement(t *testing.T) {
	t.Skip()
	dbConfigs := map[string]config.Db{
		"su0001-db1": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0001",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
		"su0001-db2": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0002",
			Topics:   []string{"DAS00004", "DAS00005", "DAS00006"},
			Default:  false,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo1",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
	}
	InitDBManagerForBeegoOrmer(dbConfigs)

	// select from db1
	dbOrmer, err := GetBeegoOrmer("su0001", "DAS000011")
	if nil != err {
		t.Errorf("Failed to get ormer, error:%++v", err)
		return
	}
	td := &Test{ID: 0}
	rerr := dbOrmer.Read(td)
	if nil != rerr {
		t.Errorf("Failed to read record from DB, error:%++v", rerr)
		return
	}

	t.Logf("[db1]The result of Test record:%++v", td)

	dmt := &DemoTable{ID: 0}
	rerr = dbOrmer.Read(dmt)
	if nil != rerr {
		t.Errorf("Failed to read record from DB, error:%++v", rerr)
		return
	}

	t.Logf("[db1]The result of DemoTable record:%++v", dmt)

	// select from db2
	dbOrmer, err = GetBeegoOrmer("su0002", "DAS00004")
	if nil != err {
		t.Errorf("Failed to get ormer, error:%++v", err)
		return
	}
	td = &Test{ID: 0}
	rerr = dbOrmer.Read(td)
	if nil != rerr {
		t.Errorf("[db2]Failed to read Test record from DB, error:%++v", rerr)
		return
	}

	t.Logf("[db2]The result of Test record:%++v", td)

	d1t := &Demo1Table{ID: 0}
	rerr = dbOrmer.Read(d1t)
	if nil != rerr {
		t.Errorf("[db2]Failed to read Demo1Table record from DB, error:%++v", rerr)
		return
	}
	t.Logf("[db2]The result of Demo1Table record:%++v", d1t)

	dbConfigs = map[string]config.Db{
		"su0001-db1": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0001",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
		"su0001-db2": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0002",
			Topics:   []string{"DAS00004", "DAS00005"},
			Default:  false,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo1",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
	}

	if e := Rotate(dbConfigs); nil != e {
		t.Errorf("Failed to rotate db configs, error:%++v", e)
		return
	}

	// select from db2
	dbOrmer, err = GetBeegoOrmer("su0002", "DAS00004")
	if nil != err {
		t.Errorf("Failed to get ormer, error:%++v", err)
		return
	}
	td = &Test{ID: 0}
	rerr = dbOrmer.Read(td)
	if nil != rerr {
		t.Errorf("[db2]Failed to read Test record from DB, error:%++v", rerr)
		return
	}

	t.Logf("[db2]The result of Test record:%++v", td)

	d2t := &Demo1Table{ID: 0}
	rerr = dbOrmer.Read(d2t)
	if nil != rerr {
		t.Errorf("[db2]Failed to read Demo1Table record from DB, error:%++v", rerr)
		return
	}
	t.Logf("[db2]The result of Demo2Table record:%++v", d2t)

	dbConfigs = map[string]config.Db{
		"su0001-db1": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0001",
			Topics:   []string{},
			Default:  true,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
		"su0001-db2": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0002",
			Topics:   []string{"DAS00004", "DAS00005"},
			Default:  false,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo1",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
		"su0003-db1": config.Db{
			Name:     "",
			Type:     "mysql",
			Su:       "su0003",
			Topics:   []string{"DAS00006", "DAS00007"},
			Default:  false,
			Addr:     "127.0.0.1:3306",
			User:     "root",
			Password: "123456",
			Database: "demo2",
			Params:   "?charset=utf8&loc=Local",
			Debug:    true,
			Pool: struct {
				MaxIdleConns int `json:"maxIdleConns"`
				MaxOpenConns int `json:"maxOpenConns"`
				MaxIdleTime  int `json:"maxIdleTime"`
				MaxLifeValue int `json:"maxLifeValue"`
			}(struct {
				MaxIdleConns      int
				MaxOpenConns      int
				MaxIdleTime 	  int
				MaxLifeValue      int
			}{
				MaxIdleConns:      30,
				MaxOpenConns:      30,
				MaxIdleTime:       30,
				MaxLifeValue:      540,
			}),
		},
	}

	if e := Rotate(dbConfigs); nil != e {
		t.Errorf("Failed to rotate db configs, error:%++v", e)
		return
	}

	// select from db2
	dbOrmer, err = GetBeegoOrmer("su0003", "DAS00007")
	if nil != err {
		t.Errorf("Failed to get ormer, error:%++v", err)
		return
	}
	td = &Test{ID: 0}
	rerr = dbOrmer.Read(td)
	if nil != rerr {
		t.Errorf("[db2]Failed to read Test record from DB, error:%++v", rerr)
		return
	}

	t.Logf("[db2]The result of Test record:%++v", td)

	d3t := &Demo2Table{ID: 0}
	rerr = dbOrmer.Read(d3t)
	if nil != rerr {
		t.Errorf("[db2]Failed to read Demo1Table record from DB, error:%++v", rerr)
		return
	}
	t.Logf("[db3]The result of Demo2Table record:%++v", d3t)
	t.Log("End of beego DB management roate")
}
