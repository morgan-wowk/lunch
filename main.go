package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type menuItemFields struct {
	WeekTitle    string
	Caterer      string
	Monday       []string
	MondayVeg    interface{}
	Tuesday      []string
	TuesdayVeg   interface{}
	Wednesday    []string
	WednesdayVeg interface{}
	Thursday     []string
	ThursdayVeg  interface{}
	Friday       []string
	FridayVeg    interface{}
}

type menuItem struct {
	Fields menuItemFields
}

type menu struct {
	Limit int
	Items []menuItem
}

func main() {
	var location string

	app := cli.NewApp()
	app.Name = "lunch"
	app.Usage = "Check the lunch for this week."
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "location, l",
			Value:       "fultz",
			Usage:       "Location you want to see the menu for.",
			Destination: &location,
		},
		cli.BoolFlag{
			Name:  "today, t",
			Usage: "Get lunch menu for today only.",
		},
	}
	app.Action = func(c *cli.Context) error {
		fultzURL := "https://cdn.contentful.com/spaces/6qqte9wlq16o/entries?access_token=bab0ec81f61331d6e29f5c0e3164d8d506c5ae6957088607c0125a71124177c7"
		smartparkURL := "https://cdn.contentful.com/spaces/sw4tprcfpvo7/entries?access_token=c04ffd4690c2804f7c772577cdd59065d5c3bfd3e9b06147905c11979a06be3c"

		var url = fultzURL
		if strings.ToLower(location) == "smartpark" {
			url = smartparkURL
		}

		client := http.Client{
			Timeout: time.Second * 5,
		}

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		res, getErr := client.Do(req)

		if getErr != nil {
			log.Fatal(getErr)
		}

		body, readErr := ioutil.ReadAll(res.Body)
		if readErr != nil {
			log.Fatal(readErr)
		}

		jsonObj := menu{}
		jsonErr := json.Unmarshal(body, &jsonObj)
		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		title := color.New(color.Bold, color.FgGreen, color.Underline).PrintlnFunc()
		label := color.New(color.Bold, color.FgMagenta).PrintFunc()
		text := color.New(color.FgCyan).PrintFunc()

		label("Week: ")
		text(jsonObj.Items[0].Fields.WeekTitle)
		text("\n")
		label("Caterer: ")
		text(jsonObj.Items[0].Fields.Caterer)
		text("\n")
		label("Location: ")
		if strings.ToLower(location) == "smartpark" {
			text("Smartpark")
		} else {
			text("Fultz")
		}
		text("\n\n")
		title("Menu")

		monday := jsonObj.Items[0].Fields.Monday
		mondayVeg := jsonObj.Items[0].Fields.MondayVeg
		tuesday := jsonObj.Items[0].Fields.Tuesday
		tuesdayVeg := jsonObj.Items[0].Fields.TuesdayVeg
		wednesday := jsonObj.Items[0].Fields.Wednesday
		wednesdayVeg := jsonObj.Items[0].Fields.WednesdayVeg
		thursday := jsonObj.Items[0].Fields.Thursday
		thursdayVeg := jsonObj.Items[0].Fields.ThursdayVeg
		friday := jsonObj.Items[0].Fields.Friday
		fridayVeg := jsonObj.Items[0].Fields.FridayVeg

		outputWeekday(monday, mondayVeg, "Monday", c.Bool("today"))
		outputWeekday(tuesday, tuesdayVeg, "Tuesday", c.Bool("today"))
		outputWeekday(wednesday, wednesdayVeg, "Wednesday", c.Bool("today"))
		outputWeekday(thursday, thursdayVeg, "Thursday", c.Bool("today"))
		outputWeekday(friday, fridayVeg, "Friday", c.Bool("today"))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func outputWeekday(lunchDay interface{}, lunchDayVeg interface{}, labelText string, todayOnly bool) {
	weekday := time.Now().Weekday().String()
	label := color.New(color.Bold, color.FgMagenta).PrintFunc()
	text := color.New(color.FgCyan).PrintFunc()
	bold := color.New(color.Bold).PrintFunc()

	if !todayOnly || weekday == labelText {
		label(labelText + ":\n\t")
		if reflect.TypeOf(lunchDay).Kind() == reflect.Slice {
			text(strings.Join(lunchDay.([]string), "\n\t"))
		} else {
			text(lunchDay)
		}
		if lunchDayVeg != nil {
			text("\n\t")
			bold("Vegetarian: ")
			if reflect.TypeOf(lunchDayVeg).Kind() == reflect.Slice {
				vegSlice := lunchDayVeg.([]interface{})
				if len(vegSlice) > 1 {
					text("\n\t\t")
				}
				vegOptions := make([]string, len(vegSlice))
				for i := 0; i < len(vegSlice); i++ {
					vegOptions[i] = lunchDayVeg.([]interface{})[i].(string)
				}
				text(strings.Join(vegOptions, "\n\t\t"))
			} else {
				text(lunchDayVeg)
			}
		}
		text("\n")
	}
}
