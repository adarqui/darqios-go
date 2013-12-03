package main

func Task_State(M *Main, S *State, TD *Task_Data) (bool) {
	Debug("Task_State:Entered\n")

	mon := MON_Gen_State(S)
	M.M<-mon

	return false
}
