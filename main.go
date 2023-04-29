package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	_ "github.com/lib/pq"

	// "./db"

	"go-crud-container-dynamo/db"

)

type User struct {
	ID    string    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// func createDBClient() (*dynamodb.DynamoDB, error) {
//     sess, err := session.NewSession(&aws.Config{
//         Region: aws.String("us-west-2"),
//         Credentials: credentials.NewStaticCredentials(
//             os.Getenv("AWS_ACCESS_KEY_ID"),
//             os.Getenv("AWS_SECRET_ACCESS_KEY"),
//             ""),
//     })
//     if err != nil {
//         return nil, err
//     }

//     svc := dynamodb.New(sess)
//     return svc, nil
// }

// createDBClient creates a new DynamoDB client
func createDBClient() (*dynamodb.DynamoDB, error) {
	// Set up a new session and DynamoDB client
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})
	if err != nil {
		return nil, err
	}

	return dynamodb.New(sess), nil
}




func getUsers(w http.ResponseWriter, r *http.Request) {
	svc, err := createDBClient()
	if err != nil {
		http.Error(w, "failed to connect to DB", http.StatusInternalServerError)
		return
	}

	result, err := svc.Scan(&dynamodb.ScanInput{
		TableName: aws.String("users"),
	})
	if err != nil {
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}

	var users []User
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &users)
	if err != nil {
		http.Error(w, "failed to unmarshal users data", http.StatusInternalServerError)
		return
	}

	// Write the users data as JSON to the response
	jsonBytes, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "failed to encode users data", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

// createUser creates a new user in the DynamoDB table
func createUser(w http.ResponseWriter, r *http.Request) {
	svc, err := createDBClient()
	if err != nil {
		http.Error(w, "failed to connect to DB", http.StatusInternalServerError)
		return
	}

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		http.Error(w, "failed to marshal user data", http.StatusInternalServerError)
		return
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// updateUser updates an existing user in the DynamoDB table
func updateUser(w http.ResponseWriter, r *http.Request) {
	svc, err := createDBClient()
	if err != nil {
		http.Error(w, "failed to connect to DB", http.StatusInternalServerError)
		return
	}

	// id := mux.Vars(r)["id"]
	var user User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}

	av, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		http.Error(w, "failed to marshal user data", http.StatusInternalServerError)
		return
	}

	_, err = svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("users"),
		Item:      av,
	})
	if err != nil {
		http.Error(w, "failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func deleteUser(w http.ResponseWriter, r *http.Request) {
	svc, err := createDBClient()
	if err != nil {
		http.Error(w, "failed to connect to DB", http.StatusInternalServerError)
		return
	}

	// Get the ID parameter from the URL
	vars := mux.Vars(r)
	id := vars["id"]

	_, err = svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("users"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		http.Error(w, "failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}







func main() {
    // Create the DynamoDB table
    err := db.createTable()
    if err != nil {
        log.Fatal("Failed to create DynamoDB table:", err)
    }

    // Set up the HTTP server and start listening for requests
    router := mux.NewRouter()
    router.HandleFunc("/users", createUser).Methods(http.MethodPost)
    router.HandleFunc("/users", getUsers).Methods(http.MethodGet)
    router.HandleFunc("/users/{id}", updateUser).Methods(http.MethodPut)
    router.HandleFunc("/users/{id}", deleteUser).Methods(http.MethodDelete)

    log.Println("Starting HTTP server on port 8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

