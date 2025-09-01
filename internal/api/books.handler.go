package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Biswa-bob/bookstore/internal/models"
	"github.com/Biswa-bob/bookstore/internal/store"
	"github.com/Biswa-bob/bookstore/internal/utils"
)

type BooksHandler struct {
	booksStore store.BooksStore
	logger     *log.Logger
}

func NewBooksHandler(booksStore store.BooksStore, logger *log.Logger) *BooksHandler {
	return &BooksHandler{
		booksStore: booksStore,
		logger:     logger,
	}
}

func (bh *BooksHandler) HandleCreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		bh.logger.Printf("ERROR: decodingCreateBook: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request sent"})
		return
	}
	// currentUser := middleware.GetUser(r)
	// if currentUser == nil || currentUser == store.AnonymousUser {
	// 	wh.logger.Printf("ERROR: decodingCreateBook: %v", err)
	// 	utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "you must be logged in"})
	// 	return
	// }
	// workout.UserID = currentUser.ID
	createdBook, err := bh.booksStore.CreateBook(&book)
	if err != nil {
		bh.logger.Printf("ERROR: createBook: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create book"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"book": createdBook})
}

func (bh *BooksHandler) HandleGetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := bh.booksStore.GetBooks()
	if err != nil {
		bh.logger.Printf("ERROR: getBooks: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"books": books})
}
