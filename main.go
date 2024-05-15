package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"

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

	// Unmarshal the XML into an interface{}
	var data map[string]interface{}
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println(data)

	// // The top-level XML element is represented as a map
	// m, ok := data.(map[string]interface{})
	// if !ok {
	// 	fmt.Println("XML data is not a map")
	// 	return
	// }

	// // Iterate over the map
	// for _, value := range m {
	// 	// Each element is represented as a slice
	// 	s, ok := value.([]interface{})
	// 	if !ok {
	// 		fmt.Println("Element is not a slice")
	// 		return
	// 	}

	// 	// Iterate over the slice
	// 	for _, item := range s {
	// 		// Each item is represented as a map
	// 		m, ok := item.(map[string]interface{})
	// 		if !ok {
	// 			fmt.Println("Item is not a map")
	// 			continue
	// 		}

	// 		// Iterate over the map
	// 		for key, value := range m {
	// 			// The text content of the elements is represented as a slice
	// 			s, ok := value.([]interface{})
	// 			if ok {
	// 				fmt.Printf("%s: %v\n", key, s[0])
	// 			}
	// 		}
	// 	}
	// }

	// // Insert a row into the users table
	// sqlStatement := `
	//     INSERT INTO users (username, email)
	//     VALUES ($1, $2)
	// `
	// _, err = db.Exec(sqlStatement, "testuser", "testuser@example.com")
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Successfully inserted!")
}
