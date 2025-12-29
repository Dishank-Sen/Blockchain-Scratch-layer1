package types

type UsersIdentity struct{
	Peers []string	`json:"peers"`
}

type Metadata struct{
	ID string `json:"id"`
}