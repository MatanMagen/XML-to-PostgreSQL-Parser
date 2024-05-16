package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"strings"

	"go.uber.org/zap"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"

	"github.com/MatanMagen/XML-to-PostgreSQL-Parser/pkg/log"
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

	// Initialize loggers.
	logLevel := zap.NewAtomicLevelAt(zap.FatalLevel)
	serverLogger := log.InitializeLogger(logLevel, true, true)
	defer log.FlushLogger(serverLogger)

	initClient()

	secret, err := secretsClient.Get(context.TODO(), "pgpassword", metaV1.GetOptions{})
	if err != nil {
		serverLogger.Fatal("can not get secret", zap.Error(err))
	}

	// Read the secret data
	password := secret.Data["PGPASSWORD"]

	// Convert password from byte slice to string
	passwordString := string(password)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwordString, dbname)

	db, err := sql.Open("postgres-bla", psqlInfo)
	if err != nil {
		serverLogger.Fatal("can not connect to DB", zap.Error(err))
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		serverLogger.Fatal("can not reach DB", zap.Error(err))
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

	// Extract user data into a map1
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

}

// Function to build a comma-separated list of placeholders
func buildPlaceholders(n int) string {
	placeholders := make([]string, n)
	for i := 0; i < n; i++ {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(placeholders, ", ")
}
