package main

import (
	"backend/util"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"

	"github.com/rs/cors"
)

type CheckerData struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	HumanSensor bool    `json:human_sensor`
	LightSensor bool    `json:light_sensor`
}

// LocationDataはフロントエンドから送信される位置情報
type LocationData struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// HomeLocationはデータベースから取得する自宅の緯度経度
type HomeLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Resultは判定結果
type Result struct {
	Location LocationData `json:"location"`
	Home     HomeLocation `json:"home"`
	Distance float64      `json:"distance"`
	Status   string       `json:"status"`
}

// Userはユーザー情報
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Claims JWTクレーム
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// JWTの秘密鍵
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

var db *sql.DB

func main() {
	// jwt keyを取得
	if jwtKey == nil {
		jwtKey = []byte("supersecretkey")
		fmt.Println("JWT key from default")
	} else {
		fmt.Println("JWT key from env")
	}

	var err error
	db, err = connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// DB接続確認
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to database")
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/signup", signup).Methods("POST")
	router.HandleFunc("/check-location", checkLocation).Methods("POST")
	router.HandleFunc("/put-user-location", putUserLocation).Methods("POST")
	router.HandleFunc("/put-home-location", putHomeLocation).Methods("POST")
	router.HandleFunc("/change-email", changeEmail).Methods("POST")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	})

	h := cors.Default().Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Server listening on port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, h))
}

func connect() (*sql.DB, error) {
	dbDriver := "postgres"
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME"))
	log.Println(dsn)
	db, err := sql.Open(dbDriver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// 指定した緯度経度からの距離を計算する関数
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // 地球の半径 (km)
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadius * c
	return distance // km
}

func putUserLocation(w http.ResponseWriter, r *http.Request) {
	var location LocationData
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received location: %+v\n", location)

	_, err = db.Exec("INSERT INTO user_location (latitude, longitude) VALUES ($1, $2)", location.Latitude, location.Longitude)
	if err != nil {
		http.Error(w, "Error inserting user location", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(nil)
}

func putHomeLocation(w http.ResponseWriter, r *http.Request) {
	var location LocationData
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received location: %+v\n", location)

	_, err = db.Exec("INSERT INTO home_location (latitude, longitude) VALUES ($1, $2)", location.Latitude, location.Longitude)
	if err != nil {
		http.Error(w, "Error inserting user location", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(nil)
}

func checkLocation(w http.ResponseWriter, r *http.Request) {

	var check CheckerData
	err := json.NewDecoder(r.Body).Decode(&check)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var location HomeLocation
	err = db.QueryRow("SELECT latitude, longitude FROM user_location ORDER BY id ASC LIMIT 1").
		Scan(&location.Latitude, &location.Longitude)
	if err != nil {
		http.Error(w, "Error fetching user location", http.StatusInternalServerError)
		return
	}

	distance := calculateDistance(location.Latitude, location.Longitude, check.Latitude, check.Longitude)

	var home HomeLocation
	home.Latitude = check.Latitude
	home.Longitude = check.Longitude

	result := Result{
		Location: LocationData(location),
		Home:     home,
		Distance: distance,
	}

	if distance <= 0.1 {
		result.Status = "home"
	} else {
		result.Status = "outside"
	}

	var email string
	err = db.QueryRow("SELECT email FROM users ORDER BY id ASC LIMIT 1").
		Scan(&email)
	if err != nil {
		http.Error(w, "Error fetching user Authentication", http.StatusInternalServerError)
		return
	}
	if result.Status == "outside" && check.HumanSensor || check.LightSensor {
		err = util.SmtpSendMail(email, "誰かが家に侵入しました！", time.Now().Format("2006-01-02 15:04:05")+"\n誰かが家に侵入しました。")
		if err != nil {
			http.Error(w, "Error Send email", http.StatusInternalServerError)
		}
	}

	fmt.Printf("Location: %+v\n", location)
	fmt.Printf("Home: %+v\n", home)
	fmt.Printf("Distance: %.2f km\n", distance)
	fmt.Printf("Status: %s\n", result.Status)

	json.NewEncoder(w).Encode(result)
}

func generateToken(email string) (string, error) {
	expirationTime := time.Now().Add(5 * time.Hour)
	claims := &Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	return tokenString, err
}

func signup(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Email already exists", http.StatusConflict)
		return
	}

	token, err := generateToken(user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func login(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err = db.QueryRow("SELECT password_hash FROM users WHERE email = $1", user.Email).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func changeEmail(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	tokenString := r.Header.Get("Authorization")
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	claims, ok := token.Claims.(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	_, err = db.Exec("UPDATE users SET email = $1 WHERE email = $2", user.Email, claims.Email)
	if err != nil {
		http.Error(w, "Failed to update email", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
