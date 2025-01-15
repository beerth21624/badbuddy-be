package models

type CourtWithVenue struct {
	Court
	VenueName     string `db:"venue_name"`
	VenueLocation string `db:"venue_location"`
	VenueStatus   string `db:"venue_status"`
}
