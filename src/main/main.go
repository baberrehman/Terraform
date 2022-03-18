package main

import (
	"Terraform/src/terraform"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	_, output := terraform.ToJSON("/home/baber/GolandProjects/scripts/aws.tf")
	fmt.Println(output)
}

func getAlbums(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, "Hello")
}
