package utils

import "go.mongodb.org/mongo-driver/v2/bson"

func ObjectIDSliceToStringSlice(ids []bson.ObjectID) []string {
	res := make([]string, len(ids))
	for i, oid := range ids {
		res[i] = oid.Hex()
	}
	return res
}

func StringSliceToObjectIDSlice(ids []string) ([]bson.ObjectID, error) {
	res := make([]bson.ObjectID, len(ids))
	for i, s := range ids {
		oid, err := bson.ObjectIDFromHex(s)
		if err != nil {
			return nil, err
		}
		res[i] = oid
	}
	return res, nil
}
