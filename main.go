package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
)

func connect() *sql.DB {
	connstring := "postgresql://postgres:123123@localhost/goodsinfo?sslmode=disable"
	db, err := sql.Open("postgres", connstring)
	if err != nil {
		panic(err)
	}
	return db
}
func getInserted() string {
	var data string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		data = scanner.Text()
	}
	return data
}
func insertIntoPostgres(db *sql.DB, code int, name string, price int, quantity int) {
	_, err := db.Exec("insert into goods(code, name, price, quantity) values($1, $2, $3, $4)", code, name, price, quantity)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Item added to the goods list")
	}
}
func updateName(db *sql.DB, code int, name string) {
	_, err := db.Exec("UPDATE goods SET name=$1 WHERE code=$2", name, code)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Name is changed")
	}
}
func updatePrice(db *sql.DB, code int, price int) {
	_, err := db.Exec("UPDATE goods SET price=$1 WHERE code=$2", price, code)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Price is changed")
	}
}
func updateQuantity(db *sql.DB, code int, quantity int) {
	_, err := db.Exec("UPDATE goods SET quantity=$1 WHERE code=$2", quantity, code)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Quantity is changed")
	}
}
func delete(db *sql.DB, code int) {
	_, err := db.Exec("DELETE FROM goods WHERE code=$1", code)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Record Deleted")
	}
}

func indexHandler(c *fiber.Ctx, db *sql.DB) error {
	var res string
	var goods []string
	rows, err := db.Query("SELECT * FROM goods")
	defer rows.Close()
	if err != nil {
		log.Fatalln(err)
		c.JSON("An error occured")
	}
	for rows.Next() {
		rows.Scan(&res)
		goods = append(goods, res)
	}
	return c.Render("index", fiber.Map{
		"Goods": goods,
	})
}

func main() {
	var option string
	var code int
	var name string
	var price int
	var quantity int
	db := connect()
	for {
		fmt.Println("Welcome to goods.kz")
		fmt.Println("1. Add new item")
		fmt.Println("2. Show list of all items")
		fmt.Println("3. Change item name")
		fmt.Println("4. Change item price")
		fmt.Println("5. Change item quantity")
		fmt.Println("6. Delete an item")
		fmt.Println("7. Exit")
		fmt.Println("8. Start Fiber")
		fmt.Scanln(&option)
		switch option {
		case "1":
			var code int
			var name string
			var price int
			var quantity int
			fmt.Println("Enter item's serial code (Only numbers)")
			ce := getInserted()
			code, _ = strconv.Atoi(ce)
			fmt.Println("Enter item's name")
			name = getInserted()
			fmt.Println("Enter item's price (Price in kzt)")
			pc := getInserted()
			price, _ = strconv.Atoi(pc)
			fmt.Println("Enter item's quantity")
			qt := getInserted()
			quantity, _ = strconv.Atoi(qt)
			insertIntoPostgres(db, code, name, price, quantity)
		case "2":
			rows, err := db.Query("select * from goods")
			if err != nil {
				panic(err)
			} else {
				fmt.Println("1.Serial code | 2.Name | 3.price | 4.quantity")
				fmt.Println("-------------------------------------")
				for rows.Next() {
					rows.Scan(&code, &name, &price, &quantity)
					fmt.Println("1.", code, "|", "2.", name, "|", "3.", price, "|", "4.", quantity)
				}
			}
		case "3":
			fmt.Println("Enter serial number of item that you want to change Name")
			ce := getInserted()
			code, _ = strconv.Atoi(ce)
			fmt.Println("Enter changed name")
			name = getInserted()
			updateName(db, code, name)
		case "4":
			fmt.Println("Enter serial number of item that you want to change Price")
			ce := getInserted()
			code, _ = strconv.Atoi(ce)
			fmt.Println("Enter changed price")
			pc := getInserted()
			price, _ = strconv.Atoi(pc)
			updatePrice(db, code, price)
		case "5":
			fmt.Println("Enter serial number of item that you want to change Quantity")
			ce := getInserted()
			code, _ = strconv.Atoi(ce)
			fmt.Println("Enter changed quantity")
			qt := getInserted()
			quantity, _ = strconv.Atoi(qt)
			updateQuantity(db, code, quantity)
		case "6":
			fmt.Println("Enter item's serial number that you want to delete")
			ce := getInserted()
			code, _ = strconv.Atoi(ce)
			delete(db, code)
		case "7":
			os.Exit(0)
		case "8":
			engine := html.New("./views", ".html")
			app := fiber.New(fiber.Config{
				Views: engine,
			})

			app.Get("/", func(c *fiber.Ctx) error {
				return indexHandler(c, db)
			})

			port := os.Getenv("PORT")
			if port == "" {
				port = "3000"
			}
			app.Static("/", "./public")
			log.Fatalln(app.Listen(fmt.Sprintf(":%v", port)))
		}
	}
}
