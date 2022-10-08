package uuid

import "go.mongodb.org/mongo-driver/bson/primitive"
func New() string {
	return primitive.NewObjectID().Hex()
}
