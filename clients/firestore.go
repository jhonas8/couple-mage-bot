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

func SaveIdsForMovieMessages(chatId int64, savedIds []int, movieTitle string) error {
	data := map[string]interface{}{
		"chatId":     chatId,
		"savedIds":   savedIds,
		"movieTitle": movieTitle,
	}

	_, err := writeData(data, "movieIds")

	return err
}

func GetIdsForMovieMessages(chatId int64, movieTitle string) ([]int, []string, error) {
	results, err := queryData("movieIds", "movieTitle", movieTitle)
	if err != nil {
		return nil, nil, err
	}

	var savedIds []int
	var documentIds []string

	for _, r := range results {
		if r["chatId"] == chatId {
			log.Printf("Found result with format %v", r)
			savedIds = r["savedIds"].([]int)
			documentIds = append(documentIds, r["id"].(string))
		}
	}

	return savedIds, documentIds, nil
}

func DeleteSavedIds(chatId int64, documentId string) error {
	return deleteDocument("movieIds", documentId)
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

func deleteDocument(collection string, documentId string) error {
	client, ctx := getClient()
	defer client.Close()

	_, err := client.Collection(collection).Doc(documentId).Delete(*ctx)

	return err
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

func queryData(collection string, field string, value string) ([]map[string]interface{}, error) {
	client, ctx := getClient()
	defer client.Close()

	iter := client.Collection(collection).Where(field, "==", value).Documents(*ctx)

	var results []map[string]interface{}

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
