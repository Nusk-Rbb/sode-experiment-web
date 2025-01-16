package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/rs/cors"
)

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

var db *sql.DB

func main() {
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
	router.HandleFunc("/check-location", checkLocation).Methods("POST")

	h := cors.Default().Handler(router)

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "OK"})
	})

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

func checkLocation(w http.ResponseWriter, r *http.Request) {
	var location LocationData
	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Printf("Received location: %+v\n", location)

	var home HomeLocation
	err = db.QueryRow("SELECT latitude, longitude FROM home_location WHERE id = 1").
		Scan(&home.Latitude, &home.Longitude)
	if err != nil {
		http.Error(w, "Error fetching home location", http.StatusInternalServerError)
		return
	}

	distance := calculateDistance(location.Latitude, location.Longitude, home.Latitude, home.Longitude)

	result := Result{
		Location: location,
		Home:     home,
		Distance: distance,
	}

	if distance <= 0.1 {
		result.Status = "home"
	} else {
		result.Status = "outside"
	}

	fmt.Printf("Location: %+v\n", location)
	fmt.Printf("Home: %+v\n", home)
	fmt.Printf("Distance: %.2f km\n", distance)
	fmt.Printf("Status: %s\n", result.Status)

	json.NewEncoder(w).Encode(result)
}
