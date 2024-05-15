package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"strings"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	_ "github.com/lib/pq"
)

const (
	host   = "postgres-service"
	port   = 5432
	user   = "postgres"
	dbname = "postgres"
)

type Element struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
}

type Generic_Struct struct {
	Elements []Element `xml:",any"`
}

type UserData map[string]interface{}

// API client for managing secrets
var secretsClient coreV1Types.SecretInterface

func initClient() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	secretsClient = clientset.CoreV1().Secrets("default")
}

func main() {

	initClient()

	secret, err := secretsClient.Get(context.TODO(), "pgpassword", metaV1.GetOptions{})
	if err != nil {
		panic(err)
	}

	// Read the secret data
	password := secret.Data["PGPASSWORD"]

	// Convert password from byte slice to string
	passwordString := string(password)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwordString, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")

	xmlData := []byte(`
        <user>
            <username>MatanMagen111</username>
            <email>Matan@example.com</email>
        </user>
    `)

	// Unmarshal the XML into Generic_Struct
	var data Generic_Struct
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	// Extract user data into a map
	userData := make(UserData)
	for _, element := range data.Elements {
		tagName := element.XMLName.Local
		tagValue := string(element.Content)
		userData[tagName] = tagValue
	}

	// Build the SQL statement dynamically
	fieldNames := make([]string, 0, len(userData))
	values := make([]interface{}, 0, len(userData))
	for key, value := range userData {
		fieldNames = append(fieldNames, key)
		values = append(values, value)
	}
	sqlStatement := fmt.Sprintf("INSERT INTO users (%s) VALUES (%s)",
		strings.Join(fieldNames, ", "),
		buildPlaceholders(len(userData)))

	println(sqlStatement)

	// Execute the SQL statement
	_, err = db.Exec(sqlStatement, values...)
	if err != nil {
		fmt.Println("Error inserting data:", err)
		return
	}

	fmt.Println("Successfully inserted!")

	// // Unmarshal the XML into an interface{}
	// var data Generic_Struct
	// err = xml.Unmarshal(xmlData, &data)
	// if err != nil {
	// 	fmt.Printf("error: %v\n", err)
	// 	return
	// }
	// // Create a map to hold the tag names and values
	// tagMap := make(map[string]string)

	// // Iterate over the elements and add them to the map
	// for _, element := range data.Elements {
	// 	tagMap[element.XMLName.Local] = string(element.Content)
	// }

	// for _, element := range tagMap {
	// 	// Insert a row into the users table
	// 	sqlStatement :=
	// 	fmt.Sprintf(
	// 		"INSERT INTO help me here !!!!"
	// 	)
	// 	_, err = db.Exec(sqlStatement, "testuser", "testuser@example.com")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println("Successfully inserted!")
	// }

}

// Function to build a comma-separated list of placeholders
func buildPlaceholders(length int) string {
	var placeholders string
	for i := 0; i < length; i++ {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
	}
	return placeholders
}
