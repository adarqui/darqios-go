package main

func (Policy *Policy) Check_Policy_For_Account(Account *Account) (bool) {
	/*
	 * get rid of regex
	 */

	DebugLn("Check_Policy_For_Account:Entered")
	if Policy.Active == false {
		Debug("ChecK_Policy_For_Account:%s is not active\n", Policy.Name)
		return false
	}

	for _, v := range(Policy.Hosts) {
		if v == Account.Host || len(v) == 0 {
			Debug("Check_Policy_For_Account:Host Match:%s\n", v)
			return true
		}
	}

	for _, v := range(Policy.Groups) {
		for _, w := range(Account.Groups) {
			if w == v {
				Debug("Check_Policy_For_Group:Group Match:%s\n", w)
				return true
			}
		}
	}

	return false
}


func (Policy *Policy) Disable() {
	Policy.Active = false
}

func (Policy *Policy) Sanitize() {
	/*
	 * Clean up a policy to get rid of invalid "numerical values" and such
	 */
	if Policy.Interval <= 0 {
		Policy.Interval = 60
	}

	if len(Policy.Thresholds) > 0 && len(Policy.Thresholds) != 3 {
		Policy.Active = false
	}

	if len(Policy.Mitigate) > 0 && len(Policy.Mitigate) != 3 {
		Policy.Active = false
	}
}
