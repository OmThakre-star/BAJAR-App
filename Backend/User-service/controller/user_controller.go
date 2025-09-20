package controller

import (
"encoding/json"
"net/http"


"user-service/database"
"user-service/models"
)


func GetProfile(w http.ResponseWriter, r *http.Request) {
uid := r.Context().Value("userID")
var user models.User
if err := database.DB.First(&user, uid).Error; err != nil {
http.Error(w, "user not found", http.StatusNotFound)
return
}
user.Password = ""
json.NewEncoder(w).Encode(user)
}


func ListUsers(w http.ResponseWriter, r *http.Request) {
var users []models.User
database.DB.Find(&users)
for i := range users {
users[i].Password = ""
}
json.NewEncoder(w).Encode(users)
}