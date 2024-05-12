package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	coreV1Types "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"

	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	user   = "postgres"
	dbname = "my_postgres"
)

// API client for managing secrets
var secretsClient coreV1Types.SecretInterface

func initClient() {
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
	password := secret.Data["POSTGRES_PASSWORD"]

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
}
