package core

import (
	"log"
	"testing"
)

func TestParseUrl(t *testing.T) {
	url,err := Parse("did:example:123/test?service=agent&relativeRef=/credentials#hh")
	if err != nil {
		log.Println(err)
	}
	log.Println(url)
}


func TestParesQuery(t *testing.T){
	v,err :=ParseQuery("service=ijkl&relativeRef=/credentials")
	if err != nil {
		log.Println(err)
	}
	log.Println(v)
}