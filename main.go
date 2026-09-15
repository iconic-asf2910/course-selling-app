package main

type Course struct {
	CourseID     int    `json:"courseid"`
	Name         string `json:"title"`
	Description  string `json:"description"`
	Price        int    `json:"price"`
	InstructorID int    `json:"instructorid"`
	Content      string `json:"content"`
}

type User struct {
	UserID   int    `json:"userid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Admin struct {
	AdminID  int    `json:"adminid"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Purchase struct {
	ID       int `json:"id"`
	UserID   int `json:"userid"`
	CourseID int `json:"courseid"`
}
