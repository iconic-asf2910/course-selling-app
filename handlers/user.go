package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iconic-asf2910/course-selling-app/middleware"
	"github.com/iconic-asf2910/course-selling-app/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"golang.org/x/crypto/bcrypt"
)

func Signup(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User

		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var existingUser models.User

		err := db.Collection("users").
			FindOne(ctx, bson.M{"email": user.Email}).
			Decode(&existingUser)

		if err == nil {
			http.Error(w, "User already exists", http.StatusConflict)
			return
		}

		if err != mongo.ErrNoDocuments {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword(
			[]byte(user.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			http.Error(w, "Could not hash password", http.StatusInternalServerError)
			return
		}

		user.Password = string(hashedPassword)
		user.UserID = int(time.Now().UnixNano() / int64(time.Millisecond))

		_, err = db.Collection("users").InsertOne(ctx, user)

		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "User created successfully",
		})
	}
}

func Login(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var user models.User

		err := db.Collection("users").
			FindOne(ctx, bson.M{"email": loginData.Email}).
			Decode(&user)

		if err == mongo.ErrNoDocuments {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(loginData.Password),
		)

		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token, err := middleware.GenerateToken(user.UserID, "user")

		if err != nil {
			http.Error(w, "Could not create token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

func GetCourses(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := db.Collection("courses").Find(ctx, bson.M{})

		if err != nil {
			http.Error(w, "Could not get courses", http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		var courses []models.Course

		if err := cursor.All(ctx, &courses); err != nil {
			http.Error(w, "Could not read courses", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(courses)
	}
}

func PurchaseCourse(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value(middleware.UserIDKey).(int)

		courseID := r.PathValue("courseID")

		var id int
		if _, err := fmt.Sscanf(courseID, "%d", &id); err != nil {
			http.Error(w, "Invalid course ID", http.StatusBadRequest)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var course models.Course

		err := db.Collection("courses").
			FindOne(ctx, bson.M{"courseId": id}).
			Decode(&course)

		if err == mongo.ErrNoDocuments {
			http.Error(w, "Course not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		purchase := models.Purchase{
			ID:       int(time.Now().UnixNano() / int64(time.Millisecond)),
			UserID:   userID,
			CourseID: id,
		}

		_, err = db.Collection("purchases").InsertOne(ctx, purchase)

		if err != nil {
			http.Error(w, "Could not purchase course", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Course purchased successfully",
		})
	}
}

func GetPurchasedCourses(db *mongo.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userID := r.Context().Value(middleware.UserIDKey).(int)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		cursor, err := db.Collection("purchases").
			Find(ctx, bson.M{"userId": userID})

		if err != nil {
			http.Error(w, "Could not get purchases", http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		var purchases []models.Purchase

		if err := cursor.All(ctx, &purchases); err != nil {
			http.Error(w, "Could not read purchases", http.StatusInternalServerError)
			return
		}

		var courses []models.Course

		for _, purchase := range purchases {

			var course models.Course

			err := db.Collection("courses").
				FindOne(ctx, bson.M{"courseId": purchase.CourseID}).
				Decode(&course)

			if err == nil {
				courses = append(courses, course)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(courses)
	}
}
