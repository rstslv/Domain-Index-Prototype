package main

import (
	"fmt"
	"time"
)

const sz = 10
const arrsz = 3

type domainT int

//domain
const (
	INT  domainT = iota
	CHAR domainT = iota
	BOOL domainT = iota
)

type attribute struct {
	domain domainT
	id     int
}

type value struct {
	number int
	name   string
	time   time.Time
}

type columnT struct {
	ckey      int
	attr      int
	value     int
	isDeleted bool
}

type familycolumnT struct {
	dkey    int
	columns [sz]columnT
	amount  int
	//size int
	//bytef *zipped_columns
}

type supercolumnT struct {
	count int
	fc    [sz]familycolumnT
}

type table struct {
	dkey   int
	ckey   int
	attrs  int
	atr    [sz]attribute
	name   value
	amount int
}

var (
	tab   table
	attrs domainT
)

func show(tab table, sc supercolumnT) {
	for j := 0; j < tab.attrs; j++ {
		fmt.Printf("|%6d|", tab.atr[j].id)
	}
	fmt.Println()
	for i := 0; i < sc.count; i++ {
		for j := 0; j < sc.fc[i].amount; j += tab.attrs - 2 {
			var values [sz]int
			for k := 0; k < tab.attrs-2; k++ {
				values[sc.fc[i].columns[j+k].attr] = sc.fc[i].columns[j+k].value
				values[tab.dkey] = sc.fc[i].dkey
				values[tab.ckey] = sc.fc[i].columns[k].ckey
			}
			if sc.fc[i].columns[j].isDeleted == false {
				for v := 0; v < tab.attrs; v++ {
					switch tab.atr[v].domain {
					case INT:
						fmt.Printf("|%6d|", values[v])
						break
					case CHAR:
						fmt.Printf("|%6d|", values[v]) //control character
						break
					case BOOL:
						if values[v] != 0 {
							fmt.Printf("|%6t|\n", true)
						} else {
							fmt.Printf("|%6t|\n", false)
						}
						break
					}
				}
			}
			fmt.Println()
		}
	}
}

func ifconsists(dk int, ck int, col supercolumnT) int {
	for i := 0; i < col.count; i++ {
		if col.fc[i].dkey == dk {
			for j := 0; j < col.fc[i].amount; j++ {
				if col.fc[i].columns[j].ckey == ck {
					return i
				}
			}
		}
	}
	return -1
}

func insert(tab table, values [arrsz]int, sc *supercolumnT) {
	if !(tab.amount+1 != sz) {
		fmt.Println("assert(tab.amount+1 != sz)")
		return
	}
	tab.amount++
	if ifconsists(values[tab.dkey], values[tab.ckey], *sc) == -1 {
		fcinit(tab.dkey, tab.ckey, sc, values)
		return
	} else if !(ifconsists(values[tab.dkey], values[tab.ckey], *sc) > 0) {
		fmt.Println("assert(ifconsists(values[tab.dkey], values[tab.ckey], *sc) > 0)")
		return
	}
}
func fcinit(dkey int, ckey int, sc *supercolumnT, values [arrsz]int) {
	c := (*sc).count
	(*sc).fc[c].dkey = values[dkey]
	(*sc).fc[c].amount = 0
	for j := 0; j < tab.attrs; j++ {
		if j != ckey && j != dkey {
			l := (*sc).fc[c].amount
			(*sc).fc[c].columns[l].ckey = values[ckey]
			(*sc).fc[c].columns[l].attr = j
			(*sc).fc[c].columns[l].value = values[j]
			(*sc).fc[c].columns[l].isDeleted = false
			(*sc).fc[c].amount++
		}
	}
	(*sc).count++
	return
}
func fcinsert(dkey int, ckey int, tab table, values [arrsz]int, sc *supercolumnT, fcindex int) {
	c := fcindex
	(*sc).fc[c].dkey = values[dkey]
	(*sc).fc[c].amount++
	for j := 0; j < tab.attrs; j++ {
		if j != ckey && j != dkey {
			l := (*sc).fc[c].amount
			(*sc).fc[c].columns[l].ckey = values[ckey]
			(*sc).fc[c].columns[l].attr = j
			(*sc).fc[c].columns[l].value = values[j]
			(*sc).fc[c].columns[l].isDeleted = false
			(*sc).fc[c].amount++
		}
	}
	(*sc).count++
	return
}
func createValues(a int, b int32, c bool, arr *[arrsz]int) { //rune = int32?? arr *int
	arr[0] = int(a)
	arr[1] = int(b)
	if c {
		arr[2] = 1
	} else {
		arr[2] = 0
	}
}
func update(tab table, values [arrsz]int, sc supercolumnT) {
	//func update(tab table, values []int, sc& supercolumnT) {
	iscons := ifconsists(values[tab.dkey], values[tab.ckey], sc)
	if iscons != -1 {
		for j := 0; j < sc.fc[iscons].amount; j++ {
			if sc.fc[iscons].columns[j].ckey == tab.ckey {
				sc.fc[iscons].columns[j].isDeleted = true
			}
		}
		fcinsert(tab.dkey, tab.ckey, tab, values, &sc, iscons)
		return
	}
	fcinit(tab.dkey, tab.ckey, &sc, values)
	return
}
func dbinit(dkey int, ckey int, name string, attrs []string, amount int) {
	tab.dkey = dkey
	tab.ckey = ckey
	tab.attrs = amount
	tab.name.name = name
	tab.name.time = time.Now()
	tab.name.number = 0
	tab.amount = 0
	//var atr [sz]attribute
	for i := 0; i < tab.attrs; i++ {
		tab.atr[i].id = i
		if attrs[i] == "CHAR" {
			tab.atr[i].domain = CHAR
		} else if attrs[i] == "BOOL" {
			tab.atr[i].domain = BOOL
		} else {
			tab.atr[i].domain = INT
		}
	}
}
func main() {
	var (
		sc  supercolumnT
		arr [arrsz]int //[3]
		//updArr  []int 	//[3]
		domains = []string{"INT", "CHAR", "BOOL"}
	)
	sc.count = 0
	dbinit(0, 1, "rnd", domains, 3)
	createValues(10, 'w', true, &arr)
	insert(tab, arr, &sc)
	createValues(18, 'm', false, &arr)
	insert(tab, arr, &sc)
	createValues(19, 'm', true, &arr)
	insert(tab, arr, &sc)
	update(tab, arr, sc)
	show(tab, sc)
}
