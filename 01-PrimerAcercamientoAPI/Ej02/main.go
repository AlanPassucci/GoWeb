/*
Ejercicio 2 - Manipulando el body
Vamos a crear un endpoint llamado /saludo.
Con una pequeña estructura con nombre y apellido que al pegarle deberá
responder en texto “Hola + nombre + apellido”
El endpoint deberá ser de método POST
Se deberá usar el package JSON para resolver el ejercicio
La respuesta deberá seguir esta estructura: “Hola Andrea Rivas”
La estructura deberá ser como esta:
{
	“nombre”: “Andrea”,
	“apellido”: “Rivas”
}
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

type Person struct {
	Name     string `json:"nombre"`
	Lastname string `json:"apellido"`
}

func main() {
	router := gin.Default()

	router.POST("/saludo", func(c *gin.Context) {
		person := Person{}
		jsonData, err := c.GetRawData()
		if err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}

		if err := json.Unmarshal(jsonData, &person); err != nil {
			c.JSON(400, gin.H{"message": err})
			return
		}

		if person.Name == "" || person.Lastname == "" {
			c.JSON(400, gin.H{"message": "parameters must be non-zero values"})
			return
		}

		c.String(200, fmt.Sprintf("Hola %s %s", person.Name, person.Lastname))
	})

	router.Run()
}
