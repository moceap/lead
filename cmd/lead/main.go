package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/calmh/lead"
)

var (
	brightness optInt
	color      rgb
	discover   = ""
	controller = ""
)

func main() {
	flag.StringVar(&discover, "discover", "", "Perform discovery on `network` (i.e., 172.16.32.0/24)")
	flag.StringVar(&controller, "controller", "", "Connect to controller at `address` (i.e., 172.16.32.185:8899)")
	flag.Var(&brightness, "brightness", "Set brightness to `N` (0..63)")
	flag.Var(&color, "color", "Set color to `R,G,B` (0..255)")
	flag.Usage = func() {
		fmt.Println("Usage:")
		fmt.Println("  lead -discover <network> [-brightness N] [-color R,G,B]")
		fmt.Println("  lead -controller <address> [-brightness N] [-color R,G,B]")
		fmt.Println("")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	if controller == "" && discover == "" {
		fmt.Println("Need one of -controller or -discover options")
		flag.Usage()
		os.Exit(2)
	}

	var cs []*lead.Controller
	if controller != "" {
		cs = append(cs, lead.NewController(controller))
	}
	if discover != "" {
		tcs, err := lead.Discover(discover)
		if err != nil {
			fmt.Println("Discovering controllers:", err)
			os.Exit(1)
		}
		fmt.Printf("Discovered %d controllers\n", len(tcs))
		for _, c := range tcs {
			fmt.Printf("  %s (%s, %s)\n", c.Address(), c.Model(), c.Serial())
		}
		cs = append(cs, tcs...)
	}

	for _, c := range cs {
		if brightness.isSet {
			if err := c.SetBrightness(brightness.val); err != nil {
				fmt.Printf("Setting brightness on %s: %v\n", c.Address(), err)
			}
		}
		if color.isSet {
			if err := c.SetRGB(color.red, color.green, color.blue); err != nil {
				fmt.Printf("Setting RGB on %s: %v\n", c.Address(), err)
			}
		}
		c.Close()
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
