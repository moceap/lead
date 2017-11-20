package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

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
	addr := kingpin.Flag("addr", "Connect to controller at ADDRESS (i.e., 172.16.32.185:8899)").PlaceHolder("ADDRESS").String()
	file := kingpin.Flag("file", "List of controllers, one per line, like the output of `discover`").PlaceHolder("FILE").ExistingFile()

	cmdDiscover := kingpin.Command("discover", "Discover controllers")
	argNetwork := cmdDiscover.Arg("network", "Network (i.e., 172.16.32.0/24) to probe").Required().String()

	cmdOn := kingpin.Command("on", "Turn on")
	cmdOff := kingpin.Command("off", "Turn on")

	cmdBrightness := kingpin.Command("brightness", "Set brightness (0-63)")
	brightness := cmdBrightness.Arg("brightness", "Brightness value").Int()

	cmdColor := kingpin.Command("color", "Set color (0-255, three octets RGB)")
	color := cmdColor.Arg("color", "Color value").Ints()

	cmd := kingpin.Parse()

	if cmd == cmdDiscover.FullCommand() {
		tcs, err := lead.Discover(*argNetwork)
		if err != nil {
			fmt.Println("Discovering controllers:", err)
			os.Exit(1)
		}
		for _, tc := range tcs {
			fmt.Printf("%s,%s,%s\n", tc.Address(), tc.Model(), tc.Serial())
		}
		os.Exit(0)
	}

	if *addr == "" && *file == "" {
		fmt.Println("Need one of --addr or --file options")
		flag.Usage()
		os.Exit(2)
	}

	var cs []*lead.Controller
	if *addr != "" {
		cs = append(cs, lead.NewController(*addr))
	}

	if *file != "" {
		fd, err := os.Open(*file)
		if err != nil {
			fmt.Println("Reading controllers:", err)
			os.Exit(1)
		}
		cr := csv.NewReader(fd)
		recs, err := cr.ReadAll()
		if err != nil {
			fmt.Println("Reading controllers:", err)
			os.Exit(1)
		}
		for _, rec := range recs {
			cs = append(cs, lead.NewController(rec[0]))
		}
	}

	for _, c := range cs {
		switch cmd {
		case cmdBrightness.FullCommand():
			if err := c.SetBrightness(*brightness); err != nil {
				fmt.Printf("Setting brightness on %s: %v\n", c.Address(), err)
			}

		case cmdColor.FullCommand():
			if err := c.SetRGB((*color)[0], (*color)[1], (*color)[2]); err != nil {
				fmt.Printf("Setting RGB on %s: %v\n", c.Address(), err)
			}

		case cmdOn.FullCommand():
			if err := c.SetOn(true); err != nil {
				fmt.Printf("Turning on %s: %v\n", c.Address(), err)
			}

		case cmdOff.FullCommand():
			if err := c.SetOn(false); err != nil {
				fmt.Printf("Turning off %s: %v\n", c.Address(), err)
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
