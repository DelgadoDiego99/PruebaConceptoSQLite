// API para registrar datos de una universidad (Solo hay una tabla para estudiantes)

package main

import (
	"Web/Proyecto1/models" // En vez de llamar la carpeta que tiene a los modelos, es mejor llamar a un archivo que contenga todos los handlers para la API
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	err := models.ConnectDatabase()
	checkErr(err)

	r := gin.Default()

	// Colocar las rutas en otro archivo y llenar el grupo a partir de una función
	v1 := r.Group("/")
	{
		v1.GET("estudiante", getEstudiantes)
		v1.GET("estudiante/:codigo", getEstudianteByID)
		v1.POST("estudiante", addEstudiante)
		v1.PUT("estudiante/:codigo", updateEstudiante)
		v1.DELETE("estudiante/:codigo", deleteEstudiante)
		v1.OPTIONS("estudiante", options)
	}

	r.Run()
}

// Colocar los handlers en otro archivo para mejor manejo
func getEstudiantes(c *gin.Context) {
	estudiantes, err := models.GetEstudiantes(100)
	checkErr(err)

	if estudiantes == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": estudiantes})
	}
}

func getEstudianteByID(c *gin.Context) {
	codigo := c.Param("codigo")
	estudiantes, err := models.GetEstudianteById(codigo)
	checkErr(err)

	if estudiantes.Nombre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No Records Found"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"data": estudiantes})
	}
}

func addEstudiante(c *gin.Context) {
	var json models.Estudiante

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	success, err := models.AddEstudiante(json)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func updateEstudiante(c *gin.Context) {
	var json models.Estudiante

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	estudianteCodigo, err := strconv.Atoi(c.Param("codigo"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Codigo invalido"})
	}

	success, err := models.UpdateEstudiante(json, estudianteCodigo)

	if success {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}
}

func deleteEstudiante(c *gin.Context) {
	codigo := c.Param("codigo")
	c.JSON(http.StatusOK, gin.H{"message": "deleteEstudiante " + codigo + " Called"})
}

func options(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "options Called"})
}

// Buscar si hay un error para detener la aplicación
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
