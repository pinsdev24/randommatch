package convert

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"

	"github.com/dimchansky/utfbom"
	"github.com/koki/randommatch/entity"
)

func csvReaderToUsers(r io.Reader) ([]entity.User, error) {
	csvReader := csv.NewReader(utfbom.SkipOnly(r))
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var header []string
	var headerCellIndex = make(map[string]int)
	if len(records) > 0 {
		// skip the header
		header = records[0]
		records = records[1:]
		for i, cell := range header {
			headerCellIndex[strings.TrimSpace(cell)] = i
		}
	}

	var users []entity.User
	for _, record := range records {
		user := entity.User{}

		if val, exists := headerCellIndex["Name"]; exists {
			user.Name = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["Email"]; exists {
			user.Email = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["Gender"]; exists {
			user.Gender = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["Birthday"]; exists {
			user.Birthday = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["PositionHeld"]; exists {
			user.PositionHeld = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["PhoneNumber"]; exists {
			user.PhoneNumber = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["Department"]; exists {
			user.Department = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["Location"]; exists {
			user.Location = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["Seniority"]; exists {
			user.Seniority = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["Role"]; exists {
			user.Role = strings.TrimSpace(record[val])
		}
		if val, exists := headerCellIndex["Tags"]; exists {
			if strings.TrimSpace(record[val]) != "" {
				user.Tags = strings.Split(strings.TrimSpace(record[val]), "-")
			}
		}

		if val, exists := headerCellIndex["Hobbies"]; exists {
			if strings.TrimSpace(record[val]) != "" {
				user.Hobbies = strings.Split(strings.TrimSpace(record[val]), "-")
			}
		}
		if val, exists := headerCellIndex["MatchPreference"]; exists {
			if strings.TrimSpace(record[val]) != "" {
				user.MatchPreference = strings.Split(strings.TrimSpace(record[val]), "-")
			}
		}

		if val, exists := headerCellIndex["MatchPreferenceTime"]; exists {
			if strings.TrimSpace(record[val]) != "" {
				user.MatchPreferenceTime = strings.Split(strings.TrimSpace(record[val]), "-")
			}
		}
		if val, exists := headerCellIndex["MultiMatch"]; exists {
			user.MultiMatch, err = strconv.ParseBool(strings.TrimSpace(record[val]))
		}
		if err != nil {
			log.Printf("Warning wrong boolean string value passed for user %v, value passed: %v\n", user.Name, user.MultiMatch)
			user.MultiMatch = false
		}

		users = append(users, user)
	}
	return users, nil
}

func CsvToUsers(csvFile *multipart.FileHeader) ([]entity.User, error) {
	openedFile, err := csvFile.Open()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return csvReaderToUsers(openedFile)
}

func ConvertRawDataToJson(filepath string) []byte {

	csvFile, err := os.Open(filepath)
	if err != nil {
		log.Println(err)
		return []byte{}
	}

	defer csvFile.Close()
	// Read data
	users, err := csvReaderToUsers(csvFile)
	if err != nil {
		log.Println(err)
		return []byte{}
	}
	// Convert to JSON
	jsonData, err := json.Marshal(users)

	if err != nil {
		log.Println(err)
		return []byte{}
	}

	return jsonData
}

func GenerateJsonFile(filename string) {

	jsonData := ConvertRawDataToJson(filename)

	jsonFile, err := os.Create("./data.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
}
