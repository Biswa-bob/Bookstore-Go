package models

func (b *Book) ToResponse() BookResponse {
	return BookResponse{
		ID:            b.ID,
		UserID:        b.UserID,
		Title:         b.Title,
		Author:        b.Author,
		ISBN:          b.ISBN,
		Category:      b.Category,
		PublishedDate: b.PublishedDate,
		Publisher:     b.Publisher,
		Pages:         b.Pages,
		Language:      b.Language,
		PriceCents:    b.PriceCents,
		Stock:         b.Stock,
		InStock:       b.Stock > 0, // computed
		Description:   b.Description,
		Tags:          b.Tags,
		Metadata:      b.Metadata,
		CreatedAt:     b.CreatedAt,
		UpdatedAt:     b.UpdatedAt,
	}
}
