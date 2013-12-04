package main

import (
	"time"
)


type Account struct {
	Id string
	Ip string
	Hash string
	Host string
	Groups []string
	Status bool
	Last time.Time
	State State_Report
	/*
	 * if true, don't notify us of anything
	 */
	Ignore bool
}

type Alert_Entry struct {
	Host string
	Ip string
	Task *Task
}

type Task struct {
	Type string
	Name string
	Idx string
	Subject string
	Body string
	Time time.Time
	Actual string
	Operator string
	Threshold string
}


type Policy struct {
	Name string
	Desc string
	Hosts []string
	Groups []string
	Interval int
	Idx string
	Level string
	Flags []string
	Params []string
	Thresholds []string
	Mitigate []string
	Active bool
	/*
	 * Runs the policy/task on the client, but just doesn't send an alert back
	 */
	Ignore bool
}

type Policies struct {
	Policy Policy
}

type Lost struct {
	Name string
	Desc string
	Hosts[]string
	Groups []string
	Params []string
	Thresholds []string
	Acitve bool
}

type Missing struct {
}

type Options struct {
}

type Alert struct {
	Scripts []string
}

type Alerts struct {
	Low *Alert
	Med *Alert
	High *Alert
}

type Policies_Config struct {
	Version int
	Name string
	Base string
	Policies []*Policy
	Missing []Missing
	Ignore []string
	Options Options
	Alerts *Alerts
}
