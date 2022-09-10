package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

const (
	// For server
	IP_Addr = "127.0.0.1"
	IP_Port = 7000
	// For connect to Redis
	Redis_IP     = "localhost:6379"
	Redis_Passwd = ""
	Redis_DB     = 0
	// Length of short string
	Short_Len = 5
	// Experation of key ~ 6 month
	// EXPERATION = 60*60*24*30*6
	EXPERATION = 0
)

var (
	rdb     *redis.Client
	ctx     context.Context
	letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

func check(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}

func main() {
	// Connect to Redis and create context
	ctx = context.Background()
	rdb = redis.NewClient(&redis.Options{
		Addr:     Redis_IP,
		Password: Redis_Passwd, // no password set
		DB:       Redis_DB,     // use default DB
	})
	defer rdb.Close()

	// Check if Redis work and we are connected
	pong, err := rdb.Ping(ctx).Result()
	check(err)
	log.Println("Redis ping :", pong)
	// For random
	rand.Seed(time.Now().UnixNano())

	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index)
	rtr.HandleFunc("/short", short).Methods("POST")
	rtr.HandleFunc("/info/{req}", info_short)
	rtr.HandleFunc("/status", status_full)
	rtr.HandleFunc("/{req}", req)

	http.Handle("/", rtr)
	addr := fmt.Sprintf("%s:%d", IP_Addr, IP_Port)
	log.Printf("Serv at addr: %s", addr)
	http.ListenAndServe(addr, nil)
}

// Process index page
func index(w http.ResponseWriter, r *http.Request) {
	log.Printf("idx %s", r.URL)
	// -- show start page
	t, err := template.ParseFiles("template/index.html", "template/footer.html", "template/header.html")
	check(err)
	t.ExecuteTemplate(w, "index", nil)
}

// Process request
func req(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	log.Printf("req %s", url)
	// Remove leading /
	short_name := url[1:]
	log.Printf("Short is %s", short_name)
	short_key := "short_" + short_name + "_url"
	short_num := "short_" + short_name + "_num"

	val, err := rdb.Get(ctx, short_key).Result()
	// Check if short name not present
	//if errors.Is(err, redis.Nil) {

	//}
	check(err)

	// Incrise counter
	rdb.Incr(ctx, short_num)

	log.Printf("Redirected to %s", val)
	http.Redirect(w, r, val, http.StatusMovedPermanently)
}

func randSeq() string {
	arr := make([]rune, Short_Len)
	for i := range arr {
		arr[i] = letters[rand.Intn(len(letters))]
	}
	return string(arr)
}

// Ganerate short link
func short(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	log.Println("URL : ", url)
	var (
		short_name string
		ex         int64 = 1
	)

	// Loop for checking key not present in DB
	for ex == 1 {
		short_name = randSeq()
		ex, _ = rdb.Exists(ctx, "short_"+short_name+"_url").Result()
	}

	// Svae in database
	err := rdb.Set(ctx, "short_"+short_name+"_url", url, EXPERATION).Err()
	check(err)
	err = rdb.Set(ctx, "short_"+short_name+"_num", 0, EXPERATION).Err()
	check(err)

	// Generate temlate
	t, err := template.ParseFiles("template/saved.html", "template/footer.html", "template/header.html")
	check(err)

	// Fill data for temlate
	data := struct {
		URL   string
		Short string
	}{URL: url, Short: short_name}
	t.ExecuteTemplate(w, "saved", data)
}

func status_full(w http.ResponseWriter, r *http.Request) {
	var info map[string]string
	info = make(map[string]string)

	db_info, _ := rdb.Info(ctx).Result()
	str := strings.Split(db_info, "\n")
	for _, j := range str {
		if strings.Index(j, ":") > 2 {
			sstr := strings.Split(j, ":")
			info[sstr[0]] = sstr[1]
		}
	}

	t, err := template.ParseFiles("template/status_full.html", "template/footer.html", "template/header.html")
	check(err)

	// Fill data for temlate
	t.ExecuteTemplate(w, "status_full", info)
}

func info_short(w http.ResponseWriter, r *http.Request) {
	var info map[string]string
	url := r.URL.Path
	log.Printf("req %s", url)
	// Remove leading /
	short_name := url[1:]

	info = make(map[string]string)

	db_info, _ := rdb.Info(ctx).Result()
	str := strings.Split(db_info, "\n")
	for _, j := range str {
		if strings.Index(j, ":") > 2 {
			sstr := strings.Split(j, ":")
			info[sstr[0]] = sstr[1]
		}
	}

	t, err := template.ParseFiles("template/status_full.html", "template/footer.html", "template/header.html")
	check(err)

	// Fill data for temlate
	t.ExecuteTemplate(w, "status_full", info)
}

func info_long(w http.ResponseWriter, r *http.Request) {
	var info map[string]string
	info = make(map[string]string)

	db_info, _ := rdb.GET(ctx).Result()
	str := strings.Split(db_info, "\n")
	for _, j := range str {
		if strings.Index(j, ":") > 2 {
			sstr := strings.Split(j, ":")
			info[sstr[0]] = sstr[1]
		}
	}

	t, err := template.ParseFiles("template/status_full.html", "template/footer.html", "template/header.html")
	check(err)

	// Fill data for temlate
	t.ExecuteTemplate(w, "status_full", info)
}
