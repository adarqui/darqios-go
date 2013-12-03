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
	err := c.Update(bson.M{"hash":Account.Hash}, bson.M{"$set":bson.M{"last":t}})
	if err != nil {
		Debug("MG_Update_Account_Last_Seen:Err:%q\n",err)
		return false, err
	}
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
	alert_entry.Task = T

	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("alerts")
	err := c.Insert(alert_entry)
	if err != nil {
		Debug("MG_Insert_Task:Err:%q\n",err)
	}
}


func (M *Main) MG_Update_Account_State(A *Account, S *State_Report) (bool,error) {

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
