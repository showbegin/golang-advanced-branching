package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type vehicle interface {
}

type car struct {
	model       string
	make        string
	typeVehicle string
}

type truck struct {
	model       string
	make        string
	typeVehicle string
}

type bike struct {
	model string
	make  string
}

// Values array for the feedback.json file
type Values struct {
	Models []Model `json: "values"`
}

// Model array for the feedback.json file
type Model struct {
	Name     string   `json: "model"`
	Feedback []string `json: "feedback"`
}

type feedbackResult struct {
	feedbackTotal    int
	feedbackPositive int
	feedbackNegative int
	feedbackNeutral  int
}

var vehicleResult map[string]feedbackResult
var inventory []vehicle

type rating float32

const (
	extraPositive rating = 1.2
	positive      rating = 0.6
	negative      rating = -0.6
	initial       rating = 5.0
	extraNegative rating = -1.2
)

func init() {

	inventory = []vehicle{
		bike{"FTR 1200", "Indian"},
		bike{"Iron 1200", "Harley"},
		car{"Sonata", "Hyundai", "Sedan"},
		car{"SantaFe", "Hyundai", "SUV"},
		car{"Civic", "Honda", "Hatchback"},
		car{"A5", "Audi", "Coupe"},
		car{"Mazda6", "Mazda", "Sedan"},
		car{"CRV", "Honda", "SUV"},
		car{"Camry", "Toyota", "Sedan"},
		truck{"F-150", "Ford", "Truck"},
		truck{"RAM1500", "Dodge", "Truck"}}

	vehicleResult = make(map[string]feedbackResult)

}

func main() {

	// Generate ratings for the different vehicles
	fmt.Println(len(vehicleResult))
	generateRating()
	fmt.Println(len(vehicleResult))
	for k, v := range vehicleResult {
		fmt.Sprintf("result for %v: %v", k, v)
	}
	// Print ratings for the different vehicles
	for _, veh := range inventory {
		switch v := veh.(type) {
		case car:
			v.carDetails()
		case bike:
			v.bikeDetails()
		case truck:
			v.truckDetails()
		default:
			fmt.Printf("Are you sure this Vehicle Type %v exists?\n", v)
		}
	}
}

func readJSONFile() Values {
	jsonFile, err := os.Open("feedback.json")

	if err != nil {
		log.Fatal("File not found")
	} else {
		fmt.Println("File found")
	}
	fmt.Println("HOHOHO")
	fi, err := jsonFile.Stat()
	if err != nil {
		log.Fatal("cannot get file info")
	} else {
		fmt.Println("File info: ", fi.Name())
	}
	fmt.Println("len:", fi.Size())
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var content Values
	err = json.Unmarshal(byteValue, &content)
	if err != nil {
		log.Fatal("cannot unmarshal json")
	} else {
		fmt.Println("Unmarshalled json")
	}
	fmt.Println("content:", len(content.Models))
	return content
}

func generateRating() {
	f := readJSONFile()
	fmt.Sprintf("\ncontent:%v", f)
	for _, v := range f.Models {
		var vehResult feedbackResult
		var vehRating rating
		for _, msg := range v.Feedback {
			fmt.Sprint("msg: %v", msg)
			if text := strings.Split(msg, " "); len(text) >= 5 {
				vehRating = 5.0
				vehResult.feedbackTotal++
				for _, word := range text {
					switch s := strings.Trim(strings.ToLower(word), " ,.,!,?,\t,\n,\r"); s {
					case "pleasure", "impressed", "wonderful", "fantastic", "splendid":
						vehRating += extraPositive
					case "help", "helpful", "thanks", "thank you", "happy":
						vehRating += positive
					case "not helpful", "sad", "angry", "improve", "annoy":
						vehRating += negative
					case "pathetic", "bad", "worse", "unfortunately", "agitated", "frustrated":
						vehRating += extraNegative
					}
				}
				switch {
				case vehRating > 8.0:
					vehResult.feedbackPositive++
				case vehRating >= 4.0 && vehRating <= 8.0:
					vehResult.feedbackNeutral++
				case vehRating < 4.0:
					vehResult.feedbackNegative++
				}
			}
		}
		fmt.Sprintf("\nvehResult:%v", vehResult)
		vehicleResult[v.Name] = vehResult
	}
}

func showRating(model string) {
	ratingFound := false
	for m, r := range vehicleResult {
		if m == model {
			fmt.Printf("Total Ratings:%v\tPositive:%v\tNegative:%v\tNeutral:%v", r.feedbackTotal, r.feedbackPositive, r.feedbackNegative, r.feedbackNeutral)
			ratingFound = true
		}
	}
	if !ratingFound {
		fmt.Printf("No rating found for this vehicle")
	}
}

func (c *car) carDetails() {
	fmt.Printf("\n%-5v: %-8v: %-12v ", "Car", c.make, c.model)
	showRating(c.model)
}

func (b *bike) bikeDetails() {
	fmt.Printf("\n%-5v: %-8v: %-12v ", "Bike", b.make, b.model)
	showRating(b.model)
}

func (t *truck) truckDetails() {
	fmt.Printf("\n%-5v: %-8v: %-12v ", "Truck", t.make, t.model)
	showRating(t.model)
}
