package main

import (
	"errors"
	"log"
	"fmt"
	"time"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type Mongo struct {
	Ses *mgo.Session
}

func MG_Init() (*Mongo) {

	MGO := new(Mongo)

	return MGO
}

func (M *Main) MG_Init() {
	M.Mongo = MG_Init()
	M.MG_Setup()
	M.MG_Setup_Schema()
}

func (M *Main) MG_Setup() {

	connect_string := fmt.Sprintf("mongodb://"+M.Startup_Config.Mongo.Addr)

	ses, err := mgo.Dial(connect_string)
	if err != nil {
		log.Fatal("dgo_mg:Setup:Open:",err)
	}

	ses.SetSafe(&mgo.Safe{})
	M.Mongo.Ses = ses

	Debug("mongo:MG_Setup:Ses=%q\n", M.Mongo.Ses)
}

func (M *Main) MG_Setup_Schema() {
	/*
	 * Setup index/mongo schema
	 */

	/* db.accounts. */
	accounts := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")
	accounts_Index := mgo.Index{
		Key: []string{"hash"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	}

	err := accounts.EnsureIndex(accounts_Index)
	if err != nil {
		log.Fatal("Setup_Schema:Accounts:EnsureIndex", err)
	}

	/* db.policies. */
	policies := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("policies")
	policies_Index := mgo.Index{
		Key: []string{"name"},
		Unique: true,
		DropDups: true,
		Background: true,
		Sparse: true,
	}

	err = policies.EnsureIndex(policies_Index)
	if err != nil {
		log.Fatal("Setup_Schema:Policies:EnsureIndex", err)
	}

	state := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("state")
	state_Index := mgo.Index{
		Key: []string{"host", "ts"},
		Unique: false,
		DropDups : false,
		Background: true,
		Sparse: true,
	}

	err = state.EnsureIndex(state_Index)
	if err != nil {
		log.Fatal("Setup_Schema:State:EnsureIndex", err)
	}

	alerts := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("alerts")
	alerts_Index := mgo.Index{
		Key: []string{"host", "ts"},
		Unique : false,
		DropDups : false,
		Background : true,
		Sparse : true,
	}

	err = alerts.EnsureIndex(alerts_Index)
	if err != nil {
		log.Fatal("Setup_Schema:Alerts:EnsureIndex", err)
	}
}


func (M *Main) MG_Lookup_Account_By_Hash(Hash string) (*Account, error) {
	Debug("HASH:%v\n", Hash)
	account := new(Account)
	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")
	err := c.Find(bson.M{"hash":Hash}).One(account)
	if err != nil {
		Debug("MG_Lookup_Account_By_Hash:Err:%q\n",err)
		return nil, err
	}

	return account, nil
}


func (M *Main) MG_Update_Account_Last_Seen(Account *Account) (bool, error) {
	t := time.Now()
	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")
	/* Update IP and last timestamp */
	err := c.Update(bson.M{"hash":Account.Hash}, bson.M{"$set":bson.M{"last":t}})
	if err != nil {
		Debug("MG_Update_Account_Last_Seen:Err:%q\n",err)
		return false, err
	}
	// FIXME ^ that should only be one query.. having trouble using $set with two values wtf?
	_ = c.Update(bson.M{"hash":Account.Hash}, bson.M{"$set":bson.M{"ip":Account.Ip}})
	return true, nil
}



func (M *Main) MG_Insert_Account(Hash string, Host string, Groups []string, Status bool) (*Account, error) {
	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")

	account := new(Account)
	account.Hash = Hash
	account.Host = Host
	account.Groups = Groups
	account.Status = Status
	err := c.Insert(account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (M *Main) MG_Insert_Task(A *Account, T *Task) {
	alert_entry := new(Alert_Entry)
	alert_entry.Host = A.Host
	alert_entry.Ip = A.Ip
	alert_entry.Ts = time.Now()
	alert_entry.Task = T

	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("alerts")
	err := c.Insert(alert_entry)
	if err != nil {
		Debug("MG_Insert_Task:Err:%q\n",err)
	}
}


func (M *Main) MG_Update_Account_State(A *Account, S *State_Report) (bool,error) {

	S.Host = A.Host

	S.STATE_Report_Sanitize()

	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")
	err := c.Update(bson.M{"hash":A.Hash},bson.M{"$set":bson.M{"state":S}})
	if err != nil {
		Debug("MG_Update_Account_State:accounts:Err:%q\n",err)
		return false,err
	}

	c = M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("state")
	err = c.Insert(S)
	if err != nil {
		Debug("MG_Update_Account_State:state:Err:%q\n",err)
		return false, err
	}

	return true,nil
}


func (M *Main) MG_Accounts_Find_Missing(Tm time.Time) ([]Account, error) {
	time_diff := Tm.Add(-60 * time.Second)
	accounts := []Account{}
	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")

	v := make([]bson.M,5)
	v[0] = bson.M{"last":bson.M{"$lt":time_diff}}
	v[1] = bson.M{"ignore":bson.M{"$ne":true}}
	err := c.Find(bson.M{"$and":v}).All(&accounts)

	if err != nil {
		return nil,err
	}
	return accounts,err
}

func (M *Main) MG_Accounts_Find_Ignore() ([]Account, error) {
	return M.MG_Accounts_Find_Field(bson.M{"ignore":true})
}

func (M *Main) MG_Accounts_Find_Enable() ([]Account, error) {
	return M.MG_Accounts_Find_Field(bson.M{"status":true})
}

func (M *Main) MG_Accounts_Find_Disable() ([]Account, error) {
	return M.MG_Accounts_Find_Field(bson.M{"status":false})
}

func (M *Main) MG_Accounts_Find_Field(Field interface{}) ([]Account, error) {
	accounts := []Account{}
	c:=M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")
//	err := c.Find(bson.M{"ignore":true}).All(&accounts)
	err := c.Find(Field).All(&accounts)

	if err != nil {
		return nil, err
	}

	return accounts,err
}


func (M *Main) MG_Query(Type string, Hosts []string, Groups []string, Ts_Start string, Ts_End string, Limit string, Filter string) (interface{}, error) {

	/*
	 * Generic query routine
	 *
	 * FIXME - support groups
	 */
	var accounts []State_Report

	ts_start, err := XTIME_Get(Ts_Start)
	if err != nil {
		return nil, err
	}

	ts_end, err := XTIME_Get(Ts_End)
	if err != nil {
		return nil, err
	}

	limit, err := LIMIT_Get(Limit)
	if err != nil {
		return nil, err
	}

//	Debug("%q %q\n", ts_start, ts_end)

//	bson_hosts := bson.M{"host":bson.M{"$in":Hosts}}
//	bson_time := bson.M{"$gt":ts_start, "$lt":ts_end}
//	q := []bson.M{}
	q := bson.M{}

	if len(Hosts) != 0 {
		q = bson.M{"host":bson.M{"$in":Hosts},"ts":bson.M{"$gt":ts_start, "$lt":ts_end}}
	} else {
		q = bson.M{"ts":bson.M{"$gt":ts_start, "$lt":ts_end}}
	}

//	Debug("%q\n", bson_hosts)

	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("state")
	err = c.Find(q).Limit(limit).All(&accounts)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (M *Main) MG_Accounts(Hosts []string, Groups []string) ([]Account, error) {
	var accounts []Account
	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")

	var err error
	if len(Hosts) != 0 && Hosts != nil {
		err = c.Find(bson.M{"host":bson.M{"$in":Hosts}}).All(&accounts)
	} else if len(Groups) != 0 && Groups != nil {
		err = c.Find(bson.M{"groups":bson.M{"$in":Groups}}).All(&accounts)
	} else {
		err = c.Find(nil).All(&accounts)
	}

	if err != nil {
		return nil, err
	}

	return accounts, nil
}


func (M *Main) MG_Accounts_Ignore(Hosts []string, Groups []string, Truth bool) (bool, error) {

	return M.MG_Set_Field(Hosts,Groups,bson.M{"ignore":Truth})
}


func (M *Main) MG_Accounts_Enable(Hosts []string, Groups []string, Truth bool) (bool, error) {
	return M.MG_Set_Field(Hosts,Groups,bson.M{"status":Truth})
}

func (M *Main) MG_Set_Field(Hosts []string, Groups []string, Fields interface{}) (bool, error) {

	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("accounts")

	var err error

	if len(Hosts) != 0 && Hosts != nil {
		err = c.Update(bson.M{"host":bson.M{"$in":Hosts}}, bson.M{"$set":Fields})
	} else if len(Groups) != 0 && Groups != nil {
		err = c.Update(bson.M{"groups":bson.M{"$in":Groups}}, bson.M{"$set":Fields})
	} else {
		return false, errors.New("specify hosts or groups")
	}

	return false, err
}
