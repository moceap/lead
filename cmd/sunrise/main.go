package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/calmh/lead"
)

var (
	brightness optInt
	color      rgb
	discover   = ""
	controller = ""
)

func main() {
	argNetwork := kingpin.Arg("network", "Network (i.e., 172.16.32.0/24) to probe").Required().String()
	sleep := kingpin.Arg("sleep", "Sleep time per increase").Default("60s").Duration()
	kingpin.Parse()

	tcs, err := lead.Discover(*argNetwork)
	if err != nil {
		fmt.Println("Discovering controllers:", err)
		os.Exit(1)
	}

	for _, c := range tcs {
		fmt.Println(c, "init")
		for i := 0; i < 5; i++ {
			if err := c.SetOn(true); err != nil {
				fmt.Printf("Turning on %s: %v\n", c.Address(), err)
			}
			time.Sleep(100 * time.Millisecond)
			if err := c.SetBrightness(2); err != nil {
				fmt.Printf("Setting brightness on %s: %v\n", c.Address(), err)
			}
			time.Sleep(100 * time.Millisecond)
			if err := c.SetRGB(255, 192, 32); err != nil {
				fmt.Printf("Setting RGB on %s: %v\n", c.Address(), err)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	for b := 3; b < 33; b++ {
		time.Sleep(*sleep)
		for _, c := range tcs {
			fmt.Println(c, "brightness", b)
			if err := c.SetBrightness(b); err != nil {
				fmt.Printf("Setting brightness on %s: %v\n", c.Address(), err)
			}
		}
	}
}

type rgb struct {
	red, green, blue int
	isSet            bool
}

func (v *rgb) Set(rgb string) error {
	fields := strings.Split(rgb, ",")
	if len(fields) != 3 {
		return fmt.Errorf("cannot parse as R,G,B")
	}

	var err error
	v.red, err = strconv.Atoi(fields[0])
	if err != nil {
		return err
	}
	v.green, err = strconv.Atoi(fields[1])
	if err != nil {
		return err
	}
	v.blue, err = strconv.Atoi(fields[2])
	if err != nil {
		return err
	}

	v.isSet = true
	return nil
}

func (v *rgb) String() string {
	if !v.isSet {
		return ""
	}
	return fmt.Sprintf("%d,%d,%d", v.red, v.green, v.blue)
}

type optInt struct {
	val   int
	isSet bool
}

func (v *optInt) Set(s string) error {
	var err error
	v.val, err = strconv.Atoi(s)
	if err != nil {
		return err
	}

	v.isSet = true
	return nil
}

func (v *optInt) String() string {
	if !v.isSet {
		return ""
	}
	return strconv.Itoa(v.val)
}
