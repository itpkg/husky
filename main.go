package main

import (
	"log"

	"github.com/itpkg/husky/web"
	_ "github.com/itpkg/husky/web/blog"
	_ "github.com/itpkg/husky/web/books"
	_ "github.com/itpkg/husky/web/cms"
	_ "github.com/itpkg/husky/web/core"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	if err := web.Run(); err != nil {
		log.Fatal(err)
	}
}
