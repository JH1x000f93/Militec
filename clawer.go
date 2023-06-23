package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	baseURL := "https://www.spip.net/fr_article884.html"
	signaturesPerPage := 10
	maxURLs := 1000

	file, err := os.Create("urls.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	page := 1
	totalURLs := 0
	for totalURLs < maxURLs {
		// Construir la URL de la página de paginación actual
		url := fmt.Sprintf("%s?debut_signatures=%d", baseURL, (page-1)*signaturesPerPage+1)

		// Realizar la solicitud HTTP y obtener el cuerpo de la respuesta
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		// Crear un nuevo documento goquery desde el cuerpo de la respuesta
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Encontrar los elementos HTML que contienen los nombres de los sitios web y sus URLs
		foundURLs := false
		doc.Find("ul li.box").Each(func(i int, s *goquery.Selection) {
			// Verificar si se ha alcanzado el límite máximo de URLs
			if totalURLs >= maxURLs {
				return
			}

			// Extraer la URL del sitio web
			url, _ := s.Find("a").Attr("href")

			// Eliminar la barra "/" al final de la URL, si es el último carácter
			if strings.HasSuffix(url, "/") {
				url = url[:len(url)-1]
			}

			// Escribir la URL en el archivo
			_, err := file.WriteString(url + "\n")
			if err != nil {
				log.Fatal(err)
			}

			foundURLs = true
			totalURLs++
		})

		if !foundURLs {
			// No se encontraron más URLs, salir del bucle
			break
		}

		page++
	}

	fmt.Printf("Se capturaron un máximo de %d URLs y se guardaron en el archivo urls.txt.\n", totalURLs)
}
