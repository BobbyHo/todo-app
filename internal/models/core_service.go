package models

type Database interface {
	Close() error
}

type Store struct {
	TodoUserStore TodoUserStore

type Service struct {
	Database      Database
	Store         *Store
}