# 00-bookstore

This project is a simple storage for books.

The books are stored with `id`, `title`, `author`, `price`, and `imageurl`

The library is stored as a `.json` file to allow persistence between calls of the binary
The name of said library is declared as a `const` in `bookstore.go`

## Usage

There are several methods to interact with the library

- *Add:* add book to the library
```bash
> ./bookstore add --id 1 --title test_book --author Me --price 0 --imageurl http://test.com
```
- *Get*: get book by `--id` or return `--all`
```bash
# By ID
./bookstore get --id 1

# Output
ID	Title		Author	Price	ImageURL
1	test_book	Me	0	http://test.com
```

- *Update*: Update book
```bash
> ./bookstore add --id 1 --title test_book --author Me --price 1000 --imageurl http://test.com
```
- *Delete*: Delete book by `id`
```bash
> ./bookstore delete --id 1
```
- *Help*: useful help function
```bash
> ./bookstore help
# Output
# --- Bookstore Help ---
# get --all: Return all books in library
# get --id X: Return book X from library
# add: Add book to library
# update: Update book in library
# delete: Delete book from library
```
![](img/add.png)

