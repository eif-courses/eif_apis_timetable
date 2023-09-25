package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"github.com/rs/cors"
	"log"
	"net/http"
	"time"
)

var (
	c        = cache.New(5*time.Minute, 10*time.Minute) // Cache data for 5 minutes with a cleanup interval of 10 minutes
	cacheKey = "all_teachers_ids"
)

var teacherCache = cache.New(5*time.Minute, 10*time.Minute)   // Cache data for 5 minutes with a cleanup interval of 10 minutes
var groupCache = cache.New(5*time.Minute, 10*time.Minute)     // Cache data for 5 minutes with a cleanup interval of 10 minutes
var classroomCache = cache.New(5*time.Minute, 10*time.Minute) // Cache data for 5 minutes with a cleanup interval of 10 minutes

var (
	clids         = cache.New(5*time.Minute, 10*time.Minute) // Cache data for 5 minutes with a cleanup interval of 10 minutes
	cacheKeyClids = "all_classroom_ids"
)
var (
	groupids         = cache.New(5*time.Minute, 10*time.Minute) // Cache data for 5 minutes with a cleanup interval of 10 minutes
	cacheKeyGroupIds = "all_group_ids"
)

func GetAllTeachersIDs(w http.ResponseWriter, r *http.Request) {

	if cachedResult, found := c.Get(cacheKey); found {
		// Cache hit: Return the cached data
		w.Header().Set("Content-Type", "application/json")
		w.Write(cachedResult.([]byte))
		return
	}
	tTeacher := time.Now().Local()
	// Calculate the date of the first day of the current week (Monday)
	startTeacher := tTeacher
	if tTeacher.Weekday() != time.Monday {
		// Subtract the appropriate number of days to get to Monday
		daysUntilMonday := int(time.Monday - tTeacher.Weekday())
		startTeacher = tTeacher.AddDate(0, 0, daysUntilMonday)
	}

	// Calculate the date of the last day of the current week (Sunday)
	endTeacher := startTeacher.AddDate(0, 0, 6)

	dateFromTeacher := startTeacher.Format("2006-01-02")
	dateToTeacher := endTeacher.Format("2006-01-02")

	fmt.Printf("nuo klases id: %s\n", dateFromTeacher)

	currentYearTeacher := time.Now().Year()

	payload := map[string]interface{}{
		"__args": []interface{}{
			nil,
			currentYearTeacher,
			map[string]interface{}{
				"vt_filter": map[string]interface{}{
					"datefrom": dateFromTeacher,
					"dateto":   dateToTeacher,
				},
			},
			map[string]interface{}{
				"op": "fetch",
				"needed_part": map[string]interface{}{
					"teachers":   []string{"short", "name", "firstname", "lastname", "subname", "code", "cb_hidden", "expired", "firstname", "lastname", "short"},
					"classes":    []string{"short", "name", "firstname", "lastname", "subname", "code", "classroomid"},
					"classrooms": []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
					"subjects":   []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
				},
				"needed_combos": map[string]interface{}{},
			},
		},
		"__gsh": "00000000",
	}

	// You need to specify the URL here
	url := "https://vikoeif.edupage.org/rpr/server/maindbi.js?__func=mainDBIAccessor"

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make an HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		teachersIDs := result["r"].(map[string]interface{})["tables"].([]interface{})[0].(map[string]interface{})["data_rows"]
		jsonResponse, err := json.Marshal(teachersIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Cache the result
		c.Set(cacheKey, jsonResponse, cache.DefaultExpiration)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		http.Error(w, "API request failed", http.StatusInternalServerError)
	}
}
func GetAllClassroomsIDs(w http.ResponseWriter, r *http.Request) {

	if cachedResult, found := clids.Get(cacheKeyClids); found {
		// Cache hit: Return the cached data
		w.Header().Set("Content-Type", "application/json")
		w.Write(cachedResult.([]byte))
		return
	}

	tTeacher := time.Now().Local()
	// Calculate the date of the first day of the current week (Monday)
	startTeacher := tTeacher
	if tTeacher.Weekday() != time.Monday {
		// Subtract the appropriate number of days to get to Monday
		daysUntilMonday := int(time.Monday - tTeacher.Weekday())
		startTeacher = tTeacher.AddDate(0, 0, daysUntilMonday)
	}

	// Calculate the date of the last day of the current week (Sunday)
	endTeacher := startTeacher.AddDate(0, 0, 6)

	dateFromTeacher := startTeacher.Format("2006-01-02")
	dateToTeacher := endTeacher.Format("2006-01-02")

	fmt.Printf("nuo klases id: %s\n", dateFromTeacher)

	currentYearTeacher := time.Now().Year()

	payload := map[string]interface{}{
		"__args": []interface{}{
			nil,
			currentYearTeacher,
			map[string]interface{}{
				"vt_filter": map[string]interface{}{
					"datefrom": dateFromTeacher,
					"dateto":   dateToTeacher,
				},
			},
			map[string]interface{}{
				"op": "fetch",
				"needed_part": map[string]interface{}{
					"teachers":   []string{"short", "name", "firstname", "lastname", "subname", "code", "cb_hidden", "expired", "firstname", "lastname", "short"},
					"classes":    []string{"short", "name", "firstname", "lastname", "subname", "code", "classroomid"},
					"classrooms": []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
					"subjects":   []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
				},
				"needed_combos": map[string]interface{}{},
			},
		},
		"__gsh": "00000000",
	}

	// You need to specify the URL here
	url := "https://vikoeif.edupage.org/rpr/server/maindbi.js?__func=mainDBIAccessor"

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make an HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		teachersIDs := result["r"].(map[string]interface{})["tables"].([]interface{})[2].(map[string]interface{})["data_rows"]
		jsonResponse, err := json.Marshal(teachersIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Cache the result
		clids.Set(cacheKeyClids, jsonResponse, cache.DefaultExpiration)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		http.Error(w, "API request failed", http.StatusInternalServerError)
	}
}
func GetAllGroupsIDs(w http.ResponseWriter, r *http.Request) {

	if cachedResult, found := groupids.Get(cacheKeyGroupIds); found {
		// Cache hit: Return the cached data
		w.Header().Set("Content-Type", "application/json")
		w.Write(cachedResult.([]byte))
		return
	}

	tTeacher := time.Now().Local()
	// Calculate the date of the first day of the current week (Monday)
	startTeacher := tTeacher
	if tTeacher.Weekday() != time.Monday {
		// Subtract the appropriate number of days to get to Monday
		daysUntilMonday := int(time.Monday - tTeacher.Weekday())
		startTeacher = tTeacher.AddDate(0, 0, daysUntilMonday)
	}

	// Calculate the date of the last day of the current week (Sunday)
	endTeacher := startTeacher.AddDate(0, 0, 6)

	dateFromTeacher := startTeacher.Format("2006-01-02")
	dateToTeacher := endTeacher.Format("2006-01-02")

	fmt.Printf("nuo klases id: %s\n", dateFromTeacher)

	currentYearTeacher := time.Now().Year()

	payload := map[string]interface{}{
		"__args": []interface{}{
			nil,
			currentYearTeacher,
			map[string]interface{}{
				"vt_filter": map[string]interface{}{
					"datefrom": dateFromTeacher,
					"dateto":   dateToTeacher,
				},
			},
			map[string]interface{}{
				"op": "fetch",
				"needed_part": map[string]interface{}{
					"teachers":   []string{"short", "name", "firstname", "lastname", "subname", "code", "cb_hidden", "expired", "firstname", "lastname", "short"},
					"classes":    []string{"short", "name", "firstname", "lastname", "subname", "code", "classroomid"},
					"classrooms": []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
					"subjects":   []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
				},
				"needed_combos": map[string]interface{}{},
			},
		},
		"__gsh": "00000000",
	}

	// You need to specify the URL here
	url := "https://vikoeif.edupage.org/rpr/server/maindbi.js?__func=mainDBIAccessor"

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Make an HTTP POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		teachersIDs := result["r"].(map[string]interface{})["tables"].([]interface{})[3].(map[string]interface{})["data_rows"]
		jsonResponse, err := json.Marshal(teachersIDs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Cache the result
		groupids.Set(cacheKeyGroupIds, jsonResponse, cache.DefaultExpiration)

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		http.Error(w, "API request failed", http.StatusInternalServerError)
	}
}

func getSingleTeacher(teacherID string) ([]map[string]interface{}, error) {
	// Define the URL
	url := "https://vikoeif.edupage.org/rpr/server/maindbi.js?__func=mainDBIAccessor"
	otherEndpoint := "https://vikoeif.edupage.org/timetable/server/currenttt.js?__func=curentttGetData" // Replace with your actual endpoint

	// Calculate the date range
	todayTeacher := time.Now()
	tTeacher := time.Now().Local()
	// Calculate the date of the first day of the current week (Monday)
	startTeacher := tTeacher
	if tTeacher.Weekday() != time.Monday {
		// Subtract the appropriate number of days to get to Monday
		daysUntilMonday := int(time.Monday - tTeacher.Weekday())
		startTeacher = tTeacher.AddDate(0, 0, daysUntilMonday)
	}

	// Calculate the date of the last day of the current week (Sunday)
	endTeacher := startTeacher.AddDate(0, 0, 6)

	dateFromTeacher := startTeacher.Format("2006-01-02")
	dateToTeacher := endTeacher.Format("2006-01-02")

	currentYearTeacher := todayTeacher.Year()

	// Create the payload
	payload := map[string]interface{}{
		"__args": []interface{}{
			nil,
			currentYearTeacher,
			map[string]interface{}{
				"vt_filter": map[string]interface{}{
					"datefrom": dateFromTeacher,
					"dateto":   dateToTeacher,
				},
			},
			map[string]interface{}{
				"op": "fetch",
				"needed_part": map[string]interface{}{
					"teachers":   []string{"short", "name", "firstname", "lastname", "subname", "code", "cb_hidden", "expired", "firstname", "lastname", "short"},
					"classes":    []string{"short", "name", "firstname", "lastname", "subname", "code", "classroomid"},
					"classrooms": []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
					"subjects":   []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
				},
				"needed_combos": map[string]interface{}{},
			},
		},
		"__gsh": "00000000",
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Make the first HTTP POST request to fetch data
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			return nil, err
		}

		teachersOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[0].(map[string]interface{})["data_rows"]
		subjectsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[1].(map[string]interface{})["data_rows"]
		classroomsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[2].(map[string]interface{})["data_rows"]
		groupsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[3].(map[string]interface{})["data_rows"]

		// Create maps to store mappings
		teachersMap := make(map[string]interface{})
		subjectsMap := make(map[string]interface{})
		classroomsMap := make(map[string]interface{})
		groupsMap := make(map[string]interface{})

		// Populate the maps
		for _, teacher := range teachersOriginal.([]interface{}) {
			teacherMap := teacher.(map[string]interface{})
			id := teacherMap["id"].(string)
			short := teacherMap["short"].(string)
			teachersMap[id] = short
		}

		for _, subject := range subjectsOriginal.([]interface{}) {
			subjectMap := subject.(map[string]interface{})
			id := subjectMap["id"].(string)
			short := subjectMap["short"].(string)
			subjectsMap[id] = short
		}

		for _, classroom := range classroomsOriginal.([]interface{}) {
			classroomMap := classroom.(map[string]interface{})
			id := classroomMap["id"].(string)
			short := classroomMap["short"].(string)
			classroomsMap[id] = short
		}

		for _, group := range groupsOriginal.([]interface{}) {
			groupMap := group.(map[string]interface{})
			id := groupMap["id"].(string)
			short := groupMap["short"].(string)
			groupsMap[id] = short
		}

		// Create the payload for the second request
		otherPayload := map[string]interface{}{
			"__args": []interface{}{
				nil,
				map[string]interface{}{
					"year":                 currentYearTeacher,
					"datefrom":             dateFromTeacher,
					"dateto":               dateToTeacher,
					"table":                "teachers",
					"id":                   teacherID,
					"showColors":           true,
					"showIgroupsInClasses": false,
					"showOrig":             true,
					"log_module":           "CurrentTTView",
				},
			},
			"__gsh": "00000000",
		}

		// Convert the otherPayload to JSON
		otherPayloadBytes, err := json.Marshal(otherPayload)
		if err != nil {
			return nil, err
		}

		// Make the second HTTP POST request to fetch teacher data by ID
		respSingleTeacher, err := http.Post(otherEndpoint, "application/json", bytes.NewBuffer(otherPayloadBytes))
		if err != nil {
			return nil, err
		}
		defer respSingleTeacher.Body.Close()

		if respSingleTeacher.StatusCode == http.StatusOK {
			var teacherData map[string]interface{}
			decoder := json.NewDecoder(respSingleTeacher.Body)
			if err := decoder.Decode(&teacherData); err != nil {
				return nil, err
			}

			// Process the teacherData as needed
			// Process the teacherData as needed
			teacherItems := teacherData["r"].(map[string]interface{})["ttitems"].([]interface{})

			// Convert the []interface{} to []map[string]interface{}
			var teacherItemsList []map[string]interface{}
			for _, item := range teacherItems {
				if teacherItemMap, ok := item.(map[string]interface{}); ok {
					teacherItemsList = append(teacherItemsList, teacherItemMap)
				}
			}

			// Replace class IDs and classroom IDs with actual data
			for i, item := range teacherItemsList {
				classIDs := item["classids"].([]interface{})
				classroomIDs := item["classroomids"].([]interface{})
				subjectID := item["subjectid"]
				teacherIDs := item["teacherids"]

				// Lookup and replace class IDs with actual data
				var classNames []string
				for _, classID := range classIDs {
					if id, ok := classID.(string); ok {
						if className, found := groupsMap[id]; found {
							classNames = append(classNames, className.(string))
						}
					}
				}
				teacherItemsList[i]["classids"] = classNames

				// Lookup and replace classroom IDs with actual data
				var classroomNames []string
				for _, classroomID := range classroomIDs {
					if id, ok := classroomID.(string); ok {
						if classroomName, found := classroomsMap[id]; found {
							classroomNames = append(classroomNames, classroomName.(string))
						}
					}
				}
				teacherItemsList[i]["classroomids"] = classroomNames

				// Check if subjectID is an array of strings
				if subjectIDs, ok := subjectID.([]interface{}); ok {
					var subjectNames []string
					for _, subjectID := range subjectIDs {
						if id, ok := subjectID.(string); ok {
							if subjectName, found := subjectsMap[id]; found {
								subjectNames = append(subjectNames, subjectName.(string))
							}
						}
					}
					teacherItemsList[i]["subjectid"] = subjectNames
				} else if id, ok := subjectID.(string); ok { // Check if subjectID is a string
					if subjectName, found := subjectsMap[id]; found {
						teacherItemsList[i]["subjectid"] = []string{subjectName.(string)}
					}
				}

				// Check if teacherIDs is an array of strings
				if teacherIDList, ok := teacherIDs.([]interface{}); ok {
					var teacherNames []string
					for _, teacherID := range teacherIDList {
						if id, ok := teacherID.(string); ok {
							if teacherName, found := teachersMap[id]; found {
								teacherNames = append(teacherNames, teacherName.(string))
							}
						}
					}
					teacherItemsList[i]["teacherids"] = teacherNames
				} else if id, ok := teacherIDs.(string); ok { // Check if teacherIDs is a string
					if teacherName, found := teachersMap[id]; found {
						teacherItemsList[i]["teacherids"] = []string{teacherName.(string)}
					}
				}

			}

			return teacherItemsList, nil
		} else {
			return nil, fmt.Errorf("HTTP request for teacher data failed with status code: %d", respSingleTeacher.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("HTTP request for initial data failed with status code: %d", resp.StatusCode)
	}
}

func GetSingleTeacherHandler(w http.ResponseWriter, r *http.Request) {
	// Get the teacher ID from the request parameters or query string
	teacherID := r.URL.Query().Get("teacher_id")

	// Check if the result is already cached
	if cachedResult, found := teacherCache.Get(teacherID); found {
		// Respond with the cached result
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cachedResult); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Call the getSingleTeacher function to fetch the teacher data
	teacherItems, err := getSingleTeacher(teacherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache the result
	teacherCache.Set(teacherID, teacherItems, cache.DefaultExpiration)

	// Respond with the teacher data as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(teacherItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getSingleGroup(groupID string) ([]map[string]interface{}, error) {
	// Define the URL
	url := "https://vikoeif.edupage.org/rpr/server/maindbi.js?__func=mainDBIAccessor"
	otherEndpoint := "https://vikoeif.edupage.org/timetable/server/currenttt.js?__func=curentttGetData" // Replace with your actual endpoint

	// Calculate the date range
	todayTeacher := time.Now()
	tTeacher := time.Now().Local()
	// Calculate the date of the first day of the current week (Monday)
	startTeacher := tTeacher
	if tTeacher.Weekday() != time.Monday {
		// Subtract the appropriate number of days to get to Monday
		daysUntilMonday := int(time.Monday - tTeacher.Weekday())
		startTeacher = tTeacher.AddDate(0, 0, daysUntilMonday)
	}

	// Calculate the date of the last day of the current week (Sunday)
	endTeacher := startTeacher.AddDate(0, 0, 6)

	dateFromTeacher := startTeacher.Format("2006-01-02")
	dateToTeacher := endTeacher.Format("2006-01-02")
	currentYearTeacher := todayTeacher.Year()

	// Create the payload
	payload := map[string]interface{}{
		"__args": []interface{}{
			nil,
			currentYearTeacher,
			map[string]interface{}{
				"vt_filter": map[string]interface{}{
					"datefrom": dateFromTeacher,
					"dateto":   dateToTeacher,
				},
			},
			map[string]interface{}{
				"op": "fetch",
				"needed_part": map[string]interface{}{
					"teachers":   []string{"short", "name", "firstname", "lastname", "subname", "code", "cb_hidden", "expired", "firstname", "lastname", "short"},
					"classes":    []string{"short", "name", "firstname", "lastname", "subname", "code", "classroomid"},
					"classrooms": []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
					"subjects":   []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
				},
				"needed_combos": map[string]interface{}{},
			},
		},
		"__gsh": "00000000",
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Make the first HTTP POST request to fetch data
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			return nil, err
		}

		teachersOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[0].(map[string]interface{})["data_rows"]
		subjectsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[1].(map[string]interface{})["data_rows"]
		classroomsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[2].(map[string]interface{})["data_rows"]
		groupsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[3].(map[string]interface{})["data_rows"]

		// Create maps to store mappings
		teachersMap := make(map[string]interface{})
		subjectsMap := make(map[string]interface{})
		classroomsMap := make(map[string]interface{})
		groupsMap := make(map[string]interface{})

		// Populate the maps
		for _, teacher := range teachersOriginal.([]interface{}) {
			teacherMap := teacher.(map[string]interface{})
			id := teacherMap["id"].(string)
			short := teacherMap["short"].(string)
			teachersMap[id] = short
		}

		for _, subject := range subjectsOriginal.([]interface{}) {
			subjectMap := subject.(map[string]interface{})
			id := subjectMap["id"].(string)
			short := subjectMap["short"].(string)
			subjectsMap[id] = short
		}

		for _, classroom := range classroomsOriginal.([]interface{}) {
			classroomMap := classroom.(map[string]interface{})
			id := classroomMap["id"].(string)
			short := classroomMap["short"].(string)
			classroomsMap[id] = short
		}

		for _, group := range groupsOriginal.([]interface{}) {
			groupMap := group.(map[string]interface{})
			id := groupMap["id"].(string)
			short := groupMap["short"].(string)
			groupsMap[id] = short
		}

		// Create the payload for the second request
		otherPayload := map[string]interface{}{
			"__args": []interface{}{
				nil,
				map[string]interface{}{
					"year":                 currentYearTeacher,
					"datefrom":             dateFromTeacher,
					"dateto":               dateToTeacher,
					"table":                "classes",
					"id":                   groupID,
					"showColors":           true,
					"showIgroupsInClasses": false,
					"showOrig":             true,
					"log_module":           "CurrentTTView",
				},
			},
			"__gsh": "00000000",
		}

		// Convert the otherPayload to JSON
		otherPayloadBytes, err := json.Marshal(otherPayload)
		if err != nil {
			return nil, err
		}

		// Make the second HTTP POST request to fetch teacher data by ID
		respSingleGroup, err := http.Post(otherEndpoint, "application/json", bytes.NewBuffer(otherPayloadBytes))
		if err != nil {
			return nil, err
		}
		defer respSingleGroup.Body.Close()

		if respSingleGroup.StatusCode == http.StatusOK {
			var groupData map[string]interface{}
			decoder := json.NewDecoder(respSingleGroup.Body)
			if err := decoder.Decode(&groupData); err != nil {
				return nil, err
			}

			// Process the teacherData as needed
			// Process the teacherData as needed
			groupItems := groupData["r"].(map[string]interface{})["ttitems"].([]interface{})

			// Convert the []interface{} to []map[string]interface{}
			var groupItemsList []map[string]interface{}
			for _, item := range groupItems {
				if groupItemMap, ok := item.(map[string]interface{}); ok {
					groupItemsList = append(groupItemsList, groupItemMap)
				}
			}

			// Replace class IDs and classroom IDs with actual data
			for i, item := range groupItemsList {
				// Replace class IDs with actual class names
				classIDs := item["classids"].([]interface{})
				var classNames []string
				for _, classID := range classIDs {
					if id, ok := classID.(string); ok {
						if className, found := groupsMap[id]; found {
							classNames = append(classNames, className.(string))
						}
					}
				}
				groupItemsList[i]["classids"] = classNames

				// Replace classroom IDs with actual classroom names
				classroomIDs := item["classroomids"].([]interface{})
				var classroomNames []string
				for _, classroomID := range classroomIDs {
					if id, ok := classroomID.(string); ok {
						if classroomName, found := classroomsMap[id]; found {
							classroomNames = append(classroomNames, classroomName.(string))
						}
					}
				}
				groupItemsList[i]["classroomids"] = classroomNames

				// Replace subject IDs with actual subject names
				subjectID := item["subjectid"]
				if subjectID != nil {
					if subjectIDs, ok := subjectID.([]interface{}); ok {
						var subjectNames []string
						for _, subjectID := range subjectIDs {
							if id, ok := subjectID.(string); ok {
								if subjectName, found := subjectsMap[id]; found {
									subjectNames = append(subjectNames, subjectName.(string))
								}
							}
						}
						groupItemsList[i]["subjectid"] = subjectNames
					} else if id, ok := subjectID.(string); ok {
						if subjectName, found := subjectsMap[id]; found {
							groupItemsList[i]["subjectid"] = []string{subjectName.(string)}
						}
					}
				}

				// Replace teacher IDs with actual teacher names
				teacherIDs := item["teacherids"]
				if teacherIDs != nil {
					if teacherIDList, ok := teacherIDs.([]interface{}); ok {
						var teacherNames []string
						for _, teacherID := range teacherIDList {
							if id, ok := teacherID.(string); ok {
								if teacherName, found := teachersMap[id]; found {
									teacherNames = append(teacherNames, teacherName.(string))
								}
							}
						}
						groupItemsList[i]["teacherids"] = teacherNames
					} else if id, ok := teacherIDs.(string); ok {
						if teacherName, found := teachersMap[id]; found {
							groupItemsList[i]["teacherids"] = []string{teacherName.(string)}
						}
					}
				}
			}

			return groupItemsList, nil
		} else {
			return nil, fmt.Errorf("HTTP request for teacher data failed with status code: %d", respSingleGroup.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("HTTP request for initial data failed with status code: %d", resp.StatusCode)
	}
}

func GetSingleGroupHandler(w http.ResponseWriter, r *http.Request) {
	// Get the teacher ID from the request parameters or query string
	groupID := r.URL.Query().Get("group_id")

	// Check if the result is already cached
	if cachedResult, found := groupCache.Get(groupID); found {
		// Respond with the cached result
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cachedResult); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Call the getSingleTeacher function to fetch the teacher data
	groupItems, err := getSingleGroup(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache the result
	groupCache.Set(groupID, groupItems, cache.DefaultExpiration)

	// Respond with the teacher data as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(groupItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getSingleClassroom(classroomID string) ([]map[string]interface{}, error) {
	// Define the URL
	url := "https://vikoeif.edupage.org/rpr/server/maindbi.js?__func=mainDBIAccessor"
	otherEndpoint := "https://vikoeif.edupage.org/timetable/server/currenttt.js?__func=curentttGetData" // Replace with your actual endpoint

	// Calculate the date range
	todayTeacher := time.Now()
	tTeacher := time.Now().Local()
	// Calculate the date of the first day of the current week (Monday)
	startTeacher := tTeacher
	if tTeacher.Weekday() != time.Monday {
		// Subtract the appropriate number of days to get to Monday
		daysUntilMonday := int(time.Monday - tTeacher.Weekday())
		startTeacher = tTeacher.AddDate(0, 0, daysUntilMonday)
	}

	// Calculate the date of the last day of the current week (Sunday)
	endTeacher := startTeacher.AddDate(0, 0, 6)

	dateFromTeacher := startTeacher.Format("2006-01-02")
	dateToTeacher := endTeacher.Format("2006-01-02")
	currentYearTeacher := todayTeacher.Year()

	// Create the payload
	payload := map[string]interface{}{
		"__args": []interface{}{
			nil,
			currentYearTeacher,
			map[string]interface{}{
				"vt_filter": map[string]interface{}{
					"datefrom": dateFromTeacher,
					"dateto":   dateToTeacher,
				},
			},
			map[string]interface{}{
				"op": "fetch",
				"needed_part": map[string]interface{}{
					"teachers":   []string{"short", "name", "firstname", "lastname", "subname", "code", "cb_hidden", "expired", "firstname", "lastname", "short"},
					"classes":    []string{"short", "name", "firstname", "lastname", "subname", "code", "classroomid"},
					"classrooms": []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
					"subjects":   []string{"short", "name", "firstname", "lastname", "subname", "code", "name", "short"},
				},
				"needed_combos": map[string]interface{}{},
			},
		},
		"__gsh": "00000000",
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Make the first HTTP POST request to fetch data
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var result map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&result); err != nil {
			return nil, err
		}

		teachersOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[0].(map[string]interface{})["data_rows"]
		subjectsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[1].(map[string]interface{})["data_rows"]
		classroomsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[2].(map[string]interface{})["data_rows"]
		groupsOriginal := result["r"].(map[string]interface{})["tables"].([]interface{})[3].(map[string]interface{})["data_rows"]

		// Create maps to store mappings
		teachersMap := make(map[string]interface{})
		subjectsMap := make(map[string]interface{})
		classroomsMap := make(map[string]interface{})
		groupsMap := make(map[string]interface{})

		// Populate the maps
		for _, teacher := range teachersOriginal.([]interface{}) {
			teacherMap := teacher.(map[string]interface{})
			id := teacherMap["id"].(string)
			short := teacherMap["short"].(string)
			teachersMap[id] = short
		}

		for _, subject := range subjectsOriginal.([]interface{}) {
			subjectMap := subject.(map[string]interface{})
			id := subjectMap["id"].(string)
			short := subjectMap["short"].(string)
			subjectsMap[id] = short
		}

		for _, classroom := range classroomsOriginal.([]interface{}) {
			classroomMap := classroom.(map[string]interface{})
			id := classroomMap["id"].(string)
			short := classroomMap["short"].(string)
			classroomsMap[id] = short
		}

		for _, group := range groupsOriginal.([]interface{}) {
			groupMap := group.(map[string]interface{})
			id := groupMap["id"].(string)
			short := groupMap["short"].(string)
			groupsMap[id] = short
		}

		// Create the payload for the second request
		otherPayload := map[string]interface{}{
			"__args": []interface{}{
				nil,
				map[string]interface{}{
					"year":                 currentYearTeacher,
					"datefrom":             dateFromTeacher,
					"dateto":               dateToTeacher,
					"table":                "classrooms",
					"id":                   classroomID,
					"showColors":           true,
					"showIgroupsInClasses": false,
					"showOrig":             true,
					"log_module":           "CurrentTTView",
				},
			},
			"__gsh": "00000000",
		}

		// Convert the otherPayload to JSON
		otherPayloadBytes, err := json.Marshal(otherPayload)
		if err != nil {
			return nil, err
		}

		// Make the second HTTP POST request to fetch teacher data by ID
		respSingleClassroom, err := http.Post(otherEndpoint, "application/json", bytes.NewBuffer(otherPayloadBytes))
		if err != nil {
			return nil, err
		}
		defer respSingleClassroom.Body.Close()

		if respSingleClassroom.StatusCode == http.StatusOK {
			var classroomData map[string]interface{}
			decoder := json.NewDecoder(respSingleClassroom.Body)
			if err := decoder.Decode(&classroomData); err != nil {
				return nil, err
			}

			// Process the teacherData as needed
			// Process the teacherData as needed
			classroomItems := classroomData["r"].(map[string]interface{})["ttitems"].([]interface{})

			// Convert the []interface{} to []map[string]interface{}
			var classroomItemsList []map[string]interface{}
			for _, item := range classroomItems {
				if classroomItemMap, ok := item.(map[string]interface{}); ok {
					classroomItemsList = append(classroomItemsList, classroomItemMap)
				}
			}

			// Replace class IDs and classroom IDs with actual data
			for i, item := range classroomItemsList {
				// Replace class IDs with actual class names
				classIDs := item["classids"].([]interface{})
				var classNames []string
				for _, classID := range classIDs {
					if id, ok := classID.(string); ok {
						if className, found := groupsMap[id]; found {
							classNames = append(classNames, className.(string))
						}
					}
				}
				classroomItemsList[i]["classids"] = classNames

				// Replace classroom IDs with actual classroom names
				classroomIDs := item["classroomids"].([]interface{})
				var classroomNames []string
				for _, classroomID := range classroomIDs {
					if id, ok := classroomID.(string); ok {
						if classroomName, found := classroomsMap[id]; found {
							classroomNames = append(classroomNames, classroomName.(string))
						}
					}
				}
				classroomItemsList[i]["classroomids"] = classroomNames

				// Replace subject IDs with actual subject names
				subjectID := item["subjectid"]
				if subjectID != nil {
					if subjectIDs, ok := subjectID.([]interface{}); ok {
						var subjectNames []string
						for _, subjectID := range subjectIDs {
							if id, ok := subjectID.(string); ok {
								if subjectName, found := subjectsMap[id]; found {
									subjectNames = append(subjectNames, subjectName.(string))
								}
							}
						}
						classroomItemsList[i]["subjectid"] = subjectNames
					} else if id, ok := subjectID.(string); ok {
						if subjectName, found := subjectsMap[id]; found {
							classroomItemsList[i]["subjectid"] = []string{subjectName.(string)}
						}
					}
				}

				// Replace teacher IDs with actual teacher names
				teacherIDs := item["teacherids"]
				if teacherIDs != nil {
					if teacherIDList, ok := teacherIDs.([]interface{}); ok {
						var teacherNames []string
						for _, teacherID := range teacherIDList {
							if id, ok := teacherID.(string); ok {
								if teacherName, found := teachersMap[id]; found {
									teacherNames = append(teacherNames, teacherName.(string))
								}
							}
						}
						classroomItemsList[i]["teacherids"] = teacherNames
					} else if id, ok := teacherIDs.(string); ok {
						if teacherName, found := teachersMap[id]; found {
							classroomItemsList[i]["teacherids"] = []string{teacherName.(string)}
						}
					}
				}
			}

			return classroomItemsList, nil
		} else {
			return nil, fmt.Errorf("HTTP request for teacher data failed with status code: %d", respSingleClassroom.StatusCode)
		}
	} else {
		return nil, fmt.Errorf("HTTP request for initial data failed with status code: %d", resp.StatusCode)
	}
}

func GetSingleClassroomHandler(w http.ResponseWriter, r *http.Request) {
	// Get the teacher ID from the request parameters or query string
	classroomID := r.URL.Query().Get("classroom_id")

	// Check if the result is already cached
	if cachedResult, found := classroomCache.Get(classroomID); found {
		// Respond with the cached result
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cachedResult); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// Call the getSingleTeacher function to fetch the teacher data
	classroomItems, err := getSingleClassroom(classroomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Cache the result
	classroomCache.Set(classroomID, classroomItems, cache.DefaultExpiration)

	// Respond with the teacher data as JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(classroomItems); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	// Create a new CORS handler
	c := cors.AllowAll() // You can customize CORS settings if needed

	// Wrap your existing handler with the CORS middleware
	handler := c.Handler(mux)

	// Define your routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	mux.HandleFunc("/api/v1/timetable/teacher", GetSingleTeacherHandler)
	mux.HandleFunc("/api/v1/timetable/group", GetSingleGroupHandler)
	mux.HandleFunc("/api/v1/timetable/classroom", GetSingleClassroomHandler)
	mux.HandleFunc("/api/v1/timetable/teachers", GetAllTeachersIDs)
	mux.HandleFunc("/api/v1/timetable/classrooms", GetAllClassroomsIDs)
	mux.HandleFunc("/api/v1/timetable/groups", GetAllGroupsIDs)

	// Start the server
	log.Println("** Service Started on Port 8080 **")
	err := http.ListenAndServe("0.0.0.0:8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}
