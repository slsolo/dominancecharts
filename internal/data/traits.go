package data

type TraitData struct {
	Name           string `redis:"name"`
	Position       int    `redis:"position"`
	FirstRelease   bool   `redis:"first_release"`
	LimitedEdition bool   `redis:"limited_edition"`
	Retired        bool   `redis:"retired"`
}
