package servo

import (
	"log"
	"testing"

	"github.com/mitghi/x/structs"
)

func TestINFOGENERAL(t *testing.T) {
	var (
		info_httpserver string
		info_server     string
	)
	info_httpserver = structs.CompileStructInfo(HTTPServer{})
	info_server = structs.CompileStructInfo(Server{})
	log.Println("info_httpserver:", info_httpserver)
	log.Println("info_server:", info_server)
}
