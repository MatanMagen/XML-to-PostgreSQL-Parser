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

type Element struct {
	XMLName xml.Name
	Content []byte `xml:",innerxml"`
}

type Generic_Struct struct {
	Elements []Element `xml:",any"`
}

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
	var data Generic_Struct
	err = xml.Unmarshal(xmlData, &data)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Println(data)

	for _, element := range data.Elements {
		fmt.Printf("Element: %s, Content: %s\n", element.XMLName.Local, string(element.Content))
	}

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
