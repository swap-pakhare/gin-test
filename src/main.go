package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"sort"
	"strconv"
)

type orderDetails struct {
	Id int
	CustomerId int
	RestaurantId int
	Amount float64
	Status string
	DEId int
	Cart string
	PaymentMode string
}

type kv struct {
	Key   int
	Value int
}
type DataHandler struct {
	jsonFile string
	fr *os.File
	dec *json.Decoder
	rests map[int]int
	sortedRests []kv
	orders []orderDetails
}

func checkError(err error)  {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (dh *DataHandler) init(jsonFilePath string) {

	dh.jsonFile = jsonFilePath

	var err error
	dh.fr, err = os.Open(jsonFilePath)
	checkError(err)

	dh.dec = json.NewDecoder(dh.fr)

	_, err = dh.dec.Token()
	checkError(err)

	dh.rests = make(map[int]int)


}

func (dh *DataHandler) close() {
	dh.fr.Close()
}

func (dh *DataHandler) processOrders()  {
	for dh.dec.More() {
		var tmp orderDetails
		err := dh.dec.Decode(&tmp)
		checkError(err)
		dh.orders = append(dh.orders, tmp)
	}
}

func (dh *DataHandler) processRestaurants() {
	for _, order := range dh.orders {
		dh.rests[order.RestaurantId]++
	}
}

func (dh *DataHandler) sortRestaurants()  {
	var ss []kv
	for k, v := range dh.rests {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	dh.sortedRests = ss
}

func (dh *DataHandler) addOrder(order orderDetails) {
	f, err := os.OpenFile(dh.jsonFile, os.O_RDWR, os.ModePerm)
	defer f.Close()
	checkError(err)

	orderJson, err := json.Marshal(order)
	checkError(err)

	orderString := string(orderJson)
	orderString = "," + orderString

	off := int64(1)
	stat, err := os.Stat(dh.jsonFile)
	fmt.Println("Size : ", stat.Size())
	start := stat.Size() - off

	tmp := []byte(orderString)
	_, err = f.WriteAt(tmp, start)
	checkError(err)

	str := []byte("]")
	_, err = f.WriteAt(str, start + int64(len(orderString)))
	checkError(err)

}

func topRestaurants(c *gin.Context) {
	str := c.Param("num")
	num, err := strconv.Atoi(str)
	checkError(err)

	c.JSON(200, dh.sortedRests[:num])
}

func topOrders(c *gin.Context) {
	str := c.Param("num")
	num, err := strconv.Atoi(str)
	checkError(err)

	c.JSON(200, dh.orders[:num])
}

func allRestaurants(c *gin.Context)  {
	c.JSON(200, dh.sortedRests)
}

func allOrders(c *gin.Context)  {
	c.JSON(200, dh.orders)
}

func createOrder(c *gin.Context)  {
	var order orderDetails
	err := c.BindJSON(&order)
	checkError(err)
	dh.addOrder(order)
	c.String(200, "OK Added!")
}

var dh DataHandler

func main() {

	jsonFilePath := "outputs.json"

	dh.init(jsonFilePath)
	defer dh.close()

	dh.processOrders()

	dh.processRestaurants()

	dh.sortRestaurants()

	router := gin.Default()

	router.GET("/orders", allOrders)

	router.GET("/orders/:num", topOrders)

	router.GET("/rest", allRestaurants)

	router.GET("/rest/:num", topRestaurants)

	router.POST("/createorder", createOrder)

	router.Run(":8080")
}