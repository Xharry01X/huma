package main

import (
    "context"
    "log"
    "net/http"
    "time"

    "github.com/danielgtaylor/huma/v2"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/bson"
    "golang.org/x/crypto/bcrypt"
)

// User represents a user in the system
type User struct {
    Username string `json:"username" bson:"username"`
    Password string `json:"password" bson:"password"`
}

var client *mongo.Client
var userCollection *mongo.Collection

func connectMongo() {
    var err error
    client, err = mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        log.Fatal(err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    err = client.Connect(ctx)
    if err != nil {
        log.Fatal(err)
    }

    userCollection = client.Database("auth").Collection("users")
}

func signup(ctx huma.Context, input User) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        huma.WriteError(ctx, http.StatusInternalServerError, "Error hashing password")
        return
    }

    input.Password = string(hashedPassword)
    _, err = userCollection.InsertOne(context.TODO(), input)
    if err != nil {
        huma.WriteError(ctx, http.StatusInternalServerError, "Error saving user")
        return
    }

    ctx.WriteModel(http.StatusCreated, input)
}

func login(ctx huma.Context, input User) {
    var foundUser User
    err := userCollection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&foundUser)
    if err != nil {
        huma.WriteError(ctx, http.StatusUnauthorized, "Invalid username or password")
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(input.Password))
    if err != nil {
        huma.WriteError(ctx, http.StatusUnauthorized, "Invalid username or password")
        return
    }

    ctx.WriteModel(http.StatusOK, foundUser)
}

func main() {
    connectMongo()

    app := huma.NewRouter("Auth API", "1.0.0")

    app.Resource("/signup").Post("signup-user", "Create a new user",
        huma.ResponseText(http.StatusCreated, "User created"),
        huma.RequestModel(User{}),
        signup,
    )

    app.Resource("/login").Post("login-user", "Login a user",
        huma.ResponseText(http.StatusOK, "User logged in"),
        huma.RequestModel(User{}),
        login,
    )

    app.Run()
}
