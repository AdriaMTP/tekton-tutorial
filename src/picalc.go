package main

import (
    "database/sql"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    _ "github.com/lib/pq"
)

const (
    host     = "10.32.0.1"
    port     = 31986
    user     = "postgresadmin"
    password = "admin123"
    dbname   = "postgresdb"
)

// Calculate pi using Gregory-Leibniz series:   (4/1) - (4/3) + (4/5) - (4/7) + (4/9) - (4/11) + (4/13) - (4/15) ...
func calculatePi(iterations int) float64 {
    var result float64 = 0.0
    var sign float64 = 1.0
    var denominator float64 = 1.0
    for i := 0; i < iterations; i++ {
        result = result + (sign * 4/denominator)
        denominator = denominator + 2
        sign = -sign
    }
    return result
}

// Purchases a 'product' equal to database id.
func purchaseProduct(product int) string {
  psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)
  db, err := sql.Open("postgres", psqlInfo)
  if err != nil {
    panic(err)
  }
  defer db.Close()

  err = db.Ping()
  if err != nil {
    panic(err)
  }

  return "Successfully connected!"
}

func handlerPicalc(w http.ResponseWriter, r *http.Request) {
    log.Print("Pi calculator received a request.")
    iterations, err := strconv.Atoi(r.URL.Query()["iterations"][0])
    if err != nil {
        fmt.Fprintf(w, "iterations parameter not valid\n")
        return
    }
    fmt.Fprintf(w, "%.10f\n", calculatePi(iterations))
}

func handlerPurchase(w http.ResponseWriter, r *http.Request) {
    log.Print("Purchase function received a request.")
    product, err := strconv.Atoi(r.URL.Query()["product"][0])
    if err != nil {
        fmt.Fprintf(w, "product parameter not valid\n")
        return
    }
    fmt.Fprintf(w, "%s\n", purchaseProduct(product))
}

func main() {
    log.Print("App started.")

    http.HandleFunc("/picalc", handlerPicalc)
    log.Print("Pi calculator is listening on '/picalc'.")

    http.HandleFunc("/purchase", handlerPurchase)
    log.Print("Product purchaser is listening on '/purchase'.")

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
