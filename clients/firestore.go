package clients

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

type Movie struct {
	Name string
}

func WriteNewMovie(m Movie) bool {
	data := map[string]interface{}{
		"name": m.Name,
	}

	writeData(data, "movies")

	return true
}

func GetAllMovies() []Movie {
	results, err := readData("movies")

	if err != nil {
		return nil
	}

	var movies []Movie

	for _, r := range results {
		m := Movie{
			Name: r["name"].(string),
		}

		movies = append(movies, m)
	}

	return movies
}

func getClient() (*firestore.Client, *context.Context) {
	// Use the application default credentials
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "couple-mage"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return client, &ctx
}

func writeData(data map[string]interface{}, collection string) (*firestore.DocumentRef, error) {
	client, ctx := getClient()
	defer client.Close()

	docRef, _, err := client.Collection(collection).Add(*ctx, data)

	if err != nil {
		return nil, err
	}

	return docRef, nil
}

func readData(collection string) ([]map[string]interface{}, error) {
	client, ctx := getClient()
	defer client.Close()

	var results []map[string]interface{}

	iter := client.Collection(collection).Documents(*ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		results = append(results, doc.Data())
	}

	return results, nil
}
