package main

import (
	
)

type User struct{
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}


