package controller

import (
"time"


"user-service/db"
"user-service/models"
"user-service/utils"
)


func Register(w http.ResponseWriter, r *http.Request) {
var input struct {
Name string `json:"name"`
Email string `json:"email"`
Password string `json:"password"`
}
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
http.Error(w, "invalid json", http.StatusBadRequest)
return
}
hash, _ := utils.HashPassword(input.Password)
user := models.User{Name: input.Name, Email: input.Email, Password: hash}
if err := database.DB.Create(&user).Error; err != nil {
http.Error(w, err.Error(), http.StatusInternalServerError)
return
}
json.NewEncoder(w).Encode(map[string]interface{}{"id": user.ID, "email": user.Email})
}


func Login(w http.ResponseWriter, r *http.Request) {
var input struct {
Email string `json:"email"`
Password string `json:"password"`
}
if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
http.Error(w, "invalid json", http.StatusBadRequest)
return
}
var user models.User
if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
http.Error(w, "invalid credentials", http.StatusUnauthorized)
return
}
if err := utils.CheckPassword(user.Password, input.Password); err != nil {
http.Error(w, "invalid credentials", http.StatusUnauthorized)
return
}
token, _ := utils.CreateToken(user.ID, 24*time.Hour)
json.NewEncoder(w).Encode(map[string]string{"token": token})
}
