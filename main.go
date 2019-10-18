package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	weekday := time.Now().Weekday().String()
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
		if location == "smartpark" {
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
		bold := color.New(color.Bold).PrintFunc()

		label("Week: ")
		text(jsonObj.Items[0].Fields.WeekTitle)
		text("\n")
		label("Caterer: ")
		text(jsonObj.Items[0].Fields.Caterer)
		text("\n\n")
		title("Menu")
		if !c.Bool("today") || weekday == "Monday" {
			label("Monday:\n\t")
			text(jsonObj.Items[0].Fields.Monday)
			if jsonObj.Items[0].Fields.MondayVeg != nil {
				text("\n\t")
				bold("Vegetarian: ")
				text(jsonObj.Items[0].Fields.MondayVeg)
			}
			text("\n")
		}
		if !c.Bool("today") || weekday == "Tuesday" {
			label("Tuesday:\n\t")
			text(jsonObj.Items[0].Fields.Tuesday)
			if jsonObj.Items[0].Fields.TuesdayVeg != nil {
				text("\n\t")
				bold("Vegetarian: ")
				text(jsonObj.Items[0].Fields.TuesdayVeg)
			}
			text("\n")
		}
		if !c.Bool("today") || weekday == "Wednesday" {
			label("Wednesday:\n\t")
			text(jsonObj.Items[0].Fields.Wednesday)
			if jsonObj.Items[0].Fields.WednesdayVeg != nil {
				text("\n\t")
				bold("Vegetarian: ")
				text(jsonObj.Items[0].Fields.WednesdayVeg)
			}
			text("\n")
		}
		if !c.Bool("today") || weekday == "Thursday" {
			label("Thursday:\n\t")
			text(jsonObj.Items[0].Fields.Thursday)
			if jsonObj.Items[0].Fields.ThursdayVeg != nil {
				text("\n\t")
				bold("Vegetarian: ")
				text(jsonObj.Items[0].Fields.ThursdayVeg)
			}
			text("\n")
		}
		if !c.Bool("today") || weekday == "Friday" {
			label("Friday:\n\t")
			text(jsonObj.Items[0].Fields.Friday)
			if jsonObj.Items[0].Fields.FridayVeg != nil {
				text("\n\t")
				bold("Vegetarian: ")
				text(jsonObj.Items[0].Fields.FridayVeg)
			}
			text("\n")
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
