package repos

import "gopkg.in/mgo.v2/bson"

//TemplateAmbiguousSelector defines format of query from templates collection
func TemplateAmbiguousSelector(channel string, commandName string) map[string]interface{} {
	return bson.M{
		"channel": channel,
		"$or": []interface{}{
			bson.M{"commandName": commandName},
			bson.M{"aliasTo": commandName}}}
}
