package main

import (
	/*
	"log"
	*/
//	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)



func (M *Main) MG_Lookup_Policies_Config_Raw(Name string) (*Policies_Config, error) {

	policy_config := new(Policies_Config)
	c := M.Mongo.Ses.DB(M.Startup_Config.Mongo.Db).C("policies")
	err := c.Find(bson.M{"name":Name}).One(policy_config)
	if err != nil {
		return nil, err
	}

	return policy_config, nil
}

func (M *Main) MG_Lookup_Policies_Config(Name string, Account *Account) (*Policies_Config, error) {

	policy_config, err := M.MG_Lookup_Policies_Config_Raw(Name)
	if err != nil {
		return nil, err
	}

	if Account != nil {
		for _, policy := range policy_config.Policies {

			truth := policy.Check_Policy_For_Account(Account)
			if truth != true {
				policy.Disable()
			}
			policy.Sanitize()
		}
	}


	return policy_config, nil
}


func (M *Main) MG_Lookup_Policies_Configs() ([]*Policies_Config, error) {
	// FIXME - add this functionality
	return nil,nil
}
