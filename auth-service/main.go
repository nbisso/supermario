package main

import (
	"auth-service/model"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mySecret []byte = []byte("AllYourBase")
var users []model.User = []model.User{
	{
		Username: "Lucho",
		Password: "123!",
	},
	{
		Username: "pepe",
		Password: "pepe",
	},
}

type NewPing struct {
	Ping string `json:”ping,omitempty”`
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://mongo-0.mongo:27017,mongo-1.mongo:27017,mongo-2.mongo:27017/auth"))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	r := gin.Default()

	r.GET("/auth/pings", func(c *gin.Context) {
		items := []NewPing{}

		collection := client.Database("auth").Collection("pings")

		cur, err := collection.Find(c, NewPing{})

		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}
		defer cur.Close(context.Background())

		for cur.Next(context.Background()) {
			// To decode into a struct, use cursor.Decode()
			result := NewPing{}

			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}
			// do something with result...

			// To get the raw bson bytes use cursor.Current
			// raw := cur.Current

			// do something with raw...

			items = append(items, result)
		}
		if err := cur.Err(); err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, items)

	})

	r.GET("/auth/ping", func(c *gin.Context) {

		name, err := os.Hostname()

		if err != nil {
			c.JSON(http.StatusBadRequest, "Error on binding body")
			return
		}

		newitem := NewPing{name}

		collection := client.Database("auth").Collection("pings")

		insertResult, err := collection.InsertOne(context.TODO(), newitem)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted post with ID:", insertResult.InsertedID)

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/auth/login", func(c *gin.Context) {
		user := model.User{}
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, "Error on binding body")
			return
		}

		result := findUser(user)

		if result != (model.User{}) {

			claims := model.SecurityClaims{
				result.Username,
				1,
				jwt.StandardClaims{
					ExpiresAt: time.Now().Local().Add(time.Minute * 30).Unix(),
					Issuer:    "SuperMario",
				},
			}

			//armo el token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			signedToken, err := token.SignedString(mySecret)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"token": signedToken,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "incorrect user or password",
		})
	})

	r.GET("/auth", func(c *gin.Context) {
		headerToken := c.GetHeader("Authorization")

		if headerToken == "" {
			c.JSON(http.StatusUnauthorized, "")
			return
		}

		resultToken, err := jwt.Parse(headerToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return mySecret, nil
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		claims, ok := resultToken.Claims.(jwt.MapClaims)

		if !ok || !resultToken.Valid {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error en verificar el token",
			})
			return
		}

		username := claims["username"]

		c.Header("X-User", fmt.Sprintf("%v", username))
		c.JSON(http.StatusOK, "")
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func findUser(user model.User) model.User {
	result := model.User{}

	for _, item := range users {
		if item == user {
			result = item
		}
	}

	return result
}
