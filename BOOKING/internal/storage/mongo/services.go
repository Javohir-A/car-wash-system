package mongostorage

import (
	"context"
	"errors"
	"time"

	"booking-service/genproto/services"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServicesMongoStorageImpl struct {
	coll *mongo.Collection
}

func NewServicesMongoStorage(coll *mongo.Collection) *ServicesMongoStorageImpl {
	return &ServicesMongoStorageImpl{
		coll: coll,
	}
}

func (s *ServicesMongoStorageImpl) CreateService(ctx context.Context, newService *services.NewService) (*services.CreateResp, error) {
	creationTime := time.Now()
	document := bson.M{
		"name":        newService.Name,
		"description": newService.Description,
		"price":       newService.Price,
		"duration":    newService.Duration,
		"created_at":  creationTime,
		"updated_at":  creationTime,
	}

	result, err := s.coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return &services.CreateResp{
		Id:        result.InsertedID.(primitive.ObjectID).Hex(),
		CreatedAt: creationTime.String(),
	}, nil
}

func (s *ServicesMongoStorageImpl) UpdateService(ctx context.Context, newData *services.NewData) (*services.UpdateResp, error) {
	updateTime := time.Now()
	update := bson.M{
		"$set": bson.M{
			"name":        newData.Name,
			"description": newData.Description,
			"price":       newData.Price,
			"duration":    newData.Duration,
			"updated_at":  updateTime,
		},
	}

	objectID, err := primitive.ObjectIDFromHex(newData.Id)
	if err != nil {
		return nil, err
	}

	_, err = s.coll.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return nil, err
	}

	return &services.UpdateResp{
		Id:        newData.Id,
		UpdatedAt: updateTime.String(),
	}, nil
}

func (s *ServicesMongoStorageImpl) DeleteService(ctx context.Context, id *services.ID) (*services.Void, error) {
	objectID, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return nil, err
	}

	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return nil, err
	}

	return &services.Void{}, nil
}

func (s *ServicesMongoStorageImpl) ListServices(ctx context.Context, pagination *services.Pagination) (*services.ServicesList, error) {
	skip := int64((pagination.Page - 1) * pagination.Limit)
	limit := int64(pagination.Limit)

	options := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}

	cursor, err := s.coll.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var serviceList []*services.Service
	for cursor.Next(ctx) {
		var service services.Service
		if err := cursor.Decode(&service); err != nil {
			return nil, err
		}
		serviceList = append(serviceList, &service)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &services.ServicesList{
		Services: serviceList,
		Page:     pagination.Page,
		Limit:    pagination.Limit,
	}, nil
}

func (s *ServicesMongoStorageImpl) SearchServices(ctx context.Context, filter *services.Filter) (*services.SearchResp, error) {
	mongoFilter := bson.M{}

	if filter.Name != "" {
		mongoFilter["name"] = bson.M{"$regex": filter.Name, "$options": "i"}
	}
	if filter.Price > 0 {
		mongoFilter["price"] = bson.M{"$gte": filter.Price}
	}
	if filter.Duration > 0 {
		mongoFilter["duration"] = bson.M{"$gte": filter.Duration}
	}
	if filter.CreatedAt != "" {
		parsedTime, err := time.Parse(time.RFC3339, filter.CreatedAt)
		if err != nil {
			return nil, err
		}
		mongoFilter["created_at"] = bson.M{"$gte": parsedTime}
	}

	cursor, err := s.coll.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var serviceList []*services.Service
	for cursor.Next(ctx) {
		var service services.Service
		if err := cursor.Decode(&service); err != nil {
			return nil, err
		}
		serviceList = append(serviceList, &service)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return &services.SearchResp{
		Services: serviceList,
	}, nil
}

func (s *ServicesMongoStorageImpl) GetServiceByID(ctx context.Context, id *services.ID) (*services.Service, error) {
	var service services.Service
	objID, err := primitive.ObjectIDFromHex(id.Id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}

	err = s.coll.FindOne(ctx, filter).Decode(&service)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil // Return nil if no document was found
		}
		return nil, err
	}
	return &service, nil
}
