package run

import "flag"

type flagSet []*flag.Flag

func (f flagSet) makeSureLengthIsSame() {
	max := 0
	for _, v := range f {
		if len(v.Name) > max {
			max = len(v.Name)
		}
	}
	for _, v := range f {
		if len(v.Name) == max {
			continue
		}
		t := max - len(v.Name)
		for i := 0; i < t; i++ {
			v.Name += " "
		}
	}
	max = 0
	for _, v := range f {
		if len(v.DefValue) > max {
			max = len(v.DefValue)
		}
	}
	for _, v := range f {
		if len(v.DefValue) == max {
			continue
		}
		t := max - len(v.DefValue)
		for i := 0; i < t; i++ {
			v.DefValue += " "
		}
	}
}
