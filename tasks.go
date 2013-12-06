package main

import (
	"fmt"
	"strconv"
)

type Task_Data struct {
	Policy *Policy
	State *State
}

func Task_Compare_Numbers(Thresholds []string, Actual float64, Which byte) (string, string) {

	/* low med high
	 * example: Thresholds[5,10,15], 6, > ... if 6 >= 15, if 6 >= 10, *if 6 > 5
	 * example: Thresholds[5,10,15], 6, < ... if 6 <= 5, *if 6 <= 10
	 * example: Thresholds[5,6,7], 6, > ... if 6 >= 7, *if 6 >= 6
	 * example: Thresholds[5,6,7], 6, < ... if 6 <= 5, *if 6 <= 6
	 * example: Thresholds[20,50,90], 40 >, ... if 40 >= 90, if 40 >= 50, *if 40 >=

	 * example: Thresholds[20,50,90], 40 <, ... if 40 <= 20, *if 40 <= 50
	 */

	Debug("Task_Compare_Numbers:Original Thresholds: %q\n", Thresholds)

	hash := make(map[int]string)

	if len(Thresholds) == 0 || len(Thresholds) > 3{
		return "", NIL_STRING
	}

	T := make([]string,3)

	hash[0] = "high"
	hash[1] = "med"
	hash[2] = "low"

	T[0] = Thresholds[2]
	T[1] = Thresholds[1]
	T[2] = Thresholds[0]

	Debug("Task_Compare_Numbers:New Thresholds: %q\n", T)

	for k,v := range(T) {
		if k > 2 {
			break
		}

		/*
		 * Passing nil to a threshold means we just skip it
		 * Added: ""
		 */
		if v == "nil" || v == "" || len(v) == 0 {
			continue
		}

		v2, err := strconv.ParseFloat(v, 64)
		if err != nil {
			continue
		}

		fmt.Printf("TESTING: threshold=%q actual=%q\n", v2, Actual)

		switch Which {
			case '>' : {
				if Actual >= v2 {
					return hash[k], v
				}
				break;
			}
			case '=' : {
				if Actual == v2 {
					return hash[k], v
				}
				break;
			}
			case '<' : {
				if Actual <= v2 {
					return hash[k], v
				}
				break;
			}
		}
	}
	return NIL_STRING, NIL_STRING
}




func Task_Get_Alert_Level(P *Policy) (string) {

	/*
	arr := strings.Split(P.Idx, ":")
	if len(arr) != 2 {
		return "high"
	}

	return arr[1]
	*/


	if P.Level != "" {
		return P.Level
	}
	return "high"
}
