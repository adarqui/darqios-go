package main

func (M *Main) NET_Client_MON() {
	DebugLn("NET_Client_MON:Initialized")

	fns := map[string]func(*Main, *State, *Task_Data) (bool) {
		"Ping" : Task_Ping,
		"Load" : Task_Load,
		"Memory" : Task_Memory,
		"Process" : Task_Process,
		"Disk" : Task_Disk,
		"Network" : Task_Network,
		"Scheduler" : Task_Scheduler,
		"State" : Task_State,
		"File" : Task_File,
		"Pipe" : Task_Pipe,
		"Custom" : Task_Custom,
	}

	state := M.STATE_Init()

	for {
		state.STATE_Sleep()

		if M.Policies_Config == nil {
			continue
		} else if M.Policies_Config.Policies == nil {
			continue
		}

		if (state.Interval_Counter % 5) == 0 {
			state.STATE_Get()
		}

		for _, policy := range M.Policies_Config.Policies {
			truth := state.STATE_Should_Run(policy)
			if truth != true {
				continue
			}

			TD := new(Task_Data)
			TD.Policy = policy

			if fn_ptr,ok := fns[policy.Name]; ok {
				go fn_ptr(M, state, TD)
			}
		}
	}
}
