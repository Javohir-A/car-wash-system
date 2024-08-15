package mongostorage

import (
	"context"
	"time"

	"booking-service/genproto/reviews"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewsStorage struct {
	collection *mongo.Collection
}

func NewReviewsStorage(coll *mongo.Collection) *ReviewsStorage {
	return &ReviewsStorage{
		collection: coll,
	}
}
func (s *ReviewsStorage) CreateReview(ctx context.Context, req *reviews.NewReview) (*reviews.CreateResp, error) {
	review := bson.M{
		"booking_id":  req.BookingId,
		"user_id":     req.UserId,
		"provider_id": req.ProviderId,
		"rating":      req.Rating,
		"comment":     req.Comment,
		"created_at":  time.Now().Format(time.RFC3339),
		"updated_at":  time.Now().Format(time.RFC3339),
	}

	result, err := s.collection.InsertOne(ctx, review)
	if err != nil {
		return nil, err
	}

	return &reviews.CreateResp{
		Id:        result.InsertedID.(primitive.ObjectID).Hex(),
		CreatedAt: review["created_at"].(string),
	}, nil
}

func (s *ReviewsStorage) UpdateReview(ctx context.Context, req *reviews.NewData) (*reviews.UpdateResp, error) {
	objID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	update := bson.M{
		"$set": bson.M{
			"rating":     req.GetRating(),
			"comment":    req.GetComment(),
			"updated_at": time.Now().Format(time.RFC3339),
		},
	}

	_, err = s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return &reviews.UpdateResp{
		Id:        req.Id,
		UpdatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

func (s *ReviewsStorage) DeleteReview(ctx context.Context, req *reviews.ID) (*reviews.Void, error) {
	objID, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}

	_, err = s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &reviews.Void{}, nil
}

func (s *ReviewsStorage) ListReviews(ctx context.Context, req *reviews.Pagination) (*reviews.ReviewsList, error) {
	findOptions := options.Find()
	findOptions.SetLimit(int64(req.Limit))
	findOptions.SetSkip(int64((req.Page - 1) * req.Limit))

	cursor, err := s.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reviewsList []*reviews.Review
	for cursor.Next(ctx) {
		var review reviews.Review
		if err := cursor.Decode(&review); err != nil {
			return nil, err
		}
		reviewsList = append(reviewsList, &review)
	}

	return &reviews.ReviewsList{
		Reviews: reviewsList,
		Page:    req.Page,
		Limit:   req.Limit,
	}, nil
}
