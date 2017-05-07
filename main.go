package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/aws/aws-sdk-go/service"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/tidwall/gjson"
	"html/template"
)


type Room struct {
	Name string
	SetPoint string
	Id string
	Power string
	Fan string
}

var tpl *template.Template
var svc *iotdataplane.IoTDataPlane

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

	sess := session.Must(session.NewSession())

	// Create the service's client with the session.
	svc = iotdataplane.New(sess,  &aws.Config {Region: aws.String("us-east-1"),
		Endpoint: aws.String("https://a1ttlrecu8vd0v.iot.us-east-1.amazonaws.com")})
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	log.Fatal(http.ListenAndServe(":8080", router))


}

func Index(w http.ResponseWriter, r *http.Request) {

	params := &iotdataplane.GetThingShadowInput{
		ThingName: aws.String("NucleoHomeServer"), // Required
	}

	resp, err := svc.GetThingShadow(params)
	//n := bytes.IndexByte(resp.Payload, 0)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	desiredRooms := []Room{}
	reportedRooms := []Room{}
	desired := gjson.GetBytes(resp.Payload, "state.desired.rooms")
	reported := gjson.GetBytes(resp.Payload, "state.reported.rooms")

	desired.ForEach(func(key, value gjson.Result) bool{
		room := Room{
			Name: value.Get("room_name").Raw,
			SetPoint: value.Get("set_point").Raw,
			Id: value.Get("room_id").Raw,
			Power: value.Get("power").Raw,
			Fan: value.Get("fan").Raw,
		}

		desiredRooms = append(desiredRooms, room)
		return true // keep iterating
	})

	reported.ForEach(func(key, value gjson.Result) bool{
		room := Room{
			Name: value.Get("room_name").Raw,
			SetPoint: value.Get("set_point").Raw,
			Id: value.Get("room_id").Raw,
			Power: value.Get("power").Raw,
			Fan: value.Get("fan").Raw,
		}

		reportedRooms = append(desiredRooms, room)
		return true // keep iterating
	})
	state := map[string][]Room{
		"desired": desiredRooms,
		"reported": reportedRooms,
	}
	tpl.Execute(w, state)
}