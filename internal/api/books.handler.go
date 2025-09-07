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

func (bh *BooksHandler) HandleGetBookById(w http.ResponseWriter, r *http.Request) {
	bookID, err := utils.ReadIDParam(r)
	if err != nil {
		bh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
	}
	book, err := bh.booksStore.GetBookById(bookID)
	if err != nil {
		bh.logger.Printf("ERROR: getBooks: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"book": book})
}

func (bh *BooksHandler) HandleUpdateBookById(w http.ResponseWriter, r *http.Request) {
	bookID, err := utils.ReadIDParam(r)
	if err != nil {
		bh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
	}
	var dto store.UpdateBookDTO
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&dto); err != nil {
		bh.logger.Printf("ERROR: decode body: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid JSON body"})
		return
	}

	// Ensure at least one field is provided
	if dto.Title == nil && dto.Author == nil && dto.PriceCents == nil && dto.PublishedAt == nil && dto.InStock == nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "no fields to update"})
		return
	}

	// (Optional) basic validation examples
	if dto.PriceCents != nil && *dto.PriceCents < 0 {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "price must be >= 0"})
		return
	}
	updated, err := bh.booksStore.UpdateBookFieldsByID(bookID, dto)
	if err != nil {
		bh.logger.Printf("ERROR: getBooks: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"updated": true, "book": updated})
}

func (bh *BooksHandler) HandleDeleteBookById(w http.ResponseWriter, r *http.Request) {
	bookID, err := utils.ReadIDParam(r)
	if err != nil {
		bh.logger.Printf("ERROR: readIDParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
	}
	err = bh.booksStore.DeleteBookById(bookID)
	if err != nil {
		bh.logger.Printf("ERROR: getBooks: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"deleted": true, "book_id": bookID})
}
