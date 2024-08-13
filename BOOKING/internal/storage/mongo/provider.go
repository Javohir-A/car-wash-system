package mongostorage

import (
	"booking-service/genproto/providers"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProvidersMongoStorageImpl struct {
	coll *mongo.Collection
}

func NewProviderMongoStorage(coll *mongo.Collection) *ProvidersMongoStorageImpl {
	return &ProvidersMongoStorageImpl{
		coll: coll,
	}
}

func (p *ProvidersMongoStorageImpl) CreateProvider(ctx context.Context, pro *providers.NewProvider) (*providers.CreateResp, error) {
	// Add the current timestamp to the document
	creationTime := time.Now()
	document := bson.M{
		"user_id":        pro.UserId,
		"company_name":   pro.CompanyName,
		"description":    pro.Description,
		"services":       pro.Services,
		"availability":   pro.Availability,
		"average_rating": pro.AverageRating,
		"location":       pro.Location,
		"created_at":     creationTime,
		"updated_at":     time.Now(),
	}

	result, err := p.coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return &providers.CreateResp{
		Id:        result.InsertedID.(primitive.ObjectID).Hex(),
		CreatedAt: creationTime.String(),
	}, nil
}

func (p *ProvidersMongoStorageImpl) SearchProviders(ctx context.Context, filter *providers.Filter) (*providers.SearchResp, error) {
	mongoFilter := bson.M{}

	if filter.Name != "" {
		mongoFilter["company_name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}
	if filter.AverageRating > 0 {
		mongoFilter["average_rating"] = bson.M{"$gte": filter.AverageRating}
	}
	if filter.CreatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, filter.CreatedAt)
		if err != nil {
			return nil, err
		}
		mongoFilter["created_at"] = bson.M{"$gte": parsedTime}
	}

	cursor, err := p.coll.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var providersList []*providers.Provider
	for cursor.Next(ctx) {
		var provider providers.Provider
		if err := cursor.Decode(&provider); err != nil {
			return nil, err
		}
		providersList = append(providersList, &provider)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &providers.SearchResp{
		Providers: providersList,
	}, nil
}
