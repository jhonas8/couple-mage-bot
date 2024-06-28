package clients

import (
	"context"
	"log"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
)

func WriteNewMovie(m OMDbMovie) error {
	data := map[string]interface{}{
		"title":  m.Title,
		"year":   m.Year,
		"imdbID": m.ImdbID,
		"type":   m.Type,
		"poster": m.Poster,
	}

	_, err := writeData(data, "movies")

	return err
}

func GetAllMovies() []OMDbMovie {
	results, err := readData("movies")

	if err != nil {
		return nil
	}

	var movies []OMDbMovie

	for _, r := range results {
		m := OMDbMovie{
			Title:  r["title"].(string),
			Year:   r["year"].(string),
			ImdbID: r["imdbID"].(string),
			Type:   r["type"].(string),
			Poster: r["poster"].(string),
		}

		movies = append(movies, m)
	}

	return movies
}

func getClient() (*firestore.Client, *context.Context) {
	// Use the application default credentials
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "linen-shape-420522"}
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
			log.Printf("Error reading the data from collection {%s}: %s", collection, err.Error())
			return nil, err
		}
		results = append(results, doc.Data())
	}

	return results, nil
}
