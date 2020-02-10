package server

import "gopkg.in/mgo.v2"
// import "log"

func (s *Server)InitDatabase() {
	if s.db == nil {
		s.db, _ = mgo.Dial("mongo:27017")
		// if err != nil {
		// 	log.Println(err.Error())
		// }
	}
}

func (s *Server)AccessDataStore() *mgo.Session {
	if s.db == nil {
		s.db, _ = mgo.Dial("mongo:27017")
	}
	return s.db.Copy()
}