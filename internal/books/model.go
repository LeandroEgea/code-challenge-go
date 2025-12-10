package books

type Book struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Author    string `json:"author"`
	UnitsSold uint   `json:"units_sold"`
	Price     uint   `json:"price"`
}

type Metrics struct {
	MeanUnitsSold        uint   `json:"mean_units_sold"`
	CheapestBook         string `json:"cheapest_book"`
	BooksWrittenByAuthor uint   `json:"books_written_by_author"`
}
